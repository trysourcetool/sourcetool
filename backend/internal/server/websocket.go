package server

import (
	"context"
	"errors"
	"net/http"

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

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 1024 // 512KB
)

func (s *Server) handleInitializeClient(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetInitializeClient()
	if in == nil {
		return errors.New("invalid message")
	}

	pageID, err := uuid.FromString(in.PageId)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	page, err := s.db.Page().Get(ctx, database.PageByID(pageID))
	if err != nil {
		return err
	}

	logger.Logger.Sugar().Infof("Page: %v", page.Name)

	ctxOrg := internal.ContextOrganization(ctx)
	if ctxOrg.ID != page.OrganizationID {
		return errdefs.ErrPermissionDenied(errors.New("organization mismatch"))
	}

	apiKey, err := s.db.APIKey().Get(ctx, database.APIKeyByID(page.APIKeyID))
	if err != nil {
		return err
	}

	env, err := s.db.Environment().Get(ctx, database.EnvironmentByID(apiKey.EnvironmentID))
	if err != nil {
		return err
	}

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

	ctxUser := internal.ContextUser(ctx)

	var sess *core.Session
	var sessionExists bool
	if internal.StringValue(in.SessionId) != "" {
		sessionID, err := uuid.FromString(internal.StringValue(in.SessionId))
		if err != nil {
			return errdefs.ErrSessionNotFound(err)
		}

		sess, err = s.db.Session().Get(ctx, database.SessionByID(sessionID))
		if err != nil {
			return err
		}
		sessionExists = true
	} else {
		sess = &core.Session{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: page.OrganizationID,
			EnvironmentID:  env.ID,
			UserID:         ctxUser.ID,
			Type:           core.SessionTypePage,
		}
		sessionExists = false
	}

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

	if err := s.sendWebSocketMessage(conn, &websocketv1.Message{
		Id: msg.Id,
		Type: &websocketv1.Message_InitializeClientCompleted{
			InitializeClientCompleted: &websocketv1.InitializeClientCompleted{
				SessionId: sess.ID.String(),
			},
		},
	}); err != nil {
		return err
	}

	s.wsManager.SetConnectedClient(sess, conn)

	if err := s.wsManager.SendToHost(ctx, onlineHostInstance.ID, &websocketv1.Message{
		Id: uuid.Must(uuid.NewV4()).String(),
		Type: &websocketv1.Message_InitializeClient{
			InitializeClient: &websocketv1.InitializeClient{
				SessionId: internal.StringPtr(sess.ID.String()),
				PageId:    page.ID.String(),
			},
		},
	}); err != nil {
		s.db.Session().Delete(ctx, sess)
		s.wsManager.DisconnectClient(sess.ID)
		logger.Logger.Sugar().Errorf("Failed to send initialize client message to host: %v", err)
		return err
	}

	return nil
}

func (s *Server) handleRenderWidget(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetRenderWidget()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return err
	}

	_, err = s.db.Session().Get(ctx, database.SessionByID(sessionID))
	if err != nil {
		return err
	}

	if err := s.wsManager.SendToClient(ctx, sessionID, msg); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send render widget message to client: %v", err)
		return err
	}

	return nil
}

