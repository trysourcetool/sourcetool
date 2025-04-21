package responses

import (
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

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

func UserFromModel(user *core.User, role core.UserOrganizationRole, org *core.Organization) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:           user.ID.String(),
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         role.String(),
		CreatedAt:    strconv.FormatInt(user.CreatedAt.Unix(), 10),
		UpdatedAt:    strconv.FormatInt(user.UpdatedAt.Unix(), 10),
		Organization: OrganizationFromModel(org),
	}
}

type UserInvitationResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

func UserInvitationFromModel(invitation *core.UserInvitation) *UserInvitationResponse {
	return &UserInvitationResponse{
		ID:        invitation.ID.String(),
		Email:     invitation.Email,
		CreatedAt: strconv.FormatInt(invitation.CreatedAt.Unix(), 10),
	}
}

type UserGroupResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	GroupID   string `json:"groupId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func UserGroupFromModel(userGroup *core.UserGroup) *UserGroupResponse {
	return &UserGroupResponse{
		ID:        userGroup.ID.String(),
		UserID:    userGroup.UserID.String(),
		GroupID:   userGroup.GroupID.String(),
		CreatedAt: strconv.FormatInt(userGroup.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(userGroup.UpdatedAt.Unix(), 10),
	}
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
