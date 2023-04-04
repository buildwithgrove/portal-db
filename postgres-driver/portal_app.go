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
	whitelistDBRow struct {
		Type         types.WhitelistType `json:"type"`
		Value        string              `json:"value"`
		BlockchainID string              `json:"chain_id"`
	}
	aatDBRow struct {
		ID              types.ProtocolAppID `json:"id"`
		ApplicationID   types.PortalAppID   `json:"application_id"`
		Address         string              `json:"address"`
		PublicKey       string              `json:"public_key"`
		ClientPublicKey string              `json:"client_public_key"`
		PrivateKey      string              `json:"private_key"`
		Signature       string              `json:"signature"`
		Version         string              `json:"version"`
	}
)

var (
	errUnmarshallingWhitelists    = errors.New("error unmarshalling whitelists")
	errUnmarshallingNotifications = errors.New("error unmarshalling notifications")
	errUnmarshallingAATs          = errors.New("error unmarshalling AATs")
)

/* ----- postgresdriver PortalApp Read Methods ----- */

// ReadApplications returns all Applications in the database as PortalApp structs
func (pg *PostgresDriver) ReadPortalApps(ctx context.Context, options types.DriverOptions) (map[types.PortalAppID]*types.PortalApp, error) {
	dbPortalApps, err := pg.SelectPortalApplications(ctx, options.IncludeDeleted)
	if err != nil {
		return nil, err
	}

	portalApps := make(map[types.PortalAppID]*types.PortalApp, len(dbPortalApps))
	for _, dbPortalApp := range dbPortalApps {
		portalApp, err := dbPortalApp.toPortalApp()
		if err != nil {
			return nil, err
		}

		portalApps[dbPortalApp.ID] = portalApp
	}

	return portalApps, nil
}

// toPortalApp converts PortalApp SELECT output to PortalApp struct
func (a *SelectPortalApplicationsRow) toPortalApp() (*types.PortalApp, error) {
	var appAATs map[types.ProtocolAppID]types.AAT
	if len(string(a.AATs)) > 2 { // length of empty JSON array in bytes
		aats, err := a.toAATs()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errUnmarshallingAATs, err)
		}
		appAATs = aats
	}

	var appWhitelists types.Whitelists
	if len(string(a.Whitelists)) > 2 { // length of empty JSON array in bytes
		whitelists, err := a.toWhitelists()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errUnmarshallingWhitelists, err)
		}
		appWhitelists = whitelists
	}

	var notifications map[types.NotificationType]types.AppNotification
	if len(string(a.Notifications)) > 2 { // length of empty JSON array in bytes
		if err := json.Unmarshal(a.Notifications, &notifications); err != nil {
			return nil, fmt.Errorf("%s: %w", errUnmarshallingNotifications, err)
		}
	}

	// TODO remove legacy fields when migration to V2 schema complete
	legacyFields := types.LegacyFields{
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
		Settings: types.Settings{
			Environment:       types.Environment(a.Environment.Environment),
			SecretKey:         a.SecretKey.String,
			SecretKeyRequired: a.SecretKeyRequired.Bool,
		},
		AATs:          appAATs,
		Whitelists:    appWhitelists,
		Notifications: notifications,
		CreatedAt:     a.CreatedAt.UTC(),
		UpdatedAt:     a.UpdatedAt.UTC(),
		Deleted:       a.Deleted,
		// TODO remove legacy fields when migration to V2 schema complete
		LegacyFields: legacyFields,
	}, nil
}

// toAATs converts AATs from DB rows to map-based PortalApp.AATs struct
func (a *SelectPortalApplicationsRow) toAATs() (map[types.ProtocolAppID]types.AAT, error) {
	var dbAATs map[types.ProtocolAppID]aatDBRow

	if err := json.Unmarshal(a.AATs, &dbAATs); err != nil {
		return nil, err
	}

	aats := make(map[types.ProtocolAppID]types.AAT, len(dbAATs))

	for protocolAppID, dbAAT := range dbAATs {
		aats[protocolAppID] = types.AAT{
			ID:              protocolAppID,
			Address:         dbAAT.Address,
			PublicKey:       dbAAT.PublicKey,
			ClientPublicKey: dbAAT.ClientPublicKey,
			PrivateKey:      dbAAT.PrivateKey,
			Signature:       dbAAT.Signature,
			Version:         dbAAT.Version,
		}
	}

	return aats, nil
}

