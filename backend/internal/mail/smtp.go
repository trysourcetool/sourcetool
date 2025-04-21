package mail

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
)

type Mail struct {
	smtpClient *smtp.Client
}

func New(smtpClient *smtp.Client) *Mail {
	return &Mail{smtpClient}
}

type MailInput struct {
	To      []string
	From    string
	Subject string
	Body    string
}

func (m *Mail) Send(ctx context.Context, input MailInput) error {
	msg := fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", input.From, input.From, strings.Join(input.To, ","), input.Subject, input.Body)

	if config.Config.Env == config.EnvLocal {
		// In local environment, just log the email content
		logger.Logger.Sugar().Debug("================= EMAIL CONTENT =================")
		logger.Logger.Sugar().Debug(msg)
		logger.Logger.Sugar().Debug("================= EMAIL CONTENT =================")

		// Don't actually send in local environment
		return nil
	}

	if err := m.smtpClient.Mail(input.From); err != nil {
		return fmt.Errorf("failed to set FROM address: %w", err)
	}

	for _, addr := range input.To {
		if err := m.smtpClient.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set TO address: %w", err)
		}
	}

	w, err := m.smtpClient.Data()
	if err != nil {
		return fmt.Errorf("failed to create message writer: %w", err)
	}

	if _, err = w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close message writer: %w", err)
	}

	return m.smtpClient.Quit()
}
