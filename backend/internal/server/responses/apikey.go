package responses

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type APIKeyResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Key         string               `json:"key"`
	CreatedAt   string               `json:"createdAt"`
	UpdatedAt   string               `json:"updatedAt"`
	Environment *EnvironmentResponse `json:"environment,omitempty"`
}

func APIKeyFromModel(apiKey *core.APIKey, env *core.Environment) *APIKeyResponse {
	if apiKey == nil {
		return nil
	}

	return &APIKeyResponse{
		ID:          apiKey.ID.String(),
		Name:        apiKey.Name,
		Key:         apiKey.Key,
		CreatedAt:   strconv.FormatInt(apiKey.CreatedAt.Unix(), 10),
		UpdatedAt:   strconv.FormatInt(apiKey.UpdatedAt.Unix(), 10),
		Environment: EnvironmentFromModel(env),
	}
}

type GetAPIKeyResponse struct {
	APIKey *APIKeyResponse `json:"apiKey"`
}

type ListAPIKeysResponse struct {
	DevKey   *APIKeyResponse   `json:"devKey"`
	LiveKeys []*APIKeyResponse `json:"liveKeys"`
}

type CreateAPIKeyResponse struct {
	APIKey *APIKeyResponse `json:"apiKey"`
}

type UpdateAPIKeyResponse struct {
	APIKey *APIKeyResponse `json:"apiKey"`
}

type DeleteAPIKeyResponse struct {
	APIKey *APIKeyResponse `json:"apiKey"`
}
