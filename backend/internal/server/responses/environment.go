package responses

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type EnvironmentResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Color     string `json:"color"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func EnvironmentFromModel(env *core.Environment) *EnvironmentResponse {
	if env == nil {
		return nil
	}

	return &EnvironmentResponse{
		ID:        env.ID.String(),
		Name:      env.Name,
		Slug:      env.Slug,
		Color:     env.Color,
		CreatedAt: strconv.FormatInt(env.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(env.UpdatedAt.Unix(), 10),
	}
}

type ListEnvironmentsResponse struct {
	Environments []*EnvironmentResponse `json:"environments"`
}

type GetEnvironmentResponse struct {
	Environment *EnvironmentResponse `json:"environment"`
}

type CreateEnvironmentResponse struct {
	Environment *EnvironmentResponse `json:"environment"`
}

type UpdateEnvironmentResponse struct {
	Environment *EnvironmentResponse `json:"environment"`
}

type DeleteEnvironmentResponse struct {
	Environment *EnvironmentResponse `json:"environment"`
}
