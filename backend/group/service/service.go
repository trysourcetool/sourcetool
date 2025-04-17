package service

import (
	"context"
	"errors"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type GroupService interface {
	Get(context.Context, dto.GetGroupInput) (*dto.GetGroupOutput, error)
	List(context.Context) (*dto.ListGroupsOutput, error)
	Create(context.Context, dto.CreateGroupInput) (*dto.CreateGroupOutput, error)
	Update(context.Context, dto.UpdateGroupInput) (*dto.UpdateGroupOutput, error)
	Delete(context.Context, dto.DeleteGroupInput) (*dto.DeleteGroupOutput, error)
}

type GroupServiceCE struct {
	*infra.Dependency
}

func NewGroupServiceCE(d *infra.Dependency) *GroupServiceCE {
	return &GroupServiceCE{Dependency: d}
}

func (s *GroupServiceCE) Get(ctx context.Context, in dto.GetGroupInput) (*dto.GetGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *GroupServiceCE) List(ctx context.Context) (*dto.ListGroupsOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *GroupServiceCE) Create(ctx context.Context, in dto.CreateGroupInput) (*dto.CreateGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *GroupServiceCE) Update(ctx context.Context, in dto.UpdateGroupInput) (*dto.UpdateGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *GroupServiceCE) Delete(ctx context.Context, in dto.DeleteGroupInput) (*dto.DeleteGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}
