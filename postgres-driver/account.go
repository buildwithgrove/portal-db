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
	userAccessDBRow struct {
		UserID          string                    `json:"user_id"`
		Email           string                    `json:"email"`
		RoleName        string                    `json:"role_name"`
		Accepted        bool                      `json:"accepted"`
		ProviderUserIDs map[types.AuthType]string `json:"provider_user_ids"`
	}
)

var (
	errInvalidRoleName           = errors.New("error invalid role name set")
	errPayPlanDoesntExist        = errors.New("error pay plan '%s' does not exist")
	errAccountDoesntExist        = errors.New("error account does not exist for account ID '%s'")
	errAccountUserDoesntExist    = errors.New("error user ID '%s' does not exist for account ID '%s'")
	errCannotDeleteOwner         = errors.New("error cannot delete user ID '%s' for account ID '%s' because this user is the current account owner")
	errCannotTransferNotAccepted = errors.New("error cannot transfer ownership to user ID '%s' for account ID '%s' because the user has not accepted their invite")
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
	var chainIDs map[types.RelayChainID]struct{}
	if len(a.PartnerChainIDs) != 0 {
		chainIDs = make(map[types.RelayChainID]struct{}, len(a.PartnerChainIDs))
		for _, chainID := range a.PartnerChainIDs {
			chainIDs[types.RelayChainID(chainID)] = struct{}{}
		}
	}

	var partnerChainIDs map[types.RelayChainID]struct{}
	if len(a.PartnerChainIDs) != 0 {

		partnerChainIDs = make(map[types.RelayChainID]struct{}, len(a.PartnerChainIDs))
		for _, chainID := range a.PartnerChainIDs {
			partnerChainIDs[types.RelayChainID(chainID)] = struct{}{}
		}
	}

	accountUsers, err := a.toAccountUsers()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshallingWhitelists, err)
	}

	return &types.Account{
		ID:                     a.ID,
		PlanType:               a.PlanType,
		Users:                  accountUsers,
		PartnerChainIDs:        partnerChainIDs,
		PartnerThroughputLimit: a.PartnerThroughputLimit.Int32,
		PartnerAppLimit:        a.PartnerApplicationLimit.Int32,
		Integrations: types.AccountIntegrations{
			CovalentAPIKeyFree: a.CovalentAPIKeyFree.String,
			CovalentAPIKeyPaid: a.CovalentAPIKeyPaid.String,
		},
		CreatedAt: a.CreatedAt.UTC(),
		UpdatedAt: a.UpdatedAt.UTC(),
		Deleted:   a.Deleted,
	}, nil
}

// toAccountUsers converts users from DB rows to map-based Account.Users struct
func (a *SelectAccountsRow) toAccountUsers() (map[types.UserID]types.AccountUserAccess, error) {
	var users map[types.UserID]types.AccountUserAccess

	var userRows []userAccessDBRow
	if err := json.Unmarshal(a.Users, &userRows); err != nil {
		return users, err
	}

	users = make(map[types.UserID]types.AccountUserAccess, len(userRows))

	for _, user := range userRows {
		users[types.UserID(user.UserID)] = types.AccountUserAccess{
			UserID:          types.UserID(user.UserID),
			Email:           types.Email(user.Email),
			RoleName:        types.RoleName(user.RoleName),
			Accepted:        user.Accepted,
			ProviderUserIDs: user.ProviderUserIDs,
		}
	}

	return users, nil
}

// /* ----- postgresdriver Account Create Methods ----- */

