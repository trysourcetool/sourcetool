package dto

import (
	"github.com/trysourcetool/sourcetool/backend/model"
)

// Organization represents organization data in DTOs
type Organization struct {
	ID        string
	Subdomain string
	CreatedAt int64
	UpdatedAt int64
}

// OrganizationFromModel converts from model.Organization to dto.Organization
func OrganizationFromModel(org *model.Organization) *Organization {
	if org == nil {
		return nil
	}

	return &Organization{
		ID:        org.ID.String(),
		Subdomain: org.Subdomain,
		CreatedAt: org.CreatedAt.Unix(),
		UpdatedAt: org.UpdatedAt.Unix(),
	}
}

// CreateOrganizationInput is the input for Create operation
type CreateOrganizationInput struct {
	Subdomain string
}

// CreateOrganizationOutput is the output for Create operation
type CreateOrganizationOutput struct {
	Organization *Organization
}

// CheckSubdomainAvailabilityInput is the input for checking subdomain availability
type CheckSubdomainAvailabilityInput struct {
	Subdomain string
}

// UpdateOrganizationUserInput is the input for updating an organization user
type UpdateOrganizationUserInput struct {
	UserID   string
	Role     *string
	GroupIDs []string
}

// UpdateOrganizationUserOutput is the output for updating an organization user
type UpdateOrganizationUserOutput struct {
	User *User // This requires dto.User to be implemented
}
