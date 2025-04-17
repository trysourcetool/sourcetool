package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto/http/requests"
	"github.com/trysourcetool/sourcetool/backend/dto/http/responses"
	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
)

// APIKeyOutputToResponse converts from output.APIKey to responses.APIKeyResponse.
func APIKeyOutputToResponse(apiKey *output.APIKey) *responses.APIKeyResponse {
	if apiKey == nil {
		return nil
	}

	result := &responses.APIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		Key:       apiKey.Key,
		CreatedAt: strconv.FormatInt(apiKey.CreatedAt, 10),
		UpdatedAt: strconv.FormatInt(apiKey.UpdatedAt, 10),
	}

	if apiKey.Environment != nil {
		result.Environment = EnvironmentOutputToResponse(apiKey.Environment)
	}

	return result
}

// GetAPIKeyRequestToInput converts from requests.GetAPIKeyRequest to input.GetAPIKeyInput.
func GetAPIKeyRequestToInput(in requests.GetAPIKeyRequest) input.GetAPIKeyInput {
	return input.GetAPIKeyInput{
		APIKeyID: in.APIKeyID,
	}
}

// GetAPIKeyOutputToResponse converts from output.GetAPIKeyOutput to responses.GetAPIKeyResponse.
func GetAPIKeyOutputToResponse(out *output.GetAPIKeyOutput) *responses.GetAPIKeyResponse {
	return &responses.GetAPIKeyResponse{
		APIKey: APIKeyOutputToResponse(out.APIKey),
	}
}

// ListAPIKeysOutputToResponse converts from output.ListAPIKeysOutput to responses.ListAPIKeysResponse.
func ListAPIKeysOutputToResponse(out *output.ListAPIKeysOutput) *responses.ListAPIKeysResponse {
	liveKeys := make([]*responses.APIKeyResponse, 0, len(out.LiveKeys))
	for _, key := range out.LiveKeys {
		liveKeys = append(liveKeys, APIKeyOutputToResponse(key))
	}

	return &responses.ListAPIKeysResponse{
		DevKey:   APIKeyOutputToResponse(out.DevKey),
		LiveKeys: liveKeys,
	}
}

// CreateAPIKeyRequestToInput converts from requests.CreateAPIKeyRequest to input.CreateAPIKeyInput.
func CreateAPIKeyRequestToInput(in requests.CreateAPIKeyRequest) input.CreateAPIKeyInput {
	return input.CreateAPIKeyInput{
		EnvironmentID: in.EnvironmentID,
		Name:          in.Name,
	}
}

// CreateAPIKeyOutputToResponse converts from output.CreateAPIKeyOutput to responses.CreateAPIKeyResponse.
func CreateAPIKeyOutputToResponse(out *output.CreateAPIKeyOutput) *responses.CreateAPIKeyResponse {
	return &responses.CreateAPIKeyResponse{
		APIKey: APIKeyOutputToResponse(out.APIKey),
	}
}

// UpdateAPIKeyRequestToInput converts from requests.UpdateAPIKeyRequest to input.UpdateAPIKeyInput.
func UpdateAPIKeyRequestToInput(in requests.UpdateAPIKeyRequest) input.UpdateAPIKeyInput {
	return input.UpdateAPIKeyInput{
		APIKeyID: in.APIKeyID,
		Name:     in.Name,
	}
}

// UpdateAPIKeyOutputToResponse converts from output.UpdateAPIKeyOutput to responses.UpdateAPIKeyResponse.
func UpdateAPIKeyOutputToResponse(out *output.UpdateAPIKeyOutput) *responses.UpdateAPIKeyResponse {
	return &responses.UpdateAPIKeyResponse{
		APIKey: APIKeyOutputToResponse(out.APIKey),
	}
}

// DeleteAPIKeyRequestToInput converts from requests.DeleteAPIKeyRequest to input.DeleteAPIKeyInput.
func DeleteAPIKeyRequestToInput(in requests.DeleteAPIKeyRequest) input.DeleteAPIKeyInput {
	return input.DeleteAPIKeyInput{
		APIKeyID: in.APIKeyID,
	}
}

// DeleteAPIKeyOutputToResponse converts from output.DeleteAPIKeyOutput to responses.DeleteAPIKeyResponse.
func DeleteAPIKeyOutputToResponse(out *output.DeleteAPIKeyOutput) *responses.DeleteAPIKeyResponse {
	return &responses.DeleteAPIKeyResponse{
		APIKey: APIKeyOutputToResponse(out.APIKey),
	}
}
