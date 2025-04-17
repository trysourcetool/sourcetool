package user

type UserOrganizationRole int

const (
	UserOrganizationRoleUnknown UserOrganizationRole = iota
	UserOrganizationRoleAdmin
	UserOrganizationRoleDeveloper
	UserOrganizationRoleMember

	userOrganizationRoleUnknown   = "unknown"
	userOrganizationRoleAdmin     = "admin"
	userOrganizationRoleDeveloper = "developer"
	userOrganizationRoleMember    = "member"
)

func (r UserOrganizationRole) String() string {
	switch r {
	case UserOrganizationRoleAdmin:
		return userOrganizationRoleAdmin
	case UserOrganizationRoleDeveloper:
		return userOrganizationRoleDeveloper
	case UserOrganizationRoleMember:
		return userOrganizationRoleMember
	}
	return userOrganizationRoleUnknown
}

func UserOrganizationRoleFromString(s string) UserOrganizationRole {
	switch s {
	case userOrganizationRoleAdmin:
		return UserOrganizationRoleAdmin
	case userOrganizationRoleDeveloper:
		return UserOrganizationRoleDeveloper
	case userOrganizationRoleMember:
		return UserOrganizationRoleMember
	}
	return UserOrganizationRoleUnknown
}
