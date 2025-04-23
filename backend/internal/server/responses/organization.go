package responses

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type OrganizationResponse struct {
	ID                string `json:"id"`
	Subdomain         string `json:"subdomain"`
	WebSocketEndpoint string `json:"webSocketEndpoint"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

func OrganizationFromModel(o *core.Organization) *OrganizationResponse {
	if o == nil {
		return nil
	}

	return &OrganizationResponse{
		ID:                o.ID.String(),
		Subdomain:         internal.StringValue(o.Subdomain),
		WebSocketEndpoint: config.Config.WebSocketOrgBaseURL(internal.StringValue(o.Subdomain)),
		CreatedAt:         strconv.FormatInt(o.CreatedAt.Unix(), 10),
		UpdatedAt:         strconv.FormatInt(o.UpdatedAt.Unix(), 10),
	}
}

type CreateOrganizationResponse struct {
	Organization *OrganizationResponse `json:"organization"`
}
