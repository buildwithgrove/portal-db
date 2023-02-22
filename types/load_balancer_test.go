package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testUserPermissions = map[UserID]UserPermissions{
	"test_user_687463gh2h72gs": {
		UserID: "test_user_687463gh2h72gs",
		LoadBalancers: map[LoadBalancerID]LoadBalancerPermissions{
			"test_lb_67774900350d9c42": {
				RoleName:    RoleOwner,
				Permissions: []PermissionsEnum{ReadEndpoint, WriteEndpoint, DeleteEndpoint, TransferEndpoint},
			},
		},
	},
	"test_user_47fhf7rwejf943": {
		UserID: "test_user_47fhf7rwejf943",
		LoadBalancers: map[LoadBalancerID]LoadBalancerPermissions{
			"test_lb_67774900350d9c42": {
				RoleName:    RoleAdmin,
				Permissions: []PermissionsEnum{ReadEndpoint, WriteEndpoint},
			},
			"test_lb_be3591c4dea8566a": {
				RoleName:    RoleMember,
				Permissions: []PermissionsEnum{ReadEndpoint},
			},
		},
	},
	"test_user_774900350d9c43": {
		UserID: "test_user_774900350d9c43",
		LoadBalancers: map[LoadBalancerID]LoadBalancerPermissions{
			"test_lb_34987u329rfn23f": {
				RoleName:    RoleOwner,
				Permissions: []PermissionsEnum{ReadEndpoint, WriteEndpoint, DeleteEndpoint, TransferEndpoint},
			},
		},
	},
}

func Test_UserPermissions_GetRole(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name           string
		userID         UserID
		loadBalancerID LoadBalancerID
		roleName       RoleName
	}{
		{
			name:           "Should return the role for a given user and load balancer ID",
			userID:         "test_user_687463gh2h72gs",
			loadBalancerID: "test_lb_67774900350d9c42",
			roleName:       RoleOwner,
		},
		{
			name:           "Should return the role for a given user and load balancer ID",
			userID:         "test_user_47fhf7rwejf943",
			loadBalancerID: "test_lb_be3591c4dea8566a",
			roleName:       RoleMember,
		},
		{
			name:           "Should return an empty string if user does not have a role for a load balancer",
			userID:         "test_user_687463gh2h72gs",
			loadBalancerID: "not_an_actual_lb",
			roleName:       "",
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		roleName := userPermissions.GetRole(test.loadBalancerID)
		c.Equal(test.roleName, roleName)
	}
}

func Test_UserPermissions_UpsertPermissions(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		loadBalancerID      LoadBalancerID
		roleName            RoleName
		expectedPermissions *UserPermissions
		err                 error
	}{
		{
			name:           "Should add permissions for an LB if the user doesn't have any for that LB",
			userID:         "test_user_687463gh2h72gs",
			loadBalancerID: "test_new_load_balancer",
			roleName:       RoleAdmin,
			expectedPermissions: &UserPermissions{
				UserID: "test_user_687463gh2h72gs",
				LoadBalancers: map[LoadBalancerID]LoadBalancerPermissions{
					"test_lb_67774900350d9c42": {
						RoleName:    RoleOwner,
						Permissions: []PermissionsEnum{ReadEndpoint, WriteEndpoint, DeleteEndpoint, TransferEndpoint},
					},
					"test_new_load_balancer": {
						RoleName:    RoleAdmin,
						Permissions: []PermissionsEnum{ReadEndpoint, WriteEndpoint},
					},
				},
			},
		},
		{
			name:           "Should update permissions for an LB if the user already has them for that LB",
			userID:         "test_user_687463gh2h72gs",
			loadBalancerID: "test_new_load_balancer",
			roleName:       RoleMember,
			expectedPermissions: &UserPermissions{
				UserID: "test_user_687463gh2h72gs",
				LoadBalancers: map[LoadBalancerID]LoadBalancerPermissions{
					"test_lb_67774900350d9c42": {
						RoleName:    RoleOwner,
						Permissions: []PermissionsEnum{ReadEndpoint, WriteEndpoint, DeleteEndpoint, TransferEndpoint},
					},
					"test_new_load_balancer": {
						RoleName:    RoleMember,
						Permissions: []PermissionsEnum{ReadEndpoint},
					},
				},
			},
		},
		{
			name:           "Should fail if passed an empty LB ID",
			userID:         "test_user_687463gh2h72gs",
			loadBalancerID: "",
			err:            ErrLBIDIsEmpty,
		},
		{
			name:           "Should fail if passed an invalid role",
			userID:         "test_user_687463gh2h72gs",
			loadBalancerID: "test_new_load_balancer",
			roleName:       RoleName("not_real"),
			err:            ErrInvalidRole,
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		updatedUserPermissions, err := userPermissions.UpsertPermissions(test.loadBalancerID, test.roleName)
		c.Equal(test.err, err)
		c.Equal(test.expectedPermissions, updatedUserPermissions)
	}
}

