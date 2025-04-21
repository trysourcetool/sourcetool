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
	Role     *string  `json:"role" validate:"oneof=admin developer member"`
	GroupIDs []string `json:"groupIds"`
}
