package signer

import (
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type SignerCE struct{}

func NewCE() *SignerCE {
	return &SignerCE{}
}

func (s *SignerCE) User() model.UserSignerCE {
	return user.NewSignerCE()
}
