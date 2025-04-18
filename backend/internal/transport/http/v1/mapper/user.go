package mapper

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/responses"
)

func UserOutputToResponse(user *dto.User) *responses.UserResponse {
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

func UserInvitationOutputToResponse(invitation *dto.UserInvitation) *responses.UserInvitationResponse {
	if invitation == nil {
		return nil
	}

	return &responses.UserInvitationResponse{
		ID:        invitation.ID,
		Email:     invitation.Email,
		CreatedAt: strconv.FormatInt(invitation.CreatedAt, 10),
	}
}

func UserGroupOutputToResponse(group *dto.UserGroup) *responses.UserGroupResponse {
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

func GetMeOutputToResponse(out *dto.GetMeOutput) *responses.GetMeResponse {
	return &responses.GetMeResponse{
		User: UserOutputToResponse(out.User),
	}
}

func ListUsersOutputToResponse(out *dto.ListUsersOutput) *responses.ListUsersResponse {
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

func UpdateMeRequestToInput(in requests.UpdateMeRequest) dto.UpdateMeInput {
	return dto.UpdateMeInput{
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

func UpdateMeOutputToResponse(out *dto.UpdateMeOutput) *responses.UpdateMeResponse {
	return &responses.UpdateMeResponse{
		User: UserOutputToResponse(out.User),
	}
}

func SendUpdateMeEmailInstructionsRequestToInput(in requests.SendUpdateMeEmailInstructionsRequest) dto.SendUpdateMeEmailInstructionsInput {
	return dto.SendUpdateMeEmailInstructionsInput{
		Email:             in.Email,
		EmailConfirmation: in.EmailConfirmation,
	}
}

func UpdateMeEmailRequestToInput(in requests.UpdateMeEmailRequest) dto.UpdateMeEmailInput {
	return dto.UpdateMeEmailInput{
		Token: in.Token,
	}
}

func UpdateMeEmailOutputToResponse(out *dto.UpdateMeEmailOutput) *responses.UpdateMeEmailResponse {
	return &responses.UpdateMeEmailResponse{
		User: UserOutputToResponse(out.User),
	}
}

func UpdateUserRequestToInput(in requests.UpdateUserRequest) dto.UpdateUserInput {
	return dto.UpdateUserInput{
		UserID:   in.UserID,
		Role:     in.Role,
		GroupIDs: in.GroupIDs,
	}
}

func UpdateUserOutputToResponse(out *dto.UpdateUserOutput) *responses.UpdateUserResponse {
	return &responses.UpdateUserResponse{
		User: UserOutputToResponse(out.User),
	}
}

func DeleteUserRequestToInput(in requests.DeleteUserRequest) dto.DeleteUserInput {
	return dto.DeleteUserInput{
		UserID: in.UserID,
	}
}

func CreateUserInvitationsRequestToInput(in requests.CreateUserInvitationsRequest) dto.CreateUserInvitationsInput {
	return dto.CreateUserInvitationsInput{
		Emails: in.Emails,
		Role:   in.Role,
	}
}

func CreateUserInvitationsOutputToResponse(out *dto.CreateUserInvitationsOutput) *responses.CreateUserInvitationsResponse {
	invitations := make([]*responses.UserInvitationResponse, 0, len(out.UserInvitations))
	for _, invitation := range out.UserInvitations {
		invitations = append(invitations, UserInvitationOutputToResponse(invitation))
	}

	return &responses.CreateUserInvitationsResponse{
		UserInvitations: invitations,
	}
}

func ResendUserInvitationRequestToInput(in requests.ResendUserInvitationRequest) dto.ResendUserInvitationInput {
	return dto.ResendUserInvitationInput{
		InvitationID: in.InvitationID,
	}
}

func ResendUserInvitationOutputToResponse(out *dto.ResendUserInvitationOutput) *responses.ResendUserInvitationResponse {
	return &responses.ResendUserInvitationResponse{
		UserInvitation: UserInvitationOutputToResponse(out.UserInvitation),
	}
}
