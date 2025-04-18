package mapper

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/responses"
)

func EnvironmentOutputToResponse(env *dto.Environment) *responses.EnvironmentResponse {
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

func GetEnvironmentRequestToInput(in requests.GetEnvironmentRequest) dto.GetEnvironmentInput {
	return dto.GetEnvironmentInput{
		EnvironmentID: in.EnvironmentID,
	}
}

func GetEnvironmentOutputToResponse(out *dto.GetEnvironmentOutput) *responses.GetEnvironmentResponse {
	return &responses.GetEnvironmentResponse{
		Environment: EnvironmentOutputToResponse(out.Environment),
	}
}

func ListEnvironmentsOutputToResponse(out *dto.ListEnvironmentsOutput) *responses.ListEnvironmentsResponse {
	envs := make([]*responses.EnvironmentResponse, 0, len(out.Environments))
	for _, env := range out.Environments {
		envs = append(envs, EnvironmentOutputToResponse(env))
	}

	return &responses.ListEnvironmentsResponse{
		Environments: envs,
	}
}

func CreateEnvironmentRequestToInput(in requests.CreateEnvironmentRequest) dto.CreateEnvironmentInput {
	return dto.CreateEnvironmentInput{
		Name:  in.Name,
		Slug:  in.Slug,
		Color: in.Color,
	}
}

func CreateEnvironmentOutputToResponse(out *dto.CreateEnvironmentOutput) *responses.CreateEnvironmentResponse {
	return &responses.CreateEnvironmentResponse{
		Environment: EnvironmentOutputToResponse(out.Environment),
	}
}

func UpdateEnvironmentRequestToInput(in requests.UpdateEnvironmentRequest) dto.UpdateEnvironmentInput {
	return dto.UpdateEnvironmentInput{
		EnvironmentID: in.EnvironmentID,
		Name:          in.Name,
		Color:         in.Color,
	}
}

func UpdateEnvironmentOutputToResponse(out *dto.UpdateEnvironmentOutput) *responses.UpdateEnvironmentResponse {
	return &responses.UpdateEnvironmentResponse{
		Environment: EnvironmentOutputToResponse(out.Environment),
	}
}

func DeleteEnvironmentRequestToInput(in requests.DeleteEnvironmentRequest) dto.DeleteEnvironmentInput {
	return dto.DeleteEnvironmentInput{
		EnvironmentID: in.EnvironmentID,
	}
}

func DeleteEnvironmentOutputToResponse(out *dto.DeleteEnvironmentOutput) *responses.DeleteEnvironmentResponse {
	return &responses.DeleteEnvironmentResponse{
		Environment: EnvironmentOutputToResponse(out.Environment),
	}
}
