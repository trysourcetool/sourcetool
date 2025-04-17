package input

// GetEnvironmentInput is the input for Get operation.
type GetEnvironmentInput struct {
	EnvironmentID string
}

// CreateEnvironmentInput is the input for Create operation.
type CreateEnvironmentInput struct {
	Name  string
	Slug  string
	Color string
}

// UpdateEnvironmentInput is the input for Update operation.
type UpdateEnvironmentInput struct {
	EnvironmentID string
	Name          *string
	Color         *string
}

// DeleteEnvironmentInput is the input for Delete operation.
type DeleteEnvironmentInput struct {
	EnvironmentID string
}
