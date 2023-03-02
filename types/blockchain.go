package types

import (
	"time"
)

type ChainAuthtype string

const (
	None      ChainAuthtype = "none"
	BasicAuth ChainAuthtype = "basicAuth"
	Bearer    ChainAuthtype = "bearer"
)

type (
	Blockchain struct {
		ID                   string                    `json:"id"`
		Blockchain           string                    `json:"blockchain"`
		ChainID              string                    `json:"chainID"`
		ChainIDCheck         string                    `json:"chainIDCheck"`
		Description          string                    `json:"description"`
		EnforceResult        string                    `json:"enforceResult"`
		Path                 string                    `json:"path"`
		Ticker               string                    `json:"ticker"`
		BlockchainAliases    []string                  `json:"blockchainAliases"`
		LogLimitBlocks       int                       `json:"logLimitBlocks"`
		RequestTimeout       int                       `json:"requestTimeout"`
		Active               bool                      `json:"active"`
		Altruists            []ChainAltruists          `json:"altruists"`
		Redirects            []ChainGigastakesRedirect `json:"redirects"`
		SyncCheckOptions     ChainSyncCheckOptions     `json:"syncCheckOptions"`
		GlobalAllowedMethods ChainGlobalAllowedMethods `json:"globalAllowedMethods"`
		CreatedAt            time.Time                 `json:"createdAt"`
		UpdatedAt            time.Time                 `json:"updatedAt"`
	}
	ChainAltruists struct {
		BlockchainID string        `json:"blockchainID,omitempty"`
		URL          string        `json:"url"`
		Auth         string        `json:"auth"`
		AuthType     ChainAuthtype `json:"authType"`
		CreatedAt    time.Time     `json:"createdAt"`
		UpdatedAt    time.Time     `json:"updatedAt"`
	}
	ChainGigastakesRedirect struct {
		BlockchainID  string    `json:"blockchainID,omitempty"`
		Alias         string    `json:"alias"`
		Domain        string    `json:"domain"`
		ProtocolAppID string    `json:"loadBalancerID"`
		CreatedAt     time.Time `json:"createdAt"`
		UpdatedAt     time.Time `json:"updatedAt"`
	}
	ChainSyncCheckOptions struct {
		BlockchainID string    `json:"blockchainID,omitempty"`
		Body         string    `json:"body"`
		ResultKey    string    `json:"resultKey"`
		Allowance    int       `json:"allowance"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}
	ChainGlobalAllowedMethods struct {
		BlockchainID string    `json:"blockchainID,omitempty"`
		Method       string    `json:"method"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}

	// Represents global blocked addresses across the entire Portal
	// TODO should this be in a separate file?
	GlobalBlockedContracts struct {
		ID               string    `json:"id"`
		BlockedAddresses []string  `json:"blockedAddresses"`
		CreatedAt        time.Time `json:"createdAt"`
		UpdatedAt        time.Time `json:"updatedAt"`
	}

	/* Update structs */
	UpdateBlockchain struct {
		Altruist          string   `json:"altruist,omitempty"`
		Blockchain        string   `json:"blockchain,omitempty"`
		ChainIDCheck      string   `json:"chainIDCheck,omitempty"`
		Description       string   `json:"description,omitempty"`
		EnforceResult     string   `json:"enforceResult,omitempty"`
		Path              string   `json:"path,omitempty"`
		Ticker            string   `json:"ticker,omitempty"`
		BlockchainAliases []string `json:"blockchainAliases,omitempty"`
		LogLimitBlocks    int      `json:"logLimitBlocks,omitempty"`
		RequestTimeout    int      `json:"requestTimeout,omitempty"`

		Body      string `json:"body,omitempty"`
		ResultKey string `json:"resultKey,omitempty"`
		Allowance *int   `json:"allowance,omitempty"`

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
	if update.Path != "" {
		b.Path = update.Path
	}
	if update.Ticker != "" {
		b.Ticker = update.Ticker
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
	if update.ResultKey != "" {
		b.SyncCheckOptions.ResultKey = update.ResultKey
	}
	if update.Allowance != nil {
		b.SyncCheckOptions.Allowance = *update.Allowance
	}
}

func (u *UpdateBlockchain) syncCheckUpdateNotNil() bool {
	return u.Body != "" || u.ResultKey != "" || u.Allowance != nil
}
