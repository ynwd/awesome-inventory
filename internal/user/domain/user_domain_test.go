package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserPermissions(t *testing.T) {
	t.Run("owner should have all permissions", func(t *testing.T) {
		owner := &User{Role: RoleOwner}

		assert.True(t, owner.HasPermission(PermissionCreate))
		assert.True(t, owner.HasPermission(PermissionRead))
		assert.True(t, owner.HasPermission(PermissionUpdate))
		assert.True(t, owner.HasPermission(PermissionDelete))
		assert.True(t, owner.HasPermission(PermissionShare))
		assert.True(t, owner.CanShare())
		assert.True(t, owner.CanEdit())
		assert.True(t, owner.CanDelete())
	})

	t.Run("editor should have create, read and update permissions", func(t *testing.T) {
		editor := &User{Role: RoleEditor}

		assert.True(t, editor.HasPermission(PermissionCreate))
		assert.True(t, editor.HasPermission(PermissionRead))
		assert.True(t, editor.HasPermission(PermissionUpdate))
		assert.False(t, editor.HasPermission(PermissionDelete))
		assert.False(t, editor.HasPermission(PermissionShare))
		assert.False(t, editor.CanShare())
		assert.True(t, editor.CanEdit())
		assert.False(t, editor.CanDelete())
	})

	t.Run("viewer should only have read permission", func(t *testing.T) {
		viewer := &User{Role: RoleViewer}

		assert.False(t, viewer.HasPermission(PermissionCreate))
		assert.True(t, viewer.HasPermission(PermissionRead))
		assert.False(t, viewer.HasPermission(PermissionUpdate))
		assert.False(t, viewer.HasPermission(PermissionDelete))
		assert.False(t, viewer.HasPermission(PermissionShare))
		assert.False(t, viewer.CanShare())
		assert.False(t, viewer.CanEdit())
		assert.False(t, viewer.CanDelete())
	})
}

func TestRoleValidation(t *testing.T) {
	t.Run("valid roles should pass validation", func(t *testing.T) {
		assert.True(t, RoleOwner.IsValid())
		assert.True(t, RoleEditor.IsValid())
		assert.True(t, RoleViewer.IsValid())
	})

	t.Run("invalid role should fail validation", func(t *testing.T) {
		invalidRole := Role("INVALID")
		assert.False(t, invalidRole.IsValid())
	})

	t.Run("empty role should fail validation", func(t *testing.T) {
		emptyRole := Role("")
		assert.False(t, emptyRole.IsValid())
	})
}

func TestPermissionChecks(t *testing.T) {
	t.Run("user with invalid role should have no permissions", func(t *testing.T) {
		invalidUser := &User{Role: Role("INVALID")}
		assert.False(t, invalidUser.HasPermission(PermissionRead))
		assert.False(t, invalidUser.CanEdit())
		assert.False(t, invalidUser.CanShare())
		assert.False(t, invalidUser.CanDelete())
	})

	t.Run("user with empty role should have no permissions", func(t *testing.T) {
		emptyRoleUser := &User{Role: Role("")}
		assert.False(t, emptyRoleUser.HasPermission(PermissionRead))
		assert.False(t, emptyRoleUser.CanEdit())
		assert.False(t, emptyRoleUser.CanShare())
		assert.False(t, emptyRoleUser.CanDelete())
	})
}
