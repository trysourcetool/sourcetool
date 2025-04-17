package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/dto/http/requests"
	"github.com/trysourcetool/sourcetool/backend/dto/http/responses"
	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
)

// OrganizationOutputToResponse converts from output.Organization to responses.OrganizationResponse.
func OrganizationOutputToResponse(org *output.Organization) *responses.OrganizationResponse {
	if org == nil {
		return nil
	}

	return &responses.OrganizationResponse{
		ID:                org.ID,
		Subdomain:         org.Subdomain,
		WebSocketEndpoint: config.Config.WebSocketOrgBaseURL(org.Subdomain),
		CreatedAt:         strconv.FormatInt(org.CreatedAt, 10),
		UpdatedAt:         strconv.FormatInt(org.UpdatedAt, 10),
	}
}

// CreateOrganizationRequestToInput converts from requests.CreateOrganizationRequest to input.CreateOrganizationInput.
func CreateOrganizationRequestToInput(in requests.CreateOrganizationRequest) input.CreateOrganizationInput {
	return input.CreateOrganizationInput{
		Subdomain: in.Subdomain,
	}
}

// CreateOrganizationOutputToResponse converts from output.CreateOrganizationOutput to responses.CreateOrganizationResponse.
func CreateOrganizationOutputToResponse(out *output.CreateOrganizationOutput) *responses.CreateOrganizationResponse {
	return &responses.CreateOrganizationResponse{
		Organization: OrganizationOutputToResponse(out.Organization),
	}
}

// CheckSubdomainAvailabilityRequestToInput converts from requests.CheckSubdomainAvailablityRequest to input.CheckSubdomainAvailabilityInput.
func CheckSubdomainAvailabilityRequestToInput(in requests.CheckSubdomainAvailablityRequest) input.CheckSubdomainAvailabilityInput {
	return input.CheckSubdomainAvailabilityInput{
		Subdomain: in.Subdomain,
	}
}
