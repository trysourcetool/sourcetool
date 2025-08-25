package sourcetool

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/agent"
	"github.com/trysourcetool/sourcetool-go/agent/models"
)

// Test parameter struct.
type TestParams struct {
	Name  string `json:"name" desc:"User name" required:"true"`
	Age   int    `json:"age,omitempty" desc:"User age" min:"0" max:"150"`
	Email string `json:"email,omitempty" desc:"Email address"`
}

func TestToolSchema(t *testing.T) {
	tool := agent.NewTool("test_tool", "Test tool for validation",
		func(ctx context.Context, params TestParams) (interface{}, error) {
			return map[string]interface{}{
				"name": params.Name,
				"age":  params.Age,
			}, nil
		},
	)

	schema, err := tool.GetSchema()
	if err != nil {
		t.Fatalf("Failed to get schema: %v", err)
	}

	// Check schema type
	if schema.Type != "object" {
		t.Errorf("Expected schema type 'object', got '%s'", schema.Type)
	}

	// Check required fields
	if len(schema.Required) != 1 || schema.Required[0] != "name" {
		t.Errorf("Expected required fields ['name'], got %v", schema.Required)
	}

	// Check properties
	nameProperty, exists := schema.Properties["name"]
	if !exists {
		t.Error("Expected 'name' property to exist")
	}
	if nameProperty.Type != "string" {
		t.Errorf("Expected name property type 'string', got '%s'", nameProperty.Type)
	}

	ageProperty, exists := schema.Properties["age"]
	if !exists {
		t.Error("Expected 'age' property to exist")
	}
	if ageProperty.Type != "integer" {
		t.Errorf("Expected age property type 'integer', got '%s'", ageProperty.Type)
	}
	if ageProperty.Minimum == nil || *ageProperty.Minimum != 0 {
		t.Error("Expected age minimum to be 0")
	}
	if ageProperty.Maximum == nil || *ageProperty.Maximum != 150 {
		t.Error("Expected age maximum to be 150")
	}
}

func TestToolExecution(t *testing.T) {
	tool := agent.NewTool("test_execution", "Test tool execution",
		func(ctx context.Context, params TestParams) (interface{}, error) {
			return map[string]interface{}{
				"greeting": "Hello " + params.Name,
				"age":      params.Age,
			}, nil
		},
	)

	// Test with valid parameters
	paramsJSON := `{"name": "John", "age": 30}`
	result, err := tool.Execute(context.Background(), json.RawMessage(paramsJSON))
	if err != nil {
		t.Fatalf("Tool execution failed: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected result to be a map")
	}

	if resultMap["greeting"] != "Hello John" {
		t.Errorf("Expected greeting 'Hello John', got '%v'", resultMap["greeting"])
	}

	// Test with invalid parameters
	invalidJSON := `{"age": "invalid"}`
	_, err = tool.Execute(context.Background(), json.RawMessage(invalidJSON))
	if err == nil {
		t.Error("Expected error with invalid parameters")
	}
}

func TestAgentValidation(t *testing.T) {
	// Valid agent
	validAgent := &Agent{
		Name:         "test_agent",
		Description:  "Test agent",
		Instructions: "You are a helpful assistant.",
		Model:        models.Anthropic("claude-3-5-sonnet"),
	}

	if err := validAgent.Validate(); err != nil {
		t.Errorf("Valid agent should not fail validation: %v", err)
	}

	// Check model is set correctly
	if validAgent.Model.Provider() != "anthropic" {
		t.Errorf("Expected model provider to be 'anthropic', got '%s'", validAgent.Model.Provider())
	}
	if validAgent.Model.Name() != "claude-3-5-sonnet" {
		t.Errorf("Expected model name to be 'claude-3-5-sonnet', got '%s'", validAgent.Model.Name())
	}
	if validAgent.Model.ID() != "anthropic/claude-3-5-sonnet" {
		t.Errorf("Expected model ID to be 'anthropic/claude-3-5-sonnet', got '%s'", validAgent.Model.ID())
	}

	// Check defaults are set
	if validAgent.MaxSteps != 10 {
		t.Errorf("Expected default MaxSteps to be 10, got %d", validAgent.MaxSteps)
	}

	// Invalid agent - missing name
	invalidAgent := &Agent{
		Description:  "Test agent",
		Instructions: "You are a helpful assistant.",
		Model:        models.Anthropic("claude-3-5-sonnet"),
	}

	if err := invalidAgent.Validate(); err == nil {
		t.Error("Agent without name should fail validation")
	}

	// Invalid agent - missing instructions
	invalidAgent2 := &Agent{
		Name:        "test_agent",
		Description: "Test agent",
		Model:       models.Anthropic("claude-3-5-sonnet"),
	}

	if err := invalidAgent2.Validate(); err == nil {
		t.Error("Agent without instructions should fail validation")
	}

	// Invalid agent - missing model
	invalidAgent3 := &Agent{
		Name:         "test_agent",
		Description:  "Test agent",
		Instructions: "You are a helpful assistant.",
	}

	if err := invalidAgent3.Validate(); err == nil {
		t.Error("Agent without model should fail validation")
	}
}

