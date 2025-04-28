package core

import "testing"

func TestCanPerform(t *testing.T) {
	tests := []struct {
		name string
		role UserOrganizationRole
		op   Operation
		want bool
	}{
		// Admin can do everything
		{"Admin can edit org", UserOrganizationRoleAdmin, OperationEditOrganization, true},
		{"Admin can edit billing", UserOrganizationRoleAdmin, OperationEditBilling, true},
		{"Admin can edit live key", UserOrganizationRoleAdmin, OperationEditLiveModeAPIKey, true},
		{"Admin can edit dev key", UserOrganizationRoleAdmin, OperationEditDevModeAPIKey, true},
		{"Admin can edit env", UserOrganizationRoleAdmin, OperationEditEnvironment, true},
		{"Admin can edit group", UserOrganizationRoleAdmin, OperationEditGroup, true},
		{"Admin can edit user", UserOrganizationRoleAdmin, OperationEditUser, true},

		// Developer can only do some
		{"Dev cannot edit org", UserOrganizationRoleDeveloper, OperationEditOrganization, false},
		{"Dev cannot edit billing", UserOrganizationRoleDeveloper, OperationEditBilling, false},
		{"Dev cannot edit live key", UserOrganizationRoleDeveloper, OperationEditLiveModeAPIKey, false},
		{"Dev can edit dev key", UserOrganizationRoleDeveloper, OperationEditDevModeAPIKey, true},
		{"Dev cannot edit env", UserOrganizationRoleDeveloper, OperationEditEnvironment, false},
		{"Dev can edit group", UserOrganizationRoleDeveloper, OperationEditGroup, true},
		{"Dev cannot edit user", UserOrganizationRoleDeveloper, OperationEditUser, false},

		// Member can do nothing
		{"Member cannot edit org", UserOrganizationRoleMember, OperationEditOrganization, false},
		{"Member cannot edit billing", UserOrganizationRoleMember, OperationEditBilling, false},
		{"Member cannot edit live key", UserOrganizationRoleMember, OperationEditLiveModeAPIKey, false},
		{"Member cannot edit dev key", UserOrganizationRoleMember, OperationEditDevModeAPIKey, false},
		{"Member cannot edit env", UserOrganizationRoleMember, OperationEditEnvironment, false},
		{"Member cannot edit group", UserOrganizationRoleMember, OperationEditGroup, false},
		{"Member cannot edit user", UserOrganizationRoleMember, OperationEditUser, false},

		// Unknown operation
		{"Unknown op returns false", UserOrganizationRoleAdmin, Operation("UNKNOWN_OP"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanPerform(tt.role, tt.op)
			if got != tt.want {
				t.Errorf("CanPerform(%q, %q) = %v, want %v", tt.role, tt.op, got, tt.want)
			}
		})
	}
}
