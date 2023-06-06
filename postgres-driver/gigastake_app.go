package postgresdriver

import (
	"context"
	"fmt"
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
	chainIDs := make(map[types.RelayChainID]struct{}, len(g.ChainIDs))
	for _, chainID := range g.ChainIDs {
		chainIDs[types.RelayChainID(chainID)] = struct{}{}
	}

	return &types.GigastakeApp{
		ID:       g.ID,
		AATID:    g.AATID,
		ChainIDs: chainIDs,
		Name:     g.Name,
		AAT: types.AAT{
			ID:              g.AATID,
			Gigastake:       true,
			Address:         g.Address,
			PublicKey:       g.PublicKey,
			ClientPublicKey: g.ClientPublicKey,
			Signature:       g.Signature,
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
	err := pg.validateGigastakeAppInput(ctx, gigastakeApp)
	if err != nil {
		return nil, err
	}

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

	// Don't return PrivateKey in response
	gigastakeApp.AAT.PrivateKey = ""

	return &gigastakeApp, nil
}

// upsertGigastakeApp performs either an insert or update on a GigastakeApp in the database
func (pg *PostgresDriver) upsertGigastakeApp(ctx context.Context, qtx *Queries, gigastakeApp types.GigastakeApp) error {
	chainIDs := []string{}
	for chainID := range gigastakeApp.ChainIDs {
		chainIDs = append(chainIDs, string(chainID))
	}
	err := qtx.UpsertGigastakeApp(ctx, UpsertGigastakeAppParams{
		AATID:     gigastakeApp.AATID,
		ChainIDs:  chainIDs,
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

// validateGigastakeAppInput performs all necessary data validation checks on incoming GigastakeApp data
func (pg *PostgresDriver) validateGigastakeAppInput(ctx context.Context, gigastakeApp types.GigastakeApp) error {
	for chainID := range gigastakeApp.ChainIDs {
		chainExists, err := pg.CheckChainExists(ctx, chainID)
		if err != nil {
			return err
		}
		if !chainExists {
			return fmt.Errorf(errChainDoesntExist.Error(), chainID)
		}
	}

	return nil
}

/* ----- Used by Listener ----- */
func (json dbGigastakeApp) toOutput() *types.GigastakeApp {
	return &types.GigastakeApp{
		ID:         json.ID,
		AATID:      json.AATID,
		Name:       json.Name,
		AAT:        json.AAT,
		CreatedAt:  json.CreatedAt,
		UpdatedAt:  json.UpdatedAt,
		Deleted:    json.Deleted,
		LegacyLBID: json.LegacyLBID,
	}
}

func (json dbAAT) toOutput() *types.AAT {
	return &types.AAT{
		ID:              json.ID,
		Gigastake:       json.Gigastake,
		Address:         json.Address,
		PublicKey:       json.PublicKey,
		ClientPublicKey: json.ClientPublicKey,
		Signature:       json.Signature,
		Version:         json.Version,
		PortalAppID:     json.PortalAppID,
	}
}

func (json dbChainGigastakeApp) toOutput() *types.ChainGigastakeApp {
	return &types.ChainGigastakeApp{
		GigastakeAppID: json.GigastakeAppID,
		ChainID:        json.ChainID,
	}
}

type dbGigastakeApp struct {
	ID         types.GigastakeAppID `json:"id"`
	AATID      types.ProtocolAppID  `json:"aat_id"`
	ChainID    types.RelayChainID   `json:"chain_id"`
	Name       string               `json:"name"`
	AAT        types.AAT            `json:"aat"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
	Deleted    bool                 `json:"deleted"`
	DeletedAt  time.Time            `json:"deleted_at"`
	LegacyLBID string               `json:"lb_id"`
}

type dbAAT struct {
	ID              types.ProtocolAppID `json:"id"`
	Gigastake       bool                `json:"gigastake"`
	Address         string              `json:"address"`
	PublicKey       string              `json:"public_key"`
	ClientPublicKey string              `json:"client_public_key"`
	Signature       string              `json:"signature"`
	Version         string              `json:"version"`
	PortalAppID     types.PortalAppID   `json:"portal_application_id"`
}

type dbChainGigastakeApp struct {
	GigastakeAppID types.GigastakeAppID `json:"gigastake_application_id"`
	ChainID        types.RelayChainID   `json:"chain_id"`
}
