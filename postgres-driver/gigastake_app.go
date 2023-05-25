package postgresdriver

import (
	"context"
	"time"

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
	return &types.GigastakeApp{
		AATID:   g.AATID,
		ChainID: types.RelayChainID(g.ChainID),
		Name:    g.Name,
		AAT: types.AAT{
			ID:              g.AATID,
			Gigastake:       true,
			Address:         g.Address,
			PublicKey:       g.PublicKey,
			ClientPublicKey: g.ClientPublicKey,
			Signature:       g.Signature,
			PrivateKey:      g.PrivateKey.String,
			Version:         g.Version,
		},
		CreatedAt: g.CreatedAt.UTC(),
		UpdatedAt: g.UpdatedAt.UTC(),
		Deleted:   g.Deleted,

		// TODO remove legacy field when migration to V2 schema complete
		LegacyLBID: g.LbID,
	}, nil
}

/* ----- postgresdriver Chain Create Methods ----- */

// WriteGigastakeApp creates a single GigastakeApp in the database
func (pg *PostgresDriver) WriteGigastakeApp(ctx context.Context, gigastakeApp types.GigastakeApp, createdAt time.Time) (*types.GigastakeApp, error) {
	protocolAppID, err := pg.generateID(ctx)
	if err != nil {
		return nil, err
	}

	gigastakeApp.AATID = types.ProtocolAppID(protocolAppID)
	gigastakeApp.AAT.ID = types.ProtocolAppID(protocolAppID)
	gigastakeApp.CreatedAt = createdAt
	gigastakeApp.UpdatedAt = createdAt

	tx, err := pg.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := pg.WithTx(tx)

	err = pg.insertGigastakeAAT(ctx, qtx, gigastakeApp)
	if err != nil {
		return nil, err
	}

	err = pg.upsertGigastakeApp(ctx, qtx, gigastakeApp)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &gigastakeApp, nil
}

// upsertGigastakeApp performs either an insert or update on a GigastakeApp in the database
func (pg *PostgresDriver) upsertGigastakeApp(ctx context.Context, qtx *Queries, gigastakeApp types.GigastakeApp) error {
	err := qtx.UpsertGigastakeApp(ctx, UpsertGigastakeAppParams{
		AATID:     gigastakeApp.AATID,
		ChainID:   gigastakeApp.ChainID,
		Name:      gigastakeApp.Name,
		CreatedAt: gigastakeApp.CreatedAt,
		UpdatedAt: gigastakeApp.UpdatedAt,
		// TODO remove legacy field when migration to V2 schema complete
		LbID: gigastakeApp.LegacyLBID,
	})
	if err != nil {
		return err
	}

	return nil
}

// upsertAAT performs an upsert operation on the AAT table in the database
func (pg *PostgresDriver) insertGigastakeAAT(ctx context.Context, qtx *Queries, gigastakeApp types.GigastakeApp) error {
	err := qtx.InsertGigastakeAAT(ctx, InsertGigastakeAATParams{
		ID:              gigastakeApp.AATID,
		Address:         gigastakeApp.AAT.Address,
		PublicKey:       gigastakeApp.AAT.PublicKey,
		ClientPublicKey: gigastakeApp.AAT.ClientPublicKey,
		PrivateKey:      newSQLNullString(gigastakeApp.AAT.PrivateKey),
		Signature:       gigastakeApp.AAT.Signature,
		Version:         gigastakeApp.AAT.Version,
	})
	if err != nil {
		return err
	}

	return nil
}
