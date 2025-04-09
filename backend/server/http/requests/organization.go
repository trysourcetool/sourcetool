package requests

type CreateOrganizationRequest struct {
	Subdomain string `json:"subdomain" validate:"required"`
}

type CheckSubdomainAvailablityRequest struct {
	Subdomain string `validate:"required"`
}

type UpdateOrganizationUserRequest struct {
	UserID   string   `json:"-" validate:"required,uuid4"`
	Role     *string  `json:"role" validate:"oneof=admin developer member"`
	GroupIDs []string `json:"groupIds"`
}

type DeleteOrganizationUserRequest struct {
	UserID string `param:"userID" validate:"required,uuid4"`
}
