package auth

import (
	"context"
	"fmt"
)

func (s *ServiceCE) sendMagicLinkEmail(ctx context.Context, email, firstName, url string) error {
	subject := "Log in to your Sourcetool account"

	content := fmt.Sprintf(`Hi %s,

Here's your magic link to log in to your Sourcetool account. Click the link below to access your account securely without a password:

%s

- This link will expire in 15 minutes for security reasons.
- If you didn't request this link, you can safely ignore this email.

Thank you for using Sourcetool!

The Sourcetool Team`, firstName, url)

	if err := s.Mailer.Send(ctx, []string{email}, "Sourcetool Team", subject, content); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *ServiceCE) sendInvitationMagicLinkEmail(ctx context.Context, email, firstName, url string) error {
	subject := "Your invitation to join Sourcetool"

	content := fmt.Sprintf(`Hi %s,

You've been invited to join Sourcetool. Click the link below to accept the invitation:

%s

This link will expire in 15 minutes.

Best regards,
The Sourcetool Team`, firstName, url)

	if err := s.Mailer.Send(ctx, []string{email}, "Sourcetool Team", subject, content); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func (s *ServiceCE) sendMultipleOrganizationsMagicLinkEmail(ctx context.Context, email, firstName string, loginURLs []string) error {
	subject := "Choose your Sourcetool organization to log in"

	urlList := ""
	for _, url := range loginURLs {
		urlList += url + "\n"
	}

	content := fmt.Sprintf(`Hi %s,

Your email, %s, is associated with multiple Sourcetool organizations. You may log in to each one by clicking its magic link below:

%s

Thank you for using Sourcetool!

The Sourcetool Team`, firstName, email, urlList)

	if err := s.Mailer.Send(ctx, []string{email}, "Sourcetool Team", subject, content); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *ServiceCE) sendMultipleOrganizationsLoginEmail(ctx context.Context, email, firstName string, loginURLs []string) error {
	subject := "Choose your Sourcetool organization to log in"

	urlList := ""
	for _, url := range loginURLs {
		urlList += url + "\n"
	}

	content := fmt.Sprintf(`Hi %s,

Your email, %s, is associated with multiple Sourcetool organizations. You may log in to each one by clicking its login link below:

%s

Thank you for using Sourcetool!

The Sourcetool Team`, firstName, email, urlList)

	if err := s.Mailer.Send(ctx, []string{email}, "Sourcetool Team", subject, content); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
