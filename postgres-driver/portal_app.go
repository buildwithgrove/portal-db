package postgresdriver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pokt-foundation/portal-db/types"
)

type (
	whitelistDBRow struct {
		Type         types.WhitelistType `json:"type"`
		Value        string              `json:"value"`
		BlockchainID string              `json:"chain_id"`
	}
)

var (
	errUnmarshallingWhitelists    = errors.New("error unmarshalling whitelists")
	errUnmarshallingNotifications = errors.New("error unmarshalling notifications")
)

/* ----- postgresdriver PortalApp Read Methods ----- */

// ReadApplications returns all Applications in the database as PortalApp structs
func (pg *PostgresDriver) ReadPortalApps(ctx context.Context) (map[types.PortalAppID]*types.PortalApp, error) {
	dbPortalApps, err := pg.SelectPortalApplications(ctx)
	if err != nil {
		return nil, err
	}

	portalApps := make(map[types.PortalAppID]*types.PortalApp, len(dbPortalApps))
	for _, dbPortalApp := range dbPortalApps {
		portalApp, err := dbPortalApp.toPortalApp()
		if err != nil {
			return nil, err
		}

		portalApps[types.PortalAppID(dbPortalApp.ID)] = portalApp
	}

	return portalApps, nil
}

// toPortalApp converts PortalApp SELECT output to PortalApp struct
func (a *SelectPortalApplicationsRow) toPortalApp() (*types.PortalApp, error) {
	var appWhitelists types.Whitelists
	if len(string(a.Whitelists)) > 2 {
		whitelists, err := a.toWhitelists()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errUnmarshallingWhitelists, err)
		}
		appWhitelists = whitelists
	}

	var notifications map[types.NotificationType]types.AppNotification
	if len(string(a.Notifications)) > 2 {
		if err := json.Unmarshal(a.Notifications, &notifications); err != nil {
			return nil, fmt.Errorf("%s: %w", errUnmarshallingNotifications, err)
		}
	}

	// TODO remove legacy fields when migration to V2 schema complete
	legacyFields := types.LegacyFields{
		ApplicationIDs:     a.ApplicationIDs,
		CustomLimit:        a.CustomLimit.Int32,
		RequestTimeout:     a.RequestTimeout.Int32,
		GigastakeRedirect:  a.GigastakeRedirect.Bool,
		FirstDateSurpassed: a.FirstDateSurpassed.Time.UTC(),
		StickyOptions: types.StickyOptions{
			Duration:      a.Duration.String,
			StickyOrigins: a.Origins,
			StickyMax:     int(a.StickyMax.Int32),
			Stickiness:    a.Stickiness.Bool,
		},
	}

	return &types.PortalApp{
		ID:        types.PortalAppID(a.ID),
		AccountID: types.AccountID(a.AccountID),
		Name:      a.Name,
		Gigastake: a.Gigastake,
		Staked:    a.Staked,
		AAT: types.AAT{
			Address:         a.Address.String,
			PublicKey:       a.PublicKey.String,
			ClientPublicKey: a.ClientPublicKey.String,
			PrivateKey:      a.PrivateKey.String,
			Signature:       a.Signature.String,
			Version:         a.Version.String,
		},
		Settings: types.Settings{
			Environment:       types.Environment(a.Environment.Environment),
			SecretKey:         a.SecretKey.String,
			SecretKeyRequired: a.SecretKeyRequired.Bool,
		},
		Notifications: notifications,
		Whitelists:    appWhitelists,
		CreatedAt:     a.CreatedAt.UTC(),
		UpdatedAt:     a.UpdatedAt.UTC(),
		// TODO remove legacy fields when migration to V2 schema complete
		LegacyFields: legacyFields,
	}, nil
}

