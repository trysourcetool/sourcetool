package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/trysourcetool/sourcetool-go"
	"github.com/trysourcetool/sourcetool-go/agent"
	"github.com/trysourcetool/sourcetool-go/agent/models"
)

// === Tool Parameter Structures ===

// SearchParams demonstrates complex parameter validation with nested structures
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

// UserParams demonstrates enum validation and required fields
type UserParams struct {
	Name     string        `json:"name" desc:"User full name" required:"true"`
	Email    string        `json:"email" desc:"User email address" required:"true"`
	Role     string        `json:"role,omitempty" desc:"User role" enum:"admin,user,guest" default:"user"`
	Active   bool          `json:"active,omitempty" desc:"User active status" default:"true"`
	Settings *UserSettings `json:"settings,omitempty" desc:"User preferences"`
}

type UserSettings struct {
	Theme         string `json:"theme,omitempty" desc:"UI theme" enum:"light,dark,auto" default:"auto"`
	Notifications bool   `json:"notifications,omitempty" desc:"Enable notifications" default:"true"`
	Language      string `json:"language,omitempty" desc:"Preferred language" default:"en"`
}

// TicketParams demonstrates validation and business logic
type TicketParams struct {
	Title       string   `json:"title" desc:"Ticket title" required:"true"`
	Description string   `json:"description" desc:"Detailed description" required:"true"`
	Priority    string   `json:"priority,omitempty" desc:"Priority level" enum:"low,normal,high,urgent" default:"normal"`
	Category    string   `json:"category,omitempty" desc:"Issue category" enum:"bug,feature,support,question"`
	Tags        []string `json:"tags,omitempty" desc:"Tags for categorization"`
	AssigneeID  string   `json:"assignee_id,omitempty" desc:"ID of assigned user"`
}

// Implement custom validation
func (t TicketParams) Validate() error {
	if len(t.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters long")
	}
	if len(t.Description) < 10 {
		return fmt.Errorf("description must be at least 10 characters long")
	}
	return nil
}

// FileParams demonstrates file operations
type FileParams struct {
	Filename string `json:"filename" desc:"File name" required:"true"`
	Content  string `json:"content" desc:"File content" required:"true"`
	Path     string `json:"path,omitempty" desc:"File path" default:"/tmp"`
}

// CalculationParams demonstrates numeric operations
type CalculationParams struct {
	Expression string `json:"expression" desc:"Mathematical expression" required:"true"`
	Precision  int    `json:"precision,omitempty" desc:"Decimal places" default:"2" min:"0" max:"10"`
}

// WeatherParams demonstrates external API simulation
type WeatherParams struct {
	Location string `json:"location" desc:"Location (city, country)" required:"true"`
	Unit     string `json:"unit,omitempty" desc:"Temperature unit" enum:"celsius,fahrenheit" default:"celsius"`
}

// === Tool Implementations ===

func main() {
	// Initialize Sourcetool with comprehensive configuration
	st := sourcetool.New(&sourcetool.Config{
		APIKey:   "dev_test_123", // Format: environment_key
		Endpoint: "http://localhost:3001",
	})

	fmt.Println("🚀 Starting Comprehensive AI Agent Demo")
	fmt.Println("========================================")

	// Create comprehensive tool suite
	tools := createComprehensiveToolSuite()

	// Demonstrate different agent configurations
	demonstrateAgentTypes(st, tools)

	// Create and test the main multi-purpose agent
	mainAgent := createMainAgent(st, tools)

	// Run comprehensive tests
	testAllFeatures(mainAgent)

	// Demonstrate error handling and edge cases
	demonstrateErrorHandling(mainAgent)

	// Demonstrate streaming capabilities
	demonstrateStreamingFeatures(mainAgent)

	// Show model provider comparisons
	compareModelProviders(st)

	// Create UI integration example
	createUIExample(st, mainAgent)

	fmt.Println("\n✅ Demo completed successfully!")
}

