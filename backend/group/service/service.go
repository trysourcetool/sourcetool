package service

import (
	"context"
	"errors"

	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type GroupService interface {
	Get(context.Context, input.GetGroupInput) (*output.GetGroupOutput, error)
	List(context.Context) (*output.ListGroupsOutput, error)
	Create(context.Context, input.CreateGroupInput) (*output.CreateGroupOutput, error)
	Update(context.Context, input.UpdateGroupInput) (*output.UpdateGroupOutput, error)
	Delete(context.Context, input.DeleteGroupInput) (*output.DeleteGroupOutput, error)
}

type GroupServiceCE struct {
	*infra.Dependency
}

func NewGroupServiceCE(d *infra.Dependency) *GroupServiceCE {
	return &GroupServiceCE{Dependency: d}
}

func (s *GroupServiceCE) Get(ctx context.Context, in input.GetGroupInput) (*output.GetGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *GroupServiceCE) List(ctx context.Context) (*output.ListGroupsOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *GroupServiceCE) Create(ctx context.Context, in input.CreateGroupInput) (*output.CreateGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *GroupServiceCE) Update(ctx context.Context, in input.UpdateGroupInput) (*output.UpdateGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *GroupServiceCE) Delete(ctx context.Context, in input.DeleteGroupInput) (*output.DeleteGroupOutput, error) {
	return nil, errors.New("group functionality is not available in CE version")
}
