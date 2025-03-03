package requests

type CreateOrganizationRequest struct {
	Subdomain string `json:"subdomain"`
}

type CheckSubdomainAvailablityRequest struct {
	Subdomain string `validate:"required"`
}

type UpdateOrganizationUserRequest struct {
	UserID   string   `json:"-" validate:"required"`
	Role     *string  `json:"role" validate:"oneof=admin developer member"`
	GroupIDs []string `json:"groupIds"`
}
