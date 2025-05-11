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
		{UserOrganizationRole(100), "unknown"}, // 範囲外
	}
	for _, c := range cases {
		if got := c.role.String(); got != c.expected {
			t.Errorf("UserOrganizationRole(%d).String() = %q, want %q", c.role, got, c.expected)
		}
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
		{"invalid", UserOrganizationRoleUnknown}, // 未知の文字列
	}
	for _, c := range cases {
		if got := UserOrganizationRoleFromString(c.input); got != c.expected {
			t.Errorf("UserOrganizationRoleFromString(%q) = %v, want %v", c.input, got, c.expected)
		}
	}
}
