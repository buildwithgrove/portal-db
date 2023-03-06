package postgresdriver

import (
	"context"
	"encoding/json"

	"github.com/pokt-foundation/portal-db/types"
)

type (
	WhitelistDBRow struct {
		Type         string `json:"type"`
		Value        string `json:"value"`
		BlockchainID string `json:"blockchain_id"`
	}
)

/* ReadApplications returns all Applications in the database */
func (p *PostgresDriver) ReadPortalApps(ctx context.Context) ([]*types.PortalApp, error) {
	dbPortalApps, err := p.SelectPortalApplications(ctx)
	if err != nil {
		return nil, err
	}

	var portalApps []*types.PortalApp
	for _, dbPortalApp := range dbPortalApps {
		portalApp, err := dbPortalApp.toApplication()
		if err != nil {
			return nil, err
		}

		portalApps = append(portalApps, portalApp)
	}

	return portalApps, nil
}

func (p *SelectPortalApplicationsRow) toApplication() (*types.PortalApp, error) {
	appWhitelists, err := p.toWhitelists()
	if err != nil {
		return nil, err
	}

	return &types.PortalApp{
		ID:   p.ID,
		Name: p.Name,

		Whitelists: appWhitelists,
		CreatedAt:  p.CreatedAt.Time,
		UpdatedAt:  p.UpdatedAt.Time,
	}, nil
}

func (p *SelectPortalApplicationsRow) toWhitelists() (types.Whitelists, error) {
	var whitelists types.Whitelists

	if len(string(p.Whitelists)) > 0 {
		var whitelistRows []WhitelistDBRow
		if err := json.Unmarshal(p.Whitelists, &whitelistRows); err != nil {
			return whitelists, err
		}

		for _, wl := range whitelistRows {
			switch wl.Type {
			case "blockchains":
				if _, ok := whitelists.Blockchains[types.BlockchainID(wl.Value)]; !ok {
					whitelists.Blockchains = make(map[types.BlockchainID]struct{})
				}
				whitelists.Blockchains[types.BlockchainID(wl.Value)] = struct{}{}
			case "origins":
				if _, ok := whitelists.Origins[types.Origin(wl.Value)]; !ok {
					whitelists.Origins = make(map[types.Origin]struct{})
				}
				whitelists.Origins[types.Origin(wl.Value)] = struct{}{}
			case "userAgents":
				if _, ok := whitelists.UserAgents[types.UserAgent(wl.Value)]; !ok {
					whitelists.UserAgents = make(map[types.UserAgent]struct{})
				}
				whitelists.UserAgents[types.UserAgent(wl.Value)] = struct{}{}
			case "contracts":
				if _, ok := whitelists.Contracts[types.BlockchainID(wl.BlockchainID)]; !ok {
					whitelists.Contracts = make(map[types.BlockchainID]map[types.Contract]struct{})
					whitelists.Contracts[types.BlockchainID(wl.BlockchainID)] = make(map[types.Contract]struct{})
				}
				whitelists.Contracts[types.BlockchainID(wl.BlockchainID)][types.Contract(wl.Value)] = struct{}{}
			case "methods":
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
