package requests

type CreateUserInvitationsRequest struct {
	Emails []string `json:"emails" validate:"required"`
	Role   string   `json:"role" validate:"required,oneof=admin developer member"`
}

type UpdateMeRequest struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

type SendUpdateMeEmailInstructionsRequest struct {
	Email             string `json:"email" validate:"required,email"`
	EmailConfirmation string `json:"emailConfirmation" validate:"required,email"`
}

type UpdateMeEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type UpdateUserRequest struct {
	UserID   string   `json:"-" validate:"required,uuid4"`
	Role     *string  `json:"role" validate:"oneof=admin developer member"`
	GroupIDs []string `json:"groupIds"`
}

type DeleteUserRequest struct {
	UserID string `param:"userID" validate:"required,uuid4"`
}

type ResendUserInvitationRequest struct {
	InvitationID string `json:"invitationId" validate:"required,uuid"`
}
