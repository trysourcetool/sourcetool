package output

import (
	"github.com/trysourcetool/sourcetool/backend/apikey"
	"github.com/trysourcetool/sourcetool/backend/environment"
)

// APIKey represents API key data in DTOs.
type APIKey struct {
	ID          string
	Name        string
	Key         string
	CreatedAt   int64
	UpdatedAt   int64
	Environment *Environment
}

// APIKeyFromModel converts from model.APIKey to dto.APIKey.
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

// GetAPIKeyOutput is the output for Get operation.
type GetAPIKeyOutput struct {
	APIKey *APIKey
}

// ListAPIKeysOutput is the output for List operation.
type ListAPIKeysOutput struct {
	DevKey   *APIKey
	LiveKeys []*APIKey
}

// CreateAPIKeyOutput is the output for Create operation.
type CreateAPIKeyOutput struct {
	APIKey *APIKey
}

// UpdateAPIKeyOutput is the output for Update operation.
type UpdateAPIKeyOutput struct {
	APIKey *APIKey
}

// DeleteAPIKeyOutput is the output for Delete operation.
type DeleteAPIKeyOutput struct {
	APIKey *APIKey
}
