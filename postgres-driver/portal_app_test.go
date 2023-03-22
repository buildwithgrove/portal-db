package postgresdriver

import (
	"context"
	"fmt"
	"time"

	"github.com/pokt-foundation/portal-db/v2/testdata"
	"github.com/pokt-foundation/portal-db/v2/types"
)

func (ts *PGDriverTestSuite) Test_ReadPortalApps() {
	tests := []struct {
		name       string
		portalApps map[types.PortalAppID]*types.PortalApp
		options    types.DriverOptions
		err        error
	}{
		{
			name: "Should return all non-deleted PortalApps from the database",
			portalApps: map[types.PortalAppID]*types.PortalApp{
				testdata.PortalApps["test_app_3487u329rfn23f9"].ID: testdata.PortalApps["test_app_3487u329rfn23f9"],
				testdata.PortalApps["test_app_2308rj09r23r9r2"].ID: testdata.PortalApps["test_app_2308rj09r23r9r2"],
				testdata.PortalApps["test_app_47fhs7j4hs7fj24"].ID: testdata.PortalApps["test_app_47fhs7j4hs7fj24"],
			},
			options: types.DriverOptions{},
			err:     nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			portalApps, err := ts.driver.ReadPortalApps(context.Background(), test.options)
			ts.Equal(test.err, err)
			ts.Equal(test.portalApps, portalApps)
		})
	}
}

