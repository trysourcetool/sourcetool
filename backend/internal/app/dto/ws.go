package dto

import "github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"

type UpdateHostInstanceStatusInput struct {
	ID     string
	Status hostinstance.HostInstanceStatus
}

type UpdateHostInstanceStatusOutput struct {
	HostInstance *HostInstance
}
