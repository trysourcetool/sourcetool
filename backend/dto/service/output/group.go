package output

import (
	"github.com/trysourcetool/sourcetool/backend/group"
)

// Group represents group data in DTOs.
type Group struct {
	ID             string
	OrganizationID string
	Name           string
	Slug           string
	CreatedAt      int64
	UpdatedAt      int64
}

// GroupFromModel converts from model.Group to dto.Group.
func GroupFromModel(group *group.Group) *Group {
	if group == nil {
		return nil
	}

	return &Group{
		ID:             group.ID.String(),
		OrganizationID: group.OrganizationID.String(),
		Name:           group.Name,
		Slug:           group.Slug,
		CreatedAt:      group.CreatedAt.Unix(),
		UpdatedAt:      group.UpdatedAt.Unix(),
	}
}

// GroupPage represents group page association in DTOs.
type GroupPage struct {
	ID        string
	GroupID   string
	PageID    string
	CreatedAt int64
	UpdatedAt int64
}

// GroupPageFromModel converts from model.GroupPage to dto.GroupPage.
func GroupPageFromModel(groupPage *group.GroupPage) *GroupPage {
	if groupPage == nil {
		return nil
	}

	return &GroupPage{
		ID:        groupPage.ID.String(),
		GroupID:   groupPage.GroupID.String(),
		PageID:    groupPage.PageID.String(),
		CreatedAt: groupPage.CreatedAt.Unix(),
		UpdatedAt: groupPage.UpdatedAt.Unix(),
	}
}

// GetGroupOutput is the output for Get operation.
type GetGroupOutput struct {
	Group *Group
}

// ListGroupsOutput is the output for List operation.
type ListGroupsOutput struct {
	Groups     []*Group
	Users      []*User
	UserGroups []*UserGroup
}

// CreateGroupOutput is the output for Create operation.
type CreateGroupOutput struct {
	Group *Group
}

// UpdateGroupOutput is the output for Update operation.
type UpdateGroupOutput struct {
	Group *Group
}

// DeleteGroupOutput is the output for Delete operation.
type DeleteGroupOutput struct {
	Group *Group
}
