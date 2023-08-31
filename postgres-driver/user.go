package postgresdriver

import (
	"context"
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
	errUserAlreadyExists        = errors.New("error user already exists for email'%s' and provider type '%s'")
	errUserHasAccount           = errors.New("error cannot delete user because they are still on an account team")
	errNoEmail                  = errors.New("error no email")
	errInvalidEmail             = errors.New("error email input is not a valid email address '%s'")
	errNoAuthProviderType       = errors.New("error no auth provider type")
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
		switch {
		case errNoRows(err):
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
		ID:               u.ID,
		IconURL:          u.IconURL.String,
		Email:            u.Email,
		SignedUp:         u.SignedUp,
		AuthProviders:    authProviders,
		UpdatesProduct:   u.UpdatesProduct.Bool,
		UpdatesMarketing: u.UpdatesMarketing.Bool,
		BetaTester:       u.BetaTester.Bool,
		CreatedAt:        u.CreatedAt.Time.UTC(),
		UpdatedAt:        u.UpdatedAt.Time.UTC(),
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
		CreatedAt:      newTimestamptz(createdAt),
		UpdatedAt:      newTimestamptz(createdAt),
	}

	tx, err := pg.DB.Begin(ctx)
	if err != nil {
		return nil, types.AccountID(""), err
	}
	defer func() { _ = tx.Rollback(ctx) }()

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
		PlanType:  newText(string(types.FreetierV0)),
		CreatedAt: newTimestamptz(createdAt),
		UpdatedAt: newTimestamptz(createdAt),
	})
	if err != nil {
		return nil, types.AccountID(""), err
	}

	// New user is Account OWNER
	_, err = qtx.InsertAccountUserAccess(ctx, InsertAccountUserAccessParams{
		AccountID: accountID,
		UserID:    createdUserID,
		Email:     user.Email,
		RoleName:  types.RoleOwner,
		Owner:     true,
		Accepted:  true,
		CreatedAt: newTimestamptz(createdAt),
		UpdatedAt: newTimestamptz(createdAt),
	})
	if err != nil {
		return nil, types.AccountID(""), err
	}

	err = tx.Commit(ctx)
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

	userExists, err := pg.CheckUserProviderExists(ctx, CheckUserProviderExistsParams{
		Email: user.Email,
		Type:  user.ProviderUserID.AuthType(),
	})
	if err != nil {
		return err
	}
	if userExists {
		return fmt.Errorf(errUserAlreadyExists.Error(), user.Email, user.ProviderUserID.AuthType())
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
	if err != nil || !userExists {
		return fmt.Errorf(errUserDoesntExist.Error(), userID)
	}

	// Check if the user is part of any accounts
	userHasAccount, err := pg.CheckUserAccountExists(ctx, userID)
	if err != nil {
		return err
	}
	if userHasAccount {
		return errUserHasAccount
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

		for _, appID := range userRoleRow.PortalApplicationIDs {
			portalAppID := types.PortalAppID(appID)

			if userPermissions, ok := userPermissionsMap[userID]; ok {
				_, err := userPermissions.UpsertPermissions(portalAppID, userRoleRow.RoleName)
				if err != nil {
					return nil, err
				}
			} else {
				emptyPermissions := types.UserPermissions{
					UserID:     userID,
					PortalApps: map[types.PortalAppID]types.PortalAppPermissions{},
				}

				permissions, err := emptyPermissions.UpsertPermissions(portalAppID, userRoleRow.RoleName)
				if err != nil {
					return nil, err
				}

				userPermissionsMap[userID] = permissions
			}
		}
	}

	return userPermissionsMap, nil
}

/* ----- Used by Listener ----- */
func (json dbUser) toOutput() *types.User {
	return &types.User{
		ID:               json.ID,
		IconURL:          json.IconURL,
		Email:            json.Email,
		SignedUp:         json.SignedUp,
		UpdatesProduct:   json.UpdatesProduct,
		UpdatesMarketing: json.UpdatesMarketing,
		BetaTester:       json.BetaTester,
		CreatedAt:        json.CreatedAt,
		UpdatedAt:        json.UpdatedAt,
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
	ID               types.UserID `json:"id"`
	IconURL          string       `json:"icon_url"`
	Email            types.Email  `json:"email"`
	SignedUp         bool         `json:"signed_up"`
	UpdatesProduct   bool         `json:"updates_product"`
	UpdatesMarketing bool         `json:"updates_marketing"`
	BetaTester       bool         `json:"beta_tester"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
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
