package sourcetool

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"

	"github.com/trysourcetool/sourcetool-go/internal/errdefs"
	"github.com/trysourcetool/sourcetool-go/internal/logger"
	agentv1 "github.com/trysourcetool/sourcetool-go/internal/pb/agent/v1"
	exceptionv1 "github.com/trysourcetool/sourcetool-go/internal/pb/exception/v1"
	pagev1 "github.com/trysourcetool/sourcetool-go/internal/pb/page/v1"
	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/pb/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/pb/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/ptrconv"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

type runtime struct {
	wsClient       websocket.Client
	sessionManager *session.SessionManager
	pageManager    *pageManager
	agentManager   *agentManager
}

func startRuntime(apiKey, endpoint string, pages map[uuid.UUID]*page, agents map[uuid.UUID]*Agent) (*runtime, error) {
	r := &runtime{
		sessionManager: session.NewSessionManager(),
		pageManager:    newPageManager(pages),
		agentManager:   newAgentManager(agents),
	}

	wsClient, err := websocket.NewClient(websocket.Config{
		URL:            endpoint,
		APIKey:         apiKey,
		InstanceID:     uuid.Must(uuid.NewV4()),
		PingInterval:   1 * time.Second,
		ReconnectDelay: 1 * time.Second,
		OnReconnecting: func() {
			logger.Log.Info("Reconnecting...")
		},
		OnReconnected: func() {
			logger.Log.Info("Reconnected!")
			r.sendInitializeHost(apiKey, pages, agents)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create websocket client: %v", err)
	}

	r.wsClient = wsClient
	wsClient.RegisterHandler(func(msg *websocketv1.Message) error {
		switch t := msg.Type.(type) {
		case *websocketv1.Message_InitializeClient:
			if err := r.handleInitializeClient(t.InitializeClient); err != nil {
				r.sendException(msg.Id, *t.InitializeClient.SessionId, err)
			}
			return nil
		case *websocketv1.Message_RerunPage:
			if err := r.handleRerunPage(t.RerunPage); err != nil {
				r.sendException(msg.Id, t.RerunPage.SessionId, err)
			}
			return nil
		case *websocketv1.Message_CloseSession:
			if err := r.handleCloseSession(t.CloseSession); err != nil {
				r.sendException(msg.Id, t.CloseSession.SessionId, err)
			}
			return nil
		case *websocketv1.Message_InitializeAgentChat:
			if err := r.handleInitializeAgentChat(t.InitializeAgentChat); err != nil {
				r.sendAgentError(ptrconv.StringValue(t.InitializeAgentChat.SessionId), t.InitializeAgentChat.AgentId, err)
			}
			return nil
		case *websocketv1.Message_SendAgentMessage:
			if err := r.handleSendAgentMessage(t.SendAgentMessage); err != nil {
				r.sendAgentError(t.SendAgentMessage.SessionId, t.SendAgentMessage.AgentId, err)
			}
			return nil
		case *websocketv1.Message_AgentResponse:
			// AgentResponse is sent from the server to the client through the SDK
			// SDK doesn't need to process it, just ignore
			return nil
		case *websocketv1.Message_AgentChatComplete:
			// AgentChatComplete is sent from the server to the client through the SDK
			// SDK doesn't need to process it, just ignore
			return nil
		case *websocketv1.Message_InitializeHostCompleted:
			// InitializeHostCompleted is a response from server confirming host initialization
			// Log for debugging purposes
			logger.Log.Info("Host initialization completed",
				zap.String("host_instance_id", t.InitializeHostCompleted.HostInstanceId))
			return nil
		default:
			return fmt.Errorf("unknown message type: %T", t)
		}
	})

	r.sendInitializeHost(apiKey, pages, agents)

	return r, nil
}

func (r *runtime) sendInitializeHost(apiKey string, pages map[uuid.UUID]*page, agents map[uuid.UUID]*Agent) {
	pagesPayload := make([]*pagev1.Page, 0, len(pages))
	for _, page := range pages {
		pagesPayload = append(pagesPayload, &pagev1.Page{
			Id:     page.id.String(),
			Name:   page.name,
			Route:  page.route,
			Path:   convertPathToInt32Slice(page.path),
			Groups: page.accessGroups,
		})
	}

	agentsPayload := make([]*agentv1.Agent, 0, len(agents))
	for _, agent := range agents {
		toolsPayload := make([]*agentv1.Tool, 0, len(agent.Tools))
		for _, tool := range agent.Tools {
			toolsPayload = append(toolsPayload, &agentv1.Tool{
				Name:        tool.GetName(),
				Description: tool.GetDescription(),
			})
		}

		agentsPayload = append(agentsPayload, &agentv1.Agent{
			Id:           agent.id.String(),
			Name:         agent.Name,
			Description:  agent.Description,
			Instructions: agent.Instructions,
			Model:        agent.Model.ID(),
			Groups:       agent.accessGroups,
			Tools:        toolsPayload,
		})
	}

	msg := &websocketv1.InitializeHost{
		ApiKey:     apiKey,
		SdkName:    "sourcetool-go",
		SdkVersion: "0.1.12",
		Pages:      pagesPayload,
		Agents:     agentsPayload,
	}

	resp, err := r.wsClient.EnqueueWithResponse(uuid.Must(uuid.NewV4()).String(), msg)
	if err != nil {
		logger.Log.Fatal("failed to send initialize host message", zap.Error(err))
	}

	if e := resp.GetException(); e != nil {
		logger.Log.Fatal("initialize host message failed", zap.String("message", e.Message))
	}

	logger.Log.Info("initialize host message sent", zap.Any("response", resp))
}

func (r *runtime) handleInitializeClient(msg *websocketv1.InitializeClient) error {
	if msg.SessionId == nil {
		return errdefs.ErrInvalidParameter(errors.New("session id is required"))
	}
	sessionID, err := uuid.FromString(ptrconv.StringValue(msg.SessionId))
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}
	pageID, err := uuid.FromString(msg.PageId)
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}

	session := session.New(sessionID, pageID)
	r.sessionManager.SetSession(session)

	page := r.pageManager.getPage(pageID)
	if page == nil {
		return errdefs.ErrInternal(fmt.Errorf("page not found: %s", pageID))
	}

	ui := &uiBuilder{
		context: context.Background(),
		runtime: r,
		session: session,
		page:    page,
		cursor:  newCursor(),
	}

	if err := page.run(ui); err != nil {
		r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.ScriptFinished{
			SessionId: sessionID.String(),
			Status:    websocketv1.ScriptFinished_STATUS_FAILURE,
		})

		return errdefs.ErrRunPage(err)
	}

	r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.ScriptFinished{
		SessionId: sessionID.String(),
		Status:    websocketv1.ScriptFinished_STATUS_SUCCESS,
	})

	return nil
}

