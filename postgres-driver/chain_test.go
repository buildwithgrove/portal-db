package postgresdriver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pokt-foundation/portal-db/v2/testdata"
	"github.com/pokt-foundation/portal-db/v2/types"
)

func (ts *PGDriverTestSuite) Test_ReadChains() {
	tests := []struct {
		name    string
		chains  map[types.ChainID]*types.Chain
		options types.DriverOptions
		err     error
	}{
		{
			name: "Should return all non-deleted Chains from the database",
			chains: map[types.ChainID]*types.Chain{
				types.ChainID("0001"): testdata.Chains[types.ChainID("0001")],
				types.ChainID("0053"): testdata.Chains[types.ChainID("0053")],
				types.ChainID("0021"): testdata.Chains[types.ChainID("0021")],
				types.ChainID("0064"): testdata.Chains[types.ChainID("0064")],
				types.ChainID("0040"): testdata.Chains[types.ChainID("0040")],
			},
			options: types.DriverOptions{},
			err:     nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			chains, err := ts.driver.ReadChains(context.Background(), test.options)
			ts.Equal(test.err, err)
			ts.Equal(test.chains, chains)
		})
	}
}

func (ts *PGDriverTestSuite) Test_WriteChain() {
	tests := []struct {
		name            string
		chain           types.Chain
		testCreatedTime time.Time
		altruistURL     types.AltruistURL
		redirectDomain  types.RedirectDomain
		redirectAccount types.AccountID
		err             error
	}{
		{
			name:            "Should create a new Chain in the database",
			chain:           *testdata.TestCreateChain,
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name:            "Should fail if chain already exists in the database",
			chain:           *testdata.Chains[types.ChainID("0064")],
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errChainExists.Error(), "0064"),
		},
		{
			name:            "Should fail if any input Altruist has an invalid URL",
			altruistURL:     "htz:/bad-domain2",
			chain:           *testdata.TestCreateChain,
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errInvalidAltruistURL.Error(), "htz:/bad-domain2"),
		},
		{
			name:            "Should fail if any input Redirect has an invalid URL",
			altruistURL:     "http://im.good.com",
			redirectDomain:  "who did this",
			chain:           *testdata.TestCreateChain,
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errInvalidRedirectDomain.Error(), "who did this"),
		},
		{
			name:            "Should fail if any input Redirect has an invalid AccountID",
			altruistURL:     "http://im.good.com",
			redirectDomain:  "im.also.good.io",
			redirectAccount: 47,
			chain:           *testdata.TestCreateChain,
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errAccountDoesntExist.Error(), 47),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			if test.altruistURL != "" {
				test.chain.Altruists[0].URL = test.altruistURL
			}
			if test.redirectDomain != "" {
				test.chain.Redirects[0].Domain = test.redirectDomain
			}
			if test.redirectAccount != 0 {
				test.chain.Redirects[0].AccountID = test.redirectAccount
			}

			createdChain, err := ts.driver.WriteChain(context.Background(), test.chain, test.testCreatedTime)
			ts.Equal(test.err, err)

			if test.err == nil {
				ts.Equal(&test.chain, createdChain)

				chains, err := ts.driver.ReadChains(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				ts.Equal(&test.chain, chains[createdChain.ID])
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_UpdateChain() {
	tests := []struct {
		name               string
		chain, updateChain types.Chain
		testCreatedTime    time.Time
		altruistURL        types.AltruistURL
		redirectDomain     types.RedirectDomain
		redirectAccount    types.AccountID
		err                error
	}{
		{
			name:            "Should update an existing chain in the database",
			chain:           testdata.UpdateChainOne,
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name:            "Should update an existing chain again in the database",
			chain:           testdata.UpdateChainTwo,
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name:            "Should not remove any subtables if they are not present in the update struct",
			chain:           testdata.UpdateChainTwo,
			updateChain:     testdata.UpdateChainThree,
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name:            "Should fail if chain doesn't exist in the database",
			chain:           testdata.UpdateChainNotExists,
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errChainDoesntExist.Error(), "0073"),
		},
		{
			name:            "Should fail if any input Altruist has an invalid URL",
			altruistURL:     "htz:/bad-domain2",
			chain:           *testdata.Chains[types.ChainID("0040")],
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errInvalidAltruistURL.Error(), "htz:/bad-domain2"),
		},
		{
			name:            "Should fail if any input Redirect has an invalid URL",
			altruistURL:     "http://im.good.com",
			redirectDomain:  "who did this",
			chain:           *testdata.Chains[types.ChainID("0040")],
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errInvalidRedirectDomain.Error(), "who did this"),
		},
		{
			name:            "Should fail if any input Redirect has an invalid AccountID",
			altruistURL:     "http://im.good.com",
			redirectDomain:  "im.also.good.io",
			redirectAccount: 47,
			chain:           *testdata.Chains[types.ChainID("0040")],
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errAccountDoesntExist.Error(), 47),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			if test.altruistURL != "" {
				test.chain.Altruists[0].URL = test.altruistURL
			}
			if test.redirectDomain != "" {
				test.chain.Redirects[0].Domain = test.redirectDomain
			}
			if test.redirectAccount != 0 {
				test.chain.Redirects[0].AccountID = test.redirectAccount
			}

			var testUpdate types.Chain
			switch {
			case test.updateChain.ID != "":
				testUpdate = test.updateChain
			default:
				testUpdate = test.chain
			}

			err := ts.driver.UpdateChain(context.Background(), testUpdate, test.testCreatedTime)
			ts.Equal(test.err, err)

			if test.err == nil {
				chains, err := ts.driver.ReadChains(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				ts.Equal(&test.chain, chains[testUpdate.ID])
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_UpdateChainJSON() {
	tests := []struct {
		name                string
		chainID             types.ChainID
		chainJSON           string
		testCreatedTime     time.Time
		description, ticker string
		aliases             []string
		requestTimeout      int32
		altruists           []types.Altruist
		redirects           []types.GigastakeRedirect
		checks              map[types.ChainCheckType]types.Check
		err                 error
	}{
		{
			name:    "Should update only fields that are present in the JSON",
			chainID: "0001",
			chainJSON: `{
				"id": "0001",
				"description": "Pocket Network Meganet",
				"ticker": "POKT-MEGA",
				"blockchainAliases": [
				  "mainnet-mega",
				  "mainnet-ultra"
				],
				"requestTimeout": 5000,
				"altruists": [
				  {
					"chainID": "",
					"url": "https://altruist-0001.com:1234",
					"auth": "test_pocket:auth123456",
					"authType": "basic_auth"
				  },
				  {
					"chainID": "",
					"url": "https://altruist-0001-2.com:1234",
					"auth": "test_pocket:auth123456",
					"authType": "basic_auth"
				  }
				],
				"redirects": [],
				"createdAt": "2023-03-16T00:00:00Z",
				"updatedAt": "2023-03-16T00:00:00Z"
			  }`,
			testCreatedTime: testdata.MockTimestamp,
			description:     "Pocket Network Meganet",
			ticker:          "POKT-MEGA",
			aliases:         []string{"mainnet-mega", "mainnet-ultra"},
			requestTimeout:  5_000,
			altruists: []types.Altruist{
				{
					URL:      "https://altruist-0001.com:1234",
					AuthType: types.ChainAuthTypeBasicAuth,
					Auth:     "test_pocket:auth123456",
				},
				{
					URL:      "https://altruist-0001-2.com:1234",
					AuthType: types.ChainAuthTypeBasicAuth,
					Auth:     "test_pocket:auth123456",
				},
			},
			redirects: []types.GigastakeRedirect{},
			checks: map[types.ChainCheckType]types.Check{
				types.ChainCheckTypeSync: {
					Type:      types.ChainCheckTypeSync,
					Payload:   `{"id":1,"jsonrpc":"2.0","method":"query"}`,
					ResultKey: "result.sync_info",
					Allowance: 1,
				},
			},
			err: nil,
		},
		{
			name:    "Should update only fields that are present in the JSON again",
			chainID: "0001",
			chainJSON: `{
				"id": "0001",
				"description": "Pocket Network Ultranet",
				"ticker": "POKT-ULTRA",
				"blockchainAliases": [
				  "mainnet-ultra"
				],
				"requestTimeout": 10000,
				"altruists": [
				  {
					"chainID": "",
					"url": "https://altruist-0001.com:1234",
					"auth": "test_pocket:auth123456",
					"authType": "basic_auth"
				  },
				  {
					"chainID": "",
					"url": "https://altruist-0001-4.com:1234",
					"auth": "test_pocket:auth123456",
					"authType": "basic_auth"
				  }
				],
				"redirects": [
					{
					  "chainID": "",
					  "accountID": 1,
					  "domain": "pokt-rpc-ultra.gateway.pokt.network",
					  "alias": "altruist-0007"
					}
				],
				"chainChecks": {},
				"createdAt": "2023-03-16T00:00:00Z",
				"updatedAt": "2023-03-16T00:00:00Z"
			  }`,
			testCreatedTime: testdata.MockTimestamp,
			description:     "Pocket Network Ultranet",
			ticker:          "POKT-ULTRA",
			aliases:         []string{"mainnet-ultra"},
			requestTimeout:  10_000,
			altruists: []types.Altruist{
				{
					URL:      "https://altruist-0001.com:1234",
					AuthType: types.ChainAuthTypeBasicAuth,
					Auth:     "test_pocket:auth123456",
				},
				{
					URL:      "https://altruist-0001-4.com:1234",
					AuthType: types.ChainAuthTypeBasicAuth,
					Auth:     "test_pocket:auth123456",
				},
			},
			redirects: []types.GigastakeRedirect{
				{AccountID: 1, Alias: "altruist-0007", Domain: "pokt-rpc-ultra.gateway.pokt.network"},
			},
			checks: map[types.ChainCheckType]types.Check{},
			err:    nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			var chain types.Chain
			err := json.Unmarshal([]byte(test.chainJSON), &chain)
			ts.NoError(test.err, err)

			err = ts.driver.UpdateChain(context.Background(), chain, test.testCreatedTime)
			ts.Equal(test.err, err)

			if test.err == nil {
				chains, err := ts.driver.ReadChains(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				chain, ok := chains[test.chainID]
				ts.True(ok)
				ts.Equal(test.description, chain.Description)
				ts.Equal(test.ticker, chain.Ticker)
				ts.Equal(test.aliases, chain.ChainAliases)
				ts.Equal(test.requestTimeout, chain.RequestTimeout)
				ts.Equal(test.altruists, chain.Altruists)
				ts.Equal(test.redirects, chain.Redirects)
				ts.Equal(test.checks, chain.Checks)
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_SetChainActiveStatus() {
	tests := []struct {
		name    string
		chainID types.ChainID
		active  bool
		err     error
	}{
		{
			name:    "Should deactivate an existing chain in the database",
			chainID: "0001",
			active:  false,
			err:     nil,
		},
		{
			name:    "Should activate an existing chain in the database",
			chainID: "0001",
			active:  true,
			err:     nil,
		},
		{
			name:    "Should fail if chain does not exist in the database",
			chainID: "7701",
			err:     fmt.Errorf(errChainDoesntExist.Error(), "7701"),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			active, err := ts.driver.SetChainActiveStatus(context.Background(), test.chainID, test.active, testdata.MockTimestamp)
			ts.Equal(test.err, err)

			if test.err == nil {
				ts.Equal(test.active, active)

				chains, err := ts.driver.ReadChains(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				ts.Equal(test.active, chains[test.chainID].Active)
			}
		})
	}
}
