package postgresdriver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pokt-foundation/portal-db/v2/types"
)

type (
	userAccessDBRow struct {
		UserID         string                               `json:"user_id"`
		Email          string                               `json:"email"`
		Owner          bool                                 `json:"owner"`
		Accepted       bool                                 `json:"accepted"`
		PortalAppRoles map[types.PortalAppID]types.RoleName `json:"portal_application_roles"`
	}
)

var (
	errNoRoleName                = errors.New("error no role name set")
	errInvalidRoleName           = errors.New("error invalid role name set")
	errPayPlanDoesntExist        = errors.New("error pay plan '%s' does not exist")
	errAccountDoesntExist        = errors.New("error account does not exist for account ID '%s'")
	errAccountUserDoesntExist    = errors.New("error user ID '%s' does not exist for portal app ID '%s'")
	errCannotDeleteOwner         = errors.New("error cannot delete user ID '%s' for account ID '%s' because this user is the current account owner")
	errCreateNoAccountID         = errors.New("error must provide account ID when creating user")
	errCreateNoPortalAppID       = errors.New("error must provide portal app ID when creating user")
	errTransferNoAccountID       = errors.New("error must provide account ID when transferring user")
	errTransferNoPortalAppID     = errors.New("error must provide portal app ID when transferring user")
	errAcceptNoPortalAppID       = errors.New("error must provide portal app ID when accepting user")
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

	users = make(map[types.UserID]types.AccountUserAccess)

	for _, user := range userRows {
		if user.UserID != "" {
			users[types.UserID(user.UserID)] = types.AccountUserAccess{
				UserID:         types.UserID(user.UserID),
				Email:          types.Email(user.Email),
				Owner:          user.Owner,
				Accepted:       user.Accepted,
				PortalAppRoles: user.PortalAppRoles,
			}
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

	err = pg.validateWriteAccountInput(ctx, creatorID, account)
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

	userEmail, err := qtx.GetUserEmail(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	// Account creator becomes Account OWNER
	_, err = qtx.InsertAccountUserAccess(ctx, InsertAccountUserAccessParams{
		AccountID: createdAccount.ID,
		UserID:    creatorID,
		Email:     userEmail,
		RoleName:  types.RoleOwner,
		Owner:     true,
		Accepted:  true,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Assign OWNER to returned Account struct
	account.Users = map[types.UserID]types.AccountUserAccess{
		types.UserID(creatorID): {
			UserID:   types.UserID(creatorID),
			Email:    types.Email(userEmail),
			Owner:    true,
			Accepted: true,
		},
	}

	return &account, nil
}

// validateWriteAccountInput validates the input to create a new Account
func (pg *PostgresDriver) validateWriteAccountInput(ctx context.Context, creatorID types.UserID, account types.Account) error {
	planExists, err := pg.CheckPlanTypeExists(ctx, account.PlanType)
	if err != nil {
		return err
	}
	if !planExists {
		return fmt.Errorf(errPayPlanDoesntExist.Error(), account.PlanType)
	}

	userExists, err := pg.CheckUserExists(ctx, creatorID)
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

/* ----- postgresdriver Account Update Methods ----- */

// UpdateAccount updates a single Account in the database's PlanType field
func (pg *PostgresDriver) UpdateAccount(ctx context.Context, update types.UpdateAccount, updatedAt time.Time) (*types.Account, error) {
	err := pg.UpdateAccountFields(ctx, UpdateAccountFieldsParams{
		ID:        update.AccountID,
		PlanType:  update.PlanType,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return nil, err
	}

	accountResult, err := pg.SelectAccount(ctx, update.AccountID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, fmt.Errorf(errAccountDoesntExist.Error(), update.AccountID)
		}
		return nil, err
	}

	accountData := &SelectAccountsRow{
		ID:                      accountResult.ID,
		PlanType:                accountResult.PlanType,
		PartnerChainIDs:         accountResult.PartnerChainIDs,
		PartnerThroughputLimit:  accountResult.PartnerThroughputLimit,
		PartnerApplicationLimit: accountResult.PartnerApplicationLimit,
		CovalentAPIKeyFree:      accountResult.CovalentAPIKeyFree,
		CovalentAPIKeyPaid:      accountResult.CovalentAPIKeyPaid,
		Users:                   accountResult.Users,
		CreatedAt:               accountResult.CreatedAt,
		UpdatedAt:               accountResult.UpdatedAt,
		Deleted:                 accountResult.Deleted,
		DeletedAt:               accountResult.DeletedAt,
	}
	account, err := accountData.toAccount()
	if err != nil {
		return nil, err
	}

	return account, nil
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
func (pg *PostgresDriver) WriteAccountUser(ctx context.Context, createAccountUser types.CreateAccountUserAccess, createdAt time.Time) (types.UserID, error) {
	err := pg.validateWriteAccountUserInput(ctx, createAccountUser)
	if err != nil {
		return "", err
	}

	// determine if user for a given email already exists
	userID, err := pg.CheckUserIDFromEmail(ctx, createAccountUser.Email)
	if err != nil {
		switch err {

		case sql.ErrNoRows:
			// user with provided email does not exist in DB so create a new User and AccountUserAccess entry
			createdUserID, err := pg.writeAccountUserAccessNoUser(ctx, createAccountUser, createdAt)
			if err != nil {
				return "", err
			}

			return createdUserID, nil

		default:
			return "", err
		}
	}

	// user with provided email already exists in DB so create a new AccountUserAccess entry
	createdUserID, err := pg.writeAccountUserAccess(ctx, userID, createAccountUser, createdAt)
	if err != nil {
		return "", err
	}

	return createdUserID, nil
}

// validateWriteAccountUserInput validates the input to create a new AccountUserAccess row
func (pg *PostgresDriver) validateWriteAccountUserInput(ctx context.Context, createAccountUser types.CreateAccountUserAccess) error {
	if createAccountUser.Email == "" {
		return errNoEmail
	}
	if !createAccountUser.Email.IsValid() {
		return fmt.Errorf(errInvalidEmail.Error(), createAccountUser.Email)
	}
	if createAccountUser.PortalAppID == "" {
		return errCreateNoPortalAppID
	}
	if createAccountUser.AccountID == "" {
		return errCreateNoAccountID
	}

	accountExists, err := pg.CheckAccountExists(ctx, createAccountUser.AccountID)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf(errAccountDoesntExist.Error(), createAccountUser.AccountID)
	}

	portalAppExists, err := pg.CheckPortalAppExists(ctx, createAccountUser.PortalAppID)
	if err != nil {
		return err
	}
	if !portalAppExists {
		return fmt.Errorf(errPortalAppDoesntExist.Error(), createAccountUser.PortalAppID)
	}

	return nil
}

// writeAccountUserAccessNoUser creates a new User in the database and then creates a new AccountUserAccess for that user
// Called when a user is invited to a new team but does not yet have a Portal Account for the provided email
func (pg *PostgresDriver) writeAccountUserAccessNoUser(
	ctx context.Context,
	createAccountUser types.CreateAccountUserAccess,
	createdAt time.Time,
) (types.UserID, error) {
	id, err := pg.generateID(ctx)
	if err != nil {
		return "", err
	}
	userID := types.UserID(id)

	params := InsertAccountUserAccessNoUserParams{
		ID:                  userID,
		AccountID:           createAccountUser.AccountID,
		PortalApplicationID: createAccountUser.PortalAppID,
		Email:               createAccountUser.Email,
		RoleName:            createAccountUser.RoleName,
		CreatedAt:           createdAt,
		UpdatedAt:           createdAt,
	}

	createdUserID, err := pg.InsertAccountUserAccessNoUser(ctx, params)
	if err != nil {
		return "", err
	}

	return createdUserID, nil
}

// writeAccountUserAccessNoUser creates a new AccountUserAccess row for an existing user
// Called when an existing Portal user is invited to a new team
func (pg *PostgresDriver) writeAccountUserAccess(
	ctx context.Context,
	userID types.UserID,
	createAccountUser types.CreateAccountUserAccess,
	createdAt time.Time,
) (types.UserID, error) {
	params := InsertAccountUserAccessParams{
		UserID:              userID,
		AccountID:           createAccountUser.AccountID,
		PortalApplicationID: string(createAccountUser.PortalAppID),
		Email:               createAccountUser.Email,
		RoleName:            createAccountUser.RoleName,
		Owner:               false,
		Accepted:            false,
		CreatedAt:           createdAt,
		UpdatedAt:           createdAt,
	}

	createdUserID, err := pg.InsertAccountUserAccess(ctx, params)
	if err != nil {
		return "", err
	}

	return createdUserID, nil
}

/* ----- postgresdriver AccountUserAccess Update Methods ----- */

// SetAccountUserRole updates the role for an existing AccountUserAccess row. If transferring ownership the account owner becomes an admin.
func (pg *PostgresDriver) SetAccountUserRole(ctx context.Context, updateAccountUser types.UpdateAccountUserRole, updatedAt time.Time) error {
	err := pg.validateSetAccountUserRoleInput(ctx, updateAccountUser)
	if err != nil {
		return err
	}

	// If transferring ownership of an account the former OWNER becomes an ADMIN for all the account's PortalApps
	if updateAccountUser.RoleName == types.RoleOwner {
		tx, err := pg.DB.Begin()
		if err != nil {
			return err
		}
		defer func() { _ = tx.Rollback() }()

		qtx := pg.WithTx(tx)

		err = qtx.UpdateOldAccountOwnerToAdmin(ctx, updateAccountUser.AccountID)
		if err != nil {
			return err
		}

		err = qtx.UpdateNewAccountOwner(ctx, UpdateNewAccountOwnerParams{
			AccountID:  updateAccountUser.AccountID,
			NewOwnerID: string(updateAccountUser.UserID),
			CreatedAt:  updatedAt,
			UpdatedAt:  updatedAt,
		})
		if err != nil {
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		return nil
	}

	// If transferring ownership of an account the former OWNER becomes an ADMIN for all the account's PortalApps
	err = pg.UpdateAccountUserRole(ctx, UpdateAccountUserRoleParams{
		PortalApplicationID: updateAccountUser.PortalAppID,
		UserID:              updateAccountUser.UserID,
		RoleName:            updateAccountUser.RoleName,
		UpdatedAt:           updatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

// validateSetAccountUserRoleInput validates the input to update an existing AccountUserAccess role
func (pg *PostgresDriver) validateSetAccountUserRoleInput(ctx context.Context, updateAccountUser types.UpdateAccountUserRole) error {
	if updateAccountUser.RoleName == "" {
		return errNoRoleName
	}
	if !updateAccountUser.RoleName.IsValid() {
		return errInvalidRoleName
	}
	if updateAccountUser.PortalAppID == "" {
		return errTransferNoPortalAppID
	}
	if updateAccountUser.AccountID == "" {
		return errTransferNoAccountID
	}

	// If transferring OWNER role
	if updateAccountUser.RoleName == types.RoleOwner {

		// cannot transfer OWNER to a user who has not accepted their invite
		acceptedParams := CheckAccountUserAcceptedParams{UserID: updateAccountUser.UserID, PortalApplicationID: updateAccountUser.PortalAppID}
		userAccepted, err := pg.CheckAccountUserAccepted(ctx, acceptedParams)
		if err != nil {
			return err
		}
		if !userAccepted {
			return fmt.Errorf(errCannotTransferNotAccepted.Error(), updateAccountUser.UserID, updateAccountUser.AccountID)
		}

		return nil
	}

	// If transferring to non-OWNER role
	portalAppExists, err := pg.CheckPortalAppExists(ctx, updateAccountUser.PortalAppID)
	if err != nil {
		return err
	}
	if !portalAppExists {
		return fmt.Errorf(errPortalAppDoesntExist.Error(), updateAccountUser.PortalAppID)
	}

	existsParams := CheckAccountUserExistsParams{UserID: updateAccountUser.UserID, PortalApplicationID: updateAccountUser.PortalAppID}
	accountUserExists, err := pg.CheckAccountUserExists(ctx, existsParams)
	if err != nil {
		return err
	}
	if !accountUserExists {
		return fmt.Errorf(errAccountUserDoesntExist.Error(), updateAccountUser.UserID, updateAccountUser.AccountID)
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
		PortalApplicationID: acceptAccountUser.PortalAppID,
		UserID:              acceptAccountUser.UserID,
		ProviderUserID:      acceptAccountUser.ProviderUserID,
		Type:                acceptAccountUser.AuthProviderType,
		Provider:            acceptAccountUser.AuthProviderType.Provider(),
		Federated:           acceptAccountUser.AuthProviderType.IsFederated(),
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
	if acceptAccountUser.AuthProviderType == "" {
		return errNoAuthProviderType
	}
	if !acceptAccountUser.AuthProviderType.IsValid() {
		return fmt.Errorf(errInvalidAuthProviderType.Error(), acceptAccountUser.AuthProviderType)
	}
	if acceptAccountUser.PortalAppID == "" {
		return errAcceptNoPortalAppID
	}

	existsParams := CheckAccountUserExistsParams{UserID: acceptAccountUser.UserID, PortalApplicationID: acceptAccountUser.PortalAppID}
	accountUserExists, err := qtx.CheckAccountUserExists(ctx, existsParams)
	if err != nil {
		return err
	}
	if !accountUserExists {
		return fmt.Errorf(errAccountUserDoesntExist.Error(), acceptAccountUser.UserID, acceptAccountUser.PortalAppID)
	}

	return nil
}

/* ----- postgresdriver AccountUserAccess Delete Methods ----- */

// RemoveAccountUser deletes a AccountUserAccess row for a given user and account ID.
func (pg *PostgresDriver) RemoveAccountUser(ctx context.Context, userID types.UserID, portalAppID types.PortalAppID, accountID types.AccountID) error {
	err := pg.validateRemoveAccountUserInput(ctx, userID, portalAppID, accountID)
	if err != nil {
		return err
	}

	err = pg.DeleteAccountUser(ctx, DeleteAccountUserParams{UserID: userID, PortalApplicationID: portalAppID})
	if err != nil {
		return err
	}

	return nil
}

// validateRemoveAccountUserInput validates the input to remove an AccountUserAccess row
func (pg *PostgresDriver) validateRemoveAccountUserInput(ctx context.Context, userID types.UserID, portalAppID types.PortalAppID, accountID types.AccountID) error {
	accountUserRole, err := pg.CheckAccountUserRole(ctx, CheckAccountUserRoleParams{
		UserID:    userID,
		AccountID: accountID,
	})
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return fmt.Errorf(errAccountUserDoesntExist.Error(), userID, portalAppID)
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
	portalAppRoles := make(map[types.PortalAppID]types.RoleName)
	if json.PortalApplicationID != "" && json.RoleName != "" {
		portalAppRoles[json.PortalApplicationID] = json.RoleName
	}

	return &types.AccountUserAccess{
		AccountID:      json.AccountID,
		UserID:         json.UserID,
		Owner:          json.Owner,
		Accepted:       json.Accepted,
		PortalAppRoles: portalAppRoles,
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
	OwnerID                 types.UserID      `json:"owner_id"`
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
	ID                  int32             `json:"id"`
	UserID              types.UserID      `json:"user_id"`
	AccountID           types.AccountID   `json:"account_id"`
	Owner               bool              `json:"owner"`
	PortalApplicationID types.PortalAppID `json:"portal_application_id"`
	RoleName            types.RoleName    `json:"role_name"`
	Accepted            bool              `json:"accepted"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}
