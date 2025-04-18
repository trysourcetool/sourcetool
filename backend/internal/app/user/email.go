package user

import (
	"context"
	"fmt"

	"github.com/samber/lo"
)

func (s *ServiceCE) sendUpdateEmailInstructions(ctx context.Context, to, firstName, url string) error {
	subject := "[Sourcetool] Confirm your new email address"
	content := fmt.Sprintf(`Hi %s,

We received a request to change the email address associated with your Sourcetool account. To ensure the security of your account, we need you to verify your new email address.

Please click the following link within the next 24 hours to confirm your email change:
%s

Thank you for being a part of the Sourcetool community!
Regards,

Sourcetool Team`,
		firstName,
		url,
	)

	return s.Mailer.Send(ctx, []string{to}, "Sourcetool Team", subject, content)
}

func (s *ServiceCE) sendInvitationEmail(ctx context.Context, invitees string, emaiURLs map[string]string) error {
	subject := "You've been invited to join Sourcetool!"
	baseContent := `You've been invited to join Sourcetool!

To accept the invitation, please create your account by clicking the URL below within 24 hours.

%s

- This URL will expire in 24 hours.
- This is a send-only email address.
- Your account will not be created unless you complete the next steps.`

	sendEmails := make([]string, 0)
	for email, url := range emaiURLs {
		if lo.Contains(sendEmails, email) {
			continue
		}

		content := fmt.Sprintf(baseContent, url)

		if err := s.Mailer.Send(ctx, []string{email}, "Sourcetool Team", subject, content); err != nil {
			return err
		}

		sendEmails = append(sendEmails, email)
	}

	return nil
}
