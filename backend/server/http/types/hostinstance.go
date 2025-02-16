package types

type HostInstancePayload struct {
	ID         string `json:"id"`
	SDKName    string `json:"sdkName"`
	SDKVersion string `json:"sdkVersion"`
	Status     string `json:"status"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type PingHostInstanceInput struct {
	PageID *string `validate:"-"`
}

type PingHostInstancePayload struct {
	HostInstance *HostInstancePayload `json:"hostInstance"`
}
