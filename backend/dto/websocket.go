package dto

import "github.com/trysourcetool/sourcetool/backend/model"

// UpdateHostInstanceStatusInput is a struct that represents the input for updating a host instance status.
type UpdateHostInstanceStatusInput struct {
	ID     string
	Status model.HostInstanceStatus
}

// UpdateHostInstanceStatusOutput is a struct that represents the output for updating a host instance status.
type UpdateHostInstanceStatusOutput struct {
	HostInstance *HostInstance
}
