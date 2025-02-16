package mailer

import (
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type MailerCE struct{}

func NewCE() *MailerCE {
	return &MailerCE{}
}

func (m *MailerCE) User() model.UserMailerCE {
	return user.NewMailerCE()
}
