package server

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

const maxUsersInCE = 5

func (s *Server) canAddUserToOrganization(ctx context.Context, organizationID uuid.UUID) error {
	users, err := s.db.User().List(ctx, database.UserByOrganizationID(organizationID))
	if err != nil {
		return err
	}

	invitations, err := s.db.User().ListInvitations(ctx, database.UserInvitationByOrganizationID(organizationID))
	if err != nil {
		return err
	}

	totalUserCount := len(users) + len(invitations)

	if totalUserCount >= maxUsersInCE {
		return errdefs.ErrUserLimitReached(
			fmt.Errorf("CE version is limited to %d users", maxUsersInCE),
		)
	}

	return nil
}