func Test_UserPermissions_DeletePermissions(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		loadBalancerID      LoadBalancerID
		roleName            RoleName
		expectedPermissions *UserPermissions
		err                 error
	}{
		{
			name:           "Should delete UserPermissions for a given user and LB",
			userID:         "test_user_774900350d9c43",
			loadBalancerID: "test_delete_lb_id",
			roleName:       RoleMember,
			expectedPermissions: &UserPermissions{
				UserID: "test_user_774900350d9c43",
				LoadBalancers: map[LoadBalancerID]LoadBalancerPermissions{
					"test_lb_34987u329rfn23f": {
						RoleName:    RoleOwner,
						Permissions: []PermissionsEnum{ReadEndpoint, WriteEndpoint, DeleteEndpoint, TransferEndpoint},
					},
				},
			},
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		updatedUserPermissions, err := userPermissions.UpsertPermissions(test.loadBalancerID, test.roleName)
		c.Equal(test.err, err)
		c.Len(updatedUserPermissions.LoadBalancers, 2)

		deletedUserPermissions := userPermissions.DeletePermissions(test.loadBalancerID)
		c.Equal(test.expectedPermissions, deletedUserPermissions)
	}
}

func Test_UserPermissions_HasPermission_Read(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name              string
		userID            UserID
		loadBalancerID    LoadBalancerID
		hasReadPermission bool
	}{
		{
			name:              "Should return a boolean indicating whether a given user has read permission for a given load balancer",
			userID:            "test_user_687463gh2h72gs",
			loadBalancerID:    "test_lb_67774900350d9c42",
			hasReadPermission: true,
		},
		{
			name:              "Should return a boolean indicating whether a given user has read permission for a given load balancer",
			userID:            "test_user_687463gh2h72gs",
			loadBalancerID:    "test_lb_34987u329rfn23f",
			hasReadPermission: false,
		},
		{
			name:              "Should return a boolean indicating whether a given user has read permission for a given load balancer",
			userID:            "test_user_47fhf7rwejf943",
			loadBalancerID:    "test_lb_be3591c4dea8566a",
			hasReadPermission: true,
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		hasReadPermission := userPermissions.HasPermission(test.loadBalancerID, ReadEndpoint)
		c.Equal(test.hasReadPermission, hasReadPermission)
	}
}

func Test_UserPermissions_HasPermission_Write(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name              string
		userID            UserID
		loadBalancerID    LoadBalancerID
		hasReadPermission bool
	}{
		{
			name:              "Should return a boolean indicating whether a given user has write permission for a given load balancer",
			userID:            "test_user_687463gh2h72gs",
			loadBalancerID:    "test_lb_67774900350d9c42",
			hasReadPermission: true,
		},
		{
			name:              "Should return a boolean indicating whether a given user has write permission for a given load balancer",
			userID:            "test_user_687463gh2h72gs",
			loadBalancerID:    "test_lb_34987u329rfn23f",
			hasReadPermission: false,
		},
		{
			name:              "Should return a boolean indicating whether a given user has write permission for a given load balancer",
			userID:            "test_user_774900350d9c43",
			loadBalancerID:    "test_lb_34987u329rfn23f",
			hasReadPermission: true,
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		hasReadPermission := userPermissions.HasPermission(test.loadBalancerID, WriteEndpoint)
		c.Equal(test.hasReadPermission, hasReadPermission)
	}
}

func Test_UserPermissions_HasPermission_Delete(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		loadBalancerID      LoadBalancerID
		hasDeletePermission bool
	}{
		{
			name:                "Should return a boolean indicating whether a given user has delete permission for a given load balancer",
			userID:              "test_user_687463gh2h72gs",
			loadBalancerID:      "test_lb_67774900350d9c42",
			hasDeletePermission: true,
		},
		{
			name:                "Should return a boolean indicating whether a given user has delete permission for a given load balancer",
			userID:              "test_user_687463gh2h72gs",
			loadBalancerID:      "test_lb_34987u329rfn23f",
			hasDeletePermission: false,
		},
		{
			name:                "Should return a boolean indicating whether a given user has delete permission for a given load balancer",
			userID:              "test_user_774900350d9c43",
			loadBalancerID:      "test_lb_34987u329rfn23f",
			hasDeletePermission: true,
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		hasReadPermission := userPermissions.HasPermission(test.loadBalancerID, DeleteEndpoint)
		c.Equal(test.hasDeletePermission, hasReadPermission)
	}
}
func Test_UserPermissions_HasPermission_Transfer(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		userID              UserID
		loadBalancerID      LoadBalancerID
		hasDeletePermission bool
	}{
		{
			name:                "Should return a boolean indicating whether a given user has transfer permission for a given load balancer",
			userID:              "test_user_687463gh2h72gs",
			loadBalancerID:      "test_lb_67774900350d9c42",
			hasDeletePermission: true,
		},
		{
			name:                "Should return a boolean indicating whether a given user has transfer permission for a given load balancer",
			userID:              "test_user_687463gh2h72gs",
			loadBalancerID:      "test_lb_34987u329rfn23f",
			hasDeletePermission: false,
		},
		{
			name:                "Should return a boolean indicating whether a given user has transfer permission for a given load balancer",
			userID:              "test_user_774900350d9c43",
			loadBalancerID:      "test_lb_34987u329rfn23f",
			hasDeletePermission: true,
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		hasReadPermission := userPermissions.HasPermission(test.loadBalancerID, TransferEndpoint)
		c.Equal(test.hasDeletePermission, hasReadPermission)
	}
}
