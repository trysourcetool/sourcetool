package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto/http/requests"
	"github.com/trysourcetool/sourcetool/backend/dto/http/responses"
	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
)

// GroupOutputToResponse converts from output.Group to responses.GroupResponse.
func GroupOutputToResponse(group *output.Group) *responses.GroupResponse {
	if group == nil {
		return nil
	}

	return &responses.GroupResponse{
		ID:        group.ID,
		Name:      group.Name,
		Slug:      group.Slug,
		CreatedAt: strconv.FormatInt(group.CreatedAt, 10),
		UpdatedAt: strconv.FormatInt(group.UpdatedAt, 10),
	}
}

// GroupPageOutputToResponse converts from output.GroupPage to responses.GroupPageResponse.
func GroupPageOutputToResponse(groupPage *output.GroupPage) *responses.GroupPageResponse {
	if groupPage == nil {
		return nil
	}

	return &responses.GroupPageResponse{
		ID:        groupPage.ID,
		GroupID:   groupPage.GroupID,
		PageID:    groupPage.PageID,
		CreatedAt: strconv.FormatInt(groupPage.CreatedAt, 10),
		UpdatedAt: strconv.FormatInt(groupPage.UpdatedAt, 10),
	}
}

// GetGroupRequestToInput converts from requests.GetGroupRequest to input.GetGroupInput.
func GetGroupRequestToInput(in requests.GetGroupRequest) input.GetGroupInput {
	return input.GetGroupInput{
		GroupID: in.GroupID,
	}
}

// GetGroupOutputToResponse converts from output.GetGroupOutput to responses.GetGroupResponse.
func GetGroupOutputToResponse(out *output.GetGroupOutput) *responses.GetGroupResponse {
	return &responses.GetGroupResponse{
		Group: GroupOutputToResponse(out.Group),
	}
}

// ListGroupsOutputToResponse converts from output.ListGroupsOutput to responses.ListGroupsResponse.
func ListGroupsOutputToResponse(out *output.ListGroupsOutput) *responses.ListGroupsResponse {
	groups := make([]*responses.GroupResponse, 0, len(out.Groups))
	for _, group := range out.Groups {
		groups = append(groups, GroupOutputToResponse(group))
	}

	users := make([]*responses.UserResponse, 0, len(out.Users))
	for _, user := range out.Users {
		users = append(users, UserOutputToResponse(user))
	}

	userGroups := make([]*responses.UserGroupResponse, 0, len(out.UserGroups))
	for _, userGroup := range out.UserGroups {
		userGroups = append(userGroups, UserGroupOutputToResponse(userGroup))
	}

	return &responses.ListGroupsResponse{
		Groups:     groups,
		Users:      users,
		UserGroups: userGroups,
	}
}

// CreateGroupRequestToInput converts from requests.CreateGroupRequest to input.CreateGroupInput.
func CreateGroupRequestToInput(in requests.CreateGroupRequest) input.CreateGroupInput {
	return input.CreateGroupInput{
		Name:    in.Name,
		Slug:    in.Slug,
		UserIDs: in.UserIDs,
	}
}

// CreateGroupOutputToResponse converts from output.CreateGroupOutput to responses.CreateGroupResponse.
func CreateGroupOutputToResponse(out *output.CreateGroupOutput) *responses.CreateGroupResponse {
	return &responses.CreateGroupResponse{
		Group: GroupOutputToResponse(out.Group),
	}
}

// UpdateGroupRequestToInput converts from requests.UpdateGroupRequest to input.UpdateGroupInput.
func UpdateGroupRequestToInput(in requests.UpdateGroupRequest) input.UpdateGroupInput {
	return input.UpdateGroupInput{
		GroupID: in.GroupID,
		Name:    in.Name,
		UserIDs: in.UserIDs,
	}
}

// UpdateGroupOutputToResponse converts from output.UpdateGroupOutput to responses.UpdateGroupResponse.
func UpdateGroupOutputToResponse(out *output.UpdateGroupOutput) *responses.UpdateGroupResponse {
	return &responses.UpdateGroupResponse{
		Group: GroupOutputToResponse(out.Group),
	}
}

// DeleteGroupRequestToInput converts from requests.DeleteGroupRequest to input.DeleteGroupInput.
func DeleteGroupRequestToInput(in requests.DeleteGroupRequest) input.DeleteGroupInput {
	return input.DeleteGroupInput{
		GroupID: in.GroupID,
	}
}

// DeleteGroupOutputToResponse converts from output.DeleteGroupOutput to responses.DeleteGroupResponse.
func DeleteGroupOutputToResponse(out *output.DeleteGroupOutput) *responses.DeleteGroupResponse {
	return &responses.DeleteGroupResponse{
		Group: GroupOutputToResponse(out.Group),
	}
}
