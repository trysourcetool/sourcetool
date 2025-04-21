package requests

type GetGroupRequest struct {
	GroupID string `json:"-" validate:"required"`
}

type CreateGroupRequest struct {
	Name    string   `json:"name" validate:"required"`
	Slug    string   `json:"slug" validate:"required"`
	UserIDs []string `json:"userIds" validate:"required"`
}

type UpdateGroupRequest struct {
	Name    *string  `json:"name" validate:"required"`
	UserIDs []string `json:"userIds" validate:"required"`
}
