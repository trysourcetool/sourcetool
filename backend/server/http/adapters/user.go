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

// UpdateMeRequestToDTOInput converts from requests.UpdateMeRequest to dto.UpdateMeInput.
func UpdateMeRequestToDTOInput(in requests.UpdateMeRequest) dto.UpdateMeInput {
	return dto.UpdateMeInput{
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

// UpdateMeOutputToResponse converts from dto.UpdateMeOutput to responses.UpdateMeResponse.
func UpdateMeOutputToResponse(out *dto.UpdateMeOutput) *responses.UpdateMeResponse {
	return &responses.UpdateMeResponse{
		User: UserDTOToResponse(out.User),
	}
}

// SendUpdateMeEmailInstructionsRequestToDTOInput converts from requests.SendUpdateMeEmailInstructionsRequest to dto.SendUpdateMeEmailInstructionsInput.
func SendUpdateMeEmailInstructionsRequestToDTOInput(in requests.SendUpdateMeEmailInstructionsRequest) dto.SendUpdateMeEmailInstructionsInput {
	return dto.SendUpdateMeEmailInstructionsInput{
		Email:             in.Email,
		EmailConfirmation: in.EmailConfirmation,
	}
}

// UpdateUserEmailRequestToDTOInput converts from requests.UpdateUserEmailRequest to dto.UpdateUserEmailInput.
func UpdateMeEmailRequestToDTOInput(in requests.UpdateMeEmailRequest) dto.UpdateMeEmailInput {
	return dto.UpdateMeEmailInput{
		Token: in.Token,
	}
}

// UpdateMeEmailOutputToResponse converts from dto.UpdateMeEmailOutput to responses.UpdateMeEmailResponse.
func UpdateMeEmailOutputToResponse(out *dto.UpdateMeEmailOutput) *responses.UpdateMeEmailResponse {
	return &responses.UpdateMeEmailResponse{
		User: UserDTOToResponse(out.User),
	}
}

// UpdateUserRequestToDTOInput converts from requests.UpdateUserRequest to dto.UpdateUserInput.
func UpdateUserRequestToDTOInput(in requests.UpdateUserRequest) dto.UpdateUserInput {
	return dto.UpdateUserInput{
		UserID:   in.UserID,
		Role:     in.Role,
		GroupIDs: in.GroupIDs,
	}
}

// UpdateUserOutputToResponse converts from dto.UpdateUserOutput to responses.UpdateUserResponse.
func UpdateUserOutputToResponse(out *dto.UpdateUserOutput) *responses.UpdateUserResponse {
	return &responses.UpdateUserResponse{
		User: UserDTOToResponse(out.User),
	}
}

// DeleteOrganizationUserRequestToDTOInput converts from requests.DeleteOrganizationUserRequest to dto.DeleteOrganizationUserInput.
func DeleteUserRequestToDTOInput(in requests.DeleteUserRequest) dto.DeleteUserInput {
	return dto.DeleteUserInput{
		UserID: in.UserID,
	}
}

// CreateUserInvitationsRequestToDTOInput converts from requests.CreateUserInvitationsRequest to dto.CreateUserInvitationsInput.
func CreateUserInvitationsRequestToDTOInput(in requests.CreateUserInvitationsRequest) dto.CreateUserInvitationsInput {
	return dto.CreateUserInvitationsInput{
		Emails: in.Emails,
		Role:   in.Role,
	}
}

// CreateUserInvitationsOutputToResponse converts from dto.CreateUserInvitationsOutput to responses.CreateUserInvitationsResponse.
func CreateUserInvitationsOutputToResponse(out *dto.CreateUserInvitationsOutput) *responses.CreateUserInvitationsResponse {
	invitations := make([]*responses.UserInvitationResponse, 0, len(out.UserInvitations))
	for _, invitation := range out.UserInvitations {
		invitations = append(invitations, UserInvitationDTOToResponse(invitation))
	}

	return &responses.CreateUserInvitationsResponse{
		UserInvitations: invitations,
	}
}

// ResendInvitationRequestToDTOInput converts from requests.ResendInvitationRequest to dto.ResendInvitationInput.
func ResendUserInvitationRequestToDTOInput(in requests.ResendUserInvitationRequest) dto.ResendUserInvitationInput {
	return dto.ResendUserInvitationInput{
		InvitationID: in.InvitationID,
	}
}

// ResendInvitationOutputToResponse converts from dto.ResendInvitationOutput to responses.ResendInvitationResponse.
func ResendUserInvitationOutputToResponse(out *dto.ResendUserInvitationOutput) *responses.ResendUserInvitationResponse {
	return &responses.ResendUserInvitationResponse{
		UserInvitation: UserInvitationDTOToResponse(out.UserInvitation),
	}
}
