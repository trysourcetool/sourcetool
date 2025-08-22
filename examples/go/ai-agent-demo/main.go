package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/trysourcetool/sourcetool-go"
	"github.com/trysourcetool/sourcetool-go/agent"
	"github.com/trysourcetool/sourcetool-go/agent/models"
)

// === Practical Tool Parameter Structures ===

// WebhookEvent represents incoming webhook data
type WebhookEvent struct {
	EventType string                 `json:"event_type" desc:"Type of webhook event" required:"true"`
	Payload   map[string]interface{} `json:"payload" desc:"Event payload data" required:"true"`
	Source    string                 `json:"source" desc:"Source system" required:"true"`
	Timestamp string                 `json:"timestamp,omitempty" desc:"Event timestamp"`
}

// SlackNotification for sending Slack messages
type SlackNotification struct {
	Channel  string `json:"channel" desc:"Slack channel (#general, @user)" required:"true"`
	Message  string `json:"message" desc:"Message content" required:"true"`
	Priority string `json:"priority,omitempty" desc:"Message priority" enum:"low,normal,high,urgent" default:"normal"`
	Mention  string `json:"mention,omitempty" desc:"User to mention (@here, @channel, @username)"`
}

// DataAnalysis for processing data
type DataAnalysis struct {
	DataSource string                 `json:"data_source" desc:"Data source identifier" required:"true"`
	Query      string                 `json:"query" desc:"Analysis query or filter" required:"true"`
	Format     string                 `json:"format,omitempty" desc:"Output format" enum:"json,csv,summary" default:"json"`
	TimeRange  *TimeRange             `json:"time_range,omitempty" desc:"Time range for analysis"`
	Metrics    []string               `json:"metrics,omitempty" desc:"Metrics to calculate"`
	GroupBy    string                 `json:"group_by,omitempty" desc:"Field to group results by"`
	Filters    map[string]interface{} `json:"filters,omitempty" desc:"Additional filters"`
}

type TimeRange struct {
	Start string `json:"start" desc:"Start time (ISO format)" required:"true"`
	End   string `json:"end" desc:"End time (ISO format)" required:"true"`
}

// APIMonitor for monitoring external APIs
type APIMonitor struct {
	URL            string            `json:"url" desc:"API endpoint URL" required:"true"`
	Method         string            `json:"method,omitempty" desc:"HTTP method" enum:"GET,POST,PUT,DELETE" default:"GET"`
	Headers        map[string]string `json:"headers,omitempty" desc:"Request headers"`
	ExpectedStatus int               `json:"expected_status,omitempty" desc:"Expected HTTP status code" default:"200"`
	Timeout        int               `json:"timeout,omitempty" desc:"Request timeout in seconds" default:"30"`
}

// TaskSchedule for scheduling tasks
type TaskSchedule struct {
	TaskName    string                 `json:"task_name" desc:"Name of the task" required:"true"`
	Schedule    string                 `json:"schedule" desc:"Cron expression or interval" required:"true"`
	Action      string                 `json:"action" desc:"Action to perform" required:"true"`
	Parameters  map[string]interface{} `json:"parameters,omitempty" desc:"Task parameters"`
	MaxRetries  int                    `json:"max_retries,omitempty" desc:"Maximum retry attempts" default:"3"`
	NotifyOnErr bool                   `json:"notify_on_error,omitempty" desc:"Send notification on error" default:"true"`
}

// LogQuery for searching logs
type LogQuery struct {
	Query     string   `json:"query" desc:"Search query" required:"true"`
	Level     string   `json:"level,omitempty" desc:"Log level filter" enum:"debug,info,warn,error,fatal"`
	Service   string   `json:"service,omitempty" desc:"Service name filter"`
	TimeRange string   `json:"time_range,omitempty" desc:"Time range (1h, 24h, 7d)" default:"1h"`
	Limit     int      `json:"limit,omitempty" desc:"Maximum results" default:"100" max:"1000"`
	Fields    []string `json:"fields,omitempty" desc:"Fields to include in results"`
}

// DatabaseBackup for database operations
type DatabaseBackup struct {
	Database   string `json:"database" desc:"Database name" required:"true"`
	BackupType string `json:"backup_type,omitempty" desc:"Backup type" enum:"full,incremental,snapshot" default:"full"`
	Compress   bool   `json:"compress,omitempty" desc:"Compress backup" default:"true"`
	Encrypt    bool   `json:"encrypt,omitempty" desc:"Encrypt backup" default:"true"`
	Location   string `json:"location,omitempty" desc:"Backup storage location" default:"s3://backups/"`
}

// === Main Application ===

