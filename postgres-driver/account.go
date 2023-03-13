package postgresdriver

import (
	"context"
	"encoding/json"
	"fmt"

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
	chainIDs := make(map[types.ChainID]struct{}, len(a.ChainIDs))
	for _, chainID := range a.ChainIDs {
		chainIDs[types.ChainID(chainID)] = struct{}{}
	}

	partnerChainIDs := make(map[types.ChainID]struct{}, len(a.PartnerChainIDs))
	for _, chainID := range a.PartnerChainIDs {
		partnerChainIDs[types.ChainID(chainID)] = struct{}{}
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
		users[types.UserID(user.ID)] = types.AccountUserAccess{
			User: types.User{
				ID:           types.UserID(user.ID),
				Email:        types.Email(user.Email),
				AuthProvider: types.AuthProviders(user.AuthProvider),
			},
			RoleName: types.RoleName(user.RoleName),
			Accepted: user.Accepted,
		}
	}

	return users, nil
}

// /* ----- postgresdriver PortalApp Create Methods ----- */

// // WritePortalApp creates a single PortalApp in the database, including its AAT and Settings rows
// // TEMP - also create its legacy StickinessOptions table (TODO remove when V2 migration completed)
// func (pg *PostgresDriver) WritePortalApp(ctx context.Context, portalApp types.PortalApp, createdAt time.Time) (*types.PortalApp, error) {
// 	id, err := generateRandomID()
// 	if err != nil {
// 		return nil, err
// 	}

// 	portalApp.ID = types.PortalAppID(id)
// 	portalApp.CreatedAt = createdAt
// 	portalApp.UpdatedAt = createdAt

// 	tx, err := pg.db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}

// 	qtx := pg.WithTx(tx)

// 	_, err = qtx.InsertPortalApplication(ctx, InsertPortalApplicationParams{
// 		ID:        portalApp.ID,
// 		AccountID: int32(portalApp.AccountID),
// 		Name:      portalApp.Name,
// 		Gigastake: portalApp.Gigastake,
// 		Staked:    portalApp.Staked,
// 		CreatedAt: portalApp.CreatedAt,
// 		UpdatedAt: portalApp.UpdatedAt,
// 		// TODO remove legacy fields when migration to V2 schema complete
// 		ApplicationIDs:     (portalApp.LegacyFields.ApplicationIDs),
// 		RequestTimeout:     newSQLNullInt32(portalApp.LegacyFields.RequestTimeout, true),
// 		CustomLimit:        newSQLNullInt32(portalApp.LegacyFields.CustomLimit, true),
// 		GigastakeRedirect:  newSQLNullBool(&portalApp.LegacyFields.GigastakeRedirect),
// 		FirstDateSurpassed: newSQLNullTime(portalApp.LegacyFields.FirstDateSurpassed),
// 	})
// 	if err != nil {
// 		_ = tx.Rollback()
// 		return nil, err
// 	}

// 	_, err = qtx.InsertPortalApplicationAAT(ctx, InsertPortalApplicationAATParams{
// 		ApplicationID:   portalApp.ID,
// 		Address:         portalApp.AAT.Address,
// 		PublicKey:       portalApp.AAT.PublicKey,
// 		PrivateKey:      portalApp.AAT.PrivateKey,
// 		ClientPublicKey: portalApp.AAT.ClientPublicKey,
// 		Signature:       portalApp.AAT.Signature,
// 		Version:         portalApp.AAT.Version,
// 	})
// 	if err != nil {
// 		_ = tx.Rollback()
// 		return nil, err
// 	}

// 	_, err = qtx.InsertPortalApplicationSetting(ctx, InsertPortalApplicationSettingParams{
// 		ApplicationID:     portalApp.ID,
// 		Environment:       portalApp.Settings.Environment,
// 		SecretKey:         portalApp.Settings.SecretKey,
// 		SecretKeyRequired: portalApp.Settings.SecretKeyRequired,
// 		MonthlyRelayLimit: portalApp.Settings.MonthlyRelayLimit,
// 	})
// 	if err != nil {
// 		_ = tx.Rollback()
// 		return nil, err
// 	}

