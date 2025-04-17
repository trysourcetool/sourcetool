package permission

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/user"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

type operation string

const (
	OperationEditOrganization   = operation("EDIT_ORGANIZATION")
	OperationEditBilling        = operation("EDIT_BILLING")
	OperationEditLiveModeAPIKey = operation("EDIT_LIVE_MODE_API_KEY")
	OperationEditDevModeAPIKey  = operation("EDIT_DEV_MODE_API_KEY")
	OperationEditEnvironment    = operation("EDIT_ENVIRONMENT")
	OperationEditGroup          = operation("EDIT_GROUP")
	OperationEditUser           = operation("EDIT_USER")
)

var rolesAllowedByOperation = map[operation][]user.UserOrganizationRole{
	OperationEditOrganization: {
		user.UserOrganizationRoleAdmin,
	},
	OperationEditBilling: {
		user.UserOrganizationRoleAdmin,
	},
	OperationEditLiveModeAPIKey: {
		user.UserOrganizationRoleAdmin,
	},
	OperationEditDevModeAPIKey: {
		user.UserOrganizationRoleAdmin,
		user.UserOrganizationRoleDeveloper,
	},
	OperationEditEnvironment: {
		user.UserOrganizationRoleAdmin,
	},
	OperationEditGroup: {
		user.UserOrganizationRoleAdmin,
	},
	OperationEditUser: {
		user.UserOrganizationRoleAdmin,
	},
}

func (c *Checker) AuthorizeOperation(ctx context.Context, o operation) error {
	currentUser := ctxutil.CurrentUser(ctx)
	currentOrg := ctxutil.CurrentOrganization(ctx)
	if currentUser == nil || currentOrg == nil {
		return errdefs.ErrPermissionDenied(errors.New("user or organization context not found"))
	}

	orgAccess, err := c.store.User().GetOrganizationAccess(ctx, user.OrganizationAccessByUserID(currentUser.ID), user.OrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil && !errdefs.IsUserOrganizationAccessNotFound(err) {
		return errdefs.ErrPermissionDenied(err)
	}

	if !slices.Contains(rolesAllowedByOperation[o], orgAccess.Role) {
		return errdefs.ErrPermissionDenied(fmt.Errorf("user role %s is not allowed to perform operation: %s", orgAccess.Role.String(), o))
	}

	return nil
}
