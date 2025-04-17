package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto/http/requests"
	"github.com/trysourcetool/sourcetool/backend/dto/http/responses"
	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
)

// EnvironmentOutputToResponse converts from output.Environment to responses.EnvironmentResponse.
func EnvironmentOutputToResponse(env *output.Environment) *responses.EnvironmentResponse {
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

// GetEnvironmentRequestToInput converts from requests.GetEnvironmentRequest to input.GetEnvironmentInput.
func GetEnvironmentRequestToInput(in requests.GetEnvironmentRequest) input.GetEnvironmentInput {
	return input.GetEnvironmentInput{
		EnvironmentID: in.EnvironmentID,
	}
}

// GetEnvironmentOutputToResponse converts from output.GetEnvironmentOutput to responses.GetEnvironmentResponse.
func GetEnvironmentOutputToResponse(out *output.GetEnvironmentOutput) *responses.GetEnvironmentResponse {
	return &responses.GetEnvironmentResponse{
		Environment: EnvironmentOutputToResponse(out.Environment),
	}
}

// ListEnvironmentsOutputToResponse converts from output.ListEnvironmentsOutput to responses.ListEnvironmentsResponse.
func ListEnvironmentsOutputToResponse(out *output.ListEnvironmentsOutput) *responses.ListEnvironmentsResponse {
	envs := make([]*responses.EnvironmentResponse, 0, len(out.Environments))
	for _, env := range out.Environments {
		envs = append(envs, EnvironmentOutputToResponse(env))
	}

	return &responses.ListEnvironmentsResponse{
		Environments: envs,
	}
}

// CreateEnvironmentRequestToInput converts from requests.CreateEnvironmentRequest to input.CreateEnvironmentInput.
func CreateEnvironmentRequestToInput(in requests.CreateEnvironmentRequest) input.CreateEnvironmentInput {
	return input.CreateEnvironmentInput{
		Name:  in.Name,
		Slug:  in.Slug,
		Color: in.Color,
	}
}

// CreateEnvironmentOutputToResponse converts from output.CreateEnvironmentOutput to responses.CreateEnvironmentResponse.
func CreateEnvironmentOutputToResponse(out *output.CreateEnvironmentOutput) *responses.CreateEnvironmentResponse {
	return &responses.CreateEnvironmentResponse{
		Environment: EnvironmentOutputToResponse(out.Environment),
	}
}

// UpdateEnvironmentRequestToInput converts from requests.UpdateEnvironmentRequest to input.UpdateEnvironmentInput.
func UpdateEnvironmentRequestToInput(in requests.UpdateEnvironmentRequest) input.UpdateEnvironmentInput {
	return input.UpdateEnvironmentInput{
		EnvironmentID: in.EnvironmentID,
		Name:          in.Name,
		Color:         in.Color,
	}
}

// UpdateEnvironmentOutputToResponse converts from output.UpdateEnvironmentOutput to responses.UpdateEnvironmentResponse.
func UpdateEnvironmentOutputToResponse(out *output.UpdateEnvironmentOutput) *responses.UpdateEnvironmentResponse {
	return &responses.UpdateEnvironmentResponse{
		Environment: EnvironmentOutputToResponse(out.Environment),
	}
}

// DeleteEnvironmentRequestToInput converts from requests.DeleteEnvironmentRequest to input.DeleteEnvironmentInput.
func DeleteEnvironmentRequestToInput(in requests.DeleteEnvironmentRequest) input.DeleteEnvironmentInput {
	return input.DeleteEnvironmentInput{
		EnvironmentID: in.EnvironmentID,
	}
}

// DeleteEnvironmentOutputToResponse converts from output.DeleteEnvironmentOutput to responses.DeleteEnvironmentResponse.
func DeleteEnvironmentOutputToResponse(out *output.DeleteEnvironmentOutput) *responses.DeleteEnvironmentResponse {
	return &responses.DeleteEnvironmentResponse{
		Environment: EnvironmentOutputToResponse(out.Environment),
	}
}
