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
	altruistDBRow struct {
		ChainID  string `json:"chain_id"`
		URL      string `json:"url"`
		Auth     string `json:"auth"`
		AuthType string `json:"auth_type"`
	}

	checkDBRow struct {
		ChainID    string `json:"chain_id"`
		Type       string `json:"type"`
		Payload    string `json:"payload"`
		ResultKey  string `json:"result_key"`
		Allowance  int32  `json:"allowance"`
		EVMChainID int32  `json:"evm_chain_id"`
	}
)

var (
	errChainCannotBeNil           = errors.New("error chain cannot be nil")
	errGigastakeAppsCannotBeEmpty = errors.New("error gigastakeApps slice cannot be empty")
	errInvalidAltruistURL         = errors.New("error altruist URL '%s' is an invalid URL")
	errInvalidDomain              = errors.New("error domain '%s' for alias '%s' is invalid")
	errChainExists                = errors.New("error chain already exists for chain ID '%s'")
	errChainDoesntExist           = errors.New("error chain does not exist for chain ID '%s'")
	errPortalAppDoesntExist       = errors.New("error portal app does not exist for ID '%s'")
	errUnmarshallingDomains       = errors.New("error unmarshalling domains: %w")
)

/* ----- postgresdriver Chain Read Methods ----- */

// ReadChains returns all Chains in the database as Chains structs
func (pg *PostgresDriver) ReadChains(ctx context.Context, options types.DriverOptions) (map[types.RelayChainID]*types.Chain, error) {
	dbChains, err := pg.SelectChains(ctx, options.IncludeDeleted)
	if err != nil {
		return nil, err
	}

	chains := make(map[types.RelayChainID]*types.Chain, len(dbChains))
	for _, dbChain := range dbChains {
		chain, err := dbChain.toChain()
		if err != nil {
			return nil, err
		}

		chains[types.RelayChainID(dbChain.ID)] = chain
	}

	return chains, nil
}

// toChain converts SelectChainsRow to Chain struct
func (c *SelectChainsRow) toChain() (*types.Chain, error) {
	altruists, err := c.toAltruists()
	if err != nil {
		return nil, err
	}
	checks, err := c.toChecks()
	if err != nil {
		return nil, err
	}
	domains, err := c.toDomains()
	if err != nil {
		return nil, err
	}

	chain := &types.Chain{
		ID:             types.RelayChainID(c.ID),
		Blockchain:     c.Blockchain,
		Description:    c.Description,
		EnforceResult:  c.EnforceResult,
		Ticker:         c.Ticker,
		AllowedMethods: c.AllowedMethods,
		Path:           c.Path.String,
		LogLimitBlocks: c.LogLimitBlocks.Int32,
		RequestTimeout: c.RequestTimeout.Int32,
		Active:         c.Active,
		Altruists:      altruists,
		Checks:         checks,
		AliasDomains:   domains,
		CreatedAt:      c.CreatedAt.UTC(),
		UpdatedAt:      c.UpdatedAt.UTC(),
	}

	return chain, nil
}

// toAltruists converts altruists from DB rows to Altruist structs
func (c *SelectChainsRow) toAltruists() (map[types.AltruistURL]types.Altruist, error) {
	var altruistRows []altruistDBRow
	if err := json.Unmarshal(c.ChainAltruists, &altruistRows); err != nil {
		return nil, err
	}

	altruists := make(map[types.AltruistURL]types.Altruist, len(altruistRows))

	for _, altruistRow := range altruistRows {
		url := types.AltruistURL(altruistRow.URL)
		altruists[url] = types.Altruist{
			URL:      types.AltruistURL(altruistRow.URL),
			Auth:     altruistRow.Auth,
			AuthType: types.ChainAuthType(altruistRow.AuthType),
		}
	}

	return altruists, nil
}

