package postgresdriver

import (
	"context"
	"time"

	"github.com/pokt-foundation/portal-db/testdata"
	"github.com/pokt-foundation/portal-db/types"
)

func (ts *PGDriverTestSuite) Test_ReadPortalApps() {
	tests := []struct {
		name       string
		portalApps map[types.PortalAppID]*types.PortalApp
		err        error
	}{
		{
			name: "Should return all PortalApps from the database",

			portalApps: map[types.PortalAppID]*types.PortalApp{
				testdata.TestPortalAppOne.ID:   &testdata.TestPortalAppOne,
				testdata.TestPortalAppTwo.ID:   &testdata.TestPortalAppTwo,
				testdata.TestPortalAppThree.ID: &testdata.TestPortalAppThree,
			},
			err: nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			portalApps, err := ts.driver.ReadPortalApps(context.Background())
			ts.Equal(test.err, err)
			ts.Equal(test.portalApps, portalApps)
		})
	}
}

func (ts *PGDriverTestSuite) Test_WritePortalApp() {
	tests := []struct {
		name            string
		portalApp       types.PortalApp
		testCreatedTime time.Time
		err             error
	}{
		{
			name:            "Should create a new PortalApp in the database",
			portalApp:       testdata.TestCreatePortalAppOne,
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			createdPortalApp, err := ts.driver.WritePortalApp(context.Background(), test.portalApp, test.testCreatedTime)
			ts.Equal(test.err, err)

			test.portalApp.ID = createdPortalApp.ID
			ts.Equal(&test.portalApp, createdPortalApp)

			portalApps, err := ts.driver.ReadPortalApps(context.Background())
			ts.Equal(test.err, err)
			ts.Equal(&test.portalApp, portalApps[createdPortalApp.ID])
		})
	}
}

func (ts *PGDriverTestSuite) Test_UpdatePortalApp() {
	tests := []struct {
		name                  string
		updatePortalApp       types.UpdatePortalApp
		testUpdateTime        time.Time
		testUpdatedWhitelists types.Whitelists
		err                   error
	}{
		{
			name:            "Should update a new PortalApp in the database",
			updatePortalApp: testdata.TestUpdatePortalAppOne,
			testUpdateTime:  testdata.MockTimestamp,
			testUpdatedWhitelists: types.Whitelists{
				Origins:     map[types.Origin]struct{}{"https://portalgun.io": {}, "https://subdomain.example.com": {}, "https://www.example.com": {}},
				UserAgents:  map[types.UserAgent]struct{}{"Brave": {}, "Google Chrome": {}, "Mozilla Firefox": {}, "Netscape Navigator": {}, "Safari": {}},
				Blockchains: map[types.BlockchainID]struct{}{"0001": {}, "0002": {}, "003E": {}, "0056": {}},
				Contracts: map[types.BlockchainID]map[types.Contract]struct{}{
					"0001": {"0xtest_2f78db6436527729929aaf6c616361de0f7": {}, "0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be": {}},
					"0002": {"0xtest_1111117dc0aa78b770fa6a738034120c302": {}, "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2": {}},
					"003E": {"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984": {}, "0xtest_f958d2ee523a2206206994597c13d831ec7": {}},
					"0056": {"0xtest_00000f279d81a1d3cc75430faa017fa5a2e": {}, "0xtest_5068778dd592e39a122f4f5a5cf09c90fe2": {}},
				},
				Methods: map[types.BlockchainID]map[types.Method]struct{}{
					"0001": {"GET": {}, "POST": {}, "PUT": {}},
					"0002": {"DELETE": {}, "GET": {}, "POST": {}, "PUT": {}},
					"003E": {"GET": {}},
					"0056": {"GET": {}, "POST": {}},
				},
			},
			err: nil,
		},
	}

	// TODO -> add rest of update tests

	for _, test := range tests {
		ts.Run(test.name, func() {
			err := ts.driver.UpdatePortalApp(context.Background(), test.updatePortalApp, test.testUpdateTime)
			ts.Equal(test.err, err)

			portalApps, err := ts.driver.ReadPortalApps(context.Background())
			ts.Equal(test.err, err)
			updatedPortalApp, ok := portalApps[test.updatePortalApp.AppID]
			ts.True(ok)
			ts.Equal(test.testUpdatedWhitelists, updatedPortalApp.Whitelists)
		})
	}
}
