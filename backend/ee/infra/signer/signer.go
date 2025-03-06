package signer

import (
	"github.com/trysourcetool/sourcetool/backend/ee/user"
	"github.com/trysourcetool/sourcetool/backend/model"
)

type signerEE struct{}

func NewEE() *signerEE {
	return &signerEE{}
}

func (s *signerEE) User() model.UserSigner {
	return user.NewSignerEE()
}
