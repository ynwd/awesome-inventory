package domain

import (
	"errors"
	"inv/pkg/domain"
)

type Role string

func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleEditor, RoleViewer:
		return true
	}
	return false
}

const (
	RoleOwner  Role = "OWNER"
	RoleEditor Role = "EDITOR"
	RoleViewer Role = "VIEWER"
)

type Permission string

const (
	PermissionCreate Permission = "CREATE"
	PermissionRead   Permission = "READ"
	PermissionUpdate Permission = "UPDATE"
	PermissionDelete Permission = "DELETE"
	PermissionShare  Permission = "SHARE"
)

var (
	ErrInvalidRole       = errors.New("invalid role")
	ErrInvalidPermission = errors.New("invalid permission")
)

var rolePermissions = map[Role][]Permission{
	RoleOwner: {
		PermissionCreate,
		PermissionRead,
		PermissionUpdate,
		PermissionDelete,
		PermissionShare,
	},
	RoleEditor: {
		PermissionCreate, // Added create permission
		PermissionRead,
		PermissionUpdate,
	},
	RoleViewer: {
		PermissionRead,
	},
}

type User struct {
	domain.BaseModel
	Email     string `gorm:"type:varchar(255);unique;not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	FirstName string `gorm:"type:varchar(100)"`
	LastName  string `gorm:"type:varchar(100)"`
	IsActive  bool   `gorm:"default:true"`
	Role      Role   `gorm:"type:varchar(20);not null;default:'VIEWER'"`
}

func (u *User) HasPermission(p Permission) bool {
	permissions, exists := rolePermissions[u.Role]
	if !exists {
		return false
	}

	for _, perm := range permissions {
		if perm == p {
			return true
		}
	}
	return false
}

func (u *User) CanShare() bool {
	return u.HasPermission(PermissionShare)
}

func (u *User) CanEdit() bool {
	return u.HasPermission(PermissionUpdate)
}

func (u *User) CanDelete() bool {
	return u.HasPermission(PermissionDelete)
}
