package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
)

// APIKeyDTOToResponse converts from dto.APIKey to responses.APIKeyResponse
func APIKeyDTOToResponse(apiKey *dto.APIKey) *responses.APIKeyResponse {
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
		result.Environment = EnvironmentDTOToResponse(apiKey.Environment)
	}

	return result
}

// GetAPIKeyRequestToDTOInput converts from requests.GetAPIKeyRequest to dto.GetAPIKeyInput
func GetAPIKeyRequestToDTOInput(in requests.GetAPIKeyRequest) dto.GetAPIKeyInput {
	return dto.GetAPIKeyInput{
		APIKeyID: in.APIKeyID,
	}
}

// GetAPIKeyOutputToResponse converts from dto.GetAPIKeyOutput to responses.GetAPIKeyResponse
func GetAPIKeyOutputToResponse(out *dto.GetAPIKeyOutput) *responses.GetAPIKeyResponse {
	return &responses.GetAPIKeyResponse{
		APIKey: APIKeyDTOToResponse(out.APIKey),
	}
}

// ListAPIKeysOutputToResponse converts from dto.ListAPIKeysOutput to responses.ListAPIKeysResponse
func ListAPIKeysOutputToResponse(out *dto.ListAPIKeysOutput) *responses.ListAPIKeysResponse {
	liveKeys := make([]*responses.APIKeyResponse, 0, len(out.LiveKeys))
	for _, key := range out.LiveKeys {
		liveKeys = append(liveKeys, APIKeyDTOToResponse(key))
	}

	return &responses.ListAPIKeysResponse{
		DevKey:   APIKeyDTOToResponse(out.DevKey),
		LiveKeys: liveKeys,
	}
}

// CreateAPIKeyRequestToDTOInput converts from requests.CreateAPIKeyRequest to dto.CreateAPIKeyInput
func CreateAPIKeyRequestToDTOInput(in requests.CreateAPIKeyRequest) dto.CreateAPIKeyInput {
	return dto.CreateAPIKeyInput{
		EnvironmentID: in.EnvironmentID,
		Name:          in.Name,
	}
}

// CreateAPIKeyOutputToResponse converts from dto.CreateAPIKeyOutput to responses.CreateAPIKeyResponse
func CreateAPIKeyOutputToResponse(out *dto.CreateAPIKeyOutput) *responses.CreateAPIKeyResponse {
	return &responses.CreateAPIKeyResponse{
		APIKey: APIKeyDTOToResponse(out.APIKey),
	}
}

// UpdateAPIKeyRequestToDTOInput converts from requests.UpdateAPIKeyRequest to dto.UpdateAPIKeyInput
func UpdateAPIKeyRequestToDTOInput(in requests.UpdateAPIKeyRequest) dto.UpdateAPIKeyInput {
	return dto.UpdateAPIKeyInput{
		APIKeyID: in.APIKeyID,
		Name:     in.Name,
	}
}

// UpdateAPIKeyOutputToResponse converts from dto.UpdateAPIKeyOutput to responses.UpdateAPIKeyResponse
func UpdateAPIKeyOutputToResponse(out *dto.UpdateAPIKeyOutput) *responses.UpdateAPIKeyResponse {
	return &responses.UpdateAPIKeyResponse{
		APIKey: APIKeyDTOToResponse(out.APIKey),
	}
}

// DeleteAPIKeyRequestToDTOInput converts from requests.DeleteAPIKeyRequest to dto.DeleteAPIKeyInput
func DeleteAPIKeyRequestToDTOInput(in requests.DeleteAPIKeyRequest) dto.DeleteAPIKeyInput {
	return dto.DeleteAPIKeyInput{
		APIKeyID: in.APIKeyID,
	}
}

// DeleteAPIKeyOutputToResponse converts from dto.DeleteAPIKeyOutput to responses.DeleteAPIKeyResponse
func DeleteAPIKeyOutputToResponse(out *dto.DeleteAPIKeyOutput) *responses.DeleteAPIKeyResponse {
	return &responses.DeleteAPIKeyResponse{
		APIKey: APIKeyDTOToResponse(out.APIKey),
	}
}
