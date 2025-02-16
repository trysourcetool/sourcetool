package types

type EnvironmentPayload struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Color     string `json:"color"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListEnvironmentsPayload struct {
	Environments []*EnvironmentPayload `json:"environments"`
}

type GetEnvironmentInput struct {
	EnvironmentID string `json:"-" validate:"required"`
}

type GetEnvironmentPayload struct {
	Environment *EnvironmentPayload `json:"environment"`
}

type CreateEnvironmentInput struct {
	Name  string `json:"name" validate:"required"`
	Slug  string `json:"slug" validate:"required"`
	Color string `json:"color" validate:"required"`
}

type CreateEnvironmentPayload struct {
	Environment *EnvironmentPayload `json:"environment"`
}

type UpdateEnvironmentInput struct {
	EnvironmentID string  `json:"-" validate:"required"`
	Name          *string `json:"name" validate:"required"`
	Color         *string `json:"color" validate:"required"`
}

type UpdateEnvironmentPayload struct {
	Environment *EnvironmentPayload `json:"environment"`
}

type DeleteEnvironmentInput struct {
	EnvironmentID string `json:"environmentId" validate:"required"`
}

type DeleteEnvironmentPayload struct {
	Environment *EnvironmentPayload `json:"environment"`
}
