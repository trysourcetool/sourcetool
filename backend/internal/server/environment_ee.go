//go:build ee
// +build ee

package server

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

func (s *Server) canCreateEnvironment(ctx context.Context, organizationID uuid.UUID) error {
	// No environment limit in EE version
	return nil
}