func (s *Server) handleRerunPage(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetRerunPage()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return err
	}

	sess, err := s.db.Session().Get(ctx, database.SessionByID(sessionID))
	if err != nil {
		return err
	}

	pageID, err := uuid.FromString(in.PageId)
	if err != nil {
		return err
	}

	page, err := s.db.Page().Get(ctx, database.PageByID(pageID), database.PageBySessionID(sess.ID))
	if err != nil {
		return err
	}

	hostInstance, err := s.db.HostInstance().Get(ctx, database.HostInstanceBySessionID(sess.ID), database.HostInstanceByStatus(core.HostInstanceStatusOnline))
	if err != nil {
		return err
	}

	if err := s.wsManager.SendToHost(ctx, hostInstance.ID, &websocketv1.Message{
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

func (s *Server) handleCloseSession(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetCloseSession()
	if in == nil {
		return errors.New("invalid message")
	}

	// Debug log for tracking CloseSession calls
	logger.Logger.Sugar().Debugf("CloseSession called - sessionId: %s, messageId: %s", in.SessionId, msg.Id)

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return errdefs.ErrAPIKeyNotFound(err)
	}

	sess, err := s.db.Session().Get(ctx, database.SessionByID(sessionID))
	if err != nil {
		return err
	}

	// Check session type and validate accordingly
	switch sess.Type {
	case core.SessionTypePage:
		pages, err := s.db.Page().List(ctx, database.PageByEnvironmentID(sess.EnvironmentID), database.PageBySessionID(sess.ID))
		if err != nil {
			logger.Logger.Sugar().Warnf("Page not found for session %s, continuing with cleanup", sess.ID)
		} else if len(pages) > 0 {
			logger.Logger.Sugar().Debugf("Closing page session: %s (environment has %d pages)", sess.ID, len(pages))
		} else {
			logger.Logger.Sugar().Warnf("No pages found in environment for page session %s", sess.ID)
		}
	case core.SessionTypeAgent:
		agents, err := s.db.Agent().List(ctx, database.AgentByEnvironmentID(sess.EnvironmentID), database.AgentBySessionID(sess.ID))
		if err != nil {
			logger.Logger.Sugar().Warnf("Failed to get agents for session %s: %v", sess.ID, err)
		} else if len(agents) > 0 {
			logger.Logger.Sugar().Debugf("Closing agent session: %s (environment has %d agents)", sess.ID, len(agents))
		} else {
			logger.Logger.Sugar().Warnf("No agents found in environment for agent session %s", sess.ID)
		}
	default:
		// Handle legacy sessions without type (assume page type)
		logger.Logger.Sugar().Debugf("Session %s has no type, treating as page session", sess.ID)
		_, err = s.db.Page().Get(ctx, database.PageByEnvironmentID(sess.EnvironmentID), database.PageBySessionID(sess.ID))
		if err != nil {
			logger.Logger.Sugar().Warnf("Page not found for legacy session %s, continuing with cleanup", sess.ID)
		}
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.Session().Delete(ctx, sess); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	hostInstances, err := s.db.HostInstance().List(ctx, database.HostInstanceBySessionID(sess.ID))
	if err != nil {
		return err
	}

	for _, hostInstance := range hostInstances {
		if err := s.wsManager.SendToHost(ctx, hostInstance.ID, &websocketv1.Message{
			Id: uuid.Must(uuid.NewV4()).String(),
			Type: &websocketv1.Message_CloseSession{
				CloseSession: &websocketv1.CloseSession{
					SessionId: sess.ID.String(),
				},
			},
		}); err != nil {
			logger.Logger.Sugar().Warnf("Failed to send close session message to host %s for session %s: %v", hostInstance.ID, sess.ID, err)
		}
	}

	s.wsManager.DisconnectClient(sess.ID)

	// Debug log for successful session cleanup
	logger.Logger.Sugar().Debugf("CloseSession completed successfully - sessionId: %s", sess.ID.String())

	return nil
}

func (s *Server) handleScriptFinished(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetScriptFinished()
	if in == nil {
		return errors.New("invalid message")
	}

	logger.Logger.Sugar().Debug("Payload: ", in)

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	_, err = s.db.Session().Get(ctx, database.SessionByID(sessionID))
	if err != nil {
		return err
	}

	if err := s.wsManager.SendToClient(ctx, sessionID, msg); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send script finished message to client: %v", err)
	}

	return nil
}

func (s *Server) handleException(ctx context.Context, conn *websocket.Conn, msg *websocketv1.Message) error {
	in := msg.GetException()
	if in == nil {
		return errors.New("invalid message")
	}

	sessionID, err := uuid.FromString(in.SessionId)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	_, err = s.db.Session().Get(ctx, database.SessionByID(sessionID))
	if err != nil {
		return err
	}

	if err := s.wsManager.SendToClient(ctx, sessionID, msg); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send exception message to client: %v", err)
	}

	return nil
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Logger.Sugar().Errorf("Failed to upgrade connection: %v", err)
		return
	}

	conn.SetReadLimit(maxMessageSize)

	ctx := internal.NewBackgroundContext(r.Context())
	done := make(chan struct{})
	defer func() {
		logger.Logger.Info("Closing connection")
		close(done)
		conn.Close()
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			logger.Logger.Sugar().Debugf("Connection closed: %v", err)
			break
		}

		var msg websocketv1.Message
		if err := proto.Unmarshal(data, &msg); err != nil {
			logger.Logger.Sugar().Errorf("Failed to unmarshal message: %v", err)
			break
		}

		switch msg.Type.(type) {
		case *websocketv1.Message_InitializeHost:
			instanceID := r.Header.Get("X-Instance-Id")
			if err := s.handleInitializeHost(ctx, conn, instanceID, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_InitializeClient:
			if err := s.handleInitializeClient(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_RenderWidget:
			if err := s.handleRenderWidget(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_RerunPage:
			if err := s.handleRerunPage(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_CloseSession:
			if err := s.handleCloseSession(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_ScriptFinished:
			if err := s.handleScriptFinished(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_Exception:
			if err := s.handleException(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		// Agent chat message handlers
		case *websocketv1.Message_InitializeAgentChat:
			if err := s.handleInitializeAgentChat(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_InitializeAgentChatCompleted:
			// This message is sent from backend to client, not from SDK
			// It should not be received here, log and skip
			logger.Logger.Debug("Received InitializeAgentChatCompleted message - skipping",
				zap.String("message_id", msg.Id))
			continue
		case *websocketv1.Message_SendAgentMessage:
			if err := s.handleSendAgentMessage(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_AgentResponse:
			if err := s.handleAgentResponse(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_AgentChatComplete:
			// This is typically sent from the host to the client
			if err := s.handleAgentResponse(ctx, conn, &msg); err != nil {
				s.sendErrWebSocketMessage(ctx, conn, msg.Id, err)
				continue
			}
		default:
			logger.Logger.Sugar().Errorf("Unknown method: %s", msg.Type)
			continue
		}
	}
}

func (s *Server) handleInitializeHostBase(ctx context.Context, conn *websocket.Conn, instanceID string, msg *websocketv1.Message) (*core.HostInstance, bool, *core.APIKey, []*core.Page, []*core.Page, []*core.Page, []*core.Agent, []*core.Agent, []*core.Agent, []*core.AgentTool, []*core.AgentTool, []*core.AgentTool, error) {
	in := msg.GetInitializeHost()
	if in == nil {
		return nil, false, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.New("invalid message")
	}

	hashedAPIKey := core.HashAPIKey(in.ApiKey)
	apikey, err := s.db.APIKey().Get(ctx, database.APIKeyByKeyHash(hashedAPIKey))
	if err != nil {
		return nil, false, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	hostInstanceID, err := uuid.FromString(instanceID)
	if err != nil {
		return nil, false, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errdefs.ErrInvalidArgument(err)
	}

	hostInstance, err := s.db.HostInstance().Get(ctx, database.HostInstanceByID(hostInstanceID))
	if err != nil && !errdefs.IsHostInstanceNotFound(err) {
		return nil, false, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
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

	existingPages, err := s.db.Page().List(ctx, database.PageByAPIKeyID(apikey.ID))
	if err != nil {
		return nil, false, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	existingPageMap := make(map[string]*core.Page)
	for _, p := range existingPages {
		existingPageMap[p.ID.String()] = p
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
				return nil, false, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
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

	// Process agents similar to pages
	existingAgents, err := s.db.Agent().List(ctx, database.AgentByAPIKeyID(apikey.ID))
	if err != nil {
		return nil, false, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	existingAgentMap := make(map[string]*core.Agent)
	for _, a := range existingAgents {
		existingAgentMap[a.ID.String()] = a
	}

	requestAgentIDs := make(map[string]struct{})
	for _, a := range in.Agents {
		requestAgentIDs[a.Id] = struct{}{}
	}

	insertAgents := make([]*core.Agent, 0)
	updateAgents := make([]*core.Agent, 0)
	deleteAgents := make([]*core.Agent, 0)
	for _, reqAgent := range in.Agents {
		if existingAgent, ok := existingAgentMap[reqAgent.Id]; ok {
			existingAgent.Name = reqAgent.Name
			existingAgent.Description = reqAgent.Description
			existingAgent.Instructions = reqAgent.Instructions
			existingAgent.Model = reqAgent.Model
			updateAgents = append(updateAgents, existingAgent)
		} else {
			parsedAgentID, err := uuid.FromString(reqAgent.Id)
			if err != nil {
				return nil, false, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			newAgent := &core.Agent{
				ID:             parsedAgentID,
				OrganizationID: apikey.OrganizationID,
				EnvironmentID:  apikey.EnvironmentID,
				APIKeyID:       apikey.ID,
				Name:           reqAgent.Name,
				Description:    reqAgent.Description,
				Instructions:   reqAgent.Instructions,
				Model:          reqAgent.Model,
			}
			insertAgents = append(insertAgents, newAgent)
		}
	}

	for _, existingAgent := range existingAgents {
		if _, exists := requestAgentIDs[existingAgent.ID.String()]; !exists {
			deleteAgents = append(deleteAgents, existingAgent)
		}
	}

	// Process agent tools
	insertAgentTools := make([]*core.AgentTool, 0)
	updateAgentTools := make([]*core.AgentTool, 0)
	deleteAgentTools := make([]*core.AgentTool, 0)

	// Build a map of all agents (existing + new) for tool processing
	allAgentsMap := make(map[string]*core.Agent)
	for _, agent := range existingAgents {
		allAgentsMap[agent.ID.String()] = agent
	}
	for _, agent := range insertAgents {
		allAgentsMap[agent.ID.String()] = agent
	}

	// Process tools for each agent in the request
	for _, reqAgent := range in.Agents {
		agentID, err := uuid.FromString(reqAgent.Id)
		if err != nil {
			return nil, false, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}

		// Get existing tools for this agent
		existingTools, err := s.db.AgentTool().List(ctx, database.AgentToolByAgentID(agentID))
		if err != nil {
			// If agent doesn't exist yet (new agent), there won't be any tools
			// Continue with empty list
			existingTools = []*core.AgentTool{}
		}

		// Create a map of existing tools by name
		existingToolMap := make(map[string]*core.AgentTool)
		for _, tool := range existingTools {
			existingToolMap[tool.Name] = tool
		}

		// Track which tools are in the request
		requestToolNames := make(map[string]struct{})
		for _, reqTool := range reqAgent.Tools {
			requestToolNames[reqTool.Name] = struct{}{}
		}

		// Process tools from request
		for _, reqTool := range reqAgent.Tools {
			if existingTool, ok := existingToolMap[reqTool.Name]; ok {
				// Update existing tool
				existingTool.Description = reqTool.Description
				updateAgentTools = append(updateAgentTools, existingTool)
			} else {
				// Insert new tool
				newTool := &core.AgentTool{
					ID:          uuid.Must(uuid.NewV4()),
					AgentID:     agentID,
					Name:        reqTool.Name,
					Description: reqTool.Description,
				}
				insertAgentTools = append(insertAgentTools, newTool)
			}
		}

		// Mark tools for deletion
		for _, existingTool := range existingTools {
			if _, exists := requestToolNames[existingTool.Name]; !exists {
				deleteAgentTools = append(deleteAgentTools, existingTool)
			}
		}
	}

	// Also check for tools that need to be deleted for agents that are being deleted
	for _, agent := range deleteAgents {
		existingTools, err := s.db.AgentTool().List(ctx, database.AgentToolByAgentID(agent.ID))
		if err != nil {
			// If error occurs, skip this agent's tools
			continue
		}
		deleteAgentTools = append(deleteAgentTools, existingTools...)
	}

	return hostInstance, hostExists, apikey, insertPages, updatePages, deletePages, insertAgents, updateAgents, deleteAgents, insertAgentTools, updateAgentTools, deleteAgentTools, nil
}
