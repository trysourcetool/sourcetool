package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/trysourcetool/sourcetool-go/agent"
	"github.com/trysourcetool/sourcetool-go/agent/models"
)

// Client represents an LLM client that can communicate with various AI providers.
type Client struct {
	httpClient  *http.Client
	timeout     time.Duration
	maxRetries  int
	retryDelay  time.Duration
	enableDebug bool
}

// ClientConfig represents configuration for the LLM client.
type ClientConfig struct {
	Timeout     time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
	EnableDebug bool
}

// NewClient creates a new LLM client.
func NewClient(timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		timeout:     timeout,
		maxRetries:  3,
		retryDelay:  1 * time.Second,
		enableDebug: false,
	}
}

// NewClientWithConfig creates a new LLM client with custom configuration.
func NewClientWithConfig(config ClientConfig) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 1 * time.Second
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		timeout:     config.Timeout,
		maxRetries:  config.MaxRetries,
		retryDelay:  config.RetryDelay,
		enableDebug: config.EnableDebug,
	}
}

// Message represents a conversation message.
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// ToolCall represents a tool function call.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall represents the function to call.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// CompletionRequest represents a request to the LLM.
type CompletionRequest struct {
	Model       string       `json:"model"`
	Messages    []Message    `json:"messages"`
	Tools       []ToolSchema `json:"tools,omitempty"`
	Temperature *float64     `json:"temperature,omitempty"`
	MaxTokens   int          `json:"max_tokens,omitempty"`
	TopP        *float64     `json:"top_p,omitempty"`
	TopK        *int         `json:"top_k,omitempty"`
	Stop        []string     `json:"stop,omitempty"`
	Seed        *int         `json:"seed,omitempty"`
	Stream      bool         `json:"stream,omitempty"`
}

// CompletionResponse represents a response from the LLM.
type CompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
}

// Choice represents a completion choice.
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage statistics.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ToolSchema represents a tool schema for the LLM.
type ToolSchema struct {
	Type     string             `json:"type"`
	Function ToolFunctionSchema `json:"function"`
}

// ToolFunctionSchema represents the schema for a tool function.
type ToolFunctionSchema struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// StreamChunk represents a streaming response chunk.
type StreamChunk struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int64          `json:"created"`
	Model             string         `json:"model"`
	Choices           []StreamChoice `json:"choices"`
	SystemFingerprint string         `json:"system_fingerprint,omitempty"`
}

// StreamChoice represents a streaming choice.
type StreamChoice struct {
	Index        int         `json:"index"`
	Delta        StreamDelta `json:"delta"`
	FinishReason *string     `json:"finish_reason"`
}

