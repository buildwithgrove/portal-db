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

	AltruistURL string
	ChainAlias  string
	ChainDomain string
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

func (r ChainDomain) IsValid() bool {
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
		ID             RelayChainID                 `json:"id"`
		Blockchain     string                       `json:"blockchain"`
		Description    string                       `json:"description"`
		EnforceResult  string                       `json:"enforceResult"`
		Path           string                       `json:"path"`
		Ticker         string                       `json:"ticker"`
		AllowedMethods []string                     `json:"allowedMethods"`
		LogLimitBlocks int32                        `json:"logLimitBlocks"`
		RequestTimeout int32                        `json:"requestTimeout"`
		Active         bool                         `json:"active"`
		Altruists      []Altruist                   `json:"altruists,omitempty"`
		Checks         map[ChainCheckType]Check     `json:"chainChecks,omitempty"`
		AliasDomains   map[ChainAlias][]ChainDomain `json:"domains"`
		CreatedAt      time.Time                    `json:"createdAt"`
		UpdatedAt      time.Time                    `json:"updatedAt"`
		Deleted        bool                         `json:"deleted"`

		// GigastakeApps are set inside PHD
		GigastakeApps []*GigastakeApp `json:"gigastakeApps,omitempty"`
	}
	Altruist struct {
		ChainID  RelayChainID  `json:"chainID,omitempty"`
		URL      AltruistURL   `json:"url"`
		Auth     string        `json:"auth"`
		AuthType ChainAuthType `json:"authType"`
	}
	Check struct {
		ChainID    RelayChainID   `json:"chainID,omitempty"`
		Type       ChainCheckType `json:"type"`
		Payload    string         `json:"payload"`
		ResultKey  string         `json:"resultKey"`
		Allowance  int32          `json:"allowance"`
		EVMChainID int32          `json:"evmChainID"`
	}
	// Used for mapping listener notification
	AliasDomains struct {
		ChainID RelayChainID  `json:"chainID,omitempty"`
		Alias   ChainAlias    `json:"alias"`
		Domains []ChainDomain `json:"domains"`
	}
)

func (c *Chain) GetChainCheck(checkType ChainCheckType) Check {
	return c.Checks[checkType]
}

func (c *Chain) UpdateBlockchain(update *Chain) {
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

func (c *Chain) updateChainChecks(update *Chain) {
	for checkType, check := range update.Checks {
		if check.Payload != "" {
			updatedCheck := Check{
				Payload:   check.Payload,
				ResultKey: check.ResultKey,
				Allowance: c.Checks[checkType].Allowance,
			}
			if check.Allowance != 0 {
				updatedCheck.Allowance = check.Allowance
			}
			if check.EVMChainID != 0 {
				updatedCheck.EVMChainID = check.EVMChainID
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

func (c *Check) Table() Table {
	return TableChainChecks
}

func (c *AliasDomains) Table() Table {
	return TableChainAliasDomains
}
