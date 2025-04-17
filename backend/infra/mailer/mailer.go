package mailer

import (
	"github.com/trysourcetool/sourcetool/backend/user"
	userMailer "github.com/trysourcetool/sourcetool/backend/user/mailer"
)

type mailerCE struct{}

func NewCE() *mailerCE {
	return &mailerCE{}
}

func (m *mailerCE) User() user.Mailer {
	return userMailer.NewUserMailerCE()
}
