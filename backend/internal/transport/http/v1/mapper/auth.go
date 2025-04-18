package mapper

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/responses"
)

func RequestMagicLinkRequestToInput(in requests.RequestMagicLinkRequest) dto.RequestMagicLinkInput {
	return dto.RequestMagicLinkInput{
		Email: in.Email,
	}
}

func AuthenticateWithMagicLinkRequestToInput(in requests.AuthenticateWithMagicLinkRequest) dto.AuthenticateWithMagicLinkInput {
	return dto.AuthenticateWithMagicLinkInput{
		Token:     in.Token,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

func AuthenticateWithMagicLinkOutputToResponse(out *dto.AuthenticateWithMagicLinkOutput) *responses.AuthenticateWithMagicLinkResponse {
	return &responses.AuthenticateWithMagicLinkResponse{
		AuthURL:         out.AuthURL,
		Token:           out.Token,
		HasOrganization: out.HasOrganization,
		IsNewUser:       out.IsNewUser,
	}
}

func RequestMagicLinkOutputToResponse(out *dto.RequestMagicLinkOutput) responses.RequestMagicLinkResponse {
	return responses.RequestMagicLinkResponse{
		Email: out.Email,
		IsNew: out.IsNew,
	}
}

func RegisterWithMagicLinkRequestToInput(in requests.RegisterWithMagicLinkRequest) dto.RegisterWithMagicLinkInput {
	return dto.RegisterWithMagicLinkInput{
		Token:     in.Token,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

func RegisterWithMagicLinkOutputToResponse(out *dto.RegisterWithMagicLinkOutput) *responses.RegisterWithMagicLinkResponse {
	return &responses.RegisterWithMagicLinkResponse{
		ExpiresAt:       out.ExpiresAt,
		HasOrganization: out.HasOrganization,
	}
}

func RequestInvitationMagicLinkRequestToInput(in requests.RequestInvitationMagicLinkRequest) dto.RequestInvitationMagicLinkInput {
	return dto.RequestInvitationMagicLinkInput{
		InvitationToken: in.InvitationToken,
	}
}

func RequestInvitationMagicLinkOutputToResponse(out *dto.RequestInvitationMagicLinkOutput) *responses.RequestInvitationMagicLinkResponse {
	return &responses.RequestInvitationMagicLinkResponse{
		Email: out.Email,
	}
}

func AuthenticateWithInvitationMagicLinkRequestToInput(in requests.AuthenticateWithInvitationMagicLinkRequest) dto.AuthenticateWithInvitationMagicLinkInput {
	return dto.AuthenticateWithInvitationMagicLinkInput{
		Token: in.Token,
	}
}

func AuthenticateWithInvitationMagicLinkOutputToResponse(out *dto.AuthenticateWithInvitationMagicLinkOutput) *responses.AuthenticateWithInvitationMagicLinkResponse {
	return &responses.AuthenticateWithInvitationMagicLinkResponse{
		AuthURL:   out.AuthURL,
		Token:     out.Token,
		IsNewUser: out.IsNewUser,
	}
}

func RegisterWithInvitationMagicLinkRequestToInput(in requests.RegisterWithInvitationMagicLinkRequest) dto.RegisterWithInvitationMagicLinkInput {
	return dto.RegisterWithInvitationMagicLinkInput{
		Token:     in.Token,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

func RegisterWithInvitationMagicLinkOutputToResponse(out *dto.RegisterWithInvitationMagicLinkOutput) *responses.RegisterWithInvitationMagicLinkResponse {
	return &responses.RegisterWithInvitationMagicLinkResponse{
		ExpiresAt: out.ExpiresAt,
	}
}

func RequestGoogleAuthLinkOutputToResponse(out *dto.RequestGoogleAuthLinkOutput) *responses.RequestGoogleAuthLinkResponse {
	return &responses.RequestGoogleAuthLinkResponse{
		AuthURL: out.AuthURL,
	}
}

func AuthenticateWithGoogleRequestToInput(req requests.AuthenticateWithGoogleRequest) dto.AuthenticateWithGoogleInput {
	return dto.AuthenticateWithGoogleInput{
		Code:  req.Code,
		State: req.State,
	}
}

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

func RegisterWithGoogleRequestToInput(in requests.RegisterWithGoogleRequest) dto.RegisterWithGoogleInput {
	return dto.RegisterWithGoogleInput{
		Token: in.Token,
	}
}

func RegisterWithGoogleOutputToResponse(out *dto.RegisterWithGoogleOutput) *responses.RegisterWithGoogleResponse {
	return &responses.RegisterWithGoogleResponse{
		AuthURL:         out.AuthURL,
		Token:           out.Token,
		HasOrganization: out.HasOrganization,
	}
}

func RequestInvitationGoogleAuthLinkRequestToInput(in requests.RequestInvitationGoogleAuthLinkRequest) dto.RequestInvitationGoogleAuthLinkInput {
	return dto.RequestInvitationGoogleAuthLinkInput{
		InvitationToken: in.InvitationToken,
	}
}

func RequestInvitationGoogleAuthLinkOutputToResponse(out *dto.RequestInvitationGoogleAuthLinkOutput) *responses.RequestInvitationGoogleAuthLinkResponse {
	return &responses.RequestInvitationGoogleAuthLinkResponse{
		AuthURL: out.AuthURL,
	}
}

func RefreshTokenRequestToInput(in requests.RefreshTokenRequest) dto.RefreshTokenInput {
	return dto.RefreshTokenInput{
		RefreshToken:    in.RefreshToken,
		XSRFTokenHeader: in.XSRFTokenHeader,
		XSRFTokenCookie: in.XSRFTokenCookie,
	}
}

func RefreshTokenOutputToResponse(out *dto.RefreshTokenOutput) *responses.RefreshTokenResponse {
	return &responses.RefreshTokenResponse{
		ExpiresAt: out.ExpiresAt,
	}
}

func SaveAuthRequestToInput(in requests.SaveAuthRequest) dto.SaveAuthInput {
	return dto.SaveAuthInput{
		Token: in.Token,
	}
}

func SaveAuthOutputToResponse(out *dto.SaveAuthOutput) *responses.SaveAuthResponse {
	return &responses.SaveAuthResponse{
		ExpiresAt:   out.ExpiresAt,
		RedirectURL: out.RedirectURL,
	}
}

func ObtainAuthTokenOutputToResponse(out *dto.ObtainAuthTokenOutput) *responses.ObtainAuthTokenResponse {
	return &responses.ObtainAuthTokenResponse{
		AuthURL: out.AuthURL,
		Token:   out.Token,
	}
}

func LogoutOutputToResponse(out *dto.LogoutOutput) *responses.LogoutResponse {
	return &responses.LogoutResponse{
		Domain: out.Domain,
	}
}
