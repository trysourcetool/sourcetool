package responses

type APIKeyResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Key         string               `json:"key"`
	CreatedAt   string               `json:"createdAt"`
	UpdatedAt   string               `json:"updatedAt"`
	Environment *EnvironmentResponse `json:"environment,omitempty"`
}

type GetAPIKeyResponse struct {
	APIKey *APIKeyResponse `json:"apiKey"`
}

type ListAPIKeysResponse struct {
	DevKey   *APIKeyResponse   `json:"devKey"`
	LiveKeys []*APIKeyResponse `json:"liveKeys"`
}

type CreateAPIKeyResponse struct {
	APIKey *APIKeyResponse `json:"apiKey"`
}

type UpdateAPIKeyResponse struct {
	APIKey *APIKeyResponse `json:"apiKey"`
}

type DeleteAPIKeyResponse struct {
	APIKey *APIKeyResponse `json:"apiKey"`
}
