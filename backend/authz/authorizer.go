package authz

import "github.com/trysourcetool/sourcetool/backend/infra"

type authorizer struct {
	store infra.ModelStore
}

func NewAuthorizer(store infra.ModelStore) *authorizer {
	return &authorizer{store}
}
