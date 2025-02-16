package infra

import "github.com/trysourcetool/sourcetool/backend/model"

type Signer interface {
	User() model.UserSignerCE
}