func createComprehensiveToolSuite() []agent.Tool {
	return []agent.Tool{
		// Search tool with complex parameters
		agent.NewTool("advanced_search", "Advanced search with filtering and sorting",
			func(ctx context.Context, params SearchParams) (interface{}, error) {
				results := []map[string]interface{}{}

				// Simulate search with different categories
				categories := []string{"docs", "tickets", "users"}
				if params.Filters != nil && params.Filters.Category != "" {
					categories = []string{params.Filters.Category}
				}

				count := 0
				for _, cat := range categories {
					if count >= params.Limit {
						break
					}
					for i := 1; i <= 3 && count < params.Limit; i++ {
						result := map[string]interface{}{
							"id":       fmt.Sprintf("%s-%d", cat, i),
							"title":    fmt.Sprintf("%s result %d for '%s'", strings.Title(cat), i, params.Query),
							"category": cat,
							"score":    0.9 - float64(i)*0.1,
							"created":  time.Now().AddDate(0, 0, -i).Format("2006-01-02"),
						}

						if params.Filters != nil && len(params.Filters.Tags) > 0 {
							result["tags"] = params.Filters.Tags
						}

						results = append(results, result)
						count++
					}
				}

				return map[string]interface{}{
					"query":     params.Query,
					"results":   results,
					"count":     len(results),
					"total":     len(results) * 2, // Simulate total available
					"filters":   params.Filters,
					"timestamp": time.Now().Format(time.RFC3339),
				}, nil
			},
		),

		// User management tool with validation
		agent.NewTool("create_user", "Create a new user with comprehensive settings",
			func(ctx context.Context, params UserParams) (interface{}, error) {
				userID := fmt.Sprintf("user_%d", time.Now().Unix())

				user := map[string]interface{}{
					"id":         userID,
					"name":       params.Name,
					"email":      params.Email,
					"role":       params.Role,
					"active":     params.Active,
					"created_at": time.Now().Format(time.RFC3339),
				}

				if params.Settings != nil {
					user["settings"] = map[string]interface{}{
						"theme":         params.Settings.Theme,
						"notifications": params.Settings.Notifications,
						"language":      params.Settings.Language,
					}
				}

				return map[string]interface{}{
					"success": true,
					"user":    user,
					"message": fmt.Sprintf("User %s created successfully with ID %s", params.Name, userID),
				}, nil
			},
		),

		// Ticket creation with custom validation
		agent.NewTool("create_ticket", "Create a support ticket with validation",
			func(ctx context.Context, params TicketParams) (interface{}, error) {
				ticketID := fmt.Sprintf("TKT-%d", time.Now().Unix())

				ticket := map[string]interface{}{
					"id":          ticketID,
					"title":       params.Title,
					"description": params.Description,
					"priority":    params.Priority,
					"category":    params.Category,
					"status":      "open",
					"created_at":  time.Now().Format(time.RFC3339),
					"tags":        params.Tags,
				}

				if params.AssigneeID != "" {
					ticket["assignee_id"] = params.AssigneeID
					ticket["assigned_at"] = time.Now().Format(time.RFC3339)
				}

				return map[string]interface{}{
					"success":   true,
					"ticket":    ticket,
					"ticket_id": ticketID,
					"message":   fmt.Sprintf("Ticket '%s' created with ID %s", params.Title, ticketID),
				}, nil
			},
		),

		// File operations tool
		agent.NewTool("file_operations", "Perform file operations",
			func(ctx context.Context, params FileParams) (interface{}, error) {
				// Simulate file operations
				fileSize := len(params.Content)
				fullPath := fmt.Sprintf("%s/%s", strings.TrimSuffix(params.Path, "/"), params.Filename)

				return map[string]interface{}{
					"success":  true,
					"filename": params.Filename,
					"path":     fullPath,
					"size":     fileSize,
					"created":  time.Now().Format(time.RFC3339),
					"message":  fmt.Sprintf("File %s created at %s (%d bytes)", params.Filename, fullPath, fileSize),
				}, nil
			},
		),

		// Mathematical calculator
		agent.NewTool("calculator", "Perform mathematical calculations",
			func(ctx context.Context, params CalculationParams) (interface{}, error) {
				// Simple expression evaluator (in real implementation, use a proper parser)
				expr := strings.ReplaceAll(params.Expression, " ", "")

				// For demo purposes, handle simple operations
				var result float64
				var err error

				if strings.Contains(expr, "+") {
					parts := strings.Split(expr, "+")
					if len(parts) == 2 {
						a, _ := strconv.ParseFloat(parts[0], 64)
						b, _ := strconv.ParseFloat(parts[1], 64)
						result = a + b
					}
				} else if strings.Contains(expr, "*") {
					parts := strings.Split(expr, "*")
					if len(parts) == 2 {
						a, _ := strconv.ParseFloat(parts[0], 64)
						b, _ := strconv.ParseFloat(parts[1], 64)
						result = a * b
					}
				} else {
					result, err = strconv.ParseFloat(expr, 64)
				}

				if err != nil {
					return nil, fmt.Errorf("invalid expression: %s", params.Expression)
				}

				// Format result with specified precision
				formatStr := fmt.Sprintf("%%.%df", params.Precision)
				formattedResult := fmt.Sprintf(formatStr, result)

				return map[string]interface{}{
					"expression": params.Expression,
					"result":     result,
					"formatted":  formattedResult,
					"precision":  params.Precision,
				}, nil
			},
		),

		// Weather API simulation
		agent.NewTool("get_weather", "Get weather information for a location",
			func(ctx context.Context, params WeatherParams) (interface{}, error) {
				// Simulate weather data
				locations := map[string]map[string]interface{}{
					"tokyo":    {"temp": 22, "condition": "sunny", "humidity": 65},
					"london":   {"temp": 15, "condition": "rainy", "humidity": 80},
					"new york": {"temp": 18, "condition": "cloudy", "humidity": 70},
				}

				locationKey := strings.ToLower(strings.Split(params.Location, ",")[0])
				weatherData, exists := locations[locationKey]
				if !exists {
					// Generate random weather for unknown locations
					weatherData = map[string]interface{}{
						"temp":      20 + (time.Now().Unix()%20 - 10),
						"condition": []string{"sunny", "cloudy", "rainy"}[time.Now().Unix()%3],
						"humidity":  50 + (time.Now().Unix() % 40),
					}
				}

				temp := weatherData["temp"].(int64)
				if params.Unit == "fahrenheit" {
					temp = temp*9/5 + 32
				}

				return map[string]interface{}{
					"location":    params.Location,
					"temperature": temp,
					"unit":        params.Unit,
					"condition":   weatherData["condition"],
					"humidity":    weatherData["humidity"],
					"timestamp":   time.Now().Format(time.RFC3339),
				}, nil
			},
		),

		// System status tool (simple tool example)
		agent.SimpleTool("system_status", "Get comprehensive system status",
			func(ctx context.Context) (interface{}, error) {
				return map[string]interface{}{
					"status":           "healthy",
					"uptime":           "5d 12h 34m",
					"version":          "v1.2.3",
					"memory_usage":     "65%",
					"cpu_usage":        "23%",
					"disk_usage":       "45%",
					"active_sessions":  42,
					"requests_per_sec": 125,
					"last_updated":     time.Now().Format(time.RFC3339),
				}, nil
			},
		),

		// Database operations simulation
		agent.NewTool("database_query", "Execute database queries",
			func(ctx context.Context, params struct {
				Query string `json:"query" desc:"SQL query to execute" required:"true"`
				Limit int    `json:"limit,omitempty" desc:"Result limit" default:"100" min:"1" max:"1000"`
			}) (interface{}, error) {
				// Simulate database response
				queryType := "SELECT"
				if strings.HasPrefix(strings.ToUpper(params.Query), "INSERT") {
					queryType = "INSERT"
				} else if strings.HasPrefix(strings.ToUpper(params.Query), "UPDATE") {
					queryType = "UPDATE"
				} else if strings.HasPrefix(strings.ToUpper(params.Query), "DELETE") {
					queryType = "DELETE"
				}

				var result interface{}
				switch queryType {
				case "SELECT":
					result = []map[string]interface{}{
						{"id": 1, "name": "John Doe", "email": "john@example.com"},
						{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
					}
				case "INSERT":
					result = map[string]interface{}{"inserted_id": 123, "rows_affected": 1}
				case "UPDATE":
					result = map[string]interface{}{"rows_affected": 3}
				case "DELETE":
					result = map[string]interface{}{"rows_affected": 1}
				}

				return map[string]interface{}{
					"query":          params.Query,
					"query_type":     queryType,
					"result":         result,
					"execution_time": "12ms",
					"rows_returned":  2,
					"timestamp":      time.Now().Format(time.RFC3339),
				}, nil
			},
		),
	}
}

func createMainAgent(st *sourcetool.Sourcetool, tools []agent.Tool) *sourcetool.Agent {
	aiAgent := &sourcetool.Agent{
		Name:        "comprehensive-ai-assistant",
		Description: "A comprehensive AI assistant with advanced capabilities",
		Instructions: `You are a highly capable AI assistant with access to a wide range of tools. Your capabilities include:

🔍 SEARCH & INFORMATION:
- Advanced search with filtering and categorization
- Weather information for any location
- System status and monitoring

👥 USER MANAGEMENT:
- Create users with detailed profiles and settings
- Manage user roles and permissions

🎫 TICKET SYSTEM:
- Create and manage support tickets
- Categorize and prioritize issues
- Assign tickets to team members

📁 FILE OPERATIONS:
- Create and manage files
- File content operations

🧮 CALCULATIONS:
- Mathematical calculations with precision control
- Expression evaluation

💾 DATABASE:
- Execute queries and data operations
- Handle various SQL operations

BEHAVIOR GUIDELINES:
- Always be helpful, accurate, and professional
- Use the most appropriate tool for each task
- Provide detailed explanations of your actions
- Handle errors gracefully and suggest alternatives
- Ask for clarification when needed
- Validate parameters before tool execution

When users request help, analyze their needs and use the appropriate tools to assist them effectively.`,
		Model: models.Anthropic("claude-3-5-sonnet",
			models.WithTemperature(0.7),
			models.WithMaxTokens(4096),
		),
		Tools:    tools,
		MaxSteps: 10,    // Allow multi-step conversations
		Timeout:  30000, // 30 second timeout
	}

	// Register the agent using router with access control
	return st.Agent("/comprehensive-chat", aiAgent)
}

func demonstrateAgentTypes(st *sourcetool.Sourcetool, tools []agent.Tool) {
	fmt.Println("\n🤖 Creating Specialized Agents")
	fmt.Println("==============================")

	agents := []struct {
		name         string
		path         string
		description  string
		instructions string
		tools        []agent.Tool
		model        models.Model
	}{
		{
			name:         "search-specialist",
			path:         "/search-agent",
			description:  "Specialized search and information retrieval agent",
			instructions: "You are a search specialist. Help users find information efficiently using advanced search capabilities.",
			tools:        []agent.Tool{tools[0], tools[6]}, // search + weather
			model:        models.OpenAI("gpt-4o-mini", models.WithTemperature(0.3)),
		},
		{
			name:         "support-specialist",
			path:         "/support-agent",
			description:  "Customer support and ticket management specialist",
			instructions: "You are a customer support specialist. Help users create tickets, manage accounts, and resolve issues.",
			tools:        []agent.Tool{tools[1], tools[2], tools[7]}, // user + ticket + status
			model:        models.Anthropic("claude-3-5-haiku", models.WithTemperature(0.5)),
		},
		{
			name:         "technical-specialist",
			path:         "/tech-agent",
			description:  "Technical operations and calculations specialist",
			instructions: "You are a technical specialist. Help users with calculations, file operations, and system queries.",
			tools:        []agent.Tool{tools[3], tools[4], tools[8]}, // file + calc + database
			model:        models.XAI("grok-beta", models.WithTemperature(0.4)),
		},
	}

	for _, spec := range agents {
		agent := &sourcetool.Agent{
			Name:         spec.name,
			Description:  spec.description,
			Instructions: spec.instructions,
			Model:        spec.model,
			Tools:        spec.tools,
			MaxSteps:     5,
		}

		registeredAgent := st.Agent(spec.path, agent)
		fmt.Printf("✓ %s registered at %s (Model: %s, Tools: %d)\n",
			registeredAgent.GetName(), spec.path, spec.model.ID(), len(spec.tools))
	}
}

func testAllFeatures(agent *sourcetool.Agent) {
	fmt.Println("\n🧪 Testing All Features")
	fmt.Println("=======================")

	testCases := []struct {
		name        string
		message     string
		expectTools []string
	}{
		{
			name:        "Advanced Search",
			message:     "Search for 'API documentation' in docs category with tags ['rest', 'api']",
			expectTools: []string{"advanced_search"},
		},
		{
			name:        "User Creation",
			message:     "Create a new admin user named Alice Johnson with email alice@company.com, dark theme, and notifications enabled",
			expectTools: []string{"create_user"},
		},
		{
			name:        "Ticket Management",
			message:     "Create a high priority bug ticket titled 'Login page not loading' with description 'Users report that the login page fails to load on mobile devices'",
			expectTools: []string{"create_ticket"},
		},
		{
			name:        "File Operations",
			message:     "Create a config file named 'app.json' with content '{\"debug\": true}' in the /etc directory",
			expectTools: []string{"file_operations"},
		},
		{
			name:        "Mathematical Calculation",
			message:     "Calculate 15.5 * 2.3 with 3 decimal places precision",
			expectTools: []string{"calculator"},
		},
		{
			name:        "Weather Information",
			message:     "What's the weather like in Tokyo in Fahrenheit?",
			expectTools: []string{"get_weather"},
		},
		{
			name:        "System Status",
			message:     "Check the current system status and performance metrics",
			expectTools: []string{"system_status"},
		},
		{
			name:        "Database Query",
			message:     "Execute a SELECT query to get all users from the database",
			expectTools: []string{"database_query"},
		},
		{
			name:        "Multi-step Task",
			message:     "Create a user named Bob Smith, then create a ticket assigned to them about a feature request",
			expectTools: []string{"create_user", "create_ticket"},
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\n--- %s ---\n", tc.name)

		response, err := agent.Generate(
			context.Background(),
			tc.message,
			sourcetool.WithUser(&sourcetool.User{
				ID:    "test-user-123",
				Name:  "Test User",
				Email: "test@demo.com",
			}),
			sourcetool.WithSessionID(fmt.Sprintf("test-session-%d", time.Now().Unix())),
		)

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("✅ Response: %s\n", truncateString(response.Message, 200))

		if response.Metadata != nil {
			fmt.Printf("📊 Agent: %v | Model: %v/%v | Tools Used: %v\n",
				response.Metadata["agent_name"],
				response.Metadata["model_provider"],
				response.Metadata["model_name"],
				response.Metadata["tools_used"])
		}

		// Show suggested actions if available
		if len(response.SuggestedActions) > 0 {
			fmt.Printf("💡 Suggested Actions:\n")
			for _, action := range response.SuggestedActions {
				fmt.Printf("   - %s: %s\n", action.Label, action.Description)
			}
		}
	}
}

func demonstrateErrorHandling(agent *sourcetool.Agent) {
	fmt.Println("\n⚠️  Error Handling & Edge Cases")
	fmt.Println("================================")

	errorCases := []struct {
		name    string
		message string
	}{
		{
			name:    "Invalid Calculation",
			message: "Calculate 'invalid expression' with 2 decimal places",
		},
		{
			name:    "Missing Required Parameters",
			message: "Create a user without providing a name or email",
		},
		{
			name:    "Validation Failure",
			message: "Create a ticket with title 'x' and description 'short'",
		},
		{
			name:    "Empty Search Query",
			message: "Search for '' with no parameters",
		},
	}

	for _, tc := range errorCases {
		fmt.Printf("\n--- %s ---\n", tc.name)

		response, err := agent.Generate(
			context.Background(),
			tc.message,
			sourcetool.WithUser(&sourcetool.User{
				ID:   "test-user-error",
				Name: "Error Test User",
			}),
		)

		if err != nil {
			fmt.Printf("❌ Expected Error: %v\n", err)
		} else {
			fmt.Printf("🔄 Handled Gracefully: %s\n", truncateString(response.Message, 150))
		}
	}
}

func demonstrateStreamingFeatures(agent *sourcetool.Agent) {
	fmt.Println("\n🌊 Streaming Capabilities")
	fmt.Println("=========================")

	streamingTests := []string{
		"Tell me a detailed story about a robot learning to paint",
		"Explain how machine learning works in simple terms",
		"Write a step-by-step guide for setting up a web server",
	}

	for i, prompt := range streamingTests {
		fmt.Printf("\n--- Streaming Test %d ---\n", i+1)
		fmt.Printf("Prompt: %s\n", prompt)
		fmt.Printf("Stream: ")

		var fullResponse strings.Builder
		err := agent.Stream(
			context.Background(),
			prompt,
			func(chunk string) error {
				fmt.Print(chunk)
				fullResponse.WriteString(chunk)
				return nil
			},
		)

		fmt.Println() // New line after streaming

		if err != nil {
			fmt.Printf("❌ Streaming Error: %v\n", err)
		} else {
			fmt.Printf("✅ Streamed %d characters successfully\n", fullResponse.Len())
		}
	}
}

func compareModelProviders(st *sourcetool.Sourcetool) {
	fmt.Println("\n🏆 Model Provider Comparison")
	fmt.Println("=============================")

	providers := []struct {
		name  string
		model models.Model
	}{
		{"Anthropic Claude", models.Anthropic("claude-3-5-sonnet", models.WithTemperature(0.7))},
		{"OpenAI GPT", models.OpenAI("gpt-4o", models.WithTemperature(0.7))},
		{"xAI Grok", models.XAI("grok-beta", models.WithTemperature(0.7))},
		{"Google Gemini", models.Google("gemini-pro", models.WithTemperature(0.7))},
	}

	testPrompt := "Explain quantum computing in exactly 3 sentences."

	for _, provider := range providers {
		fmt.Printf("\n--- %s ---\n", provider.name)

		agent := &sourcetool.Agent{
			Name:         fmt.Sprintf("comparison-%s", strings.ToLower(strings.ReplaceAll(provider.name, " ", "-"))),
			Description:  fmt.Sprintf("Test agent using %s", provider.name),
			Instructions: "You are a helpful AI assistant. Respond clearly and concisely.",
			Model:        provider.model,
		}

		testAgent := st.Agent(fmt.Sprintf("/test-%s", agent.Name), agent)

		start := time.Now()
		response, err := testAgent.Generate(context.Background(), testPrompt)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("✅ Model: %s\n", provider.model.ID())
		fmt.Printf("⏱️  Response Time: %v\n", duration)
		fmt.Printf("📝 Response: %s\n", truncateString(response.Message, 250))

		if response.Metadata != nil {
			fmt.Printf("📊 Tokens: %v\n", response.Metadata["tokens_used"])
		}
	}
}

func createUIExample(st *sourcetool.Sourcetool, agent *sourcetool.Agent) {
	fmt.Println("\n🎨 Creating UI Integration Example")
	fmt.Println("===================================")

	// Create a UI page that integrates with the agent
	st.Page("/ai-demo", "AI Agent Demo", func(ui sourcetool.UIBuilder) error {
		ui.Markdown("# 🤖 Comprehensive AI Agent Demo")
		ui.Markdown("This demo showcases the full capabilities of the AI agent system.")

		// Agent info section
		ui.Markdown("## Agent Information")
		ui.Markdown(fmt.Sprintf("**Name:** %s", agent.GetName()))
		ui.Markdown(fmt.Sprintf("**Model:** %s", agent.Model.ID()))
		ui.Markdown(fmt.Sprintf("**Tools Available:** %d", len(agent.ListTools())))

		// Interactive chat section
		ui.Markdown("## Interactive Chat")
		userInput := ui.TextInput("Ask the AI assistant anything...")

		if ui.Button("Send Message") && userInput != "" {
			// In a real implementation, this would call agent.Generate()
			ui.Markdown(fmt.Sprintf("**You:** %s", userInput))
			ui.Markdown("**AI:** I would help you with that request using my available tools!")
		}

		// Tool showcase
		ui.Markdown("## Available Tools")
		tools := agent.ListTools()
		for _, tool := range tools {
			ui.Markdown(fmt.Sprintf("- **%s**: %s", tool.Name, tool.Description))
		}

		// Quick actions
		ui.Markdown("## Quick Actions")
		if ui.Button("Check System Status") {
			ui.Markdown("✅ System status check initiated!")
		}
		if ui.Button("Create Sample User") {
			ui.Markdown("👤 Sample user creation started!")
		}
		if ui.Button("Run Search Demo") {
			ui.Markdown("🔍 Search demonstration launched!")
		}

		return nil
	})

	fmt.Printf("✅ UI page created at /ai-demo\n")
	fmt.Printf("🌐 Agent endpoints:\n")
	fmt.Printf("   - Main Chat: /comprehensive-chat\n")
	fmt.Printf("   - Search Agent: /search-agent\n")
	fmt.Printf("   - Support Agent: /support-agent\n")
	fmt.Printf("   - Technical Agent: /tech-agent\n")
}

// === Utility Functions ===

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// === Todo: Mark first task as complete ===
