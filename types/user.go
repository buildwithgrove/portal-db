package types

import (
	"errors"
	"fmt"
	"time"
)

var (
	errAccountIDIsEmpty error = errors.New("account ID is empty")
	errInvalidRole      error = errors.New("invalid role provided")
)

/* Enums */
type (
	AuthProvider string
	AuthType     string
	Permissions  string
	RoleName     string
)

const (
	AuthProviderAuth0 AuthProvider = "auth0"

	AuthTypeAuth0Github   AuthType = "auth0_github"
	AuthTypeAuth0Username AuthType = "auth0_username"

	PermReadEndpoint     Permissions = "read:endpoint"
	PermWriteEndpoint    Permissions = "write:endpoint"
	PermDeleteEndpoint   Permissions = "delete:endpoint"
	PermTransferEndpoint Permissions = "transfer:endpoint"

	RoleOwner  RoleName = "OWNER"
	RoleAdmin  RoleName = "ADMIN"
	RoleMember RoleName = "MEMBER"
)

func (a AuthType) IsValid() bool {
	switch a {
	case AuthTypeAuth0Github, AuthTypeAuth0Username:
		return true
	default:
		return false
	}
}

func (a AuthType) IsFederated() bool {
	switch a {
	case AuthTypeAuth0Username:
		return false
	case AuthTypeAuth0Github:
		return true
	default:
		return false
	}
}

func (a AuthType) Provider() AuthProvider {
	switch a {
	case AuthTypeAuth0Username, AuthTypeAuth0Github:
		return AuthProviderAuth0
	default:
		return ""
	}
}

func (p Permissions) IsValid() bool {
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
		ID            UserID                        `json:"id"`
		Email         Email                         `json:"email"`
		SignedUp      bool                          `json:"signedUp"`
		AuthProviders map[AuthType]UserAuthProvider `json:"authProviders"`
		CreatedAt     time.Time                     `json:"createdAt"`
		UpdatedAt     time.Time                     `json:"updatedAt"`
	}
	// UserAuthProvider represents a single auth provider for a user (eg. Auth0)
	UserAuthProvider struct {
		ProviderUserID string       `json:"providerUserID"`
		Type           AuthType     `json:"type"`
		Provider       AuthProvider `json:"provider"`
		Federated      bool         `json:"federated"`
	}

	CreateUser struct {
		Email            Email    `json:"email"`
		AuthProviderType AuthType `json:"type"`
		ProviderUserID   string   `json:"providerUserID"`
	}
)

/* UserPermissions Struct Definition and Methods */

type (
	// UserPermissions stores all roles and read/write permissions for all Accounts for a given user
	UserPermissions struct {
		UserID   UserID                           `json:"userID"`
		Accounts map[AccountID]AccountPermissions `json:"loadBalancers"`
	}
	// AccountPermissions stores user role and permissions for a given PortalApp
	AccountPermissions struct {
		RoleName    RoleName      `json:"roleName"`
		Permissions []Permissions `json:"permissions"`
	}
)

var permissionsList = map[RoleName][]Permissions{
	RoleOwner:  {PermReadEndpoint, PermWriteEndpoint, PermDeleteEndpoint, PermTransferEndpoint},
	RoleAdmin:  {PermReadEndpoint, PermWriteEndpoint},
	RoleMember: {PermReadEndpoint},
}

func (u *UserPermissions) IsEmpty() bool {
	if u.UserID == UserID(0) || len(u.Accounts) == 0 {
		return true
	}
	return false
}

func (u *UserPermissions) GetRole(accountID AccountID) RoleName {
	app, ok := u.Accounts[accountID]
	if !ok {
		return RoleName("")
	}

	return app.RoleName
}

func (u *UserPermissions) UpsertPermissions(accountID AccountID, role RoleName) (*UserPermissions, error) {
	if accountID == 0 {
		return nil, errAccountIDIsEmpty
	}
	if !role.IsValid() {
		return nil, errInvalidRole
	}

	u.Accounts[accountID] = AccountPermissions{
		RoleName:    role,
		Permissions: permissionsList[role],
	}

	return u, nil
}

func (u *UserPermissions) DeletePermissions(accountID AccountID) *UserPermissions {
	delete(u.Accounts, accountID)

	return u
}

func (u *UserPermissions) HasPermission(accountID AccountID, permission Permissions) bool {
	app, ok := u.Accounts[accountID]
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

func (p *Permissions) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*p = Permissions(s)
	case string:
		*p = Permissions(s)
	default:
		return fmt.Errorf("unsupported scan type for Permissions: %T", src)
	}
	return nil
}

func (u *User) Table() Table {
	return TableUsers
}

func (u *UserAuthProvider) Table() Table {
	return TableUserAuthProviders
}