// StreamDelta represents the delta in a streaming response.
type StreamDelta struct {
	Role      string     `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Complete sends a completion request to the appropriate LLM provider.
func (c *Client) Complete(ctx context.Context, model models.Model, messages []Message, tools []agent.Tool) (*CompletionResponse, error) {
	// Convert tools to LLM format
	toolSchemas := c.convertTools(tools)

	// Create request
	req := CompletionRequest{
		Model:    model.Name(),
		Messages: messages,
		Tools:    toolSchemas,
	}

	// Apply model configuration
	config := model.Config()
	if config.Temperature != nil {
		req.Temperature = config.Temperature
	}
	if config.MaxTokens > 0 {
		req.MaxTokens = config.MaxTokens
	}
	if config.TopP != nil {
		req.TopP = config.TopP
	}
	if config.TopK != nil {
		req.TopK = config.TopK
	}
	if len(config.Stop) > 0 {
		req.Stop = config.Stop
	}
	if config.Seed != nil {
		req.Seed = config.Seed
	}

	// Route to appropriate provider
	switch model.Provider() {
	case "anthropic":
		return c.callAnthropic(ctx, req)
	case "openai":
		return c.callOpenAI(ctx, req)
	case "xai":
		return c.callXAI(ctx, req)
	case "google":
		return c.callGoogle(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported model provider: %s", model.Provider())
	}
}

// Stream sends a streaming completion request.
func (c *Client) Stream(ctx context.Context, model models.Model, messages []Message, tools []agent.Tool, callback func(chunk StreamChunk) error) error {
	// Convert tools to LLM format
	toolSchemas := c.convertTools(tools)

	// Create request
	req := CompletionRequest{
		Model:    model.Name(),
		Messages: messages,
		Tools:    toolSchemas,
		Stream:   true,
	}

	// Apply model configuration
	config := model.Config()
	if config.Temperature != nil {
		req.Temperature = config.Temperature
	}
	if config.MaxTokens > 0 {
		req.MaxTokens = config.MaxTokens
	}
	if config.TopP != nil {
		req.TopP = config.TopP
	}
	if config.TopK != nil {
		req.TopK = config.TopK
	}
	if len(config.Stop) > 0 {
		req.Stop = config.Stop
	}
	if config.Seed != nil {
		req.Seed = config.Seed
	}

	// Route to appropriate provider
	switch model.Provider() {
	case "anthropic":
		return c.streamAnthropic(ctx, req, callback)
	case "openai":
		return c.streamOpenAI(ctx, req, callback)
	case "xai":
		return c.streamXAI(ctx, req, callback)
	case "google":
		return c.streamGoogle(ctx, req, callback)
	default:
		return fmt.Errorf("unsupported model provider: %s", model.Provider())
	}
}

// convertTools converts agent tools to LLM tool schemas.
func (c *Client) convertTools(tools []agent.Tool) []ToolSchema {
	if len(tools) == 0 {
		return nil
	}

	schemas := make([]ToolSchema, 0, len(tools))
	for _, tool := range tools {
		schema, err := tool.GetSchema()
		if err != nil {
			continue // Skip tools with invalid schemas
		}

		llmSchema := ToolSchema{
			Type: "function",
			Function: ToolFunctionSchema{
				Name:        tool.GetName(),
				Description: tool.GetDescription(),
				Parameters: map[string]interface{}{
					"type":       schema.Type,
					"properties": schema.Properties,
					"required":   schema.Required,
				},
			},
		}
		schemas = append(schemas, llmSchema)
	}

	return schemas
}

// callOpenAI makes a request to OpenAI API.
func (c *Client) callOpenAI(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	url := "https://api.openai.com/v1/chat/completions"

	// Convert to OpenAI format (remove unsupported parameters)
	openaiReq := c.convertToOpenAIFormat(req)

	jsonData, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+getAPIKey("OPENAI_API_KEY"))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// callAnthropic makes a request to Anthropic API.
func (c *Client) callAnthropic(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	// Convert to Anthropic format
	anthropicReq := c.convertToAnthropicFormat(req)

	url := "https://api.anthropic.com/v1/messages"

	jsonData, err := json.Marshal(anthropicReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Anthropic-specific headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", getAPIKey("ANTHROPIC_API_KEY"))
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse Anthropic response and convert to standard format
	var anthropicResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return c.convertFromAnthropicFormat(anthropicResp), nil
}

// callXAI makes a request to xAI API.
func (c *Client) callXAI(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	// xAI uses OpenAI-compatible format but with some limitations
	url := "https://api.x.ai/v1/chat/completions"

	// Convert to xAI format (remove unsupported parameters for reasoning models)
	xaiReq := c.convertToXAIFormat(req)

	jsonData, err := json.Marshal(xaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+getAPIKey("XAI_API_KEY"))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// callGoogle makes a request to Google AI API.
func (c *Client) callGoogle(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	// Convert to Google format
	googleReq := c.convertToGoogleFormat(req)

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/%s:generateContent", req.Model)

	jsonData, err := json.Marshal(googleReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", getAPIKey("GOOGLE_API_KEY"))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse Google response and convert to standard format
	var googleResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return c.convertFromGoogleFormat(googleResp), nil
}

// Streaming implementations.
func (c *Client) streamOpenAI(ctx context.Context, req CompletionRequest, callback func(chunk StreamChunk) error) error {
	url := "https://api.openai.com/v1/chat/completions"

	// Convert to OpenAI format (remove unsupported parameters)
	openaiReq := c.convertToOpenAIFormat(req)

	jsonData, err := json.Marshal(openaiReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+getAPIKey("OPENAI_API_KEY"))
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse streaming response
	return c.parseStreamingResponse(resp.Body, callback)
}

func (c *Client) streamAnthropic(ctx context.Context, req CompletionRequest, callback func(chunk StreamChunk) error) error {
	// Implement Anthropic streaming
	return fmt.Errorf("streaming not yet implemented for Anthropic")
}

func (c *Client) streamXAI(ctx context.Context, req CompletionRequest, callback func(chunk StreamChunk) error) error {
	// xAI uses OpenAI-compatible streaming
	return c.streamOpenAI(ctx, req, callback)
}

func (c *Client) streamGoogle(ctx context.Context, req CompletionRequest, callback func(chunk StreamChunk) error) error {
	// Implement Google streaming
	return fmt.Errorf("streaming not yet implemented for Google")
}

// Helper functions for format conversion

// convertToOpenAIFormat removes unsupported parameters for OpenAI.
func (c *Client) convertToOpenAIFormat(req CompletionRequest) map[string]interface{} {
	openaiReq := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
	}

	// Add supported parameters
	if req.Temperature != nil {
		openaiReq["temperature"] = *req.Temperature
	}
	if req.MaxTokens > 0 {
		openaiReq["max_tokens"] = req.MaxTokens
	}
	if req.TopP != nil {
		openaiReq["top_p"] = *req.TopP
	}
	// Note: OpenAI doesn't support top_k
	if len(req.Stop) > 0 {
		openaiReq["stop"] = req.Stop
	}
	if req.Seed != nil {
		openaiReq["seed"] = *req.Seed
	}
	if len(req.Tools) > 0 {
		openaiReq["tools"] = req.Tools
	}
	if req.Stream {
		openaiReq["stream"] = req.Stream
	}

	return openaiReq
}

// convertToXAIFormat removes unsupported parameters for xAI (reasoning models).
func (c *Client) convertToXAIFormat(req CompletionRequest) map[string]interface{} {
	xaiReq := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
	}

	// Add supported parameters
	if req.Temperature != nil {
		xaiReq["temperature"] = *req.Temperature
	}
	if req.MaxTokens > 0 {
		xaiReq["max_tokens"] = req.MaxTokens
	}
	if req.TopP != nil {
		xaiReq["top_p"] = *req.TopP
	}
	if req.TopK != nil {
		xaiReq["top_k"] = *req.TopK
	}
	// Note: Grok 4 reasoning models don't support stop sequences
	if req.Seed != nil {
		xaiReq["seed"] = *req.Seed
	}
	if len(req.Tools) > 0 {
		xaiReq["tools"] = req.Tools
	}
	if req.Stream {
		xaiReq["stream"] = req.Stream
	}

	return xaiReq
}

func (c *Client) convertToAnthropicFormat(req CompletionRequest) map[string]interface{} {
	// Convert OpenAI format to Anthropic format
	messages := make([]map[string]interface{}, 0, len(req.Messages))

	var systemMessage string
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			systemMessage = msg.Content
			continue
		}

		messages = append(messages, map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	anthropicReq := map[string]interface{}{
		"model":      req.Model,
		"messages":   messages,
		"max_tokens": 4096, // Anthropic requires max_tokens
	}

	if systemMessage != "" {
		anthropicReq["system"] = systemMessage
	}
	if req.Temperature != nil {
		anthropicReq["temperature"] = *req.Temperature
	}
	if req.MaxTokens > 0 {
		anthropicReq["max_tokens"] = req.MaxTokens
	}
	if req.TopP != nil {
		anthropicReq["top_p"] = *req.TopP
	}
	if len(req.Stop) > 0 {
		anthropicReq["stop_sequences"] = req.Stop
	}
	// Note: Anthropic doesn't support TopK and Seed

	return anthropicReq
}

func (c *Client) convertFromAnthropicFormat(resp map[string]interface{}) *CompletionResponse {
	// Convert Anthropic response to standard format
	content := ""
	if contentArray, ok := resp["content"].([]interface{}); ok && len(contentArray) > 0 {
		if contentObj, ok := contentArray[0].(map[string]interface{}); ok {
			if text, ok := contentObj["text"].(string); ok {
				content = text
			}
		}
	}

	return &CompletionResponse{
		ID:      getString(resp, "id"),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   getString(resp, "model"),
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: getString(resp, "stop_reason"),
			},
		},
		Usage: Usage{
			PromptTokens:     getInt(resp, "usage.input_tokens"),
			CompletionTokens: getInt(resp, "usage.output_tokens"),
			TotalTokens:      getInt(resp, "usage.input_tokens") + getInt(resp, "usage.output_tokens"),
		},
	}
}

func (c *Client) convertToGoogleFormat(req CompletionRequest) map[string]interface{} {
	// Convert to Google AI format
	contents := make([]map[string]interface{}, 0, len(req.Messages))

	for _, msg := range req.Messages {
		role := msg.Role
		switch role {
		case "assistant":
			role = "model"
		case "system":
			role = "user" // Google treats system as user message
		}

		contents = append(contents, map[string]interface{}{
			"role": role,
			"parts": []map[string]interface{}{
				{"text": msg.Content},
			},
		})
	}

	googleReq := map[string]interface{}{
		"contents": contents,
	}

	// Add generation config
	if req.Temperature != nil || req.MaxTokens > 0 || req.TopP != nil || req.TopK != nil || len(req.Stop) > 0 || req.Seed != nil {
		generationConfig := make(map[string]interface{})
		if req.Temperature != nil {
			generationConfig["temperature"] = *req.Temperature
		}
		if req.MaxTokens > 0 {
			generationConfig["maxOutputTokens"] = req.MaxTokens
		}
		if req.TopP != nil {
			generationConfig["topP"] = *req.TopP
		}
		if req.TopK != nil {
			generationConfig["topK"] = *req.TopK
		}
		if len(req.Stop) > 0 {
			generationConfig["stopSequences"] = req.Stop
		}
		if req.Seed != nil {
			generationConfig["seed"] = *req.Seed
		}
		googleReq["generationConfig"] = generationConfig
	}

	return googleReq
}

func (c *Client) convertFromGoogleFormat(resp map[string]interface{}) *CompletionResponse {
	// Convert Google response to standard format
	content := ""
	if candidates, ok := resp["candidates"].([]interface{}); ok && len(candidates) > 0 {
		if candidate, ok := candidates[0].(map[string]interface{}); ok {
			if contentObj, ok := candidate["content"].(map[string]interface{}); ok {
				if parts, ok := contentObj["parts"].([]interface{}); ok && len(parts) > 0 {
					if part, ok := parts[0].(map[string]interface{}); ok {
						if text, ok := part["text"].(string); ok {
							content = text
						}
					}
				}
			}
		}
	}

	return &CompletionResponse{
		ID:      "google-" + fmt.Sprintf("%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "google-model",
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: "stop",
			},
		},
	}
}

func (c *Client) parseStreamingResponse(body io.Reader, callback func(chunk StreamChunk) error) error {
	// Parse Server-Sent Events format
	// This is a simplified implementation
	// In production, you'd want to use a proper SSE parser

	// For now, return an error indicating streaming needs full implementation
	return fmt.Errorf("streaming response parsing not fully implemented")
}

// Utility functions.
func getAPIKey(envVar string) string {
	// In a real implementation, get from environment variables
	// For demo purposes, return placeholder
	return "placeholder-api-key"
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	// Support nested keys like "usage.input_tokens"
	keys := strings.Split(key, ".")
	current := m

	for i, k := range keys {
		if i == len(keys)-1 {
			if val, ok := current[k]; ok {
				if num, ok := val.(float64); ok {
					return int(num)
				}
				if num, ok := val.(int); ok {
					return num
				}
			}
		} else {
			if val, ok := current[k].(map[string]interface{}); ok {
				current = val
			} else {
				break
			}
		}
	}
	return 0
}