func (r *runtime) handleRerunPage(msg *websocketv1.RerunPage) error {
	sessionID, err := uuid.FromString(msg.SessionId)
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}
	sess := r.sessionManager.GetSession(sessionID)
	if sess == nil {
		return errdefs.ErrSessionNotFound(fmt.Errorf("session not found: %s", sessionID))
	}

	pageID, err := uuid.FromString(msg.PageId)
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}
	page := r.pageManager.getPage(pageID)
	if page == nil {
		return errdefs.ErrPageNotFound(fmt.Errorf("page not found: %s", pageID))
	}

	if sess.PageID != pageID {
		sess.State.ResetStates()
	}

	newWidgetStates := make(map[uuid.UUID]session.WidgetState)
	for _, widget := range msg.States {
		id, err := uuid.FromString(widget.Id)
		if err != nil {
			return errdefs.ErrInvalidParameter(err)
		}
		switch t := widget.Type.(type) {
		case *widgetv1.Widget_TextInput:
			newWidgetStates[id] = convertTextInputProtoToState(id, t.TextInput)
		case *widgetv1.Widget_NumberInput:
			newWidgetStates[id] = convertNumberInputProtoToState(id, t.NumberInput)
		case *widgetv1.Widget_DateInput:
			state, err := convertDateInputProtoToState(id, t.DateInput, time.Local)
			if err != nil {
				return errdefs.ErrInvalidParameter(err)
			}
			newWidgetStates[id] = state
		case *widgetv1.Widget_DateTimeInput:
			state, err := convertDateTimeInputProtoToState(id, t.DateTimeInput, time.Local)
			if err != nil {
				return errdefs.ErrInvalidParameter(err)
			}
			newWidgetStates[id] = state
		case *widgetv1.Widget_TimeInput:
			state, err := convertTimeInputProtoToState(id, t.TimeInput, time.Local)
			if err != nil {
				return errdefs.ErrInvalidParameter(err)
			}
			newWidgetStates[id] = state
		case *widgetv1.Widget_Form:
			newWidgetStates[id] = convertFormProtoToState(id, t.Form)
		case *widgetv1.Widget_Button:
			newWidgetStates[id] = convertButtonProtoToState(id, t.Button)
		case *widgetv1.Widget_Markdown:
			newWidgetStates[id] = convertMarkdownProtoToState(id, t.Markdown)
		case *widgetv1.Widget_Columns:
			newWidgetStates[id] = convertColumnsProtoToState(id, t.Columns)
		case *widgetv1.Widget_ColumnItem:
			newWidgetStates[id] = convertColumnItemProtoToState(id, t.ColumnItem)
		case *widgetv1.Widget_Table:
			newWidgetStates[id] = convertTableProtoToState(id, t.Table)
		case *widgetv1.Widget_Selectbox:
			newWidgetStates[id] = convertSelectboxProtoToState(id, t.Selectbox)
		case *widgetv1.Widget_MultiSelect:
			newWidgetStates[id] = convertMultiSelectProtoToState(id, t.MultiSelect)
		case *widgetv1.Widget_Checkbox:
			newWidgetStates[id] = convertCheckboxProtoToState(id, t.Checkbox)
		case *widgetv1.Widget_CheckboxGroup:
			newWidgetStates[id] = convertCheckboxGroupProtoToState(id, t.CheckboxGroup)
		case *widgetv1.Widget_Radio:
			newWidgetStates[id] = convertRadioProtoToState(id, t.Radio)
		case *widgetv1.Widget_TextArea:
			newWidgetStates[id] = convertTextAreaProtoToState(id, t.TextArea)
		default:
			return errdefs.ErrInvalidParameter(fmt.Errorf("unknown widget type: %T", t))
		}
	}

	sess.State.SetStates(newWidgetStates)

	ui := &uiBuilder{
		context: context.Background(),
		runtime: r,
		session: sess,
		page:    page,
		cursor:  newCursor(),
	}

	if err := page.run(ui); err != nil {
		r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.ScriptFinished{
			SessionId: sessionID.String(),
			Status:    websocketv1.ScriptFinished_STATUS_FAILURE,
		})

		return errdefs.ErrRunPage(err)
	}

	r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.ScriptFinished{
		SessionId: sessionID.String(),
		Status:    websocketv1.ScriptFinished_STATUS_SUCCESS,
	})

	sess.State.ResetButtons()

	return nil
}