func TestAgentContextToolExecution(t *testing.T) {
	// Create a test tool
	testTool := agent.NewTool("test_tool", "Test tool",
		func(ctx context.Context, params TestParams) (interface{}, error) {
			return map[string]interface{}{
				"result": "success",
				"name":   params.Name,
			}, nil
		},
	)

	// Test tool execution directly
	paramsJSON := json.RawMessage(`{"name": "TestUser", "age": 25}`)
	result, err := testTool.Execute(context.Background(), paramsJSON)
	if err != nil {
		t.Fatalf("Tool execution failed: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected result to be a map")
	}

	if resultMap["result"] != "success" {
		t.Errorf("Expected result 'success', got '%v'", resultMap["result"])
	}

	// Test schema generation
	schema, err := testTool.GetSchema()
	if err != nil {
		t.Fatalf("Schema generation failed: %v", err)
	}

	if schema.Type != "object" {
		t.Errorf("Expected schema type 'object', got '%s'", schema.Type)
	}
}

func TestAgentGenerate(t *testing.T) {
	// Create a simple agent
	testAgent := &Agent{
		Name:         "test_generate",
		Description:  "Test agent for generation",
		Instructions: "You are a helpful test assistant.",
		Model:        models.OpenAI("gpt-4o-mini"),
	}

	// Initialize agent with test data
	testID, _ := uuid.NewV4()
	testAgent = newAgent(testAgent, testID, "/test", []string{})

	// Test Generate method
	response, err := testAgent.Generate(context.Background(), "Hello!")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if response == nil {
		t.Fatal("Response should not be nil")
	}

	if response.Message == "" {
		t.Error("Response message should not be empty")
	}

	// Check metadata
	if response.Metadata == nil {
		t.Error("Response metadata should not be nil")
	}

	agentName, ok := response.Metadata["agent_name"]
	if !ok || agentName != "test_generate" {
		t.Errorf("Expected agent_name in metadata to be 'test_generate', got %v", agentName)
	}

	modelProvider, ok := response.Metadata["model_provider"]
	if !ok || modelProvider != "openai" {
		t.Errorf("Expected model_provider in metadata to be 'openai', got %v", modelProvider)
	}

	modelName, ok := response.Metadata["model_name"]
	if !ok || modelName != "gpt-4o-mini" {
		t.Errorf("Expected model_name in metadata to be 'gpt-4o-mini', got %v", modelName)
	}
}

