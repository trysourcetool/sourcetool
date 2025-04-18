package mailer

import (
	"github.com/trysourcetool/sourcetool/backend/ee/user/mailer"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type mailerEE struct{}

func NewMailerEE() *mailerEE {
	return &mailerEE{}
}

func (m *mailerEE) User() user.Mailer {
	return mailer.NewUserMailerEE()
}
