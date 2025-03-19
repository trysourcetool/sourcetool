package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
)

// PageDTOToResponse converts from dto.Page to responses.PageResponse.
func PageDTOToResponse(page *dto.Page) *responses.PageResponse {
	if page == nil {
		return nil
	}

	return &responses.PageResponse{
		ID:        page.ID,
		Name:      page.Name,
		Route:     page.Route,
		CreatedAt: strconv.FormatInt(page.CreatedAt, 10),
		UpdatedAt: strconv.FormatInt(page.UpdatedAt, 10),
	}
}

// ListPagesRequestToDTOInput converts from requests.ListPagesRequest to dto.ListPagesInput.
func ListPagesRequestToDTOInput(in *requests.ListPagesRequest) *dto.ListPagesInput {
	return &dto.ListPagesInput{
		EnvironmentID: in.EnvironmentID,
	}
}

// ListPagesOutputToResponse converts from dto.ListPagesOutput to responses.ListPagesResponse.
func ListPagesOutputToResponse(out *dto.ListPagesOutput) *responses.ListPagesResponse {
	pages := make([]*responses.PageResponse, 0, len(out.Pages))
	for _, page := range out.Pages {
		pages = append(pages, PageDTOToResponse(page))
	}

	groups := make([]*responses.GroupResponse, 0, len(out.Groups))
	for _, group := range out.Groups {
		groups = append(groups, GroupDTOToResponse(group))
	}

	groupPages := make([]*responses.GroupPageResponse, 0, len(out.GroupPages))
	for _, groupPage := range out.GroupPages {
		groupPages = append(groupPages, GroupPageDTOToResponse(groupPage))
	}

	users := make([]*responses.UserResponse, 0, len(out.Users))
	for _, user := range out.Users {
		users = append(users, UserDTOToResponse(user))
	}

	userGroups := make([]*responses.UserGroupResponse, 0, len(out.UserGroups))
	for _, userGroup := range out.UserGroups {
		userGroups = append(userGroups, UserGroupDTOToResponse(userGroup))
	}

	return &responses.ListPagesResponse{
		Pages:      pages,
		Groups:     groups,
		GroupPages: groupPages,
		Users:      users,
		UserGroups: userGroups,
	}
}
