package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
)

type agentResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Instructions string `json:"instructions"`
	Model        string `json:"model"`
}

type listAgentsResponse struct {
	Agents []agentResponse `json:"agents"`
}

func (s *Server) handleListAgents(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctxOrg := internal.ContextOrganization(ctx)

	environmentID := r.URL.Query().Get("environmentId")
	if environmentID == "" {
		return errdefs.ErrInvalidArgument(errors.New("environmentId is required"))
	}

	envUUID, err := uuid.FromString(environmentID)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	// Verify environment belongs to organization
	env, err := s.db.Environment().Get(ctx, database.EnvironmentByID(envUUID))
	if err != nil {
		return err
	}
	if env.OrganizationID != ctxOrg.ID {
		return errdefs.ErrPermissionDenied(nil)
	}

	// Get agents for environment
	agents, err := s.db.Agent().List(ctx, database.AgentByEnvironmentID(envUUID))
	if err != nil {
		return err
	}

	response := listAgentsResponse{
		Agents: make([]agentResponse, len(agents)),
	}

	for i, agent := range agents {
		response.Agents[i] = agentResponse{
			ID:           agent.ID.String(),
			Name:         agent.Name,
			Description:  agent.Description,
			Instructions: agent.Instructions,
			Model:        agent.Model,
		}
	}

	return s.renderJSON(w, http.StatusOK, response)
}

func (s *Server) handleGetAgent(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctxOrg := internal.ContextOrganization(ctx)

	agentIDStr := chi.URLParam(r, "agentID")
	agentID, err := uuid.FromString(agentIDStr)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	agent, err := s.db.Agent().Get(ctx, database.AgentByID(agentID))
	if err != nil {
		return err
	}

	// Verify agent belongs to organization
	if agent.OrganizationID != ctxOrg.ID {
		return errdefs.ErrPermissionDenied(nil)
	}

	response := agentResponse{
		ID:           agent.ID.String(),
		Name:         agent.Name,
		Description:  agent.Description,
		Instructions: agent.Instructions,
		Model:        agent.Model,
	}

	return s.renderJSON(w, http.StatusOK, response)
}

// Agent chat session management - Now handled by database and SDK
// The Backend only manages WebSocket connections and message routing

