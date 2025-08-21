package sourcetool

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/agent"
	"github.com/trysourcetool/sourcetool-go/agent/llm"
	"github.com/trysourcetool/sourcetool-go/agent/models"
)

// Agent defines an AI agent with its configuration and execution capabilities.
type Agent struct {
	Name         string
	Description  string
	Instructions string
	Model        models.Model
	Tools        []agent.Tool

	// Optional configuration
	Timeout  int // seconds
	MaxSteps int // Maximum number of tool execution steps

	// Internal fields (managed automatically)
	id           uuid.UUID
	route        string
	toolsMap     map[string]agent.Tool
	accessGroups []string
	llmClient    *llm.Client
}

// GenerateRequest represents a request to generate a response.
type GenerateRequest struct {
	Message      string
	Context      context.Context
	Conversation []Message
	User         *User
	SessionID    string
}

// StreamCallback is called for each streaming chunk.
type StreamCallback func(chunk string) error

// GenerateOption configures generation requests.
type GenerateOption func(*GenerateRequest)

// WithConversation adds conversation history to the request.
func WithConversation(conversation []Message) GenerateOption {
	return func(req *GenerateRequest) {
		req.Conversation = conversation
	}
}

// WithUser adds user context to the request.
func WithUser(user *User) GenerateOption {
	return func(req *GenerateRequest) {
		req.User = user
	}
}

// WithSessionID adds session ID to the request.
func WithSessionID(sessionID string) GenerateOption {
	return func(req *GenerateRequest) {
		req.SessionID = sessionID
	}
}

// Message represents a conversation message.
type Message struct {
	Role      string                 `json:"role"` // "user", "assistant", "system", "tool"
	Content   string                 `json:"content"`
	ToolCalls []ToolCall             `json:"tool_calls,omitempty"`
	ToolID    string                 `json:"tool_id,omitempty"` // For tool response messages
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// User represents the current user context.
type User struct {
	ID       string
	Email    string
	Name     string
	Groups   []string
	Metadata map[string]interface{}
}

// Response represents the agent's response.
type Response struct {
	Message          string                 `json:"message"`
	ToolCalls        []ToolCall             `json:"tool_calls,omitempty"`
	SuggestedActions []Action               `json:"suggested_actions,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// ToolCall represents a tool invocation request or result.
type ToolCall struct {
	ID     string                 `json:"id"`
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params,omitempty"`
	Result interface{}            `json:"result,omitempty"`
	Error  string                 `json:"error,omitempty"`
}

// Action represents a suggested action for the user.
type Action struct {
	Label       string                 `json:"label"`
	Tool        string                 `json:"tool,omitempty"`
	Params      map[string]interface{} `json:"params,omitempty"`
	Description string                 `json:"description,omitempty"`
}

// Context provides context for agent execution.
type Context struct {
	Context      context.Context
	Message      string
	Conversation []Message
	User         *User
	SessionID    string

	// Tool execution
	tools map[string]agent.Tool
}

// ExecuteTool executes a tool by name with the given parameters.
func (ac *Context) ExecuteTool(name string, params interface{}) (interface{}, error) {
	tool, exists := ac.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}

	// Convert params to JSON
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters: %w", err)
	}

	return tool.Execute(ac.Context, jsonParams)
}

// ExecuteToolJSON executes a tool with raw JSON parameters.
func (ac *Context) ExecuteToolJSON(name string, params json.RawMessage) (interface{}, error) {
	tool, exists := ac.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}

	return tool.Execute(ac.Context, params)
}

// GetAvailableTools returns information about all available tools.
func (ac *Context) GetAvailableTools() []agent.ToolInfo {
	tools := make([]agent.ToolInfo, 0, len(ac.tools))
	for _, tool := range ac.tools {
		schema, _ := tool.GetSchema()
		tools = append(tools, agent.ToolInfo{
			Name:        tool.GetName(),
			Description: tool.GetDescription(),
			Schema:      schema,
		})
	}
	return tools
}

// GetTool returns a specific tool by name.
func (ac *Context) GetTool(name string) (agent.Tool, bool) {
	tool, exists := ac.tools[name]
	return tool, exists
}

// HasTool checks if a tool exists.
func (ac *Context) HasTool(name string) bool {
	_, exists := ac.tools[name]
	return exists
}

// SetTools sets the available tools (for testing).
func (ac *Context) SetTools(tools map[string]agent.Tool) {
	ac.tools = tools
}

// AddMessage adds a message to the conversation history.
func (ac *Context) AddMessage(role, content string) {
	ac.Conversation = append(ac.Conversation, Message{
		Role:    role,
		Content: content,
	})
}

