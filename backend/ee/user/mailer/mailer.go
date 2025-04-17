package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/user/mailer"
)

type userMailerEE struct {
	auth    smtp.Auth
	addr    string
	from    string
	tlsConf *tls.Config
	host    string
	*mailer.UserMailerCE
}

func NewUserMailerEE() *userMailerEE {
	cfg := config.Config.SMTP
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	tlsConf := &tls.Config{
		ServerName: cfg.Host,
		MinVersion: tls.VersionTLS12,
	}
	return &userMailerEE{
		auth:         auth,
		addr:         addr,
		from:         cfg.FromEmail,
		tlsConf:      tlsConf,
		host:         cfg.Host,
		UserMailerCE: mailer.NewUserMailerCE(),
	}
}
