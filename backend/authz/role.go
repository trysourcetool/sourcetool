package authz

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
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

var rolesAllowedByOperation = map[operation][]model.UserOrganizationRole{
	OperationEditOrganization: {
		model.UserOrganizationRoleAdmin,
	},
	OperationEditBilling: {
		model.UserOrganizationRoleAdmin,
	},
	OperationEditLiveModeAPIKey: {
		model.UserOrganizationRoleAdmin,
	},
	OperationEditDevModeAPIKey: {
		model.UserOrganizationRoleAdmin,
		model.UserOrganizationRoleDeveloper,
	},
	OperationEditEnvironment: {
		model.UserOrganizationRoleAdmin,
	},
	OperationEditGroup: {
		model.UserOrganizationRoleAdmin,
	},
	OperationEditUser: {
		model.UserOrganizationRoleAdmin,
	},
}

func (a *authorizer) AuthorizeOperation(ctx context.Context, o operation) error {
	currentUser := ctxutil.CurrentUser(ctx)
	currentOrg := ctxutil.CurrentOrganization(ctx)
	if currentUser == nil || currentOrg == nil {
		return errdefs.ErrPermissionDenied(errors.New("user or organization context not found"))
	}

	orgAccess, err := a.store.User().GetOrganizationAccess(ctx, storeopts.UserOrganizationAccessByUserID(currentUser.ID), storeopts.UserOrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil && !errdefs.IsUserOrganizationAccessNotFound(err) {
		return errdefs.ErrPermissionDenied(err)
	}

	if !slices.Contains(rolesAllowedByOperation[o], orgAccess.Role) {
		return errdefs.ErrPermissionDenied(fmt.Errorf("user role %s is not allowed to perform operation: %s", orgAccess.Role.String(), o))
	}

	return nil
}
