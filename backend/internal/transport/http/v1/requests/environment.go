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
	EnvironmentID string  `json:"-" validate:"required"`
	Name          *string `json:"name" validate:"required"`
	Color         *string `json:"color" validate:"required"`
}

type DeleteEnvironmentRequest struct {
	EnvironmentID string `json:"environmentId" validate:"required"`
}
