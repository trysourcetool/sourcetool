package permission

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/postgres"
)

type Operation string

const (
	OperationEditOrganization   Operation = "EDIT_ORGANIZATION"
	OperationEditBilling        Operation = "EDIT_BILLING"
	OperationEditLiveModeAPIKey Operation = "EDIT_LIVE_MODE_API_KEY" // #nosec G101
	OperationEditDevModeAPIKey  Operation = "EDIT_DEV_MODE_API_KEY"  // #nosec G101
	OperationEditEnvironment    Operation = "EDIT_ENVIRONMENT"
	OperationEditGroup          Operation = "EDIT_GROUP"
	OperationEditUser           Operation = "EDIT_USER"
)

var rolesAllowedByOperation = map[Operation][]core.UserOrganizationRole{
	OperationEditOrganization:   {core.UserOrganizationRoleAdmin},
	OperationEditBilling:        {core.UserOrganizationRoleAdmin},
	OperationEditLiveModeAPIKey: {core.UserOrganizationRoleAdmin},
	OperationEditDevModeAPIKey:  {core.UserOrganizationRoleAdmin, core.UserOrganizationRoleDeveloper},
	OperationEditEnvironment:    {core.UserOrganizationRoleAdmin},
	OperationEditGroup:          {core.UserOrganizationRoleAdmin},
	OperationEditUser:           {core.UserOrganizationRoleAdmin},
}

func canPerform(role core.UserOrganizationRole, op Operation) bool {
	allowed, ok := rolesAllowedByOperation[op]
	if !ok {
		return false
	}
	return slices.Contains(allowed, role)
}

type Checker struct {
	db *postgres.DB
}

func NewChecker(db *postgres.DB) *Checker {
	return &Checker{db: db}
}

func (c *Checker) AuthorizeOperation(ctx context.Context, op Operation) error {
	currentUser := internal.CurrentUser(ctx)
	currentOrg := internal.CurrentOrganization(ctx)
	if currentUser == nil || currentOrg == nil {
		return errdefs.ErrPermissionDenied(errors.New("user or organization context not found"))
	}
	orgAccess, err := c.db.GetUserOrganizationAccess(
		ctx,
		postgres.UserOrganizationAccessByUserID(currentUser.ID),
		postgres.UserOrganizationAccessByOrganizationID(currentOrg.ID),
	)
	if err != nil && !errdefs.IsUserOrganizationAccessNotFound(err) {
		return errdefs.ErrPermissionDenied(err)
	}
	if !canPerform(orgAccess.Role, op) {
		return errdefs.ErrPermissionDenied(fmt.Errorf("user role %s is not allowed to perform operation: %s", orgAccess.Role.String(), op))
	}
	return nil
}