func main() {
	// Initialize Sourcetool
	st := sourcetool.New(&sourcetool.Config{
		APIKey:   "development_GEBWX4LsqzGRMBI0orlaMNb5tnycTeXLGEBWX4LsqzGRMBI0orl",
		Endpoint: "ws://localhost:3000",
	})

	log.Println("🚀 Starting Practical AI Agent Service")
	log.Println("========================================")

	// Create and register practical agents
	agents := createPracticalAgents(st)

	// Register all agents
	for _, agent := range agents {
		st.Agent(agent.Name, agent)
	}

	if err := st.Listen(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// === Practical Agents Creation ===

func createPracticalAgents(st *sourcetool.Sourcetool) []*sourcetool.Agent {
	agents := []*sourcetool.Agent{
		// Webhook Processing Agent
		createWebhookAgent(st),

		// Data Analysis Agent
		createDataAnalysisAgent(st),

		// DevOps Agent
		createDevOpsAgent(st),

		// Customer Support Agent
		createSupportAgent(st),

		// System Monitor Agent
		createMonitorAgent(st),
	}

	return agents
}

func createWebhookAgent(st *sourcetool.Sourcetool) *sourcetool.Agent {
	tools := []agent.Tool{
		// Process webhook events
		agent.NewTool("process_webhook", "Process incoming webhook events",
			func(ctx context.Context, params WebhookEvent) (interface{}, error) {
				log.Printf("Processing webhook: %s from %s", params.EventType, params.Source)

				// Process based on event type
				var result map[string]interface{}
				switch params.EventType {
				case "deployment":
					result = handleDeploymentEvent(params.Payload)
				case "alert":
					result = handleAlertEvent(params.Payload)
				case "user_signup":
					result = handleUserSignupEvent(params.Payload)
				default:
					result = map[string]interface{}{
						"status": "processed",
						"type":   params.EventType,
						"source": params.Source,
					}
				}

				return map[string]interface{}{
					"success":   true,
					"event_id":  fmt.Sprintf("evt_%d", time.Now().Unix()),
					"processed": result,
					"timestamp": time.Now().Format(time.RFC3339),
				}, nil
			},
		),

		// Send Slack notifications
		agent.NewTool("send_slack", "Send notification to Slack",
			func(ctx context.Context, params SlackNotification) (interface{}, error) {
				log.Printf("Sending Slack notification to %s: %s", params.Channel, truncateString(params.Message, 50))

				// In production, this would use actual Slack API
				notification := map[string]interface{}{
					"channel":   params.Channel,
					"message":   params.Message,
					"priority":  params.Priority,
					"timestamp": time.Now().Format(time.RFC3339),
				}

				if params.Mention != "" {
					notification["mention"] = params.Mention
				}

				return map[string]interface{}{
					"success":      true,
					"message_id":   fmt.Sprintf("slack_%d", time.Now().Unix()),
					"notification": notification,
				}, nil
			},
		),

		// Log events
		agent.SimpleTool("log_event", "Log webhook events for audit",
			func(ctx context.Context) (interface{}, error) {
				return map[string]interface{}{
					"logged":    true,
					"log_id":    fmt.Sprintf("log_%d", time.Now().Unix()),
					"timestamp": time.Now().Format(time.RFC3339),
				}, nil
			},
		),
	}

	return &sourcetool.Agent{
		Name:        "webhook-processor",
		Description: "Webhook processing and notification agent",
		Instructions: `You are a webhook processing agent responsible for:
- Processing incoming webhook events from various sources
- Routing events to appropriate handlers
- Sending notifications via Slack for important events
- Logging all events for audit purposes
- Handling errors gracefully and retrying when necessary

Always acknowledge webhook receipt quickly and process asynchronously when needed.`,
		Tools: tools,
		Model: models.OpenAI("gpt-4o-mini", models.WithTemperature(0.3)),
	}
}

func createDataAnalysisAgent(st *sourcetool.Sourcetool) *sourcetool.Agent {
	tools := []agent.Tool{
		// Analyze data
		agent.NewTool("analyze_data", "Perform data analysis",
			func(ctx context.Context, params DataAnalysis) (interface{}, error) {
				log.Printf("Analyzing data from %s with query: %s", params.DataSource, params.Query)

				// Simulate data analysis
				results := generateSampleAnalysisResults(params)

				return map[string]interface{}{
					"success":     true,
					"data_source": params.DataSource,
					"results":     results,
					"format":      params.Format,
					"timestamp":   time.Now().Format(time.RFC3339),
				}, nil
			},
		),

		// Generate reports
		agent.SimpleTool("generate_report", "Generate analysis report",
			func(ctx context.Context) (interface{}, error) {
				return map[string]interface{}{
					"report_id": fmt.Sprintf("report_%d", time.Now().Unix()),
					"status":    "generated",
					"format":    "pdf",
					"location":  "/reports/latest.pdf",
					"timestamp": time.Now().Format(time.RFC3339),
				}, nil
			},
		),

		// Export data
		agent.SimpleTool("export_data", "Export analysis results",
			func(ctx context.Context) (interface{}, error) {
				return map[string]interface{}{
					"export_id": fmt.Sprintf("export_%d", time.Now().Unix()),
					"format":    "csv",
					"location":  "/exports/data.csv",
					"rows":      1250,
					"size":      "2.4MB",
				}, nil
			},
		),
	}

	return &sourcetool.Agent{
		Name:        "data-analyst",
		Description: "Data analysis and reporting agent",
		Instructions: `You are a data analysis agent specializing in:
- Processing and analyzing large datasets
- Generating insights and trends
- Creating reports and visualizations
- Exporting data in various formats
- Performing statistical calculations

Focus on accuracy and provide clear, actionable insights.`,
		Tools: tools,
		Model: models.Anthropic("claude-3-5-sonnet", models.WithTemperature(0.2)),
	}
}

func createDevOpsAgent(st *sourcetool.Sourcetool) *sourcetool.Agent {
	tools := []agent.Tool{
		// Monitor APIs
		agent.NewTool("monitor_api", "Monitor API health",
			func(ctx context.Context, params APIMonitor) (interface{}, error) {
				log.Printf("Monitoring API: %s %s", params.Method, params.URL)

				// Simulate API monitoring
				status := "healthy"
				responseTime := 150 + time.Now().Unix()%200

				if responseTime > 300 {
					status = "slow"
				}

				return map[string]interface{}{
					"url":           params.URL,
					"status":        status,
					"response_time": fmt.Sprintf("%dms", responseTime),
					"status_code":   params.ExpectedStatus,
					"timestamp":     time.Now().Format(time.RFC3339),
				}, nil
			},
		),

		// Schedule tasks
		agent.NewTool("schedule_task", "Schedule automated tasks",
			func(ctx context.Context, params TaskSchedule) (interface{}, error) {
				log.Printf("Scheduling task: %s with schedule: %s", params.TaskName, params.Schedule)

				return map[string]interface{}{
					"task_id":     fmt.Sprintf("task_%d", time.Now().Unix()),
					"task_name":   params.TaskName,
					"schedule":    params.Schedule,
					"status":      "scheduled",
					"next_run":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
					"max_retries": params.MaxRetries,
				}, nil
			},
		),

		// Backup database
		agent.NewTool("backup_database", "Create database backup",
			func(ctx context.Context, params DatabaseBackup) (interface{}, error) {
				log.Printf("Creating %s backup for database: %s", params.BackupType, params.Database)

				backupSize := "245MB"
				if params.Compress {
					backupSize = "48MB"
				}

				return map[string]interface{}{
					"backup_id":  fmt.Sprintf("backup_%d", time.Now().Unix()),
					"database":   params.Database,
					"type":       params.BackupType,
					"size":       backupSize,
					"compressed": params.Compress,
					"encrypted":  params.Encrypt,
					"location":   params.Location,
					"completed":  time.Now().Format(time.RFC3339),
				}, nil
			},
		),

		// Query logs
		agent.NewTool("query_logs", "Search application logs",
			func(ctx context.Context, params LogQuery) (interface{}, error) {
				log.Printf("Querying logs: %s (level: %s, service: %s)", params.Query, params.Level, params.Service)

				// Simulate log search results
				logs := generateSampleLogs(params)

				return map[string]interface{}{
					"query":      params.Query,
					"matches":    len(logs),
					"logs":       logs,
					"time_range": params.TimeRange,
					"timestamp":  time.Now().Format(time.RFC3339),
				}, nil
			},
		),
	}

	return &sourcetool.Agent{
		Name:        "devops-assistant",
		Description: "DevOps automation and monitoring agent",
		Instructions: `You are a DevOps agent responsible for:
- Monitoring API health and performance
- Scheduling and managing automated tasks
- Creating and managing database backups
- Searching and analyzing application logs
- Alerting on system issues
- Performing routine maintenance tasks

Prioritize system stability and quick issue resolution.`,
		Tools: tools,
		Model: models.OpenAI("gpt-4o", models.WithTemperature(0.2)),
	}
}

func createSupportAgent(st *sourcetool.Sourcetool) *sourcetool.Agent {
	tools := []agent.Tool{
		// Create support ticket
		agent.SimpleTool("create_ticket", "Create customer support ticket",
			func(ctx context.Context) (interface{}, error) {
				return map[string]interface{}{
					"ticket_id": fmt.Sprintf("TICKET-%d", time.Now().Unix()),
					"status":    "open",
					"priority":  "normal",
					"assigned":  "support-team",
					"created":   time.Now().Format(time.RFC3339),
				}, nil
			},
		),

		// Search knowledge base
		agent.SimpleTool("search_kb", "Search knowledge base",
			func(ctx context.Context) (interface{}, error) {
				return map[string]interface{}{
					"articles": []map[string]string{
						{"id": "KB001", "title": "Getting Started Guide"},
						{"id": "KB002", "title": "Common Issues and Solutions"},
						{"id": "KB003", "title": "API Documentation"},
					},
					"total": 3,
				}, nil
			},
		),

		// Send customer email
		agent.SimpleTool("send_email", "Send email to customer",
			func(ctx context.Context) (interface{}, error) {
				return map[string]interface{}{
					"email_id": fmt.Sprintf("email_%d", time.Now().Unix()),
					"status":   "sent",
					"queued":   false,
				}, nil
			},
		),
	}

	return &sourcetool.Agent{
		Name:        "support-agent",
		Description: "Customer support automation agent",
		Instructions: `You are a customer support agent that helps with:
- Creating and managing support tickets
- Searching the knowledge base for solutions
- Sending automated responses to customers
- Escalating issues when necessary
- Tracking customer satisfaction

Always be helpful, professional, and empathetic.`,
		Tools: tools,
		Model: models.Anthropic("claude-3-5-haiku", models.WithTemperature(0.5)),
	}
}

func createMonitorAgent(st *sourcetool.Sourcetool) *sourcetool.Agent {
	tools := []agent.Tool{
		// System metrics
		agent.SimpleTool("get_metrics", "Get system metrics",
			func(ctx context.Context) (interface{}, error) {
				return map[string]interface{}{
					"cpu_usage":    "45%",
					"memory_usage": "62%",
					"disk_usage":   "71%",
					"network_io":   "125 MB/s",
					"uptime":       "15d 4h 23m",
					"timestamp":    time.Now().Format(time.RFC3339),
				}, nil
			},
		),

		// Check services
		agent.SimpleTool("check_services", "Check service health",
			func(ctx context.Context) (interface{}, error) {
				return map[string]interface{}{
					"services": map[string]string{
						"api":      "healthy",
						"database": "healthy",
						"cache":    "degraded",
						"queue":    "healthy",
					},
					"overall": "operational",
				}, nil
			},
		),

		// Alert management
		agent.SimpleTool("manage_alerts", "Manage system alerts",
			func(ctx context.Context) (interface{}, error) {
				return map[string]interface{}{
					"active_alerts": 2,
					"alerts": []map[string]string{
						{"id": "ALT001", "severity": "warning", "message": "High memory usage"},
						{"id": "ALT002", "severity": "info", "message": "Scheduled maintenance"},
					},
				}, nil
			},
		),
	}

	return &sourcetool.Agent{
		Name:        "system-monitor",
		Description: "System monitoring and alerting agent",
		Instructions: `You are a system monitoring agent that:
- Tracks system metrics and performance
- Monitors service health
- Manages alerts and notifications
- Provides real-time system status
- Predicts potential issues

Focus on proactive monitoring and early issue detection.`,
		Tools: tools,
		Model: models.OpenAI("gpt-4o-mini", models.WithTemperature(0.1)),
	}
}

// === Helper Functions ===

func handleDeploymentEvent(payload map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":      "deployment_processed",
		"environment": payload["environment"],
		"version":     payload["version"],
		"status":      "success",
	}
}

func handleAlertEvent(payload map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":   "alert_processed",
		"severity": payload["severity"],
		"message":  payload["message"],
		"notified": true,
	}
}