// AddToolMessage adds a tool response message to the conversation.
func (ac *Context) AddToolMessage(toolID, content string) {
	ac.Conversation = append(ac.Conversation, Message{
		Role:    "tool",
		Content: content,
		ToolID:  toolID,
	})
}

// GetLastMessage returns the most recent message in the conversation.
func (ac *Context) GetLastMessage() *Message {
	if len(ac.Conversation) == 0 {
		return nil
	}
	return &ac.Conversation[len(ac.Conversation)-1]
}

// GetMessageHistory returns messages with an optional limit.
func (ac *Context) GetMessageHistory(limit int) []Message {
	if limit <= 0 || limit > len(ac.Conversation) {
		return ac.Conversation
	}
	return ac.Conversation[len(ac.Conversation)-limit:]
}

// Validate validates the agent.
func (a *Agent) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("agent name is required")
	}
	if a.Instructions == "" {
		return fmt.Errorf("agent instructions are required")
	}
	if a.Model == nil {
		return fmt.Errorf("agent model is required")
	}

	// Set defaults
	if a.MaxSteps == 0 {
		a.MaxSteps = 10 // Default max steps to prevent infinite loops
	}

	return nil
}

// newAgent creates and initializes an agent with internal fields.
func newAgent(agentDef *Agent, id uuid.UUID, route string, accessGroups []string) *Agent {
	// Create tools map from agent tools array
	toolsMap := make(map[string]agent.Tool)
	for _, tool := range agentDef.Tools {
		toolsMap[tool.GetName()] = tool
	}

	// Copy and initialize the agent
	agent := &Agent{
		Name:         agentDef.Name,
		Description:  agentDef.Description,
		Instructions: agentDef.Instructions,
		Model:        agentDef.Model,
		Tools:        agentDef.Tools,
		Timeout:      agentDef.Timeout,
		MaxSteps:     agentDef.MaxSteps,
		id:           id,
		route:        route,
		toolsMap:     toolsMap,
		accessGroups: accessGroups,
		llmClient:    llm.NewClient(30 * time.Second),
	}

	// Set defaults
	if agent.MaxSteps == 0 {
		agent.MaxSteps = 10
	}

	return agent
}

// Agent methods

// GetID returns the agent's ID.
func (a *Agent) GetID() uuid.UUID {
	return a.id
}

// GetName returns the agent's name (convenience method).
func (a *Agent) GetName() string {
	return a.Name
}

// GetDescription returns the agent's description (convenience method).
func (a *Agent) GetDescription() string {
	return a.Description
}

// hasAccess checks if user has access to this agent based on groups.
// TODO: Will be used for access control implementation.
//
//nolint:unused
func (a *Agent) hasAccess(userGroups []string) bool {
	if len(a.accessGroups) == 0 {
		return true
	}

	for _, requiredGroup := range a.accessGroups {
		if slices.Contains(userGroups, requiredGroup) {
			return true
		}
	}
	return false
}

// Generate creates a response using the agent's configuration.
func (a *Agent) Generate(ctx context.Context, message string, options ...GenerateOption) (*Response, error) {
	req := &GenerateRequest{
		Message: message,
		Context: ctx,
	}

	// Apply options
	for _, opt := range options {
		opt(req)
	}

	return a.generateResponse(req)
}

// Stream creates a streaming response using the agent's configuration.
func (a *Agent) Stream(ctx context.Context, message string, callback StreamCallback, options ...GenerateOption) error {
	req := &GenerateRequest{
		Message: message,
		Context: ctx,
	}

	// Apply options
	for _, opt := range options {
		opt(req)
	}

	return a.streamResponse(req, callback)
}

// ListTools returns all available tools for this agent.
func (a *Agent) ListTools() []agent.ToolInfo {
	tools := make([]agent.ToolInfo, 0, len(a.toolsMap))
	for _, tool := range a.toolsMap {
		schema, _ := tool.GetSchema()
		tools = append(tools, agent.ToolInfo{
			Name:        tool.GetName(),
			Description: tool.GetDescription(),
			Schema:      schema,
		})
	}
	return tools
}

// Internal implementation methods

// generateResponse processes a generation request and returns a response.
func (a *Agent) generateResponse(req *GenerateRequest) (*Response, error) {
	// Create execution context
	execCtx := &Context{
		Context:      req.Context,
		Message:      req.Message,
		Conversation: req.Conversation,
		User:         req.User,
		SessionID:    req.SessionID,
		tools:        a.toolsMap,
	}

	// Build the system prompt with instructions and tool information
	systemPrompt := a.buildSystemPrompt()

	// Add system message to conversation if not already present
	conversation := a.prepareConversation(req.Conversation, systemPrompt)

	// Add current user message
	conversation = append(conversation, Message{
		Role:    "user",
		Content: req.Message,
	})

	// Execute with tool calling support
	return a.executeWithTools(execCtx, conversation)
}

