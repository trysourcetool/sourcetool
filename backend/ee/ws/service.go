package ws

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/generated/proto/websocket/v1"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
	"github.com/trysourcetool/sourcetool/backend/ws"
)

type serviceEE struct {
	*infra.Dependency
	*ws.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE:  ws.NewServiceCE(infra.NewDependency(d.Store, d.Mailer)),
	}
}

func (s *serviceEE) InitializeHost(ctx context.Context, instanceID string, msg *websocketv1.Message) (*model.HostInstance, error) {
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

	var allGroupSlugs []string
	for _, p := range in.Pages {
		allGroupSlugs = append(allGroupSlugs, p.Groups...)
	}
	groups, err := s.Store.Group().List(ctx, storeopts.GroupByOrganizationID(apikey.OrganizationID), storeopts.GroupBySlugs(allGroupSlugs))
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

		existingGroupPages, err := tx.Group().ListPages(ctx, storeopts.GroupPageByPageIDs(pageIDs))
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

	ws.GetConnManager().SetConnectedHost(hostInstance, apikey, s.GetConn())

	if err := ws.SendResponse(s.GetConn(), &websocketv1.Message{
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