// toChecks converts checks from DB rows to Check structs
func (c *SelectChainsRow) toChecks() (map[types.ChainCheckType]types.Check, error) {
	var checkRows []checkDBRow
	if err := json.Unmarshal(c.ChainChecks, &checkRows); err != nil {
		return nil, err
	}

	checks := make(map[types.ChainCheckType]types.Check, len(checkRows))

	for _, checkRow := range checkRows {
		checkType := types.ChainCheckType(checkRow.Type)
		checks[checkType] = types.Check{
			Type:       checkType,
			Payload:    checkRow.Payload,
			ResultKey:  checkRow.ResultKey,
			Allowance:  checkRow.Allowance,
			EVMChainID: checkRow.EVMChainID,
		}
	}

	return checks, nil
}

// toDomains converts chain aliases and domains from DB rows to Check structs
func (c *SelectChainsRow) toDomains() (map[types.ChainAlias][]types.ChainDomain, error) {
	var domains map[types.ChainAlias][]types.ChainDomain
	if len(string(c.AliasDomainsMap)) > 2 { // length of empty JSON array in bytes
		if err := json.Unmarshal(c.AliasDomainsMap, &domains); err != nil {
			return nil, fmt.Errorf("%s: %w", errUnmarshallingDomains, err)
		}
	}

	return domains, nil
}

// /* ----- postgresdriver Chain Create Methods ----- */

