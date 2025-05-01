//go:build ee
// +build ee

package server

import (
	"errors"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func (s *Server) handleListPages(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	envIDReq := r.URL.Query().Get("environmentId")
	if envIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("environmentId is required"))
	}

	ctxOrg := internal.ContextOrganization(ctx)

	envID, err := uuid.FromString(envIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	env, err := s.db.Environment().Get(ctx, database.EnvironmentByID(envID))
	if err != nil {
		return err
	}

	pagesOut, usersOut, userGroupsOut, err := s.handleListPagesBase(ctx, env, ctxOrg)
	if err != nil {
		return err
	}

	groups, err := s.db.Group().List(ctx, database.GroupByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}

	groupPages, err := s.db.Group().ListPages(ctx, database.GroupPageByOrganizationID(ctxOrg.ID), database.GroupPageByEnvironmentID(env.ID))
	if err != nil {
		return err
	}

	groupsOut := make([]*groupResponse, 0, len(groups))
	for _, group := range groups {
		groupsOut = append(groupsOut, s.groupFromModel(group))
	}

	groupPagesOut := make([]*groupPageResponse, 0, len(groupPages))
	for _, groupPage := range groupPages {
		groupPagesOut = append(groupPagesOut, s.groupPageFromModel(groupPage))
	}

	return s.renderJSON(w, http.StatusOK, listPagesResponse{
		Pages:      pagesOut,
		Groups:     groupsOut,
		GroupPages: groupPagesOut,
		Users:      usersOut,
		UserGroups: userGroupsOut,
	})
}
