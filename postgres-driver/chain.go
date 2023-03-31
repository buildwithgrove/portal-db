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

	gigastakeRedirectDBRow struct {
		ChainID              string `json:"chain_id"`
		AccountID            string `json:"account_id"`
		Alias                string `json:"alias"`
		Domain               string `json:"domain"`
		LegacyLoadBalancerID string `json:"lb_id"`
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
	errInvalidAltruistURL    = errors.New("error altruist URL '%s' is an invalid URL")
	errInvalidRedirectDomain = errors.New("error redirect domain '%s' is an invalid domain")
	errChainExists           = errors.New("error chain already exists for chain ID '%s'")
	errChainDoesntExist      = errors.New("error chain does not exist for chain ID '%s'")
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
	redirects, err := c.toGigastakeRedirects()
	if err != nil {
		return nil, err
	}
	checks, err := c.toChecks()
	if err != nil {
		return nil, err
	}

	chain := &types.Chain{
		ID:             types.RelayChainID(c.ID),
		Blockchain:     c.Blockchain,
		Description:    c.Description,
		EnforceResult:  c.EnforceResult,
		Ticker:         c.Ticker,
		ChainAliases:   c.ChainAliases,
		AllowedMethods: c.AllowedMethods,
		Path:           c.Path.String,
		LogLimitBlocks: c.LogLimitBlocks.Int32,
		RequestTimeout: c.RequestTimeout.Int32,
		Active:         c.Active,
		Altruists:      altruists,
		Redirects:      redirects,
		Checks:         checks,
		CreatedAt:      c.CreatedAt.UTC(),
		UpdatedAt:      c.UpdatedAt.UTC(),
	}

	return chain, nil
}

// toAltruists converts altruists from DB rows to Altruist structs
func (c *SelectChainsRow) toAltruists() ([]types.Altruist, error) {
	var altruistRows []altruistDBRow
	if err := json.Unmarshal(c.ChainAltruists, &altruistRows); err != nil {
		return nil, err
	}

	altruists := make([]types.Altruist, len(altruistRows))

	for i, altruistRow := range altruistRows {
		altruists[i] = types.Altruist{
			URL:      types.AltruistURL(altruistRow.URL),
			Auth:     altruistRow.Auth,
			AuthType: types.ChainAuthType(altruistRow.AuthType),
		}
	}

	return altruists, nil
}

