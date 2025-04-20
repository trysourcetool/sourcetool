package mapper

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/responses"
)

func OrganizationOutputToResponse(org *dto.Organization) *responses.OrganizationResponse {
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

func CreateOrganizationRequestToInput(in requests.CreateOrganizationRequest) dto.CreateOrganizationInput {
	return dto.CreateOrganizationInput{
		Subdomain: in.Subdomain,
	}
}

func CreateOrganizationOutputToResponse(out *dto.CreateOrganizationOutput) *responses.CreateOrganizationResponse {
	return &responses.CreateOrganizationResponse{
		Organization: OrganizationOutputToResponse(out.Organization),
	}
}

func CheckSubdomainAvailabilityRequestToInput(in requests.CheckSubdomainAvailablityRequest) dto.CheckSubdomainAvailabilityInput {
	return dto.CheckSubdomainAvailabilityInput{
		Subdomain: in.Subdomain,
	}
}