// WriteAccount creates a single Account in the database, including its OWNER's AccountUserAccess row
func (pg *PostgresDriver) WriteAccount(ctx context.Context, creatorID types.UserID, account types.Account, createdAt time.Time) (*types.Account, error) {
	tx, err := pg.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := pg.WithTx(tx)

	err = pg.validateWriteAccountInput(ctx, qtx, creatorID, account)
	if err != nil {
		return nil, err
	}

	id, err := pg.generateID(ctx)
	if err != nil {
		return nil, err
	}
	account.ID = types.AccountID(id)

	createdAccount, err := qtx.InsertAccount(ctx, InsertAccountParams{
		ID:        account.ID,
		PlanType:  account.PlanType,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	})
	if err != nil {
		return nil, err
	}

	account.ID = createdAccount.ID
	account.CreatedAt = createdAt
	account.UpdatedAt = createdAt

	// Account creator becomes Account OWNER
	owner, err := qtx.InsertAccountUserAccess(ctx, InsertAccountUserAccessParams{
		AccountID: createdAccount.ID,
		UserID:    creatorID,
		RoleName:  types.RoleOwner,
		Accepted:  true,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	})
	if err != nil {
		return nil, err
	}

	var providerUserIDs map[types.AuthType]string
	if err := json.Unmarshal(owner.ProviderUserIDs, &providerUserIDs); err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Assign OWNER to returned Account struct
	account.Users = map[types.UserID]types.AccountUserAccess{
		types.UserID(owner.UserID): {
			UserID:          types.UserID(owner.UserID),
			Email:           types.Email(owner.UserEmail),
			RoleName:        owner.RoleName,
			Accepted:        owner.Accepted,
			ProviderUserIDs: providerUserIDs,
		},
	}

	return &account, nil
}

// validateWriteAccountInput validates the input to create a new Account
func (pg *PostgresDriver) validateWriteAccountInput(ctx context.Context, qtx *Queries, creatorID types.UserID, account types.Account) error {
	planExists, err := qtx.CheckPlanTypeExists(ctx, account.PlanType)
	if err != nil {
		return err
	}
	if !planExists {
		return fmt.Errorf(errPayPlanDoesntExist.Error(), account.PlanType)
	}

	userExists, err := qtx.CheckUserExists(ctx, creatorID)
	if err != nil {
		return err
	}
	if !userExists {
		return fmt.Errorf(errUserDoesntExist.Error(), creatorID)
	}

	return nil
}

/* UpsertAccountIntegration saves or updates input AccountIntegrations in the database */
func (pg *PostgresDriver) UpsertAccountIntegration(ctx context.Context, integrations types.AccountIntegrations) (*types.AccountIntegrations, error) {
	time := time.Now()

	accountIntegrations, err := pg.UpsertAccountIntegrations(ctx, UpsertAccountIntegrationsParams{
		AccountID:          integrations.AccountID,
		CovalentAPIKeyFree: newSQLNullString(integrations.CovalentAPIKeyFree),
		CovalentAPIKeyPaid: newSQLNullString(integrations.CovalentAPIKeyPaid),
		CreatedAt:          newSQLNullTime(time),
		UpdatedAt:          newSQLNullTime(time),
	})
	if err != nil {
		return nil, err
	}

	return &types.AccountIntegrations{
		AccountID:          accountIntegrations.AccountID,
		CovalentAPIKeyFree: accountIntegrations.CovalentAPIKeyFree.String,
		CovalentAPIKeyPaid: accountIntegrations.CovalentAPIKeyPaid.String,
	}, nil
}

/* ----- postgresdriver Account Delete Methods ----- */