func (r *runtime) handleCloseSession(msg *websocketv1.CloseSession) error {
	sessionID, err := uuid.FromString(msg.SessionId)
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}

	r.sessionManager.DisconnectSession(sessionID)

	return nil
}

func (r *runtime) sendException(id, sessionID string, err error) {
	e, ok := err.(*errdefs.Error)
	if !ok {
		v := errdefs.ErrInternal(err)
		e = v.(*errdefs.Error)
	}

	exception := &exceptionv1.Exception{
		Title:      e.Title,
		Message:    e.Message,
		StackTrace: e.StackTrace(),
		SessionId:  sessionID,
	}

	r.wsClient.Enqueue(id, exception)
}

func (r *runtime) Close() error {
	err := r.wsClient.Close()
	r.wsClient = nil
	return err
}

// Agent handling methods

// AgentSession represents an agent chat session in the SDK.
type AgentSession struct {
	ID           uuid.UUID
	AgentID      uuid.UUID
	Conversation []Message // Using SDK's Message type
	StartedAt    time.Time
	LastActivity time.Time
}

// agentSessions stores active agent sessions.
var (
	agentSessions   = make(map[uuid.UUID]*AgentSession)
	agentSessionsMu sync.RWMutex
)

// handleInitializeAgentChat handles agent chat initialization from backend.
func (r *runtime) handleInitializeAgentChat(msg *websocketv1.InitializeAgentChat) error {
	if msg.SessionId == nil {
		return errdefs.ErrInvalidParameter(errors.New("session id is required"))
	}

	sessionID, err := uuid.FromString(ptrconv.StringValue(msg.SessionId))
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}

	agentID, err := uuid.FromString(msg.AgentId)
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}

	// Get the agent from the agent manager
	agent := r.agentManager.getAgent(agentID)
	if agent == nil {
		return errdefs.ErrInternal(fmt.Errorf("agent not found: %s", agentID))
	}

	// Create a new agent session
	agentSession := &AgentSession{
		ID:           sessionID,
		AgentID:      agentID,
		Conversation: []Message{},
		StartedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	// Store the agent session
	agentSessionsMu.Lock()
	agentSessions[sessionID] = agentSession
	agentSessionsMu.Unlock()

	// Send completion message back to backend
	r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.InitializeAgentChatCompleted{
		SessionId: sessionID.String(),
		AgentId:   agentID.String(),
	})

	logger.Log.Info("Agent chat session initialized",
		zap.String("session_id", sessionID.String()),
		zap.String("agent_id", agentID.String()))

	return nil
}

