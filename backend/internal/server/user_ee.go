//go:build ee
// +build ee

package server

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

func (s *Server) canAddUserToOrganization(ctx context.Context, organizationID uuid.UUID) error {
	// No limit in EE version
	return nil
}
