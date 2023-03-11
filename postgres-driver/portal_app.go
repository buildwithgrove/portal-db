package postgresdriver

import (
	"context"
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

// ReadApplications returns all Applications in the database
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

	// if a.ID == types.PortalAppID("test_app_3487u329rfn23f9") {
	// 	PrettyString("ROWS HERE in toPortalApp", appWhitelists)
	// }

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
		CreatedAt:     a.CreatedAt.Time.UTC(),
		UpdatedAt:     a.UpdatedAt.Time.UTC(),
		// TODO remove legacy fields when migration to V2 schema complete
		LegacyFields: legacyFields,
	}, nil
}

func (a *SelectPortalApplicationsRow) toWhitelists() (types.Whitelists, error) {
	whitelists := types.Whitelists{
		Origins:     make(map[types.Origin]struct{}),
		UserAgents:  make(map[types.UserAgent]struct{}),
		Blockchains: make(map[types.BlockchainID]struct{}),
		Contracts:   make(map[types.BlockchainID]map[types.Contract]struct{}),
		Methods:     make(map[types.BlockchainID]map[types.Method]struct{}),
	}

	var whitelistRows []whitelistDBRow
	if err := json.Unmarshal(a.Whitelists, &whitelistRows); err != nil {
		return whitelists, err
	}

	// if a.ID == types.PortalAppID("test_app_3487u329rfn23f9") {
	// 	PrettyString("ROWS HERE in toWhitelists", whitelistRows)
	// }

	for _, wl := range whitelistRows {
		switch wl.Type {

		case types.WhitelistTypeBlockchains:
			whitelists.Blockchains[types.BlockchainID(wl.Value)] = struct{}{}

		case types.WhitelistTypeOrigins:
			whitelists.Origins[types.Origin(wl.Value)] = struct{}{}

		case types.WhitelistTypeUserAgents:
			whitelists.UserAgents[types.UserAgent(wl.Value)] = struct{}{}

		case types.WhitelistTypeContracts:
			if _, ok := whitelists.Contracts[types.BlockchainID(wl.BlockchainID)]; !ok {
				whitelists.Contracts[types.BlockchainID(wl.BlockchainID)] = make(map[types.Contract]struct{})
			}
			whitelists.Contracts[types.BlockchainID(wl.BlockchainID)][types.Contract(wl.Value)] = struct{}{}

		case types.WhitelistTypeMethods:
			if _, ok := whitelists.Methods[types.BlockchainID(wl.BlockchainID)]; !ok {
				whitelists.Methods[types.BlockchainID(wl.BlockchainID)] = make(map[types.Method]struct{})
			}
			whitelists.Methods[types.BlockchainID(wl.BlockchainID)][types.Method(wl.Value)] = struct{}{}
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

		CreatedAt: newSQLNullTime(portalApp.CreatedAt),
		UpdatedAt: newSQLNullTime(portalApp.UpdatedAt),
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
// TEMP - also update its legacy StickinessOptions table (TODO remove when V2 migration completed)
func (pg *PostgresDriver) UpdatePortalApp(ctx context.Context, update types.UpdatePortalApp, updatedAt time.Time) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	qtx := pg.WithTx(tx)

	// TODO -> add rest of updates

	if update.Whitelists != nil {
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
					updateWhitelists.ChainIDs = append(updateWhitelists.ChainIDs, blockchainValues.BlockchainID)
					updateWhitelists.Values = append(updateWhitelists.Values, whitelistValue)
				}
			}
		}

		// Insert all whitelist rows in update struct that don't exist in DB
		err := qtx.UpdateInsertWhitelists(ctx, updateWhitelists)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		// Delete all whitelist rows in DB that don't exists in update struct
		err = qtx.UpdateDeletePortalAppWhitelists(ctx, UpdateDeletePortalAppWhitelistsParams{
			ApplicationID: updateWhitelists.ApplicationID,
			Types:         updateWhitelists.Types,
			Values:        updateWhitelists.Values,
			ChainIDs:      updateWhitelists.ChainIDs,
		})
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