func TestModelConfiguration(t *testing.T) {
	// Test Anthropic model
	anthropicModel := models.Anthropic("claude-3-5-sonnet",
		models.WithTemperature(0.7),
		models.WithMaxTokens(2048),
	)

	if anthropicModel.Provider() != "anthropic" {
		t.Errorf("Expected provider 'anthropic', got '%s'", anthropicModel.Provider())
	}
	if anthropicModel.Name() != "claude-3-5-sonnet" {
		t.Errorf("Expected name 'claude-3-5-sonnet', got '%s'", anthropicModel.Name())
	}
	if anthropicModel.ID() != "anthropic/claude-3-5-sonnet" {
		t.Errorf("Expected ID 'anthropic/claude-3-5-sonnet', got '%s'", anthropicModel.ID())
	}

	config := anthropicModel.Config()
	if config.Temperature == nil || *config.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %v", config.Temperature)
	}
	if config.MaxTokens != 2048 {
		t.Errorf("Expected max_tokens 2048, got %d", config.MaxTokens)
	}

	// Test OpenAI model
	openaiModel := models.OpenAI("gpt-4o",
		models.WithTopP(0.9),
		models.WithStop("END"),
	)

	if openaiModel.Provider() != "openai" {
		t.Errorf("Expected provider 'openai', got '%s'", openaiModel.Provider())
	}
	if openaiModel.ID() != "openai/gpt-4o" {
		t.Errorf("Expected ID 'openai/gpt-4o', got '%s'", openaiModel.ID())
	}

	openaiConfig := openaiModel.Config()
	if openaiConfig.TopP == nil || *openaiConfig.TopP != 0.9 {
		t.Errorf("Expected top_p 0.9, got %v", openaiConfig.TopP)
	}
	if len(openaiConfig.Stop) != 1 || openaiConfig.Stop[0] != "END" {
		t.Errorf("Expected stop sequence ['END'], got %v", openaiConfig.Stop)
	}

	// Test xAI model
	xaiModel := models.XAI("grok-beta",
		models.WithTemperature(0.8),
		models.WithTopK(40),
	)

	if xaiModel.Provider() != "xai" {
		t.Errorf("Expected provider 'xai', got '%s'", xaiModel.Provider())
	}
	if xaiModel.Name() != "grok-beta" {
		t.Errorf("Expected name 'grok-beta', got '%s'", xaiModel.Name())
	}
	if xaiModel.ID() != "xai/grok-beta" {
		t.Errorf("Expected ID 'xai/grok-beta', got '%s'", xaiModel.ID())
	}

	xaiConfig := xaiModel.Config()
	if xaiConfig.Temperature == nil || *xaiConfig.Temperature != 0.8 {
		t.Errorf("Expected temperature 0.8, got %v", xaiConfig.Temperature)
	}
	if xaiConfig.TopK == nil || *xaiConfig.TopK != 40 {
		t.Errorf("Expected top_k 40, got %v", xaiConfig.TopK)
	}

	// Test Google/Gemini model
	geminiModel := models.Google("gemini-pro",
		models.WithSeed(12345),
	)

	if geminiModel.Provider() != "google" {
		t.Errorf("Expected provider 'google', got '%s'", geminiModel.Provider())
	}
	if geminiModel.Name() != "gemini-pro" {
		t.Errorf("Expected name 'gemini-pro', got '%s'", geminiModel.Name())
	}
	if geminiModel.ID() != "google/gemini-pro" {
		t.Errorf("Expected ID 'google/gemini-pro', got '%s'", geminiModel.ID())
	}

	geminiConfig := geminiModel.Config()
	if geminiConfig.Seed == nil || *geminiConfig.Seed != 12345 {
		t.Errorf("Expected seed 12345, got %v", geminiConfig.Seed)
	}

	// Test Google function (should be equivalent to Gemini)
	googleModel := models.Google("gemini-flash")
	if googleModel.Provider() != "google" {
		t.Errorf("Expected provider 'google', got '%s'", googleModel.Provider())
	}
	if googleModel.ID() != "google/gemini-flash" {
		t.Errorf("Expected ID 'google/gemini-flash', got '%s'", googleModel.ID())
	}
}

// Benchmark tests.
func BenchmarkToolSchemaGeneration(b *testing.B) {
	tool := agent.NewTool("benchmark_tool", "Benchmark tool",
		func(ctx context.Context, params TestParams) (interface{}, error) {
			return params, nil
		},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tool.GetSchema()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkToolExecution(b *testing.B) {
	tool := agent.NewTool("benchmark_execution", "Benchmark execution",
		func(ctx context.Context, params TestParams) (interface{}, error) {
			return params, nil
		},
	)

	paramsJSON := json.RawMessage(`{"name": "BenchmarkUser", "age": 30}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tool.Execute(context.Background(), paramsJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}
