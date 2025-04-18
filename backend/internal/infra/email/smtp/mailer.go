package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/logger"
)

type mailerCE struct {
	auth      smtp.Auth
	addr      string
	fromEmail string
	tlsConf   *tls.Config
	host      string
}

func NewMailerCE() *mailerCE {
	cfg := config.Config.SMTP
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	tlsConf := &tls.Config{
		ServerName: cfg.Host,
		MinVersion: tls.VersionTLS12,
	}
	return &mailerCE{
		auth:      auth,
		addr:      addr,
		fromEmail: cfg.FromEmail,
		tlsConf:   tlsConf,
		host:      cfg.Host,
	}
}

func (c *mailerCE) Send(ctx context.Context, to []string, from, subject, body string) error {
	conn, err := tls.Dial("tcp", c.addr, c.tlsConf)
	if err != nil {
		return fmt.Errorf("failed to create TLS connection: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, c.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err = client.Auth(c.auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err = client.Mail(c.fromEmail); err != nil {
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

	msg := fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", from, c.fromEmail, strings.Join(to, ","), subject, body)

	if _, err = w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close message writer: %w", err)
	}

	return sendWithLogging(ctx, msg, func() error {
		return client.Quit()
	})
}

// SendWithLogging sends an email in production environments and logs the content
// in local development environment. Useful for debugging emails without actually
// sending them in the local environment.
//
// Parameters:
//   - ctx: The context
//   - content: The email content as string, used for local environment display
//   - sendFunc: A function that performs the actual email sending
//
// Returns:
//   - error: Any error that occurred during the process
func sendWithLogging(ctx context.Context, content string, sendFunc func() error) error {
	if config.Config.Env == config.EnvLocal {
		// In local environment, just log the email content
		logger.Logger.Sugar().Debug("================= EMAIL CONTENT =================")
		logger.Logger.Sugar().Debug(content)
		logger.Logger.Sugar().Debug("================= EMAIL CONTENT =================")

		// Don't actually send in local environment
		return nil
	}

	// In non-local environments, perform the email sending
	if err := sendFunc(); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