// toWhitelists converts whitelists from DB rows to map-based PortalApp.Whitelists struct
func (a *SelectPortalApplicationsRow) toWhitelists() (types.Whitelists, error) {
	whitelists := types.Whitelists{
		Origins:     make(map[types.Origin]struct{}),
		UserAgents:  make(map[types.UserAgent]struct{}),
		Blockchains: make(map[types.ChainID]struct{}),
		Contracts:   make(map[types.ChainID]map[types.Contract]struct{}),
		Methods:     make(map[types.ChainID]map[types.Method]struct{}),
	}

	var whitelistRows []whitelistDBRow
	if err := json.Unmarshal(a.Whitelists, &whitelistRows); err != nil {
		return whitelists, err
	}

	for _, wl := range whitelistRows {
		switch wl.Type {

		case types.WhitelistTypeBlockchains:
			whitelists.Blockchains[types.ChainID(wl.Value)] = struct{}{}
		case types.WhitelistTypeOrigins:
			whitelists.Origins[types.Origin(wl.Value)] = struct{}{}
		case types.WhitelistTypeUserAgents:
			whitelists.UserAgents[types.UserAgent(wl.Value)] = struct{}{}

		case types.WhitelistTypeContracts:
			if _, ok := whitelists.Contracts[types.ChainID(wl.BlockchainID)]; !ok {
				whitelists.Contracts[types.ChainID(wl.BlockchainID)] = make(map[types.Contract]struct{})
			}
			whitelists.Contracts[types.ChainID(wl.BlockchainID)][types.Contract(wl.Value)] = struct{}{}
		case types.WhitelistTypeMethods:
			if _, ok := whitelists.Methods[types.ChainID(wl.BlockchainID)]; !ok {
				whitelists.Methods[types.ChainID(wl.BlockchainID)] = make(map[types.Method]struct{})
			}
			whitelists.Methods[types.ChainID(wl.BlockchainID)][types.Method(wl.Value)] = struct{}{}
		}
	}

	return whitelists, nil
}

/* ----- postgresdriver PortalApp Create Methods ----- */

