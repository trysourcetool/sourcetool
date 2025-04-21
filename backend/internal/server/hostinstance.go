package server

import (
	"errors"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/postgres"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

func (s *Server) pingHostInstance(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	pageIDReq := internal.NilValue(r.URL.Query().Get("pageId"))
	currentOrg := internal.CurrentOrganization(ctx)
	if currentOrg == nil {
		return errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	hostInstanceOpts := []postgres.HostInstanceQuery{
		postgres.HostInstanceByOrganizationID(currentOrg.ID),
	}
	if pageIDReq != nil {
		pageID, err := uuid.FromString(internal.SafeValue(pageIDReq))
		if err != nil {
			return errdefs.ErrInvalidArgument(err)
		}

		page, err := s.db.GetPage(ctx, postgres.PageByID(pageID))
		if err != nil {
			return err
		}

		apiKey, err := s.db.GetAPIKey(ctx, postgres.APIKeyByID(page.APIKeyID))
		if err != nil {
			return err
		}

		hostInstanceOpts = append(hostInstanceOpts, postgres.HostInstanceByAPIKeyID(apiKey.ID))
	}

	hostInstances, err := s.db.ListHostInstances(ctx, hostInstanceOpts...)
	if err != nil {
		return err
	}

	var onlineHostInstance *core.HostInstance
	for _, hostInstance := range hostInstances {
		if hostInstance.Status == core.HostInstanceStatusOnline {
			if err := s.wsManager.PingConnectedHost(hostInstance.ID); err != nil {
				hostInstance.Status = core.HostInstanceStatusOffline

				if err := s.db.UpdateHostInstance(ctx, nil, hostInstance); err != nil {
					return err
				}

				continue
			}

			onlineHostInstance = hostInstance
			break
		}
	}

	if onlineHostInstance == nil {
		return errdefs.ErrHostInstanceStatusNotOnline(errors.New("host instance status is not online"))
	}

	return s.renderJSON(w, http.StatusOK, responses.PingHostInstanceResponse{
		HostInstance: responses.HostInstanceFromModel(onlineHostInstance),
	})
}
