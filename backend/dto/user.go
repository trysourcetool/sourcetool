package dto

import (
	"github.com/trysourcetool/sourcetool/backend/organization"
	"github.com/trysourcetool/sourcetool/backend/user"
)

// User represents user data in DTOs.
type User struct {
	ID           string
	Email        string
	FirstName    string
	LastName     string
	Role         string
	CreatedAt    int64
	UpdatedAt    int64
	Organization *Organization
}

// UserFromModel converts from model.User to dto.User.
func UserFromModel(user *user.User, org *organization.Organization, role user.UserOrganizationRole) *User {
	if user == nil {
		return nil
	}

	result := &User{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      role.String(),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}

	if org != nil {
		result.Organization = OrganizationFromModel(org)
	}

	return result
}

// UserInvitation represents user invitation data in DTOs.
type UserInvitation struct {
	ID        string
	Email     string
	CreatedAt int64
}

// UserInvitationFromModel converts from model.UserInvitation to dto.UserInvitation.
func UserInvitationFromModel(invitation *user.UserInvitation) *UserInvitation {
	if invitation == nil {
		return nil
	}

	return &UserInvitation{
		ID:        invitation.ID.String(),
		Email:     invitation.Email,
		CreatedAt: invitation.CreatedAt.Unix(),
	}
}

// UserGroup represents user group data in DTOs.
type UserGroup struct {
	ID        string
	UserID    string
	GroupID   string
	CreatedAt int64
	UpdatedAt int64
}

// UserGroupFromModel converts from model.UserGroup to dto.UserGroup.
func UserGroupFromModel(group *user.UserGroup) *UserGroup {
	if group == nil {
		return nil
	}

	return &UserGroup{
		ID:        group.ID.String(),
		UserID:    group.UserID.String(),
		GroupID:   group.GroupID.String(),
		CreatedAt: group.CreatedAt.Unix(),
		UpdatedAt: group.UpdatedAt.Unix(),
	}
}

// GetMeOutput is the output for GetMe operation.
type GetMeOutput struct {
	User *User
}

// UpdateMeInput is the input for Update Me operation.
type UpdateMeInput struct {
	FirstName *string
	LastName  *string
}

// UpdateMeOutput is the output for Update Me operation.
type UpdateMeOutput struct {
	User *User
}

// SendUpdateMeEmailInstructionsInput is the input for Send Update Me Email Instructions operation.
type SendUpdateMeEmailInstructionsInput struct {
	Email             string
	EmailConfirmation string
}

// UpdateMeEmailInput is the input for Update Me Email operation.
type UpdateMeEmailInput struct {
	Token string
}

// UpdateMeEmailOutput is the output for Update Me Email operation.
type UpdateMeEmailOutput struct {
	User *User
}

// ListUsersOutput is the output for List operation.
type ListUsersOutput struct {
	Users           []*User
	UserInvitations []*UserInvitation
}

// UpdateUserInput is the input for Update User operation.
type UpdateUserInput struct {
	UserID   string
	Role     *string
	GroupIDs []string
}

// UpdateUserOutput is the output for Update User operation.
type UpdateUserOutput struct {
	User *User
}

// DeleteUserInput defines the input for deleting a user.
type DeleteUserInput struct {
	UserID string
}

// CreateUserInvitationsInput is the input for Create User Invitations operation.
type CreateUserInvitationsInput struct {
	Emails []string
	Role   string
}

// CreateUserInvitationsOutput is the output for Create User Invitations operation.
type CreateUserInvitationsOutput struct {
	UserInvitations []*UserInvitation
}

// ResendUserInvitationInput is the input for Resend User Invitation operation.
type ResendUserInvitationInput struct {
	InvitationID string
}

// ResendUserInvitationOutput is the output for Resend User Invitation operation.
type ResendUserInvitationOutput struct {
	UserInvitation *UserInvitation
}
