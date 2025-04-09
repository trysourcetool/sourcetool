package user

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/utils/emailutil"
)

type MailerCE struct {
	emailClient *emailutil.EmailClient
}

func NewMailerCE() *MailerCE {
	return &MailerCE{
		emailClient: emailutil.NewEmailClient(),
	}
}

func (m *MailerCE) SendUpdateEmailInstructions(ctx context.Context, in *model.SendUpdateUserEmailInstructions) error {
	subject := "[Sourcetool] Confirm your new email address"
	content := fmt.Sprintf(`Hi %s,

We received a request to change the email address associated with your Sourcetool account. To ensure the security of your account, we need you to verify your new email address.

Please click the following link within the next 24 hours to confirm your email change:
%s

Thank you for being a part of the Sourcetool community!
Regards,

Sourcetool Team`,
		in.FirstName,
		in.URL,
	)

	msg := fmt.Sprintf("From: Sourcetool Team <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", m.emailClient.GetFromAddress(), in.To, subject, content)

	return emailutil.SendWithLogging(ctx, msg, func() error {
		if err := m.emailClient.SendMail([]string{in.To}, []byte(msg)); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
		return nil
	})
}

func (m *MailerCE) SendInvitationEmail(ctx context.Context, in *model.SendInvitationEmail) error {
	subject := "You've been invited to join Sourcetool!"
	baseContent := `You've been invited to join Sourcetool!

To accept the invitation, please create your account by clicking the URL below within 24 hours.

%s

- This URL will expire in 24 hours.
- This is a send-only email address.
- Your account will not be created unless you complete the next steps.`

	sendEmails := make([]string, 0)
	for email, url := range in.URLs {
		if lo.Contains(sendEmails, email) {
			continue
		}

		content := fmt.Sprintf(baseContent, url)
		msg := fmt.Sprintf("From: Sourcetool Team <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n", m.emailClient.GetFromAddress(), email, subject, content)

		err := emailutil.SendWithLogging(ctx, msg, func() error {
			if err := m.emailClient.SendMail([]string{email}, []byte(msg)); err != nil {
				return fmt.Errorf("failed to send email: %w", err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		sendEmails = append(sendEmails, email)
	}

	return nil
}

func (m *MailerCE) SendMultipleOrganizationsEmail(ctx context.Context, in *model.SendMultipleOrganizationsEmail) error {
	subject := "Choose your Sourcetool organization to log in"

	urlList := ""
	for _, url := range in.LoginURLs {
		urlList += url + "\n"
	}

	content := fmt.Sprintf(`Hi %s,

Your email, %s, is associated with multiple Sourcetool organizations. You may log in to each one by visiting its login page below:

%s
If you have any questions, encounter any issues, or need further assistance, please don't hesitate to reach out to support@sourcetool.com.

Thank you!

The Sourcetool Team`, in.FirstName, in.Email, urlList)

	msg := fmt.Sprintf("From: Sourcetool Team <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", m.emailClient.GetFromAddress(), in.To, subject, content)

	return emailutil.SendWithLogging(ctx, msg, func() error {
		if err := m.emailClient.SendMail([]string{in.To}, []byte(msg)); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
		return nil
	})
}

func (m *MailerCE) SendMagicLinkEmail(ctx context.Context, in *model.SendMagicLinkEmail) error {
	subject := "Log in to your Sourcetool account"
	content := fmt.Sprintf(`Hi %s,

Here's your magic link to sign in to your Sourcetool account. Click the link below to access your account securely without a password:

%s

- This link will expire in 15 minutes for security reasons.
- If you didn't request this link, you can safely ignore this email.

Thank you for using Sourcetool!

The Sourcetool Team`, in.FirstName, in.URL)

	msg := fmt.Sprintf("From: Sourcetool Team <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", m.emailClient.GetFromAddress(), in.To, subject, content)

	return emailutil.SendWithLogging(ctx, msg, func() error {
		if err := m.emailClient.SendMail([]string{in.To}, []byte(msg)); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
		return nil
	})
}

func (m *MailerCE) SendMultipleOrganizationsMagicLinkEmail(ctx context.Context, in *model.SendMultipleOrganizationsMagicLinkEmail) error {
	subject := "Log in to your Sourcetool organizations"

	urlList := ""
	for _, url := range in.LoginURLs {
		urlList += url + "\n"
	}

	content := fmt.Sprintf(`Hi %s,

Your email, %s, is associated with multiple Sourcetool organizations. You may sign in to each one by clicking its magic link below:

%s

- These links will expire in 15 minutes for security reasons.
- If you didn't request these links, you can safely ignore this email.

Thank you for using Sourcetool!

The Sourcetool Team`, in.FirstName, in.Email, urlList)

	msg := fmt.Sprintf("From: Sourcetool Team <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", m.emailClient.GetFromAddress(), in.To, subject, content)

	return emailutil.SendWithLogging(ctx, msg, func() error {
		if err := m.emailClient.SendMail([]string{in.To}, []byte(msg)); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
		return nil
	})
}

func (m *MailerCE) SendInvitationMagicLinkEmail(ctx context.Context, in *model.SendInvitationMagicLinkEmail) error {
	subject := "Your invitation to join Sourcetool"
	content := fmt.Sprintf(`Hi %s,

You've been invited to join Sourcetool. Click the link below to accept the invitation:

%s

This link will expire in 15 minutes.

Best regards,
The Sourcetool Team`, in.FirstName, in.URL)

	msg := fmt.Sprintf("From: Sourcetool Team <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", m.emailClient.GetFromAddress(), in.To, subject, content)

	return emailutil.SendWithLogging(ctx, msg, func() error {
		if err := m.emailClient.SendMail([]string{in.To}, []byte(msg)); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
		return nil
	})
}
