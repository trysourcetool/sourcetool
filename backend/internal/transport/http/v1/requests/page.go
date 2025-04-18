package requests

type ListPagesRequest struct {
	EnvironmentID string `json:"environment_id" validate:"required"`
}
