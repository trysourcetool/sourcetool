package server

type groupResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type groupPageResponse struct {
	ID        string `json:"id"`
	GroupID   string `json:"groupId"`
	PageID    string `json:"pageId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
