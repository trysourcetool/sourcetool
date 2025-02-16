package types

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/model"
)

type HostInstancePayload struct {
	ID         string `json:"id"`
	SDKName    string `json:"sdkName"`
	SDKVersion string `json:"sdkVersion"`
	Status     string `json:"status"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

func HostInstanceToPayload(m *model.HostInstance) *HostInstancePayload {
	return &HostInstancePayload{
		ID:         m.ID.String(),
		SDKName:    m.SDKName,
		SDKVersion: m.SDKVersion,
		Status:     m.Status.String(),
		CreatedAt:  strconv.FormatInt(m.CreatedAt.Unix(), 10),
		UpdatedAt:  strconv.FormatInt(m.UpdatedAt.Unix(), 10),
	}
}

type UpdateHostInstanceStatusInput struct {
	ID     string
	Status model.HostInstanceStatus
}

type UpdateHostInstanceStatusPayload struct {
	HostInstance *HostInstancePayload
}
