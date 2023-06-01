package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GetChainAliases(t *testing.T) {
	tests := []struct {
		name     string
		chain    Chain
		expected []ChainAlias
	}{
		{
			name: "Should return all aliases of the chain",
			chain: Chain{
				AliasDomains: map[ChainAlias][]ChainDomain{
					"pokt-mainnet": {"pokt-rpc.gateway.pokt.network"},
					"pokt-testnet": {"pokt-rpc-test.gateway.pokt.network"},
				},
			},
			expected: []ChainAlias{"pokt-mainnet", "pokt-testnet"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aliases := test.chain.GetChainAliases()
			assert.ElementsMatch(t, test.expected, aliases)
		})
	}
}

func Test_GetChainDomains(t *testing.T) {
	tests := []struct {
		name     string
		chain    Chain
		expected []ChainDomain
	}{
		{
			name: "Should return all domains of the chain",
			chain: Chain{
				AliasDomains: map[ChainAlias][]ChainDomain{
					"pokt-mainnet": {"pokt-rpc.gateway.pokt.network", "pokt-rpc-2.gateway.pokt.network"},
					"pokt-testnet": {"pokt-rpc-test.gateway.pokt.network"},
				},
			},
			expected: []ChainDomain{"pokt-rpc.gateway.pokt.network", "pokt-rpc-2.gateway.pokt.network", "pokt-rpc-test.gateway.pokt.network"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			domains := test.chain.GetChainDomains()
			assert.ElementsMatch(t, test.expected, domains)
		})
	}
}