func handleUserSignupEvent(payload map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":             "user_signup_processed",
		"user_id":            payload["user_id"],
		"welcome_sent":       true,
		"onboarding_started": true,
	}
}

func generateSampleAnalysisResults(params DataAnalysis) interface{} {
	if params.Format == "summary" {
		return map[string]interface{}{
			"total_records": 15234,
			"avg_value":     127.45,
			"max_value":     892.31,
			"min_value":     12.08,
			"trend":         "increasing",
		}
	}

	// Return sample data rows
	return []map[string]interface{}{
		{"id": 1, "value": 123.45, "category": "A", "date": "2024-01-15"},
		{"id": 2, "value": 234.56, "category": "B", "date": "2024-01-15"},
		{"id": 3, "value": 345.67, "category": "A", "date": "2024-01-14"},
	}
}

func generateSampleLogs(params LogQuery) []map[string]interface{} {
	logs := []map[string]interface{}{}
	levels := []string{"info", "warn", "error"}

	for i := 0; i < 5 && i < params.Limit; i++ {
		level := params.Level
		if level == "" {
			level = levels[i%3]
		}

		logs = append(logs, map[string]interface{}{
			"timestamp": time.Now().Add(-time.Duration(i) * time.Minute).Format(time.RFC3339),
			"level":     level,
			"service":   params.Service,
			"message":   fmt.Sprintf("Sample log entry matching query: %s", params.Query),
			"trace_id":  fmt.Sprintf("trace_%d", time.Now().Unix()+int64(i)),
		})
	}

	return logs
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
