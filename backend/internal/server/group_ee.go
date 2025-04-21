//go:build ee
// +build ee

package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/permission"
	"github.com/trysourcetool/sourcetool/backend/internal/postgres"
	"github.com/trysourcetool/sourcetool/backend/internal/server/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

func (s *Server) getGroup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	groupIDReq := chi.URLParam(r, "groupID")
	if groupIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("groupID is required"))
	}

	groupID, err := uuid.FromString(groupIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	currentOrg := internal.CurrentOrganization(ctx)
	group, err := s.db.GetGroup(ctx, postgres.GroupByOrganizationID(currentOrg.ID), postgres.GroupByID(groupID))
	if err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, responses.GetGroupResponse{
		Group: responses.GroupFromModel(group),
	})
}

func (s *Server) listGroups(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	currentOrg := internal.CurrentOrganization(ctx)
	groups, err := s.db.ListGroups(ctx, postgres.GroupByOrganizationID(currentOrg.ID))
	if err != nil {
		return err
	}

	users, err := s.db.ListUsers(ctx, postgres.UserByOrganizationID(currentOrg.ID))
	if err != nil {
		return err
	}

	userGroups, err := s.db.ListUserGroups(ctx, postgres.UserGroupByOrganizationID(currentOrg.ID))
	if err != nil {
		return err
	}

	userIDs := make([]uuid.UUID, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	orgAccesses, err := s.db.ListUserOrganizationAccesses(ctx, postgres.UserOrganizationAccessByUserIDs(userIDs))
	if err != nil {
		return err
	}

	orgAccessesMap := make(map[uuid.UUID]*core.UserOrganizationAccess)
	for _, orgAccess := range orgAccesses {
		orgAccessesMap[orgAccess.UserID] = orgAccess
	}

	groupsOut := make([]*responses.GroupResponse, 0, len(groups))
	for _, group := range groups {
		groupsOut = append(groupsOut, responses.GroupFromModel(group))
	}

	usersOut := make([]*responses.UserResponse, 0, len(users))
	for _, u := range users {
		var role core.UserOrganizationRole
		orgAccess, ok := orgAccessesMap[u.ID]
		if ok {
			role = orgAccess.Role
		}
		usersOut = append(usersOut, responses.UserFromModel(u, role, nil))
	}

	userGroupsOut := make([]*responses.UserGroupResponse, 0, len(userGroups))
	for _, userGroup := range userGroups {
		userGroupsOut = append(userGroupsOut, responses.UserGroupFromModel(userGroup))
	}

	return s.renderJSON(w, http.StatusOK, responses.ListGroupsResponse{
		Groups:     groupsOut,
		Users:      usersOut,
		UserGroups: userGroupsOut,
	})
}

func (s *Server) createGroup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req requests.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	if err := s.checker.AuthorizeOperation(ctx, permission.OperationEditGroup); err != nil {
		return err
	}

	currentOrg := internal.CurrentOrganization(ctx)

	slugExists, err := s.db.IsGroupSlugExistsInOrganization(ctx, currentOrg.ID, req.Slug)
	if err != nil {
		return err
	}
	if slugExists {
		return errdefs.ErrGroupSlugAlreadyExists(errors.New("slug already exists"))
	}

	if !validateSlug(req.Slug) {
		return errdefs.ErrInvalidArgument(errors.New("invalid slug"))
	}

	g := &core.Group{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: currentOrg.ID,
		Name:           req.Name,
		Slug:           req.Slug,
	}

	userIDs := make([]uuid.UUID, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		id, err := uuid.FromString(userID)
		if err != nil {
			return errdefs.ErrInvalidArgument(err)
		}
		userIDs = append(userIDs, id)
	}

	userGroups := make([]*core.UserGroup, 0, len(userIDs))
	for _, userID := range userIDs {
		userGroups = append(userGroups, &core.UserGroup{
			ID:      uuid.Must(uuid.NewV4()),
			UserID:  userID,
			GroupID: g.ID,
		})
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.db.CreateGroup(ctx, tx, g); err != nil {
		return err
	}

	if err := s.db.BulkInsertUserGroups(ctx, tx, userGroups); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	g, _ = s.db.GetGroup(ctx, postgres.GroupByID(g.ID))

	return s.renderJSON(w, http.StatusOK, responses.CreateGroupResponse{
		Group: responses.GroupFromModel(g),
	})
}

func (s *Server) updateGroup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	groupIDReq := chi.URLParam(r, "groupID")
	if groupIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("groupID is required"))
	}

	var req requests.UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	if err := s.checker.AuthorizeOperation(ctx, permission.OperationEditGroup); err != nil {
		return err
	}

	currentOrg := internal.CurrentOrganization(ctx)
	groupID, err := uuid.FromString(groupIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	g, err := s.db.GetGroup(ctx, postgres.GroupByID(groupID), postgres.GroupByOrganizationID(currentOrg.ID))
	if err != nil {
		return err
	}

	if req.Name != nil {
		g.Name = internal.SafeValue(req.Name)
	}

	userIDs := make([]uuid.UUID, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		id, err := uuid.FromString(userID)
		if err != nil {
			return errdefs.ErrInvalidArgument(err)
		}
		userIDs = append(userIDs, id)
	}

	userGroups := make([]*core.UserGroup, 0, len(userIDs))
	for _, userID := range userIDs {
		userGroups = append(userGroups, &core.UserGroup{
			ID:      uuid.Must(uuid.NewV4()),
			UserID:  userID,
			GroupID: g.ID,
		})
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.db.UpdateGroup(ctx, tx, g); err != nil {
		return err
	}

	existingGroups, err := s.db.ListUserGroups(ctx, postgres.UserGroupByGroupID(g.ID))
	if err != nil {
		return err
	}

	if err := s.db.BulkDeleteUserGroups(ctx, tx, existingGroups); err != nil {
		return err
	}

	if err := s.db.BulkInsertUserGroups(ctx, tx, userGroups); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	g, _ = s.db.GetGroup(ctx, postgres.GroupByID(g.ID))

	return s.renderJSON(w, http.StatusOK, responses.UpdateGroupResponse{
		Group: responses.GroupFromModel(g),
	})
}

func (s *Server) deleteGroup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	groupIDReq := chi.URLParam(r, "groupID")
	if groupIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("groupID is required"))
	}

	groupID, err := uuid.FromString(groupIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	currentOrg := internal.CurrentOrganization(ctx)

	g, err := s.db.GetGroup(ctx, postgres.GroupByID(groupID), postgres.GroupByOrganizationID(currentOrg.ID))
	if err != nil {
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.db.DeleteGroup(ctx, tx, g); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, responses.DeleteGroupResponse{
		Group: responses.GroupFromModel(g),
	})
}
