package responses

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type PageResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Route     string `json:"route"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func PageFromModel(page *core.Page) *PageResponse {
	return &PageResponse{
		ID:        page.ID.String(),
		Name:      page.Name,
		Route:     page.Route,
		CreatedAt: strconv.FormatInt(page.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(page.UpdatedAt.Unix(), 10),
	}
}

type ListPagesResponse struct {
	Pages      []*PageResponse      `json:"pages"`
	Groups     []*GroupResponse     `json:"groups"`
	GroupPages []*GroupPageResponse `json:"groupPages"`
	Users      []*UserResponse      `json:"users"`
	UserGroups []*UserGroupResponse `json:"userGroups"`
}
