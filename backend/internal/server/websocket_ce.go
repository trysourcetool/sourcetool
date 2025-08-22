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
	hostInstance, hostExists, apikey, insertPages, updatePages, deletePages, insertAgents, updateAgents, deleteAgents, agentTools, err := s.handleInitializeHostBase(ctx, conn, instanceID, msg)
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

		// Handle pages
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

		// Handle agents
		if len(deleteAgents) > 0 {
			if err := tx.Agent().BulkDelete(ctx, deleteAgents); err != nil {
				return err
			}
		}
		if len(updateAgents) > 0 {
			if err := tx.Agent().BulkUpdate(ctx, updateAgents); err != nil {
				return err
			}
		}
		if len(insertAgents) > 0 {
			if err := tx.Agent().BulkInsert(ctx, insertAgents); err != nil {
				return err
			}
		}

		// Delete existing agent tools for agents being updated or deleted
		allAgentIDs := make([]uuid.UUID, 0)
		for _, agent := range updateAgents {
			allAgentIDs = append(allAgentIDs, agent.ID)
		}
		for _, agent := range deleteAgents {
			allAgentIDs = append(allAgentIDs, agent.ID)
		}
		if len(allAgentIDs) > 0 {
			existingAgentTools, err := tx.AgentTool().List(ctx, database.AgentToolByAgentID(allAgentIDs[0]))
			if err != nil {
				return err
			}
			var toolsToDelete []*core.AgentTool
			for _, tool := range existingAgentTools {
				for _, agentID := range allAgentIDs {
					if tool.AgentID == agentID {
						toolsToDelete = append(toolsToDelete, tool)
						break
					}
				}
			}
			if len(toolsToDelete) > 0 {
				if err := tx.AgentTool().BulkDelete(ctx, toolsToDelete); err != nil {
					return err
				}
			}
		}

		// Insert new agent tools
		if len(agentTools) > 0 {
			if err := tx.AgentTool().BulkInsert(ctx, agentTools); err != nil {
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
