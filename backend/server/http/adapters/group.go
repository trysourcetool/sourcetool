package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
)

// GroupDTOToResponse converts from dto.Group to responses.GroupResponse.
func GroupDTOToResponse(group *dto.Group) *responses.GroupResponse {
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

// GroupPageDTOToResponse converts from dto.GroupPage to responses.GroupPageResponse.
func GroupPageDTOToResponse(groupPage *dto.GroupPage) *responses.GroupPageResponse {
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

// GetGroupRequestToDTOInput converts from requests.GetGroupRequest to dto.GetGroupInput.
func GetGroupRequestToDTOInput(in requests.GetGroupRequest) dto.GetGroupInput {
	return dto.GetGroupInput{
		GroupID: in.GroupID,
	}
}

// GetGroupOutputToResponse converts from dto.GetGroupOutput to responses.GetGroupResponse.
func GetGroupOutputToResponse(out *dto.GetGroupOutput) *responses.GetGroupResponse {
	return &responses.GetGroupResponse{
		Group: GroupDTOToResponse(out.Group),
	}
}

// ListGroupsOutputToResponse converts from dto.ListGroupsOutput to responses.ListGroupsResponse.
func ListGroupsOutputToResponse(out *dto.ListGroupsOutput) *responses.ListGroupsResponse {
	groups := make([]*responses.GroupResponse, 0, len(out.Groups))
	for _, group := range out.Groups {
		groups = append(groups, GroupDTOToResponse(group))
	}

	users := make([]*responses.UserResponse, 0, len(out.Users))
	for _, user := range out.Users {
		users = append(users, UserDTOToResponse(user))
	}

	userGroups := make([]*responses.UserGroupResponse, 0, len(out.UserGroups))
	for _, userGroup := range out.UserGroups {
		userGroups = append(userGroups, UserGroupDTOToResponse(userGroup))
	}

	return &responses.ListGroupsResponse{
		Groups:     groups,
		Users:      users,
		UserGroups: userGroups,
	}
}

// CreateGroupRequestToDTOInput converts from requests.CreateGroupRequest to dto.CreateGroupInput.
func CreateGroupRequestToDTOInput(in requests.CreateGroupRequest) dto.CreateGroupInput {
	return dto.CreateGroupInput{
		Name:    in.Name,
		Slug:    in.Slug,
		UserIDs: in.UserIDs,
	}
}

// CreateGroupOutputToResponse converts from dto.CreateGroupOutput to responses.CreateGroupResponse.
func CreateGroupOutputToResponse(out *dto.CreateGroupOutput) *responses.CreateGroupResponse {
	return &responses.CreateGroupResponse{
		Group: GroupDTOToResponse(out.Group),
	}
}

// UpdateGroupRequestToDTOInput converts from requests.UpdateGroupRequest to dto.UpdateGroupInput.
func UpdateGroupRequestToDTOInput(in requests.UpdateGroupRequest) dto.UpdateGroupInput {
	return dto.UpdateGroupInput{
		GroupID: in.GroupID,
		Name:    in.Name,
		UserIDs: in.UserIDs,
	}
}

// UpdateGroupOutputToResponse converts from dto.UpdateGroupOutput to responses.UpdateGroupResponse.
func UpdateGroupOutputToResponse(out *dto.UpdateGroupOutput) *responses.UpdateGroupResponse {
	return &responses.UpdateGroupResponse{
		Group: GroupDTOToResponse(out.Group),
	}
}

// DeleteGroupRequestToDTOInput converts from requests.DeleteGroupRequest to dto.DeleteGroupInput.
func DeleteGroupRequestToDTOInput(in requests.DeleteGroupRequest) dto.DeleteGroupInput {
	return dto.DeleteGroupInput{
		GroupID: in.GroupID,
	}
}

// DeleteGroupOutputToResponse converts from dto.DeleteGroupOutput to responses.DeleteGroupResponse.
func DeleteGroupOutputToResponse(out *dto.DeleteGroupOutput) *responses.DeleteGroupResponse {
	return &responses.DeleteGroupResponse{
		Group: GroupDTOToResponse(out.Group),
	}
}
