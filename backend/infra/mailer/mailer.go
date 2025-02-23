package mailer

import (
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type mailerCE struct{}

func NewCE() *mailerCE {
	return &mailerCE{}
}

func (m *mailerCE) User() model.UserMailer {
	return user.NewMailerCE()
}