// streamResponse processes a streaming generation request.
func (a *Agent) streamResponse(req *GenerateRequest, callback StreamCallback) error {
	// Build the system prompt with instructions and tool information
	systemPrompt := a.buildSystemPrompt()

	// Add system message to conversation if not already present
	conversation := a.prepareConversation(req.Conversation, systemPrompt)

	// Add current user message
	conversation = append(conversation, Message{
		Role:    "user",
		Content: req.Message,
	})

	// Convert conversation to LLM format
	llmMessages := a.convertConversationToLLM(conversation)

	// Add system message with instructions and available tools
	systemMessage := a.buildSystemMessage()
	if systemMessage != "" {
		llmMessages = append([]llm.Message{{
			Role:    "system",
			Content: systemMessage,
		}}, llmMessages...)
	}

	// Use streaming API
	return a.llmClient.Stream(
		req.Context,
		a.Model,
		llmMessages,
		a.Tools,
		func(chunk llm.StreamChunk) error {
			// Extract content from streaming chunk
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				return callback(chunk.Choices[0].Delta.Content)
			}
			return nil
		},
	)
}

// buildSystemPrompt creates the system prompt with instructions and tool info.
func (a *Agent) buildSystemPrompt() string {
	prompt := a.Instructions

	if len(a.toolsMap) > 0 {
		prompt += "\n\nYou have access to the following tools:\n"
		for _, tool := range a.toolsMap {
			prompt += fmt.Sprintf("- %s: %s\n", tool.GetName(), tool.GetDescription())
		}
		prompt += "\nUse tools when appropriate to help the user."
	}

	return prompt
}

// prepareConversation ensures the conversation has a system message.
func (a *Agent) prepareConversation(conversation []Message, systemPrompt string) []Message {
	if len(conversation) == 0 || conversation[0].Role != "system" {
		// Prepend system message
		return append([]Message{{
			Role:    "system",
			Content: systemPrompt,
		}}, conversation...)
	}

	// Update existing system message
	conversation[0].Content = systemPrompt
	return conversation
}

// executeWithTools executes the agent with tool calling support.
func (a *Agent) executeWithTools(ctx *Context, conversation []Message) (*Response, error) {
	maxSteps := a.MaxSteps
	if maxSteps == 0 {
		maxSteps = 10
	}

	currentConversation := make([]Message, len(conversation))
	copy(currentConversation, conversation)

	for step := 0; step < maxSteps; step++ {
		// Simulate LLM call with current conversation
		llmResponse, err := a.callLLM(ctx, currentConversation)
		if err != nil {
			return nil, fmt.Errorf("LLM call failed at step %d: %w", step, err)
		}

		// Add LLM response to conversation
		currentConversation = append(currentConversation, Message{
			Role:      "assistant",
			Content:   llmResponse.Message,
			ToolCalls: llmResponse.ToolCalls,
		})

		// Check if LLM wants to use tools
		if len(llmResponse.ToolCalls) == 0 {
			// No tool calls, we're done
			return &Response{
				Message: llmResponse.Message,
				Metadata: map[string]interface{}{
					"agent_name":     a.Name,
					"model_id":       a.Model.ID(),
					"model_provider": a.Model.Provider(),
					"model_name":     a.Model.Name(),
					"tools":          len(a.toolsMap),
					"steps_taken":    step + 1,
				},
			}, nil
		}

		// Execute tool calls
		for _, toolCall := range llmResponse.ToolCalls {
			result, err := a.executeToolCall(ctx, toolCall)
			if err != nil {
				result.Error = err.Error()
			}

			// Add tool result to conversation
			currentConversation = append(currentConversation, Message{
				Role:    "tool",
				Content: a.formatToolResult(result),
				ToolID:  result.ID,
			})
		}

		// Continue the loop for next LLM response
	}

	// Max steps reached
	return &Response{
		Message: "I've reached the maximum number of steps. Please try again with a simpler request.",
		Metadata: map[string]interface{}{
			"agent_name":        a.Name,
			"model_id":          a.Model.ID(),
			"model_provider":    a.Model.Provider(),
			"model_name":        a.Model.Name(),
			"tools":             len(a.toolsMap),
			"steps_taken":       maxSteps,
			"max_steps_reached": true,
		},
	}, nil
}

