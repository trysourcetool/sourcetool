package user

import "github.com/trysourcetool/sourcetool/backend/user"

type signerEE struct {
	*user.SignerCE
}

func NewSignerEE() *signerEE {
	return &signerEE{
		SignerCE: user.NewSignerCE(),
	}
}
