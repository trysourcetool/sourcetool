package dto

import "github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"

type PingHostInstanceInput struct {
	PageID *string
}

type HostInstance struct {
	ID             string
	OrganizationID string
	APIKeyID       string
	SDKName        string
	SDKVersion     string
	Status         string
	CreatedAt      int64
	UpdatedAt      int64
}

func HostInstanceFromModel(instance *hostinstance.HostInstance) *HostInstance {
	if instance == nil {
		return nil
	}

	return &HostInstance{
		ID:             instance.ID.String(),
		OrganizationID: instance.OrganizationID.String(),
		APIKeyID:       instance.APIKeyID.String(),
		SDKName:        instance.SDKName,
		SDKVersion:     instance.SDKVersion,
		Status:         instance.Status.String(),
		CreatedAt:      instance.CreatedAt.Unix(),
		UpdatedAt:      instance.UpdatedAt.Unix(),
	}
}

type PingHostInstanceOutput struct {
	HostInstance *HostInstance
}