// toWhitelists converts whitelists from DB rows to map-based PortalApp.Whitelists struct
func (a *SelectPortalApplicationsRow) toWhitelists() (types.Whitelists, error) {
	whitelists := types.Whitelists{
		Origins:     make(map[types.Origin]struct{}),
		UserAgents:  make(map[types.UserAgent]struct{}),
		Blockchains: make(map[types.RelayChainID]struct{}),
		Contracts:   make(map[types.RelayChainID]map[types.Contract]struct{}),
		Methods:     make(map[types.RelayChainID]map[types.Method]struct{}),
	}

	var whitelistRows []whitelistDBRow
	if err := json.Unmarshal(a.Whitelists, &whitelistRows); err != nil {
		return whitelists, err
	}

	for _, wl := range whitelistRows {
		switch wl.Type {

		case types.WhitelistTypeBlockchains:
			whitelists.Blockchains[types.RelayChainID(wl.Value)] = struct{}{}
		case types.WhitelistTypeOrigins:
			whitelists.Origins[types.Origin(wl.Value)] = struct{}{}
		case types.WhitelistTypeUserAgents:
			whitelists.UserAgents[types.UserAgent(wl.Value)] = struct{}{}

		case types.WhitelistTypeContracts:
			if _, ok := whitelists.Contracts[types.RelayChainID(wl.BlockchainID)]; !ok {
				whitelists.Contracts[types.RelayChainID(wl.BlockchainID)] = make(map[types.Contract]struct{})
			}
			whitelists.Contracts[types.RelayChainID(wl.BlockchainID)][types.Contract(wl.Value)] = struct{}{}
		case types.WhitelistTypeMethods:
			if _, ok := whitelists.Methods[types.RelayChainID(wl.BlockchainID)]; !ok {
				whitelists.Methods[types.RelayChainID(wl.BlockchainID)] = make(map[types.Method]struct{})
			}
			whitelists.Methods[types.RelayChainID(wl.BlockchainID)][types.Method(wl.Value)] = struct{}{}
		}
	}

	return whitelists, nil
}

/* ----- postgresdriver PortalApp Create Methods ----- */

