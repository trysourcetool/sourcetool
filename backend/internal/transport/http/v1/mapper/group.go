package mapper

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/responses"
)

func GroupOutputToResponse(group *dto.Group) *responses.GroupResponse {
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

func GroupPageOutputToResponse(groupPage *dto.GroupPage) *responses.GroupPageResponse {
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

func GetGroupRequestToInput(in requests.GetGroupRequest) dto.GetGroupInput {
	return dto.GetGroupInput{
		GroupID: in.GroupID,
	}
}

func GetGroupOutputToResponse(out *dto.GetGroupOutput) *responses.GetGroupResponse {
	return &responses.GetGroupResponse{
		Group: GroupOutputToResponse(out.Group),
	}
}

func ListGroupsOutputToResponse(out *dto.ListGroupsOutput) *responses.ListGroupsResponse {
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

func CreateGroupRequestToInput(in requests.CreateGroupRequest) dto.CreateGroupInput {
	return dto.CreateGroupInput{
		Name:    in.Name,
		Slug:    in.Slug,
		UserIDs: in.UserIDs,
	}
}

func CreateGroupOutputToResponse(out *dto.CreateGroupOutput) *responses.CreateGroupResponse {
	return &responses.CreateGroupResponse{
		Group: GroupOutputToResponse(out.Group),
	}
}

func UpdateGroupRequestToInput(in requests.UpdateGroupRequest) dto.UpdateGroupInput {
	return dto.UpdateGroupInput{
		GroupID: in.GroupID,
		Name:    in.Name,
		UserIDs: in.UserIDs,
	}
}

func UpdateGroupOutputToResponse(out *dto.UpdateGroupOutput) *responses.UpdateGroupResponse {
	return &responses.UpdateGroupResponse{
		Group: GroupOutputToResponse(out.Group),
	}
}

func DeleteGroupRequestToInput(in requests.DeleteGroupRequest) dto.DeleteGroupInput {
	return dto.DeleteGroupInput{
		GroupID: in.GroupID,
	}
}

func DeleteGroupOutputToResponse(out *dto.DeleteGroupOutput) *responses.DeleteGroupResponse {
	return &responses.DeleteGroupResponse{
		Group: GroupOutputToResponse(out.Group),
	}
}
