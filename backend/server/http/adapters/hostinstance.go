package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto/http/requests"
	"github.com/trysourcetool/sourcetool/backend/dto/http/responses"
	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
)

// HostInstanceOutputToResponse converts from output.HostInstance to responses.HostInstanceResponse.
func HostInstanceOutputToResponse(instance *output.HostInstance) *responses.HostInstanceResponse {
	if instance == nil {
		return nil
	}

	return &responses.HostInstanceResponse{
		ID:         instance.ID,
		SDKName:    instance.SDKName,
		SDKVersion: instance.SDKVersion,
		Status:     instance.Status,
		CreatedAt:  strconv.FormatInt(instance.CreatedAt, 10),
		UpdatedAt:  strconv.FormatInt(instance.UpdatedAt, 10),
	}
}

// PingHostInstanceRequestToInput converts from requests.PingHostInstanceRequest to input.PingHostInstanceInput.
func PingHostInstanceRequestToInput(in requests.PingHostInstanceRequest) input.PingHostInstanceInput {
	return input.PingHostInstanceInput{
		PageID: in.PageID,
	}
}

// PingHostInstanceOutputToResponse converts from output.PingHostInstanceOutput to responses.PingHostInstanceResponse.
func PingHostInstanceOutputToResponse(out *output.PingHostInstanceOutput) *responses.PingHostInstanceResponse {
	return &responses.PingHostInstanceResponse{
		HostInstance: HostInstanceOutputToResponse(out.HostInstance),
	}
}
