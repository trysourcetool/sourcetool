package service

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/apikey"
	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/page"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
	"github.com/trysourcetool/sourcetool/backend/ws/conn"
)

type HostInstanceService interface {
	Ping(context.Context, input.PingHostInstanceInput) (*output.PingHostInstanceOutput, error)
}

type HostInstanceServiceCE struct {
	*infra.Dependency
}

func NewHostInstanceServiceCE(d *infra.Dependency) *HostInstanceServiceCE {
	return &HostInstanceServiceCE{Dependency: d}
}

func (s *HostInstanceServiceCE) Ping(ctx context.Context, in input.PingHostInstanceInput) (*output.PingHostInstanceOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)
	if currentOrg == nil {
		return nil, errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	hostInstanceOpts := []hostinstance.StoreOption{
		hostinstance.ByOrganizationID(currentOrg.ID),
	}
	if in.PageID != nil {
		pageID, err := uuid.FromString(conv.SafeValue(in.PageID))
		if err != nil {
			return nil, errdefs.ErrInvalidArgument(err)
		}

		page, err := s.Store.Page().Get(ctx, page.ByID(pageID))
		if err != nil {
			return nil, err
		}

		apiKey, err := s.Store.APIKey().Get(ctx, apikey.ByID(page.APIKeyID))
		if err != nil {
			return nil, err
		}

		hostInstanceOpts = append(hostInstanceOpts, hostinstance.ByAPIKeyID(apiKey.ID))
	}

	hostInstances, err := s.Store.HostInstance().List(ctx, hostInstanceOpts...)
	if err != nil {
		return nil, err
	}

	var onlineHostInstance *hostinstance.HostInstance
	for _, hostInstance := range hostInstances {
		if hostInstance.Status == hostinstance.HostInstanceStatusOnline {
			connManager := conn.GetConnManager()
			if err := connManager.PingHost(hostInstance.ID); err != nil {
				hostInstance.Status = hostinstance.HostInstanceStatusOffline
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

	return &output.PingHostInstanceOutput{
		HostInstance: output.HostInstanceFromModel(onlineHostInstance),
	}, nil
}
