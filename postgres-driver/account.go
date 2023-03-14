package postgresdriver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pokt-foundation/portal-db/types"
)

type (
	userAccessDBRow struct {
		ID           string `json:"user_id"`
		Email        string `json:"email"`
		AuthProvider string `json:"auth_provider"`
		RoleName     string `json:"role_name"`
		Accepted     bool   `json:"accepted"`
	}
)

var (
	errAccountMustHavePlanTypeSet = errors.New("error account input does not have a plan type set")
	errUserDoesNotExist           = errors.New("error account creator does not exist in the database")
)

/* ----- postgresdriver Account Read Methods ----- */

// ReadAccounts returns all Accounts in the database as Accounts structs
func (pg *PostgresDriver) ReadAccounts(ctx context.Context, options types.DriverOptions) (map[types.AccountID]*types.Account, error) {
	dbAccounts, err := pg.SelectAccounts(ctx, options.IncludeDeleted)
	if err != nil {
		return nil, err
	}

	accounts := make(map[types.AccountID]*types.Account, len(dbAccounts))
	for _, dbAccount := range dbAccounts {
		account, err := dbAccount.toAccount()
		if err != nil {
			return nil, err
		}

		accounts[dbAccount.ID] = account
	}

	return accounts, nil
}

// toAccount converts Account SELECT output to Account struct
func (a *SelectAccountsRow) toAccount() (*types.Account, error) {
	var chainIDs map[types.ChainID]struct{}
	if len(a.ChainIDs) != 0 {
		chainIDs = make(map[types.ChainID]struct{}, len(a.ChainIDs))
		for _, chainID := range a.ChainIDs {
			chainIDs[types.ChainID(chainID)] = struct{}{}
		}
	}

	var partnerChainIDs map[types.ChainID]struct{}
	if len(a.PartnerChainIDs) != 0 {
		partnerChainIDs = make(map[types.ChainID]struct{}, len(a.PartnerChainIDs))
		for _, chainID := range a.PartnerChainIDs {
			partnerChainIDs[types.ChainID(chainID)] = struct{}{}
		}
	}

	accountUsers, err := a.toAccountUsers()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshallingWhitelists, err)
	}

	return &types.Account{
		ID: a.ID,
		Plan: types.Plan{
			Type:              types.PayPlanType(a.PlanType),
			ChainIDs:          chainIDs,
			MonthlyRelayLimit: a.MonthlyRelayLimit.Int32,
			ThroughputLimit:   a.ThroughputLimit.Int32,
			AppLimit:          a.ApplicationLimit.Int32,
			LegacyDailyLimit:  a.DailyLimit.Int32,
		},
		Users:                  accountUsers,
		PartnerChainIDs:        partnerChainIDs,
		PartnerThroughputLimit: a.PartnerThroughputLimit.Int32,
		PartnerAppLimit:        a.PartnerApplicationLimit.Int32,
		CreatedAt:              a.CreatedAt.UTC(),
		UpdatedAt:              a.UpdatedAt.UTC(),
		Deleted:                a.Deleted,
	}, nil
}

// toAccountUsers converts users from DB rows to map-based Account.Users struct
func (a *SelectAccountsRow) toAccountUsers() (map[types.Email]types.AccountUserAccess, error) {
	var users map[types.Email]types.AccountUserAccess

	var userRows []userAccessDBRow
	if err := json.Unmarshal(a.Users, &userRows); err != nil {
		return users, err
	}

	users = make(map[types.Email]types.AccountUserAccess, len(userRows))

	for _, user := range userRows {
		users[types.Email(user.Email)] = types.AccountUserAccess{
			UserID:   types.UserID(user.ID),
			Email:    types.Email(user.Email),
			RoleName: types.RoleName(user.RoleName),
			Accepted: user.Accepted,
		}
	}

	return users, nil
}

// /* ----- postgresdriver Account Create Methods ----- */

// WriteAccount creates a single Account in the database, including its OWNER's AccountUserAccess row
func (pg *PostgresDriver) WriteAccount(ctx context.Context, creatorID types.UserID, account types.Account, createdAt time.Time) (*types.Account, error) {
	if account.Plan.Type == types.PayPlanType("") {
		return nil, errAccountMustHavePlanTypeSet
	}

	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	qtx := pg.WithTx(tx)

	userEmail, err := qtx.CheckUserEmail(ctx, creatorID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "no rows in result set"):
			return nil, errUserDoesNotExist
		default:
			return nil, err
		}
	}

	// Account created with only PlanType
	createdAccount, err := qtx.InsertAccount(ctx, InsertAccountParams{
		PlanType:  account.Plan.Type,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	account.ID = createdAccount.ID
	account.CreatedAt = createdAt
	account.UpdatedAt = createdAt

	// Account creator becomes Account OWNER
	owner, err := qtx.InsertAccountUserAccess(ctx, InsertAccountUserAccessParams{
		AccountID: createdAccount.ID,
		UserEmail: userEmail,
		RoleName:  types.RoleOwner,
		Accepted:  true,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Assign OWNER to returned Account struct
	account.Users = map[types.Email]types.AccountUserAccess{
		owner.UserEmail: {
			UserID:   types.UserID(owner.UserID),
			Email:    owner.UserEmail,
			RoleName: owner.RoleName,
			Accepted: owner.Accepted,
		},
	}

	return &account, nil
}

/* ----- postgresdriver Account Delete Methods ----- */

// SetAccountDeleted updates a single Account in the database's Deleted field to true
func (pg *PostgresDriver) SetAccountDeleted(ctx context.Context, accountID types.AccountID, deletedAt time.Time) error {
	params := DeleteAccountParams{ID: accountID, DeletedAt: newSQLNullTime(deletedAt)}

	err := pg.DeleteAccount(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

/* ----- postgresdriver AccountUserAccess Write Methods ----- */

// WriteAccountUser saves a single input AccountUserAccess to the database.
func (pg *PostgresDriver) WriteAccountUser(ctx context.Context, accountUser types.AccountUserAccess, createdAt time.Time) (*types.AccountUserAccess, error) {
	params := InsertAccountUserAccessParams{
		AccountID: accountUser.AccountID,
		UserEmail: accountUser.Email,
		RoleName:  accountUser.RoleName,
		Accepted:  false,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	user, err := pg.InsertAccountUserAccess(ctx, params)
	if err != nil {
		return nil, err
	}

	return &types.AccountUserAccess{
		UserID:   types.UserID(user.UserID),
		Email:    user.UserEmail,
		RoleName: user.RoleName,
		Accepted: user.Accepted,
	}, nil
}
