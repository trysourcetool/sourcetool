package responses

type HostInstanceResponse struct {
	ID         string `json:"id"`
	SDKName    string `json:"sdkName"`
	SDKVersion string `json:"sdkVersion"`
	Status     string `json:"status"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type PingHostInstanceResponse struct {
	HostInstance *HostInstanceResponse `json:"hostInstance"`
}
