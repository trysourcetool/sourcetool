package types

type APIKeyPayload struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Key         string              `json:"key"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
	Environment *EnvironmentPayload `json:"environment,omitempty"`
}

type GetAPIKeyInput struct {
	APIKeyID string `json:"-" validate:"required"`
}

type GetAPIKeyPayload struct {
	APIKey *APIKeyPayload `json:"apiKey"`
}

type ListAPIKeysPayload struct {
	DevKey   *APIKeyPayload   `json:"devKey"`
	LiveKeys []*APIKeyPayload `json:"liveKeys"`
}

type CreateAPIKeyInput struct {
	EnvironmentID string `json:"environmentId" validate:"required"`
	Name          string `json:"name" validate:"required"`
}

type CreateAPIKeyPayload struct {
	APIKey *APIKeyPayload `json:"apiKey"`
}

type UpdateAPIKeyInput struct {
	APIKeyID string  `json:"-" validate:"required"`
	Name     *string `json:"name" validate:"-"`
}

type UpdateAPIKeyPayload struct {
	APIKey *APIKeyPayload `json:"apiKey"`
}

type DeleteAPIKeyInput struct {
	APIKeyID string `json:"-" validate:"required"`
}

type DeleteAPIKeyPayload struct {
	APIKey *APIKeyPayload `json:"apiKey"`
}
