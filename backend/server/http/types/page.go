package types

type PagePayload struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Route     string `json:"route"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListPagesPayload struct {
	Pages      []*PagePayload      `json:"pages"`
	Groups     []*GroupPayload     `json:"groups"`
	GroupPages []*GroupPagePayload `json:"groupPages"`
	Users      []*UserPayload      `json:"users"`
	UserGroups []*UserGroupPayload `json:"userGroups"`
}
