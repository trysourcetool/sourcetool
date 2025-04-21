package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = pingPeriod * 2

	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 1024 // 512KB

	// Maximum time to wait for connection recovery.
	maxRecoveryWait = 6 * time.Hour
)

func (s *Server) handleInitializeClient(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetInitializeClient()
	if in == nil {
		return errors.New("invalid message")
	}

	pageID, err := uuid.FromString(in.PageId)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	page, err := s.db.Page().Get(ctx, database.PageByID(pageID))
	if err != nil {
		return err
	}

	currentOrg := internal.CurrentOrganization(ctx)
	if currentOrg.ID != page.OrganizationID {
		return errdefs.ErrPermissionDenied(errors.New("organization mismatch"))
	}

	apiKey, err := s.db.APIKey().Get(ctx, database.APIKeyByID(page.APIKeyID))
	if err != nil {
		return err
	}

	hostInstances, err := s.db.HostInstance().List(ctx, database.HostInstanceByAPIKeyID(apiKey.ID))
	if err != nil {
		return err
	}

	// Try to find an online host that responds to ping
	var onlineHostInstance *core.HostInstance
	for _, hostInstance := range hostInstances {
		if hostInstance.Status == core.HostInstanceStatusOnline {
			if err := s.wsManager.PingConnectedHost(hostInstance.ID); err != nil {
				// Update host status to offline if ping fails
				hostInstance.Status = core.HostInstanceStatusOffline
				if err := s.db.HostInstance().Update(ctx, hostInstance); err != nil {
					logger.Logger.Sugar().Errorf("Failed to update host status: %v", err)
				}
				continue
			}

			onlineHostInstance = hostInstance
			break
		}
	}

	// If no online host found, try hosts that might be unreachable
	if onlineHostInstance == nil {
		for _, hostInstance := range hostInstances {
			if hostInstance.Status == core.HostInstanceStatusUnreachable {
				if err := s.wsManager.PingConnectedHost(hostInstance.ID); err == nil {
					// Host is actually reachable, update its status
					hostInstance.Status = core.HostInstanceStatusOnline
					if err := s.db.HostInstance().Update(ctx, hostInstance); err != nil {
						logger.Logger.Sugar().Errorf("Failed to update host status: %v", err)
						continue
					}
					onlineHostInstance = hostInstance
					break
				}
			}
		}
	}

	if onlineHostInstance == nil {
		return errdefs.ErrHostInstanceStatusNotOnline(errors.New("no available host instances"))
	}

	currentUser := internal.CurrentUser(ctx)

	var sess *core.Session
	var sessionExists bool
	if internal.SafeValue(in.SessionId) != "" {
		sessionID, err := uuid.FromString(internal.SafeValue(in.SessionId))
		if err != nil {
			return errdefs.ErrSessionNotFound(err)
		}

		sess, err = s.db.Session().Get(ctx, database.SessionByID(sessionID))
		if err != nil {
			return err
		}
		sessionExists = true
	} else {
		sess = &core.Session{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: page.OrganizationID,
			APIKeyID:       page.APIKeyID,
			HostInstanceID: onlineHostInstance.ID,
			UserID:         currentUser.ID,
		}
		sessionExists = false
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if !sessionExists {
			if err := tx.Session().Create(ctx, sess); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if err := s.sendWebSocketMessage(conn, &websocketv1.Message{
		Id: msg.Id,
		Type: &websocketv1.Message_InitializeClientCompleted{
			InitializeClientCompleted: &websocketv1.InitializeClientCompleted{
				SessionId: sess.ID.String(),
			},
		},
	}); err != nil {
		return err
	}

	s.wsManager.SetConnectedClient(sess, conn)

	if err := s.wsManager.SendToHost(ctx, onlineHostInstance.ID, &websocketv1.Message{
		Id: uuid.Must(uuid.NewV4()).String(),
		Type: &websocketv1.Message_InitializeClient{
			InitializeClient: &websocketv1.InitializeClient{
				SessionId: internal.NilValue(sess.ID.String()),
				PageId:    page.ID.String(),
			},
		},
	}); err != nil {
		s.db.Session().Delete(ctx, sess)
		s.wsManager.DisconnectClient(sess.ID)
		logger.Logger.Sugar().Errorf("Failed to send initialize client message to host: %v", err)
		return err
	}

	return nil
}

func (s *Server) handleRenderWidget(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetRenderWidget()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return err
	}

	_, err = s.db.Session().Get(ctx, database.SessionByID(sessionID))
	if err != nil {
		return err
	}

	if err := s.wsManager.SendToClient(ctx, sessionID, msg); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send render widget message to client: %v", err)
		return err
	}

	return nil
}

