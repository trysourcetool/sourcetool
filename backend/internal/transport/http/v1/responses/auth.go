package responses

type RequestMagicLinkResponse struct {
	Email string `json:"email"`
	IsNew bool   `json:"isNew"`
}

type AuthenticateWithMagicLinkResponse struct {
	AuthURL         string `json:"authUrl"`
	Token           string `json:"token"`
	HasOrganization bool   `json:"hasOrganization"`
	IsNewUser       bool   `json:"isNewUser"`
}

type RegisterWithMagicLinkResponse struct {
	HasOrganization bool   `json:"hasOrganization"`
	ExpiresAt       string `json:"expiresAt"`
}

type RequestInvitationMagicLinkResponse struct {
	Email string `json:"email"`
}

type AuthenticateWithInvitationMagicLinkResponse struct {
	AuthURL   string `json:"authUrl"`
	Token     string `json:"token"`
	IsNewUser bool   `json:"isNewUser"`
}

type RegisterWithInvitationMagicLinkResponse struct {
	ExpiresAt string `json:"expiresAt"`
}

type RequestGoogleAuthLinkResponse struct {
	AuthURL string `json:"authUrl"`
}

type AuthenticateWithGoogleResponse struct {
	FirstName                string `json:"firstName,omitempty"`
	LastName                 string `json:"lastName,omitempty"`
	AuthURL                  string `json:"authUrl"`
	Token                    string `json:"token"`
	HasOrganization          bool   `json:"hasOrganization"`
	HasMultipleOrganizations bool   `json:"hasMultipleOrganizations"`
	IsNewUser                bool   `json:"isNewUser"`
}

type RegisterWithGoogleResponse struct {
	AuthURL         string `json:"authUrl"`
	Token           string `json:"token"`
	HasOrganization bool   `json:"hasOrganization"`
}

type RequestInvitationGoogleAuthLinkResponse struct {
	AuthURL string `json:"authUrl"`
}

type RefreshTokenResponse struct {
	ExpiresAt string `json:"expiresAt"`
}

type SaveAuthResponse struct {
	ExpiresAt   string `json:"expiresAt"`
	RedirectURL string `json:"redirectUrl"`
}

type ObtainAuthTokenResponse struct {
	AuthURL string `json:"authUrl"`
	Token   string `json:"token"`
}

type LogoutResponse struct {
	Domain string
}
