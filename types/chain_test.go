package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ChainDomain_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		domain   ChainDomain
		expected bool
	}{
		{
			name:     "Should return true for a valid domain without port",
			domain:   ChainDomain("pokt-rpc.gateway.pokt.network"),
			expected: true,
		},
		{
			name:     "Should return true for a valid domain with port",
			domain:   ChainDomain("pokt-rpc.gateway.pokt.network:80"),
			expected: true,
		},
		{
			name:     "Should return false for a domain without TLD",
			domain:   ChainDomain("pokt-rpc.gateway"),
			expected: false,
		},
		{
			name:     "Should return false for an invalid domain",
			domain:   ChainDomain("pokt-rpc..gateway.pokt.network"),
			expected: false,
		},
		{
			name:     "Should return false for an empty domain",
			domain:   ChainDomain(""),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := test.domain.IsValid()
			assert.Equal(t, test.expected, isValid)
		})
	}
}

func Test_ChainDomain_IsPublicEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		domain   ChainDomain
		expected bool
	}{
		{
			name:     "Should return true for public endpoint",
			domain:   ChainDomain("pokt-rpc.gateway.pokt.network"),
			expected: true,
		},
		{
			name:     "Should return false for non-public endpoint",
			domain:   ChainDomain("pokt.gateway.pokt.network"),
			expected: false,
		},
		{
			name:     "Should return false for empty domain",
			domain:   ChainDomain(""),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isPublic := test.domain.IsPublicEndpoint()
			assert.Equal(t, test.expected, isPublic)
		})
	}
}

func Test_ChainDomain_IsWildcardDomain(t *testing.T) {
	tests := []struct {
		name     string
		domain   ChainDomain
		expected bool
	}{
		{
			name:     "Should return true for wildcard domain",
			domain:   ChainDomain("*.gateway.pokt.network"),
			expected: true,
		},
		{
			name:     "Should return false for non-wildcard domain",
			domain:   ChainDomain("pokt-rpc.gateway.pokt.network"),
			expected: false,
		},
		{
			name:     "Should return false for empty domain",
			domain:   ChainDomain(""),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isWildcard := test.domain.IsWildcardDomain()
			assert.Equal(t, test.expected, isWildcard)
		})
	}
}

func Test_ChainDomain_GetAlias(t *testing.T) {
	tests := []struct {
		name     string
		domain   ChainDomain
		expected ChainAlias
	}{
		{
			name:     "Should return the alias for valid domain",
			domain:   ChainDomain("pokt-mainnet.gateway.pokt.network"),
			expected: ChainAlias("pokt-mainnet"),
		},
		{
			name:     "Should return empty alias for empty domain",
			domain:   ChainDomain(""),
			expected: ChainAlias(""),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			alias := test.domain.GetAlias()
			assert.Equal(t, test.expected, alias)
		})
	}
}

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

func Test_GetGigastakeApps(t *testing.T) {
	tests := []struct {
		name     string
		chain    Chain
		expected []GigastakeApp
	}{
		{
			name: "Should return all gigastake apps of the chain",
			chain: Chain{
				GigastakeApps: map[GigastakeAppID]*GigastakeApp{
					"aat_1": {Version: "0.0.1", ClientPublicKey: "clientPublicKey1", PublicKey: "applicationPublicKey1", Signature: "applicationSignature1"},
					"aat_2": {Version: "0.0.2", ClientPublicKey: "clientPublicKey2", PublicKey: "applicationPublicKey2", Signature: "applicationSignature2"},
				},
			},
			expected: []GigastakeApp{
				{Version: "0.0.1", ClientPublicKey: "clientPublicKey1", PublicKey: "applicationPublicKey1", Signature: "applicationSignature1"},
				{Version: "0.0.2", ClientPublicKey: "clientPublicKey2", PublicKey: "applicationPublicKey2", Signature: "applicationSignature2"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aats := test.chain.GetGigastakeApps()
			assert.ElementsMatch(t, test.expected, aats)
		})
	}
}
