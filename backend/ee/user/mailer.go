package user

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/trysourcetool/sourcetool/backend/ee/config"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type mailerEE struct {
	auth    smtp.Auth
	addr    string
	from    string
	tlsConf *tls.Config
	host    string
	*user.MailerCE
}

func NewMailerEE() *mailerEE {
	cfg := config.Config.SMTP
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	tlsConf := &tls.Config{
		ServerName: cfg.Host,
		MinVersion: tls.VersionTLS12,
	}
	return &mailerEE{
		auth:     auth,
		addr:     addr,
		from:     cfg.FromEmail,
		tlsConf:  tlsConf,
		host:     cfg.Host,
		MailerCE: user.NewMailerCE(),
	}
}
