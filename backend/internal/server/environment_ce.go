//go:build !ee
// +build !ee

package server

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

const maxEnvironmentsInCE = 2

func (s *Server) canCreateEnvironment(ctx context.Context, organizationID uuid.UUID) error {
	environments, err := s.db.Environment().List(ctx, database.EnvironmentByOrganizationID(organizationID))
	if err != nil {
		return err
	}

	if len(environments) >= maxEnvironmentsInCE {
		return errdefs.ErrEnvironmentLimitReached(
			fmt.Errorf("CE version is limited to %d environments", maxEnvironmentsInCE),
		)
	}

	return nil
}
