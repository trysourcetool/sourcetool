package dto

import (
	"github.com/gofrs/uuid/v5"
	"github.com/trysourcetool/sourcetool/backend/model"
)

// Group represents group data in DTOs
type Group struct {
	ID             string
	OrganizationID string
	Name           string
	Slug           string
	CreatedAt      int64
	UpdatedAt      int64
}

// GroupFromModel converts from model.Group to dto.Group
func GroupFromModel(group *model.Group) *Group {
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

// ToGroupID converts string ID to uuid.UUID
func ToGroupID(id string) (uuid.UUID, error) {
	return uuid.FromString(id)
}

// GroupPage represents group page association in DTOs
type GroupPage struct {
	ID        string
	GroupID   string
	PageID    string
	CreatedAt int64
	UpdatedAt int64
}

// GroupPageFromModel converts from model.GroupPage to dto.GroupPage
func GroupPageFromModel(groupPage *model.GroupPage) *GroupPage {
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

// GetGroupInput is the input for Get operation
type GetGroupInput struct {
	GroupID string
}

// GetGroupOutput is the output for Get operation
type GetGroupOutput struct {
	Group *Group
}

// ListGroupsOutput is the output for List operation
type ListGroupsOutput struct {
	Groups     []*Group
	Users      []*User
	UserGroups []*UserGroup
}

// CreateGroupInput is the input for Create operation
type CreateGroupInput struct {
	Name    string
	Slug    string
	UserIDs []string
}

// CreateGroupOutput is the output for Create operation
type CreateGroupOutput struct {
	Group *Group
}

// UpdateGroupInput is the input for Update operation
type UpdateGroupInput struct {
	GroupID string
	Name    *string
	UserIDs []string
}

// UpdateGroupOutput is the output for Update operation
type UpdateGroupOutput struct {
	Group *Group
}

// DeleteGroupInput is the input for Delete operation
type DeleteGroupInput struct {
	GroupID string
}

// DeleteGroupOutput is the output for Delete operation
type DeleteGroupOutput struct {
	Group *Group
}
