package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

type hostInstanceResponse struct {
	ID         string `json:"id"`
	SDKName    string `json:"sdkName"`
	SDKVersion string `json:"sdkVersion"`
	Status     string `json:"status"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

func (s *Server) hostInstanceFromModel(hostInstance *core.HostInstance) *hostInstanceResponse {
	if hostInstance == nil {
		return nil
	}

	return &hostInstanceResponse{
		ID:         hostInstance.ID.String(),
		SDKName:    hostInstance.SDKName,
		SDKVersion: hostInstance.SDKVersion,
		Status:     hostInstance.Status.String(),
		CreatedAt:  strconv.FormatInt(hostInstance.CreatedAt.Unix(), 10),
		UpdatedAt:  strconv.FormatInt(hostInstance.UpdatedAt.Unix(), 10),
	}
}

type pingHostInstanceResponse struct {
	HostInstance *hostInstanceResponse `json:"hostInstance"`
}

func (s *Server) handlePingHostInstance(w http.ResponseWriter, r *http.Request) error {
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

	return s.renderJSON(w, http.StatusOK, pingHostInstanceResponse{
		HostInstance: s.hostInstanceFromModel(onlineHostInstance),
	})
}
