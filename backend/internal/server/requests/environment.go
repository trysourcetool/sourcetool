package requests

type GetEnvironmentRequest struct {
	EnvironmentID string `json:"-" validate:"required"`
}

type CreateEnvironmentRequest struct {
	Name  string `json:"name" validate:"required"`
	Slug  string `json:"slug" validate:"required"`
	Color string `json:"color" validate:"required"`
}

type UpdateEnvironmentRequest struct {
	Name  *string `json:"name" validate:"required"`
	Color *string `json:"color" validate:"required"`
}
