package dto

import (
	"github.com/trysourcetool/sourcetool/backend/model"
)

// APIKey represents API key data in DTOs
type APIKey struct {
	ID          string
	Name        string
	Key         string
	CreatedAt   int64
	UpdatedAt   int64
	Environment *Environment
}

// APIKeyFromModel converts from model.APIKey to dto.APIKey
func APIKeyFromModel(apiKey *model.APIKey, env *model.Environment) *APIKey {
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

// GetAPIKeyInput is the input for Get operation
type GetAPIKeyInput struct {
	APIKeyID string
}

// GetAPIKeyOutput is the output for Get operation
type GetAPIKeyOutput struct {
	APIKey *APIKey
}

// ListAPIKeysOutput is the output for List operation
type ListAPIKeysOutput struct {
	DevKey   *APIKey
	LiveKeys []*APIKey
}

// CreateAPIKeyInput is the input for Create operation
type CreateAPIKeyInput struct {
	EnvironmentID string
	Name          string
}

// CreateAPIKeyOutput is the output for Create operation
type CreateAPIKeyOutput struct {
	APIKey *APIKey
}

// UpdateAPIKeyInput is the input for Update operation
type UpdateAPIKeyInput struct {
	APIKeyID string
	Name     *string
}

// UpdateAPIKeyOutput is the output for Update operation
type UpdateAPIKeyOutput struct {
	APIKey *APIKey
}

// DeleteAPIKeyInput is the input for Delete operation
type DeleteAPIKeyInput struct {
	APIKeyID string
}

// DeleteAPIKeyOutput is the output for Delete operation
type DeleteAPIKeyOutput struct {
	APIKey *APIKey
}
