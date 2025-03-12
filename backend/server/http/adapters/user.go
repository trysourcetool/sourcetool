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

// UpdateUserPasswordRequestToDTOInput converts from requests.UpdateUserPasswordRequest to dto.UpdateUserPasswordInput.
func UpdateUserPasswordRequestToDTOInput(in requests.UpdateUserPasswordRequest) dto.UpdateUserPasswordInput {
	return dto.UpdateUserPasswordInput{
		CurrentPassword:      in.CurrentPassword,
		Password:             in.Password,
		PasswordConfirmation: in.PasswordConfirmation,
	}
}

// UpdateUserPasswordOutputToResponse converts from dto.UpdateUserPasswordOutput to responses.UpdateUserPasswordResponse.
func UpdateUserPasswordOutputToResponse(out *dto.UpdateUserPasswordOutput) *responses.UpdateUserPasswordResponse {
	return &responses.UpdateUserPasswordResponse{
		User: UserDTOToResponse(out.User),
	}
}

// SignInRequestToDTOInput converts from requests.SignInRequest to dto.SignInInput.
func SignInRequestToDTOInput(in requests.SignInRequest) dto.SignInInput {
	return dto.SignInInput{
		Email:    in.Email,
		Password: in.Password,
	}
}

// SignInOutputToResponse converts from dto.SignInOutput to responses.SignInResponse.
func SignInOutputToResponse(out *dto.SignInOutput) *responses.SignInResponse {
	return &responses.SignInResponse{
		AuthURL:              out.AuthURL,
		Token:                out.Token,
		IsOrganizationExists: out.IsOrganizationExists,
	}
}

// SignInWithGoogleRequestToDTOInput converts from requests.SignInWithGoogleRequest to dto.SignInWithGoogleInput.
func SignInWithGoogleRequestToDTOInput(in requests.SignInWithGoogleRequest) dto.SignInWithGoogleInput {
	return dto.SignInWithGoogleInput{
		SessionToken: in.SessionToken,
	}
}

// SignInWithGoogleOutputToResponse converts from dto.SignInWithGoogleOutput to responses.SignInWithGoogleResponse.
func SignInWithGoogleOutputToResponse(out *dto.SignInWithGoogleOutput) *responses.SignInWithGoogleResponse {
	return &responses.SignInWithGoogleResponse{
		AuthURL:              out.AuthURL,
		Token:                out.Token,
		IsOrganizationExists: out.IsOrganizationExists,
	}
}

// SendSignUpInstructionsRequestToDTOInput converts from requests.SendSignUpInstructionsRequest to dto.SendSignUpInstructionsInput.
func SendSignUpInstructionsRequestToDTOInput(in requests.SendSignUpInstructionsRequest) dto.SendSignUpInstructionsInput {
	return dto.SendSignUpInstructionsInput{
		Email: in.Email,
	}
}

// SendSignUpInstructionsOutputToResponse converts from dto.SendSignUpInstructionsOutput to responses.SendSignUpInstructionsResponse.
func SendSignUpInstructionsOutputToResponse(out *dto.SendSignUpInstructionsOutput) *responses.SendSignUpInstructionsResponse {
	return &responses.SendSignUpInstructionsResponse{
		Email: out.Email,
	}
}

// SignUpRequestToDTOInput converts from requests.SignUpRequest to dto.SignUpInput.
func SignUpRequestToDTOInput(in requests.SignUpRequest) dto.SignUpInput {
	return dto.SignUpInput{
		Token:                in.Token,
		FirstName:            in.FirstName,
		LastName:             in.LastName,
		Password:             in.Password,
		PasswordConfirmation: in.PasswordConfirmation,
	}
}

// SignUpOutputToResponse converts from dto.SignUpOutput to responses.SignUpResponse.
func SignUpOutputToResponse(out *dto.SignUpOutput) *responses.SignUpResponse {
	return &responses.SignUpResponse{
		Token:     out.Token,
		XSRFToken: out.XSRFToken,
	}
}

// SignUpWithGoogleRequestToDTOInput converts from requests.SignUpWithGoogleRequest to dto.SignUpWithGoogleInput.
func SignUpWithGoogleRequestToDTOInput(in requests.SignUpWithGoogleRequest) dto.SignUpWithGoogleInput {
	return dto.SignUpWithGoogleInput{
		SessionToken: in.SessionToken,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
	}
}

