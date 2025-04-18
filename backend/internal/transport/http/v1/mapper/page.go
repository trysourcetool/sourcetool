package mapper

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/responses"
)

func PageOutputToResponse(page *dto.Page) *responses.PageResponse {
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

func ListPagesRequestToInput(in requests.ListPagesRequest) dto.ListPagesInput {
	return dto.ListPagesInput{
		EnvironmentID: in.EnvironmentID,
	}
}

func ListPagesOutputToResponse(out *dto.ListPagesOutput) *responses.ListPagesResponse {
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
