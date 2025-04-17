package output

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

// UpdateMeOutput is the output for Update Me operation.
type UpdateMeOutput struct {
	User *User
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

// UpdateUserOutput is the output for Update User operation.
type UpdateUserOutput struct {
	User *User
}

// CreateUserInvitationsOutput is the output for Create User Invitations operation.
type CreateUserInvitationsOutput struct {
	UserInvitations []*UserInvitation
}

// ResendUserInvitationOutput is the output for Resend User Invitation operation.
type ResendUserInvitationOutput struct {
	UserInvitation *UserInvitation
}