// WritePortalApp creates a single PortalApp in the database, including its AAT and Settings rows
// TEMP - also create its legacy StickinessOptions table (TODO remove when V2 migration completed)
func (pg *PostgresDriver) WritePortalApp(ctx context.Context, portalApp types.PortalApp, createdAt time.Time) (*types.PortalApp, error) {
	id, err := generateRandomID()
	if err != nil {
		return nil, err
	}

	portalApp.ID = types.PortalAppID(id)
	portalApp.CreatedAt = createdAt
	portalApp.UpdatedAt = createdAt

	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	qtx := pg.WithTx(tx)

	_, err = qtx.InsertPortalApplication(ctx, InsertPortalApplicationParams{
		ID:        portalApp.ID,
		AccountID: int64(portalApp.AccountID),
		Name:      portalApp.Name,
		Gigastake: portalApp.Gigastake,
		Staked:    portalApp.Staked,

		// TODO remove legacy fields when migration to V2 schema complete
		ApplicationIDs:     (portalApp.LegacyFields.ApplicationIDs),
		RequestTimeout:     newSQLNullInt32(portalApp.LegacyFields.RequestTimeout, true),
		CustomLimit:        newSQLNullInt32(portalApp.LegacyFields.CustomLimit, true),
		GigastakeRedirect:  newSQLNullBool(&portalApp.LegacyFields.GigastakeRedirect),
		FirstDateSurpassed: newSQLNullTime(portalApp.LegacyFields.FirstDateSurpassed),

		CreatedAt: portalApp.CreatedAt,
		UpdatedAt: portalApp.UpdatedAt,
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	_, err = qtx.InsertPortalApplicationAAT(ctx, InsertPortalApplicationAATParams{
		ApplicationID:   portalApp.ID,
		Address:         portalApp.AAT.Address,
		PublicKey:       portalApp.AAT.PublicKey,
		PrivateKey:      portalApp.AAT.PrivateKey,
		ClientPublicKey: portalApp.AAT.ClientPublicKey,
		Signature:       portalApp.AAT.Signature,
		Version:         portalApp.AAT.Version,
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	_, err = qtx.InsertPortalApplicationSetting(ctx, InsertPortalApplicationSettingParams{
		ApplicationID:     portalApp.ID,
		Environment:       portalApp.Settings.Environment,
		SecretKey:         portalApp.Settings.SecretKey,
		SecretKeyRequired: portalApp.Settings.SecretKeyRequired,
		MonthlyRelayLimit: portalApp.Settings.MonthlyRelayLimit,
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// TODO remove legacy fields when migration to V2 schema complete
	_, err = qtx.InsertStickinessOption(ctx, InsertStickinessOptionParams{
		ApplicationID: portalApp.ID,
		Duration:      newSQLNullString(portalApp.LegacyFields.StickyOptions.Duration),
		StickyMax:     newSQLNullInt32(int32(portalApp.LegacyFields.StickyOptions.StickyMax), true),
		Stickiness:    newSQLNullBool(&portalApp.LegacyFields.StickyOptions.Stickiness),
		Origins:       portalApp.LegacyFields.StickyOptions.StickyOrigins,
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &portalApp, nil
}

/* ----- postgresdriver PortalApp Update Methods ----- */

// UpdatePortalApp updates a single PortalApp in the database: Name field and its Notifications, Whitelists and Settings
func (pg *PostgresDriver) UpdatePortalApp(ctx context.Context, update types.UpdatePortalApp, updatedAt time.Time) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	qtx := pg.WithTx(tx)

	if update.Name != "" {
		err := qtx.UpdatePortalAppName(ctx, UpdatePortalAppNameParams{
			ID: update.AppID, Name: update.Name, UpdatedAt: updatedAt,
		})
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if update.Settings != nil {
		err := pg.updateSettings(ctx, tx, qtx, update, updatedAt)
		if err != nil {
			return err
		}
	}

	// TODO uncomment once notifications figured out
	// if update.Notifications != nil && len(update.Notifications) > 0 {
	// 	err := pg.updateNotifications(ctx, tx, qtx, update, updatedAt)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	if update.Whitelists != nil {
		err := pg.updateWhitelists(ctx, tx, qtx, update, updatedAt)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// updateSettings updates the PortalApp's settings row in the portal_application_settings table
func (pg *PostgresDriver) updateSettings(ctx context.Context, tx *sql.Tx, qtx *Queries, update types.UpdatePortalApp, updatedAt time.Time) error {
	updateSettings := UpdatePortalAppSettingsParams{
		ApplicationID:     update.AppID,
		SecretKey:         update.Settings.SecretKey,
		SecretKeyRequired: update.Settings.SecretKeyRequired,
		MonthlyRelayLimit: update.Settings.MonthlyRelayLimit,
		Environment:       update.Settings.Environment,
		FavoritedChainIDs: update.Settings.FavoritedChainIDs,
		UpdatedAt:         newSQLNullTime(updatedAt),
	}

	err := qtx.UpdatePortalAppSettings(ctx, updateSettings)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// TODO uncomment once notifications figured out
// updateNotifications updates the PortalApp's notifications rows in the portal_application_notifications table
// func (pg *PostgresDriver) updateNotifications(ctx context.Context, tx *sql.Tx, qtx *Queries, update types.UpdatePortalApp, updatedAt time.Time) error {
// 	updateNotifications := UpsertPortalAppNotificationsParams{ApplicationID: update.AppID, UpdatedAt: newSQLNullTime(updatedAt)}

// 	for _, appNotification := range update.Notifications {
// 		updateNotifications.Active = append(updateNotifications.Active, appNotification.Active)
// 		updateNotifications.Types = append(updateNotifications.Types, appNotification.NotificationType)
// 		updateNotifications.Destination = append(updateNotifications.Destination, appNotification.Destination)
// 		updateNotifications.Trigger = append(updateNotifications.Trigger, appNotification.Trigger)
// 		// TODO need to update here (events should be slice of slices, figure out how to make SQLC represent as such)
// 		// updateNotifications.Events = append(updateNotifications.Events, appNotification.Events)
// 	}

// 	err := qtx.UpsertPortalAppNotifications(ctx, updateNotifications)
// 	if err != nil {
// 		_ = tx.Rollback()
// 		return err
// 	}

// 	return nil
// }

// updateWhitelists updates the PortalApp's whitelists rows in the portal_application_whitelists table
func (pg *PostgresDriver) updateWhitelists(ctx context.Context, tx *sql.Tx, qtx *Queries, update types.UpdatePortalApp, updatedAt time.Time) error {
	// Map all whitelist rows from update struct to insert query params
	updateWhitelists := UpdateInsertWhitelistsParams{ApplicationID: update.AppID, CreatedAt: updatedAt}
	for _, appWhitelist := range update.Whitelists.AppWhitelists {
		for _, whitelistValue := range appWhitelist.Values {
			updateWhitelists.Types = append(updateWhitelists.Types, appWhitelist.Type)
			updateWhitelists.Values = append(updateWhitelists.Values, whitelistValue)
			updateWhitelists.ChainIDs = append(updateWhitelists.ChainIDs, "")
		}
	}
	for _, chainWhitelist := range update.Whitelists.ChainWhitelists {
		for _, blockchainValues := range chainWhitelist.Values {
			for _, whitelistValue := range blockchainValues.Values {
				updateWhitelists.Types = append(updateWhitelists.Types, chainWhitelist.Type)
				updateWhitelists.ChainIDs = append(updateWhitelists.ChainIDs, blockchainValues.ChainID)
				updateWhitelists.Values = append(updateWhitelists.Values, whitelistValue)
			}
		}
	}

	// Insert all whitelist rows for application_id in update struct that are not in DB
	err := qtx.UpdateInsertWhitelists(ctx, updateWhitelists)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Delete all whitelist rows for application_id in DB that are not in update struct
	err = qtx.UpdateDeleteWhitelists(ctx, UpdateDeleteWhitelistsParams{
		ApplicationID: updateWhitelists.ApplicationID,
		Types:         updateWhitelists.Types,
		Values:        updateWhitelists.Values,
		ChainIDs:      updateWhitelists.ChainIDs,
	})
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
