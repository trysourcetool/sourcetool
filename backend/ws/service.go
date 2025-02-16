package ws

import (
	"context"
	"errors"
	"strconv"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool/backend/conv"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/server/ws/types"
	websocketv1 "github.com/trysourcetool/sourcetool/proto/go/websocket/v1"
)

type Service interface {
	InitializeClient(context.Context, *websocketv1.Message) error
	InitializeHost(context.Context, string, *websocketv1.Message) (*model.HostInstance, error)
	RerunPage(context.Context, *websocketv1.Message) error
	RenderWidget(context.Context, *websocketv1.Message) error
	CloseSession(context.Context, *websocketv1.Message) error
	ScriptFinished(context.Context, *websocketv1.Message) error
	UpdateStatus(context.Context, types.UpdateHostInstanceStatusInput) (*types.UpdateHostInstanceStatusPayload, error)
}

type service struct {
	conn *websocket.Conn
	*infra.Dependency
}

func NewService(conn *websocket.Conn, d *infra.Dependency) Service {
	return &service{conn: conn, Dependency: d}
}

func (s *service) InitializeClient(ctx context.Context, msg *websocketv1.Message) error {
	in := msg.GetInitializeClient()
	if in == nil {
		return errors.New("invalid message")
	}

	pageID, err := uuid.FromString(in.PageId)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	page, err := s.Store.Page().Get(ctx, model.PageByID(pageID))
	if err != nil {
		return err
	}

	currentOrg := ctxutils.CurrentOrganization(ctx)
	if currentOrg.ID != page.OrganizationID {
		return errdefs.ErrPermissionDenied(errors.New("organization mismatch"))
	}

	apiKey, err := s.Store.APIKey().Get(ctx, model.APIKeyByID(page.APIKeyID))
	if err != nil {
		return err
	}

	hostInstances, err := s.Store.HostInstance().List(ctx, model.HostInstanceByAPIKeyID(apiKey.ID))
	if err != nil {
		return err
	}

	var onlineHostInstance *model.HostInstance
	for _, hostInstance := range hostInstances {
		if hostInstance.Status == model.HostInstanceStatusOnline {
			connManager := GetConnManager()
			if err := connManager.PingHost(hostInstance.ID); err != nil {
				hostInstance.Status = model.HostInstanceStatusOffline
				if err := s.Store.HostInstance().Update(ctx, hostInstance); err != nil {
					return err
				}
				continue
			}

			onlineHostInstance = hostInstance
			break
		}
	}

	if onlineHostInstance == nil {
		return errdefs.ErrHostInstanceStatusNotOnline(errors.New("host instance status is not online"))
	}

	currentUser := ctxutils.CurrentUser(ctx)

	var sess *model.Session
	var sessionExists bool
	if conv.SafeValue(in.SessionId) != "" {
		sessionID, err := uuid.FromString(conv.SafeValue(in.SessionId))
		if err != nil {
			return errdefs.ErrSessionNotFound(err)
		}

		sess, err = s.Store.Session().Get(ctx, model.SessionByID(sessionID))
		if err != nil {
			return err
		}
		sessionExists = true
	} else {
		sess = &model.Session{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: page.OrganizationID,
			PageID:         page.ID,
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

	sess, err = s.Store.Session().Get(ctx, model.SessionByID(sess.ID))
	if err != nil {
		return err
	}

	if err := SendResponse(s.conn, &websocketv1.Message{
		Id: msg.Id,
		Type: &websocketv1.Message_InitializeClientCompleted{
			InitializeClientCompleted: &websocketv1.InitializeClientCompleted{
				SessionId: sess.ID.String(),
			},
		},
	}); err != nil {
		return err
	}

	GetConnManager().SetConnectedClient(sess, s.conn)

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

func (s *service) InitializeHost(ctx context.Context, instanceID string, msg *websocketv1.Message) (*model.HostInstance, error) {
	in := msg.GetInitializeHost()
	if in == nil {
		return nil, errors.New("invalid message")
	}

	apikey, err := s.Store.APIKey().Get(ctx, model.APIKeyByKey(in.ApiKey))
	if err != nil {
		return nil, err
	}

	hostInstanceID, err := uuid.FromString(instanceID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	hostInstance, err := s.Store.HostInstance().Get(ctx, model.HostInstanceByID(hostInstanceID))
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

	existingPages, err := s.Store.Page().List(ctx, model.PageByAPIKeyID(apikey.ID))
	if err != nil {
		return nil, err
	}

	existingPageMap := make(map[string]*model.Page)
	for _, p := range existingPages {
		existingPageMap[p.ID.String()] = p
	}

	var allGroupSlugs []string
	for _, p := range in.Pages {
		allGroupSlugs = append(allGroupSlugs, p.Groups...)
	}
	groups, err := s.Store.Group().List(ctx, model.GroupByOrganizationID(apikey.OrganizationID), model.GroupBySlugs(allGroupSlugs))
	if err != nil {
		return nil, err
	}
	groupMap := make(map[string]*model.Group)
	for _, g := range groups {
		groupMap[g.Slug] = g
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

		if len(insertPages) > 0 {
			if err := tx.Page().BulkInsert(ctx, insertPages); err != nil {
				return err
			}
		}
		if len(updatePages) > 0 {
			if err := tx.Page().BulkUpdate(ctx, updatePages); err != nil {
				return err
			}
		}
		if len(deletePages) > 0 {
			if err := tx.Page().BulkDelete(ctx, deletePages); err != nil {
				return err
			}
		}

		var pageIDs []uuid.UUID
		pageGroupMap := make(map[uuid.UUID][]string) // pageID -> group slugs
		for _, reqPage := range in.Pages {
			pageID, err := uuid.FromString(reqPage.Id)
			if err != nil {
				return err
			}
			pageIDs = append(pageIDs, pageID)
			pageGroupMap[pageID] = reqPage.Groups
		}

		existingGroupPages, err := tx.Group().ListPages(ctx, model.GroupPageByPageIDs(pageIDs))
		if err != nil {
			return err
		}

		if len(existingGroupPages) > 0 {
			if err := tx.Group().BulkDeletePages(ctx, existingGroupPages); err != nil {
				return err
			}
		}

		var newGroupPages []*model.GroupPage
		for pageID, groupSlugs := range pageGroupMap {
			for _, slug := range groupSlugs {
				group, ok := groupMap[slug]
				if !ok {
					continue
				}
				newGroupPages = append(newGroupPages, &model.GroupPage{
					ID:      uuid.Must(uuid.NewV4()),
					GroupID: group.ID,
					PageID:  pageID,
				})
			}
		}

		if len(newGroupPages) > 0 {
			if err := tx.Group().BulkInsertPages(ctx, newGroupPages); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	GetConnManager().SetConnectedHost(hostInstance, apikey, s.conn)

	if err := SendResponse(s.conn, &websocketv1.Message{
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

func (s *service) RerunPage(ctx context.Context, msg *websocketv1.Message) error {
	in := msg.GetRerunPage()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return err
	}

	sess, err := s.Store.Session().Get(ctx, model.SessionByID(sessionID))
	if err != nil {
		return err
	}

	pageID, err := uuid.FromString(in.PageId)
	if err != nil {
		return err
	}

	page, err := s.Store.Page().Get(ctx, model.PageByID(pageID), model.PageBySessionID(sess.ID))
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

func (s *service) RenderWidget(ctx context.Context, msg *websocketv1.Message) error {
	in := msg.GetRenderWidget()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return err
	}

	_, err = s.Store.Session().Get(ctx, model.SessionByID(sessionID))
	if err != nil {
		return err
	}

	if err := GetConnManager().SendToClient(ctx, sessionID, msg); err != nil {
		return err
	}

	return nil
}

func (s *service) CloseSession(ctx context.Context, msg *websocketv1.Message) error {
	in := msg.GetCloseSession()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return errdefs.ErrAPIKeyNotFound(err)
	}

	sess, err := s.Store.Session().Get(ctx, model.SessionByID(sessionID))
	if err != nil {
		return err
	}

	_, err = s.Store.Page().Get(ctx, model.PageByID(sess.PageID), model.PageBySessionID(sess.ID))
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

func (s *service) ScriptFinished(ctx context.Context, msg *websocketv1.Message) error {
	in := msg.GetScriptFinished()
	if in == nil {
		return errors.New("invalid message")
	}

	logger.Logger.Sugar().Debug("Payload: ", in)

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return err
	}

	_, err = s.Store.Session().Get(ctx, model.SessionByID(sessionID))
	if err != nil {
		return err
	}

	if err := GetConnManager().SendToClient(ctx, sessionID, msg); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateStatus(ctx context.Context, in types.UpdateHostInstanceStatusInput) (*types.UpdateHostInstanceStatusPayload, error) {
	hostInstanceID, err := uuid.FromString(in.ID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	host, err := s.Store.HostInstance().Get(ctx, model.HostInstanceByID(hostInstanceID))
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

	return &types.UpdateHostInstanceStatusPayload{
		HostInstance: &types.HostInstancePayload{
			ID:         host.ID.String(),
			SDKName:    host.SDKName,
			SDKVersion: host.SDKVersion,
			Status:     host.Status.String(),
			CreatedAt:  strconv.FormatInt(host.CreatedAt.Unix(), 10),
			UpdatedAt:  strconv.FormatInt(host.UpdatedAt.Unix(), 10),
		},
	}, nil
}
