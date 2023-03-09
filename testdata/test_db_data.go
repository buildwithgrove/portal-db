package testdata

import (
	"time"

	"github.com/pokt-foundation/portal-db/types"
)

/*
This file contains the Go struct equivalent to the seed test data for the test database container defined in `seed_test_db.sql`.

This mock data can be used anywhere that used the test database Docker container `pocketfoundation/test-portal-postgres`.

If/when adding new data to the `seed_test_db.sql` file it should also be updated here.
This means all E2E tests that use `pocketfoundation/test-portal-postgres` only need to be kept updated in one place.
*/

var (
	MockTimestamp = time.Date(2022, time.November, 11, 11, 11, 11, 0, time.UTC)

	TestPortalAppOne = types.PortalApp{
		ID:        "test_app_3487u329rfn23f9",
		AccountID: 1,
		Name:      "pokt_app_123",
		Gigastake: true,
		Staked:    false,
		AAT: types.AAT{
			Address:         "test_34715cae753e67c75fbb340442e7de8e",
			PublicKey:       "test_34715cae753e67c75fbb340442e7de8e",
			ClientPublicKey: "test_11b8d394ca331d7c7a71ca1896d630f6",
			PrivateKey:      "test_89a3af6a587aec02cfade6f5000424c2",
			Signature:       "test_1dc39a2e5a84a35bf030969a0b3231f7",
			Version:         "0.0.1",
		},
		Settings: types.Settings{
			Environment:       types.EnvironmentProduction,
			SecretKey:         "test_40f482d91a5ef2300ebb4e2308c",
			SecretKeyRequired: true,
		},
		Whitelists: types.Whitelists{
			Origins:     map[types.Origin]struct{}{"https://example.com": {}},
			UserAgents:  map[types.UserAgent]struct{}{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)": {}},
			Blockchains: map[types.BlockchainID]struct{}{"0053": {}},
			Contracts: map[types.BlockchainID]map[types.Contract]struct{}{
				"0001": {"0x1234567890abcdef": {}},
			},
			Methods: map[types.BlockchainID]map[types.Method]struct{}{
				"0001": {"GET": {}},
			},
		},
		Notifications: map[types.NotificationType]types.AppNotification{
			types.NotificationTypeEmail: {
				Active:      true,
				Destination: "test@test.com",
				Trigger:     "trigger123",
				Events: map[types.NotificationEvent]bool{
					types.NotificationEventFull:          true,
					types.NotificationEventQuarter:       true,
					types.NotificationEventThreeQuarters: true,
				},
			},
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
		// TODO remove legacy fields when migration to V2 schema complete
		LegacyFields: types.LegacyFields{
			ApplicationID:      "test_app_47hfnths73j2se7",
			CustomLimit:        0,
			RequestTimeout:     5000,
			GigastakeRedirect:  true,
			FirstDateSurpassed: MockTimestamp,
			StickyOptions: types.StickyOptions{
				Duration:      "60",
				StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
				StickyMax:     300,
				Stickiness:    true,
			},
		},
	}

	TestPortalAppTwo = types.PortalApp{
		ID:        "test_app_2308rj09r23r9r",
		AccountID: 1,
		Name:      "pokt_app_456",
		Gigastake: false,
		Staked:    true,
		AAT: types.AAT{
			Address:         "test_8237c72345f12d1b1a8b64a1a7f66fa4",
			PublicKey:       "test_8237c72345f12d1b1a8b64a1a7f66fa4",
			ClientPublicKey: "test_2e83c836a29b423a47d8e18c779fd422",
			PrivateKey:      "test_04c71d90a92f40416b6f1d7d8af17e02",
			Signature:       "test_f48d33b30ddaf60a1e5bb50d2ba8da5a",
			Version:         "0.0.1",
		},
		Settings: types.Settings{
			Environment:       types.EnvironmentProduction,
			SecretKey:         "test_9c9e3b193cfba5348f93bb2f3e3fb794",
			SecretKeyRequired: false,
		},
		Whitelists: types.Whitelists{
			Origins:     map[types.Origin]struct{}{"https://test.com": {}},
			UserAgents:  map[types.UserAgent]struct{}{"Mozilla/5.0 (Linux; Android 10; SM-A205U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36": {}},
			Blockchains: map[types.BlockchainID]struct{}{"0021": {}},
			Contracts: map[types.BlockchainID]map[types.Contract]struct{}{
				"0064": {"0x0987654321abcdef": {}},
			},
			Methods: map[types.BlockchainID]map[types.Method]struct{}{
				"0064": {"POST": {}},
			},
		},
		Notifications: map[types.NotificationType]types.AppNotification{
			types.NotificationTypeEmail: {
				Active:      true,
				Destination: "email@pokt.network",
				Trigger:     "trigger456",
				Events: map[types.NotificationEvent]bool{
					types.NotificationEventHalf: true,
					types.NotificationEventFull: true,
				},
			},
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
		// TODO remove legacy fields when migration to V2 schema complete
		LegacyFields: types.LegacyFields{
			ApplicationID:      "test_app_43jr9304urj30fj",
			CustomLimit:        0,
			RequestTimeout:     10000,
			GigastakeRedirect:  false,
			FirstDateSurpassed: MockTimestamp,
			StickyOptions: types.StickyOptions{
				Duration:      "30",
				StickyOrigins: []string{"https://example.com", "https://test.com"},
				StickyMax:     600,
				Stickiness:    true,
			},
		},
	}
)
