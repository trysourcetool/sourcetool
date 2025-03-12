package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
)

// EnvironmentDTOToResponse converts from dto.Environment to responses.EnvironmentResponse.
func EnvironmentDTOToResponse(env *dto.Environment) *responses.EnvironmentResponse {
	if env == nil {
		return nil
	}

	return &responses.EnvironmentResponse{
		ID:        env.ID,
		Name:      env.Name,
		Slug:      env.Slug,
		Color:     env.Color,
		CreatedAt: strconv.FormatInt(env.CreatedAt, 10),
		UpdatedAt: strconv.FormatInt(env.UpdatedAt, 10),
	}
}

// GetEnvironmentRequestToDTOInput converts from requests.GetEnvironmentRequest to dto.GetEnvironmentInput.
func GetEnvironmentRequestToDTOInput(in requests.GetEnvironmentRequest) dto.GetEnvironmentInput {
	return dto.GetEnvironmentInput{
		EnvironmentID: in.EnvironmentID,
	}
}

// GetEnvironmentOutputToResponse converts from dto.GetEnvironmentOutput to responses.GetEnvironmentResponse.
func GetEnvironmentOutputToResponse(out *dto.GetEnvironmentOutput) *responses.GetEnvironmentResponse {
	return &responses.GetEnvironmentResponse{
		Environment: EnvironmentDTOToResponse(out.Environment),
	}
}

// ListEnvironmentsOutputToResponse converts from dto.ListEnvironmentsOutput to responses.ListEnvironmentsResponse.
func ListEnvironmentsOutputToResponse(out *dto.ListEnvironmentsOutput) *responses.ListEnvironmentsResponse {
	envs := make([]*responses.EnvironmentResponse, 0, len(out.Environments))
	for _, env := range out.Environments {
		envs = append(envs, EnvironmentDTOToResponse(env))
	}

	return &responses.ListEnvironmentsResponse{
		Environments: envs,
	}
}

// CreateEnvironmentRequestToDTOInput converts from requests.CreateEnvironmentRequest to dto.CreateEnvironmentInput.
func CreateEnvironmentRequestToDTOInput(in requests.CreateEnvironmentRequest) dto.CreateEnvironmentInput {
	return dto.CreateEnvironmentInput{
		Name:  in.Name,
		Slug:  in.Slug,
		Color: in.Color,
	}
}

// CreateEnvironmentOutputToResponse converts from dto.CreateEnvironmentOutput to responses.CreateEnvironmentResponse.
func CreateEnvironmentOutputToResponse(out *dto.CreateEnvironmentOutput) *responses.CreateEnvironmentResponse {
	return &responses.CreateEnvironmentResponse{
		Environment: EnvironmentDTOToResponse(out.Environment),
	}
}

// UpdateEnvironmentRequestToDTOInput converts from requests.UpdateEnvironmentRequest to dto.UpdateEnvironmentInput.
func UpdateEnvironmentRequestToDTOInput(in requests.UpdateEnvironmentRequest) dto.UpdateEnvironmentInput {
	return dto.UpdateEnvironmentInput{
		EnvironmentID: in.EnvironmentID,
		Name:          in.Name,
		Color:         in.Color,
	}
}

// UpdateEnvironmentOutputToResponse converts from dto.UpdateEnvironmentOutput to responses.UpdateEnvironmentResponse.
func UpdateEnvironmentOutputToResponse(out *dto.UpdateEnvironmentOutput) *responses.UpdateEnvironmentResponse {
	return &responses.UpdateEnvironmentResponse{
		Environment: EnvironmentDTOToResponse(out.Environment),
	}
}

// DeleteEnvironmentRequestToDTOInput converts from requests.DeleteEnvironmentRequest to dto.DeleteEnvironmentInput.
func DeleteEnvironmentRequestToDTOInput(in requests.DeleteEnvironmentRequest) dto.DeleteEnvironmentInput {
	return dto.DeleteEnvironmentInput{
		EnvironmentID: in.EnvironmentID,
	}
}

// DeleteEnvironmentOutputToResponse converts from dto.DeleteEnvironmentOutput to responses.DeleteEnvironmentResponse.
func DeleteEnvironmentOutputToResponse(out *dto.DeleteEnvironmentOutput) *responses.DeleteEnvironmentResponse {
	return &responses.DeleteEnvironmentResponse{
		Environment: EnvironmentDTOToResponse(out.Environment),
	}
}
