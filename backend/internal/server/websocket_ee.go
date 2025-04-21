//go:build ee
// +build ee

package server

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
)

func (s *Server) handleInitializeHost(ctx context.Context, conn *websocket.Conn, instanceID string, msg *websocketv1.Message) (*core.HostInstance, error) {
	in := msg.GetInitializeHost()
	if in == nil {
		return nil, errors.New("invalid message")
	}

	hostInstance, hostExists, apikey, insertPages, updatePages, deletePages, err := s.handleInitializeHostBase(ctx, conn, instanceID, msg)
	if err != nil {
		return nil, err
	}

	var allGroupSlugs []string
	for _, p := range in.Pages {
		allGroupSlugs = append(allGroupSlugs, p.Groups...)
	}
	groups, err := s.db.Group().List(ctx, database.GroupByOrganizationID(apikey.OrganizationID), database.GroupBySlugs(allGroupSlugs))
	if err != nil {
		return nil, err
	}
	groupMap := make(map[string]*core.Group)
	for _, g := range groups {
		groupMap[g.Slug] = g
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if hostExists {
			if err := tx.HostInstance().Update(ctx, hostInstance); err != nil {
				return err
			}
		} else {
			if err := s.db.HostInstance().Create(ctx, hostInstance); err != nil {
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

		existingGroupPages, err := s.db.Group().ListPages(ctx, database.GroupPageByPageIDs(pageIDs))
		if err != nil {
			return err
		}

		if len(existingGroupPages) > 0 {
			if err := tx.Group().BulkDeletePages(ctx, existingGroupPages); err != nil {
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
			if err := s.db.Group().BulkInsertPages(ctx, newGroupPages); err != nil {
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