// WriteChainAndGigastakeApp creates a single Chain along with its GigastakeApps in the database as a single transaction.
// Used specifically for running the `new-chains-ci` and adding a new Chain along with one or more GigastakeApps.
func (pg *PostgresDriver) WriteChainAndGigastakeApps(ctx context.Context, input types.NewChainInput, createdAt time.Time) (*types.NewChainInput, error) {
	err := validateWriteChainAndGigastakeApps(input)
	if err != nil {
		return nil, err
	}

	tx, err := pg.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := pg.WithTx(tx)

	chain := input.Chain
	chain.CreatedAt = createdAt
	chain.UpdatedAt = createdAt

	err = pg.upsertChain(ctx, qtx, *chain, false)
	if err != nil {
		return nil, err
	}

	for _, gigastakeApp := range input.GigastakeApps {
		protocolAppID, err := pg.generateID(ctx)
		if err != nil {
			return nil, err
		}

		gigastakeApp.AATID = types.ProtocolAppID(protocolAppID)
		gigastakeApp.AAT.ID = types.ProtocolAppID(protocolAppID)
		gigastakeApp.ChainID = chain.ID
		gigastakeApp.CreatedAt = createdAt
		gigastakeApp.UpdatedAt = createdAt

		err = pg.insertGigastakeAAT(ctx, qtx, *gigastakeApp)
		if err != nil {
			return nil, err
		}

		err = pg.upsertGigastakeApp(ctx, qtx, *gigastakeApp)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &input, nil
}

// Validate checks that the fields in NewChainInput are correct
func validateWriteChainAndGigastakeApps(input types.NewChainInput) error {
	if input.Chain == nil {
		return errChainCannotBeNil
	}

	if len(input.GigastakeApps) < 1 {
		return errGigastakeAppsCannotBeEmpty
	}

	return nil
}

// WriteChain creates a single Chain in the database
func (pg *PostgresDriver) WriteChain(ctx context.Context, chain types.Chain, createdAt time.Time) (*types.Chain, error) {
	tx, err := pg.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := pg.WithTx(tx)

	chain.CreatedAt = createdAt
	chain.UpdatedAt = createdAt

	err = pg.upsertChain(ctx, qtx, chain, false)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &chain, nil
}

// UpdateChain updates a single Chain in the database
func (pg *PostgresDriver) UpdateChain(ctx context.Context, chain types.Chain, updatedAt time.Time) error {
	tx, err := pg.DB.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := pg.WithTx(tx)

	chain.UpdatedAt = updatedAt

	err = pg.upsertChain(ctx, qtx, chain, true)
	if err != nil {
		return err
	}

	err = pg.removeUnusedChainRows(ctx, qtx, chain)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// upsertChain performs either an insert or update on a chain in the database
func (pg *PostgresDriver) upsertChain(ctx context.Context, qtx *Queries, chain types.Chain, update bool) error {
	err := pg.validateChainInput(ctx, qtx, chain, update)
	if err != nil {
		return err
	}

	createdChainID, err := qtx.UpsertChain(ctx, UpsertChainParams{
		ID:             chain.ID,
		Blockchain:     chain.Blockchain,
		Description:    chain.Description,
		EnforceResult:  chain.EnforceResult,
		Ticker:         chain.Ticker,
		Path:           newSQLNullString(chain.Path),
		RequestTimeout: newSQLNullInt32(chain.RequestTimeout, true),
		LogLimitBlocks: newSQLNullInt32(chain.LogLimitBlocks, true),
		AllowedMethods: chain.AllowedMethods,
		CreatedAt:      chain.CreatedAt,
		UpdatedAt:      chain.UpdatedAt,
	})
	if err != nil {
		return err
	}

	for _, altruist := range chain.Altruists {
		err := qtx.UpsertChainAltruist(ctx, UpsertChainAltruistParams{
			ChainID:   createdChainID,
			URL:       altruist.URL,
			AuthType:  altruist.AuthType,
			Auth:      newSQLNullString(altruist.Auth),
			CreatedAt: chain.CreatedAt,
			UpdatedAt: chain.UpdatedAt,
		})
		if err != nil {
			return err
		}
	}
	for checkType, check := range chain.Checks {
		err := qtx.UpsertChainCheck(ctx, UpsertChainCheckParams{
			ChainID:    createdChainID,
			Type:       checkType,
			Payload:    newSQLNullString(check.Payload),
			ResultKey:  newSQLNullString(check.ResultKey),
			Allowance:  newSQLNullInt32(check.Allowance, false),
			EVMChainID: newSQLNullInt32(check.EVMChainID, false),
			CreatedAt:  chain.CreatedAt,
			UpdatedAt:  chain.UpdatedAt,
		})
		if err != nil {
			return err
		}
	}
	for alias, domains := range chain.AliasDomains {
		err := qtx.UpsertChainAliasDomains(ctx, UpsertChainAliasDomainsParams{
			ChainID:   createdChainID,
			Alias:     alias,
			Domains:   domains,
			UpdatedAt: chain.UpdatedAt,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// validateChainInput performs all necessary data validation checks on incoming Chain data, either for insert or update
func (pg *PostgresDriver) validateChainInput(ctx context.Context, qtx *Queries, chain types.Chain, update bool) error {
	for url := range chain.Altruists {
		if !url.IsValid() {
			return fmt.Errorf(errInvalidAltruistURL.Error(), url)
		}
	}
	for alias, domains := range chain.AliasDomains {
		for _, domain := range domains {
			if !domain.IsValid() {
				return fmt.Errorf(errInvalidDomain.Error(), domain, alias)
			}
		}
	}

	chainExists, err := qtx.CheckChainExists(ctx, chain.ID)
	if err != nil {
		return err
	}
	switch {
	case !update && chainExists:
		return fmt.Errorf(errChainExists.Error(), chain.ID)
	case update && !chainExists:
		return fmt.Errorf(errChainDoesntExist.Error(), chain.ID)
	}

	return nil
}

// removeUnusedChainRows removes chain subtables (altruists, redirects or checks)
// on a Chain update if they are not present in the update data.
// For example:
// - Chain.Redirects = [<REDIRECT_1>, <REDIRECT_2>] - all redirects except these two will be deleted
// - Chain.Redirects = [] - all redirects for the chain will be deleted
// - Chain.Redirects = nil - no changes to chain redirects
func (pg *PostgresDriver) removeUnusedChainRows(ctx context.Context, qtx *Queries, chain types.Chain) error {
	if chain.Altruists != nil {
		deleteAltruistParams := DeleteUnusedChainAltruistsParams{ChainID: chain.ID}
		for _, altruist := range chain.Altruists {
			deleteAltruistParams.URLs = append(deleteAltruistParams.URLs, string(altruist.URL))
		}
		err := qtx.DeleteUnusedChainAltruists(ctx, deleteAltruistParams)
		if err != nil {
			return err
		}
	}

	if chain.Checks != nil {
		deleteCheckParams := DeleteUnusedChainChecksParams{ChainID: chain.ID}
		for checkType := range chain.Checks {
			deleteCheckParams.Types = append(deleteCheckParams.Types, checkType)
		}
		err := qtx.DeleteUnusedChainChecks(ctx, deleteCheckParams)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg *PostgresDriver) SetChainActiveStatus(ctx context.Context, chainID types.RelayChainID, active bool, updatedAt time.Time) (bool, error) {
	chainExists, err := pg.CheckChainExists(ctx, chainID)
	if err != nil {
		return false, err
	}
	if !chainExists {
		return false, fmt.Errorf(errChainDoesntExist.Error(), chainID)
	}

	params := UpdateChainActiveParams{ID: chainID, Active: active, UpdatedAt: updatedAt}

	activeStatus, err := pg.UpdateChainActive(ctx, params)
	if err != nil {
		return false, err
	}

	return activeStatus, nil
}

/* ----- Used by Listener ----- */
func (json dbChain) toOutput() *types.Chain {
	return &types.Chain{
		ID:             json.ID,
		Blockchain:     json.Blockchain,
		Description:    json.Description,
		EnforceResult:  json.EnforceResult,
		Path:           json.Path,
		Ticker:         json.Ticker,
		AllowedMethods: json.AllowedMethods,
		LogLimitBlocks: json.LogLimitBlocks,
		RequestTimeout: json.RequestTimeout,
		Active:         json.Active,
		CreatedAt:      json.CreatedAt,
		UpdatedAt:      json.UpdatedAt,
	}
}

func (json dbChainAltruist) toOutput() *types.Altruist {
	return &types.Altruist{
		ChainID:  json.ChainID,
		URL:      json.URL,
		Auth:     json.Auth,
		AuthType: json.AuthType,
	}
}

func (json dbChainCheck) toOutput() *types.Check {
	return &types.Check{
		ChainID:    json.ChainID,
		Type:       json.Type,
		Payload:    json.Payload,
		ResultKey:  json.ResultKey,
		Allowance:  json.Allowance,
		EVMChainID: json.EVMChainID,
	}
}

func (json dbChainAliasDomains) toOutput() *types.AliasDomains {
	return &types.AliasDomains{
		ChainID: json.ChainID,
		Alias:   json.Alias,
		Domains: json.Domains,
	}
}

type dbChain struct {
	ID                       types.RelayChainID  `json:"id"`
	Blockchain               string              `json:"blockchain"`
	Description              string              `json:"description"`
	EnforceResult            string              `json:"enforce_result"`
	Ticker                   string              `json:"ticker"`
	Path                     string              `json:"path"`
	RequestTimeout           int32               `json:"request_timeout"`
	LogLimitBlocks           int32               `json:"log_limit_blocks"`
	ChainAliases             []string            `json:"chain_aliases"`
	AllowedMethods           []string            `json:"allowed_methods"`
	GigastakeRedirectDomains []types.ChainDomain `json:"gigastake_redirect_domains"`
	Active                   bool                `json:"active"`
	CreatedAt                time.Time           `json:"created_at"`
	UpdatedAt                time.Time           `json:"updated_at"`
	Deleted                  bool                `json:"deleted"`
	DeletedAt                time.Time           `json:"deleted_at"`
}

type dbChainAltruist struct {
	ID        int32               `json:"id"`
	ChainID   types.RelayChainID  `json:"chain_id"`
	URL       types.AltruistURL   `json:"url"`
	Auth      string              `json:"auth"`
	AuthType  types.ChainAuthType `json:"auth_type"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type dbChainCheck struct {
	ID         int32                `json:"id"`
	ChainID    types.RelayChainID   `json:"chain_id"`
	Type       types.ChainCheckType `json:"type"`
	Payload    string               `json:"payload"`
	ResultKey  string               `json:"result_key"`
	Allowance  int32                `json:"allowance"`
	EVMChainID int32                `json:"evm_chain_id"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

type dbChainAliasDomains struct {
	ChainID types.RelayChainID  `json:"chain_id"`
	Alias   types.ChainAlias    `json:"alias"`
	Domains []types.ChainDomain `json:"domains"`
}
