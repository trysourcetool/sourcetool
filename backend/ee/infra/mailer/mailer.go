package mailer

import (
	"github.com/trysourcetool/sourcetool/backend/ee/user"
	"github.com/trysourcetool/sourcetool/backend/model"
)

type mailerEE struct{}

func NewEE() *mailerEE {
	return &mailerEE{}
}

func (m *mailerEE) User() model.UserMailer {
	return user.NewMailerEE()
}
