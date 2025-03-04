package requests

type GetAPIKeyRequest struct {
	APIKeyID string `json:"-" validate:"required"`
}

type CreateAPIKeyRequest struct {
	EnvironmentID string `json:"environmentId" validate:"required"`
	Name          string `json:"name" validate:"required"`
}

type UpdateAPIKeyRequest struct {
	APIKeyID string  `json:"-" validate:"required"`
	Name     *string `json:"name" validate:"-"`
}

type DeleteAPIKeyRequest struct {
	APIKeyID string `json:"-" validate:"required"`
}
