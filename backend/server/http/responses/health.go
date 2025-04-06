package responses

type HealthResponse struct {
	Status    string            `json:"status"`
	Uptime    string            `json:"uptime"`
	Timestamp string            `json:"timestamp"`
	Details   map[string]string `json:"details"`
}