// 	// TODO remove legacy fields when migration to V2 schema complete
// 	_, err = qtx.InsertStickinessOption(ctx, InsertStickinessOptionParams{
// 		LbID:       portalApp.ID,
// 		Duration:   newSQLNullString(portalApp.LegacyFields.StickyOptions.Duration),
// 		StickyMax:  newSQLNullInt32(int32(portalApp.LegacyFields.StickyOptions.StickyMax), true),
// 		Stickiness: newSQLNullBool(&portalApp.LegacyFields.StickyOptions.Stickiness),
// 		Origins:    portalApp.LegacyFields.StickyOptions.StickyOrigins,
// 	})
// 	if err != nil {
// 		_ = tx.Rollback()
// 		return nil, err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &portalApp, nil
// }

// /* ----- postgresdriver PortalApp Update Methods ----- */

// // UpdatePortalApp updates a single PortalApp in the database: Name field and its Notifications, Whitelists and Settings
// func (pg *PostgresDriver) UpdatePortalApp(ctx context.Context, update types.UpdatePortalApp, updatedAt time.Time) error {
// 	tx, err := pg.db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	qtx := pg.WithTx(tx)

// 	if update.Name != "" {
// 		err := qtx.UpdatePortalAppName(ctx, UpdatePortalAppNameParams{
// 			ID: update.AppID, Name: update.Name, UpdatedAt: updatedAt,
// 		})
// 		if err != nil {
// 			_ = tx.Rollback()
// 			return err
// 		}
// 	}
// 	if update.Settings != nil {
// 		err := pg.updateSettings(ctx, tx, qtx, update, updatedAt)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if update.Notifications != nil && len(update.Notifications) > 0 {
// 		err := pg.updateNotifications(ctx, tx, qtx, update, updatedAt)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if update.Whitelists != nil {
// 		err := pg.updateWhitelists(ctx, tx, qtx, update, updatedAt)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // updateSettings updates the PortalApp's settings row in the portal_application_settings table
// func (pg *PostgresDriver) updateSettings(ctx context.Context, tx *sql.Tx, qtx *Queries, update types.UpdatePortalApp, updatedAt time.Time) error {
// 	updateSettings := UpdatePortalAppSettingsParams{
// 		ApplicationID:     update.AppID,
// 		SecretKey:         update.Settings.SecretKey,
// 		SecretKeyRequired: update.Settings.SecretKeyRequired,
// 		MonthlyRelayLimit: update.Settings.MonthlyRelayLimit,
// 		Environment:       update.Settings.Environment,
// 		FavoritedChainIDs: update.Settings.FavoritedChainIDs,
// 		UpdatedAt:         newSQLNullTime(updatedAt),
// 	}

// 	err := qtx.UpdatePortalAppSettings(ctx, updateSettings)
// 	if err != nil {
// 		_ = tx.Rollback()
// 		return err
// 	}

// 	return nil
// }

// // updateNotifications updates the PortalApp's notifications rows in the portal_application_notifications table
// func (pg *PostgresDriver) updateNotifications(ctx context.Context, tx *sql.Tx, qtx *Queries, update types.UpdatePortalApp, updatedAt time.Time) error {
// 	for _, appNotification := range update.Notifications {
// 		updateNotification := UpdateUpsertPortalAppNotificationParams{
// 			ApplicationID: update.AppID,
// 			Type:          appNotification.NotificationType,
// 			Active:        appNotification.Active,
// 			Destination:   appNotification.Destination,
// 			Trigger:       appNotification.Trigger,
// 			Events:        appNotification.Events,
// 			UpdatedAt:     newSQLNullTime(updatedAt),
// 		}

// 		// Upsert notification row for application_id & type if active: true in update struct
// 		err := qtx.UpdateUpsertPortalAppNotification(ctx, updateNotification)
// 		if err != nil {
// 			_ = tx.Rollback()
// 			return err
// 		}

