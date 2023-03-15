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

	/* Read Data */

	TestPortalAppOne = types.PortalApp{
		ID:        "test_app_3487u329rfn23f9",
		AccountID: 1,
		Name:      "pokt_app_123",
		Gigastake: true,
		Staked:    false,
		AAT: types.AAT{
			Address:         "test_34715cae753e67c75fbb340442e7de8e",
			PublicKey:       "test_34715cae753e67c75fbb340442e7de8e",
			ClientPublicKey: "test_89a3af6a587aec02cfade6f5000424c2",
			PrivateKey:      "test_11b8d394ca331d7c7a71ca1896d630f6",
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
			Blockchains: map[types.ChainID]struct{}{"0053": {}},
			Contracts: map[types.ChainID]map[types.Contract]struct{}{
				"0001": {"0x1234567890abcdef": {}},
			},
			Methods: map[types.ChainID]map[types.Method]struct{}{
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
			ApplicationIDs:     []string{"test_app_47hfnths73j2se7"},
			CustomLimit:        0,
			RequestTimeout:     5_000,
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
			ClientPublicKey: "test_04c71d90a92f40416b6f1d7d8af17e02",
			PrivateKey:      "test_2e83c836a29b423a47d8e18c779fd422",
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
			Blockchains: map[types.ChainID]struct{}{"0021": {}},
			Contracts: map[types.ChainID]map[types.Contract]struct{}{
				"0064": {"0x0987654321abcdef": {}},
			},
			Methods: map[types.ChainID]map[types.Method]struct{}{
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
			ApplicationIDs:     []string{"test_app_43jr9304urj30fj"},
			CustomLimit:        0,
			RequestTimeout:     10_000,
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

	TestPortalAppThree = types.PortalApp{
		ID:        "test_app_47fhs7j4hs7fj2",
		AccountID: 1,
		Name:      "pokt_app_789",
		Gigastake: false,
		Staked:    true,
		AAT: types.AAT{
			Address:         "test_b5e07928fc80083c13ad0201b81bae9b",
			PublicKey:       "test_f608500e4fe3e09014fe2411b4a560b5",
			ClientPublicKey: "test_328a9cf1b35085eeaa669aa858f6fba9",
			PrivateKey:      "test_8663e187c19f3c6e27317eab4ed6d7d5",
			Signature:       "test_c3cd8be16ba32e24dd49fdb0247fc9b8",
			Version:         "0.0.1",
		},
		Settings: types.Settings{
			Environment:       types.EnvironmentProduction,
			SecretKey:         "test_9f48b13e2bc5fd31ab367841f11495c1",
			SecretKeyRequired: false,
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
		// TODO remove legacy fields when migration to V2 schema complete
		LegacyFields: types.LegacyFields{
			ApplicationIDs:     []string{"test_app_43jr947fh23dfg4"},
			CustomLimit:        0,
			RequestTimeout:     10_000,
			GigastakeRedirect:  false,
			FirstDateSurpassed: MockTimestamp,
		},
	}

	/* Write/Create Data */

	TestCreatePortalAppOne = types.PortalApp{
		ID:        "test_app_rj09fjw208r23r",
		AccountID: 1,
		Name:      "create_pokt_app_1",
		Gigastake: true,
		Staked:    false,
		AAT: types.AAT{
			Address:         "test_1a8b64a1a7f66fa48237c72345f12dgr",
			PublicKey:       "test_8237c72345f1a7f66fa41b1b8b644g2f",
			ClientPublicKey: "test_d4222e83c836a29b423a47d8e18c779f",
			PrivateKey:      "test_a92f40416b6f1d7d8af17e0204c71d90",
			Signature:       "test_da5af48d33b30ddaf60a1e5bb50d2b8f",
			Version:         "0.0.1",
		},
		Settings: types.Settings{
			Environment:       types.EnvironmentProduction,
			SecretKey:         "test_3e3fb7949c9e3b193cfba5348f93bb2f",
			SecretKeyRequired: true,
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
		// TODO remove legacy fields when migration to V2 schema complete
		LegacyFields: types.LegacyFields{
			ApplicationIDs:     []string{"test_app_30fj43jr94j9304"},
			CustomLimit:        750_000,
			RequestTimeout:     15_000,
			GigastakeRedirect:  true,
			FirstDateSurpassed: MockTimestamp,
			StickyOptions: types.StickyOptions{
				Duration:      "60",
				StickyOrigins: []string{"https://pokt.network", "https://example.com"},
				StickyMax:     1200,
				Stickiness:    true,
			},
		},
	}

	/* Update Data */

	TestCreateAppForUpdate = types.PortalApp{
		ID:        "test_app_539a081b03ca84c",
		AccountID: 1,
		Gigastake: true,
		Staked:    false,
		AAT: types.AAT{
			Address:         "test_7d0cd2743543a6200e41224594954b06",
			PublicKey:       "test_7d0cd2743543a6200e41224594954b06",
			ClientPublicKey: "test_3d2b1cf05bd9b479b6fd65b9ffdf1976",
			PrivateKey:      "test_9c59143368436aeee593c2e6cdbda57b",
			Signature:       "test_a8546957653d23e3b2e76bb718099e7a",
			Version:         "0.0.1",
		},
		Settings: types.Settings{
			Environment:       types.EnvironmentProduction,
			SecretKey:         "test_849c1397586f9fb6f902576120d0d10f",
			SecretKeyRequired: true,
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
		LegacyFields: types.LegacyFields{
			ApplicationIDs:     []string{"test_app_bb162dd67c99615"},
			CustomLimit:        0,
			RequestTimeout:     5_000,
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

	TestPortalAppName     = "portal-app-updated"
	TestPortalAppSettings = &types.UpdateAppSettings{
		Environment:       types.EnvironmentProduction,
		SecretKey:         "test_9d07c8a96ad53e7c288b0e86f37c5680",
		SecretKeyRequired: true,
		MonthlyRelayLimit: 2_500_000,
		FavoritedChainIDs: []string{"0003", "0009", "00H3"},
	}
	TestPortalAppNotifications = []types.UpdateAppNotifications{
		{
			NotificationType: types.NotificationTypeEmail,
			Active:           true,
			Destination:      "user@example.com",
			Trigger:          "daily",
			Events: []types.NotificationEvent{
				types.NotificationEventSignedUp,
				types.NotificationEventHalf,
				types.NotificationEventQuarter,
				types.NotificationEventThreeQuarters,
				types.NotificationEventFull,
			},
		},
		{
			NotificationType: types.NotificationTypeWebhook,
			Active:           true,
			Destination:      "https://example.com/webhook",
			Trigger:          "hourly",
			Events: []types.NotificationEvent{
				types.NotificationEventHalf,
				types.NotificationEventFull,
			},
		},
		{NotificationType: types.NotificationTypePortal, Active: false},
	}
	TestPortalAppWhitelists = &types.WhitelistsObject{
		AppWhitelists: [3]types.ApplicationWhitelists{
			{Type: "origins", Values: []string{"https://portalgun.io", "https://subdomain.example.com", "https://www.example.com"}},
			{Type: "userAgents", Values: []string{"Brave", "Google Chrome", "Mozilla Firefox", "Netscape Navigator", "Safari"}},
			{Type: "blockchains", Values: []string{"0001", "0002", "003E", "0056"}},
		},
		ChainWhitelists: [2]types.ChainWhitelists{
			{Type: "contracts", Values: []types.BlockchainIDWhitelists{
				{ChainID: "0001", Values: []string{"0xtest_2f78db6436527729929aaf6c616361de0f7", "0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be"}},
				{ChainID: "0002", Values: []string{"0xtest_1111117dc0aa78b770fa6a738034120c302", "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2"}},
				{ChainID: "003E", Values: []string{"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984", "0xtest_f958d2ee523a2206206994597c13d831ec7"}},
				{ChainID: "0056", Values: []string{"0xtest_00000f279d81a1d3cc75430faa017fa5a2e", "0xtest_5068778dd592e39a122f4f5a5cf09c90fe2"}},
			},
			},
			{Type: "methods", Values: []types.BlockchainIDWhitelists{
				{ChainID: "0001", Values: []string{"GET", "POST", "PUT"}},
				{ChainID: "0002", Values: []string{"DELETE", "GET", "POST", "PUT"}},
				{ChainID: "003E", Values: []string{"GET"}},
				{ChainID: "0056", Values: []string{"GET", "POST"}},
			},
			},
		},
	}
)
