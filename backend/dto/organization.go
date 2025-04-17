package dto

import (
	"github.com/trysourcetool/sourcetool/backend/organization"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
)

// Organization represents organization data in DTOs.
type Organization struct {
	ID        string
	Subdomain string
	CreatedAt int64
	UpdatedAt int64
}

// OrganizationFromModel converts from model.Organization to dto.Organization.
func OrganizationFromModel(org *organization.Organization) *Organization {
	if org == nil {
		return nil
	}

	return &Organization{
		ID:        org.ID.String(),
		Subdomain: conv.SafeValue(org.Subdomain),
		CreatedAt: org.CreatedAt.Unix(),
		UpdatedAt: org.UpdatedAt.Unix(),
	}
}

// CreateOrganizationInput is the input for Create operation.
type CreateOrganizationInput struct {
	Subdomain string
}

// CreateOrganizationOutput is the output for Create operation.
type CreateOrganizationOutput struct {
	Organization *Organization
}

// CheckSubdomainAvailabilityInput is the input for checking subdomain availability.
type CheckSubdomainAvailabilityInput struct {
	Subdomain string
}
