package responses

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type HostInstanceResponse struct {
	ID         string `json:"id"`
	SDKName    string `json:"sdkName"`
	SDKVersion string `json:"sdkVersion"`
	Status     string `json:"status"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

func HostInstanceFromModel(hostInstance *core.HostInstance) *HostInstanceResponse {
	if hostInstance == nil {
		return nil
	}

	return &HostInstanceResponse{
		ID:         hostInstance.ID.String(),
		SDKName:    hostInstance.SDKName,
		SDKVersion: hostInstance.SDKVersion,
		Status:     hostInstance.Status.String(),
		CreatedAt:  strconv.FormatInt(hostInstance.CreatedAt.Unix(), 10),
		UpdatedAt:  strconv.FormatInt(hostInstance.UpdatedAt.Unix(), 10),
	}
}

type PingHostInstanceResponse struct {
	HostInstance *HostInstanceResponse `json:"hostInstance"`
}
