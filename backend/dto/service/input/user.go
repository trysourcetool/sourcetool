package input

// UpdateMeInput is the input for Update Me operation.
type UpdateMeInput struct {
	FirstName *string
	LastName  *string
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

// UpdateUserInput is the input for Update User operation.
type UpdateUserInput struct {
	UserID   string
	Role     *string
	GroupIDs []string
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

// ResendUserInvitationInput is the input for Resend User Invitation operation.
type ResendUserInvitationInput struct {
	InvitationID string
}
