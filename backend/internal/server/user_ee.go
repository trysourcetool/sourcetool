package server

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

func (s *Server) canAddUserToOrganization(ctx context.Context, organizationID uuid.UUID) error {
	return nil
}
