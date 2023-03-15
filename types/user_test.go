package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testUserPermissions = map[UserID]UserPermissions{
	1: {
		UserID: 1,
		PortalApps: map[PortalAppID]PortalAppPermissions{
			"test_app_6774900350d9c42": {
				RoleName:    RoleOwner,
				Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
			},
		},
	},
	2: {
		UserID: 2,
		PortalApps: map[PortalAppID]PortalAppPermissions{
			"test_app_6774900350d9c42": {
				RoleName:    RoleAdmin,
				Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint},
			},
			"test_app_be591c4dea8566a": {
				RoleName:    RoleMember,
				Permissions: []Permissions{PermReadEndpoint},
			},
		},
	},
	3: {
		UserID: 3,
		PortalApps: map[PortalAppID]PortalAppPermissions{
			"test_app_3487u329rfn23f": {
				RoleName:    RoleOwner,
				Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
			},
		},
	},
}

func Test_UserPermissions_GetRole(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name        string
		userID      UserID
		portalAppID PortalAppID
		roleName    RoleName
	}{
		{
			name:        "Should return the role for a given user and load balancer ID",
			userID:      1,
			portalAppID: "test_app_6774900350d9c42",
			roleName:    RoleOwner,
		},
		{
			name:        "Should return the role for a given user and load balancer ID",
			userID:      2,
			portalAppID: "test_app_be591c4dea8566a",
			roleName:    RoleMember,
		},
		{
			name:        "Should return an empty string if user does not have a role for a load balancer",
			userID:      1,
			portalAppID: "not_an_actual_app",
			roleName:    "",
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		roleName := userPermissions.GetRole(test.portalAppID)
		c.Equal(test.roleName, roleName)
	}
}

func Test_UserPermissions_UpsertPermissions(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		portalAppID         PortalAppID
		roleName            RoleName
		expectedPermissions *UserPermissions
		err                 error
	}{
		{
			name:        "Should add permissions for an App if the user doesn't have any for that App",
			userID:      1,
			portalAppID: "test_new_portal_app",
			roleName:    RoleAdmin,
			expectedPermissions: &UserPermissions{
				UserID: 1,
				PortalApps: map[PortalAppID]PortalAppPermissions{
					"test_app_6774900350d9c42": {
						RoleName:    RoleOwner,
						Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
					},
					"test_new_portal_app": {
						RoleName:    RoleAdmin,
						Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint},
					},
				},
			},
		},
		{
			name:        "Should update permissions for an App if the user already has them for that App",
			userID:      1,
			portalAppID: "test_new_portal_app",
			roleName:    RoleMember,
			expectedPermissions: &UserPermissions{
				UserID: 1,
				PortalApps: map[PortalAppID]PortalAppPermissions{
					"test_app_6774900350d9c42": {
						RoleName:    RoleOwner,
						Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
					},
					"test_new_portal_app": {
						RoleName:    RoleMember,
						Permissions: []Permissions{PermReadEndpoint},
					},
				},
			},
		},
		{
			name:        "Should fail if passed an empty App ID",
			userID:      1,
			portalAppID: "",
			err:         ErrAppIDIsEmpty,
		},
		{
			name:        "Should fail if passed an invalid role",
			userID:      1,
			portalAppID: "test_new_portal_app",
			roleName:    RoleName("not_real"),
			err:         ErrInvalidRole,
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		updatedUserPermissions, err := userPermissions.UpsertPermissions(test.portalAppID, test.roleName)
		c.Equal(test.err, err)
		c.Equal(test.expectedPermissions, updatedUserPermissions)
	}
}

func Test_UserPermissions_DeletePermissions(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		portalAppID         PortalAppID
		roleName            RoleName
		expectedPermissions *UserPermissions
		err                 error
	}{
		{
			name:        "Should delete UserPermissions for a given user and App",
			userID:      3,
			portalAppID: "test_delete_app_id",
			roleName:    RoleMember,
			expectedPermissions: &UserPermissions{
				UserID: 3,
				PortalApps: map[PortalAppID]PortalAppPermissions{
					"test_app_3487u329rfn23f": {
						RoleName:    RoleOwner,
						Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
					},
				},
			},
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		updatedUserPermissions, err := userPermissions.UpsertPermissions(test.portalAppID, test.roleName)
		c.Equal(test.err, err)
		c.Len(updatedUserPermissions.PortalApps, 2)

		deletedUserPermissions := userPermissions.DeletePermissions(test.portalAppID)
		c.Equal(test.expectedPermissions, deletedUserPermissions)
	}
}

