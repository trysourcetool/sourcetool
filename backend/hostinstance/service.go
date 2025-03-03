package hostinstance

import (
	"context"
	"errors"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/conv"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
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
	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]

	o, err := s.Store.Organization().Get(ctx, model.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	hostInstanceConds := []any{
		model.HostInstanceByOrganizationID(o.ID),
	}
	if in.PageID != nil {
		pageID, err := uuid.FromString(conv.SafeValue(in.PageID))
		if err != nil {
			return nil, errdefs.ErrInvalidArgument(err)
		}

		page, err := s.Store.Page().Get(ctx, model.PageByID(pageID))
		if err != nil {
			return nil, err
		}

		apiKey, err := s.Store.APIKey().Get(ctx, model.APIKeyByID(page.APIKeyID))
		if err != nil {
			return nil, err
		}

		hostInstanceConds = append(hostInstanceConds, model.HostInstanceByAPIKeyID(apiKey.ID))
	}

	hostInstances, err := s.Store.HostInstance().List(ctx, hostInstanceConds...)
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
