package postgresdriver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pokt-foundation/portal-db/types"
)

var (
	ErrInvalidRedirectJSON = errors.New("error: redirect JSON is invalid")
)

/* ReadBlockchains returns all blockchains in the database and marshals to types struct */
func (p *PostgresDriver) ReadBlockchains(ctx context.Context) ([]*types.Blockchain, error) {
	dbBlockchains, err := p.SelectBlockchains(ctx)
	if err != nil {
		return nil, err
	}

	var blockchains []*types.Blockchain
	for _, dbBlockchain := range dbBlockchains {
		blockchain, err := dbBlockchain.toBlockchain()
		if err != nil {
			return nil, err
		}

		blockchains = append(blockchains, blockchain)
	}

	return blockchains, nil
}

func (b *SelectBlockchainsRow) toBlockchain() (*types.Blockchain, error) {
	blockchain := types.Blockchain{
		ID:                b.BlockchainID,
		Altruist:          b.Altruist.String,
		Blockchain:        b.Blockchain.String,
		ChainID:           b.ChainID.String,
		ChainIDCheck:      b.ChainIDCheck.String,
		Description:       b.Description.String,
		EnforceResult:     b.EnforceResult.String,
		Network:           b.Network.String,
		Path:              b.Path.String,
		SyncCheck:         b.SSyncCheck.String,
		Ticker:            b.Ticker.String,
		BlockchainAliases: b.BlockchainAliases,
		LogLimitBlocks:    int(b.LogLimitBlocks.Int32),
		RequestTimeout:    int(b.RequestTimeout.Int32),
		Active:            b.Active.Bool,

		SyncCheckOptions: types.SyncCheckOptions{
			Body:      b.SBody.String,
			ResultKey: b.SResultKey.String,
			Path:      b.SPath.String,
			Allowance: int(b.SAllowance.Int32),
		},

		CreatedAt: b.CreatedAt.Time,
		UpdatedAt: b.UpdatedAt.Time,
	}

	// Unmarshal Blockchain Redirects JSON into []types.Redirects
	err := json.Unmarshal(b.Redirects, &blockchain.Redirects)
	if err != nil {
		return &types.Blockchain{}, fmt.Errorf("%w: %s", ErrInvalidRedirectJSON, err)
	}

	return &blockchain, nil
}

