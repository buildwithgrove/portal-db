package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testUserPermissions = map[UserID]UserPermissions{
	1: {
		UserID: 1,
		Accounts: map[AccountID]AccountPermissions{
			1: {
				RoleName:    RoleOwner,
				Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
			},
		},
	},
	2: {
		UserID: 2,
		Accounts: map[AccountID]AccountPermissions{
			1: {
				RoleName:    RoleAdmin,
				Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint},
			},
			2: {
				RoleName:    RoleMember,
				Permissions: []Permissions{PermReadEndpoint},
			},
		},
	},
	3: {
		UserID: 3,
		Accounts: map[AccountID]AccountPermissions{
			3: {
				RoleName:    RoleOwner,
				Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
			},
		},
	},
}

func Test_UserPermissions_GetRole(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name      string
		userID    UserID
		accountID AccountID
		roleName  RoleName
	}{
		{
			name:      "Should return the role for a given user and portal application ID",
			userID:    1,
			accountID: 1,
			roleName:  RoleOwner,
		},
		{
			name:      "Should return the role for a given user and portal application ID",
			userID:    2,
			accountID: 2,
			roleName:  RoleMember,
		},
		{
			name:      "Should return an empty string if user does not have a role for a portal application",
			userID:    1,
			accountID: 77,
			roleName:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userPermissions := testUserPermissions[test.userID]

			roleName := userPermissions.GetRole(test.accountID)
			c.Equal(test.roleName, roleName)
		})
	}
}

func Test_UserPermissions_UpsertPermissions(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		accountID           AccountID
		roleName            RoleName
		expectedPermissions *UserPermissions
		err                 error
	}{
		{
			name:      "Should add permissions for an App if the user doesn't have any for that App",
			userID:    1,
			accountID: 4,
			roleName:  RoleAdmin,
			expectedPermissions: &UserPermissions{
				UserID: 1,
				Accounts: map[AccountID]AccountPermissions{
					1: {
						RoleName:    RoleOwner,
						Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
					},
					4: {
						RoleName:    RoleAdmin,
						Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint},
					},
				},
			},
		},
		{
			name:      "Should update permissions for an App if the user already has them for that App",
			userID:    1,
			accountID: 4,
			roleName:  RoleMember,
			expectedPermissions: &UserPermissions{
				UserID: 1,
				Accounts: map[AccountID]AccountPermissions{
					1: {
						RoleName:    RoleOwner,
						Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
					},
					4: {
						RoleName:    RoleMember,
						Permissions: []Permissions{PermReadEndpoint},
					},
				},
			},
		},
		{
			name:      "Should fail if passed an empty App ID",
			userID:    1,
			accountID: 0,
			err:       errAccountIDIsEmpty,
		},
		{
			name:      "Should fail if passed an invalid role",
			userID:    1,
			accountID: 4,
			roleName:  RoleName("not_real"),
			err:       errInvalidRole,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userPermissions := testUserPermissions[test.userID]

			updatedUserPermissions, err := userPermissions.UpsertPermissions(test.accountID, test.roleName)
			c.Equal(test.err, err)
			c.Equal(test.expectedPermissions, updatedUserPermissions)
		})
	}
}

func Test_UserPermissions_DeletePermissions(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		accountID           AccountID
		roleName            RoleName
		expectedPermissions *UserPermissions
		err                 error
	}{
		{
			name:      "Should delete UserPermissions for a given user and App",
			userID:    3,
			accountID: 5,
			roleName:  RoleMember,
			expectedPermissions: &UserPermissions{
				UserID: 3,
				Accounts: map[AccountID]AccountPermissions{
					3: {
						RoleName:    RoleOwner,
						Permissions: []Permissions{PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userPermissions := testUserPermissions[test.userID]

			updatedUserPermissions, err := userPermissions.UpsertPermissions(test.accountID, test.roleName)
			c.Equal(test.err, err)
			c.Len(updatedUserPermissions.Accounts, 2)

			deletedUserPermissions := userPermissions.DeletePermissions(test.accountID)
			c.Equal(test.expectedPermissions, deletedUserPermissions)
		})
	}
}

func Test_UserPermissions_HasPermission_Read(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name              string
		userID            UserID
		accountID         AccountID
		hasReadPermission bool
	}{
		{
			name:              "Should return a boolean indicating whether a given user has read permission for a given portal application",
			userID:            1,
			accountID:         1,
			hasReadPermission: true,
		},
		{
			name:              "Should return a boolean indicating whether a given user has read permission for a given portal application",
			userID:            1,
			accountID:         89,
			hasReadPermission: false,
		},
		{
			name:              "Should return a boolean indicating whether a given user has read permission for a given portal application",
			userID:            2,
			accountID:         2,
			hasReadPermission: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userPermissions := testUserPermissions[test.userID]

			hasReadPermission := userPermissions.HasPermission(test.accountID, PermReadEndpoint)
			c.Equal(test.hasReadPermission, hasReadPermission)
		})
	}
}

func Test_UserPermissions_HasPermission_Write(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name               string
		userID             UserID
		accountID          AccountID
		hasWritePermission bool
	}{
		{
			name:               "Should return a boolean indicating whether a given user has write permission for a given portal application",
			userID:             1,
			accountID:          1,
			hasWritePermission: true,
		},
		{
			name:               "Should return a boolean indicating whether a given user has write permission for a given portal application",
			userID:             1,
			accountID:          89,
			hasWritePermission: false,
		},
		{
			name:               "Should return a boolean indicating whether a given user has write permission for a given portal application",
			userID:             3,
			accountID:          3,
			hasWritePermission: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userPermissions := testUserPermissions[test.userID]

			hasWritePermission := userPermissions.HasPermission(test.accountID, PermWriteEndpoint)
			c.Equal(test.hasWritePermission, hasWritePermission)
		})
	}
}

func Test_UserPermissions_HasPermission_Delete(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		accountID           AccountID
		hasDeletePermission bool
	}{
		{
			name:                "Should return a boolean indicating whether a given user has delete permission for a given portal application",
			userID:              1,
			accountID:           1,
			hasDeletePermission: true,
		},
		{
			name:                "Should return a boolean indicating whether a given user has delete permission for a given portal application",
			userID:              1,
			accountID:           89,
			hasDeletePermission: false,
		},
		{
			name:                "Should return a boolean indicating whether a given user has delete permission for a given portal application",
			userID:              3,
			accountID:           3,
			hasDeletePermission: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userPermissions := testUserPermissions[test.userID]

			hasReadPermission := userPermissions.HasPermission(test.accountID, PermDeleteEndpoint)
			c.Equal(test.hasDeletePermission, hasReadPermission)
		})
	}
}
func Test_UserPermissions_HasPermission_Transfer(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		accountID           AccountID
		hasDeletePermission bool
	}{
		{
			name:                "Should return a boolean indicating whether a given user has transfer permission for a given portal application",
			userID:              1,
			accountID:           1,
			hasDeletePermission: true,
		},
		{
			name:                "Should return a boolean indicating whether a given user has transfer permission for a given portal application",
			userID:              1,
			accountID:           89,
			hasDeletePermission: false,
		},
		{
			name:                "Should return a boolean indicating whether a given user has transfer permission for a given portal application",
			userID:              3,
			accountID:           3,
			hasDeletePermission: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userPermissions := testUserPermissions[test.userID]

			hasReadPermission := userPermissions.HasPermission(test.accountID, PermTransferEndpoint)
			c.Equal(test.hasDeletePermission, hasReadPermission)
		})
	}
}
