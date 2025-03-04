package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
)

// OrganizationDTOToResponse converts from dto.Organization to responses.OrganizationResponse
func OrganizationDTOToResponse(org *dto.Organization) *responses.OrganizationResponse {
	if org == nil {
		return nil
	}

	return &responses.OrganizationResponse{
		ID:        org.ID,
		Subdomain: org.Subdomain,
		CreatedAt: strconv.FormatInt(org.CreatedAt, 10),
		UpdatedAt: strconv.FormatInt(org.UpdatedAt, 10),
	}
}

// CreateOrganizationRequestToDTOInput converts from requests.CreateOrganizationRequest to dto.CreateOrganizationInput
func CreateOrganizationRequestToDTOInput(in requests.CreateOrganizationRequest) dto.CreateOrganizationInput {
	return dto.CreateOrganizationInput{
		Subdomain: in.Subdomain,
	}
}

// CreateOrganizationOutputToResponse converts from dto.CreateOrganizationOutput to responses.CreateOrganizationResponse
func CreateOrganizationOutputToResponse(out *dto.CreateOrganizationOutput) *responses.CreateOrganizationResponse {
	return &responses.CreateOrganizationResponse{
		Organization: OrganizationDTOToResponse(out.Organization),
	}
}

// CheckSubdomainAvailabilityRequestToDTOInput converts from requests.CheckSubdomainAvailablityRequest to dto.CheckSubdomainAvailabilityInput
func CheckSubdomainAvailabilityRequestToDTOInput(in requests.CheckSubdomainAvailablityRequest) dto.CheckSubdomainAvailabilityInput {
	return dto.CheckSubdomainAvailabilityInput{
		Subdomain: in.Subdomain,
	}
}

// UpdateOrganizationUserRequestToDTOInput converts from requests.UpdateOrganizationUserRequest to dto.UpdateOrganizationUserInput
func UpdateOrganizationUserRequestToDTOInput(in requests.UpdateOrganizationUserRequest) dto.UpdateOrganizationUserInput {
	return dto.UpdateOrganizationUserInput{
		UserID:   in.UserID,
		Role:     in.Role,
		GroupIDs: in.GroupIDs,
	}
}

// UpdateOrganizationUserOutputToResponse converts from dto.UpdateOrganizationUserOutput to responses.UpdateUserResponse
func UpdateOrganizationUserOutputToResponse(out *dto.UpdateOrganizationUserOutput) *responses.UpdateUserResponse {
	return &responses.UpdateUserResponse{
		User: UserDTOToResponse(out.User),
	}
}