func (r *runtime) handleSendAgentMessage(msg *websocketv1.SendAgentMessage) error {
	agentID, err := uuid.FromString(msg.AgentId)
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}

	agent := r.agentManager.getAgent(agentID)
	if agent == nil {
		return errdefs.ErrInternal(fmt.Errorf("agent not found: %s", agentID))
	}

	// Get or create the agent session
	sessionID, err := uuid.FromString(msg.SessionId)
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}

	agentSessionsMu.RLock()
	agentSession, exists := agentSessions[sessionID]
	agentSessionsMu.RUnlock()

	if !exists {
		// If session doesn't exist, create a new one (backward compatibility)
		agentSession = &AgentSession{
			ID:           sessionID,
			AgentID:      agentID,
			Conversation: []Message{},
			StartedAt:    time.Now(),
			LastActivity: time.Now(),
		}
		agentSessionsMu.Lock()
		agentSessions[sessionID] = agentSession
		agentSessionsMu.Unlock()
	}

	// Update last activity
	agentSession.LastActivity = time.Now()

	// Add the new user message to conversation
	agentSession.Conversation = append(agentSession.Conversation, Message{
		Role:    "user",
		Content: msg.Message,
	})

	// Also include any conversation history sent from backend
	conversation := agentSession.Conversation
	if len(msg.ConversationHistory) > 0 {
		// If backend sends conversation history, use it instead
		conversation = make([]Message, len(msg.ConversationHistory))
		for i, chatMsg := range msg.ConversationHistory {
			var role string
			switch chatMsg.Role {
			case websocketv1.ChatMessage_ROLE_USER:
				role = "user"
			case websocketv1.ChatMessage_ROLE_ASSISTANT:
				role = "assistant"
			case websocketv1.ChatMessage_ROLE_SYSTEM:
				role = "system"
			case websocketv1.ChatMessage_ROLE_TOOL:
				role = "tool"
			default:
				role = "user"
			}

			toolCalls := make([]ToolCall, len(chatMsg.ToolCalls))
			for j, tc := range chatMsg.ToolCalls {
				toolCalls[j] = ToolCall{
					ID:     tc.Id,
					Name:   tc.Name,
					Params: make(map[string]interface{}),
				}
				// Parse JSON arguments if needed
				// TODO: Add proper JSON parsing
			}

			conversation[i] = Message{
				Role:      role,
				Content:   chatMsg.Content,
				ToolCalls: toolCalls,
			}
		}
		// Update the session's conversation
		agentSession.Conversation = conversation
	}

	// Create context for agent execution
	ctx := &Context{
		Context:      context.Background(),
		Message:      msg.Message,
		Conversation: conversation,
		SessionID:    msg.SessionId,
		tools:        agent.toolsMap,
	}

	// Start streaming response
	go r.streamAgentResponse(ctx, agent, msg.SessionId, msg.AgentId)

	return nil
}

func (r *runtime) streamAgentResponse(ctx *Context, agent *Agent, sessionID, agentID string) {
	defer func() {
		if rec := recover(); rec != nil {
			logger.Log.Error("Agent execution panic",
				zap.String("session_id", sessionID),
				zap.String("agent_id", agentID),
				zap.Any("panic", rec))
			r.sendAgentError(sessionID, agentID, fmt.Errorf("agent execution panic: %v", rec))
		}
	}()

	logger.Log.Info("Starting agent streaming response",
		zap.String("session_id", sessionID),
		zap.String("agent_id", agentID),
		zap.String("message", ctx.Message))

	// Buffer to accumulate the full response for session storage
	var fullResponse strings.Builder
	chunkCount := 0

	// Create tool notification callback
	agent.SetToolNotifier(func(toolID, toolName, parameters string) {
		r.sendToolCallStart(sessionID, agentID, toolID, toolName, parameters)
	}, func(toolID, result string, duration int64) {
		r.sendToolCallComplete(sessionID, agentID, toolID, result, duration)
	}, func(toolID, errorMessage string, duration int64) {
		r.sendToolCallError(sessionID, agentID, toolID, errorMessage, duration)
	})

	// Stream response using the agent's streaming capability
	err := agent.Stream(ctx.Context, ctx.Message, func(chunk string) error {
		chunkCount++
		logger.Log.Debug("Streaming chunk received",
			zap.String("session_id", sessionID),
			zap.Int("chunk_count", chunkCount),
			zap.Int("chunk_length", len(chunk)),
			zap.String("chunk_preview", func() string {
				if len(chunk) > 50 {
					return chunk[:50] + "..."
				}
				return chunk
			}()))

		// Accumulate full response
		fullResponse.WriteString(chunk)

		// Send each chunk immediately to the client
		r.sendAgentTextChunk(sessionID, agentID, chunk)

		return nil
	}, WithConversation(ctx.Conversation))
	if err != nil {
		logger.Log.Error("Agent streaming failed",
			zap.String("session_id", sessionID),
			zap.String("agent_id", agentID),
			zap.Error(err))
		r.sendAgentError(sessionID, agentID, err)
		return
	}

	// Get the complete response text
	completeResponse := fullResponse.String()

	logger.Log.Info("Agent streaming completed",
		zap.String("session_id", sessionID),
		zap.String("agent_id", agentID),
		zap.Int("total_chunks", chunkCount),
		zap.Int("total_length", len(completeResponse)))

	// Update the agent session with the assistant's complete response
	sessionUUID, _ := uuid.FromString(sessionID)
	agentSessionsMu.RLock()
	if agentSession, exists := agentSessions[sessionUUID]; exists {
		agentSession.Conversation = append(agentSession.Conversation, Message{
			Role:    "assistant",
			Content: completeResponse,
		})
		agentSession.LastActivity = time.Now()
		logger.Log.Debug("Updated agent session conversation",
			zap.String("session_id", sessionID),
			zap.Int("conversation_length", len(agentSession.Conversation)))
	}
	agentSessionsMu.RUnlock()

	// Send completion message
	logger.Log.Info("Sending agent completion message",
		zap.String("session_id", sessionID),
		zap.String("agent_id", agentID))
	r.sendAgentComplete(sessionID, agentID, completeResponse)
}