func Test_UserPermissions_HasPermission_Read(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name              string
		userID            UserID
		portalAppID       PortalAppID
		hasReadPermission bool
	}{
		{
			name:              "Should return a boolean indicating whether a given user has read permission for a given load balancer",
			userID:            1,
			portalAppID:       "test_app_6774900350d9c42",
			hasReadPermission: true,
		},
		{
			name:              "Should return a boolean indicating whether a given user has read permission for a given load balancer",
			userID:            1,
			portalAppID:       "test_app_3487u329rfn23f9",
			hasReadPermission: false,
		},
		{
			name:              "Should return a boolean indicating whether a given user has read permission for a given load balancer",
			userID:            2,
			portalAppID:       "test_app_be591c4dea8566a",
			hasReadPermission: true,
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		hasReadPermission := userPermissions.HasPermission(test.portalAppID, PermReadEndpoint)
		c.Equal(test.hasReadPermission, hasReadPermission)
	}
}

func Test_UserPermissions_HasPermission_Write(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name              string
		userID            UserID
		portalAppID       PortalAppID
		hasReadPermission bool
	}{
		{
			name:              "Should return a boolean indicating whether a given user has write permission for a given load balancer",
			userID:            1,
			portalAppID:       "test_app_6774900350d9c42",
			hasReadPermission: true,
		},
		{
			name:              "Should return a boolean indicating whether a given user has write permission for a given load balancer",
			userID:            1,
			portalAppID:       "test_app_3487u329rfn23f9",
			hasReadPermission: false,
		},
		{
			name:              "Should return a boolean indicating whether a given user has write permission for a given load balancer",
			userID:            3,
			portalAppID:       "test_app_3487u329rfn23f9",
			hasReadPermission: true,
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		hasReadPermission := userPermissions.HasPermission(test.portalAppID, PermWriteEndpoint)
		c.Equal(test.hasReadPermission, hasReadPermission)
	}
}

func Test_UserPermissions_HasPermission_Delete(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		portalAppID         PortalAppID
		hasDeletePermission bool
	}{
		{
			name:                "Should return a boolean indicating whether a given user has delete permission for a given load balancer",
			userID:              1,
			portalAppID:         "test_app_6774900350d9c42",
			hasDeletePermission: true,
		},
		{
			name:                "Should return a boolean indicating whether a given user has delete permission for a given load balancer",
			userID:              1,
			portalAppID:         "test_app_3487u329rfn23f9",
			hasDeletePermission: false,
		},
		{
			name:                "Should return a boolean indicating whether a given user has delete permission for a given load balancer",
			userID:              3,
			portalAppID:         "test_app_3487u329rfn23f9",
			hasDeletePermission: true,
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		hasReadPermission := userPermissions.HasPermission(test.portalAppID, PermDeleteEndpoint)
		c.Equal(test.hasDeletePermission, hasReadPermission)
	}
}
func Test_UserPermissions_HasPermission_Transfer(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		portalAppID         PortalAppID
		hasDeletePermission bool
	}{
		{
			name:                "Should return a boolean indicating whether a given user has transfer permission for a given load balancer",
			userID:              1,
			portalAppID:         "test_app_6774900350d9c42",
			hasDeletePermission: true,
		},
		{
			name:                "Should return a boolean indicating whether a given user has transfer permission for a given load balancer",
			userID:              1,
			portalAppID:         "test_app_3487u329rfn23f9",
			hasDeletePermission: false,
		},
		{
			name:                "Should return a boolean indicating whether a given user has transfer permission for a given load balancer",
			userID:              3,
			portalAppID:         "test_app_3487u329rfn23f9",
			hasDeletePermission: true,
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		hasReadPermission := userPermissions.HasPermission(test.portalAppID, PermTransferEndpoint)
		c.Equal(test.hasDeletePermission, hasReadPermission)
	}
}
