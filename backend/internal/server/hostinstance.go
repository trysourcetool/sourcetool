package server

import (
	"errors"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

func (s *Server) pingHostInstance(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	pageIDReq := internal.StringPtr(r.URL.Query().Get("pageId"))
	ctxOrg := internal.ContextOrganization(ctx)
	if ctxOrg == nil {
		return errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	hostInstanceOpts := []database.HostInstanceQuery{
		database.HostInstanceByOrganizationID(ctxOrg.ID),
	}
	if pageIDReq != nil {
		pageID, err := uuid.FromString(internal.StringValue(pageIDReq))
		if err != nil {
			return errdefs.ErrInvalidArgument(err)
		}

		page, err := s.db.Page().Get(ctx, database.PageByID(pageID))
		if err != nil {
			return err
		}

		apiKey, err := s.db.APIKey().Get(ctx, database.APIKeyByID(page.APIKeyID))
		if err != nil {
			return err
		}

		hostInstanceOpts = append(hostInstanceOpts, database.HostInstanceByAPIKeyID(apiKey.ID))
	}

	hostInstances, err := s.db.HostInstance().List(ctx, hostInstanceOpts...)
	if err != nil {
		return err
	}

	var onlineHostInstance *core.HostInstance
	for _, hostInstance := range hostInstances {
		if hostInstance.Status == core.HostInstanceStatusOnline {
			if err := s.wsManager.PingConnectedHost(hostInstance.ID); err != nil {
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