func (r *runtime) sendAgentTextChunk(sessionID, agentID, content string) {
	msg := &websocketv1.AgentResponse{
		SessionId: sessionID,
		AgentId:   agentID,
		Type:      websocketv1.AgentResponse_RESPONSE_TYPE_TEXT_CHUNK,
		Content: &websocketv1.AgentResponse_TextChunk{
			TextChunk: content,
		},
	}

	r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), msg)
}

func (r *runtime) sendAgentComplete(sessionID, agentID, finalMessage string) {
	msg := &websocketv1.AgentChatComplete{
		SessionId:    sessionID,
		AgentId:      agentID,
		FinalMessage: finalMessage,
		TotalTokens:  0, // TODO: Add token counting
		DurationMs:   0, // TODO: Add duration tracking
	}

	r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), msg)
}

func (r *runtime) sendAgentError(sessionID, agentID string, err error) {
	msg := &websocketv1.AgentResponse{
		SessionId: sessionID,
		AgentId:   agentID,
		Type:      websocketv1.AgentResponse_RESPONSE_TYPE_ERROR,
		Content: &websocketv1.AgentResponse_ErrorMessage{
			ErrorMessage: err.Error(),
		},
	}

	r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), msg)
}

func (r *runtime) sendToolCallStart(sessionID, agentID, toolID, toolName, parameters string) {
	msg := &websocketv1.AgentResponse{
		SessionId: sessionID,
		AgentId:   agentID,
		Type:      websocketv1.AgentResponse_RESPONSE_TYPE_TOOL_CALL_START,
		Content: &websocketv1.AgentResponse_ToolCallInfo{
			ToolCallInfo: &websocketv1.ToolCallInfo{
				ToolId:     toolID,
				ToolName:   toolName,
				Parameters: parameters,
			},
		},
	}

	r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), msg)
}

func (r *runtime) sendToolCallComplete(sessionID, agentID, toolID, result string, duration int64) {
	msg := &websocketv1.AgentResponse{
		SessionId: sessionID,
		AgentId:   agentID,
		Type:      websocketv1.AgentResponse_RESPONSE_TYPE_TOOL_CALL_COMPLETE,
		Content: &websocketv1.AgentResponse_ToolCallInfo{
			ToolCallInfo: &websocketv1.ToolCallInfo{
				ToolId:     toolID,
				Result:     result,
				DurationMs: duration,
			},
		},
	}

	r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), msg)
}

func (r *runtime) sendToolCallError(sessionID, agentID, toolID, errorMessage string, duration int64) {
	msg := &websocketv1.AgentResponse{
		SessionId: sessionID,
		AgentId:   agentID,
		Type:      websocketv1.AgentResponse_RESPONSE_TYPE_TOOL_CALL_ERROR,
		Content: &websocketv1.AgentResponse_ToolCallInfo{
			ToolCallInfo: &websocketv1.ToolCallInfo{
				ToolId:       toolID,
				ErrorMessage: errorMessage,
				DurationMs:   duration,
			},
		},
	}

	r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), msg)
}
