package postgresdriver

import (
	"context"
	"fmt"
	"time"

	"github.com/pokt-foundation/portal-db/v2/testdata"
	"github.com/pokt-foundation/portal-db/v2/types"
)

func (ts *PGDriverTestSuite) Test_ReadChains() {
	tests := []struct {
		name    string
		chains  map[types.RelayChainID]*types.Chain
		options types.DriverOptions
		err     error
	}{
		{
			name: "Should return all non-deleted Chains from the database",
			chains: map[types.RelayChainID]*types.Chain{
				types.RelayChainID("0001"): testdata.Chains[types.RelayChainID("0001")],
				types.RelayChainID("0053"): testdata.Chains[types.RelayChainID("0053")],
				types.RelayChainID("0021"): testdata.Chains[types.RelayChainID("0021")],
				types.RelayChainID("0064"): testdata.Chains[types.RelayChainID("0064")],
				types.RelayChainID("0040"): testdata.Chains[types.RelayChainID("0040")],
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

func (ts *PGDriverTestSuite) Test_WriteChainAndGigastakeApps() {
	tests := []struct {
		name            string
		newChainInput   types.NewChainInput
		testCreatedTime time.Time
		err             error
	}{
		{
			name:            "Should create a new Chain along with its GigastakeApps in the database",
			newChainInput:   testdata.TestCreateNewChainInput,
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name: "Should fail if chain is nil",
			newChainInput: types.NewChainInput{
				Chain:         nil,
				GigastakeApps: []*types.GigastakeApp{{}},
			},
			err: errChainCannotBeNil,
		},
		{
			name: "Should fail if GigastakeApps is empty",
			newChainInput: types.NewChainInput{
				Chain:         &types.Chain{},
				GigastakeApps: []*types.GigastakeApp{},
			},
			err: errGigastakeAppsCannotBeEmpty,
		},
		{
			name: "Should fail if altruist URL is invalid",
			newChainInput: types.NewChainInput{
				Chain: &types.Chain{
					ID:         "testID",
					Blockchain: "testBlockchain",
					Altruists:  map[types.AltruistURL]types.Altruist{"invalid_url": {URL: "invalid_url"}},
				},
				GigastakeApps: []*types.GigastakeApp{{}},
			},
			err: fmt.Errorf(errInvalidAltruistURL.Error(), "invalid_url"),
		},
		{
			name: "Should fail if domain is invalid",
			newChainInput: types.NewChainInput{
				Chain: &types.Chain{
					ID:         "testID",
					Blockchain: "testBlockchain",
					AliasDomains: map[types.ChainAlias][]types.ChainDomain{
						"testAlias": {"invalid_domain"},
					},
				},
				GigastakeApps: []*types.GigastakeApp{{}},
			},
			err: fmt.Errorf(errInvalidDomain.Error(), "invalid_domain", "testAlias"),
		},
		{
			name: "Should fail if trying to insert an existing chain",
			newChainInput: types.NewChainInput{
				Chain: &types.Chain{
					ID:         "0001",
					Blockchain: "testBlockchain",
				},
				GigastakeApps: []*types.GigastakeApp{{}},
			},
			err: fmt.Errorf(errChainExists.Error(), "0001"),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			createdChainData, err := ts.driver.WriteChainAndGigastakeApps(context.Background(), test.newChainInput, test.testCreatedTime)
			ts.Equal(test.err, err)

			if test.err == nil {
				ts.NotEmpty(createdChainData.Chain.ID)
				for i, gigastakeApp := range createdChainData.GigastakeApps {
					ts.NotEmpty(gigastakeApp.AATID)
					test.newChainInput.GigastakeApps[i].AATID = gigastakeApp.AATID
					break
				}

				chains, err := ts.driver.ReadChains(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				exists := false
				for chainID, chain := range chains {
					if chainID == test.newChainInput.Chain.ID {
						exists = true
						ts.Equal(test.newChainInput.Chain, chain)
						break
					}
				}
				ts.True(exists)

				gigastakeApps, err := ts.driver.ReadGigastakeApps(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				exists = false
				for aatID, gigastakeApp := range gigastakeApps {
					for _, testGigastakeApp := range test.newChainInput.GigastakeApps {
						testGigastakeApp.AAT.PrivateKey = "" // Private key is never read from the DB
						if aatID == testGigastakeApp.AATID {
							exists = true
							ts.Equal(testGigastakeApp, gigastakeApp)
							break
						}
					}
				}
				ts.True(exists)
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_WriteChain() {
	tests := []struct {
		name            string
		chain           types.Chain
		testCreatedTime time.Time
		altruistURL     types.AltruistURL
		aliasDomains    map[types.ChainAlias][]types.ChainDomain
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
			chain:           *testdata.Chains[types.RelayChainID("0064")],
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
			name:        "Should fail if any input alias has an invalid domain",
			altruistURL: "http://www.good-domain.com",
			aliasDomains: map[types.ChainAlias][]types.ChainDomain{
				"sol-mainnet": {"im-not-a-domain"},
			},
			chain:           *testdata.TestCreateChain,
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errInvalidDomain.Error(), "im-not-a-domain", "sol-mainnet"),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			testChain := test.chain

			if test.altruistURL != "" {
				altruist := testChain.Altruists[test.altruistURL]
				altruist.URL = test.altruistURL
				testChain.Altruists = make(map[types.AltruistURL]types.Altruist)
				testChain.Altruists[test.altruistURL] = altruist
			}
			if len(test.aliasDomains) > 0 {
				testChain.AliasDomains = test.aliasDomains
			}

			createdChain, err := ts.driver.WriteChain(context.Background(), testChain, test.testCreatedTime)
			ts.Equal(test.err, err)

			if test.err == nil {
				ts.Equal(&testChain, createdChain)

				chains, err := ts.driver.ReadChains(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				ts.Equal(&testChain, chains[createdChain.ID])
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
		aliasDomains       map[types.ChainAlias][]types.ChainDomain
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
			chain:           *testdata.Chains[types.RelayChainID("0040")],
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errInvalidAltruistURL.Error(), "htz:/bad-domain2"),
		},
		{
			name: "Should fail if any input alias has an invalid domain",
			aliasDomains: map[types.ChainAlias][]types.ChainDomain{
				"sol-mainnet": {"im-not-a-domain"},
			},
			chain:           *testdata.TestCreateChain,
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errInvalidDomain.Error(), "im-not-a-domain", "sol-mainnet"),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			testChain := test.chain

			if test.altruistURL != "" {
				altruist := testChain.Altruists[test.altruistURL]
				altruist.URL = test.altruistURL
				testChain.Altruists[test.altruistURL] = altruist
			}
			if len(test.aliasDomains) > 0 {
				testChain.AliasDomains = test.aliasDomains
			}

			var testUpdate types.Chain
			switch {
			case test.updateChain.ID != "":
				testUpdate = test.updateChain
			default:
				testUpdate = testChain
			}

			err := ts.driver.UpdateChain(context.Background(), testUpdate, test.testCreatedTime)
			ts.Equal(test.err, err)

			if test.err == nil {
				chains, err := ts.driver.ReadChains(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				ts.Equal(&testChain, chains[testUpdate.ID])
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_SetChainActiveStatus() {
	tests := []struct {
		name    string
		chainID types.RelayChainID
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
