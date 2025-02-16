package types

type OrganizationPayload struct {
	ID        string `json:"id"`
	Subdomain string `json:"subdomain"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateOrganizationInput struct {
	Subdomain string `json:"subdomain"`
}

type CreateOrganizationPayload struct {
	Organization *OrganizationPayload `json:"organization"`
}

type CheckSubdomainAvailablityInput struct {
	Subdomain string `validate:"required"`
}

type UpdateOrganizationUserInput struct {
	UserID   string   `json:"-" validate:"required"`
	Role     *string  `json:"role" validate:"oneof=admin developer member"`
	GroupIDs []string `json:"groupIds"`
}
