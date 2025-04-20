package group

import (
	"context"
	"errors"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
)

type Service interface {
	Get(context.Context, dto.GetGroupInput) (*dto.GetGroupOutput, error)
	List(context.Context) (*dto.ListGroupsOutput, error)
	Create(context.Context, dto.CreateGroupInput) (*dto.CreateGroupOutput, error)
	Update(context.Context, dto.UpdateGroupInput) (*dto.UpdateGroupOutput, error)
	Delete(context.Context, dto.DeleteGroupInput) (*dto.DeleteGroupOutput, error)
}

type ServiceCE struct {
	*port.Dependencies
}

func NewServiceCE(d *port.Dependencies) *ServiceCE {
	return &ServiceCE{Dependencies: d}
}

func (s *ServiceCE) Get(ctx context.Context, in dto.GetGroupInput) (*dto.GetGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *ServiceCE) List(ctx context.Context) (*dto.ListGroupsOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *ServiceCE) Create(ctx context.Context, in dto.CreateGroupInput) (*dto.CreateGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *ServiceCE) Update(ctx context.Context, in dto.UpdateGroupInput) (*dto.UpdateGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *ServiceCE) Delete(ctx context.Context, in dto.DeleteGroupInput) (*dto.DeleteGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}