// 		// Delete notification row for application_id & type in DB if active: false in update struct
// 		if !appNotification.Active {
// 			err := qtx.UpdateDeletePortalAppNotification(ctx, UpdateDeletePortalAppNotificationParams{
// 				ApplicationID: update.AppID, Type: appNotification.NotificationType,
// 			})
// 			if err != nil {
// 				_ = tx.Rollback()
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// // updateWhitelists updates the PortalApp's whitelists rows in the portal_application_whitelists table
// func (pg *PostgresDriver) updateWhitelists(ctx context.Context, tx *sql.Tx, qtx *Queries, update types.UpdatePortalApp, updatedAt time.Time) error {
// 	// Map all whitelist rows from update struct to insert query params
// 	updateWhitelists := UpdateInsertWhitelistsParams{ApplicationID: update.AppID, CreatedAt: updatedAt}
// 	for _, appWhitelist := range update.Whitelists.AppWhitelists {
// 		for _, whitelistValue := range appWhitelist.Values {
// 			updateWhitelists.Types = append(updateWhitelists.Types, appWhitelist.Type)
// 			updateWhitelists.Values = append(updateWhitelists.Values, whitelistValue)
// 			updateWhitelists.ChainIDs = append(updateWhitelists.ChainIDs, "")
// 		}
// 	}
// 	for _, chainWhitelist := range update.Whitelists.ChainWhitelists {
// 		for _, blockchainValues := range chainWhitelist.Values {
// 			for _, whitelistValue := range blockchainValues.Values {
// 				updateWhitelists.Types = append(updateWhitelists.Types, chainWhitelist.Type)
// 				updateWhitelists.ChainIDs = append(updateWhitelists.ChainIDs, blockchainValues.ChainID)
// 				updateWhitelists.Values = append(updateWhitelists.Values, whitelistValue)
// 			}
// 		}
// 	}

// 	// Insert all whitelist rows for application_id in update struct that are not in DB
// 	err := qtx.UpdateInsertWhitelists(ctx, updateWhitelists)
// 	if err != nil {
// 		_ = tx.Rollback()
// 		return err
// 	}

// 	// Delete all whitelist rows for application_id in DB that are not in update struct
// 	err = qtx.UpdateDeleteWhitelists(ctx, UpdateDeleteWhitelistsParams{
// 		ApplicationID: updateWhitelists.ApplicationID,
// 		Types:         updateWhitelists.Types,
// 		Values:        updateWhitelists.Values,
// 		ChainIDs:      updateWhitelists.ChainIDs,
// 	})
// 	if err != nil {
// 		_ = tx.Rollback()
// 		return err
// 	}

// 	return nil
// }

// // UpdatePortalAppsFirstDateSurpassed updates multiple PortalApps' LegacyFields.FirstDateSurpassed fields
// // TODO legacy method - determine if still needed and remove if not when V2 migration completed
// func (pg *PostgresDriver) UpdatePortalAppsFirstDateSurpassed(ctx context.Context, update *types.UpdateFirstDateSurpassed) error {
// 	params := UpdateFirstDatesSurpassedParams{
// 		ApplicationIDs:     update.PortalAppIDs,
// 		FirstDateSurpassed: newSQLNullTime(update.FirstDateSurpassed),
// 	}

// 	err := pg.UpdateFirstDatesSurpassed(ctx, params)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// /* ----- postgresdriver PortalApp Delete Methods ----- */

// // SetPortalAppDeleted updates a single PortalApp in the database's Deleted field to true
// func (pg *PostgresDriver) SetPortalAppDeleted(ctx context.Context, portalAppID types.PortalAppID, deletedAt time.Time) error {
// 	params := DeletePortalAppParams{ID: portalAppID, DeletedAt: newSQLNullTime(deletedAt)}

// 	err := pg.DeletePortalApp(ctx, params)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
