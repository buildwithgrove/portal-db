package postgresdriver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pokt-foundation/portal-db/types"
)

type (
	whitelistDBRow struct {
		Type         types.WhitelistType `json:"type"`
		Value        string              `json:"value"`
		BlockchainID string              `json:"blockchain_id"`
	}
)

var (
	errUnmarshallingWhitelists    = errors.New("error unmarshalling whitelists")
	errUnmarshallingNotifications = errors.New("error unmarshalling notifications")
)

/* ReadApplications returns all Applications in the database */
func (p *PostgresDriver) ReadPortalApps(ctx context.Context) (map[types.PortalAppID]*types.PortalApp, error) {
	dbPortalApps, err := p.SelectPortalApplications(ctx)
	if err != nil {
		return nil, err
	}

	portalApps := make(map[types.PortalAppID]*types.PortalApp, len(dbPortalApps))
	for _, dbPortalApp := range dbPortalApps {
		portalApp, err := dbPortalApp.toApplication()
		if err != nil {
			return nil, err
		}

		portalApps[types.PortalAppID(dbPortalApp.ID)] = portalApp
	}

	return portalApps, nil
}

func (p *SelectPortalApplicationsRow) toApplication() (*types.PortalApp, error) {
	appWhitelists, err := p.toWhitelists()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshallingWhitelists, err)
	}

	var notifications map[types.NotificationType]types.AppNotification
	if err := json.Unmarshal(p.Notifications, &notifications); err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshallingNotifications, err)
	}

	// TODO remove legacy fields when migration to V2 schema complete
	legacyFields := types.LegacyFields{
		ApplicationID:      p.ApplicationID.String,
		CustomLimit:        int(p.CustomLimit.Int32),
		RequestTimeout:     int(p.RequestTimeout.Int32),
		GigastakeRedirect:  p.GigastakeRedirect.Bool,
		FirstDateSurpassed: p.FirstDateSurpassed.Time.UTC(),
		StickyOptions: types.StickyOptions{
			Duration:      p.Duration.String,
			StickyOrigins: p.Origins,
			StickyMax:     int(p.StickyMax.Int32),
			Stickiness:    p.Stickiness.Bool,
		},
	}

	return &types.PortalApp{
		ID:        types.PortalAppID(p.ID),
		AccountID: types.AccountID(p.AccountID),
		Name:      p.Name,
		Gigastake: p.Gigastake,
		Staked:    p.Staked,
		AAT: types.AAT{
			Address:         p.Address.String,
			PublicKey:       p.PublicKey.String,
			ClientPublicKey: p.PrivateKey.String,
			PrivateKey:      p.ClientPublicKey.String,
			Signature:       p.Signature.String,
			Version:         p.Version.String,
		},
		Settings: types.Settings{
			Environment:       types.Environment(p.Environment.Environment),
			SecretKey:         p.SecretKey.String,
			SecretKeyRequired: p.SecretKeyRequired.Bool,
		},
		Notifications: notifications,
		Whitelists:    appWhitelists,
		CreatedAt:     p.CreatedAt.Time.UTC(),
		UpdatedAt:     p.UpdatedAt.Time.UTC(),
		// TODO remove legacy fields when migration to V2 schema complete
		LegacyFields: legacyFields,
	}, nil
}

func (p *SelectPortalApplicationsRow) toWhitelists() (types.Whitelists, error) {
	var whitelists types.Whitelists

	if len(string(p.Whitelists)) > 0 {
		var whitelistRows []whitelistDBRow
		if err := json.Unmarshal(p.Whitelists, &whitelistRows); err != nil {
			return whitelists, err
		}

		for _, wl := range whitelistRows {
			switch wl.Type {
			case types.WhitelistTypeBlockchains:
				if _, ok := whitelists.Blockchains[types.BlockchainID(wl.Value)]; !ok {
					whitelists.Blockchains = make(map[types.BlockchainID]struct{})
				}
				whitelists.Blockchains[types.BlockchainID(wl.Value)] = struct{}{}
			case types.WhitelistTypeOrigins:
				if _, ok := whitelists.Origins[types.Origin(wl.Value)]; !ok {
					whitelists.Origins = make(map[types.Origin]struct{})
				}
				whitelists.Origins[types.Origin(wl.Value)] = struct{}{}
			case types.WhitelistTypeUserAgents:
				if _, ok := whitelists.UserAgents[types.UserAgent(wl.Value)]; !ok {
					whitelists.UserAgents = make(map[types.UserAgent]struct{})
				}
				whitelists.UserAgents[types.UserAgent(wl.Value)] = struct{}{}
			case types.WhitelistTypeContracts:
				if _, ok := whitelists.Contracts[types.BlockchainID(wl.BlockchainID)]; !ok {
					whitelists.Contracts = make(map[types.BlockchainID]map[types.Contract]struct{})
					whitelists.Contracts[types.BlockchainID(wl.BlockchainID)] = make(map[types.Contract]struct{})
				}
				whitelists.Contracts[types.BlockchainID(wl.BlockchainID)][types.Contract(wl.Value)] = struct{}{}
			case types.WhitelistTypeMethods:
				if _, ok := whitelists.Methods[types.BlockchainID(wl.BlockchainID)]; !ok {
					whitelists.Methods = make(map[types.BlockchainID]map[types.Method]struct{})
					whitelists.Methods[types.BlockchainID(wl.BlockchainID)] = make(map[types.Method]struct{})
				}
				whitelists.Methods[types.BlockchainID(wl.BlockchainID)][types.Method(wl.Value)] = struct{}{}
			}
		}
	}

	return whitelists, nil
}
