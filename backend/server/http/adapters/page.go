package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto/http/requests"
	"github.com/trysourcetool/sourcetool/backend/dto/http/responses"
	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
)

// PageOutputToResponse converts from output.Page to responses.PageResponse.
func PageOutputToResponse(page *output.Page) *responses.PageResponse {
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

// ListPagesRequestToInput converts from requests.ListPagesRequest to input.ListPagesInput.
func ListPagesRequestToInput(in requests.ListPagesRequest) input.ListPagesInput {
	return input.ListPagesInput{
		EnvironmentID: in.EnvironmentID,
	}
}

// ListPagesOutputToResponse converts from output.ListPagesOutput to responses.ListPagesResponse.
func ListPagesOutputToResponse(out *output.ListPagesOutput) *responses.ListPagesResponse {
	pages := make([]*responses.PageResponse, 0, len(out.Pages))
	for _, page := range out.Pages {
		pages = append(pages, PageOutputToResponse(page))
	}

	groups := make([]*responses.GroupResponse, 0, len(out.Groups))
	for _, group := range out.Groups {
		groups = append(groups, GroupOutputToResponse(group))
	}

	groupPages := make([]*responses.GroupPageResponse, 0, len(out.GroupPages))
	for _, groupPage := range out.GroupPages {
		groupPages = append(groupPages, GroupPageOutputToResponse(groupPage))
	}

	users := make([]*responses.UserResponse, 0, len(out.Users))
	for _, user := range out.Users {
		users = append(users, UserOutputToResponse(user))
	}

	userGroups := make([]*responses.UserGroupResponse, 0, len(out.UserGroups))
	for _, userGroup := range out.UserGroups {
		userGroups = append(userGroups, UserGroupOutputToResponse(userGroup))
	}

	return &responses.ListPagesResponse{
		Pages:      pages,
		Groups:     groups,
		GroupPages: groupPages,
		Users:      users,
		UserGroups: userGroups,
	}
}
