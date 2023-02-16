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
		BlockchainID   string    `json:"blockchainID"`
		Alias          string    `json:"alias"`
		Domain         string    `json:"domain"`
		LoadBalancerID string    `json:"loadBalancerID"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}
	SyncCheckOptions struct {
		BlockchainID string `json:"blockchainID"`
		Body         string `json:"body"`
		Path         string `json:"path"`
		ResultKey    string `json:"resultKey"`
		Allowance    int    `json:"allowance"`
	}
	/* Update structs */
	UpdateBlockchain struct {
		BlockchainID      string   `json:"blockchainID"`
		Altruist          string   `json:"altruist,omitempty"`
		Blockchain        string   `json:"blockchain,omitempty"`
		BlockchainAliases []string `json:"blockchainAliases,omitempty"`
		ChainIDCheck      string   `json:"chainIDCheck,omitempty"`
		Description       string   `json:"description,omitempty"`
		EnforceResult     string   `json:"enforceResult,omitempty"`
		LogLimitBlocks    int32    `json:"logLimitBlocks,omitempty"`
		Network           string   `json:"network,omitempty"`
		Path              string   `json:"path,omitempty"`
		RequestTimeout    int32    `json:"requestTimeout,omitempty"`
		Ticker            string   `json:"ticker,omitempty"`

		Synccheck     string `json:"synccheck,omitempty"`
		Allowance     *int32 `json:"allowance,omitempty"`
		Body          string `json:"body,omitempty"`
		SyncCheckPath string `json:"sync_check_path,omitempty"`
		ResultKey     string `json:"resultKey,omitempty"`

		UpdatedAt time.Time `json:"updatedAt"`
	}
)
