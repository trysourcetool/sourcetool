package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
)

// HostInstanceDTOToResponse converts from dto.HostInstance to responses.HostInstanceResponse
func HostInstanceDTOToResponse(instance *dto.HostInstance) *responses.HostInstanceResponse {
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

// PingHostInstanceRequestToDTOInput converts from requests.PingHostInstanceRequest to dto.PingHostInstanceInput
func PingHostInstanceRequestToDTOInput(in requests.PingHostInstanceRequest) dto.PingHostInstanceInput {
	return dto.PingHostInstanceInput{
		PageID: in.PageID,
	}
}

// PingHostInstanceOutputToResponse converts from dto.PingHostInstanceOutput to responses.PingHostInstanceResponse
func PingHostInstanceOutputToResponse(out *dto.PingHostInstanceOutput) *responses.PingHostInstanceResponse {
	return &responses.PingHostInstanceResponse{
		HostInstance: HostInstanceDTOToResponse(out.HostInstance),
	}
}
