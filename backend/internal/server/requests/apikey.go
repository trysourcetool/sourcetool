package requests

type CreateAPIKeyRequest struct {
	EnvironmentID string `json:"environmentId" validate:"required"`
	Name          string `json:"name" validate:"required"`
}

type UpdateAPIKeyRequest struct {
	Name *string `json:"name" validate:"-"`
}
