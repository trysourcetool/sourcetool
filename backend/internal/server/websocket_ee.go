//go:build ee
// +build ee

package server

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
	"github.com/trysourcetool/sourcetool/backend/internal/postgres"
)

func (s *Server) wsInitializeHost(ctx context.Context, conn *websocket.Conn, instanceID string, msg *websocketv1.Message) (*core.HostInstance, error) {
	in := msg.GetInitializeHost()
	if in == nil {
		return nil, errors.New("invalid message")
	}

	apikey, err := s.db.GetAPIKey(ctx, postgres.APIKeyByKey(in.ApiKey))
	if err != nil {
		return nil, err
	}

	hostInstanceID, err := uuid.FromString(instanceID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	hostInstance, err := s.db.GetHostInstance(ctx, postgres.HostInstanceByID(hostInstanceID))
	if err != nil && !errdefs.IsHostInstanceNotFound(err) {
		return nil, err
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

	existingPages, err := s.db.ListPages(ctx, postgres.PageByAPIKeyID(apikey.ID))
	if err != nil {
		return nil, err
	}

	existingPageMap := make(map[string]*core.Page)
	for _, p := range existingPages {
		existingPageMap[p.ID.String()] = p
	}

	var allGroupSlugs []string
	for _, p := range in.Pages {
		allGroupSlugs = append(allGroupSlugs, p.Groups...)
	}
	groups, err := s.db.ListGroups(ctx, postgres.GroupByOrganizationID(apikey.OrganizationID), postgres.GroupBySlugs(allGroupSlugs))
	if err != nil {
		return nil, err
	}
	groupMap := make(map[string]*core.Group)
	for _, g := range groups {
		groupMap[g.Slug] = g
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
				return nil, err
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

	if err := s.db.WithTx(ctx, func(tx *sqlx.Tx) error {
		if hostExists {
			if err := s.db.UpdateHostInstance(ctx, tx, hostInstance); err != nil {
				return err
			}
		} else {
			if err := s.db.CreateHostInstance(ctx, tx, hostInstance); err != nil {
				return err
			}
		}

		if len(deletePages) > 0 {
			if err := s.db.BulkDeletePages(ctx, tx, deletePages); err != nil {
				return err
			}
		}
		if len(updatePages) > 0 {
			if err := s.db.BulkUpdatePages(ctx, tx, updatePages); err != nil {
				return err
			}
		}
		if len(insertPages) > 0 {
			if err := s.db.BulkInsertPages(ctx, tx, insertPages); err != nil {
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

		existingGroupPages, err := s.db.ListGroupPages(ctx, postgres.GroupPageByPageIDs(pageIDs))
		if err != nil {
			return err
		}

		if len(existingGroupPages) > 0 {
			if err := s.db.BulkDeleteGroupPages(ctx, tx, existingGroupPages); err != nil {
				return err
			}
		}

		var newGroupPages []*core.GroupPage
		for pageID, groupSlugs := range pageGroupMap {
			for _, slug := range groupSlugs {
				g, ok := groupMap[slug]
				if !ok {
					continue
				}
				newGroupPages = append(newGroupPages, &core.GroupPage{
					ID:      uuid.Must(uuid.NewV4()),
					GroupID: g.ID,
					PageID:  pageID,
				})
			}
		}

		if len(newGroupPages) > 0 {
			if err := s.db.BulkInsertGroupPages(ctx, tx, newGroupPages); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	s.wsManager.SetConnectedHost(hostInstance, apikey, conn)

	if err := s.sendWebSocketMessage(conn, &websocketv1.Message{
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
