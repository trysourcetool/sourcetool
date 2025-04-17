package dto

import "github.com/trysourcetool/sourcetool/backend/hostinstance"

// UpdateHostInstanceStatusInput is a struct that represents the input for updating a host instance status.
type UpdateHostInstanceStatusInput struct {
	ID     string
	Status hostinstance.HostInstanceStatus
}

// UpdateHostInstanceStatusOutput is a struct that represents the output for updating a host instance status.
type UpdateHostInstanceStatusOutput struct {
	HostInstance *HostInstance
}
