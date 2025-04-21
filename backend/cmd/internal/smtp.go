package internal

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
)

func OpenSMTP() (*smtp.Client, error) {
	cfg := config.Config.SMTP
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	tlsConf := &tls.Config{
		ServerName: cfg.Host,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", addr, tlsConf)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS connection: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to create SMTP client: %w", err)
	}

	if err = client.Auth(auth); err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	return client, nil
}
