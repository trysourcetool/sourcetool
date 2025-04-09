package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
)

// UserDTOToResponse converts from dto.User to responses.UserResponse.
func UserDTOToResponse(user *dto.User) *responses.UserResponse {
	if user == nil {
		return nil
	}

	result := &responses.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: strconv.FormatInt(user.CreatedAt, 10),
		UpdatedAt: strconv.FormatInt(user.UpdatedAt, 10),
	}

	if user.Organization != nil {
		result.Organization = OrganizationDTOToResponse(user.Organization)
	}

	return result
}

// UserInvitationDTOToResponse converts from dto.UserInvitation to responses.UserInvitationResponse.
func UserInvitationDTOToResponse(invitation *dto.UserInvitation) *responses.UserInvitationResponse {
	if invitation == nil {
		return nil
	}

	return &responses.UserInvitationResponse{
		ID:        invitation.ID,
		Email:     invitation.Email,
		CreatedAt: strconv.FormatInt(invitation.CreatedAt, 10),
	}
}

// UserGroupDTOToResponse converts from dto.UserGroup to responses.UserGroupResponse.
func UserGroupDTOToResponse(group *dto.UserGroup) *responses.UserGroupResponse {
	if group == nil {
		return nil
	}

	return &responses.UserGroupResponse{
		ID:        group.ID,
		UserID:    group.UserID,
		GroupID:   group.GroupID,
		CreatedAt: strconv.FormatInt(group.CreatedAt, 10),
		UpdatedAt: strconv.FormatInt(group.UpdatedAt, 10),
	}
}

// GetMeOutputToResponse converts from dto.GetMeOutput to responses.GetMeResponse.
func GetMeOutputToResponse(out *dto.GetMeOutput) *responses.GetMeResponse {
	return &responses.GetMeResponse{
		User: UserDTOToResponse(out.User),
	}
}

// ListUsersOutputToResponse converts from dto.ListUsersOutput to responses.ListUsersResponse.
func ListUsersOutputToResponse(out *dto.ListUsersOutput) *responses.ListUsersResponse {
	users := make([]*responses.UserResponse, 0, len(out.Users))
	for _, user := range out.Users {
		users = append(users, UserDTOToResponse(user))
	}

	invitations := make([]*responses.UserInvitationResponse, 0, len(out.UserInvitations))
	for _, invitation := range out.UserInvitations {
		invitations = append(invitations, UserInvitationDTOToResponse(invitation))
	}

	return &responses.ListUsersResponse{
		Users:           users,
		UserInvitations: invitations,
	}
}

// UpdateUserRequestToDTOInput converts from requests.UpdateUserRequest to dto.UpdateUserInput.
func UpdateUserRequestToDTOInput(in requests.UpdateUserRequest) dto.UpdateUserInput {
	return dto.UpdateUserInput{
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

// UpdateUserOutputToResponse converts from dto.UpdateUserOutput to responses.UpdateUserResponse.
func UpdateUserOutputToResponse(out *dto.UpdateUserOutput) *responses.UpdateUserResponse {
	return &responses.UpdateUserResponse{
		User: UserDTOToResponse(out.User),
	}
}

// SendUpdateUserEmailInstructionsRequestToDTOInput converts from requests.SendUpdateUserEmailInstructionsRequest to dto.SendUpdateUserEmailInstructionsInput.
func SendUpdateUserEmailInstructionsRequestToDTOInput(in requests.SendUpdateUserEmailInstructionsRequest) dto.SendUpdateUserEmailInstructionsInput {
	return dto.SendUpdateUserEmailInstructionsInput{
		Email:             in.Email,
		EmailConfirmation: in.EmailConfirmation,
	}
}

// UpdateUserEmailRequestToDTOInput converts from requests.UpdateUserEmailRequest to dto.UpdateUserEmailInput.
func UpdateUserEmailRequestToDTOInput(in requests.UpdateUserEmailRequest) dto.UpdateUserEmailInput {
	return dto.UpdateUserEmailInput{
		Token: in.Token,
	}
}

// UpdateUserEmailOutputToResponse converts from dto.UpdateUserEmailOutput to responses.UpdateUserEmailResponse.
func UpdateUserEmailOutputToResponse(out *dto.UpdateUserEmailOutput) *responses.UpdateUserEmailResponse {
	return &responses.UpdateUserEmailResponse{
		User: UserDTOToResponse(out.User),
	}
}

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
		AuthURL:              out.AuthURL,
		Token:                out.Token,
		IsOrganizationExists: out.IsOrganizationExists,
		IsNewUser:            out.IsNewUser,
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
		ExpiresAt: out.ExpiresAt,
	}
}

// RefreshTokenRequestToDTOInput converts from requests.RefreshTokenRequest to dto.RefreshTokenInput.
func RefreshTokenRequestToDTOInput(in requests.RefreshTokenRequest) dto.RefreshTokenInput {
	return dto.RefreshTokenInput{
		Secret:          in.Secret,
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

// InviteUsersRequestToDTOInput converts from requests.InviteUsersRequest to dto.InviteUsersInput.
func InviteUsersRequestToDTOInput(in requests.InviteUsersRequest) dto.InviteUsersInput {
	return dto.InviteUsersInput{
		Emails: in.Emails,
		Role:   in.Role,
	}
}

// InviteUsersOutputToResponse converts from dto.InviteUsersOutput to responses.InviteUsersResponse.
func InviteUsersOutputToResponse(out *dto.InviteUsersOutput) *responses.InviteUsersResponse {
	invitations := make([]*responses.UserInvitationResponse, 0, len(out.UserInvitations))
	for _, invitation := range out.UserInvitations {
		invitations = append(invitations, UserInvitationDTOToResponse(invitation))
	}

	return &responses.InviteUsersResponse{
		UserInvitations: invitations,
	}
}

// ResendInvitationRequestToDTOInput converts from requests.ResendInvitationRequest to dto.ResendInvitationInput.
func ResendInvitationRequestToDTOInput(in requests.ResendInvitationRequest) dto.ResendInvitationInput {
	return dto.ResendInvitationInput{
		InvitationID: in.InvitationID,
	}
}

// ResendInvitationOutputToResponse converts from dto.ResendInvitationOutput to responses.ResendInvitationResponse.
func ResendInvitationOutputToResponse(out *dto.ResendInvitationOutput) *responses.ResendInvitationResponse {
	return &responses.ResendInvitationResponse{
		UserInvitation: UserInvitationDTOToResponse(out.UserInvitation),
	}
}

// SignOutOutputToResponse converts from dto.SignOutOutput to responses.SignOutResponse.
func SignOutOutputToResponse(out *dto.SignOutOutput) *responses.SignOutResponse {
	return &responses.SignOutResponse{
		Domain: out.Domain,
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
		FirstName:            out.FirstName,
		LastName:             out.LastName,
		AuthURL:              out.AuthURL,
		Token:                out.Token,
		IsOrganizationExists: out.IsOrganizationExists,
		IsNewUser:            out.IsNewUser,
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
		AuthURL:              out.AuthURL,
		Token:                out.Token,
		IsOrganizationExists: out.IsOrganizationExists,
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

// AuthenticateWithInvitationGoogleAuthLinkRequestToDTOInput converts request to DTO input.
func AuthenticateWithInvitationGoogleAuthLinkRequestToDTOInput(in requests.AuthenticateWithInvitationGoogleAuthLinkRequest) dto.AuthenticateWithInvitationGoogleAuthLinkInput {
	return dto.AuthenticateWithInvitationGoogleAuthLinkInput{
		Code:  in.Code,
		State: in.State,
	}
}
