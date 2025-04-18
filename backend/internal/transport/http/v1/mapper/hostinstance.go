package mapper

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/responses"
)

func HostInstanceOutputToResponse(instance *dto.HostInstance) *responses.HostInstanceResponse {
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

func PingHostInstanceRequestToInput(in requests.PingHostInstanceRequest) dto.PingHostInstanceInput {
	return dto.PingHostInstanceInput{
		PageID: in.PageID,
	}
}

func PingHostInstanceOutputToResponse(out *dto.PingHostInstanceOutput) *responses.PingHostInstanceResponse {
	return &responses.PingHostInstanceResponse{
		HostInstance: HostInstanceOutputToResponse(out.HostInstance),
	}
}
