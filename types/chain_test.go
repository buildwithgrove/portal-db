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
		LogLimitBlocks:    100_000,
		Active:            true,
		Altruists: []ChainAltruist{
			{
				URL:      "https://user:test_123@pokt-test.us-1.pokt.network:1234",
				Auth:     "auth_123",
				AuthType: ChainAuthBearer,
			},
		},
		Redirects: []ChainGigastakesRedirect{testRedirect},
		SyncCheckOptions: ChainSyncCheckOptions{
			Body:      `{}`,
			ResultKey: "testing",
			Allowance: 1,
		},
		GlobalAllowedMethods: ChainGlobalAllowedMethods{
			Methods: []string{"GET", "POST", "PUT"},
		},
		CreatedAt: time.Date(2023, time.February, 14, 11, 11, 11, 0, time.UTC),
		UpdatedAt: time.Date(2023, time.February, 27, 13, 13, 13, 0, time.UTC),
	}

	testRedirect = ChainGigastakesRedirect{
		Alias:         "mainnet",
		Domain:        "pokt.test.com",
		ProtocolAppID: "test_5416bb8d696386455b8",
	}
)
