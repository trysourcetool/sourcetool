package dto

import (
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
)

type UpdateMeInput struct {
	FirstName *string
	LastName  *string
}

type SendUpdateMeEmailInstructionsInput struct {
	Email             string
	EmailConfirmation string
}

type UpdateMeEmailInput struct {
	Token string
}

type UpdateUserInput struct {
	UserID   string
	Role     *string
	GroupIDs []string
}

type DeleteUserInput struct {
	UserID string
}

type CreateUserInvitationsInput struct {
	Emails []string
	Role   string
}

type ResendUserInvitationInput struct {
	InvitationID string
}

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

type UserInvitation struct {
	ID        string
	Email     string
	CreatedAt int64
}

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

type UserGroup struct {
	ID        string
	UserID    string
	GroupID   string
	CreatedAt int64
	UpdatedAt int64
}

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

type GetMeOutput struct {
	User *User
}

type UpdateMeOutput struct {
	User *User
}

type UpdateMeEmailOutput struct {
	User *User
}

type ListUsersOutput struct {
	Users           []*User
	UserInvitations []*UserInvitation
}

type UpdateUserOutput struct {
	User *User
}

type CreateUserInvitationsOutput struct {
	UserInvitations []*UserInvitation
}

type ResendUserInvitationOutput struct {
	UserInvitation *UserInvitation
}
