package dto

import (
	"github.com/trysourcetool/sourcetool/backend/internal/domain/environment"
)

type GetEnvironmentInput struct {
	EnvironmentID string
}

type CreateEnvironmentInput struct {
	Name  string
	Slug  string
	Color string
}

type UpdateEnvironmentInput struct {
	EnvironmentID string
	Name          *string
	Color         *string
}

type DeleteEnvironmentInput struct {
	EnvironmentID string
}

type Environment struct {
	ID        string
	Name      string
	Slug      string
	Color     string
	CreatedAt int64
	UpdatedAt int64
}

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

type GetEnvironmentOutput struct {
	Environment *Environment
}

type ListEnvironmentsOutput struct {
	Environments []*Environment
}

type CreateEnvironmentOutput struct {
	Environment *Environment
}

type UpdateEnvironmentOutput struct {
	Environment *Environment
}

type DeleteEnvironmentOutput struct {
	Environment *Environment
}
