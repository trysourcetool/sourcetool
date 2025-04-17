package output

import (
	"github.com/trysourcetool/sourcetool/backend/environment"
)

// Environment represents environment data in DTOs.
type Environment struct {
	ID        string
	Name      string
	Slug      string
	Color     string
	CreatedAt int64
	UpdatedAt int64
}

// EnvironmentFromModel converts from model.Environment to dto.Environment.
func EnvironmentFromModel(env *environment.Environment) *Environment {
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

// GetEnvironmentOutput is the output for Get operation.
type GetEnvironmentOutput struct {
	Environment *Environment
}

// ListEnvironmentsOutput is the output for List operation.
type ListEnvironmentsOutput struct {
	Environments []*Environment
}

// CreateEnvironmentOutput is the output for Create operation.
type CreateEnvironmentOutput struct {
	Environment *Environment
}

// UpdateEnvironmentOutput is the output for Update operation.
type UpdateEnvironmentOutput struct {
	Environment *Environment
}

// DeleteEnvironmentOutput is the output for Delete operation.
type DeleteEnvironmentOutput struct {
	Environment *Environment
}