// SetAccountDeleted updates a single Account in the database's Deleted field to true
func (pg *PostgresDriver) SetAccountDeleted(ctx context.Context, accountID types.AccountID, deletedAt time.Time) error {
	err := pg.validateDeleteAccountInput(ctx, accountID)
	if err != nil {
		return err
	}

	params := DeleteAccountParams{ID: accountID, DeletedAt: newSQLNullTime(deletedAt)}

	err = pg.DeleteAccount(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

// validateDeleteAccountInput validates the input to delete an existing Account
func (pg *PostgresDriver) validateDeleteAccountInput(ctx context.Context, accountID types.AccountID) error {
	accountExists, err := pg.CheckAccountExists(ctx, accountID)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf(errAccountDoesntExist.Error(), accountID)
	}

	return nil
}

/* ----- postgresdriver AccountUserAccess Write Methods ----- */

// WriteAccountUser saves a single input AccountUserAccess to the database.
func (pg *PostgresDriver) WriteAccountUser(ctx context.Context, createAccountUser types.CreateAccountUserAccess, createdAt time.Time) (*types.AccountUserAccess, error) {
	err := pg.validateWriteAccountUserInput(ctx, createAccountUser)
	if err != nil {
		return nil, err
	}

	// determine if user for a given email already exists
	userID, err := pg.CheckUserIDFromEmail(ctx, createAccountUser.Email)
	if err != nil {
		switch err {

		case sql.ErrNoRows:
			// user with provided email does not exist in DB so create a new User and AccountUserAccess entry
			accountUser, err := pg.writeAccountUserAccessNoUser(ctx, createAccountUser, createdAt)
			if err != nil {
				return nil, err
			}

			return accountUser, nil

		default:
			return nil, err
		}
	}

	// user with provided email already exists in DB so create a new AccountUserAccess entry
	accountUser, err := pg.writeAccountUserAccess(ctx, userID, createAccountUser, createdAt)
	if err != nil {
		return nil, err
	}

	return accountUser, nil
}

// validateWriteAccountUserInput validates the input to create a new AccountUserAccess row
func (pg *PostgresDriver) validateWriteAccountUserInput(ctx context.Context, createAccountUser types.CreateAccountUserAccess) error {
	if !createAccountUser.Email.IsValid() {
		return fmt.Errorf(errInvalidEmail.Error(), createAccountUser.Email)
	}

	accountExists, err := pg.CheckAccountExists(ctx, createAccountUser.AccountID)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf(errAccountDoesntExist.Error(), createAccountUser.AccountID)
	}

	return nil
}

