package input

// GetGroupInput is the input for Get operation.
type GetGroupInput struct {
	GroupID string
}

// CreateGroupInput is the input for Create operation.
type CreateGroupInput struct {
	Name    string
	Slug    string
	UserIDs []string
}

// UpdateGroupInput is the input for Update operation.
type UpdateGroupInput struct {
	GroupID string
	Name    *string
	UserIDs []string
}

// DeleteGroupInput is the input for Delete operation.
type DeleteGroupInput struct {
	GroupID string
}
