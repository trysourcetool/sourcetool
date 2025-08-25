# 🚀 Comprehensive AI Agent Demo

This example showcases the complete capabilities of the Sourcetool AI agent system with a comprehensive demonstration of all features, tools, and integrations.

## 🎯 What This Demo Covers

### 🤖 Agent Types
- **Main Multi-Purpose Agent**: Comprehensive assistant with all tools
- **Search Specialist**: Advanced search and information retrieval
- **Support Specialist**: Customer support and ticket management
- **Technical Specialist**: Calculations, file operations, and database queries

### 🔧 Tool Showcase
1. **Advanced Search** - Complex parameter validation with nested structures
2. **User Management** - User creation with detailed settings and validation
3. **Ticket System** - Support ticket creation with custom validation
4. **File Operations** - File creation and content management
5. **Calculator** - Mathematical calculations with precision control
6. **Weather API** - Location-based weather information
7. **System Status** - Comprehensive system monitoring
8. **Database Operations** - SQL query simulation

### 🌟 Advanced Features Demonstrated
- **Multiple Model Providers**: Anthropic, OpenAI, xAI, Google/Gemini
- **Complex Parameter Validation**: Nested structures, custom validation, enums
- **Error Handling**: Graceful error recovery and user feedback
- **Streaming Responses**: Real-time response streaming
- **UI Integration**: Web interface with agent interactions
- **Multi-step Conversations**: Complex task execution across multiple steps
- **Session Management**: User context and session tracking

## 🏃‍♂️ Running the Demo

### Prerequisites
```bash
go mod tidy
```

### Basic Execution
```bash
go run main.go
```

### What You'll See
The demo runs through several phases:

1. **🤖 Agent Creation**: Creates multiple specialized agents
2. **🧪 Feature Testing**: Tests all major capabilities
3. **⚠️ Error Handling**: Demonstrates error recovery
4. **🌊 Streaming**: Shows real-time response streaming
5. **🏆 Model Comparison**: Compares different AI models
6. **🎨 UI Integration**: Creates web interface

## 📋 Example Outputs

### Agent Registration
```
🤖 Creating Specialized Agents
==============================
✓ search-specialist registered at /search-agent (Model: openai/gpt-4o-mini, Tools: 2)
✓ support-specialist registered at /support-agent (Model: anthropic/claude-3-5-haiku, Tools: 3)
✓ technical-specialist registered at /tech-agent (Model: xai/grok-beta, Tools: 3)
```

### Feature Testing
```
🧪 Testing All Features
=======================

--- Advanced Search ---
✅ Response: I found several relevant results for 'API documentation' in the docs category...
📊 Agent: comprehensive-ai-assistant | Model: anthropic/claude-3-5-sonnet | Tools Used: ["advanced_search"]

--- User Creation ---
✅ Response: I've successfully created an admin user named Alice Johnson...
📊 Agent: comprehensive-ai-assistant | Model: anthropic/claude-3-5-sonnet | Tools Used: ["create_user"]
```

## 🛠 Tool Parameter Examples

### Complex Search Parameters
```go
type SearchParams struct {
    Query   string        `json:"query" desc:"Search query string" required:"true"`
    Limit   int           `json:"limit,omitempty" desc:"Maximum results" default:"10" min:"1" max:"100"`
    Filters *SearchFilter `json:"filters,omitempty" desc:"Optional search filters"`
}

type SearchFilter struct {
    Category     string   `json:"category,omitempty" desc:"Filter by category" enum:"docs,tickets,users"`
    CreatedAfter string   `json:"created_after,omitempty" desc:"Filter by creation date (ISO format)"`
    Tags         []string `json:"tags,omitempty" desc:"Filter by tags"`
}
```

### User Creation with Settings
```go
type UserParams struct {
    Name     string        `json:"name" desc:"User full name" required:"true"`
    Email    string        `json:"email" desc:"User email address" required:"true"`
    Role     string        `json:"role,omitempty" desc:"User role" enum:"admin,user,guest" default:"user"`
    Active   bool          `json:"active,omitempty" desc:"User active status" default:"true"`
    Settings *UserSettings `json:"settings,omitempty" desc:"User preferences"`
}
```

