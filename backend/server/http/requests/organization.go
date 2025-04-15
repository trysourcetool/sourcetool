package requests

type CreateOrganizationRequest struct {
	Subdomain string `json:"subdomain" validate:"required"`
}

type CheckSubdomainAvailablityRequest struct {
	Subdomain string `validate:"required"`
}
