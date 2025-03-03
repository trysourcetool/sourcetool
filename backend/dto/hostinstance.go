package dto

import (
	"github.com/gofrs/uuid/v5"
	"github.com/trysourcetool/sourcetool/backend/model"
)

// HostInstance represents host instance data in DTOs
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

// HostInstanceFromModel converts from model.HostInstance to dto.HostInstance
func HostInstanceFromModel(instance *model.HostInstance) *HostInstance {
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

// ToHostInstanceID converts string ID to uuid.UUID
func ToHostInstanceID(id string) (uuid.UUID, error) {
	return uuid.FromString(id)
}

// PingHostInstanceInput is the input for Ping operation
type PingHostInstanceInput struct {
	PageID *string
}

// PingHostInstanceOutput is the output for Ping operation
type PingHostInstanceOutput struct {
	HostInstance *HostInstance
}
