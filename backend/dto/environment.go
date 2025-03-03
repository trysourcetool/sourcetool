package dto

import (
	"github.com/trysourcetool/sourcetool/backend/model"
)

// Environment represents environment data in DTOs
type Environment struct {
	ID        string
	Name      string
	Slug      string
	Color     string
	CreatedAt int64
	UpdatedAt int64
}

// EnvironmentFromModel converts from model.Environment to dto.Environment
func EnvironmentFromModel(env *model.Environment) *Environment {
	if env == nil {
		return nil
	}

	return &Environment{
		ID:        env.ID.String(),
		Name:      env.Name,
		Slug:      env.Slug,
		Color:     env.Color,
		CreatedAt: env.CreatedAt.Unix(),
		UpdatedAt: env.UpdatedAt.Unix(),
	}
}

// GetEnvironmentInput is the input for Get operation
type GetEnvironmentInput struct {
	EnvironmentID string
}

// GetEnvironmentOutput is the output for Get operation
type GetEnvironmentOutput struct {
	Environment *Environment
}

// ListEnvironmentsOutput is the output for List operation
type ListEnvironmentsOutput struct {
	Environments []*Environment
}

// CreateEnvironmentInput is the input for Create operation
type CreateEnvironmentInput struct {
	Name  string
	Slug  string
	Color string
}

// CreateEnvironmentOutput is the output for Create operation
type CreateEnvironmentOutput struct {
	Environment *Environment
}

// UpdateEnvironmentInput is the input for Update operation
type UpdateEnvironmentInput struct {
	EnvironmentID string
	Name          *string
	Color         *string
}

// UpdateEnvironmentOutput is the output for Update operation
type UpdateEnvironmentOutput struct {
	Environment *Environment
}

// DeleteEnvironmentInput is the input for Delete operation
type DeleteEnvironmentInput struct {
	EnvironmentID string
}

// DeleteEnvironmentOutput is the output for Delete operation
type DeleteEnvironmentOutput struct {
	Environment *Environment
}
