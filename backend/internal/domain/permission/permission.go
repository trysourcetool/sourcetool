package permission

import (
	"golang.org/x/exp/slices"

	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
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

var rolesAllowedByOperation = map[Operation][]user.UserOrganizationRole{
	OperationEditOrganization:   {user.UserOrganizationRoleAdmin},
	OperationEditBilling:        {user.UserOrganizationRoleAdmin},
	OperationEditLiveModeAPIKey: {user.UserOrganizationRoleAdmin},
	OperationEditDevModeAPIKey:  {user.UserOrganizationRoleAdmin, user.UserOrganizationRoleDeveloper},
	OperationEditEnvironment:    {user.UserOrganizationRoleAdmin},
	OperationEditGroup:          {user.UserOrganizationRoleAdmin},
	OperationEditUser:           {user.UserOrganizationRoleAdmin},
}

func CanPerform(role user.UserOrganizationRole, op Operation) bool {
	allowed, ok := rolesAllowedByOperation[op]
	if !ok {
		return false
	}
	return slices.Contains(allowed, role)
}
