package types

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrAppIDIsEmpty error = errors.New("load balancer ID is empty")
	ErrNoOwner      error = errors.New("load balancer does not have an owner")
	ErrInvalidRole  error = errors.New("invalid role provided")
)

/* User Struct Definition and Methods */
type (
	// User represents a single Portal user
	User struct {
		ID           string       `json:"id"`
		Email        string       `json:"email"`
		AuthProvider AuthProvider `json:"authProvider"`
		CreatedAt    time.Time    `json:"createdAt"`
		UpdatedAt    time.Time    `json:"updatedAt"`
	}

	AuthProvider    string
	RoleName        string
	PermissionsEnum string
)

const (
	ProviderAuth0 AuthProvider = "Auth0"

	RoleOwner  RoleName = "OWNER"
	RoleAdmin  RoleName = "ADMIN"
	RoleMember RoleName = "MEMBER"

	PermReadEndpoint     PermissionsEnum = "read:endpoint"
	PermWriteEndpoint    PermissionsEnum = "write:endpoint"
	PermDeleteEndpoint   PermissionsEnum = "delete:endpoint"
	PermTransferEndpoint PermissionsEnum = "transfer:endpoint"
)

var (
	ValidRoleNames = map[RoleName]bool{
		RoleOwner:  true,
		RoleAdmin:  true,
		RoleMember: true,
	}

	ValidPermissions = map[PermissionsEnum]bool{
		PermReadEndpoint:  true,
		PermWriteEndpoint: true,
	}

	permissionsList = map[RoleName][]PermissionsEnum{
		RoleOwner:  {PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
		RoleAdmin:  {PermReadEndpoint, PermWriteEndpoint},
		RoleMember: {PermReadEndpoint},
	}
)

func (app *PortalApp) GetOwnerEmail() (string, error) {
	for _, userAccess := range app.Account.Users {
		if userAccess.RoleName == RoleOwner {
			return userAccess.User.Email, nil
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

// UserPermissions stores all load balancer read/write permissions for a given user
type (
	UserPermissions struct {
		UserID     UserID                                 `json:"userID"`
		PortalApps map[ApplicationID]PortalAppPermissions `json:"loadBalancers"`
	}

	PortalAppPermissions struct {
		RoleName    RoleName          `json:"roleName"`
		Permissions []PermissionsEnum `json:"permissions"`
	}
)

func (u *UserPermissions) IsEmpty() bool {
	if u.UserID == UserID("") || len(u.PortalApps) == 0 {
		return true
	}
	return false
}

func (u *UserPermissions) GetRole(appID ApplicationID) RoleName {
	app, ok := u.PortalApps[appID]
	if !ok {
		return RoleName("")
	}

	return app.RoleName
}

func (u *UserPermissions) UpsertPermissions(appID ApplicationID, role RoleName) (*UserPermissions, error) {
	if appID == "" {
		return nil, ErrAppIDIsEmpty
	}
	if !ValidRoleNames[role] {
		return nil, ErrInvalidRole
	}

	u.PortalApps[appID] = PortalAppPermissions{
		RoleName:    role,
		Permissions: permissionsList[role],
	}

	return u, nil
}

func (u *UserPermissions) DeletePermissions(appID ApplicationID) *UserPermissions {
	delete(u.PortalApps, appID)

	return u
}

func (u *UserPermissions) HasPermission(appID ApplicationID, permission PermissionsEnum) bool {
	app, ok := u.PortalApps[appID]
	if !ok {
		return false
	}

	for _, portalAppPermission := range app.Permissions {
		if portalAppPermission == permission {
			return true
		}
	}

	return false
}

func (t *User) Table() Table {
	return TableUsers
}
