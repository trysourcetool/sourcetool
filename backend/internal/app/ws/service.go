package ws

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/ctxutil"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/page"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/session"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/ws/message"
	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
	"github.com/trysourcetool/sourcetool/backend/pkg/ptrconv"
)

type Service interface {
	InitializeClient(context.Context, *websocket.Conn, *websocketv1.Message) error
	InitializeHost(context.Context, *websocket.Conn, string, *websocketv1.Message) (*hostinstance.HostInstance, error)
	RerunPage(context.Context, *websocket.Conn, *websocketv1.Message) error
	RenderWidget(context.Context, *websocket.Conn, *websocketv1.Message) error
	CloseSession(context.Context, *websocket.Conn, *websocketv1.Message) error
	ScriptFinished(context.Context, *websocket.Conn, *websocketv1.Message) error
	Exception(context.Context, *websocket.Conn, *websocketv1.Message) error
	UpdateStatus(context.Context, dto.UpdateHostInstanceStatusInput) (*dto.UpdateHostInstanceStatusOutput, error)
}

type ServiceCE struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) InitializeClient(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetInitializeClient()
	if in == nil {
		return errors.New("invalid message")
	}

	pageID, err := uuid.FromString(in.PageId)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	page, err := s.Repository.Page().Get(ctx, page.ByID(pageID))
	if err != nil {
		return err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)
	if currentOrg.ID != page.OrganizationID {
		return errdefs.ErrPermissionDenied(errors.New("organization mismatch"))
	}

	apiKey, err := s.Repository.APIKey().Get(ctx, apikey.ByID(page.APIKeyID))
	if err != nil {
		return err
	}

	hostInstances, err := s.Repository.HostInstance().List(ctx, hostinstance.ByAPIKeyID(apiKey.ID))
	if err != nil {
		return err
	}

	// Try to find an online host that responds to ping
	var onlineHostInstance *hostinstance.HostInstance
	for _, hostInstance := range hostInstances {
		if hostInstance.Status == hostinstance.HostInstanceStatusOnline {
			if err := s.WSManager.PingConnectedHost(hostInstance.ID); err != nil {
				// Update host status to offline if ping fails
				hostInstance.Status = hostinstance.HostInstanceStatusOffline
				if err := s.Repository.HostInstance().Update(ctx, hostInstance); err != nil {
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
			if hostInstance.Status == hostinstance.HostInstanceStatusUnreachable {
				if err := s.WSManager.PingConnectedHost(hostInstance.ID); err == nil {
					// Host is actually reachable, update its status
					hostInstance.Status = hostinstance.HostInstanceStatusOnline
					if err := s.Repository.HostInstance().Update(ctx, hostInstance); err != nil {
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

	currentUser := ctxutil.CurrentUser(ctx)

	var sess *session.Session
	var sessionExists bool
	if ptrconv.SafeValue(in.SessionId) != "" {
		sessionID, err := uuid.FromString(ptrconv.SafeValue(in.SessionId))
		if err != nil {
			return errdefs.ErrSessionNotFound(err)
		}

		sess, err = s.Repository.Session().Get(ctx, session.ByID(sessionID))
		if err != nil {
			return err
		}
		sessionExists = true
	} else {
		sess = &session.Session{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: page.OrganizationID,
			APIKeyID:       page.APIKeyID,
			HostInstanceID: onlineHostInstance.ID,
			UserID:         currentUser.ID,
		}
		sessionExists = false
	}

	if err := s.Repository.RunTransaction(func(tx db.Transaction) error {
		if !sessionExists {
			if err := tx.Session().Create(ctx, sess); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	sess, err = s.Repository.Session().Get(ctx, session.ByID(sess.ID))
	if err != nil {
		return err
	}

	if err := message.SendResponse(conn, &websocketv1.Message{
		Id: msg.Id,
		Type: &websocketv1.Message_InitializeClientCompleted{
			InitializeClientCompleted: &websocketv1.InitializeClientCompleted{
				SessionId: sess.ID.String(),
			},
		},
	}); err != nil {
		return err
	}

	s.WSManager.SetConnectedClient(sess, conn)

	if err := s.WSManager.SendToHost(ctx, onlineHostInstance.ID, &websocketv1.Message{
		Id: uuid.Must(uuid.NewV4()).String(),
		Type: &websocketv1.Message_InitializeClient{
			InitializeClient: &websocketv1.InitializeClient{
				SessionId: ptrconv.NilValue(sess.ID.String()),
				PageId:    page.ID.String(),
			},
		},
	}); err != nil {
		s.Repository.Session().Delete(ctx, sess)
		s.WSManager.DisconnectClient(sess.ID)
		logger.Logger.Sugar().Errorf("Failed to send initialize client message to host: %v", err)
		return err
	}

	return nil
}

func (s *ServiceCE) InitializeHost(ctx context.Context, conn *websocket.Conn, instanceID string, msg *websocketv1.Message) (*hostinstance.HostInstance, error) {
	in := msg.GetInitializeHost()
	if in == nil {
		return nil, errors.New("invalid message")
	}

	apikey, err := s.Repository.APIKey().Get(ctx, apikey.ByKey(in.ApiKey))
	if err != nil {
		return nil, err
	}

	hostInstanceID, err := uuid.FromString(instanceID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	hostInstance, err := s.Repository.HostInstance().Get(ctx, hostinstance.ByID(hostInstanceID))
	if err != nil && !errdefs.IsHostInstanceNotFound(err) {
		return nil, err
	}

	hostExists := hostInstance != nil

	if !hostExists {
		hostInstance = &hostinstance.HostInstance{
			ID:             hostInstanceID,
			OrganizationID: apikey.OrganizationID,
			APIKeyID:       apikey.ID,
		}
	}

	hostInstance.SDKName = in.SdkName
	hostInstance.SDKVersion = in.SdkVersion
	hostInstance.Status = hostinstance.HostInstanceStatusOnline

	existingPages, err := s.Repository.Page().List(ctx, page.ByAPIKeyID(apikey.ID))
	if err != nil {
		return nil, err
	}

	existingPageMap := make(map[string]*page.Page)
	for _, p := range existingPages {
		existingPageMap[p.ID.String()] = p
	}

	requestPageIDs := make(map[string]struct{})
	for _, p := range in.Pages {
		requestPageIDs[p.Id] = struct{}{}
	}

	insertPages := make([]*page.Page, 0)
	updatePages := make([]*page.Page, 0)
	deletePages := make([]*page.Page, 0)
	for _, reqPage := range in.Pages {
		if existingPage, ok := existingPageMap[reqPage.Id]; ok {
			existingPage.Name = reqPage.Name
			existingPage.Route = reqPage.Route
			existingPage.Path = reqPage.Path
			updatePages = append(updatePages, existingPage)
		} else {
			pageID, err := uuid.FromString(reqPage.Id)
			if err != nil {
				return nil, err
			}
			newPage := &page.Page{
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

	if err := s.Repository.RunTransaction(func(tx db.Transaction) error {
		if hostExists {
			if err := tx.HostInstance().Update(ctx, hostInstance); err != nil {
				return err
			}
		} else {
			if err := tx.HostInstance().Create(ctx, hostInstance); err != nil {
				return err
			}
		}

		if len(deletePages) > 0 {
			if err := tx.Page().BulkDelete(ctx, deletePages); err != nil {
				return err
			}
		}
		if len(updatePages) > 0 {
			if err := tx.Page().BulkUpdate(ctx, updatePages); err != nil {
				return err
			}
		}
		if len(insertPages) > 0 {
			if err := tx.Page().BulkInsert(ctx, insertPages); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	s.WSManager.SetConnectedHost(hostInstance, apikey, conn)

	if err := message.SendResponse(conn, &websocketv1.Message{
		Id: msg.Id,
		Type: &websocketv1.Message_InitializeHostCompleted{
			InitializeHostCompleted: &websocketv1.InitializeHostCompleted{
				HostInstanceId: hostInstance.ID.String(),
			},
		},
	}); err != nil {
		return nil, err
	}

	return hostInstance, nil
}

func (s *ServiceCE) RerunPage(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetRerunPage()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return err
	}

	sess, err := s.Repository.Session().Get(ctx, session.ByID(sessionID))
	if err != nil {
		return err
	}

	pageID, err := uuid.FromString(in.PageId)
	if err != nil {
		return err
	}

	page, err := s.Repository.Page().Get(ctx, page.ByID(pageID), page.BySessionID(sess.ID))
	if err != nil {
		return err
	}

	if err := s.WSManager.SendToHost(ctx, sess.HostInstanceID, &websocketv1.Message{
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

func (s *ServiceCE) RenderWidget(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetRenderWidget()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return err
	}

	_, err = s.Repository.Session().Get(ctx, session.ByID(sessionID))
	if err != nil {
		return err
	}

	if err := s.WSManager.SendToClient(ctx, sessionID, msg); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send render widget message to client: %v", err)
		return err
	}

	return nil
}

func (s *ServiceCE) CloseSession(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetCloseSession()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return errdefs.ErrAPIKeyNotFound(err)
	}

	sess, err := s.Repository.Session().Get(ctx, session.ByID(sessionID))
	if err != nil {
		return err
	}

	_, err = s.Repository.Page().Get(ctx, page.ByAPIKeyID(sess.APIKeyID), page.BySessionID(sess.ID))
	if err != nil {
		return err
	}

	if err := s.Repository.RunTransaction(func(tx db.Transaction) error {
		if err := tx.Session().Delete(ctx, sess); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if err := s.WSManager.SendToHost(ctx, sess.HostInstanceID, &websocketv1.Message{
		Id: uuid.Must(uuid.NewV4()).String(),
		Type: &websocketv1.Message_CloseSession{
			CloseSession: &websocketv1.CloseSession{
				SessionId: sess.ID.String(),
			},
		},
	}); err != nil {
		logger.Logger.Sugar().Warnf("Failed to send close session message to host %s for session %s: %v", sess.HostInstanceID, sess.ID, err)
	}

	s.WSManager.DisconnectClient(sess.ID)

	return nil
}

func (s *ServiceCE) ScriptFinished(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetScriptFinished()
	if in == nil {
		return errors.New("invalid message")
	}

	logger.Logger.Sugar().Debug("Payload: ", in)

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	_, err = s.Repository.Session().Get(ctx, session.ByID(sessionID))
	if err != nil {
		return err
	}

	if err := s.WSManager.SendToClient(ctx, sessionID, msg); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send script finished message to client: %v", err)
	}

	return nil
}

func (s *ServiceCE) Exception(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetException()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	_, err = s.Repository.Session().Get(ctx, session.ByID(sessionID))
	if err != nil {
		return err
	}

	if err := s.WSManager.SendToClient(ctx, sessionID, msg); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send exception message to client: %v", err)
	}

	return nil
}

func (s *ServiceCE) UpdateStatus(ctx context.Context, in dto.UpdateHostInstanceStatusInput) (*dto.UpdateHostInstanceStatusOutput, error) {
	hostInstanceID, err := uuid.FromString(in.ID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	host, err := s.Repository.HostInstance().Get(ctx, hostinstance.ByID(hostInstanceID))
	if err != nil {
		return nil, err
	}

	host.Status = in.Status

	if err := s.Repository.RunTransaction(func(tx db.Transaction) error {
		if err := tx.HostInstance().Update(ctx, host); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &dto.UpdateHostInstanceStatusOutput{
		HostInstance: dto.HostInstanceFromModel(host),
	}, nil
}
