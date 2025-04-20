package dto

import (
	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
)

type CreateOrganizationInput struct {
	Subdomain string
}

type CheckSubdomainAvailabilityInput struct {
	Subdomain string
}

type Organization struct {
	ID        string
	Subdomain string
	CreatedAt int64
	UpdatedAt int64
}

func OrganizationFromModel(org *organization.Organization) *Organization {
	if org == nil {
		return nil
	}

	return &Organization{
		ID:        org.ID.String(),
		Subdomain: internal.SafeValue(org.Subdomain),
		CreatedAt: org.CreatedAt.Unix(),
		UpdatedAt: org.UpdatedAt.Unix(),
	}
}

type CreateOrganizationOutput struct {
	Organization *Organization
}
