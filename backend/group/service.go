package group

import (
	"context"
	"errors"

	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
)

type Service interface {
	Get(context.Context, types.GetGroupInput) (*types.GetGroupPayload, error)
	List(context.Context) (*types.ListGroupsPayload, error)
	Create(context.Context, types.CreateGroupInput) (*types.CreateGroupPayload, error)
	Update(context.Context, types.UpdateGroupInput) (*types.UpdateGroupPayload, error)
	Delete(context.Context, types.DeleteGroupInput) (*types.DeleteGroupPayload, error)
}

type ServiceCE struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) Get(ctx context.Context, in types.GetGroupInput) (*types.GetGroupPayload, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *ServiceCE) List(ctx context.Context) (*types.ListGroupsPayload, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *ServiceCE) Create(ctx context.Context, in types.CreateGroupInput) (*types.CreateGroupPayload, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *ServiceCE) Update(ctx context.Context, in types.UpdateGroupInput) (*types.UpdateGroupPayload, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *ServiceCE) Delete(ctx context.Context, in types.DeleteGroupInput) (*types.DeleteGroupPayload, error) {
	return nil, errors.New("group functionality is not available in CE version")
}
