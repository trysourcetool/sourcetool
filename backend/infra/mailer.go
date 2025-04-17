package infra

import "github.com/trysourcetool/sourcetool/backend/user"

type Mailer interface {
	User() user.Mailer
}
