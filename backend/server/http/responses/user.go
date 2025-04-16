package responses

type CreateUserInvitationsResponse struct {
	UserInvitations []*UserInvitationResponse `json:"userInvitations"`
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

type UpdateMeResponse struct {
	User *UserResponse `json:"user"`
}

type UpdateMeEmailResponse struct {
	User *UserResponse `json:"user"`
}

type UpdateUserResponse struct {
	User *UserResponse `json:"user"`
}

type ResendUserInvitationResponse struct {
	UserInvitation *UserInvitationResponse `json:"userInvitation"`
}

type UserGroupResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	GroupID   string `json:"groupId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
