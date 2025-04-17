package adapters

import (
	"github.com/trysourcetool/sourcetool/backend/dto/http/requests"
	"github.com/trysourcetool/sourcetool/backend/dto/http/responses"
	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
)

// RequestMagicLinkRequestToInput converts from requests.RequestMagicLinkRequest to input.RequestMagicLinkInput.
func RequestMagicLinkRequestToInput(in requests.RequestMagicLinkRequest) input.RequestMagicLinkInput {
	return input.RequestMagicLinkInput{
		Email: in.Email,
	}
}

// AuthenticateWithMagicLinkRequestToInput converts from requests.AuthenticateWithMagicLinkRequest to input.AuthenticateWithMagicLinkInput.
func AuthenticateWithMagicLinkRequestToInput(in requests.AuthenticateWithMagicLinkRequest) input.AuthenticateWithMagicLinkInput {
	return input.AuthenticateWithMagicLinkInput{
		Token:     in.Token,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

// AuthenticateWithMagicLinkOutputToResponse converts from output.AuthenticateWithMagicLinkOutput to responses.AuthenticateWithMagicLinkResponse.
func AuthenticateWithMagicLinkOutputToResponse(out *output.AuthenticateWithMagicLinkOutput) *responses.AuthenticateWithMagicLinkResponse {
	return &responses.AuthenticateWithMagicLinkResponse{
		AuthURL:         out.AuthURL,
		Token:           out.Token,
		HasOrganization: out.HasOrganization,
		IsNewUser:       out.IsNewUser,
	}
}

// RequestMagicLinkOutputToResponse converts from output.RequestMagicLinkOutput to responses.RequestMagicLinkResponse.
func RequestMagicLinkOutputToResponse(out *output.RequestMagicLinkOutput) responses.RequestMagicLinkResponse {
	return responses.RequestMagicLinkResponse{
		Email: out.Email,
		IsNew: out.IsNew,
	}
}

// RegisterWithMagicLinkRequestToInput converts from requests.RegisterWithMagicLinkRequest to input.RegisterWithMagicLinkInput.
func RegisterWithMagicLinkRequestToInput(in requests.RegisterWithMagicLinkRequest) input.RegisterWithMagicLinkInput {
	return input.RegisterWithMagicLinkInput{
		Token:     in.Token,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

// RegisterWithMagicLinkOutputToResponse converts from output.RegisterWithMagicLinkOutput to responses.RegisterWithMagicLinkResponse.
func RegisterWithMagicLinkOutputToResponse(out *output.RegisterWithMagicLinkOutput) *responses.RegisterWithMagicLinkResponse {
	return &responses.RegisterWithMagicLinkResponse{
		ExpiresAt:       out.ExpiresAt,
		HasOrganization: out.HasOrganization,
	}
}

// RequestInvitationMagicLinkRequestToInput converts from requests.RequestInvitationMagicLinkRequest to input.RequestInvitationMagicLinkInput.
func RequestInvitationMagicLinkRequestToInput(in requests.RequestInvitationMagicLinkRequest) input.RequestInvitationMagicLinkInput {
	return input.RequestInvitationMagicLinkInput{
		InvitationToken: in.InvitationToken,
	}
}

// RequestInvitationMagicLinkOutputToResponse converts from output.RequestInvitationMagicLinkOutput to responses.RequestInvitationMagicLinkResponse.
func RequestInvitationMagicLinkOutputToResponse(out *output.RequestInvitationMagicLinkOutput) *responses.RequestInvitationMagicLinkResponse {
	return &responses.RequestInvitationMagicLinkResponse{
		Email: out.Email,
	}
}

// AuthenticateWithInvitationMagicLinkRequestToInput converts from requests.AuthenticateWithInvitationMagicLinkRequest to input.AuthenticateWithInvitationMagicLinkInput.
func AuthenticateWithInvitationMagicLinkRequestToInput(in requests.AuthenticateWithInvitationMagicLinkRequest) input.AuthenticateWithInvitationMagicLinkInput {
	return input.AuthenticateWithInvitationMagicLinkInput{
		Token: in.Token,
	}
}

// AuthenticateWithInvitationMagicLinkOutputToResponse converts from output.AuthenticateWithInvitationMagicLinkOutput to responses.AuthenticateWithInvitationMagicLinkResponse.
func AuthenticateWithInvitationMagicLinkOutputToResponse(out *output.AuthenticateWithInvitationMagicLinkOutput) *responses.AuthenticateWithInvitationMagicLinkResponse {
	return &responses.AuthenticateWithInvitationMagicLinkResponse{
		AuthURL:   out.AuthURL,
		Token:     out.Token,
		IsNewUser: out.IsNewUser,
	}
}

// RegisterWithInvitationMagicLinkRequestToInput converts from requests.RegisterWithInvitationMagicLinkRequest to input.RegisterWithInvitationMagicLinkInput.
func RegisterWithInvitationMagicLinkRequestToInput(in requests.RegisterWithInvitationMagicLinkRequest) input.RegisterWithInvitationMagicLinkInput {
	return input.RegisterWithInvitationMagicLinkInput{
		Token:     in.Token,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

// RegisterWithInvitationMagicLinkOutputToResponse converts from output.RegisterWithInvitationMagicLinkOutput to responses.RegisterWithInvitationMagicLinkResponse.
func RegisterWithInvitationMagicLinkOutputToResponse(out *output.RegisterWithInvitationMagicLinkOutput) *responses.RegisterWithInvitationMagicLinkResponse {
	return &responses.RegisterWithInvitationMagicLinkResponse{
		ExpiresAt: out.ExpiresAt,
	}
}

// RequestGoogleAuthLinkOutputToResponse converts from output.RequestGoogleAuthLinkOutput to responses.RequestGoogleAuthLinkResponse.
func RequestGoogleAuthLinkOutputToResponse(out *output.RequestGoogleAuthLinkOutput) *responses.RequestGoogleAuthLinkResponse {
	return &responses.RequestGoogleAuthLinkResponse{
		AuthURL: out.AuthURL,
	}
}

// AuthenticateWithGoogleRequestToInput converts requests.AuthenticateWithGoogleRequest to input.AuthenticateWithGoogleInput.
func AuthenticateWithGoogleRequestToInput(req requests.AuthenticateWithGoogleRequest) input.AuthenticateWithGoogleInput {
	return input.AuthenticateWithGoogleInput{
		Code:  req.Code,
		State: req.State,
	}
}

// AuthenticateWithGoogleOutputToResponse converts output.AuthenticateWithGoogleOutput to responses.AuthenticateWithGoogleResponse.
func AuthenticateWithGoogleOutputToResponse(out *output.AuthenticateWithGoogleOutput) *responses.AuthenticateWithGoogleResponse {
	return &responses.AuthenticateWithGoogleResponse{
		FirstName:                out.FirstName,
		LastName:                 out.LastName,
		AuthURL:                  out.AuthURL,
		Token:                    out.Token,
		HasOrganization:          out.HasOrganization,
		HasMultipleOrganizations: out.HasMultipleOrganizations,
		IsNewUser:                out.IsNewUser,
	}
}

// RegisterWithGoogleRequestToInput converts from requests.RegisterWithGoogleRequest to input.RegisterWithGoogleInput.
func RegisterWithGoogleRequestToInput(in requests.RegisterWithGoogleRequest) input.RegisterWithGoogleInput {
	return input.RegisterWithGoogleInput{
		Token: in.Token,
	}
}

// RegisterWithGoogleOutputToResponse converts from output.RegisterWithGoogleOutput to responses.RegisterWithGoogleResponse.
func RegisterWithGoogleOutputToResponse(out *output.RegisterWithGoogleOutput) *responses.RegisterWithGoogleResponse {
	return &responses.RegisterWithGoogleResponse{
		AuthURL:         out.AuthURL,
		Token:           out.Token,
		HasOrganization: out.HasOrganization,
	}
}

// RequestInvitationGoogleAuthLinkRequestToInput converts request to input.
func RequestInvitationGoogleAuthLinkRequestToInput(in requests.RequestInvitationGoogleAuthLinkRequest) input.RequestInvitationGoogleAuthLinkInput {
	return input.RequestInvitationGoogleAuthLinkInput{
		InvitationToken: in.InvitationToken,
	}
}

// RequestInvitationGoogleAuthLinkOutputToResponse converts output.RequestInvitationGoogleAuthLinkOutput to responses.RequestInvitationGoogleAuthLinkResponse.
func RequestInvitationGoogleAuthLinkOutputToResponse(out *output.RequestInvitationGoogleAuthLinkOutput) *responses.RequestInvitationGoogleAuthLinkResponse {
	return &responses.RequestInvitationGoogleAuthLinkResponse{
		AuthURL: out.AuthURL,
	}
}

// RefreshTokenRequestToInput converts from requests.RefreshTokenRequest to input.RefreshTokenInput.
func RefreshTokenRequestToInput(in requests.RefreshTokenRequest) input.RefreshTokenInput {
	return input.RefreshTokenInput{
		RefreshToken:    in.RefreshToken,
		XSRFTokenHeader: in.XSRFTokenHeader,
		XSRFTokenCookie: in.XSRFTokenCookie,
	}
}

// RefreshTokenOutputToResponse converts from output.RefreshTokenOutput to responses.RefreshTokenResponse.
func RefreshTokenOutputToResponse(out *output.RefreshTokenOutput) *responses.RefreshTokenResponse {
	return &responses.RefreshTokenResponse{
		ExpiresAt: out.ExpiresAt,
	}
}

// SaveAuthRequestToInput converts from requests.SaveAuthRequest to input.SaveAuthInput.
func SaveAuthRequestToInput(in requests.SaveAuthRequest) input.SaveAuthInput {
	return input.SaveAuthInput{
		Token: in.Token,
	}
}

// SaveAuthOutputToResponse converts from output.SaveAuthOutput to responses.SaveAuthResponse.
func SaveAuthOutputToResponse(out *output.SaveAuthOutput) *responses.SaveAuthResponse {
	return &responses.SaveAuthResponse{
		ExpiresAt:   out.ExpiresAt,
		RedirectURL: out.RedirectURL,
	}
}

// ObtainAuthTokenOutputToResponse converts from output.ObtainAuthTokenOutput to responses.ObtainAuthTokenResponse.
func ObtainAuthTokenOutputToResponse(out *output.ObtainAuthTokenOutput) *responses.ObtainAuthTokenResponse {
	return &responses.ObtainAuthTokenResponse{
		AuthURL: out.AuthURL,
		Token:   out.Token,
	}
}

// LogoutOutputToResponse converts from output.LogoutOutput to responses.LogoutResponse.
func LogoutOutputToResponse(out *output.LogoutOutput) *responses.LogoutResponse {
	return &responses.LogoutResponse{
		Domain: out.Domain,
	}
}
