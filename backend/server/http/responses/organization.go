package responses

type OrganizationResponse struct {
	ID        string `json:"id"`
	Subdomain string `json:"subdomain"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateOrganizationResponse struct {
	Organization *OrganizationResponse `json:"organization"`
}
