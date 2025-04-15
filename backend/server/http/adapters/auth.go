package adapters

import (
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
)

// RequestMagicLinkRequestToDTOInput converts from requests.RequestMagicLinkRequest to dto.RequestMagicLinkInput.
func RequestMagicLinkRequestToDTOInput(in requests.RequestMagicLinkRequest) dto.RequestMagicLinkInput {
	return dto.RequestMagicLinkInput{
		Email: in.Email,
	}
}

// AuthenticateWithMagicLinkRequestToDTOInput converts from requests.AuthenticateWithMagicLinkRequest to dto.AuthenticateWithMagicLinkInput.
func AuthenticateWithMagicLinkRequestToDTOInput(in requests.AuthenticateWithMagicLinkRequest) dto.AuthenticateWithMagicLinkInput {
	return dto.AuthenticateWithMagicLinkInput{
		Token:     in.Token,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

// AuthenticateWithMagicLinkOutputToResponse converts from dto.AuthenticateWithMagicLinkOutput to responses.AuthenticateWithMagicLinkResponse.
func AuthenticateWithMagicLinkOutputToResponse(out *dto.AuthenticateWithMagicLinkOutput) *responses.AuthenticateWithMagicLinkResponse {
	return &responses.AuthenticateWithMagicLinkResponse{
		AuthURL:         out.AuthURL,
		Token:           out.Token,
		HasOrganization: out.HasOrganization,
		IsNewUser:       out.IsNewUser,
	}
}

// RequestMagicLinkOutputToResponse converts from dto.RequestMagicLinkOutput to responses.RequestMagicLinkResponse.
func RequestMagicLinkOutputToResponse(out *dto.RequestMagicLinkOutput) responses.RequestMagicLinkResponse {
	return responses.RequestMagicLinkResponse{
		Email: out.Email,
		IsNew: out.IsNew,
	}
}

// RegisterWithMagicLinkRequestToDTOInput converts from requests.RegisterWithMagicLinkRequest to dto.RegisterWithMagicLinkInput.
func RegisterWithMagicLinkRequestToDTOInput(in requests.RegisterWithMagicLinkRequest) dto.RegisterWithMagicLinkInput {
	return dto.RegisterWithMagicLinkInput{
		Token:     in.Token,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

// RegisterWithMagicLinkOutputToResponse converts from dto.RegisterWithMagicLinkOutput to responses.RegisterWithMagicLinkResponse.
func RegisterWithMagicLinkOutputToResponse(out *dto.RegisterWithMagicLinkOutput) *responses.RegisterWithMagicLinkResponse {
	return &responses.RegisterWithMagicLinkResponse{
		ExpiresAt:       out.ExpiresAt,
		HasOrganization: out.HasOrganization,
	}
}

// RequestInvitationMagicLinkRequestToDTOInput converts from requests.RequestInvitationMagicLinkRequest to dto.RequestInvitationMagicLinkInput.
func RequestInvitationMagicLinkRequestToDTOInput(in requests.RequestInvitationMagicLinkRequest) dto.RequestInvitationMagicLinkInput {
	return dto.RequestInvitationMagicLinkInput{
		InvitationToken: in.InvitationToken,
	}
}

// RequestInvitationMagicLinkOutputToResponse converts from dto.RequestInvitationMagicLinkOutput to responses.RequestInvitationMagicLinkResponse.
func RequestInvitationMagicLinkOutputToResponse(out *dto.RequestInvitationMagicLinkOutput) *responses.RequestInvitationMagicLinkResponse {
	return &responses.RequestInvitationMagicLinkResponse{
		Email: out.Email,
	}
}

// AuthenticateWithInvitationMagicLinkRequestToDTOInput converts from requests.AuthenticateWithInvitationMagicLinkRequest to dto.AuthenticateWithInvitationMagicLinkInput.
func AuthenticateWithInvitationMagicLinkRequestToDTOInput(in requests.AuthenticateWithInvitationMagicLinkRequest) dto.AuthenticateWithInvitationMagicLinkInput {
	return dto.AuthenticateWithInvitationMagicLinkInput{
		Token: in.Token,
	}
}

// AuthenticateWithInvitationMagicLinkOutputToResponse converts from dto.AuthenticateWithInvitationMagicLinkOutput to responses.AuthenticateWithInvitationMagicLinkResponse.
func AuthenticateWithInvitationMagicLinkOutputToResponse(out *dto.AuthenticateWithInvitationMagicLinkOutput) *responses.AuthenticateWithInvitationMagicLinkResponse {
	return &responses.AuthenticateWithInvitationMagicLinkResponse{
		AuthURL:   out.AuthURL,
		Token:     out.Token,
		IsNewUser: out.IsNewUser,
	}
}

// RegisterWithInvitationMagicLinkRequestToDTOInput converts from requests.RegisterWithInvitationMagicLinkRequest to dto.RegisterWithInvitationMagicLinkInput.
func RegisterWithInvitationMagicLinkRequestToDTOInput(in requests.RegisterWithInvitationMagicLinkRequest) dto.RegisterWithInvitationMagicLinkInput {
	return dto.RegisterWithInvitationMagicLinkInput{
		Token:     in.Token,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

// RegisterWithInvitationMagicLinkOutputToResponse converts from dto.RegisterWithInvitationMagicLinkOutput to responses.RegisterWithInvitationMagicLinkResponse.
func RegisterWithInvitationMagicLinkOutputToResponse(out *dto.RegisterWithInvitationMagicLinkOutput) *responses.RegisterWithInvitationMagicLinkResponse {
	return &responses.RegisterWithInvitationMagicLinkResponse{
		ExpiresAt: out.ExpiresAt,
	}
}

// RequestGoogleAuthLinkOutputToResponse converts from dto.RequestGoogleAuthLinkOutput to responses.RequestGoogleAuthLinkResponse.
func RequestGoogleAuthLinkOutputToResponse(out *dto.RequestGoogleAuthLinkOutput) *responses.RequestGoogleAuthLinkResponse {
	return &responses.RequestGoogleAuthLinkResponse{
		AuthURL: out.AuthURL,
	}
}

// AuthenticateWithGoogleRequestToDTOInput converts requests.AuthenticateWithGoogleRequest to dto.AuthenticateWithGoogleInput.
func AuthenticateWithGoogleRequestToDTOInput(req requests.AuthenticateWithGoogleRequest) dto.AuthenticateWithGoogleInput {
	return dto.AuthenticateWithGoogleInput{
		Code:  req.Code,
		State: req.State,
	}
}

// AuthenticateWithGoogleOutputToResponse converts dto.AuthenticateWithGoogleOutput to responses.AuthenticateWithGoogleResponse.
func AuthenticateWithGoogleOutputToResponse(out *dto.AuthenticateWithGoogleOutput) *responses.AuthenticateWithGoogleResponse {
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

// RegisterWithGoogleRequestToDTOInput converts from requests.RegisterWithGoogleRequest to dto.RegisterWithGoogleInput.
func RegisterWithGoogleRequestToDTOInput(in requests.RegisterWithGoogleRequest) dto.RegisterWithGoogleInput {
	return dto.RegisterWithGoogleInput{
		Token: in.Token,
	}
}

// RegisterWithGoogleOutputToResponse converts from dto.RegisterWithGoogleOutput to responses.RegisterWithGoogleResponse.
func RegisterWithGoogleOutputToResponse(out *dto.RegisterWithGoogleOutput) *responses.RegisterWithGoogleResponse {
	return &responses.RegisterWithGoogleResponse{
		AuthURL:         out.AuthURL,
		Token:           out.Token,
		HasOrganization: out.HasOrganization,
	}
}

// RequestInvitationGoogleAuthLinkRequestToDTOInput converts request to DTO input.
func RequestInvitationGoogleAuthLinkRequestToDTOInput(in requests.RequestInvitationGoogleAuthLinkRequest) dto.RequestInvitationGoogleAuthLinkInput {
	return dto.RequestInvitationGoogleAuthLinkInput{
		InvitationToken: in.InvitationToken,
	}
}

// RequestInvitationGoogleAuthLinkOutputToResponse converts DTO output to response.
func RequestInvitationGoogleAuthLinkOutputToResponse(out *dto.RequestInvitationGoogleAuthLinkOutput) *responses.RequestInvitationGoogleAuthLinkResponse {
	return &responses.RequestInvitationGoogleAuthLinkResponse{
		AuthURL: out.AuthURL,
	}
}

// RefreshTokenRequestToDTOInput converts from requests.RefreshTokenRequest to dto.RefreshTokenInput.
func RefreshTokenRequestToDTOInput(in requests.RefreshTokenRequest) dto.RefreshTokenInput {
	return dto.RefreshTokenInput{
		RefreshToken:    in.RefreshToken,
		XSRFTokenHeader: in.XSRFTokenHeader,
		XSRFTokenCookie: in.XSRFTokenCookie,
	}
}

// RefreshTokenOutputToResponse converts from dto.RefreshTokenOutput to responses.RefreshTokenResponse.
func RefreshTokenOutputToResponse(out *dto.RefreshTokenOutput) *responses.RefreshTokenResponse {
	return &responses.RefreshTokenResponse{
		ExpiresAt: out.ExpiresAt,
	}
}

// SaveAuthRequestToDTOInput converts from requests.SaveAuthRequest to dto.SaveAuthInput.
func SaveAuthRequestToDTOInput(in requests.SaveAuthRequest) dto.SaveAuthInput {
	return dto.SaveAuthInput{
		Token: in.Token,
	}
}

// SaveAuthOutputToResponse converts from dto.SaveAuthOutput to responses.SaveAuthResponse.
func SaveAuthOutputToResponse(out *dto.SaveAuthOutput) *responses.SaveAuthResponse {
	return &responses.SaveAuthResponse{
		ExpiresAt:   out.ExpiresAt,
		RedirectURL: out.RedirectURL,
	}
}

// ObtainAuthTokenOutputToResponse converts from dto.ObtainAuthTokenOutput to responses.ObtainAuthTokenResponse.
func ObtainAuthTokenOutputToResponse(out *dto.ObtainAuthTokenOutput) *responses.ObtainAuthTokenResponse {
	return &responses.ObtainAuthTokenResponse{
		AuthURL: out.AuthURL,
		Token:   out.Token,
	}
}

// LogoutOutputToResponse converts from dto.LogoutOutput to responses.LogoutResponse.
func LogoutOutputToResponse(out *dto.LogoutOutput) *responses.LogoutResponse {
	return &responses.LogoutResponse{
		Domain: out.Domain,
	}
}
