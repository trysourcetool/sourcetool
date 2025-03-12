package dto

import (
	"github.com/trysourcetool/sourcetool/backend/model"
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
func UserFromModel(user *model.User, org *model.Organization, role model.UserOrganizationRole) *User {
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
func UserInvitationFromModel(invitation *model.UserInvitation) *UserInvitation {
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
func UserGroupFromModel(group *model.UserGroup) *UserGroup {
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

// ListUsersOutput is the output for List operation.
type ListUsersOutput struct {
	Users           []*User
	UserInvitations []*UserInvitation
}

// UpdateUserInput is the input for Update operation.
type UpdateUserInput struct {
	FirstName *string
	LastName  *string
}

// UpdateUserOutput is the output for Update operation.
type UpdateUserOutput struct {
	User *User
}

// SendUpdateUserEmailInstructionsInput is the input for Send Update User Email Instructions operation.
type SendUpdateUserEmailInstructionsInput struct {
	Email             string
	EmailConfirmation string
}

// UpdateUserEmailInput is the input for Update User Email operation.
type UpdateUserEmailInput struct {
	Token string
}

// UpdateUserEmailOutput is the output for Update User Email operation.
type UpdateUserEmailOutput struct {
	User *User
}

// UpdateUserPasswordInput is the input for Update User Password operation.
type UpdateUserPasswordInput struct {
	CurrentPassword      string
	Password             string
	PasswordConfirmation string
}

// UpdateUserPasswordOutput is the output for Update User Password operation.
type UpdateUserPasswordOutput struct {
	User *User
}

// SignInInput is the input for Sign In operation.
type SignInInput struct {
	Email    string
	Password string
}

// SignInOutput is the output for Sign In operation.
type SignInOutput struct {
	AuthURL              string
	Token                string
	IsOrganizationExists bool
	Secret               string
	XSRFToken            string
	Domain               string
}

// SignInWithGoogleInput is the input for Sign In With Google operation.
type SignInWithGoogleInput struct {
	SessionToken string
}

// SignInWithGoogleOutput is the output for Sign In With Google operation.
type SignInWithGoogleOutput struct {
	AuthURL              string
	Token                string
	IsOrganizationExists bool
	Secret               string
	XSRFToken            string
	Domain               string
}

// SendSignUpInstructionsInput is the input for Send SignUp Instructions operation.
type SendSignUpInstructionsInput struct {
	Email string
}

// SendSignUpInstructionsOutput is the output for Send SignUp Instructions operation.
type SendSignUpInstructionsOutput struct {
	Email string
}

// SignUpInput is the input for Sign Up operation.
type SignUpInput struct {
	Token                string
	FirstName            string
	LastName             string
	Password             string
	PasswordConfirmation string
}

// SignUpOutput is the output for Sign Up operation.
type SignUpOutput struct {
	Token     string
	XSRFToken string
}

// SignUpWithGoogleInput is the input for Sign Up With Google operation.
type SignUpWithGoogleInput struct {
	SessionToken string
	FirstName    string
	LastName     string
}

// SignUpWithGoogleOutput is the output for Sign Up With Google operation.
type SignUpWithGoogleOutput struct {
	Token     string
	XSRFToken string
}

// RefreshTokenInput is the input for Refresh Token operation.
type RefreshTokenInput struct {
	Secret          string
	XSRFTokenHeader string
	XSRFTokenCookie string
}

// RefreshTokenOutput is the output for Refresh Token operation.
type RefreshTokenOutput struct {
	Token     string
	Secret    string
	XSRFToken string
	ExpiresAt string
	Domain    string
}

// SaveAuthInput is the input for Save Auth operation.
type SaveAuthInput struct {
	Token string
}

// SaveAuthOutput is the output for Save Auth operation.
type SaveAuthOutput struct {
	Token       string
	Secret      string
	XSRFToken   string
	ExpiresAt   string
	RedirectURL string
	Domain      string
}

// ObtainAuthTokenOutput is the output for Obtain Auth Token operation.
type ObtainAuthTokenOutput struct {
	AuthURL   string
	Token     string
	Secret    string
	XSRFToken string
	Domain    string
}

// InviteUsersInput is the input for Invite Users operation.
type InviteUsersInput struct {
	Emails []string
	Role   string
}

// InviteUsersOutput is the output for Invite Users operation.
type InviteUsersOutput struct {
	UserInvitations []*UserInvitation
}

// SignInInvitationInput is the input for Sign In Invitation operation.
type SignInInvitationInput struct {
	InvitationToken string
	Password        string
}

// SignInInvitationOutput is the output for Sign In Invitation operation.
type SignInInvitationOutput struct {
	Token     string
	Secret    string
	XSRFToken string
	ExpiresAt string
	Domain    string
}

// SignUpInvitationInput is the input for Sign Up Invitation operation.
type SignUpInvitationInput struct {
	InvitationToken      string
	FirstName            string
	LastName             string
	Password             string
	PasswordConfirmation string
}

// SignUpInvitationOutput is the output for Sign Up Invitation operation.
type SignUpInvitationOutput struct {
	Token     string
	Secret    string
	XSRFToken string
	ExpiresAt string
	Domain    string
}

// GetGoogleAuthCodeURLOutput is the output for Get Google Auth Code URL operation.
type GetGoogleAuthCodeURLOutput struct {
	URL string
}

// GoogleOAuthCallbackInput is the input for Google OAuth Callback operation.
type GoogleOAuthCallbackInput struct {
	State string
	Code  string
}

// GoogleOAuthCallbackOutput is the output for Google OAuth Callback operation.
type GoogleOAuthCallbackOutput struct {
	SessionToken string
	IsUserExists bool
	FirstName    string
	LastName     string
	Domain       string
	Invited      bool
}

// GetGoogleAuthCodeURLInvitationInput is the input for Get Google Auth Code URL Invitation operation.
type GetGoogleAuthCodeURLInvitationInput struct {
	InvitationToken string
}

// GetGoogleAuthCodeURLInvitationOutput is the output for Get Google Auth Code URL Invitation operation.
type GetGoogleAuthCodeURLInvitationOutput struct {
	URL string
}

// SignInWithGoogleInvitationInput is the input for Sign In With Google Invitation operation.
type SignInWithGoogleInvitationInput struct {
	SessionToken string
}

// SignInWithGoogleInvitationOutput is the output for Sign In With Google Invitation operation.
type SignInWithGoogleInvitationOutput struct {
	Token     string
	Secret    string
	XSRFToken string
	ExpiresAt string
	Domain    string
}

// SignUpWithGoogleInvitationInput is the input for Sign Up With Google Invitation operation.
type SignUpWithGoogleInvitationInput struct {
	SessionToken string
	FirstName    string
	LastName     string
}

// SignUpWithGoogleInvitationOutput is the output for Sign Up With Google Invitation operation.
type SignUpWithGoogleInvitationOutput struct {
	Token     string
	Secret    string
	XSRFToken string
	ExpiresAt string
	Domain    string
}

// SignOutOutput is the output for Sign Out operation.
type SignOutOutput struct {
	Domain string
}

// ResendInvitationInput is the input for Resend Invitation operation.
type ResendInvitationInput struct {
	InvitationID string
}

// ResendInvitationOutput is the output for Resend Invitation operation.
type ResendInvitationOutput struct {
	UserInvitation *UserInvitation
}
