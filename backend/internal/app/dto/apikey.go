package dto

import (
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/environment"
)

type GetAPIKeyInput struct {
	APIKeyID string
}

type CreateAPIKeyInput struct {
	EnvironmentID string
	Name          string
}

type UpdateAPIKeyInput struct {
	APIKeyID string
	Name     *string
}

type DeleteAPIKeyInput struct {
	APIKeyID string
}

type APIKey struct {
	ID          string
	Name        string
	Key         string
	CreatedAt   int64
	UpdatedAt   int64
	Environment *Environment
}

func APIKeyFromModel(apiKey *apikey.APIKey, env *environment.Environment) *APIKey {
	if apiKey == nil {
		return nil
	}

	result := &APIKey{
		ID:        apiKey.ID.String(),
		Name:      apiKey.Name,
		Key:       apiKey.Key,
		CreatedAt: apiKey.CreatedAt.Unix(),
		UpdatedAt: apiKey.UpdatedAt.Unix(),
	}

	if env != nil {
		result.Environment = EnvironmentFromModel(env)
	}

	return result
}

type GetAPIKeyOutput struct {
	APIKey *APIKey
}

type ListAPIKeysOutput struct {
	DevKey   *APIKey
	LiveKeys []*APIKey
}

type CreateAPIKeyOutput struct {
	APIKey *APIKey
}

type UpdateAPIKeyOutput struct {
	APIKey *APIKey
}

type DeleteAPIKeyOutput struct {
	APIKey *APIKey
}
