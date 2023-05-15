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
	errUserDoesntExist          = errors.New("error user does not exist for portal ID '%s'")
	errInvalidEmail             = errors.New("error email input is not a valid email address '%s'")
	errInvalidAuthProviderType  = errors.New("error invalid auth provider type '%s'")
	errAuthProviderTypeNotFound = errors.New("error no auth provider type found")
)

// /* ----- postgresdriver User Read Methods ----- */

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
func (pg *PostgresDriver) WriteUserNewSignUp(ctx context.Context, user types.CreateUser, createdAt time.Time) (*types.User, types.AccountID, error) {
	err := pg.validateWriteUserNewSignUpInput(ctx, user)
	if err != nil {
		return nil, types.AccountID(""), err
	}

	id, err := pg.generateID(ctx)
	if err != nil {
		return nil, types.AccountID(""), err
	}
	userID := types.UserID(id)

	authType := user.ProviderUserID.AuthType()

	params := CreateUserNewSignUpParams{
		ID:             userID,
		Email:          user.Email,
		ProviderUserID: user.ProviderUserID,
		Type:           authType,
		Provider:       authType.Provider(),
		Federated:      authType.IsFederated(),
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
	}

	tx, err := pg.DB.Begin()
	if err != nil {
		return nil, types.AccountID(""), err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := pg.WithTx(tx)

	createdUserID, err := qtx.CreateUserNewSignUp(ctx, params)
	if err != nil {
		return nil, types.AccountID(""), err
	}

	// When a new user is created, also create a default Account for them using the free plan
	id, err = pg.generateID(ctx)
	if err != nil {
		return nil, types.AccountID(""), err
	}
	accountID := types.AccountID(id)

	account, err := qtx.InsertAccount(ctx, InsertAccountParams{
		ID:        accountID,
		PlanType:  types.FreetierV0,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	})
	if err != nil {
		return nil, types.AccountID(""), err
	}

	// New user is Account OWNER
	owner, err := qtx.InsertAccountUserAccess(ctx, InsertAccountUserAccessParams{
		AccountID: accountID,
		UserID:    createdUserID,
		Email:     user.Email,
		RoleName:  types.RoleOwner,
		Accepted:  true,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	})
	if err != nil {
		return nil, types.AccountID(""), err
	}

	var providerUserIDs map[types.AuthType]string
	if err := json.Unmarshal(owner.ProviderUserIDs, &providerUserIDs); err != nil {
		return nil, types.AccountID(""), err
	}

	err = tx.Commit()
	if err != nil {
		return nil, types.AccountID(""), err
	}

	createdUser, err := pg.ReadUserByUserID(ctx, createdUserID)
	if err != nil {
		return nil, types.AccountID(""), err
	}

	return createdUser, account.ID, nil
}

// validateWriteUserNewSignUpInput validates the input to create a new User and User Auth Provider
func (pg *PostgresDriver) validateWriteUserNewSignUpInput(ctx context.Context, user types.CreateUser) error {
	if !user.Email.IsValid() {
		return fmt.Errorf(errInvalidEmail.Error(), user.Email)
	}
	if !user.ProviderUserID.AuthType().IsValid() {
		if user.ProviderUserID.AuthType() != types.AuthType("") {
			return fmt.Errorf(errInvalidAuthProviderType.Error(), user.ProviderUserID.AuthType())
		}
		return errAuthProviderTypeNotFound
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
func (json dbUser) toOutput() *types.User {
	return &types.User{
		ID:        json.ID,
		Email:     json.Email,
		SignedUp:  json.SignedUp,
		CreatedAt: json.CreatedAt,
		UpdatedAt: json.UpdatedAt,
	}
}

func (json dbUserAuthProvider) toOutput() *types.UserAuthProvider {
	return &types.UserAuthProvider{
		UserID:         json.UserID,
		ProviderUserID: json.ProviderUserID,
		Type:           json.Type,
		Provider:       json.Provider,
		Federated:      json.Federated,
	}
}

type dbUser struct {
	ID        types.UserID `json:"id"`
	Email     types.Email  `json:"email"`
	SignedUp  bool         `json:"signed_up"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type dbUserAuthProvider struct {
	ID             int32                `json:"id"`
	UserID         types.UserID         `json:"user_id"`
	Type           types.AuthType       `json:"type"`
	Provider       types.AuthProvider   `json:"provider"`
	ProviderUserID types.ProviderUserID `json:"provider_user_id"`
	Federated      bool                 `json:"federated"`
	CreatedAt      time.Time            `json:"created_at"`
}
