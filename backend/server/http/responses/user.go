package responses

type SendSignUpInstructionsResponse struct {
	Email string `json:"email"`
}

type SignInResponse struct {
	AuthURL              string `json:"authUrl"`
	Token                string `json:"token"`
	IsOrganizationExists bool   `json:"isOrganizationExists"`
}

type SignInWithGoogleResponse struct {
	AuthURL              string `json:"authUrl"`
	Token                string `json:"token"`
	IsOrganizationExists bool   `json:"isOrganizationExists"`
}

type SignUpResponse struct {
	Token     string
	XSRFToken string
}

type SignUpWithGoogleResponse struct {
	Token     string
	XSRFToken string
}

type RefreshTokenResponse struct {
	ExpiresAt string `json:"expiresAt"`
}

type SaveAuthResponse struct {
	ExpiresAt   string `json:"expiresAt"`
	RedirectURL string `json:"redirectUrl"`
}

type ObtainAuthTokenResponse struct {
	AuthURL string `json:"authUrl"`
	Token   string `json:"token"`
}

type InviteUsersResponse struct {
	UserInvitations []*UserInvitationResponse `json:"userInvitations"`
}

type SignInInvitationResponse struct {
	ExpiresAt string `json:"expiresAt"`
}

type SignUpInvitationResponse struct {
	ExpiresAt string `json:"expiresAt"`
}

type GetGoogleAuthCodeURLResponse struct {
	URL string `json:"url"`
}

type GetGoogleAuthCodeURLInvitationResponse struct {
	URL string `json:"url"`
}

type SignInWithGoogleInvitationResponse struct {
	ExpiresAt string `json:"expiresAt"`
}

type SignUpWithGoogleInvitationResponse struct {
	ExpiresAt string `json:"expiresAt"`
}

type SignOutResponse struct {
	Domain string
}

type UserResponse struct {
	ID           string                `json:"id"`
	Email        string                `json:"email"`
	FirstName    string                `json:"firstName"`
	LastName     string                `json:"lastName"`
	Role         string                `json:"role"`
	CreatedAt    string                `json:"createdAt"`
	UpdatedAt    string                `json:"updatedAt"`
	Organization *OrganizationResponse `json:"organization"`
}

type UserInvitationResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type ListUsersResponse struct {
	Users           []*UserResponse           `json:"users"`
	UserInvitations []*UserInvitationResponse `json:"userInvitations"`
}

type GetMeResponse struct {
	User *UserResponse `json:"user"`
}

type UpdateUserResponse struct {
	User *UserResponse `json:"user"`
}

type UpdateUserPasswordResponse struct {
	User *UserResponse `json:"user"`
}

type UpdateUserEmailResponse struct {
	User *UserResponse `json:"user"`
}

type ResendInvitationResponse struct {
	UserInvitation *UserInvitationResponse `json:"userInvitation"`
}

type UserGroupResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	GroupID   string `json:"groupId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
