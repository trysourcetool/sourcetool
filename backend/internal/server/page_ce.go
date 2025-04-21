//go:build !ee
// +build !ee

package server

import (
	"errors"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

func (s *Server) listPages(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	envIDReq := r.URL.Query().Get("environmentId")
	if envIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("environmentId is required"))
	}

	o := internal.CurrentOrganization(ctx)

	envID, err := uuid.FromString(envIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	env, err := s.db.Environment().Get(ctx, database.EnvironmentByID(envID))
	if err != nil {
		return err
	}

	pagesOut, usersOut, userGroupsOut, err := s.listPagesBase(ctx, env, o)
	if err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, responses.ListPagesResponse{
		Pages:      pagesOut,
		Groups:     make([]*responses.GroupResponse, 0),
		GroupPages: make([]*responses.GroupPageResponse, 0),
		Users:      usersOut,
		UserGroups: userGroupsOut,
	})
}
