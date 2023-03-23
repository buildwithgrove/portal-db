package types

import (
	"net/url"
	"regexp"
	"time"
)

/* Enums */
type (
	ChainAuthType  string
	ChainCheckType string

	AltruistURL    string
	RedirectDomain string
)

const (
	ChainAuthTypeBasicAuth   ChainAuthType = "basic_auth"
	ChainAuthTypeBearerToken ChainAuthType = "bearer_token"
	ChainAuthTypeNone        ChainAuthType = "none"

	ChainCheckTypeArchival ChainCheckType = "archival"
	ChainCheckTypeChain    ChainCheckType = "chain"
	ChainCheckTypeMerge    ChainCheckType = "merge"
	ChainCheckTypeSync     ChainCheckType = "sync"
)

func (c ChainAuthType) IsValid() bool {
	switch c {
	case ChainAuthTypeBasicAuth, ChainAuthTypeBearerToken, ChainAuthTypeNone:
		return true
	default:
		return false
	}
}

func (c ChainCheckType) IsValid() bool {
	switch c {
	case ChainCheckTypeArchival, ChainCheckTypeChain, ChainCheckTypeMerge, ChainCheckTypeSync:
		return true
	default:
		return false
	}
}

func (a AltruistURL) IsValid() bool {
	parsedURL, err := url.Parse(string(a))
	if err != nil {
		return false
	}
	if parsedURL.Scheme == "http" || parsedURL.Scheme == "https" {
		return true
	}
	return false
}

func (r RedirectDomain) IsValid() bool {
	pattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, string(r))
	if err != nil {
		return false
	}
	return matched
}

/* Chain Struct and Methods */
type (
	Chain struct {
		ID             ChainID                  `json:"id"`
		Blockchain     string                   `json:"blockchain"`
		Description    string                   `json:"description"`
		EnforceResult  string                   `json:"enforceResult"`
		Path           string                   `json:"path"`
		Ticker         string                   `json:"ticker"`
		ChainAliases   []string                 `json:"blockchainAliases"`
		AllowedMethods []string                 `json:"allowedMethods"`
		BlockchainID   int32                    `json:"blockchainID"`
		LogLimitBlocks int32                    `json:"logLimitBlocks"`
		RequestTimeout int32                    `json:"requestTimeout"`
		Active         bool                     `json:"active"`
		Altruists      []Altruist               `json:"altruists,omitempty"`
		Redirects      []GigastakeRedirect      `json:"redirects,omitempty"`
		Checks         map[ChainCheckType]Check `json:"chainChecks,omitempty"`
		CreatedAt      time.Time                `json:"createdAt"`
		UpdatedAt      time.Time                `json:"updatedAt"`
	}
	Altruist struct {
		ChainID  ChainID       `json:"chainID,omitempty"`
		URL      AltruistURL   `json:"url"`
		Auth     string        `json:"auth"`
		AuthType ChainAuthType `json:"authType"`
	}
	GigastakeRedirect struct {
		ChainID   ChainID        `json:"chainID,omitempty"`
		AccountID AccountID      `json:"accountID"`
		Domain    RedirectDomain `json:"domain"`
		Alias     string         `json:"alias"`

		// TODO - remove when v2 migration finished
		// LegacyLoadBalancerID is the load balancer ID that the account was migrated from
		LegacyLoadBalancerID string `json:"legacyLoadBalancerID"`
	}
	Check struct {
		ChainID   ChainID        `json:"chainID,omitempty"`
		Type      ChainCheckType `json:"type"`
		Payload   string         `json:"payload"`
		ResultKey string         `json:"resultKey"`
		Allowance int32          `json:"allowance"`
	}

	/* Update structs */
	UpdateChain struct {
		Blockchain     string     `json:"blockchain,omitempty"`
		Description    string     `json:"description,omitempty"`
		EnforceResult  string     `json:"enforceResult,omitempty"`
		Path           string     `json:"path,omitempty"`
		Ticker         string     `json:"ticker,omitempty"`
		ChainAliases   []string   `json:"blockchainAliases,omitempty"`
		LogLimitBlocks int32      `json:"logLimitBlocks,omitempty"`
		RequestTimeout int32      `json:"requestTimeout,omitempty"`
		Altruists      []Altruist `json:"altruists,omitempty"`

		Checks map[ChainCheckType]UpdateCheck `json:"chainChecks"`

		UpdatedAt time.Time `json:"updatedAt"`
	}
	UpdateCheck struct {
		ChainID   string `json:"chainID,omitempty"`
		Payload   string `json:"payload"`
		ResultKey string `json:"resultKey"`
		Allowance *int32 `json:"allowance"` // must be able to set allowance to 0
	}
)

func (c *Chain) GetChainCheck(checkType ChainCheckType) Check {
	return c.Checks[checkType]
}

func (c *Chain) UpdateBlockchain(update *UpdateChain) {
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
	if update.ChainAliases != nil {
		c.ChainAliases = update.ChainAliases
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

func (c *Chain) Table() Table {
	return TableChains
}

func (c *Altruist) Table() Table {
	return TableChainAltruists
}

func (c *GigastakeRedirect) Table() Table {
	return TableChainGigastakeRedirects
}

func (c *Check) Table() Table {
	return TableChainChecks
}