func (s *Server) handleRerunPage(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetRerunPage()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return err
	}

	sess, err := s.db.Session().Get(ctx, database.SessionByID(sessionID))
	if err != nil {
		return err
	}

	pageID, err := uuid.FromString(in.PageId)
	if err != nil {
		return err
	}

	page, err := s.db.Page().Get(ctx, database.PageByID(pageID), database.PageBySessionID(sess.ID))
	if err != nil {
		return err
	}

	if err := s.wsManager.SendToHost(ctx, sess.HostInstanceID, &websocketv1.Message{
		Id: msg.Id,
		Type: &websocketv1.Message_RerunPage{
			RerunPage: &websocketv1.RerunPage{
				SessionId: sess.ID.String(),
				PageId:    page.ID.String(),
				States:    in.States,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) handleCloseSession(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetCloseSession()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return errdefs.ErrAPIKeyNotFound(err)
	}

	sess, err := s.db.Session().Get(ctx, database.SessionByID(sessionID))
	if err != nil {
		return err
	}

	_, err = s.db.Page().Get(ctx, database.PageByAPIKeyID(sess.APIKeyID), database.PageBySessionID(sess.ID))
	if err != nil {
		return err
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.Session().Delete(ctx, sess); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := s.wsManager.SendToHost(ctx, sess.HostInstanceID, &websocketv1.Message{
		Id: uuid.Must(uuid.NewV4()).String(),
		Type: &websocketv1.Message_CloseSession{
			CloseSession: &websocketv1.CloseSession{
				SessionId: sess.ID.String(),
			},
		},
	}); err != nil {
		logger.Logger.Sugar().Warnf("Failed to send close session message to host %s for session %s: %v", sess.HostInstanceID, sess.ID, err)
	}

	s.wsManager.DisconnectClient(sess.ID)

	return nil
}

func (s *Server) handleScriptFinished(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetScriptFinished()
	if in == nil {
		return errors.New("invalid message")
	}

	logger.Logger.Sugar().Debug("Payload: ", in)

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	_, err = s.db.Session().Get(ctx, database.SessionByID(sessionID))
	if err != nil {
		return err
	}

	if err := s.wsManager.SendToClient(ctx, sessionID, msg); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send script finished message to client: %v", err)
	}

	return nil
}

func (s *Server) handleException(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetException()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	_, err = s.db.Session().Get(ctx, database.SessionByID(sessionID))
	if err != nil {
		return err
	}

	if err := s.wsManager.SendToClient(ctx, sessionID, msg); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send exception message to client: %v", err)
	}

	return nil
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Logger.Sugar().Errorf("Failed to upgrade connection: %v", err)
		return
	}

	conn.SetReadLimit(maxMessageSize)
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	ctx := internal.NewBackgroundContext(r.Context())

	done := make(chan struct{})
	defer func() {
		logger.Logger.Info("Closing connection")
		close(done)
		conn.Close()
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			logger.Logger.Sugar().Debugf("Connection closed: %v", err)
			break
		}

		var msg websocketv1.Message
		if err := proto.Unmarshal(data, &msg); err != nil {
			logger.Logger.Sugar().Errorf("Failed to unmarshal message: %v", err)
			break
		}

		switch msg.Type.(type) {
		case *websocketv1.Message_InitializeHost:
			instanceID := r.Header.Get("X-Instance-Id")
			hostInstance, err := s.handleInitializeHost(ctx, conn, instanceID, &msg)
			if err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}

			defer func() {
				if err := s.updateHostInstanceStatus(ctx, hostInstance.ID, core.HostInstanceStatusOffline); err != nil {
					logger.Logger.Sugar().Errorf("Failed to update host instance status offline: %v", err)
				}
				s.wsManager.DisconnectHost(hostInstance.ID)
			}()

			go s.pingPongHostInstanceLoop(ctx, conn, done, hostInstance)
		case *websocketv1.Message_InitializeClient:
			if err := s.handleInitializeClient(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_RenderWidget:
			if err := s.handleRenderWidget(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_RerunPage:
			if err := s.handleRerunPage(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_CloseSession:
			if err := s.handleCloseSession(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_ScriptFinished:
			if err := s.handleScriptFinished(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_Exception:
			if err := s.handleException(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		default:
			logger.Logger.Sugar().Errorf("Unknown method: %s", msg.Type)
			continue
		}
	}
}

func (s *Server) updateHostInstanceStatus(ctx context.Context, hostInstanceID uuid.UUID, status core.HostInstanceStatus) error {
	host, err := s.db.HostInstance().Get(ctx, database.HostInstanceByID(hostInstanceID))
	if err != nil {
		return err
	}

	host.Status = status

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.HostInstance().Update(ctx, host); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) pingPongHostInstanceLoop(ctx context.Context, conn *websocket.Conn, done chan struct{}, hostInstance *core.HostInstance) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		logger.Logger.Info("Stopped ping ticker")
		ticker.Stop()
	}()

	var firstFailureTime *time.Time
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
				logger.Logger.Sugar().Errorf("Failed to write ping message: %v", err)
				// Record the first failure time if not already set
				if firstFailureTime == nil {
					now := time.Now()
					firstFailureTime = &now
					logger.Logger.Sugar().Infof("Recording first ping failure time: %v", now)
				}

				// Check if we've exceeded the maximum recovery wait time
				if time.Since(*firstFailureTime) > maxRecoveryWait {
					logger.Logger.Sugar().Infof("Connection unrecoverable after %v", maxRecoveryWait)
					return
				}
				if hostInstance.Status != core.HostInstanceStatusUnreachable {
					if err := s.updateHostInstanceStatus(ctx, hostInstance.ID, core.HostInstanceStatusUnreachable); err != nil {
						logger.Logger.Sugar().Errorf("Failed to update host instance status unreachable: %v", err)
					}
				}
				continue
			}
			// Reset failure time if ping succeeds
			if firstFailureTime != nil {
				logger.Logger.Info("Connection recovered, resetting failure time")
				firstFailureTime = nil
			}
			if hostInstance.Status != core.HostInstanceStatusOnline {
				if err := s.updateHostInstanceStatus(ctx, hostInstance.ID, core.HostInstanceStatusOnline); err != nil {
					logger.Logger.Sugar().Errorf("Failed to update host instance status online: %v", err)
				}
			}
		}
	}
}

