package infra

import "github.com/trysourcetool/sourcetool/backend/model"

type Mailer interface {
	User() model.UserMailerCE
}
