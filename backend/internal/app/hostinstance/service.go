package hostinstance

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	"github.com/trysourcetool/sourcetool/backend/internal/ctxutil"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/page"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
	"github.com/trysourcetool/sourcetool/backend/pkg/ptrconv"
)

type Service interface {
	Ping(context.Context, dto.PingHostInstanceInput) (*dto.PingHostInstanceOutput, error)
}

type ServiceCE struct {
	*port.Dependencies
}

func NewServiceCE(d *port.Dependencies) *ServiceCE {
	return &ServiceCE{Dependencies: d}
}

func (s *ServiceCE) Ping(ctx context.Context, in dto.PingHostInstanceInput) (*dto.PingHostInstanceOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)
	if currentOrg == nil {
		return nil, errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	hostInstanceOpts := []hostinstance.Query{
		hostinstance.ByOrganizationID(currentOrg.ID),
	}
	if in.PageID != nil {
		pageID, err := uuid.FromString(ptrconv.SafeValue(in.PageID))
		if err != nil {
			return nil, errdefs.ErrInvalidArgument(err)
		}

		page, err := s.Repository.Page().Get(ctx, page.ByID(pageID))
		if err != nil {
			return nil, err
		}

		apiKey, err := s.Repository.APIKey().Get(ctx, apikey.ByID(page.APIKeyID))
		if err != nil {
			return nil, err
		}

		hostInstanceOpts = append(hostInstanceOpts, hostinstance.ByAPIKeyID(apiKey.ID))
	}

	hostInstances, err := s.Repository.HostInstance().List(ctx, hostInstanceOpts...)
	if err != nil {
		return nil, err
	}

	var onlineHostInstance *hostinstance.HostInstance
	for _, hostInstance := range hostInstances {
		if hostInstance.Status == hostinstance.HostInstanceStatusOnline {
			if err := s.WSManager.PingConnectedHost(hostInstance.ID); err != nil {
				hostInstance.Status = hostinstance.HostInstanceStatusOffline
				if err := s.Repository.HostInstance().Update(ctx, hostInstance); err != nil {
					return nil, err
				}
				continue
			}

			onlineHostInstance = hostInstance
			break
		}
	}

	if onlineHostInstance == nil {
		return nil, errdefs.ErrHostInstanceStatusNotOnline(errors.New("host instance status is not online"))
	}

	return &dto.PingHostInstanceOutput{
		HostInstance: dto.HostInstanceFromModel(onlineHostInstance),
	}, nil
}