// callLLM makes an actual request to the LLM API and returns a response with potential tool calls.
func (a *Agent) callLLM(ctx *Context, conversation []Message) (*Response, error) {
	// Convert conversation to LLM format
	llmMessages := a.convertConversationToLLM(conversation)

	// Add system message with instructions and available tools
	systemMessage := a.buildSystemMessage()
	if systemMessage != "" {
		llmMessages = append([]llm.Message{{
			Role:    "system",
			Content: systemMessage,
		}}, llmMessages...)
	}

	// Call the LLM API
	llmResponse, err := a.llmClient.Complete(
		ctx.Context,
		a.Model,
		llmMessages,
		a.Tools,
	)
	if err != nil {
		return nil, fmt.Errorf("LLM API call failed: %w", err)
	}

	if len(llmResponse.Choices) == 0 {
		return nil, fmt.Errorf("LLM returned no choices")
	}

	choice := llmResponse.Choices[0]
	message := choice.Message

	// Convert LLM response to our format
	response := &Response{
		Message: message.Content,
		Metadata: map[string]interface{}{
			"agent_name":     a.Name,
			"model_id":       a.Model.ID(),
			"model_provider": a.Model.Provider(),
			"model_name":     a.Model.Name(),
			"tools":          len(a.toolsMap),
			"finish_reason":  choice.FinishReason,
		},
	}

	// Add token usage info if available
	if llmResponse.Usage.TotalTokens > 0 {
		response.Metadata["tokens_used"] = llmResponse.Usage.TotalTokens
		response.Metadata["prompt_tokens"] = llmResponse.Usage.PromptTokens
		response.Metadata["completion_tokens"] = llmResponse.Usage.CompletionTokens
	}

	// Convert tool calls if present
	if len(message.ToolCalls) > 0 {
		response.ToolCalls = make([]ToolCall, len(message.ToolCalls))
		for i, tc := range message.ToolCalls {
			// Parse tool call arguments
			var params map[string]interface{}
			if tc.Function.Arguments != "" {
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &params); err != nil {
					return nil, fmt.Errorf("failed to parse tool call arguments: %w", err)
				}
			} else {
				params = make(map[string]interface{})
			}

			response.ToolCalls[i] = ToolCall{
				ID:     tc.ID,
				Name:   tc.Function.Name,
				Params: params,
			}
		}
	}

	return response, nil
}

// executeToolCall executes a single tool call and returns the result.
func (a *Agent) executeToolCall(ctx *Context, toolCall ToolCall) (ToolCall, error) {
	result := toolCall // Copy the original tool call

	// Execute the tool
	toolResult, err := ctx.ExecuteTool(toolCall.Name, toolCall.Params)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	result.Result = toolResult
	return result, nil
}

// formatToolResult formats a tool result for inclusion in conversation.
func (a *Agent) formatToolResult(toolResult ToolCall) string {
	if toolResult.Error != "" {
		return fmt.Sprintf("Tool '%s' failed: %s", toolResult.Name, toolResult.Error)
	}

	// Convert result to JSON for the conversation
	if toolResult.Result != nil {
		resultJSON, err := json.Marshal(toolResult.Result)
		if err != nil {
			return fmt.Sprintf("Tool '%s' completed but result formatting failed: %v", toolResult.Name, err)
		}
		return fmt.Sprintf("Tool '%s' result: %s", toolResult.Name, string(resultJSON))
	}

	return fmt.Sprintf("Tool '%s' completed successfully", toolResult.Name)
}

// Helper methods for LLM integration

// convertConversationToLLM converts internal conversation format to LLM format.
func (a *Agent) convertConversationToLLM(conversation []Message) []llm.Message {
	llmMessages := make([]llm.Message, len(conversation))

	for i, msg := range conversation {
		llmMsg := llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}

		// Handle tool call responses
		if msg.Role == "tool" && msg.ToolID != "" {
			llmMsg.ToolCallID = msg.ToolID
		}

		llmMessages[i] = llmMsg
	}

	return llmMessages
}

// buildSystemMessage creates the system prompt with instructions and tool descriptions.
func (a *Agent) buildSystemMessage() string {
	var systemParts []string

	// Add basic instructions
	if a.Instructions != "" {
		systemParts = append(systemParts, a.Instructions)
	}

	// Add tool descriptions if tools are available
	if len(a.Tools) > 0 {
		systemParts = append(systemParts, "\nYou have access to the following tools:")
		for _, tool := range a.Tools {
			toolDesc := fmt.Sprintf("- %s: %s", tool.GetName(), tool.GetDescription())
			systemParts = append(systemParts, toolDesc)
		}
		systemParts = append(systemParts, "\nUse these tools when appropriate to help the user. Always explain what you're doing when using tools.")
	}

	return strings.Join(systemParts, "\n")
}

type agentManager struct {
	agents map[uuid.UUID]*Agent
	mu     sync.RWMutex
}

func newAgentManager(agents map[uuid.UUID]*Agent) *agentManager {
	return &agentManager{
		agents: agents,
	}
}

func (am *agentManager) getAgent(id uuid.UUID) *Agent {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.agents[id]
}
