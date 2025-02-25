package user

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/model"
)

type MailerCE struct {
	auth    smtp.Auth
	addr    string
	from    string
	tlsConf *tls.Config
	host    string
}

func NewMailerCE() *MailerCE {
	cfg := config.Config.SMTP
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	tlsConf := &tls.Config{
		ServerName: cfg.Host,
		MinVersion: tls.VersionTLS12,
	}
	return &MailerCE{
		auth:    auth,
		addr:    addr,
		from:    cfg.FromEmail,
		tlsConf: tlsConf,
		host:    cfg.Host,
	}
}

func (m *MailerCE) SendSignUpInstructions(ctx context.Context, in *model.SendSignUpInstructions) error {
	subject := "Activate your Sourcetool account"
	content := fmt.Sprintf(`Welcome to Sourcetool!

To complete your registration for our service, please create your account by clicking the URL below within 24 hours.

%s 

- This URL will expire in 24 hours.
- This is a send-only email address. 
- Your account will not be created unless you complete the next steps.`, in.URL)

	msg := fmt.Sprintf("From: Sourcetool Team <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", m.from, in.To, subject, content)

	if err := m.sendMail([]string{in.To}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
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
		"%s\r\n", m.from, in.To, subject, content)

	if err := m.sendMail([]string{in.To}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
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
			"%s\r\n", m.from, email, subject, content)

		if err := m.sendMail([]string{email}, []byte(msg)); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}

		sendEmails = append(sendEmails, email)
	}

	return nil
}

func (m *MailerCE) sendMail(to []string, msg []byte) error {
	conn, err := tls.Dial("tcp", m.addr, m.tlsConf)
	if err != nil {
		return fmt.Errorf("failed to create TLS connection: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err = client.Auth(m.auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err = client.Mail(m.from); err != nil {
		return fmt.Errorf("failed to set FROM address: %w", err)
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set TO address: %w", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to create message writer: %w", err)
	}

	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close message writer: %w", err)
	}

	return client.Quit()
}