// writeAccountUserAccessNoUser creates a new User in the database and then creates a new AccountUserAccess for that user
// Called when a user is invited to a new team but does not yet have a Portal Account for the provided email
func (pg *PostgresDriver) writeAccountUserAccessNoUser(
	ctx context.Context,
	createAccountUser types.CreateAccountUserAccess,
	createdAt time.Time,
) (*types.AccountUserAccess, error) {
	id, err := pg.generateID(ctx)
	if err != nil {
		return nil, err
	}
	userID := types.UserID(id)

	params := InsertAccountUserAccessNoUserParams{
		ID:        userID,
		AccountID: createAccountUser.AccountID,
		Email:     createAccountUser.Email,
		RoleName:  createAccountUser.RoleName,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	user, err := pg.InsertAccountUserAccessNoUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &types.AccountUserAccess{
		UserID:   types.UserID(user.UserID),
		Email:    types.Email(user.UserEmail.String),
		RoleName: user.RoleName,
		Accepted: user.Accepted,
	}, nil
}

// writeAccountUserAccessNoUser creates a new AccountUserAccess row for an existing user
// Called when an existing Portal user is invited to a new team
func (pg *PostgresDriver) writeAccountUserAccess(
	ctx context.Context,
	userID types.UserID,
	createAccountUser types.CreateAccountUserAccess,
	createdAt time.Time,
) (*types.AccountUserAccess, error) {
	params := InsertAccountUserAccessParams{
		UserID:    userID,
		AccountID: createAccountUser.AccountID,
		RoleName:  createAccountUser.RoleName,
		Accepted:  false,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	user, err := pg.InsertAccountUserAccess(ctx, params)
	if err != nil {
		return nil, err
	}

	var providerUserIDs map[types.AuthType]string
	if err := json.Unmarshal(user.ProviderUserIDs, &providerUserIDs); err != nil {
		return nil, err
	}

	return &types.AccountUserAccess{
		UserID:          types.UserID(user.UserID),
		Email:           types.Email(user.UserEmail),
		RoleName:        user.RoleName,
		Accepted:        user.Accepted,
		ProviderUserIDs: providerUserIDs,
	}, nil
}

/* ----- postgresdriver AccountUserAccess Update Methods ----- */

// SetAccountUserRole updates the role for an existing AccountUserAccess row. If transferring ownership the account owner becomes an admin.
func (pg *PostgresDriver) SetAccountUserRole(ctx context.Context, updateAccountUser types.UpdateAccountUserRole, updatedAt time.Time) error {
	tx, err := pg.DB.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := pg.WithTx(tx)

	err = pg.validateSetAccountUserRoleInput(ctx, qtx, updateAccountUser)
	if err != nil {
		return err
	}

	if updateAccountUser.RoleName == types.RoleOwner {
		// if transferring ownership of an account the former OWNER becomes an ADMIN
		err = qtx.UpdateAccountOwnerToAdmin(ctx, updateAccountUser.AccountID)
		if err != nil {

			return err
		}
	}

	params := UpdateAccountUserRoleParams{
		AccountID: updateAccountUser.AccountID,
		UserID:    updateAccountUser.UserID,
		RoleName:  updateAccountUser.RoleName,
		UpdatedAt: updatedAt,
	}

	err = qtx.UpdateAccountUserRole(ctx, params)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// validateSetAccountUserRoleInput validates the input to update an existing AccountUserAccess role
func (pg *PostgresDriver) validateSetAccountUserRoleInput(ctx context.Context, qtx *Queries, updateAccountUser types.UpdateAccountUserRole) error {
	if !updateAccountUser.RoleName.IsValid() {
		return errInvalidRoleName
	}

	existsParams := CheckAccountUserExistsParams{UserID: updateAccountUser.UserID, AccountID: updateAccountUser.AccountID}
	accountUserExists, err := qtx.CheckAccountUserExists(ctx, existsParams)
	if err != nil {
		return err
	}
	if !accountUserExists {
		return fmt.Errorf(errAccountUserDoesntExist.Error(), updateAccountUser.UserID, updateAccountUser.AccountID)
	}

	// cannot transfer OWNER to a user who has not accepted their invite
	if updateAccountUser.RoleName == types.RoleOwner {
		acceptedParams := CheckAccountUserAcceptedParams{UserID: updateAccountUser.UserID, AccountID: updateAccountUser.AccountID}
		userAccepted, err := qtx.CheckAccountUserAccepted(ctx, acceptedParams)
		if err != nil {
			return err
		}
		if !userAccepted {
			return fmt.Errorf(errCannotTransferNotAccepted.Error(), updateAccountUser.UserID, updateAccountUser.AccountID)
		}
	}

	return nil
}

// UpdateAcceptAccountUser creates a new portal UserAuthProvider in the DB when a user accepts their team invite.
// Also updates User.SignedUp and AccountUserAccess.Accepted fields to true.
func (pg *PostgresDriver) UpdateAcceptAccountUser(ctx context.Context, acceptAccountUser types.UpdateAcceptAccountUser, updatedAt time.Time) error {
	tx, err := pg.DB.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := pg.WithTx(tx)

	err = pg.validateUpdateAcceptAccountUserInput(ctx, qtx, acceptAccountUser)
	if err != nil {
		return err
	}

	params := UpdateUserAcceptedInviteParams{
		AccountID:      acceptAccountUser.AccountID,
		UserID:         acceptAccountUser.UserID,
		ProviderUserID: acceptAccountUser.ProviderUserID,
		Type:           acceptAccountUser.AuthProviderType,
		Provider:       acceptAccountUser.AuthProviderType.Provider(),
		Federated:      acceptAccountUser.AuthProviderType.IsFederated(),
	}

	err = qtx.UpdateUserAcceptedInvite(ctx, params)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// validateUpdateAcceptAccountUserInput validates the input to set an AccountUserAccess role to Accepted
func (pg *PostgresDriver) validateUpdateAcceptAccountUserInput(ctx context.Context, qtx *Queries, acceptAccountUser types.UpdateAcceptAccountUser) error {
	if !acceptAccountUser.AuthProviderType.IsValid() {
		return fmt.Errorf(errInvalidAuthProviderType.Error(), acceptAccountUser.AuthProviderType)
	}

	existsParams := CheckAccountUserExistsParams{UserID: acceptAccountUser.UserID, AccountID: acceptAccountUser.AccountID}
	accountUserExists, err := qtx.CheckAccountUserExists(ctx, existsParams)
	if err != nil {
		return err
	}
	if !accountUserExists {
		return fmt.Errorf(errAccountUserDoesntExist.Error(), acceptAccountUser.UserID, acceptAccountUser.AccountID)
	}

	return nil
}

/* ----- postgresdriver AccountUserAccess Delete Methods ----- */

// RemoveAccountUser deletes a AccountUserAccess row for a given user and account ID.
func (pg *PostgresDriver) RemoveAccountUser(ctx context.Context, userID types.UserID, accountID types.AccountID) error {
	err := pg.validateRemoveAccountUserInput(ctx, userID, accountID)
	if err != nil {
		return err
	}

	err = pg.DeleteAccountUser(ctx, DeleteAccountUserParams{UserID: userID, AccountID: accountID})
	if err != nil {
		return err
	}

	return nil
}

// validateRemoveAccountUserInput validates the input to remove an AccountUserAccess row
func (pg *PostgresDriver) validateRemoveAccountUserInput(ctx context.Context, userID types.UserID, accountID types.AccountID) error {
	isOwnerParams := CheckAccountUserRoleParams{UserID: userID, AccountID: accountID}
	accountUserRole, err := pg.CheckAccountUserRole(ctx, isOwnerParams)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return fmt.Errorf(errAccountUserDoesntExist.Error(), userID, accountID)
		default:
			return err
		}
	}
	if accountUserRole == types.RoleOwner {
		return fmt.Errorf(errCannotDeleteOwner.Error(), userID, accountID)
	}
	return nil
}

/* ----- Used by Listener ----- */
func (json dbAccount) toOutput() *types.Account {
	var partnerChainIDs map[types.RelayChainID]struct{}
	if len(json.PartnerChainIDs) != 0 {
		partnerChainIDs = make(map[types.RelayChainID]struct{})
		for _, chainID := range json.PartnerChainIDs {
			partnerChainIDs[types.RelayChainID(chainID)] = struct{}{}
		}
	}

	return &types.Account{
		ID:                     json.ID,
		PlanType:               json.PlanType,
		PartnerChainIDs:        partnerChainIDs,
		PartnerThroughputLimit: json.PartnerThroughputLimit,
		PartnerAppLimit:        json.PartnerApplicationLimit,
		CreatedAt:              json.CreatedAt,
		UpdatedAt:              json.UpdatedAt,
		Deleted:                json.Deleted,
	}
}

func (json dbAccountUserAccess) toOutput() *types.AccountUserAccess {
	return &types.AccountUserAccess{
		AccountID: json.AccountID,
		UserID:    json.UserID,
		RoleName:  json.RoleName,
		Accepted:  json.Accepted,
	}
}

func (j dbAccountIntegration) toOutput() *types.AccountIntegrations {
	return &types.AccountIntegrations{
		AccountID:          j.AccountID,
		CovalentAPIKeyFree: j.CovalentAPIKeyFree,
		CovalentAPIKeyPaid: j.CovalentAPIKeyPaid,
	}
}

type dbAccount struct {
	ID                      types.AccountID   `json:"id"`
	PlanType                types.PayPlanType `json:"plan_type"`
	PartnerChainIDs         []string          `json:"partner_chain_ids"`
	PartnerThroughputLimit  int32             `json:"partner_throughput_limit"`
	PartnerApplicationLimit int32             `json:"partner_application_limit"`
	CreatedAt               time.Time         `json:"created_at"`
	UpdatedAt               time.Time         `json:"updated_at"`
	Deleted                 bool              `json:"deleted"`
	DeletedAt               time.Time         `json:"deleted_at"`
}

type dbAccountIntegration struct {
	ID                 int32           `json:"id"`
	AccountID          types.AccountID `json:"account_id"`
	CovalentAPIKeyFree string          `json:"covalent_api_key_free"`
	CovalentAPIKeyPaid string          `json:"covalent_api_key_paid"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type dbAccountUserAccess struct {
	ID        int32           `json:"id"`
	AccountID types.AccountID `json:"account_id"`
	UserID    types.UserID    `json:"user_id"`
	RoleName  types.RoleName  `json:"role_name"`
	Accepted  bool            `json:"accepted"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
