package requests

type CreateOrganizationRequest struct {
	Subdomain string `json:"subdomain" validate:"required"`
}
