package core

import "testing"

func TestUserOrganizationRole_String(t *testing.T) {
	cases := []struct {
		role     UserOrganizationRole
		expected string
	}{
		{UserOrganizationRoleUnknown, "unknown"},
		{UserOrganizationRoleAdmin, "admin"},
		{UserOrganizationRoleDeveloper, "developer"},
		{UserOrganizationRoleMember, "member"},
		{UserOrganizationRole(100), "unknown"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.expected, func(t *testing.T) {
			t.Parallel()
			if got := c.role.String(); got != c.expected {
				t.Errorf("UserOrganizationRole(%d).String() = %q, want %q", c.role, got, c.expected)
			}
		})
	}
}

func TestUserOrganizationRoleFromString(t *testing.T) {
	cases := []struct {
		input    string
		expected UserOrganizationRole
	}{
		{"unknown", UserOrganizationRoleUnknown},
		{"admin", UserOrganizationRoleAdmin},
		{"developer", UserOrganizationRoleDeveloper},
		{"member", UserOrganizationRoleMember},
		{"invalid", UserOrganizationRoleUnknown},
	}
	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			t.Parallel()
			if got := UserOrganizationRoleFromString(c.input); got != c.expected {
				t.Errorf("UserOrganizationRoleFromString(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}
