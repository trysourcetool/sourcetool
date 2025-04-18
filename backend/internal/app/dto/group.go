package dto

import "github.com/trysourcetool/sourcetool/backend/internal/domain/group"

type GetGroupInput struct {
	GroupID string
}

type CreateGroupInput struct {
	Name    string
	Slug    string
	UserIDs []string
}

type UpdateGroupInput struct {
	GroupID string
	Name    *string
	UserIDs []string
}

type DeleteGroupInput struct {
	GroupID string
}

type Group struct {
	ID             string
	OrganizationID string
	Name           string
	Slug           string
	CreatedAt      int64
	UpdatedAt      int64
}

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

type GroupPage struct {
	ID        string
	GroupID   string
	PageID    string
	CreatedAt int64
	UpdatedAt int64
}

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

type GetGroupOutput struct {
	Group *Group
}

type ListGroupsOutput struct {
	Groups     []*Group
	Users      []*User
	UserGroups []*UserGroup
}

type CreateGroupOutput struct {
	Group *Group
}

type UpdateGroupOutput struct {
	Group *Group
}

type DeleteGroupOutput struct {
	Group *Group
}