// handleInitializeAgentChat initializes a new agent chat session
// Similar to handleInitializeClient, it creates a database session and forwards to SDK
func (s *Server) handleInitializeAgentChat(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	logger.Logger.Info("Initializing agent chat",
		zap.String("message_id", msg.Id))

	in := msg.GetInitializeAgentChat()
	if in == nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("invalid initialize agent chat message"))
	}

	// Parse agent ID
	agentID, err := uuid.FromString(in.AgentId)
	if err != nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("invalid agent ID: %w", err))
	}

	// Get agent from database
	agent, err := s.db.Agent().Get(ctx, database.AgentByID(agentID))
	if err != nil {
		if errdefs.IsAgentNotFound(err) {
			return errdefs.ErrAgentNotFound(err)
		}
		return err
	}

	// Verify organization matches
	ctxOrg := internal.ContextOrganization(ctx)
	if ctxOrg.ID != agent.OrganizationID {
		return errdefs.ErrPermissionDenied(errors.New("organization mismatch"))
	}

	// Get API key for the agent
	apiKey, err := s.db.APIKey().Get(ctx, database.APIKeyByID(agent.APIKeyID))
	if err != nil {
		return err
	}

	// Get environment
	env, err := s.db.Environment().Get(ctx, database.EnvironmentByID(apiKey.EnvironmentID))
	if err != nil {
		return err
	}

	// Find online host instances (same logic as handleInitializeClient)
	hostInstances, err := s.db.HostInstance().List(ctx, database.HostInstanceByAPIKeyID(apiKey.ID))
	if err != nil {
		return err
	}

	// Try to find an online host that responds to ping
	var onlineHostInstance *core.HostInstance
	for _, hostInstance := range hostInstances {
		if hostInstance.Status == core.HostInstanceStatusOnline {
			if err := s.wsManager.PingConnectedHost(hostInstance.ID); err != nil {
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

	// Get user from context
	user := internal.ContextUser(ctx)
	if user == nil {
		return errdefs.ErrUnauthenticated(fmt.Errorf("user not found in context"))
	}

	// Create or retrieve session
	var sess *core.Session
	var sessionExists bool
	if in.SessionId != nil && *in.SessionId != "" {
		sessionID, err := uuid.FromString(*in.SessionId)
		if err != nil {
			return errdefs.ErrInvalidArgument(fmt.Errorf("invalid session ID: %w", err))
		}
		sess, err = s.db.Session().Get(ctx, database.SessionByID(sessionID))
		if err != nil {
			return err
		}
		sessionExists = true
	} else {
		sess = &core.Session{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: agent.OrganizationID,
			EnvironmentID:  env.ID,
			UserID:         user.ID,
			Type:           core.SessionTypeAgent,
		}
		sessionExists = false
	}

	// Create session in database (similar to handleInitializeClient)
	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if !sessionExists {
			if err := tx.Session().Create(ctx, sess); err != nil {
				return err
			}

			if err := tx.Session().CreateHostInstance(ctx, &core.SessionHostInstance{
				ID:             uuid.Must(uuid.NewV4()),
				SessionID:      sess.ID,
				HostInstanceID: onlineHostInstance.ID,
			}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// Send completion message to client first
	if err := s.sendWebSocketMessage(conn, &websocketv1.Message{
		Id: msg.Id,
		Type: &websocketv1.Message_InitializeAgentChatCompleted{
			InitializeAgentChatCompleted: &websocketv1.InitializeAgentChatCompleted{
				SessionId: sess.ID.String(),
				AgentId:   agent.ID.String(),
			},
		},
	}); err != nil {
		return err
	}

	// Register the client connection
	s.wsManager.SetConnectedClient(sess, conn)

	// Forward the initialization request to the SDK host
	if err := s.wsManager.SendToHost(ctx, onlineHostInstance.ID, &websocketv1.Message{
		Id: uuid.Must(uuid.NewV4()).String(),
		Type: &websocketv1.Message_InitializeAgentChat{
			InitializeAgentChat: &websocketv1.InitializeAgentChat{
				SessionId: internal.StringPtr(sess.ID.String()),
				AgentId:   agent.ID.String(),
			},
		},
	}); err != nil {
		s.db.Session().Delete(ctx, sess)
		s.wsManager.DisconnectClient(sess.ID)
		logger.Logger.Sugar().Errorf("Failed to send initialize agent chat message to host: %v", err)
		return err
	}

	logger.Logger.Info("Agent chat initialized and forwarded to SDK",
		zap.String("session_id", sess.ID.String()),
		zap.String("agent_id", agent.ID.String()),
		zap.String("user_id", user.ID.String()),
		zap.String("host_instance_id", onlineHostInstance.ID.String()))

	return nil
}

// handleSendAgentMessage handles incoming messages from the user to the agent
func (s *Server) handleSendAgentMessage(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	logger.Logger.Info("Processing agent message",
		zap.String("message_id", msg.Id))

	in := msg.GetSendAgentMessage()
	if in == nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("invalid send agent message"))
	}

	// Get session from database
	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("invalid session ID: %w", err))
	}

	sess, err := s.db.Session().Get(ctx, database.SessionByID(sessionID))
	if err != nil {
		return errdefs.ErrSessionNotFound(fmt.Errorf("session not found: %s", sessionID))
	}

	// Get the host instance associated with this session
	hostInstance, err := s.db.HostInstance().Get(ctx, 
		database.HostInstanceBySessionID(sess.ID),
		database.HostInstanceByStatus(core.HostInstanceStatusOnline))
	if err != nil {
		// If no host instance found for session, send error response
		errorMsg := &websocketv1.Message{
			Id: uuid.Must(uuid.NewV4()).String(),
			Type: &websocketv1.Message_AgentResponse{
				AgentResponse: &websocketv1.AgentResponse{
					SessionId: in.SessionId,
					AgentId:   in.AgentId,
					Type:      websocketv1.AgentResponse_RESPONSE_TYPE_ERROR,
					Content: &websocketv1.AgentResponse_ErrorMessage{
						ErrorMessage: "No host instance available for this session",
					},
				},
			},
		}

		data, _ := proto.Marshal(errorMsg)
		conn.WriteMessage(websocket.BinaryMessage, data)
		
		logger.Logger.Error("Failed to get host instance for session",
			zap.Error(err),
			zap.String("session_id", sess.ID.String()))
		return nil
	}

	// Verify the host is still online
	if err := s.wsManager.PingConnectedHost(hostInstance.ID); err != nil {
		// Host is not responding, send error
		errorMsg := &websocketv1.Message{
			Id: uuid.Must(uuid.NewV4()).String(),
			Type: &websocketv1.Message_AgentResponse{
				AgentResponse: &websocketv1.AgentResponse{
					SessionId: in.SessionId,
					AgentId:   in.AgentId,
					Type:      websocketv1.AgentResponse_RESPONSE_TYPE_ERROR,
					Content: &websocketv1.AgentResponse_ErrorMessage{
						ErrorMessage: "Host instance is not responding",
					},
				},
			},
		}

		data, _ := proto.Marshal(errorMsg)
		conn.WriteMessage(websocket.BinaryMessage, data)
		
		logger.Logger.Error("Host instance not responding",
			zap.String("host_instance_id", hostInstance.ID.String()),
			zap.String("session_id", sess.ID.String()))
		return nil
	}

	// Forward the message to the host via WebSocket manager
	hostMsg := &websocketv1.Message{
		Id: uuid.Must(uuid.NewV4()).String(),
		Type: &websocketv1.Message_SendAgentMessage{
			SendAgentMessage: in,
		},
	}

	if err := s.wsManager.SendToHost(ctx, hostInstance.ID, hostMsg); err != nil {
		// Send error response if forwarding fails
		errorMsg := &websocketv1.Message{
			Id: uuid.Must(uuid.NewV4()).String(),
			Type: &websocketv1.Message_AgentResponse{
				AgentResponse: &websocketv1.AgentResponse{
					SessionId: in.SessionId,
					AgentId:   in.AgentId,
					Type:      websocketv1.AgentResponse_RESPONSE_TYPE_ERROR,
					Content: &websocketv1.AgentResponse_ErrorMessage{
						ErrorMessage: "Failed to forward message to host: " + err.Error(),
					},
				},
			},
		}

		data, _ := proto.Marshal(errorMsg)
		conn.WriteMessage(websocket.BinaryMessage, data)

		logger.Logger.Error("Failed to forward agent message to host",
			zap.Error(err),
			zap.String("host_instance_id", hostInstance.ID.String()))
		return nil
	}

	// The host will send back AgentResponse messages which will be forwarded to the client
	// This is handled in the main WebSocket message handler

	logger.Logger.Info("Agent message forwarded to host",
		zap.String("session_id", sessionID.String()),
		zap.String("agent_id", in.AgentId),
		zap.String("host_instance_id", hostInstance.ID.String()))

	return nil
}

// handleAgentResponse forwards agent responses from host to client
func (s *Server) handleAgentResponse(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	logger.Logger.Debug("Forwarding agent response",
		zap.String("message_id", msg.Id))

	// Get the session ID from the message
	var sessionID string
	if response := msg.GetAgentResponse(); response != nil {
		sessionID = response.SessionId
	} else if complete := msg.GetAgentChatComplete(); complete != nil {
		sessionID = complete.SessionId
	} else {
		return fmt.Errorf("invalid agent response message")
	}

	// Parse session ID
	sessionUUID, err := uuid.FromString(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	// Forward the response to the client
	if err := s.wsManager.SendToClient(ctx, sessionUUID, msg); err != nil {
		return fmt.Errorf("failed to forward agent response to client: %w", err)
	}

	// Log completion for debugging
	if complete := msg.GetAgentChatComplete(); complete != nil {
		logger.Logger.Info("Agent chat completed",
			zap.String("session_id", complete.SessionId),
			zap.String("agent_id", complete.AgentId),
			zap.Int32("total_tokens", complete.TotalTokens),
			zap.Int64("duration_ms", complete.DurationMs))
	}

	return nil
}