func Test_GetGigastakeAATs(t *testing.T) {
	tests := []struct {
		name     string
		chain    Chain
		expected []AAT
	}{
		{
			name: "Should return all AATs of the chain",
			chain: Chain{
				GigastakeApps: map[ProtocolAppID]*GigastakeApp{
					"aat_1": {AAT: AAT{Version: "0.0.1", ClientPublicKey: "clientPublicKey1", PublicKey: "applicationPublicKey1", Signature: "applicationSignature1"}},
					"aat_2": {AAT: AAT{Version: "0.0.2", ClientPublicKey: "clientPublicKey2", PublicKey: "applicationPublicKey2", Signature: "applicationSignature2"}},
				},
			},
			expected: []AAT{
				{Version: "0.0.1", ClientPublicKey: "clientPublicKey1", PublicKey: "applicationPublicKey1", Signature: "applicationSignature1"},
				{Version: "0.0.2", ClientPublicKey: "clientPublicKey2", PublicKey: "applicationPublicKey2", Signature: "applicationSignature2"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aats := test.chain.GetGigastakeAATs()
			assert.ElementsMatch(t, test.expected, aats)
		})
	}
}

func Test_UpdateBlockchain(t *testing.T) {
	mockTime := time.Now()

	tests := []struct {
		name     string
		chain    Chain
		update   Chain
		expected Chain
	}{
		{
			name: "Should update all fields",
			chain: Chain{
				ID:             "0001",
				Blockchain:     "pokt-mainnet",
				Description:    "Pocket Network Mainnet",
				EnforceResult:  "JSON",
				Path:           "/v1/query/height",
				Ticker:         "POKT",
				Active:         true,
				AllowedMethods: []string{"GET", "POST"},
				LogLimitBlocks: 10,
				RequestTimeout: 5,
				Altruists: []Altruist{
					{
						URL:      "https://altruist-0001.com:1234",
						AuthType: ChainAuthTypeBasicAuth,
						Auth:     "test_pocket:auth123456",
					},
				},
				Checks: map[ChainCheckType]Check{
					ChainCheckTypeSync: {
						Type:      ChainCheckTypeSync,
						Payload:   `{"id":1,"jsonrpc":"2.0","method":"query"}`,
						ResultKey: "result.sync_info",
						Allowance: 1,
					},
				},
				AliasDomains: map[ChainAlias][]ChainDomain{
					"pokt-mainnet": {"pokt-rpc.gateway.pokt.network"},
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Deleted:   false,
			},
			update: Chain{
				Blockchain:     "updated-blockchain",
				Description:    "updated-description",
				EnforceResult:  "updated-result",
				Path:           "/updated/path",
				Ticker:         "updated-ticker",
				Active:         false,
				AllowedMethods: []string{"PUT"},
				LogLimitBlocks: 20,
				RequestTimeout: 10,
				Altruists: []Altruist{
					{
						URL:      "https://altruist-updated.com:5678",
						AuthType: ChainAuthTypeBearerToken,
						Auth:     "updated_auth",
					},
				},
				Checks: map[ChainCheckType]Check{
					ChainCheckTypeSync: {
						Type:      ChainCheckTypeSync,
						Payload:   `{"id":2,"jsonrpc":"2.0","method":"query"}`,
						ResultKey: "result.sync_info_updated",
						Allowance: 2,
					},
				},
				AliasDomains: map[ChainAlias][]ChainDomain{
					"updated-alias": {"updated.domain.com"},
				},
				UpdatedAt: mockTime.Add(10 * time.Minute),
				Deleted:   true,
			},
			expected: Chain{
				ID:             "0001",
				Blockchain:     "updated-blockchain",
				Description:    "updated-description",
				EnforceResult:  "updated-result",
				Path:           "/updated/path",
				Ticker:         "updated-ticker",
				Active:         false,
				AllowedMethods: []string{"PUT"},
				LogLimitBlocks: 20,
				RequestTimeout: 10,
				Altruists: []Altruist{
					{
						URL:      "https://altruist-updated.com:5678",
						AuthType: ChainAuthTypeBearerToken,
						Auth:     "updated_auth",
					},
				},
				Checks: map[ChainCheckType]Check{
					ChainCheckTypeSync: {
						Type:      ChainCheckTypeSync,
						Payload:   `{"id":2,"jsonrpc":"2.0","method":"query"}`,
						ResultKey: "result.sync_info_updated",
						Allowance: 2,
					},
				},
				AliasDomains: map[ChainAlias][]ChainDomain{
					"updated-alias": {"updated.domain.com"},
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime.Add(10 * time.Minute),
				Deleted:   true,
			},
		},
		{
			name: "Should update blockchain",
			chain: Chain{
				Blockchain: "Old Blockchain",
			},
			update: Chain{
				Blockchain: "New Blockchain",
			},
			expected: Chain{
				Blockchain: "New Blockchain",
			},
		},
		{
			name: "Should update description",
			chain: Chain{
				Description: "Old Description",
			},
			update: Chain{
				Description: "New Description",
			},
			expected: Chain{
				Description: "New Description",
			},
		},
		{
			name: "Should update enforceResult",
			chain: Chain{
				EnforceResult: "Old Result",
			},
			update: Chain{
				EnforceResult: "New Result",
			},
			expected: Chain{
				EnforceResult: "New Result",
			},
		},
		{
			name: "Should update path",
			chain: Chain{
				Path: "Old Path",
			},
			update: Chain{
				Path: "New Path",
			},
			expected: Chain{
				Path: "New Path",
			},
		},
		{
			name: "Should update ticker",
			chain: Chain{
				Ticker: "Old Ticker",
			},
			update: Chain{
				Ticker: "New Ticker",
			},
			expected: Chain{
				Ticker: "New Ticker",
			},
		},
		{
			name: "Should update active",
			chain: Chain{
				Active: true,
			},
			update: Chain{
				Active: false,
			},
			expected: Chain{
				Active: false,
			},
		},
		{
			name: "Should update logLimitBlocks",
			chain: Chain{
				LogLimitBlocks: 5,
			},
			update: Chain{
				LogLimitBlocks: 10,
			},
			expected: Chain{
				LogLimitBlocks: 10,
			},
		},
		{
			name: "Should update requestTimeout",
			chain: Chain{
				RequestTimeout: 10,
			},
			update: Chain{
				RequestTimeout: 20,
			},
			expected: Chain{
				RequestTimeout: 20,
			},
		},
		{
			name: "Should update altruists",
			chain: Chain{
				Altruists: []Altruist{
					{
						URL:      "https://old-altruist.com:1234",
						AuthType: ChainAuthTypeBasicAuth,
						Auth:     "old_auth",
					},
				},
			},
			update: Chain{
				Altruists: []Altruist{
					{
						URL:      "https://new-altruist.com:5678",
						AuthType: ChainAuthTypeBearerToken,
						Auth:     "new_auth",
					},
				},
			},
			expected: Chain{
				Altruists: []Altruist{
					{
						URL:      "https://new-altruist.com:5678",
						AuthType: ChainAuthTypeBearerToken,
						Auth:     "new_auth",
					},
				},
			},
		},
		{
			name: "Should update checks",
			chain: Chain{
				Checks: map[ChainCheckType]Check{
					ChainCheckTypeSync: {
						Payload:   "Old Payload",
						ResultKey: "Old Key",
						Allowance: 1,
					},
				},
			},
			update: Chain{
				Checks: map[ChainCheckType]Check{
					ChainCheckTypeSync: {
						Payload:   "New Payload",
						ResultKey: "New Key",
						Allowance: 2,
					},
				},
			},
			expected: Chain{
				Checks: map[ChainCheckType]Check{
					ChainCheckTypeSync: {
						Payload:   "New Payload",
						ResultKey: "New Key",
						Allowance: 2,
					},
				},
			},
		},
		{
			name: "Should update aliasDomains",
			chain: Chain{
				AliasDomains: map[ChainAlias][]ChainDomain{
					"old-alias": {"old.domain.com"},
				},
			},
			update: Chain{
				AliasDomains: map[ChainAlias][]ChainDomain{
					"new-alias": {"new.domain.com"},
				},
			},
			expected: Chain{
				AliasDomains: map[ChainAlias][]ChainDomain{
					"new-alias": {"new.domain.com"},
				},
			},
		},
		{
			name: "Should update deleted",
			chain: Chain{
				Deleted: false,
			},
			update: Chain{
				Deleted: true,
			},
			expected: Chain{
				Deleted: true,
			},
		},

		{
			name: "Should not update with empty update",
			chain: Chain{
				Blockchain:  "Old Blockchain",
				Description: "Old Description",
			},
			update: Chain{},
			expected: Chain{
				Blockchain:  "Old Blockchain",
				Description: "Old Description",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.chain.UpdateBlockchain(&test.update)
			assert.Equal(t, test.expected, test.chain)
		})
	}
}