/* WriteBlockchain saves input Blockchain struct to the database */
func (p *PostgresDriver) WriteBlockchain(ctx context.Context, blockchain *types.Blockchain) (*types.Blockchain, error) {
	time := time.Now()
	blockchain.CreatedAt = time
	blockchain.UpdatedAt = time

	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := p.WithTx(tx)

	err = qtx.InsertBlockchain(ctx, extractInsertDBBlockchain(blockchain))
	if err != nil {
		return nil, err
	}

	syncCheckOptionsParams := extractInsertSyncCheckOptions(blockchain)
	if syncCheckOptionsParams.isNotNull() {
		err = qtx.InsertSyncCheckOptions(ctx, syncCheckOptionsParams)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return blockchain, nil
}

func extractInsertDBBlockchain(blockchain *types.Blockchain) InsertBlockchainParams {
	return InsertBlockchainParams{
		BlockchainID:      blockchain.ID,
		Altruist:          newSQLNullString(blockchain.Altruist),
		Blockchain:        newSQLNullString(blockchain.Blockchain),
		ChainID:           newSQLNullString(blockchain.ChainID),
		ChainIDCheck:      newSQLNullString(blockchain.ChainIDCheck),
		Path:              newSQLNullString(blockchain.Path),
		Description:       newSQLNullString(blockchain.Description),
		EnforceResult:     newSQLNullString(blockchain.EnforceResult),
		Network:           newSQLNullString(blockchain.Network),
		Ticker:            newSQLNullString(blockchain.Ticker),
		BlockchainAliases: blockchain.BlockchainAliases,
		LogLimitBlocks:    newSQLNullInt32(int32(blockchain.LogLimitBlocks), false),
		RequestTimeout:    newSQLNullInt32(int32(blockchain.RequestTimeout), false),
		Active:            newSQLNullBool(&blockchain.Active),
		CreatedAt:         newSQLNullTime(blockchain.CreatedAt),
		UpdatedAt:         newSQLNullTime(blockchain.UpdatedAt),
	}
}

func extractInsertSyncCheckOptions(blockchain *types.Blockchain) InsertSyncCheckOptionsParams {
	return InsertSyncCheckOptionsParams{
		BlockchainID: blockchain.ID,
		Synccheck:    newSQLNullString(blockchain.SyncCheck),
		Body:         newSQLNullString(blockchain.SyncCheckOptions.Body),
		Path:         newSQLNullString(blockchain.SyncCheckOptions.Path),
		ResultKey:    newSQLNullString(blockchain.SyncCheckOptions.ResultKey),
		Allowance:    newSQLNullInt32(int32(blockchain.SyncCheckOptions.Allowance), false),
	}
}

func (i *InsertSyncCheckOptionsParams) isNotNull() bool {
	return i.Synccheck.Valid || i.Body.Valid || i.Path.Valid || i.ResultKey.Valid || i.Allowance.Valid
}

/* WriteRedirect saves input Redirect struct to the database.
It must be called separately from WriteBlockchain due to how new chains are added to the dB */
func (p *PostgresDriver) WriteRedirect(ctx context.Context, redirect *types.Redirect) (*types.Redirect, error) {
	time := time.Now()
	redirect.CreatedAt = time
	redirect.UpdatedAt = time

	err := p.InsertRedirect(ctx, extractInsertDBRedirect(redirect))
	if err != nil {
		return nil, err
	}

	return redirect, nil
}

func extractInsertDBRedirect(redirect *types.Redirect) InsertRedirectParams {
	return InsertRedirectParams{
		BlockchainID: redirect.BlockchainID,
		Alias:        redirect.Alias,
		Loadbalancer: redirect.LoadBalancerID,
		Domain:       redirect.Domain,
		CreatedAt:    newSQLNullTime(redirect.CreatedAt),
		UpdatedAt:    newSQLNullTime(redirect.UpdatedAt),
	}
}

/* ActivateChain toggles chain.active field on or off */
func (p *PostgresDriver) ActivateChain(ctx context.Context, id string, active bool) error {
	params := ActivateBlockchainParams{
		BlockchainID: id,
		Active:       newSQLNullBool(&active),
		UpdatedAt:    newSQLNullTime(time.Now()),
	}

	err := p.ActivateBlockchain(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

/* Used by Listener */
type (
	dbBlockchainJSON struct {
		BlockchainID      string   `json:"blockchain_id"`
		Altruist          string   `json:"altruist"`
		Blockchain        string   `json:"blockchain"`
		ChainID           string   `json:"chain_id"`
		ChainIDCheck      string   `json:"chain_id_check"`
		ChainPath         string   `json:"path"`
		Description       string   `json:"description"`
		EnforceResult     string   `json:"enforce_result"`
		Network           string   `json:"network"`
		Ticker            string   `json:"ticker"`
		BlockchainAliases []string `json:"blockchain_aliases"`
		LogLimitBlocks    int      `json:"log_limit_blocks"`
		RequestTimeout    int      `json:"request_timeout"`
		Active            bool     `json:"active"`
		CreatedAt         string   `json:"created_at"`
		UpdatedAt         string   `json:"updated_at"`
	}
	dbSyncCheckOptionsJSON struct {
		BlockchainID string `json:"blockchain_id"`
		SyncCheck    string `json:"synccheck"`
		Body         string `json:"body"`
		Path         string `json:"path"`
		ResultKey    string `json:"result_key"`
		Allowance    int    `json:"allowance"`
	}
	dbRedirectJSON struct {
		BlockchainID   string `json:"blockchain_id"`
		Alias          string `json:"alias"`
		LoadBalancerID string `json:"loadbalancer"`
		Domain         string `json:"domain"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at"`
	}
)

func (j dbBlockchainJSON) toOutput() *types.Blockchain {
	return &types.Blockchain{
		ID:                j.BlockchainID,
		Altruist:          j.Altruist,
		Blockchain:        j.Blockchain,
		ChainID:           j.ChainID,
		ChainIDCheck:      j.ChainIDCheck,
		Path:              j.ChainPath,
		Description:       j.Description,
		EnforceResult:     j.EnforceResult,
		Network:           j.Network,
		Ticker:            j.Ticker,
		BlockchainAliases: j.BlockchainAliases,
		LogLimitBlocks:    j.LogLimitBlocks,
		RequestTimeout:    j.RequestTimeout,
		Active:            j.Active,
		CreatedAt:         psqlDateToTime(j.CreatedAt),
		UpdatedAt:         psqlDateToTime(j.UpdatedAt),
	}
}
func (j dbSyncCheckOptionsJSON) toOutput() *types.SyncCheckOptions {
	return &types.SyncCheckOptions{
		BlockchainID: j.BlockchainID,
		Body:         j.Body,
		Path:         j.Path,
		ResultKey:    j.ResultKey,
		Allowance:    j.Allowance,
	}
}

func (j dbRedirectJSON) toOutput() *types.Redirect {
	return &types.Redirect{
		BlockchainID:   j.BlockchainID,
		Alias:          j.Alias,
		LoadBalancerID: j.LoadBalancerID,
		Domain:         j.Domain,
		CreatedAt:      psqlDateToTime(j.CreatedAt),
		UpdatedAt:      psqlDateToTime(j.UpdatedAt),
	}
}
