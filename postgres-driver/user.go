package postgresdriver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pokt-foundation/portal-db/v2/types"
)

type (
	authProviderDBRow struct {
		UserID         types.UserID         `json:"user_id"`
		Type           types.AuthType       `json:"type"`
		Provider       types.AuthProvider   `json:"provider"`
		ProviderUserID types.ProviderUserID `json:"provider_user_id"`
		Federated      bool                 `json:"federated"`
	}
)

var (
	errUserIDDoesntExist       = errors.New("error user ID does not exist for auth provider ID '%s'")
	errUserDoesntExist         = errors.New("error user does not exist for portal ID '%s'")
	errInvalidEmail            = errors.New("error email input is not a valid email address '%s'")
	errInvalidAuthProviderType = errors.New("error invalid auth provider type '%s'")
)

// /* ----- postgresdriver User Read Methods ----- */

// GetPortalUserIDFromProviderID takes a user's auth provider ID and returns the Portal UserID
func (pg *PostgresDriver) GetPortalUserIDFromProviderID(ctx context.Context, providerUserID types.ProviderUserID) (types.UserID, error) {
	userID, err := pg.GetPortalUserID(ctx, providerUserID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return types.UserID(""), fmt.Errorf(errUserIDDoesntExist.Error(), providerUserID)
		default:
			return types.UserID(""), err
		}
	}

	return userID, nil
}

// ReadUserIDsMap returns all user IDs in the database mapped by their provider user IDs
func (pg *PostgresDriver) ReadUserIDsMap(ctx context.Context) (map[types.ProviderUserID]types.UserID, error) {
	userIDs, err := pg.SelectUserIDs(ctx)
	if err != nil {
		return nil, err
	}

	userIDsMap := make(map[types.ProviderUserID]types.UserID, len(userIDs))
	for _, userIDRow := range userIDs {
		userIDsMap[userIDRow.ProviderUserID] = userIDRow.UserID
	}

	return userIDsMap, nil
}

// ReadUserByUserID takes a portal UserID and returns a single user in the database as a User struct
func (pg *PostgresDriver) ReadUserByUserID(ctx context.Context, userID types.UserID) (*types.User, error) {
	userData, err := pg.GetUserDataFromPortalUserID(ctx, userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, fmt.Errorf(errUserDoesntExist.Error(), userID)
		default:
			return nil, err
		}
	}

	user, err := userData.toUser()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// toUser converts User SELECT output to User struct
func (u *GetUserDataFromPortalUserIDRow) toUser() (*types.User, error) {
	var providerRows []authProviderDBRow
	if err := json.Unmarshal(u.AuthProviders, &providerRows); err != nil {
		return nil, err
	}

	authProviders := make(map[types.AuthType]types.UserAuthProvider, len(providerRows))
	for _, provider := range providerRows {
		authProviders[provider.Type] = types.UserAuthProvider{
			ProviderUserID: provider.ProviderUserID,
			Type:           provider.Type,
			Provider:       provider.Provider,
			Federated:      provider.Federated,
		}
	}

	return &types.User{
		ID:            u.ID,
		Email:         u.Email,
		SignedUp:      u.SignedUp,
		AuthProviders: authProviders,
		CreatedAt:     u.CreatedAt.UTC(),
		UpdatedAt:     u.UpdatedAt.UTC(),
	}, nil
}

/* ----- postgresdriver User Create Methods ----- */

