package hostinstance

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
	"github.com/trysourcetool/sourcetool/backend/ws"
)

type Service interface {
	Ping(context.Context, dto.PingHostInstanceInput) (*dto.PingHostInstanceOutput, error)
}

type ServiceCE struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) Ping(ctx context.Context, in dto.PingHostInstanceInput) (*dto.PingHostInstanceOutput, error) {
	subdomain, err := httputils.GetSubdomainFromHost(ctxutils.HTTPHost(ctx))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	o, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	hostInstanceOpts := []storeopts.HostInstanceOption{
		storeopts.HostInstanceByOrganizationID(o.ID),
	}
	if in.PageID != nil {
		pageID, err := uuid.FromString(conv.SafeValue(in.PageID))
		if err != nil {
			return nil, errdefs.ErrInvalidArgument(err)
		}

		page, err := s.Store.Page().Get(ctx, storeopts.PageByID(pageID))
		if err != nil {
			return nil, err
		}

		apiKey, err := s.Store.APIKey().Get(ctx, storeopts.APIKeyByID(page.APIKeyID))
		if err != nil {
			return nil, err
		}

		hostInstanceOpts = append(hostInstanceOpts, storeopts.HostInstanceByAPIKeyID(apiKey.ID))
	}

	hostInstances, err := s.Store.HostInstance().List(ctx, hostInstanceOpts...)
	if err != nil {
		return nil, err
	}

	var onlineHostInstance *model.HostInstance
	for _, hostInstance := range hostInstances {
		if hostInstance.Status == model.HostInstanceStatusOnline {
			connManager := ws.GetConnManager()
			if err := connManager.PingHost(hostInstance.ID); err != nil {
				hostInstance.Status = model.HostInstanceStatusOffline
				if err := s.Store.HostInstance().Update(ctx, hostInstance); err != nil {
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
