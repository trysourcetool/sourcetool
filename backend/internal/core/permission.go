package core

import "slices"

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

var rolesAllowedByOperation = map[Operation][]UserOrganizationRole{
	OperationEditOrganization:   {UserOrganizationRoleAdmin},
	OperationEditBilling:        {UserOrganizationRoleAdmin},
	OperationEditLiveModeAPIKey: {UserOrganizationRoleAdmin},
	OperationEditDevModeAPIKey:  {UserOrganizationRoleAdmin, UserOrganizationRoleDeveloper},
	OperationEditEnvironment:    {UserOrganizationRoleAdmin},
	OperationEditGroup:          {UserOrganizationRoleAdmin},
	OperationEditUser:           {UserOrganizationRoleAdmin},
}

func CanPerform(role UserOrganizationRole, op Operation) bool {
	allowed, ok := rolesAllowedByOperation[op]
	if !ok {
		return false
	}
	return slices.Contains(allowed, role)
}