// WritePortalApp creates a single PortalApp in the database, including its AAT and Settings rows
// TEMP - also create its legacy StickinessOptions table (TODO remove when V2 migration completed)
func (pg *PostgresDriver) WritePortalApp(ctx context.Context, portalApp types.PortalApp, aat types.AAT, createdAt time.Time) (*types.PortalApp, error) {
	portalAppID, protocolAppID, err := pg.generatePortalAppIDs(ctx)
	if err != nil {
		return nil, err
	}

	portalApp.ID = portalAppID
	portalApp.CreatedAt = createdAt
	portalApp.UpdatedAt = createdAt

	aat.ID = protocolAppID
	portalApp.AATs = map[types.ProtocolAppID]types.AAT{
		protocolAppID: aat,
	}

	tx, err := pg.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := pg.WithTx(tx)

	_, err = qtx.InsertPortalApplication(ctx, InsertPortalApplicationParams{
		ID:        portalApp.ID,
		AccountID: portalApp.AccountID,
		Name:      portalApp.Name,
		Gigastake: portalApp.Gigastake,
		Staked:    portalApp.Staked,
		CreatedAt: portalApp.CreatedAt,
		UpdatedAt: portalApp.UpdatedAt,
		// TODO remove legacy fields when migration to V2 schema complete
		RequestTimeout:     newSQLNullInt32(portalApp.LegacyFields.RequestTimeout, true),
		CustomLimit:        newSQLNullInt32(portalApp.LegacyFields.CustomLimit, true),
		GigastakeRedirect:  newSQLNullBool(&portalApp.LegacyFields.GigastakeRedirect),
		FirstDateSurpassed: newSQLNullTime(portalApp.LegacyFields.FirstDateSurpassed),
	})
	if err != nil {
		return nil, err
	}
	_, err = qtx.InsertPortalApplicationAAT(ctx, InsertPortalApplicationAATParams{
		ApplicationID:   portalApp.ID,
		ID:              aat.ID,
		Address:         aat.Address,
		PublicKey:       aat.PublicKey,
		PrivateKey:      aat.PrivateKey,
		ClientPublicKey: aat.ClientPublicKey,
		Signature:       aat.Signature,
		Version:         aat.Version,
	})
	if err != nil {
		return nil, err
	}
	_, err = qtx.InsertPortalApplicationSetting(ctx, InsertPortalApplicationSettingParams{
		ApplicationID:     portalApp.ID,
		Environment:       portalApp.Settings.Environment,
		SecretKey:         newSQLNullString(portalApp.Settings.SecretKey),
		SecretKeyRequired: newSQLNullBool(&portalApp.Settings.SecretKeyRequired),
		MonthlyRelayLimit: portalApp.Settings.MonthlyRelayLimit,
	})
	if err != nil {
		return nil, err
	}
	// TODO remove legacy fields when migration to V2 schema complete
	_, err = qtx.InsertStickinessOption(ctx, InsertStickinessOptionParams{
		LbID:       portalApp.ID,
		Duration:   newSQLNullString(portalApp.LegacyFields.StickyOptions.Duration),
		StickyMax:  newSQLNullInt32(int32(portalApp.LegacyFields.StickyOptions.StickyMax), true),
		Stickiness: newSQLNullBool(&portalApp.LegacyFields.StickyOptions.Stickiness),
		Origins:    portalApp.LegacyFields.StickyOptions.StickyOrigins,
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &portalApp, nil
}

func (pg *PostgresDriver) generatePortalAppIDs(ctx context.Context) (types.PortalAppID, types.ProtocolAppID, error) {
	portalAppID, err := pg.generateID(ctx)
	if err != nil {
		return "", "", err
	}
	protocolAppID, err := pg.generateID(ctx)
	if err != nil {
		return "", "", err
	}
	for protocolAppID == portalAppID { // Ensure the protocol app ID cannot match the portal app ID
		protocolAppID, err = pg.generateID(ctx)
		if err != nil {
			return "", "", err
		}
	}

	return types.PortalAppID(portalAppID), types.ProtocolAppID(portalAppID), nil
}

/* ----- postgresdriver PortalApp Update Methods ----- */

// UpdatePortalApp updates a single PortalApp in the database: Name field and its Notifications, Whitelists and Settings
func (pg *PostgresDriver) UpdatePortalApp(ctx context.Context, update types.UpdatePortalApp, updatedAt time.Time) error {
	tx, err := pg.DB.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := pg.WithTx(tx)

	if update.Name != "" {
		err := qtx.UpdatePortalAppName(ctx, UpdatePortalAppNameParams{
			ID: update.AppID, Name: update.Name, UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}
	}
	if update.Settings != nil {
		err := pg.updateSettings(ctx, qtx, update, updatedAt)
		if err != nil {
			return err
		}
	}
	if update.Notifications != nil && len(update.Notifications) > 0 {
		err := pg.updateNotifications(ctx, qtx, update, updatedAt)
		if err != nil {
			return err
		}
	}
	if update.Whitelists != nil {
		err := pg.updateWhitelists(ctx, qtx, update, updatedAt)
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
func (pg *PostgresDriver) updateSettings(ctx context.Context, qtx *Queries, update types.UpdatePortalApp, updatedAt time.Time) error {
	updateSettings := UpdatePortalAppSettingsParams{
		ApplicationID:     update.AppID,
		SecretKey:         newSQLNullString(update.Settings.SecretKey),
		SecretKeyRequired: newSQLNullBool(&update.Settings.SecretKeyRequired),
		MonthlyRelayLimit: update.Settings.MonthlyRelayLimit,
		Environment:       update.Settings.Environment,
		FavoritedChainIDs: update.Settings.FavoritedChainIDs,
		UpdatedAt:         updatedAt,
	}

	err := qtx.UpdatePortalAppSettings(ctx, updateSettings)
	if err != nil {
		return err
	}

	return nil
}

// updateNotifications updates the PortalApp's notifications rows in the portal_application_notifications table
func (pg *PostgresDriver) updateNotifications(ctx context.Context, qtx *Queries, update types.UpdatePortalApp, updatedAt time.Time) error {
	for _, appNotification := range update.Notifications {
		updateNotification := UpdateUpsertPortalAppNotificationParams{
			ApplicationID: update.AppID,
			Type:          appNotification.NotificationType,
			Active:        appNotification.Active,
			Destination:   appNotification.Destination,
			Trigger:       appNotification.Trigger,
			Events:        appNotification.Events,
			UpdatedAt:     updatedAt,
		}

		// Upsert notification row for application_id & type if active: true in update struct
		err := qtx.UpdateUpsertPortalAppNotification(ctx, updateNotification)
		if err != nil {
			return err
		}

		// Delete notification row for application_id & type in DB if active: false in update struct
		if !appNotification.Active {
			err := qtx.UpdateDeletePortalAppNotification(ctx, UpdateDeletePortalAppNotificationParams{
				ApplicationID: update.AppID, Type: appNotification.NotificationType,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// updateWhitelists updates the PortalApp's whitelists rows in the portal_application_whitelists table
func (pg *PostgresDriver) updateWhitelists(ctx context.Context, qtx *Queries, update types.UpdatePortalApp, updatedAt time.Time) error {
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
		return err
	}

	return nil
}

// UpdatePortalAppsFirstDateSurpassed updates multiple PortalApps' LegacyFields.FirstDateSurpassed fields
// TODO legacy method - determine if still needed and remove if not when V2 migration completed
func (pg *PostgresDriver) UpdatePortalAppsFirstDateSurpassed(ctx context.Context, update *types.UpdateFirstDateSurpassed) error {
	params := UpdateFirstDatesSurpassedParams{
		ApplicationIDs:     update.PortalAppIDs,
		FirstDateSurpassed: newSQLNullTime(update.FirstDateSurpassed),
	}

	err := pg.UpdateFirstDatesSurpassed(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

/* ----- postgresdriver PortalApp Delete Methods ----- */

// SetPortalAppDeleted updates a single PortalApp in the database's Deleted field to true
func (pg *PostgresDriver) SetPortalAppDeleted(ctx context.Context, portalAppID types.PortalAppID, deletedAt time.Time) error {
	params := DeletePortalAppParams{ID: portalAppID, DeletedAt: newSQLNullTime(deletedAt)}

	err := pg.DeletePortalApp(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

/* ----- Used by Listener ----- */
func (json PortalApplication) toOutput() *types.PortalApp {
	return &types.PortalApp{
		ID:        json.ID,
		AccountID: types.AccountID(json.AccountID),
		Name:      json.Name,
		Gigastake: json.Gigastake,
		Staked:    json.Staked,
		CreatedAt: json.CreatedAt,
		UpdatedAt: json.UpdatedAt,
		Deleted:   json.Deleted,
		LegacyFields: types.LegacyFields{
			RequestTimeout:     json.RequestTimeout.Int32,
			GigastakeRedirect:  json.GigastakeRedirect.Bool,
			FirstDateSurpassed: json.FirstDateSurpassed.Time,
			CustomLimit:        json.CustomLimit.Int32,
		},
	}
}

func (json PortalApplicationAat) toOutput() *types.AAT {
	return &types.AAT{
		AppID:           json.ApplicationID,
		Address:         json.Address,
		PublicKey:       json.PublicKey,
		ClientPublicKey: json.ClientPublicKey,
		PrivateKey:      json.PrivateKey,
		Signature:       json.Signature,
		Version:         json.Version,
	}
}

func (json PortalApplicationSetting) toOutput() *types.Settings {
	var favoritedChainIDs map[types.RelayChainID]struct{}
	if len(json.FavoritedChainIDs) != 0 {
		favoritedChainIDs = make(map[types.RelayChainID]struct{})
		for _, chainID := range json.FavoritedChainIDs {
			favoritedChainIDs[types.RelayChainID(chainID)] = struct{}{}
		}

	}

	return &types.Settings{
		AppID:             json.ApplicationID,
		Environment:       types.Environment(json.Environment),
		SecretKey:         json.SecretKey.String,
		SecretKeyRequired: json.SecretKeyRequired.Bool,
		FavoritedChainIDs: favoritedChainIDs,
		MonthlyRelayLimit: json.MonthlyRelayLimit,
	}
}

func (json PortalApplicationWhitelist) toOutput() *types.Whitelist {
	return &types.Whitelist{
		AppID:   json.ApplicationID,
		Type:    json.Type,
		Value:   json.Value,
		ChainID: types.RelayChainID(json.ChainID.String),
	}
}

func (json PortalApplicationNotification) toOutput() *types.AppNotification {
	return &types.AppNotification{
		AppID:       json.ApplicationID,
		Type:        json.Type,
		Active:      json.Active,
		Destination: json.Destination.String,
		Trigger:     json.Trigger.String,
		Events:      json.mapEvents(),
	}
}

func (json PortalApplicationNotification) mapEvents() map[types.NotificationEvent]bool {
	events := make(map[types.NotificationEvent]bool)
	for _, event := range json.Events {
		events[event] = true
	}
	return events
}

// TODO - remove when v2 migration finished
// Fields required for compatibility with the old Portal API and Services (temporary)
func (json StickinessOption) toOutput() *types.StickyOptions {
	return &types.StickyOptions{
		ID:            string(json.LbID),
		Duration:      json.Duration.String,
		StickyOrigins: json.Origins,
		StickyMax:     int(json.StickyMax.Int32),
		Stickiness:    json.Stickiness.Bool,
	}
}
