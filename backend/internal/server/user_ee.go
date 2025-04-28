//go:build ee
// +build ee

package server

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

func (s *Server) canAddUserToOrganization(ctx context.Context, organizationID uuid.UUID) error {
	return s.canAddUsersToOrganization(ctx, organizationID, 1)
}

func (s *Server) canAddUsersToOrganization(ctx context.Context, organizationID uuid.UUID, newUserCount int) error {
	// No limit in EE version
	return nil
}
