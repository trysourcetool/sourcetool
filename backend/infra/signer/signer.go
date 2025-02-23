package signer

import (
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type signerCE struct{}

func NewCE() *signerCE {
	return &signerCE{}
}

func (s *signerCE) User() model.UserSigner {
	return user.NewSignerCE()
}
