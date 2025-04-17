package user

import "context"

type SendUpdateUserEmailInstructions struct {
	To        string
	FirstName string
	URL       string
}

type SendWelcomeEmail struct {
	To string
}

type SendInvitationEmail struct {
	Invitees string
	URLs     map[string]string // email -> url
}

// Email structure for sending magic link email.
type SendMagicLinkEmail struct {
	To        string
	FirstName string
	URL       string
}

// Email structure for sending multiple organizations magic link email.
type SendMultipleOrganizationsMagicLinkEmail struct {
	To        string
	FirstName string
	Email     string
	LoginURLs []string
}

// Email structure for sending multiple organizations login email.
type SendMultipleOrganizationsLoginEmail struct {
	To        string
	FirstName string
	Email     string
	LoginURLs []string
}

// SendInvitationMagicLinkEmail represents the data needed to send an invitation magic link email.
type SendInvitationMagicLinkEmail struct {
	To        string
	URL       string
	FirstName string
}

type UserMailer interface {
	SendUpdateEmailInstructions(ctx context.Context, in *SendUpdateUserEmailInstructions) error
	SendInvitationEmail(ctx context.Context, in *SendInvitationEmail) error
	SendMagicLinkEmail(ctx context.Context, in *SendMagicLinkEmail) error
	SendMultipleOrganizationsMagicLinkEmail(ctx context.Context, in *SendMultipleOrganizationsMagicLinkEmail) error
	SendMultipleOrganizationsLoginEmail(ctx context.Context, in *SendMultipleOrganizationsLoginEmail) error
	SendInvitationMagicLinkEmail(ctx context.Context, in *SendInvitationMagicLinkEmail) error
}
