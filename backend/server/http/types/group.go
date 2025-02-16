package types

type GroupPayload struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type GroupPagePayload struct {
	ID        string `json:"id"`
	GroupID   string `json:"groupId"`
	PageID    string `json:"pageId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListGroupsPayload struct {
	Groups     []*GroupPayload     `json:"groups"`
	Users      []*UserPayload      `json:"users"`
	UserGroups []*UserGroupPayload `json:"userGroups"`
}

type GetGroupInput struct {
	GroupID string `json:"-" validate:"required"`
}

type GetGroupPayload struct {
	Group *GroupPayload `json:"group"`
}

type CreateGroupInput struct {
	Name    string   `json:"name" validate:"required"`
	Slug    string   `json:"slug" validate:"required"`
	UserIDs []string `json:"userIds" validate:"required"`
}

type CreateGroupPayload struct {
	Group *GroupPayload `json:"group"`
}

type UpdateGroupInput struct {
	GroupID string   `json:"-" validate:"required"`
	Name    *string  `json:"name" validate:"required"`
	UserIDs []string `json:"userIds" validate:"required"`
}

type UpdateGroupPayload struct {
	Group *GroupPayload `json:"group"`
}

type DeleteGroupInput struct {
	GroupID string `json:"groupId" validate:"required"`
}

type DeleteGroupPayload struct {
	Group *GroupPayload `json:"group"`
}