// WriteUserNewSignUp creates a new portal User and UserAuthProviderin the DB from a CreateUser input when a new user signs up.
func (pg *PostgresDriver) WriteUserNewSignUp(ctx context.Context, user types.CreateUser, createdAt time.Time) (types.UserID, error) {
	err := pg.validateWriteUserNewSignUpInput(ctx, user)
	if err != nil {
		return types.UserID(""), err
	}

	id, err := pg.generateID(ctx)
	if err != nil {
		return types.UserID(""), err
	}
	userID := types.UserID(id)

	params := CreateUserNewSignUpParams{
		ID:             userID,
		Email:          user.Email,
		ProviderUserID: user.ProviderUserID,
		Type:           user.AuthProviderType,
		Provider:       user.AuthProviderType.Provider(),
		Federated:      user.AuthProviderType.IsFederated(),
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
	}

	createdUserID, err := pg.CreateUserNewSignUp(ctx, params)
	if err != nil {
		return types.UserID(""), err
	}

	return types.UserID(createdUserID), nil
}

// validateWriteUserNewSignUpInput validates the input to create a new User and User Auth Provider
func (pg *PostgresDriver) validateWriteUserNewSignUpInput(ctx context.Context, user types.CreateUser) error {
	if !user.Email.IsValid() {
		return fmt.Errorf(errInvalidEmail.Error(), user.Email)
	}
	if !user.AuthProviderType.IsValid() {
		return fmt.Errorf(errInvalidAuthProviderType.Error(), user.AuthProviderType)
	}

	return nil
}

/* ----- postgresdriver User Delete Methods ----- */

// DeletePortalUser deletes a portal User from the DB. WARNING will do a hard delete.
// Will also delete the user's `account_user_access` and `user_auth_providers` rows.
func (pg *PostgresDriver) DeletePortalUser(ctx context.Context, userID types.UserID) (types.UserID, error) {
	err := pg.validateDeletePortalUserInput(ctx, userID)
	if err != nil {
		return types.UserID(""), err
	}

	deletedUserID, err := pg.DeleteUser(ctx, userID)
	if err != nil {
		return types.UserID(""), err
	}

	return types.UserID(deletedUserID), nil
}

// validateDeletePortalUserInput validates the input to delete a User
func (pg *PostgresDriver) validateDeletePortalUserInput(ctx context.Context, userID types.UserID) error {
	userExists, err := pg.CheckUserExists(ctx, userID)
	if err != nil {
		return err
	}
	if !userExists {
		return fmt.Errorf(errUserDoesntExist.Error(), userID)
	}

	return nil
}

/* ----- postgresdriver UserPermissions Read Methods ----- */

/* ReadUserPermissions returns all UserPermissions in the database as a map that takes the form map[types.UserID]*types.UserPermissions */
func (pg *PostgresDriver) ReadUserPermissions(ctx context.Context) (map[types.UserID]*types.UserPermissions, error) {
	userRoles, err := pg.SelectUserPermissions(ctx)
	if err != nil {
		return nil, err
	}

	userPermissionsMap := make(map[types.UserID]*types.UserPermissions)

	for _, userRoleRow := range userRoles {
		userID := userRoleRow.UserID
		accountID := userRoleRow.AccountID

		if userPermissions, ok := userPermissionsMap[userID]; ok {
			_, err := userPermissions.UpsertPermissions(accountID, userRoleRow.RoleName)
			if err != nil {
				return nil, err
			}
		} else {
			emptyPermissions := types.UserPermissions{
				UserID:   userID,
				Accounts: map[types.AccountID]types.AccountPermissions{},
			}

			permissions, err := emptyPermissions.UpsertPermissions(accountID, userRoleRow.RoleName)
			if err != nil {
				return nil, err
			}

			userPermissionsMap[userID] = permissions
		}
	}

	return userPermissionsMap, nil
}

/* ----- Used by Listener ----- */
func (json User) toOutput() *types.User {
	return &types.User{
		ID:        json.ID,
		Email:     json.Email,
		SignedUp:  json.SignedUp,
		CreatedAt: json.CreatedAt,
		UpdatedAt: json.UpdatedAt,
	}
}

func (json UserAuthProvider) toOutput() *types.UserAuthProvider {
	return &types.UserAuthProvider{
		ProviderUserID: json.ProviderUserID,
		Type:           json.Type,
		Provider:       json.Provider,
		Federated:      json.Federated,
	}
}
