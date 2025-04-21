package permission

import (
	"context"
	"errors"
	"fmt"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

type Checker struct {
	db database.DB
}

func NewChecker(db database.DB) *Checker {
	return &Checker{db: db}
}

func (c *Checker) AuthorizeOperation(ctx context.Context, op core.Operation) error {
	currentUser := internal.CurrentUser(ctx)
	currentOrg := internal.CurrentOrganization(ctx)
	if currentUser == nil || currentOrg == nil {
		return errdefs.ErrPermissionDenied(errors.New("user or organization context not found"))
	}
	orgAccess, err := c.db.User().GetOrganizationAccess(
		ctx,
		database.UserOrganizationAccessByUserID(currentUser.ID),
		database.UserOrganizationAccessByOrganizationID(currentOrg.ID),
	)
	if err != nil && !errdefs.IsUserOrganizationAccessNotFound(err) {
		return errdefs.ErrPermissionDenied(err)
	}
	if !core.CanPerform(orgAccess.Role, op) {
		return errdefs.ErrPermissionDenied(fmt.Errorf("user role %s is not allowed to perform operation: %s", orgAccess.Role.String(), op))
	}
	return nil
}
