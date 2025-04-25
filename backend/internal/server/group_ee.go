//go:build ee
// +build ee

package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func groupFromModel(g *core.Group) *groupResponse {
	if g == nil {
		return nil
	}

	return &groupResponse{
		ID:        g.ID.String(),
		Name:      g.Name,
		Slug:      g.Slug,
		CreatedAt: strconv.FormatInt(g.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(g.UpdatedAt.Unix(), 10),
	}
}

func groupPageFromModel(g *core.GroupPage) *groupPageResponse {
	if g == nil {
		return nil
	}

	return &groupPageResponse{
		ID:        g.ID.String(),
		GroupID:   g.GroupID.String(),
		PageID:    g.PageID.String(),
		CreatedAt: strconv.FormatInt(g.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(g.UpdatedAt.Unix(), 10),
	}
}

type getGroupResponse struct {
	Group *groupResponse `json:"group"`
}

func (s *Server) handleGetGroup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	groupIDReq := chi.URLParam(r, "groupID")
	if groupIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("groupID is required"))
	}

	groupID, err := uuid.FromString(groupIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	ctxOrg := internal.ContextOrganization(ctx)
	group, err := s.db.Group().Get(ctx, database.GroupByOrganizationID(ctxOrg.ID), database.GroupByID(groupID))
	if err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, getGroupResponse{
		Group: groupFromModel(group),
	})
}

type listGroupsResponse struct {
	Groups     []*groupResponse     `json:"groups"`
	Users      []*userResponse      `json:"users"`
	UserGroups []*userGroupResponse `json:"userGroups"`
}

func (s *Server) handleListGroups(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	ctxOrg := internal.ContextOrganization(ctx)
	groups, err := s.db.Group().List(ctx, database.GroupByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}

	users, err := s.db.User().List(ctx, database.UserByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}

	userGroups, err := s.db.User().ListGroups(ctx, database.UserGroupByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}

	userIDs := make([]uuid.UUID, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	orgAccesses, err := s.db.User().ListOrganizationAccesses(ctx, database.UserOrganizationAccessByUserIDs(userIDs))
	if err != nil {
		return err
	}

	orgAccessesMap := make(map[uuid.UUID]*core.UserOrganizationAccess)
	for _, orgAccess := range orgAccesses {
		orgAccessesMap[orgAccess.UserID] = orgAccess
	}

	groupsOut := make([]*groupResponse, 0, len(groups))
	for _, group := range groups {
		groupsOut = append(groupsOut, groupFromModel(group))
	}

	usersOut := make([]*userResponse, 0, len(users))
	for _, u := range users {
		var role core.UserOrganizationRole
		orgAccess, ok := orgAccessesMap[u.ID]
		if ok {
			role = orgAccess.Role
		}
		usersOut = append(usersOut, userFromModel(u, role, ctxOrg))
	}

	userGroupsOut := make([]*userGroupResponse, 0, len(userGroups))
	for _, userGroup := range userGroups {
		userGroupsOut = append(userGroupsOut, userGroupFromModel(userGroup))
	}

	return s.renderJSON(w, http.StatusOK, listGroupsResponse{
		Groups:     groupsOut,
		Users:      usersOut,
		UserGroups: userGroupsOut,
	})
}

type createGroupRequest struct {
	Name    string   `json:"name" validate:"required"`
	Slug    string   `json:"slug" validate:"required"`
	UserIDs []string `json:"userIds" validate:"required"`
}

type createGroupResponse struct {
	Group *groupResponse `json:"group"`
}

func (s *Server) handleCreateGroup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req createGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	if err := s.checker.AuthorizeOperation(ctx, core.OperationEditGroup); err != nil {
		return err
	}

	ctxOrg := internal.ContextOrganization(ctx)

	slugExists, err := s.db.Group().IsSlugExistsInOrganization(ctx, ctxOrg.ID, req.Slug)
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
		OrganizationID: ctxOrg.ID,
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

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.Group().Create(ctx, g); err != nil {
			return err
		}

		if err := tx.User().BulkInsertGroups(ctx, userGroups); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	g, _ = s.db.Group().Get(ctx, database.GroupByID(g.ID))

	return s.renderJSON(w, http.StatusOK, createGroupResponse{
		Group: groupFromModel(g),
	})
}

type updateGroupRequest struct {
	Name    *string  `json:"name" validate:"required"`
	UserIDs []string `json:"userIds" validate:"required"`
}

type updateGroupResponse struct {
	Group *groupResponse `json:"group"`
}

func (s *Server) handleUpdateGroup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	groupIDReq := chi.URLParam(r, "groupID")
	if groupIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("groupID is required"))
	}

	var req updateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	ctxOrg := internal.ContextOrganization(ctx)
	groupID, err := uuid.FromString(groupIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	g, err := s.db.Group().Get(ctx, database.GroupByID(groupID), database.GroupByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}

	if req.Name != nil {
		g.Name = internal.StringValue(req.Name)
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

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.Group().Update(ctx, g); err != nil {
			return err
		}

		existingGroups, err := tx.User().ListGroups(ctx, database.UserGroupByGroupID(g.ID))
		if err != nil {
			return err
		}

		if err := tx.User().BulkDeleteGroups(ctx, existingGroups); err != nil {
			return err
		}

		if err := tx.User().BulkInsertGroups(ctx, userGroups); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	g, _ = s.db.Group().Get(ctx, database.GroupByID(g.ID))

	return s.renderJSON(w, http.StatusOK, updateGroupResponse{
		Group: groupFromModel(g),
	})
}

type deleteGroupResponse struct {
	Group *groupResponse `json:"group"`
}

func (s *Server) handleDeleteGroup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	groupIDReq := chi.URLParam(r, "groupID")
	if groupIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("groupID is required"))
	}

	groupID, err := uuid.FromString(groupIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	ctxOrg := internal.ContextOrganization(ctx)

	g, err := s.db.Group().Get(ctx, database.GroupByID(groupID), database.GroupByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.Group().Delete(ctx, g); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, deleteGroupResponse{
		Group: groupFromModel(g),
	})
}
