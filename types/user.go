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

/* Enums */
type (
	AuthProviders   string
	AuthSignIn      string
	PermissionsEnum string
	RoleName        string
)

const (
	AuthProviderAuth0 AuthProviders = "auth0"

	AuthSignInGitHub   AuthSignIn = "github"
	AuthSignInUsername AuthSignIn = "username"

	PermReadEndpoint     PermissionsEnum = "read:endpoint"
	PermWriteEndpoint    PermissionsEnum = "write:endpoint"
	PermDeleteEndpoint   PermissionsEnum = "delete:endpoint"
	PermTransferEndpoint PermissionsEnum = "transfer:endpoint"

	RoleOwner  RoleName = "OWNER"
	RoleAdmin  RoleName = "ADMIN"
	RoleMember RoleName = "MEMBER"
)

func (a AuthProviders) IsValid() bool {
	switch a {
	case AuthProviderAuth0:
		return true
	default:
		return false
	}
}

func (a AuthSignIn) IsValid() bool {
	switch a {
	case AuthSignInGitHub, AuthSignInUsername:
		return true
	default:
		return false
	}
}

func (p PermissionsEnum) IsValid() bool {
	switch p {
	case PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint:
		return true
	default:
		return false
	}
}

func (r RoleName) IsValid() bool {
	switch r {
	case RoleOwner, RoleAdmin, RoleMember:
		return true
	default:
		return false
	}
}

/* User Struct Definition and Methods */
type (
	// User represents a single Portal user
	User struct {
		ID           string        `json:"id"`
		Email        string        `json:"email"`
		AuthProvider AuthProviders `json:"authProvider"`
		CreatedAt    time.Time     `json:"createdAt"`
		UpdatedAt    time.Time     `json:"updatedAt"`
	}
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

// UserPermissions stores all load balancer read/write permissions for a given user
type (
	UserPermissions struct {
		UserID     UserID                               `json:"userID"`
		PortalApps map[PortalAppID]PortalAppPermissions `json:"loadBalancers"`
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

func (u *UserPermissions) GetRole(appID PortalAppID) RoleName {
	app, ok := u.PortalApps[appID]
	if !ok {
		return RoleName("")
	}

	return app.RoleName
}

func (u *UserPermissions) UpsertPermissions(appID PortalAppID, role RoleName) (*UserPermissions, error) {
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

func (u *UserPermissions) DeletePermissions(appID PortalAppID) *UserPermissions {
	delete(u.PortalApps, appID)

	return u
}

func (u *UserPermissions) HasPermission(appID PortalAppID, permission PermissionsEnum) bool {
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

func (u *User) Table() Table {
	return TableUsers
}
