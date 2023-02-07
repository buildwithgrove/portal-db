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
				Permissions: []PermissionsEnum{ReadEndpoint, WriteEndpoint},
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
				Permissions: []PermissionsEnum{ReadEndpoint, WriteEndpoint},
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
		RoleName       RoleName
	}{
		{
			name:           "Should return the role for a given user and load balancer ID",
			userID:         "test_user_687463gh2h72gs",
			loadBalancerID: "test_lb_67774900350d9c42",
			RoleName:       RoleOwner,
		},
		{
			name:           "Should return the role for a given user and load balancer ID",
			userID:         "test_user_47fhf7rwejf943",
			loadBalancerID: "test_lb_be3591c4dea8566a",
			RoleName:       RoleMember,
		},
		{
			name:           "Should return an empty string if user does not have a role for a load balancer",
			userID:         "test_user_687463gh2h72gs",
			loadBalancerID: "not_an_actual_lb",
			RoleName:       "",
		},
	}

	for _, test := range tests {
		userPermissions := testUserPermissions[test.userID]

		RoleName := userPermissions.GetRole(test.loadBalancerID)
		c.Equal(test.RoleName, RoleName)
	}
}

func Test_UserPermissions_HasReadPermission(t *testing.T) {
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

		hasReadPermission := userPermissions.HasReadPermission(test.loadBalancerID)
		c.Equal(test.hasReadPermission, hasReadPermission)
	}
}

func Test_UserPermissions_HasWritePermission(t *testing.T) {
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

		hasReadPermission := userPermissions.HasWritePermission(test.loadBalancerID)
		c.Equal(test.hasReadPermission, hasReadPermission)
	}
}
