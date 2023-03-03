package types

import (
	"time"
)

/* Enums */
type (
	ChainAuthType  string
	ChainCheckType string
)

const (
	ChainAuthNone      ChainAuthType = "none"
	ChainAuthBasicAuth ChainAuthType = "basicAuth"
	ChainAuthBearer    ChainAuthType = "bearer"

	CheckSync     ChainCheckType = "sync"
	CheckChain    ChainCheckType = "chain"
	CheckArchival ChainCheckType = "archival"
	CheckMerge    ChainCheckType = "merge"
)

/* Chain Struct and Methods */
type (
	Chain struct {
		ID                BlockchainID             `json:"id"`
		Blockchain        string                   `json:"blockchain"`
		ChainID           string                   `json:"chainID"`
		Description       string                   `json:"description"`
		EnforceResult     string                   `json:"enforceResult"`
		Path              string                   `json:"path"`
		Ticker            string                   `json:"ticker"`
		BlockchainAliases []string                 `json:"blockchainAliases"`
		AllowedMethods    []string                 `json:"allowedMethods"`
		LogLimitBlocks    int                      `json:"logLimitBlocks"`
		RequestTimeout    int                      `json:"requestTimeout"`
		Active            bool                     `json:"active"`
		Altruists         []Altruist               `json:"altruists"`
		Redirects         []GigastakeRedirect      `json:"redirects"`
		Checks            map[ChainCheckType]Check `json:"chainChecks"`
		CreatedAt         time.Time                `json:"createdAt"`
		UpdatedAt         time.Time                `json:"updatedAt"`
	}
	Altruist struct {
		BlockchainID string        `json:"blockchainID,omitempty"`
		URL          string        `json:"url"`
		Auth         string        `json:"auth"`
		AuthType     ChainAuthType `json:"authType"`
	}
	GigastakeRedirect struct {
		BlockchainID  string `json:"blockchainID,omitempty"`
		Alias         string `json:"alias"`
		Domain        string `json:"domain"`
		ProtocolAppID string `json:"loadBalancerID"`
	}
	Check struct {
		BlockchainID string `json:"blockchainID,omitempty"`
		Payload      string `json:"payload"`
		ResultKey    string `json:"resultKey"`
		Allowance    int    `json:"allowance"`
	}

	// Represents global blocked addresses across the entire Portal
	// TODO should this be in a separate file?
	GlobalBlockedContracts struct {
		ID               string   `json:"id"`
		BlockedAddresses []string `json:"blockedAddresses"`
	}

	/* Update structs */
	UpdateChain struct {
		Blockchain        string     `json:"blockchain,omitempty"`
		Description       string     `json:"description,omitempty"`
		EnforceResult     string     `json:"enforceResult,omitempty"`
		Path              string     `json:"path,omitempty"`
		Ticker            string     `json:"ticker,omitempty"`
		BlockchainAliases []string   `json:"blockchainAliases,omitempty"`
		LogLimitBlocks    int        `json:"logLimitBlocks,omitempty"`
		RequestTimeout    int        `json:"requestTimeout,omitempty"`
		Altruists         []Altruist `json:"altruists,omitempty"`

		Checks map[ChainCheckType]UpdateCheck `json:"chainChecks"`

		UpdatedAt time.Time `json:"updatedAt"`
	}
	UpdateCheck struct {
		BlockchainID string `json:"blockchainID,omitempty"`
		Payload      string `json:"payload"`
		ResultKey    string `json:"resultKey"`
		Allowance    *int   `json:"allowance"` // must be able to set allowance to 0
	}
)

func (c *Chain) GetChainCheck(checkType ChainCheckType) Check {
	return c.Checks[checkType]
}

func (c *Chain) UpdateBlockchain(update *UpdateChain) *Chain {
	if update.Blockchain != "" {
		c.Blockchain = update.Blockchain
	}
	if update.Description != "" {
		c.Description = update.Description
	}
	if update.EnforceResult != "" {
		c.EnforceResult = update.EnforceResult
	}
	if update.Path != "" {
		c.Path = update.Path
	}
	if update.Ticker != "" {
		c.Ticker = update.Ticker
	}
	if update.BlockchainAliases != nil {
		c.BlockchainAliases = update.BlockchainAliases
	}
	if update.LogLimitBlocks != 0 {
		c.LogLimitBlocks = update.LogLimitBlocks
	}
	if update.RequestTimeout != 0 {
		c.RequestTimeout = update.RequestTimeout
	}
	if update.Altruists != nil && len(update.Altruists) > 0 {
		c.Altruists = update.Altruists
	}
	if len(update.Checks) > 0 {
		c.updateChainChecks(update)
	}
	return c
}

func (c *Chain) updateChainChecks(update *UpdateChain) {
	for checkType, check := range update.Checks {
		if check.Payload != "" {
			updatedCheck := Check{
				Payload:   check.Payload,
				ResultKey: check.ResultKey,
				Allowance: c.Checks[checkType].Allowance,
			}
			if check.Allowance != nil {
				updatedCheck.Allowance = *check.Allowance
			}
			c.Checks[checkType] = updatedCheck
		}
	}
}
