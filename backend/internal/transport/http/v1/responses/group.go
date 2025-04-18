package responses

type GroupResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type GroupPageResponse struct {
	ID        string `json:"id"`
	GroupID   string `json:"groupId"`
	PageID    string `json:"pageId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListGroupsResponse struct {
	Groups     []*GroupResponse     `json:"groups"`
	Users      []*UserResponse      `json:"users"`
	UserGroups []*UserGroupResponse `json:"userGroups"`
}

type GetGroupResponse struct {
	Group *GroupResponse `json:"group"`
}

type CreateGroupResponse struct {
	Group *GroupResponse `json:"group"`
}

type UpdateGroupResponse struct {
	Group *GroupResponse `json:"group"`
}

type DeleteGroupResponse struct {
	Group *GroupResponse `json:"group"`
}
