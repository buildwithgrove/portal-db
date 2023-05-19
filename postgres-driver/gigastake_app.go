package postgresdriver

import (
	"context"

	"github.com/pokt-foundation/portal-db/v2/types"
)

/* ----- postgresdriver GigastakeApp Read Methods ----- */

// ReadGigastakeApps returns all gigastake applications in the database as GigastakeApp structs
func (pg *PostgresDriver) ReadGigastakeApps(ctx context.Context, options types.DriverOptions) (map[types.ProtocolAppID]*types.GigastakeApp, error) {
	dbGigastakeApps, err := pg.SelectGigastakeApplications(ctx)
	if err != nil {
		return nil, err
	}

	gigastakeApps := make(map[types.ProtocolAppID]*types.GigastakeApp, len(dbGigastakeApps))
	for _, dbGigastakeApp := range dbGigastakeApps {
		gigastakeApp, err := dbGigastakeApp.toGigastakeApp()
		if err != nil {
			return nil, err
		}

		gigastakeApps[dbGigastakeApp.AATID] = gigastakeApp
	}

	return gigastakeApps, nil
}

// toGigastakeApp converts GigastakeApp SELECT output to GigastakeApp struct
func (g *SelectGigastakeApplicationsRow) toGigastakeApp() (*types.GigastakeApp, error) {
	aat := types.AAT{
		ID:              g.AATID,
		Gigastake:       g.Gigastake,
		Address:         g.Address,
		PublicKey:       g.PublicKey,
		ClientPublicKey: g.ClientPublicKey,
		Signature:       g.Signature,
		PrivateKey:      g.PrivateKey.String,
		Version:         g.Version,
	}

	return &types.GigastakeApp{
		AATID:      g.AATID,
		ChainID:    types.RelayChainID(g.ChainID),
		ChainAlias: g.ChainAlias,
		Name:       g.Name,
		AAT:        aat,
		CreatedAt:  g.CreatedAt.UTC(),
		UpdatedAt:  g.UpdatedAt.UTC(),
		Deleted:    g.Deleted,
	}, nil
}
