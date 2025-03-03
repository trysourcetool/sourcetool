package responses

type PageResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Route     string `json:"route"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListPagesResponse struct {
	Pages      []*PageResponse      `json:"pages"`
	Groups     []*GroupResponse     `json:"groups"`
	GroupPages []*GroupPageResponse `json:"groupPages"`
	Users      []*UserResponse      `json:"users"`
	UserGroups []*UserGroupResponse `json:"userGroups"`
}