// toGigastakeRedirects converts gigastake redirects from DB rows to GigastakeRedirect structs
func (c *SelectChainsRow) toGigastakeRedirects() ([]types.GigastakeRedirect, error) {
	var redirectRows []gigastakeRedirectDBRow
	if err := json.Unmarshal(c.ChainGigastakeRedirects, &redirectRows); err != nil {
		return nil, err
	}

	redirects := make([]types.GigastakeRedirect, len(redirectRows))

	for i, redirectRow := range redirectRows {
		redirects[i] = types.GigastakeRedirect{
			AccountID: types.AccountID(redirectRow.AccountID),
			Domain:    types.RedirectDomain(redirectRow.Domain),
			Alias:     redirectRow.Alias,

			// TODO - remove when v2 migration finished
			// LegacyLoadBalancerID is the load balancer ID that the account was migrated from
			LegacyLoadBalancerID: redirectRow.LegacyLoadBalancerID,
		}
	}

	return redirects, nil
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

// /* ----- postgresdriver Chain Create Methods ----- */

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
		ChainAliases:   chain.ChainAliases,
		AllowedMethods: chain.AllowedMethods,
		CreatedAt:      chain.CreatedAt,
		UpdatedAt:      chain.CreatedAt,
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
			UpdatedAt: chain.CreatedAt,
		})
		if err != nil {
			return err
		}
	}
	for _, redirect := range chain.Redirects {
		err = qtx.UpsertChainGigastakeRedirect(ctx, UpsertChainGigastakeRedirectParams{
			ChainID:   createdChainID,
			AccountID: redirect.AccountID,
			Alias:     redirect.Alias,
			Domain:    redirect.Domain,
			CreatedAt: chain.CreatedAt,
			UpdatedAt: chain.CreatedAt,
			// TODO remove legacy fields when migration to V2 schema complete
			LbID: redirect.LegacyLoadBalancerID,
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
			UpdatedAt:  chain.CreatedAt,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// validateChainInput performs all necessary data validation checks on incoming Chain data, either for insert or update
func (pg *PostgresDriver) validateChainInput(ctx context.Context, qtx *Queries, chain types.Chain, update bool) error {
	for _, altruist := range chain.Altruists {
		if !altruist.URL.IsValid() {
			return fmt.Errorf(errInvalidAltruistURL.Error(), altruist.URL)
		}
	}

	for _, redirect := range chain.Redirects {
		if !redirect.Domain.IsValid() {
			return fmt.Errorf(errInvalidRedirectDomain.Error(), redirect.Domain)
		}
		accountExists, err := qtx.CheckAccountExists(ctx, redirect.AccountID)
		if err != nil {
			return err
		}
		if !accountExists {
			return fmt.Errorf(errAccountDoesntExist.Error(), redirect.AccountID)
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

	if chain.Redirects != nil {
		deleteRedirectParams := DeleteUnusedChainGigastakeRedirectsParams{ChainID: chain.ID}
		for _, redirect := range chain.Redirects {
			deleteRedirectParams.AccountIDs = append(deleteRedirectParams.AccountIDs, string(redirect.AccountID))
		}
		err := qtx.DeleteUnusedChainGigastakeRedirects(ctx, deleteRedirectParams)
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

func (pg *PostgresDriver) RemoveGigastakeRedirect(ctx context.Context, chainID types.RelayChainID, accountID types.AccountID, domain types.RedirectDomain) error {
	redirectExists, err := pg.CheckRedirectExists(ctx, CheckRedirectExistsParams{
		ChainID:   chainID,
		AccountID: accountID,
		Domain:    domain,
	})
	if err != nil {
		return err
	}
	if !redirectExists {
		return fmt.Errorf("Redirect with chain ID '%s', account ID '%s' and domain '%s' doesn't exist", chainID, accountID, domain)
	}

	err = pg.DeleteGigastakeRedirect(ctx, DeleteGigastakeRedirectParams{ChainID: chainID, AccountID: accountID, Domain: domain})
	if err != nil {
		return err
	}

	return nil
}

/* ----- Used by Listener ----- */
func (json Chain) toOutput() *types.Chain {
	return &types.Chain{
		ID:             json.ID,
		Blockchain:     json.Blockchain,
		Description:    json.Description,
		EnforceResult:  json.EnforceResult,
		Path:           json.Path.String,
		Ticker:         json.Ticker,
		ChainAliases:   json.ChainAliases,
		AllowedMethods: json.AllowedMethods,
		LogLimitBlocks: json.LogLimitBlocks.Int32,
		RequestTimeout: json.RequestTimeout.Int32,
		Active:         json.Active,
		CreatedAt:      json.CreatedAt,
		UpdatedAt:      json.UpdatedAt,
	}
}

func (json ChainAltruist) toOutput() *types.Altruist {
	return &types.Altruist{
		ChainID:  json.ChainID,
		URL:      json.URL,
		Auth:     json.Auth.String,
		AuthType: json.AuthType,
	}
}

func (json ChainCheck) toOutput() *types.Check {
	return &types.Check{
		ChainID:    json.ChainID,
		Type:       json.Type,
		Payload:    json.Payload.String,
		ResultKey:  json.ResultKey.String,
		Allowance:  json.Allowance.Int32,
		EVMChainID: json.EVMChainID.Int32,
	}
}

func (r ChainGigastakeRedirect) toOutput() *types.GigastakeRedirect {
	return &types.GigastakeRedirect{
		ChainID:              r.ChainID,
		AccountID:            r.AccountID,
		Domain:               r.Domain,
		Alias:                r.Alias,
		LegacyLoadBalancerID: r.LbID,
	}
}
