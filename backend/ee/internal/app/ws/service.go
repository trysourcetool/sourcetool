package ws

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	gorillaws "github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	wsSvc "github.com/trysourcetool/sourcetool/backend/internal/app/ws"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/group"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/page"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/ws/message"
)

type serviceEE struct {
	*port.Dependencies
	*wsSvc.ServiceCE
}

func NewServiceEE(d *port.Dependencies) *serviceEE {
	return &serviceEE{
		Dependencies: d,
		ServiceCE:    wsSvc.NewServiceCE(port.NewDependencies(d.Repository, d.Mailer, d.PubSub, d.WSManager)),
	}
}

func (s *serviceEE) InitializeHost(ctx context.Context, conn *gorillaws.Conn, instanceID string, msg *websocketv1.Message) (*hostinstance.HostInstance, error) {
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

	var allGroupSlugs []string
	for _, p := range in.Pages {
		allGroupSlugs = append(allGroupSlugs, p.Groups...)
	}
	groups, err := s.Repository.Group().List(ctx, group.ByOrganizationID(apikey.OrganizationID), group.BySlugs(allGroupSlugs))
	if err != nil {
		return nil, err
	}
	groupMap := make(map[string]*group.Group)
	for _, g := range groups {
		groupMap[g.Slug] = g
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

	if err := s.Repository.RunTransaction(func(tx port.Transaction) error {
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

		existingGroupPages, err := tx.Group().ListPages(ctx, group.PageByPageIDs(pageIDs))
		if err != nil {
			return err
		}

		if len(existingGroupPages) > 0 {
			if err := tx.Group().BulkDeletePages(ctx, existingGroupPages); err != nil {
				return err
			}
		}

		var newGroupPages []*group.GroupPage
		for pageID, groupSlugs := range pageGroupMap {
			for _, slug := range groupSlugs {
				g, ok := groupMap[slug]
				if !ok {
					continue
				}
				newGroupPages = append(newGroupPages, &group.GroupPage{
					ID:      uuid.Must(uuid.NewV4()),
					GroupID: g.ID,
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