func (ts *PGDriverTestSuite) Test_SetPortalAppDeleted() {
	tests := []struct {
		name                                          string
		deleteParams                                  DeletePortalAppParams
		portalAppsBeforeDelete, portalAppsAfterDelete map[types.PortalAppID]*types.PortalApp
		err                                           error
	}{
		{
			name: "Should set a PortalApp's deleted field to true, causing it to not appear in the ReadPortalApps query",
			deleteParams: DeletePortalAppParams{
				ID: testdata.PortalApps["test_app_47fhs7j4hs7fj24"].ID, DeletedAt: newSQLNullTime(testdata.MockTimestamp),
			},
			portalAppsBeforeDelete: map[types.PortalAppID]*types.PortalApp{
				testdata.PortalApps["test_app_3487u329rfn23f9"].ID: testdata.PortalApps["test_app_3487u329rfn23f9"],
				testdata.PortalApps["test_app_2308rj09r23r9r2"].ID: testdata.PortalApps["test_app_2308rj09r23r9r2"],
				testdata.PortalApps["test_app_47fhs7j4hs7fj24"].ID: testdata.PortalApps["test_app_47fhs7j4hs7fj24"],
			},
			portalAppsAfterDelete: map[types.PortalAppID]*types.PortalApp{
				testdata.PortalApps["test_app_3487u329rfn23f9"].ID: testdata.PortalApps["test_app_3487u329rfn23f9"],
				testdata.PortalApps["test_app_2308rj09r23r9r2"].ID: testdata.PortalApps["test_app_2308rj09r23r9r2"],
			},
			err: nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			// Check all PortalApps exist before delete
			portalApps, err := ts.driver.ReadPortalApps(context.Background(), types.DriverOptions{IncludeDeleted: false})
			ts.Equal(test.err, err)
			ts.Equal(test.portalAppsBeforeDelete, portalApps)

			// Delete PortalApp
			err = ts.driver.SetPortalAppDeleted(context.Background(), test.deleteParams.ID, test.deleteParams.DeletedAt.Time)
			ts.Equal(test.err, err)

			// Check PortalApp was deleted
			portalApps, err = ts.driver.ReadPortalApps(context.Background(), types.DriverOptions{IncludeDeleted: false})
			ts.Equal(test.err, err)
			ts.Equal(test.portalAppsAfterDelete, portalApps)

			// Check PortalApp still appears if IncludeDeleted: true
			portalApps, err = ts.driver.ReadPortalApps(context.Background(), types.DriverOptions{IncludeDeleted: true})
			ts.Equal(test.err, err)
			testDeletedApp, ok := test.portalAppsBeforeDelete[test.deleteParams.ID]
			ts.True(ok)
			testDeletedApp.Deleted = true
			ts.Equal(test.portalAppsBeforeDelete, portalApps)
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
			portalApp:       *testdata.PortalApps["test_app_create_208r23r"],
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

			portalApps, err := ts.driver.ReadPortalApps(context.Background(), types.DriverOptions{})
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
		testUpdatedNotifications map[types.NotificationType]types.AppNotification
		testUpdatedWhitelists    types.Whitelists
		err                      error
	}{
		{
			name: "Should update a new PortalApp in the database with all fields",
			updatePortalApp: types.UpdatePortalApp{
				Name:          testdata.UpdatePortalAppName,
				Settings:      testdata.UpdatePortalAppSettings,
				Notifications: testdata.UpdatePortalAppNotifications,
				Whitelists:    testdata.UpdatePortalAppWhitelists,
			},
			testUpdateTime:  testdata.MockTimestamp,
			testUpdatedName: "portal-app-updated",
			testUpdatedSettings: types.Settings{
				Environment:       types.EnvironmentProduction,
				SecretKey:         "test_9d07c8a96ad53e7c288b0e86f37c5680",
				SecretKeyRequired: true,
				MonthlyRelayLimit: 2_500_000,
				FavoritedChainIDs: map[types.ChainID]struct{}{"0003": {}, "0009": {}, "00H3": {}},
			},
			testUpdatedNotifications: map[types.NotificationType]types.AppNotification{
				types.NotificationTypeEmail: {
					Active:      true,
					Destination: "user@example.com",
					Trigger:     "daily",
					Events: map[types.NotificationEvent]bool{
						types.NotificationEventSignedUp:      true,
						types.NotificationEventHalf:          true,
						types.NotificationEventQuarter:       true,
						types.NotificationEventThreeQuarters: true,
						types.NotificationEventFull:          true,
					},
				},
				types.NotificationTypeWebhook: {
					Active:      true,
					Destination: "https://example.com/webhook",
					Trigger:     "hourly",
					Events: map[types.NotificationEvent]bool{
						types.NotificationEventHalf: true,
						types.NotificationEventFull: true,
					},
				},
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
			name: "Should update a new PortalApp in the database with only a new Name",
			updatePortalApp: types.UpdatePortalApp{
				Name: testdata.UpdatePortalAppName,
			},
			testUpdateTime:  testdata.MockTimestamp,
			testUpdatedName: "portal-app-updated",
			err:             nil,
		},
		{
			name: "Should update a new PortalApp in the database with only new Settings",
			updatePortalApp: types.UpdatePortalApp{
				Settings: testdata.UpdatePortalAppSettings,
			},
			testUpdateTime: testdata.MockTimestamp,
			testUpdatedSettings: types.Settings{
				Environment:       types.EnvironmentProduction,
				SecretKey:         "test_9d07c8a96ad53e7c288b0e86f37c5680",
				SecretKeyRequired: true,
				MonthlyRelayLimit: 2_500_000,
				FavoritedChainIDs: map[types.ChainID]struct{}{"0003": {}, "0009": {}, "00H3": {}},
			},
			err: nil,
		},
		{
			name: "Should update a new PortalApp in the database with only new Notifications",
			updatePortalApp: types.UpdatePortalApp{
				Notifications: testdata.UpdatePortalAppNotifications,
			},
			testUpdateTime: testdata.MockTimestamp,
			testUpdatedNotifications: map[types.NotificationType]types.AppNotification{
				types.NotificationTypeEmail: {
					Active:      true,
					Destination: "user@example.com",
					Trigger:     "daily",
					Events: map[types.NotificationEvent]bool{
						types.NotificationEventSignedUp:      true,
						types.NotificationEventHalf:          true,
						types.NotificationEventQuarter:       true,
						types.NotificationEventThreeQuarters: true,
						types.NotificationEventFull:          true,
					},
				},
				types.NotificationTypeWebhook: {
					Active:      true,
					Destination: "https://example.com/webhook",
					Trigger:     "hourly",
					Events: map[types.NotificationEvent]bool{
						types.NotificationEventHalf: true,
						types.NotificationEventFull: true,
					},
				},
			},
			err: nil,
		},
		{
			name: "Should update a new PortalApp in the database with only new Whitelists",
			updatePortalApp: types.UpdatePortalApp{
				Whitelists: testdata.UpdatePortalAppWhitelists,
			},
			testUpdateTime: testdata.MockTimestamp,
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

	for i, test := range tests {
		ts.Run(test.name, func() {
			// Create new portal app for test case
			createApp := *testdata.PortalApps["test_app_update_b03ca84c"]
			createApp.Name = fmt.Sprintf("test-update-portal-app-%d", i+1)
			createdPortalApp, err := ts.driver.WritePortalApp(context.Background(), createApp, testdata.MockTimestamp)
			ts.NoError(err)

			// Update created portal app for test case
			updateApp := test.updatePortalApp
			updateApp.AppID = createdPortalApp.ID
			err = ts.driver.UpdatePortalApp(context.Background(), updateApp, test.testUpdateTime)
			ts.Equal(test.err, err)

			// Get all portal apps from DB
			portalApps, err := ts.driver.ReadPortalApps(context.Background(), types.DriverOptions{})
			ts.Equal(test.err, err)
			updatedPortalApp, ok := portalApps[createdPortalApp.ID]
			ts.True(ok)

			// Check update changes
			if test.testUpdatedName != "" {
				ts.Equal(test.testUpdatedName, updatedPortalApp.Name)
			} else {
				ts.Equal(createdPortalApp.Name, updatedPortalApp.Name)
			}

			if test.testUpdatedSettings.Environment != "" {
				ts.Equal(test.testUpdatedWhitelists, updatedPortalApp.Whitelists)
			} else {
				ts.Equal(createdPortalApp.Settings, updatedPortalApp.Settings)
			}

			if len(test.testUpdatedNotifications) != 0 {
				ts.Equal(test.testUpdatedNotifications, updatedPortalApp.Notifications)
			} else {
				ts.Equal(createdPortalApp.Notifications, updatedPortalApp.Notifications)
			}

			if len(test.testUpdatedWhitelists.Origins) != 0 {
				ts.Equal(test.testUpdatedWhitelists, updatedPortalApp.Whitelists)
			} else {
				ts.Equal(createdPortalApp.Whitelists, updatedPortalApp.Whitelists)
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_UpdatePortalAppsFirstDateSurpassed() {
	tests := []struct {
		name   string
		update *types.UpdateFirstDateSurpassed
		err    error
	}{
		{
			name: "Should update multiple PortalApps in the database with a LegacyFields.FirstDateSurpassed timestamp ",
			update: &types.UpdateFirstDateSurpassed{
				PortalAppIDs:       []string{string(testdata.PortalApps["test_app_3487u329rfn23f9"].ID), string(testdata.PortalApps["test_app_2308rj09r23r9r2"].ID)},
				FirstDateSurpassed: time.Date(2023, time.February, 14, 0, 0, 0, 0, time.UTC),
			},
			err: nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			// Update FirstDateSurpassed fields for PortalApps
			err := ts.driver.UpdatePortalAppsFirstDateSurpassed(context.Background(), test.update)
			ts.Equal(test.err, err)

			// Check FirstDateSurpassed fields for PortalApps in the DB
			portalApps, err := ts.driver.ReadPortalApps(context.Background(), types.DriverOptions{})
			ts.Equal(test.err, err)
			for _, appID := range test.update.PortalAppIDs {
				portalApp, ok := portalApps[types.PortalAppID(appID)]
				ts.True(ok)
				ts.Equal(test.update.FirstDateSurpassed, portalApp.LegacyFields.FirstDateSurpassed)
			}
		})
	}
}
