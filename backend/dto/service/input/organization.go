package input

// CreateOrganizationInput is the input for Create operation.
type CreateOrganizationInput struct {
	Subdomain string
}

// CheckSubdomainAvailabilityInput is the input for checking subdomain availability.
type CheckSubdomainAvailabilityInput struct {
	Subdomain string
}
