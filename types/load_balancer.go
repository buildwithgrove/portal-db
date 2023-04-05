package types

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrLBIDIsEmpty error = errors.New("load balancer ID is empty")
	ErrNoOwner     error = errors.New("load balancer does not have an owner")
	ErrInvalidRole error = errors.New("invalid role provided")
)

/* LB Apps Table represents DB relationship of LBs and apps */
// do not change the tags, they're snake_case on purpose
type LbApp struct {
	LbID  string `json:"lb_id"`
	AppID string `json:"app_id"`
}

/* Load Balancers Table */
type (
	LoadBalancer struct {
		ID                string              `json:"id"`
		Name              string              `json:"name"`
		UserID            string              `json:"userID"`
		ApplicationIDs    []string            `json:"applicationIDs,omitempty"`
		RequestTimeout    int                 `json:"requestTimeout"`
		Gigastake         bool                `json:"gigastake"`
		GigastakeRedirect bool                `json:"gigastakeRedirect"`
		StickyOptions     StickyOptions       `json:"stickinessOptions"`
		Integrations      AccountIntegrations `json:"integrations"`
		Applications      []*Application      `json:"applications"`
		Users             []UserAccess        `json:"users"`
		CreatedAt         time.Time           `json:"createdAt"`
		UpdatedAt         time.Time           `json:"updatedAt"`
		AccountID         string              `json:"accountID"`
	}
	AccountIntegrations struct {
		AccountID          string `json:"id,omitempty"`
		CovalentAPIKeyFree string `json:"covalentAPIKeyFree"`
		CovalentAPIKeyPaid string `json:"covalentAPIKeyPaid"`
	}
	StickyOptions struct {
		ID            string   `json:"id,omitempty"`
		Duration      string   `json:"duration"`
		StickyOrigins []string `json:"stickyOrigins"`
		StickyMax     int      `json:"stickyMax"`
		Stickiness    bool     `json:"stickiness"`
	}
	UserAccess struct {
		ID       string   `json:"id,omitempty"`
		UserID   string   `json:"userID"`
		RoleName RoleName `json:"roleName"`
		Email    string   `json:"email"`
		Accepted bool     `json:"accepted"`
	}
	/* Update structs */
	UpdateLoadBalancer struct {
		Name          string               `json:"name,omitempty"`
		StickyOptions *UpdateStickyOptions `json:"stickinessOptions,omitempty"`
		Remove        bool                 `json:"remove,omitempty"`
	}
	UpdateStickyOptions struct {
		ID            string   `json:"id,omitempty"`
		Duration      string   `json:"duration"`
		StickyOrigins []string `json:"stickyOrigins"`
		StickyMax     int      `json:"stickyMax"`
		Stickiness    *bool    `json:"stickiness"`
	}
	UpdateUserAccess struct {
		ID           string   `json:"id,omitempty"`
		UserID       string   `json:"userID"`
		Email        string   `json:"email"`
		UpdaterEmail string   `json:"updaterEmail,omitempty"`
		RoleName     RoleName `json:"roleName"`
	}

	RoleName        string
	PermissionsEnum string
)

const (
	RoleOwner  RoleName = "OWNER"
	RoleAdmin  RoleName = "ADMIN"
	RoleMember RoleName = "MEMBER"

	ReadEndpoint     PermissionsEnum = "read:endpoint"
	WriteEndpoint    PermissionsEnum = "write:endpoint"
	DeleteEndpoint   PermissionsEnum = "delete:endpoint"
	TransferEndpoint PermissionsEnum = "transfer:endpoint"
)

var (
	ValidRoleNames = map[RoleName]bool{
		RoleOwner:  true,
		RoleAdmin:  true,
		RoleMember: true,
	}

	ValidPermissions = map[PermissionsEnum]bool{
		ReadEndpoint:  true,
		WriteEndpoint: true,
	}

	permissionsList = map[RoleName][]PermissionsEnum{
		RoleOwner:  {ReadEndpoint, WriteEndpoint, DeleteEndpoint, TransferEndpoint},
		RoleAdmin:  {ReadEndpoint, WriteEndpoint},
		RoleMember: {ReadEndpoint},
	}
)

func (lb *LoadBalancer) GetOwnerEmail() (string, error) {
	for _, user := range lb.Users {
		if user.RoleName == RoleOwner {
			return user.Email, nil
		}
	}

	return "", ErrNoOwner
}

func (e *PermissionsEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PermissionsEnum(s)
	case string:
		*e = PermissionsEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for PermissionsEnum: %T", src)
	}
	return nil
}

func (s *StickyOptions) IsEmpty() bool {
	if !s.Stickiness {
		return true
	}
	return len(s.StickyOrigins) == 0
}

// UserPermissions stores all load balancer read/write permissions for a given user
type (
	UserID         string
	LoadBalancerID string

	UserPermissions struct {
		UserID        UserID                                     `json:"userID"`
		LoadBalancers map[LoadBalancerID]LoadBalancerPermissions `json:"loadBalancers"`
	}

	LoadBalancerPermissions struct {
		RoleName    RoleName          `json:"roleName"`
		Permissions []PermissionsEnum `json:"permissions"`
	}
)

func (u *UserPermissions) IsEmpty() bool {
	if u.UserID == UserID("") || len(u.LoadBalancers) == 0 {
		return true
	}
	return false
}

func (u *UserPermissions) GetRole(loadBalancerID LoadBalancerID) RoleName {
	lb, ok := u.LoadBalancers[loadBalancerID]
	if !ok {
		return RoleName("")
	}

	return lb.RoleName
}

func (u *UserPermissions) UpsertPermissions(loadBalancerID LoadBalancerID, role RoleName) (*UserPermissions, error) {
	if loadBalancerID == "" {
		return nil, ErrLBIDIsEmpty
	}
	if !ValidRoleNames[role] {
		return nil, ErrInvalidRole
	}

	u.LoadBalancers[loadBalancerID] = LoadBalancerPermissions{
		RoleName:    role,
		Permissions: permissionsList[role],
	}

	return u, nil
}

func (u *UserPermissions) DeletePermissions(loadBalancerID LoadBalancerID) *UserPermissions {
	delete(u.LoadBalancers, loadBalancerID)

	return u
}

func (u *UserPermissions) HasPermission(loadBalancerID LoadBalancerID, permission PermissionsEnum) bool {
	lb, ok := u.LoadBalancers[loadBalancerID]
	if !ok {
		return false
	}

	for _, loadBalancerPermission := range lb.Permissions {
		if loadBalancerPermission == permission {
			return true
		}
	}

	return false
}
