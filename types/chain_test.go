package types

import "time"

var (
	testChain = Chain{
		ID:                "0001",
		Blockchain:        "pokt-mainnet",
		Description:       "POKT Network Mainnet",
		EnforceResult:     "JSON",
		Path:              "/wow/test",
		Ticker:            "POKT",
		BlockchainAliases: []string{"pokt-mainnet"},
		AllowedMethods:    []string{"GET", "POST", "PUT"},
		LogLimitBlocks:    100_000,
		Active:            true,
		Altruists: []Altruist{
			{
				URL:      "https://user:test_123@pokt-test.us-1.pokt.network:1234",
				Auth:     "auth_123",
				AuthType: ChainAuthTypeBearerToken,
			},
		},
		Redirects: []GigastakeRedirect{testRedirect},
		Checks: map[ChainCheckType]Check{
			ChainCheckTypeSync: {
				Payload:   `{"method":"eth_blockNumber","id":1,"jsonrpc":"2.0"}`,
				ResultKey: "testing",
				Allowance: 1,
			},
			ChainCheckTypeChain: {Payload: `{"method":"eth_chainId","id":1,"jsonrpc":"2.0"}`},
		},

		CreatedAt: time.Date(2023, time.February, 14, 11, 11, 11, 0, time.UTC),
		UpdatedAt: time.Date(2023, time.February, 27, 13, 13, 13, 0, time.UTC),
	}

	testRedirect = GigastakeRedirect{
		Alias:         "mainnet",
		Domain:        "pokt.test.com",
		ProtocolAppID: "test_5416bb8d696386455b8",
	}
)
