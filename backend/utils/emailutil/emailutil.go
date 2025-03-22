package emailutil

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/logger"
)

// EmailClient represents an email client with SMTP configuration.
type EmailClient struct {
	auth    smtp.Auth
	addr    string
	from    string
	tlsConf *tls.Config
	host    string
}

// NewEmailClient creates a new email client using SMTP configuration.
func NewEmailClient() *EmailClient {
	cfg := config.Config.SMTP
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	tlsConf := &tls.Config{
		ServerName: cfg.Host,
		MinVersion: tls.VersionTLS12,
	}
	return &EmailClient{
		auth:    auth,
		addr:    addr,
		from:    cfg.FromEmail,
		tlsConf: tlsConf,
		host:    cfg.Host,
	}
}

// SendMail sends an email through SMTP with TLS.
func (c *EmailClient) SendMail(to []string, msg []byte) error {
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

	if err = client.Mail(c.from); err != nil {
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

// GetFromAddress returns the sender email address.
func (c *EmailClient) GetFromAddress() string {
	return c.from
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
func SendWithLogging(ctx context.Context, content string, sendFunc func() error) error {
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
