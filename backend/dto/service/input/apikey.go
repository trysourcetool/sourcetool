package input

// GetAPIKeyInput is the input for Get operation.
type GetAPIKeyInput struct {
	APIKeyID string
}

// CreateAPIKeyInput is the input for Create operation.
type CreateAPIKeyInput struct {
	EnvironmentID string
	Name          string
}

// UpdateAPIKeyInput is the input for Update operation.
type UpdateAPIKeyInput struct {
	APIKeyID string
	Name     *string
}

// DeleteAPIKeyInput is the input for Delete operation.
type DeleteAPIKeyInput struct {
	APIKeyID string
}
