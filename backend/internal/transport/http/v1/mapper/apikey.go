package mapper

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/responses"
)

func APIKeyOutputToResponse(apiKey *dto.APIKey) *responses.APIKeyResponse {
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

func GetAPIKeyRequestToInput(in requests.GetAPIKeyRequest) dto.GetAPIKeyInput {
	return dto.GetAPIKeyInput{
		APIKeyID: in.APIKeyID,
	}
}

func GetAPIKeyOutputToResponse(out *dto.GetAPIKeyOutput) *responses.GetAPIKeyResponse {
	return &responses.GetAPIKeyResponse{
		APIKey: APIKeyOutputToResponse(out.APIKey),
	}
}

func ListAPIKeysOutputToResponse(out *dto.ListAPIKeysOutput) *responses.ListAPIKeysResponse {
	liveKeys := make([]*responses.APIKeyResponse, 0, len(out.LiveKeys))
	for _, key := range out.LiveKeys {
		liveKeys = append(liveKeys, APIKeyOutputToResponse(key))
	}

	return &responses.ListAPIKeysResponse{
		DevKey:   APIKeyOutputToResponse(out.DevKey),
		LiveKeys: liveKeys,
	}
}

func CreateAPIKeyRequestToInput(in requests.CreateAPIKeyRequest) dto.CreateAPIKeyInput {
	return dto.CreateAPIKeyInput{
		EnvironmentID: in.EnvironmentID,
		Name:          in.Name,
	}
}

func CreateAPIKeyOutputToResponse(out *dto.CreateAPIKeyOutput) *responses.CreateAPIKeyResponse {
	return &responses.CreateAPIKeyResponse{
		APIKey: APIKeyOutputToResponse(out.APIKey),
	}
}

func UpdateAPIKeyRequestToInput(in requests.UpdateAPIKeyRequest) dto.UpdateAPIKeyInput {
	return dto.UpdateAPIKeyInput{
		APIKeyID: in.APIKeyID,
		Name:     in.Name,
	}
}

func UpdateAPIKeyOutputToResponse(out *dto.UpdateAPIKeyOutput) *responses.UpdateAPIKeyResponse {
	return &responses.UpdateAPIKeyResponse{
		APIKey: APIKeyOutputToResponse(out.APIKey),
	}
}

func DeleteAPIKeyRequestToInput(in requests.DeleteAPIKeyRequest) dto.DeleteAPIKeyInput {
	return dto.DeleteAPIKeyInput{
		APIKeyID: in.APIKeyID,
	}
}

func DeleteAPIKeyOutputToResponse(out *dto.DeleteAPIKeyOutput) *responses.DeleteAPIKeyResponse {
	return &responses.DeleteAPIKeyResponse{
		APIKey: APIKeyOutputToResponse(out.APIKey),
	}
}
