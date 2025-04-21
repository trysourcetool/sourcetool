package responses

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type GroupResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func GroupFromModel(g *core.Group) *GroupResponse {
	return &GroupResponse{
		ID:        g.ID.String(),
		Name:      g.Name,
		Slug:      g.Slug,
		CreatedAt: strconv.FormatInt(g.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(g.UpdatedAt.Unix(), 10),
	}
}

type GroupPageResponse struct {
	ID        string `json:"id"`
	GroupID   string `json:"groupId"`
	PageID    string `json:"pageId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func GroupPageFromModel(g *core.GroupPage) *GroupPageResponse {
	return &GroupPageResponse{
		ID:        g.ID.String(),
		GroupID:   g.GroupID.String(),
		PageID:    g.PageID.String(),
		CreatedAt: strconv.FormatInt(g.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(g.UpdatedAt.Unix(), 10),
	}
}

type ListGroupsResponse struct {
	Groups     []*GroupResponse     `json:"groups"`
	Users      []*UserResponse      `json:"users"`
	UserGroups []*UserGroupResponse `json:"userGroups"`
}

type GetGroupResponse struct {
	Group *GroupResponse `json:"group"`
}

type CreateGroupResponse struct {
	Group *GroupResponse `json:"group"`
}

type UpdateGroupResponse struct {
	Group *GroupResponse `json:"group"`
}

type DeleteGroupResponse struct {
	Group *GroupResponse `json:"group"`
}
