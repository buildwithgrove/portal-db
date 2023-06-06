package types

import (
	"testing"

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
