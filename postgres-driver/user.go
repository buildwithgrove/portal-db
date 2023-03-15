package postgresdriver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pokt-foundation/portal-db/types"
)

type (
	authProviderDBRow struct {
		UserID         types.UserID       `json:"user_id"`
		Type           types.AuthType     `json:"type"`
		Provider       types.AuthProvider `json:"provider"`
		ProviderUserID string             `json:"provider_user_id"`
		Federated      bool               `json:"federated"`
	}
)

var (
	errUserIDDoesntExist = errors.New("error user ID does not exist for auth provider ID '%s'")
	errUserDoesntExist   = errors.New("error user does not exist for portal ID '%d'")
)

// /* ----- postgresdriver Account Read Methods ----- */

// GetPortalUserIDFromProviderID takes a user's auth provider ID and returns the Portal UserID
func (pg *PostgresDriver) GetPortalUserIDFromProviderID(ctx context.Context, providerUserID string) (types.UserID, error) {
	userID, err := pg.GetPortalUserID(ctx, providerUserID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return types.UserID(0), fmt.Errorf(errUserIDDoesntExist.Error(), providerUserID)
		default:
			return types.UserID(0), err
		}
	}

	return userID, nil
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

// // toAccountUsers converts users from DB rows to map-based Account.Users struct
// func (a *SelectAccountsRow) toAccountUsers() (map[types.Email]types.AccountUserAccess, error) {
// 	var users map[types.Email]types.AccountUserAccess

// 	var userRows []userAccessDBRow
// 	if err := json.Unmarshal(a.Users, &userRows); err != nil {
// 		return users, err
// 	}

// 	users = make(map[types.Email]types.AccountUserAccess, len(userRows))

// 	for _, user := range userRows {
// 		users[types.Email(user.Email)] = types.AccountUserAccess{
// 			UserID:   types.UserID(user.ID),
// 			Email:    types.Email(user.Email),
// 			RoleName: types.RoleName(user.RoleName),
// 			Accepted: user.Accepted,
// 		}
// 	}

// 	return users, nil
// }

// // /* ----- postgresdriver Account Create Methods ----- */

// // WriteAccount creates a single Account in the database, including its OWNER's AccountUserAccess row
// func (pg *PostgresDriver) WriteAccount(ctx context.Context, creatorID types.UserID, account types.Account, createdAt time.Time) (*types.Account, error) {
// 	if account.Plan.Type == types.PayPlanType("") {
// 		return nil, errAccountMustHavePlanTypeSet
// 	}

// 	tx, err := pg.db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}

// 	qtx := pg.WithTx(tx)

// 	userEmail, err := qtx.CheckUserEmail(ctx, creatorID)
// 	if err != nil {
// 		switch {
// 		case strings.Contains(err.Error(), "no rows in result set"):
// 			return nil, errUserDoesNotExist
// 		default:
// 			return nil, err
// 		}
// 	}

// 	// Account created with only PlanType
// 	createdAccount, err := qtx.InsertAccount(ctx, InsertAccountParams{
// 		PlanType:  account.Plan.Type,
// 		CreatedAt: createdAt,
// 		UpdatedAt: createdAt,
// 	})
// 	if err != nil {
// 		_ = tx.Rollback()
// 		return nil, err
// 	}

// 	account.ID = createdAccount.ID
// 	account.CreatedAt = createdAt
// 	account.UpdatedAt = createdAt

// 	// Account creator becomes Account OWNER
// 	owner, err := qtx.InsertAccountUserAccess(ctx, InsertAccountUserAccessParams{
// 		AccountID: createdAccount.ID,
// 		UserEmail: userEmail,
// 		RoleName:  types.RoleOwner,
// 		Accepted:  true,
// 		CreatedAt: createdAt,
// 		UpdatedAt: createdAt,
// 	})
// 	if err != nil {
// 		_ = tx.Rollback()
// 		return nil, err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Assign OWNER to returned Account struct
// 	account.Users = map[types.Email]types.AccountUserAccess{
// 		owner.UserEmail: {
// 			UserID:   types.UserID(owner.UserID),
// 			Email:    owner.UserEmail,
// 			RoleName: owner.RoleName,
// 			Accepted: owner.Accepted,
// 		},
// 	}

// 	return &account, nil
// }

// /* ----- postgresdriver Account Delete Methods ----- */

// // SetAccountDeleted updates a single Account in the database's Deleted field to true
// func (pg *PostgresDriver) SetAccountDeleted(ctx context.Context, accountID types.AccountID, deletedAt time.Time) error {
// 	params := DeleteAccountParams{ID: accountID, DeletedAt: newSQLNullTime(deletedAt)}

// 	err := pg.DeleteAccount(ctx, params)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// /* ----- postgresdriver AccountUserAccess Write Methods ----- */

// // WriteAccountUser saves a single input AccountUserAccess to the database.
// func (pg *PostgresDriver) WriteAccountUser(ctx context.Context, accountUser types.AccountUserAccess, createdAt time.Time) (*types.AccountUserAccess, error) {
// 	params := InsertAccountUserAccessParams{
// 		AccountID: accountUser.AccountID,
// 		UserEmail: accountUser.Email,
// 		RoleName:  accountUser.RoleName,
// 		Accepted:  false,
// 		CreatedAt: createdAt,
// 		UpdatedAt: createdAt,
// 	}

// 	user, err := pg.InsertAccountUserAccess(ctx, params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &types.AccountUserAccess{
// 		UserID:   types.UserID(user.UserID),
// 		Email:    user.UserEmail,
// 		RoleName: user.RoleName,
// 		Accepted: user.Accepted,
// 	}, nil
// }
