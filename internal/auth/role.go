package auth

import (
	"fmt"
	"strconv"
)

// Role ...
type Role uint8

// Roles ...
const (
	RoleAdmin Role = iota + 1
	RoleAuditor
)

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "Admin"
	case RoleAuditor:
		return "Auditor"
	default:
		return "Unknown"
	}
}

var roles = map[string]Role{
	// RoleAdmin
	strconv.Itoa(int(RoleAdmin)): RoleAdmin,
	fmt.Sprintf("%s", RoleAdmin): RoleAdmin,
	// RoleAuditor
	strconv.Itoa(int(RoleAuditor)): RoleAuditor,
	fmt.Sprintf("%s", RoleAuditor): RoleAuditor,
}

// ParseRole ...
func ParseRole(roleStr string) (Role, error) {
	role, ok := roles[roleStr]
	if ok {
		return role, nil
	}
	return 0, fmt.Errorf("unknown role %s", roleStr)
}
