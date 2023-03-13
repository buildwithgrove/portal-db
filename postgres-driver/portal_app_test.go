package postgresdriver

import (
	"context"
	"fmt"
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
		name                     string
		updatePortalApp          types.UpdatePortalApp
		testUpdateTime           time.Time
		testUpdatedName          string
		testUpdatedSettings      types.Settings
		testUpdatedNotifications map[NotificationType]types.AppNotification
		testUpdatedWhitelists    types.Whitelists
		err                      error
	}{
		{
			name:            "Should update a new PortalApp in the database with all fields",
			updatePortalApp: testdata.TestUpdatePortalAppAll,
			testUpdateTime:  testdata.MockTimestamp,
			testUpdatedName: "portal-app-updated",
			testUpdatedSettings: types.Settings{
				Environment:       types.EnvironmentProduction,
				SecretKey:         "test_9d07c8a96ad53e7c288b0e86f37c5680",
				SecretKeyRequired: true,
				MonthlyRelayLimit: 2_500_000,
				FavoritedChainIDs: map[types.ChainID]struct{}{"0003": {}, "0009": {}, "00H3": {}},
			},
			testUpdatedWhitelists: types.Whitelists{
				Origins:     map[types.Origin]struct{}{"https://portalgun.io": {}, "https://subdomain.example.com": {}, "https://www.example.com": {}},
				UserAgents:  map[types.UserAgent]struct{}{"Brave": {}, "Google Chrome": {}, "Mozilla Firefox": {}, "Netscape Navigator": {}, "Safari": {}},
				Blockchains: map[types.ChainID]struct{}{"0001": {}, "0002": {}, "003E": {}, "0056": {}},
				Contracts: map[types.ChainID]map[types.Contract]struct{}{
					"0001": {"0xtest_2f78db6436527729929aaf6c616361de0f7": {}, "0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be": {}},
					"0002": {"0xtest_1111117dc0aa78b770fa6a738034120c302": {}, "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2": {}},
					"003E": {"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984": {}, "0xtest_f958d2ee523a2206206994597c13d831ec7": {}},
					"0056": {"0xtest_00000f279d81a1d3cc75430faa017fa5a2e": {}, "0xtest_5068778dd592e39a122f4f5a5cf09c90fe2": {}},
				},
				Methods: map[types.ChainID]map[types.Method]struct{}{
					"0001": {"GET": {}, "POST": {}, "PUT": {}},
					"0002": {"DELETE": {}, "GET": {}, "POST": {}, "PUT": {}},
					"003E": {"GET": {}},
					"0056": {"GET": {}, "POST": {}},
				},
			},
			err: nil,
		},
		{
			name:            "Should update a new PortalApp in the database with only a new Name",
			updatePortalApp: testdata.TestUpdatePortalAppName,
			testUpdateTime:  testdata.MockTimestamp,
			testUpdatedName: "portal-app-updated",
			err:             nil,
		},
		{
			name:            "Should update a new PortalApp in the database with only new Settings",
			updatePortalApp: testdata.TestUpdatePortalAppSettings,
			testUpdateTime:  testdata.MockTimestamp,
			testUpdatedSettings: types.Settings{
				Environment:       types.EnvironmentProduction,
				SecretKey:         "test_9d07c8a96ad53e7c288b0e86f37c5680",
				SecretKeyRequired: true,
				MonthlyRelayLimit: 2_500_000,
				FavoritedChainIDs: map[types.ChainID]struct{}{"0003": {}, "0009": {}, "00H3": {}},
			},
			err: nil,
		},
		// TODO -> test notifications
		// {
		// 	name:                     "Should update a new PortalApp in the database with only new Notifications",
		// 	updatePortalApp:          testdata.TestUpdatePortalAppSettings,
		// 	testUpdateTime:           testdata.MockTimestamp,
		// 	testUpdatedNotifications: map[NotificationType]types.AppNotification{},
		// 	err:                      nil,
		// },
		{
			name:            "Should update a new PortalApp in the database with only new Whitelists",
			updatePortalApp: testdata.TestUpdatePortalAppWhitelists,
			testUpdateTime:  testdata.MockTimestamp,
			testUpdatedWhitelists: types.Whitelists{
				Origins:     map[types.Origin]struct{}{"https://portalgun.io": {}, "https://subdomain.example.com": {}, "https://www.example.com": {}},
				UserAgents:  map[types.UserAgent]struct{}{"Brave": {}, "Google Chrome": {}, "Mozilla Firefox": {}, "Netscape Navigator": {}, "Safari": {}},
				Blockchains: map[types.ChainID]struct{}{"0001": {}, "0002": {}, "003E": {}, "0056": {}},
				Contracts: map[types.ChainID]map[types.Contract]struct{}{
					"0001": {"0xtest_2f78db6436527729929aaf6c616361de0f7": {}, "0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be": {}},
					"0002": {"0xtest_1111117dc0aa78b770fa6a738034120c302": {}, "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2": {}},
					"003E": {"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984": {}, "0xtest_f958d2ee523a2206206994597c13d831ec7": {}},
					"0056": {"0xtest_00000f279d81a1d3cc75430faa017fa5a2e": {}, "0xtest_5068778dd592e39a122f4f5a5cf09c90fe2": {}},
				},
				Methods: map[types.ChainID]map[types.Method]struct{}{
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

	for i, test := range tests {
		ts.Run(test.name, func() {
			// Create new portal app for test case
			createApp := testdata.TestCreateAppForUpdate
			createApp.Name = fmt.Sprintf("test-update-portal-app-%d", i+1)
			createdPortalApp, err := ts.driver.WritePortalApp(context.Background(), createApp, testdata.MockTimestamp)
			ts.NoError(err)

			// Update created portal app for test case
			updateApp := test.updatePortalApp
			updateApp.AppID = createdPortalApp.ID
			err = ts.driver.UpdatePortalApp(context.Background(), updateApp, test.testUpdateTime)
			ts.Equal(test.err, err)

			// Get all portal apps from DB
			portalApps, err := ts.driver.ReadPortalApps(context.Background())
			ts.Equal(test.err, err)
			updatedPortalApp, ok := portalApps[createdPortalApp.ID]
			ts.True(ok)

			// Check update changes
			if test.testUpdatedName != "" {
				ts.Equal(test.testUpdatedName, updatedPortalApp.Name)
			} else {
				ts.Equal(createdPortalApp.Name, updatedPortalApp.Name)
			}

			// TODO -> test notifications

			if test.testUpdatedSettings.Environment != "" {
				ts.Equal(test.testUpdatedWhitelists, updatedPortalApp.Whitelists)
			} else {
				ts.Equal(createdPortalApp.Settings, updatedPortalApp.Settings)
			}

			if len(test.testUpdatedWhitelists.Origins) != 0 {
				ts.Equal(test.testUpdatedWhitelists, updatedPortalApp.Whitelists)
			} else {
				ts.Equal(createdPortalApp.Whitelists, updatedPortalApp.Whitelists)
			}

		})
	}
}
