package adapters

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/dto/http/requests"
	"github.com/trysourcetool/sourcetool/backend/dto/http/responses"
	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
)

// UserOutputToResponse converts from output.User to responses.UserResponse.
func UserOutputToResponse(user *output.User) *responses.UserResponse {
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
		result.Organization = OrganizationOutputToResponse(user.Organization)
	}

	return result
}

// UserInvitationOutputToResponse converts from output.UserInvitation to responses.UserInvitationResponse.
func UserInvitationOutputToResponse(invitation *output.UserInvitation) *responses.UserInvitationResponse {
	if invitation == nil {
		return nil
	}

	return &responses.UserInvitationResponse{
		ID:        invitation.ID,
		Email:     invitation.Email,
		CreatedAt: strconv.FormatInt(invitation.CreatedAt, 10),
	}
}

// UserGroupOutputToResponse converts from output.UserGroup to responses.UserGroupResponse.
func UserGroupOutputToResponse(group *output.UserGroup) *responses.UserGroupResponse {
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

// GetMeOutputToResponse converts from output.GetMeOutput to responses.GetMeResponse.
func GetMeOutputToResponse(out *output.GetMeOutput) *responses.GetMeResponse {
	return &responses.GetMeResponse{
		User: UserOutputToResponse(out.User),
	}
}

// ListUsersOutputToResponse converts from output.ListUsersOutput to responses.ListUsersResponse.
func ListUsersOutputToResponse(out *output.ListUsersOutput) *responses.ListUsersResponse {
	users := make([]*responses.UserResponse, 0, len(out.Users))
	for _, user := range out.Users {
		users = append(users, UserOutputToResponse(user))
	}

	invitations := make([]*responses.UserInvitationResponse, 0, len(out.UserInvitations))
	for _, invitation := range out.UserInvitations {
		invitations = append(invitations, UserInvitationOutputToResponse(invitation))
	}

	return &responses.ListUsersResponse{
		Users:           users,
		UserInvitations: invitations,
	}
}

// UpdateMeRequestToInput converts from requests.UpdateMeRequest to input.UpdateMeInput.
func UpdateMeRequestToInput(in requests.UpdateMeRequest) input.UpdateMeInput {
	return input.UpdateMeInput{
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

// UpdateMeOutputToResponse converts from output.UpdateMeOutput to responses.UpdateMeResponse.
func UpdateMeOutputToResponse(out *output.UpdateMeOutput) *responses.UpdateMeResponse {
	return &responses.UpdateMeResponse{
		User: UserOutputToResponse(out.User),
	}
}

// SendUpdateMeEmailInstructionsRequestToInput converts from requests.SendUpdateMeEmailInstructionsRequest to input.SendUpdateMeEmailInstructionsInput.
func SendUpdateMeEmailInstructionsRequestToInput(in requests.SendUpdateMeEmailInstructionsRequest) input.SendUpdateMeEmailInstructionsInput {
	return input.SendUpdateMeEmailInstructionsInput{
		Email:             in.Email,
		EmailConfirmation: in.EmailConfirmation,
	}
}

// UpdateMeEmailRequestToInput converts from requests.UpdateMeEmailRequest to input.UpdateMeEmailInput.
func UpdateMeEmailRequestToInput(in requests.UpdateMeEmailRequest) input.UpdateMeEmailInput {
	return input.UpdateMeEmailInput{
		Token: in.Token,
	}
}

// UpdateUserEmailOutputToResponse converts from output.UpdateUserEmailOutput to responses.UpdateUserEmailResponse.
func UpdateMeEmailOutputToResponse(out *output.UpdateMeEmailOutput) *responses.UpdateMeEmailResponse {
	return &responses.UpdateMeEmailResponse{
		User: UserOutputToResponse(out.User),
	}
}

// UpdateUserRequestToInput converts from requests.UpdateUserRequest to input.UpdateUserInput.
func UpdateUserRequestToInput(in requests.UpdateUserRequest) input.UpdateUserInput {
	return input.UpdateUserInput{
		UserID:   in.UserID,
		Role:     in.Role,
		GroupIDs: in.GroupIDs,
	}
}

// UpdateUserOutputToResponse converts from output.UpdateUserOutput to responses.UpdateUserResponse.
func UpdateUserOutputToResponse(out *output.UpdateUserOutput) *responses.UpdateUserResponse {
	return &responses.UpdateUserResponse{
		User: UserOutputToResponse(out.User),
	}
}

// DeleteUserRequestToInput converts from requests.DeleteUserRequest to input.DeleteUserInput.
func DeleteUserRequestToInput(in requests.DeleteUserRequest) input.DeleteUserInput {
	return input.DeleteUserInput{
		UserID: in.UserID,
	}
}

// CreateUserInvitationsRequestToInput converts from requests.CreateUserInvitationsRequest to input.CreateUserInvitationsInput.
func CreateUserInvitationsRequestToInput(in requests.CreateUserInvitationsRequest) input.CreateUserInvitationsInput {
	return input.CreateUserInvitationsInput{
		Emails: in.Emails,
		Role:   in.Role,
	}
}

// CreateUserInvitationsOutputToResponse converts from output.CreateUserInvitationsOutput to responses.CreateUserInvitationsResponse.
func CreateUserInvitationsOutputToResponse(out *output.CreateUserInvitationsOutput) *responses.CreateUserInvitationsResponse {
	invitations := make([]*responses.UserInvitationResponse, 0, len(out.UserInvitations))
	for _, invitation := range out.UserInvitations {
		invitations = append(invitations, UserInvitationOutputToResponse(invitation))
	}

	return &responses.CreateUserInvitationsResponse{
		UserInvitations: invitations,
	}
}

// ResendInvitationRequestToInput converts from requests.ResendInvitationRequest to input.ResendInvitationInput.
func ResendUserInvitationRequestToInput(in requests.ResendUserInvitationRequest) input.ResendUserInvitationInput {
	return input.ResendUserInvitationInput{
		InvitationID: in.InvitationID,
	}
}

// ResendInvitationOutputToResponse converts from output.ResendUserInvitationOutput to responses.ResendUserInvitationResponse.
func ResendUserInvitationOutputToResponse(out *output.ResendUserInvitationOutput) *responses.ResendUserInvitationResponse {
	return &responses.ResendUserInvitationResponse{
		UserInvitation: UserInvitationOutputToResponse(out.UserInvitation),
	}
}
