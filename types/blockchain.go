package types

import (
	"time"
)

type (
	Blockchain struct {
		ID                string           `json:"id"`
		Altruist          string           `json:"altruist"`
		Blockchain        string           `json:"blockchain"`
		ChainID           string           `json:"chainID"`
		ChainIDCheck      string           `json:"chainIDCheck"`
		Description       string           `json:"description"`
		EnforceResult     string           `json:"enforceResult"`
		Network           string           `json:"network"`
		Path              string           `json:"path"`
		SyncCheck         string           `json:"syncCheck"`
		Ticker            string           `json:"ticker"`
		BlockchainAliases []string         `json:"blockchainAliases"`
		LogLimitBlocks    int              `json:"logLimitBlocks"`
		RequestTimeout    int              `json:"requestTimeout"`
		SyncAllowance     int              `json:"syncAllowance"`
		Active            bool             `json:"active"`
		Redirects         []Redirect       `json:"redirects"`
		SyncCheckOptions  SyncCheckOptions `json:"syncCheckOptions"`
		CreatedAt         time.Time        `json:"createdAt"`
		UpdatedAt         time.Time        `json:"updatedAt"`
	}
	Redirect struct {
		BlockchainID   string    `json:"blockchainID,omitempty"`
		Alias          string    `json:"alias"`
		Domain         string    `json:"domain"`
		LoadBalancerID string    `json:"loadBalancerID"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}
	SyncCheckOptions struct {
		BlockchainID string `json:"blockchainID,omitempty"`
		Body         string `json:"body"`
		Path         string `json:"path"`
		ResultKey    string `json:"resultKey"`
		Allowance    int    `json:"allowance"`
	}
	/* Update structs */
	UpdateBlockchain struct {
		Altruist          string   `json:"altruist,omitempty"`
		Blockchain        string   `json:"blockchain,omitempty"`
		ChainIDCheck      string   `json:"chainIDCheck,omitempty"`
		Description       string   `json:"description,omitempty"`
		EnforceResult     string   `json:"enforceResult,omitempty"`
		Network           string   `json:"network,omitempty"`
		Path              string   `json:"path,omitempty"`
		Ticker            string   `json:"ticker,omitempty"`
		BlockchainAliases []string `json:"blockchainAliases,omitempty"`
		LogLimitBlocks    int      `json:"logLimitBlocks,omitempty"`
		RequestTimeout    int      `json:"requestTimeout,omitempty"`

		Synccheck     string `json:"synccheck,omitempty"`
		Body          string `json:"body,omitempty"`
		SyncCheckPath string `json:"sync_check_path,omitempty"`
		ResultKey     string `json:"resultKey,omitempty"`
		Allowance     *int   `json:"allowance,omitempty"`

		UpdatedAt time.Time `json:"updatedAt"`
	}
)

func (b *Blockchain) UpdateBlockchain(update *UpdateBlockchain) *Blockchain {
	if update.Altruist != "" {
		b.Altruist = update.Altruist
	}
	if update.Blockchain != "" {
		b.Blockchain = update.Blockchain
	}
	if update.ChainIDCheck != "" {
		b.ChainIDCheck = update.ChainIDCheck
	}
	if update.Description != "" {
		b.Description = update.Description
	}
	if update.EnforceResult != "" {
		b.EnforceResult = update.EnforceResult
	}
	if update.Network != "" {
		b.Network = update.Network
	}
	if update.Path != "" {
		b.Path = update.Path
	}
	if update.Ticker != "" {
		b.Ticker = update.Ticker
	}
	if update.Synccheck != "" {
		b.SyncCheck = update.Synccheck
	}
	if update.BlockchainAliases != nil {
		b.BlockchainAliases = update.BlockchainAliases
	}
	if update.LogLimitBlocks != 0 {
		b.LogLimitBlocks = update.LogLimitBlocks
	}
	if update.RequestTimeout != 0 {
		b.RequestTimeout = update.RequestTimeout
	}

	if update.syncCheckUpdateNotNil() {
		b.updateSyncCheckOptions(update)
	}

	return b
}

func (b *Blockchain) updateSyncCheckOptions(update *UpdateBlockchain) {
	if update.Body != "" {
		b.SyncCheckOptions.Body = update.Body
	}
	if update.SyncCheckPath != "" {
		b.SyncCheckOptions.Path = update.SyncCheckPath
	}
	if update.ResultKey != "" {
		b.SyncCheckOptions.ResultKey = update.ResultKey
	}

	if update.Allowance != nil {
		b.SyncCheckOptions.Allowance = *update.Allowance
	}
}

func (u *UpdateBlockchain) syncCheckUpdateNotNil() bool {
	return u.Synccheck != "" || u.Body != "" || u.SyncCheckPath != "" || u.ResultKey != ""
}
