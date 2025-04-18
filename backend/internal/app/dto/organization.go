package dto

import (
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	"github.com/trysourcetool/sourcetool/backend/pkg/ptrconv"
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
		Subdomain: ptrconv.SafeValue(org.Subdomain),
		CreatedAt: org.CreatedAt.Unix(),
		UpdatedAt: org.UpdatedAt.Unix(),
	}
}

type CreateOrganizationOutput struct {
	Organization *Organization
}