// SignUpWithGoogleOutputToResponse converts from dto.SignUpWithGoogleOutput to responses.SignUpWithGoogleResponse.
func SignUpWithGoogleOutputToResponse(out *dto.SignUpWithGoogleOutput) *responses.SignUpWithGoogleResponse {
	return &responses.SignUpWithGoogleResponse{
		Token:     out.Token,
		XSRFToken: out.XSRFToken,
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

// SignInInvitationRequestToDTOInput converts from requests.SignInInvitationRequest to dto.SignInInvitationInput.
func SignInInvitationRequestToDTOInput(in requests.SignInInvitationRequest) dto.SignInInvitationInput {
	return dto.SignInInvitationInput{
		InvitationToken: in.InvitationToken,
		Password:        in.Password,
	}
}

// SignInInvitationOutputToResponse converts from dto.SignInInvitationOutput to responses.SignInInvitationResponse.
func SignInInvitationOutputToResponse(out *dto.SignInInvitationOutput) *responses.SignInInvitationResponse {
	return &responses.SignInInvitationResponse{
		ExpiresAt: out.ExpiresAt,
	}
}

// SignUpInvitationRequestToDTOInput converts from requests.SignUpInvitationRequest to dto.SignUpInvitationInput.
func SignUpInvitationRequestToDTOInput(in requests.SignUpInvitationRequest) dto.SignUpInvitationInput {
	return dto.SignUpInvitationInput{
		InvitationToken:      in.InvitationToken,
		FirstName:            in.FirstName,
		LastName:             in.LastName,
		Password:             in.Password,
		PasswordConfirmation: in.PasswordConfirmation,
	}
}

// SignUpInvitationOutputToResponse converts from dto.SignUpInvitationOutput to responses.SignUpInvitationResponse.
func SignUpInvitationOutputToResponse(out *dto.SignUpInvitationOutput) *responses.SignUpInvitationResponse {
	return &responses.SignUpInvitationResponse{
		ExpiresAt: out.ExpiresAt,
	}
}

// GoogleOAuthCallbackRequestToDTOInput converts from requests.GoogleOAuthCallbackRequest to dto.GoogleOAuthCallbackInput.
func GoogleOAuthCallbackRequestToDTOInput(in requests.GoogleOAuthCallbackRequest) dto.GoogleOAuthCallbackInput {
	return dto.GoogleOAuthCallbackInput{
		State: in.State,
		Code:  in.Code,
	}
}

// GetGoogleAuthCodeURLInvitationRequestToDTOInput converts from requests.GetGoogleAuthCodeURLInvitationRequest to dto.GetGoogleAuthCodeURLInvitationInput.
func GetGoogleAuthCodeURLInvitationRequestToDTOInput(in requests.GetGoogleAuthCodeURLInvitationRequest) dto.GetGoogleAuthCodeURLInvitationInput {
	return dto.GetGoogleAuthCodeURLInvitationInput{
		InvitationToken: in.InvitationToken,
	}
}

// GetGoogleAuthCodeURLInvitationOutputToResponse converts from dto.GetGoogleAuthCodeURLInvitationOutput to responses.GetGoogleAuthCodeURLInvitationResponse.
func GetGoogleAuthCodeURLInvitationOutputToResponse(out *dto.GetGoogleAuthCodeURLInvitationOutput) *responses.GetGoogleAuthCodeURLInvitationResponse {
	return &responses.GetGoogleAuthCodeURLInvitationResponse{
		URL: out.URL,
	}
}

// SignInWithGoogleInvitationRequestToDTOInput converts from requests.SignInWithGoogleInvitationRequest to dto.SignInWithGoogleInvitationInput.
func SignInWithGoogleInvitationRequestToDTOInput(in requests.SignInWithGoogleInvitationRequest) dto.SignInWithGoogleInvitationInput {
	return dto.SignInWithGoogleInvitationInput{
		SessionToken: in.SessionToken,
	}
}

// SignInWithGoogleInvitationOutputToResponse converts from dto.SignInWithGoogleInvitationOutput to responses.SignInWithGoogleInvitationResponse.
func SignInWithGoogleInvitationOutputToResponse(out *dto.SignInWithGoogleInvitationOutput) *responses.SignInWithGoogleInvitationResponse {
	return &responses.SignInWithGoogleInvitationResponse{
		ExpiresAt: out.ExpiresAt,
	}
}

// SignUpWithGoogleInvitationRequestToDTOInput converts from requests.SignUpWithGoogleInvitationRequest to dto.SignUpWithGoogleInvitationInput.
func SignUpWithGoogleInvitationRequestToDTOInput(in requests.SignUpWithGoogleInvitationRequest) dto.SignUpWithGoogleInvitationInput {
	return dto.SignUpWithGoogleInvitationInput{
		SessionToken: in.SessionToken,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
	}
}

// SignUpWithGoogleInvitationOutputToResponse converts from dto.SignUpWithGoogleInvitationOutput to responses.SignUpWithGoogleInvitationResponse.
func SignUpWithGoogleInvitationOutputToResponse(out *dto.SignUpWithGoogleInvitationOutput) *responses.SignUpWithGoogleInvitationResponse {
	return &responses.SignUpWithGoogleInvitationResponse{
		ExpiresAt: out.ExpiresAt,
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
