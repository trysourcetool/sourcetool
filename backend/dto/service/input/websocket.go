package input

import "github.com/trysourcetool/sourcetool/backend/hostinstance"

// UpdateHostInstanceStatusInput is a struct that represents the input for updating a host instance status.
type UpdateHostInstanceStatusInput struct {
	ID     string
	Status hostinstance.HostInstanceStatus
}
