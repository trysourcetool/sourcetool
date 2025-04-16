package ws

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/model"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/pb/go/websocket/v1"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

type Service interface {
	InitializeClient(context.Context, *websocket.Conn, *websocketv1.Message) error
	InitializeHost(context.Context, *websocket.Conn, string, *websocketv1.Message) (*model.HostInstance, error)
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

	page, err := s.Store.Page().Get(ctx, storeopts.PageByID(pageID))
	if err != nil {
		return err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)
	if currentOrg.ID != page.OrganizationID {
		return errdefs.ErrPermissionDenied(errors.New("organization mismatch"))
	}

	apiKey, err := s.Store.APIKey().Get(ctx, storeopts.APIKeyByID(page.APIKeyID))
	if err != nil {
		return err
	}

	hostInstances, err := s.Store.HostInstance().List(ctx, storeopts.HostInstanceByAPIKeyID(apiKey.ID))
	if err != nil {
		return err
	}

	// Try to find an online host that responds to ping
	var onlineHostInstance *model.HostInstance
	connManager := GetConnManager()

	// First, try hosts that are already marked as online
	for _, hostInstance := range hostInstances {
		if hostInstance.Status == model.HostInstanceStatusOnline {
			if err := connManager.PingHost(hostInstance.ID); err != nil {
				// Update host status to offline if ping fails
				hostInstance.Status = model.HostInstanceStatusOffline
				if err := s.Store.HostInstance().Update(ctx, hostInstance); err != nil {
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
			if hostInstance.Status == model.HostInstanceStatusUnreachable {
				if err := connManager.PingHost(hostInstance.ID); err == nil {
					// Host is actually reachable, update its status
					hostInstance.Status = model.HostInstanceStatusOnline
					if err := s.Store.HostInstance().Update(ctx, hostInstance); err != nil {
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

	var sess *model.Session
	var sessionExists bool
	if conv.SafeValue(in.SessionId) != "" {
		sessionID, err := uuid.FromString(conv.SafeValue(in.SessionId))
		if err != nil {
			return errdefs.ErrSessionNotFound(err)
		}

		sess, err = s.Store.Session().Get(ctx, storeopts.SessionByID(sessionID))
		if err != nil {
			return err
		}
		sessionExists = true
	} else {
		sess = &model.Session{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: page.OrganizationID,
			APIKeyID:       page.APIKeyID,
			HostInstanceID: onlineHostInstance.ID,
			UserID:         currentUser.ID,
		}
		sessionExists = false
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if !sessionExists {
			if err := tx.Session().Create(ctx, sess); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	sess, err = s.Store.Session().Get(ctx, storeopts.SessionByID(sess.ID))
	if err != nil {
		return err
	}

	if err := SendResponse(conn, &websocketv1.Message{
		Id: msg.Id,
		Type: &websocketv1.Message_InitializeClientCompleted{
			InitializeClientCompleted: &websocketv1.InitializeClientCompleted{
				SessionId: sess.ID.String(),
			},
		},
	}); err != nil {
		return err
	}

	GetConnManager().SetConnectedClient(sess, conn)

	if err := GetConnManager().SendToHost(ctx, onlineHostInstance.ID, &websocketv1.Message{
		Id: uuid.Must(uuid.NewV4()).String(),
		Type: &websocketv1.Message_InitializeClient{
			InitializeClient: &websocketv1.InitializeClient{
				SessionId: conv.NilValue(sess.ID.String()),
				PageId:    page.ID.String(),
			},
		},
	}); err != nil {
		s.Store.Session().Delete(ctx, sess)
		GetConnManager().DisconnectClient(sess.ID)
		logger.Logger.Sugar().Errorf("Failed to send initialize client message to host: %v", err)
		return err
	}

	return nil
}

func (s *ServiceCE) InitializeHost(ctx context.Context, conn *websocket.Conn, instanceID string, msg *websocketv1.Message) (*model.HostInstance, error) {
	in := msg.GetInitializeHost()
	if in == nil {
		return nil, errors.New("invalid message")
	}

	apikey, err := s.Store.APIKey().Get(ctx, storeopts.APIKeyByKey(in.ApiKey))
	if err != nil {
		return nil, err
	}

	hostInstanceID, err := uuid.FromString(instanceID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	hostInstance, err := s.Store.HostInstance().Get(ctx, storeopts.HostInstanceByID(hostInstanceID))
	if err != nil && !errdefs.IsHostInstanceNotFound(err) {
		return nil, err
	}

	hostExists := hostInstance != nil

	if !hostExists {
		hostInstance = &model.HostInstance{
			ID:             hostInstanceID,
			OrganizationID: apikey.OrganizationID,
			APIKeyID:       apikey.ID,
		}
	}

	hostInstance.SDKName = in.SdkName
	hostInstance.SDKVersion = in.SdkVersion
	hostInstance.Status = model.HostInstanceStatusOnline

	existingPages, err := s.Store.Page().List(ctx, storeopts.PageByAPIKeyID(apikey.ID))
	if err != nil {
		return nil, err
	}

	existingPageMap := make(map[string]*model.Page)
	for _, p := range existingPages {
		existingPageMap[p.ID.String()] = p
	}

	requestPageIDs := make(map[string]struct{})
	for _, p := range in.Pages {
		requestPageIDs[p.Id] = struct{}{}
	}

	insertPages := make([]*model.Page, 0)
	updatePages := make([]*model.Page, 0)
	deletePages := make([]*model.Page, 0)
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
			newPage := &model.Page{
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

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
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

	GetConnManager().SetConnectedHost(hostInstance, apikey, conn)

	if err := SendResponse(conn, &websocketv1.Message{
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

	sess, err := s.Store.Session().Get(ctx, storeopts.SessionByID(sessionID))
	if err != nil {
		return err
	}

	pageID, err := uuid.FromString(in.PageId)
	if err != nil {
		return err
	}

	page, err := s.Store.Page().Get(ctx, storeopts.PageByID(pageID), storeopts.PageBySessionID(sess.ID))
	if err != nil {
		return err
	}

	if err := GetConnManager().SendToHost(ctx, sess.HostInstanceID, &websocketv1.Message{
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

	_, err = s.Store.Session().Get(ctx, storeopts.SessionByID(sessionID))
	if err != nil {
		return err
	}

	if err := GetConnManager().SendToClient(ctx, sessionID, msg); err != nil {
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

	sess, err := s.Store.Session().Get(ctx, storeopts.SessionByID(sessionID))
	if err != nil {
		return err
	}

	_, err = s.Store.Page().Get(ctx, storeopts.PageByAPIKeyID(sess.APIKeyID), storeopts.PageBySessionID(sess.ID))
	if err != nil {
		return err
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.Session().Delete(ctx, sess); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if err := GetConnManager().SendToHost(ctx, sess.HostInstanceID, &websocketv1.Message{
		Id: msg.Id,
		Type: &websocketv1.Message_CloseSession{
			CloseSession: &websocketv1.CloseSession{
				SessionId: sess.ID.String(),
			},
		},
	}); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send close session message to host: %v", err)
		return err
	}

	GetConnManager().DisconnectClient(sess.ID)

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
		return err
	}

	_, err = s.Store.Session().Get(ctx, storeopts.SessionByID(sessionID))
	if err != nil {
		return err
	}

	if err := GetConnManager().SendToClient(ctx, sessionID, msg); err != nil {
		return err
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
		return err
	}

	_, err = s.Store.Session().Get(ctx, storeopts.SessionByID(sessionID))
	if err != nil {
		return err
	}

	if err := GetConnManager().SendToClient(ctx, sessionID, msg); err != nil {
		return err
	}

	return nil
}

func (s *ServiceCE) UpdateStatus(ctx context.Context, in dto.UpdateHostInstanceStatusInput) (*dto.UpdateHostInstanceStatusOutput, error) {
	hostInstanceID, err := uuid.FromString(in.ID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	host, err := s.Store.HostInstance().Get(ctx, storeopts.HostInstanceByID(hostInstanceID))
	if err != nil {
		return nil, err
	}

	host.Status = in.Status

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
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