func (s *Server) handleInitializeHostBase(ctx context.Context, conn *websocket.Conn, instanceID string, msg *websocketv1.Message) (*core.HostInstance, bool, *core.APIKey, []*core.Page, []*core.Page, []*core.Page, error) {
	in := msg.GetInitializeHost()
	if in == nil {
		return nil, false, nil, nil, nil, nil, errors.New("invalid message")
	}

	apikey, err := s.db.APIKey().Get(ctx, database.APIKeyByKey(in.ApiKey))
	if err != nil {
		return nil, false, nil, nil, nil, nil, err
	}

	hostInstanceID, err := uuid.FromString(instanceID)
	if err != nil {
		return nil, false, nil, nil, nil, nil, errdefs.ErrInvalidArgument(err)
	}

	hostInstance, err := s.db.HostInstance().Get(ctx, database.HostInstanceByID(hostInstanceID))
	if err != nil && !errdefs.IsHostInstanceNotFound(err) {
		return nil, false, nil, nil, nil, nil, err
	}

	hostExists := hostInstance != nil

	if !hostExists {
		hostInstance = &core.HostInstance{
			ID:             hostInstanceID,
			OrganizationID: apikey.OrganizationID,
			APIKeyID:       apikey.ID,
		}
	}

	hostInstance.SDKName = in.SdkName
	hostInstance.SDKVersion = in.SdkVersion
	hostInstance.Status = core.HostInstanceStatusOnline

	existingPages, err := s.db.Page().List(ctx, database.PageByAPIKeyID(apikey.ID))
	if err != nil {
		return nil, false, nil, nil, nil, nil, err
	}

	existingPageMap := make(map[string]*core.Page)
	for _, p := range existingPages {
		existingPageMap[p.ID.String()] = p
	}

	requestPageIDs := make(map[string]struct{})
	for _, p := range in.Pages {
		requestPageIDs[p.Id] = struct{}{}
	}

	insertPages := make([]*core.Page, 0)
	updatePages := make([]*core.Page, 0)
	deletePages := make([]*core.Page, 0)
	for _, reqPage := range in.Pages {
		if existingPage, ok := existingPageMap[reqPage.Id]; ok {
			existingPage.Name = reqPage.Name
			existingPage.Route = reqPage.Route
			existingPage.Path = reqPage.Path
			updatePages = append(updatePages, existingPage)
		} else {
			pageID, err := uuid.FromString(reqPage.Id)
			if err != nil {
				return nil, false, nil, nil, nil, nil, err
			}
			newPage := &core.Page{
				ID:             pageID,
				OrganizationID: apikey.OrganizationID,
				EnvironmentID:  apikey.EnvironmentID,
				APIKeyID:       apikey.ID,
				Name:           reqPage.Name,
				Route:          reqPage.Route,
				Path:           reqPage.Path,
			}
			insertPages = append(insertPages, newPage)
		}
	}

	for _, existingPage := range existingPages {
		if _, exists := requestPageIDs[existingPage.ID.String()]; !exists {
			deletePages = append(deletePages, existingPage)
		}
	}

	return hostInstance, hostExists, apikey, insertPages, updatePages, deletePages, nil
}
