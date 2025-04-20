package permission

import (
	"context"
	"errors"
	"fmt"

	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	"github.com/trysourcetool/sourcetool/backend/internal/ctxutil"
	domainperm "github.com/trysourcetool/sourcetool/backend/internal/domain/permission"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
)

// Checker handles authorization logic in the application layer.
type Checker struct {
	repo port.Repositories
}

// NewChecker creates a new permission Checker.
func NewChecker(repo port.Repositories) *Checker {
	return &Checker{repo: repo}
}

// AuthorizeOperation checks if the current user can perform the given operation.
func (c *Checker) AuthorizeOperation(ctx context.Context, op domainperm.Operation) error {
	currentUser := ctxutil.CurrentUser(ctx)
	currentOrg := ctxutil.CurrentOrganization(ctx)
	if currentUser == nil || currentOrg == nil {
		return errdefs.ErrPermissionDenied(errors.New("user or organization context not found"))
	}
	orgAccess, err := c.repo.User().GetOrganizationAccess(
		ctx,
		user.OrganizationAccessByUserID(currentUser.ID),
		user.OrganizationAccessByOrganizationID(currentOrg.ID),
	)
	if err != nil && !errdefs.IsUserOrganizationAccessNotFound(err) {
		return errdefs.ErrPermissionDenied(err)
	}
	if !domainperm.CanPerform(orgAccess.Role, op) {
		return errdefs.ErrPermissionDenied(fmt.Errorf("user role %s is not allowed to perform operation: %s", orgAccess.Role.String(), op))
	}
	return nil
}