### Custom Validation
```go
// Implement custom validation for business logic
func (t TicketParams) Validate() error {
    if len(t.Title) < 3 {
        return fmt.Errorf("title must be at least 3 characters long")
    }
    if len(t.Description) < 10 {
        return fmt.Errorf("description must be at least 10 characters long")
    }
    return nil
}
```

## 🌐 Generated Endpoints

After running the demo, you'll have these endpoints available:

- **Main Chat**: `/comprehensive-chat` - Full-featured AI assistant
- **Search Agent**: `/search-agent` - Specialized search assistant
- **Support Agent**: `/support-agent` - Customer support specialist
- **Technical Agent**: `/tech-agent` - Technical operations assistant
- **UI Demo**: `/ai-demo` - Interactive web interface

## 🔍 Key Learning Points

### 1. Agent Architecture
```go
agent := &sourcetool.Agent{
    Name:        "comprehensive-ai-assistant",
    Description: "A comprehensive AI assistant with advanced capabilities",
    Instructions: `Detailed instructions with tool capabilities...`,
    Model:       models.Anthropic("claude-3-5-sonnet", models.WithTemperature(0.7)),
    Tools:       tools,
    MaxSteps:    10,    // Allow multi-step conversations
    Timeout:     30000, // 30 second timeout
}
```

### 2. Tool Creation
```go
// Type-safe tool with complex parameters
tool := agent.NewTool("advanced_search", "Advanced search with filtering",
    func(ctx context.Context, params SearchParams) (interface{}, error) {
        // Implementation with full parameter validation
        return results, nil
    },
)
```

### 3. Router Registration
```go
// Register agent with access control and routing
myAgent := st.Agent("/comprehensive-chat", agent)
```

### 4. Generation with Context
```go
response, err := agent.Generate(
    context.Background(),
    message,
    sourcetool.WithUser(&sourcetool.User{
        ID:    "test-user-123",
        Name:  "Test User", 
        Email: "test@demo.com",
    }),
    sourcetool.WithSessionID("session-456"),
)
```

### 5. Streaming Responses
```go
err := agent.Stream(context.Background(), prompt, func(chunk string) error {
    fmt.Print(chunk)
    return nil
})
```

## 🎛 Configuration Options

### Model Configuration
```go
// Different model providers with specific settings
models.Anthropic("claude-3-5-sonnet", models.WithTemperature(0.7), models.WithMaxTokens(4096))
models.OpenAI("gpt-4o", models.WithTemperature(0.9), models.WithTopP(0.95))
models.XAI("grok-beta", models.WithTopK(50))
models.Google("gemini-pro", models.WithSeed(42))
```

### Agent Settings
```go
agent := &sourcetool.Agent{
    // ... other fields
    MaxSteps: 10,    // Multi-step conversation limit
    Timeout:  30000, // 30 second timeout
}
```

## 🚨 Error Handling Examples

The demo includes comprehensive error handling scenarios:

- Invalid calculation expressions
- Missing required parameters
- Custom validation failures
- Empty search queries
- Network timeouts
- Tool execution errors

## 📊 Performance Monitoring

The demo includes performance measurement:

- Response time tracking
- Token usage monitoring
- Tool execution statistics
- Model comparison metrics

## 🔧 Customization

You can easily customize this demo:

1. **Add New Tools**: Implement additional tool functions
2. **Modify Agents**: Change instructions, models, or tool combinations
3. **Extend UI**: Add more interactive elements
4. **Add Validation**: Implement custom parameter validation
5. **Integration**: Connect to real APIs and databases

## 📚 Further Reading

- [Sourcetool Documentation](https://docs.sourcetool.io)
- [Agent API Reference](https://docs.sourcetool.io/agents)
- [Tool Creation Guide](https://docs.sourcetool.io/tools)
- [Model Configuration](https://docs.sourcetool.io/models)

---

This comprehensive demo serves as a complete reference implementation for building sophisticated AI agent systems with Sourcetool. Use it as a starting point for your own agent implementations! 🚀