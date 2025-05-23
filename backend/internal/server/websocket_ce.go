//go:build !ee
// +build !ee

package server

import (
	"context"

	"github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool/backend/internal/database"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
)

func (s *Server) handleInitializeHost(ctx context.Context, conn *websocket.Conn, instanceID string, msg *websocketv1.Message) error {
	hostInstance, hostExists, apikey, insertPages, updatePages, deletePages, err := s.handleInitializeHostBase(ctx, conn, instanceID, msg)
	if err != nil {
		return err
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
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
		return err
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
		return err
	}

	return nil
}
