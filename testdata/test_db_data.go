package testdata

import (
	"time"

	"github.com/pokt-foundation/portal-db/v2/types"
)

/*
This file contains the Go struct equivalent to the seed test data for the test database container defined in `seed_test_db.sql`.

This mock data can be used anywhere that used the test database Docker container `pocketfoundation/test-portal-postgres`.

If/when adding new data to the `seed_test_db.sql` file it should also be updated here.
This means all E2E tests that use `pocketfoundation/test-portal-postgres` only need to be kept updated in one place.
*/

var (
	MockTimestamp = time.Date(2022, time.November, 11, 11, 11, 11, 0, time.UTC)

	/* ----- Read Data ----- */

	PayPlans = map[types.PayPlanType]*types.Plan{
		types.PayPlanType("basic_plan"): {
			Type:              types.PayPlanType("basic_plan"),
			ChainIDs:          map[types.RelayChainID]struct{}{"0001": {}, "0053": {}},
			MonthlyRelayLimit: 5_000_000,
			ThroughputLimit:   5_000,
			AppLimit:          2,
			LegacyDailyLimit:  1_000,
			CreatedAt:         MockTimestamp,
		},
		types.PayPlanType("pro_plan"): {
			Type:              types.PayPlanType("pro_plan"),
			ChainIDs:          map[types.RelayChainID]struct{}{"0001": {}, "0053": {}, "0021": {}, "0064": {}},
			MonthlyRelayLimit: 10_000_000,
			ThroughputLimit:   10_000,
			AppLimit:          5,
			LegacyDailyLimit:  5_000,
			CreatedAt:         MockTimestamp,
		},
		types.PayPlanType("enterprise_plan"): {
			Type:              types.PayPlanType("enterprise_plan"),
			ChainIDs:          map[types.RelayChainID]struct{}{"0001": {}, "0053": {}, "0021": {}, "0064": {}, "0034": {}},
			MonthlyRelayLimit: 20_000_000,
			ThroughputLimit:   20_000,
			AppLimit:          10,
			LegacyDailyLimit:  10_000,
			CreatedAt:         MockTimestamp,
		},
		types.PayPlanType("developer_plan"): {
			Type:              types.PayPlanType("developer_plan"),
			ChainIDs:          map[types.RelayChainID]struct{}{"0001": {}, "0053": {}, "0021": {}, "0034": {}},
			MonthlyRelayLimit: 500_000,
			ThroughputLimit:   500,
			AppLimit:          1,
			LegacyDailyLimit:  100,
			CreatedAt:         MockTimestamp,
		},
		types.PayPlanType("startup_plan"): {
			Type:              types.PayPlanType("startup_plan"),
			ChainIDs:          map[types.RelayChainID]struct{}{"0001": {}, "0053": {}, "0064": {}, "0034": {}},
			MonthlyRelayLimit: 1_000_000,
			ThroughputLimit:   1_000,
			AppLimit:          5,
			LegacyDailyLimit:  500,
			CreatedAt:         MockTimestamp,
		},
	}

	Accounts = map[types.AccountID]*types.Account{
		"account_1": {
			ID:       "account_1",
			PlanType: types.PayPlanType("basic_plan"),
			Users: map[types.UserID]types.AccountUserAccess{
				"user_1": AccountUserAccess[1],
				"user_2": AccountUserAccess[2],
				"user_8": AccountUserAccess[8],
			},
			Integrations: types.AccountIntegrations{
				CovalentAPIKeyFree: "covalent_api_key_1",
			},
			PartnerChainIDs:        map[types.RelayChainID]struct{}{"0001": {}, "0053": {}},
			PartnerThroughputLimit: 2_000,
			PartnerAppLimit:        1,
			CreatedAt:              MockTimestamp,
			UpdatedAt:              MockTimestamp,
		},
		"account_2": {
			ID:       "account_2",
			PlanType: types.PayPlanType("pro_plan"),
			Users: map[types.UserID]types.AccountUserAccess{
				"user_3": AccountUserAccess[3],
				"user_4": AccountUserAccess[4],
				"user_9": AccountUserAccess[9],
				"user_2": AccountUserAccess[10],
			},
			Integrations: types.AccountIntegrations{
				CovalentAPIKeyFree: "covalent_api_key_2",
			},
			PartnerChainIDs:        map[types.RelayChainID]struct{}{"0001": {}, "0053": {}, "0021": {}, "0064": {}},
			PartnerThroughputLimit: 5_000,
			PartnerAppLimit:        3,
			CreatedAt:              MockTimestamp,
			UpdatedAt:              MockTimestamp,
		},
		"account_3": {
			ID:       "account_3",
			PlanType: types.PayPlanType("startup_plan"),
			Users: map[types.UserID]types.AccountUserAccess{
				"user_5":  AccountUserAccess[5],
				"user_6":  AccountUserAccess[6],
				"user_7":  AccountUserAccess[7],
				"user_10": AccountUserAccess[12],
			},
			Integrations: types.AccountIntegrations{
				CovalentAPIKeyFree: "covalent_api_key_3",
			},
			PartnerChainIDs:        map[types.RelayChainID]struct{}{"0001": {}, "0053": {}, "0064": {}, "0034": {}},
			PartnerThroughputLimit: 1_000,
			PartnerAppLimit:        2,
			CreatedAt:              MockTimestamp,
			UpdatedAt:              MockTimestamp,
		},
		"account_4": {
			ID:       "account_4",
			PlanType: types.PayPlanType("enterprise_plan"),
			Users: map[types.UserID]types.AccountUserAccess{
				"user_4": AccountUserAccess[11],
			},
			Integrations: types.AccountIntegrations{
				CovalentAPIKeyFree: "covalent_api_key_4",
			},
			PartnerChainIDs:        map[types.RelayChainID]struct{}{"0001": {}},
			PartnerThroughputLimit: 1_000,
			PartnerAppLimit:        2,
			CreatedAt:              MockTimestamp,
			UpdatedAt:              MockTimestamp,
		},
		"account_5": {
			ID:       "account_5",
			PlanType: types.PayPlanType("basic_plan"),
			Users: map[types.UserID]types.AccountUserAccess{
				"user_4": AccountUserAccess[11],
			},
			PartnerChainIDs:        map[types.RelayChainID]struct{}{"0006": {}, "0040": {}},
			PartnerThroughputLimit: 6_000,
			PartnerAppLimit:        1,
			CreatedAt:              MockTimestamp,
			UpdatedAt:              MockTimestamp,
		},
	}

	// TestCreateAccount account used to test creation of Accounts
	TestCreateAccount = &types.Account{
		PlanType:  types.PayPlanType("developer_plan"),
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
	}

	AccountUserAccess = map[int]types.AccountUserAccess{
		1: {UserID: "user_1", Email: "james.holden123@test.com", RoleName: types.RoleOwner, Accepted: true, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|james_holden", types.AuthTypeAuth0Github: "github|james_holden"}},
		2: {UserID: "user_2", Email: "paul.atreides456@test.com", RoleName: types.RoleAdmin, Accepted: true, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|paul_atreides", types.AuthTypeAuth0Github: "github|paul_atreides"}},
		3: {UserID: "user_3", Email: "ellen.ripley789@test.com", RoleName: types.RoleOwner, Accepted: true, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|ellen_ripley"}},
		4: {UserID: "user_4", Email: "ulfric.stormcloak123@test.com", RoleName: types.RoleMember, Accepted: true, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|ulfric_stormcloak"}},
		5: {UserID: "user_5", Email: "chrisjen.avasarala1@test.com", RoleName: types.RoleOwner, Accepted: true, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|chrisjen_avasarala"}},
		6: {UserID: "user_6", Email: "amos.burton789@test.com", RoleName: types.RoleAdmin, Accepted: true, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|amos_burton"}},
		7: {UserID: "user_7", Email: "frodo.baggins123@test.com", RoleName: types.RoleMember, Accepted: true, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|frodo_baggins"}},
		8: {UserID: "user_8", Email: "rick.deckard456@test.com", RoleName: types.RoleAdmin, Accepted: false, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|rick_deckard"}},
		9: {UserID: "user_9", Email: "tyrion.lannister789@test.com", RoleName: types.RoleMember, Accepted: false, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|tyrion_lannister"}},
		// Paul is an admin of Account 1 as well as a member of Account 2
		10: {UserID: "user_2", Email: "paul.atreides456@test.com", RoleName: types.RoleMember, Accepted: true, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|paul_atreides", types.AuthTypeAuth0Github: "github|paul_atreides"}},
		// Ulfric is an admin of Account 2 as well as the owner of Accounts 4 and 5
		11: {UserID: "user_4", Email: "ulfric.stormcloak123@test.com", RoleName: types.RoleOwner, Accepted: true, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|ulfric_stormcloak"}},
		// Daenerys has not signed up with an auth provider yet and is a member of Account 3
		12: {UserID: "user_10", Email: "daenerys.targaryen123@test.com", RoleName: types.RoleMember, Accepted: false},
		// Bernard is an existing user and is used to create a new AccountUserAccess row
		13: {UserID: "user_11", Email: "bernard.marx@test.com", Accepted: false, RoleName: types.RoleMember, ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|bernard_marx"}},
		// Winston has not signed up yet and is used to create a new AccountUserAccess row
		14: {Email: "winston.smith@test.com", RoleName: types.RoleAdmin, Accepted: false},
	}

	Users = map[types.UserID]*types.User{
		"user_1": {
			ID:       "user_1",
			Email:    "james.holden123@test.com",
			SignedUp: true,
			AuthProviders: map[types.AuthType]types.UserAuthProvider{
				types.AuthTypeAuth0Username: {
					ProviderUserID: "auth0|james_holden",
					Type:           types.AuthTypeAuth0Username,
					Provider:       "auth0",
					Federated:      false,
				},
				types.AuthTypeAuth0Github: {
					ProviderUserID: "github|james_holden",
					Type:           types.AuthTypeAuth0Github,
					Provider:       "auth0",
					Federated:      true,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"user_2": {
			ID:       "user_2",
			Email:    "paul.atreides456@test.com",
			SignedUp: true,
			AuthProviders: map[types.AuthType]types.UserAuthProvider{
				types.AuthTypeAuth0Username: {
					ProviderUserID: "auth0|paul_atreides",
					Type:           types.AuthTypeAuth0Username,
					Provider:       "auth0",
					Federated:      false,
				},
				types.AuthTypeAuth0Github: {
					ProviderUserID: "github|paul_atreides",
					Type:           types.AuthTypeAuth0Github,
					Provider:       "auth0",
					Federated:      true,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"user_3": {
			ID:       "user_3",
			Email:    "ellen.ripley789@test.com",
			SignedUp: true,
			AuthProviders: map[types.AuthType]types.UserAuthProvider{
				types.AuthTypeAuth0Username: {
					ProviderUserID: "auth0|ellen_ripley",
					Type:           types.AuthTypeAuth0Username,
					Provider:       "auth0",
					Federated:      false,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"user_4": {
			ID:       "user_4",
			Email:    "ulfric.stormcloak123@test.com",
			SignedUp: true,
			AuthProviders: map[types.AuthType]types.UserAuthProvider{
				types.AuthTypeAuth0Username: {
					ProviderUserID: "auth0|ulfric_stormcloak",
					Type:           types.AuthTypeAuth0Username,
					Provider:       "auth0",
					Federated:      false,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"user_5": {
			ID:       "user_5",
			Email:    "chrisjen.avasarala1@test.com",
			SignedUp: true,
			AuthProviders: map[types.AuthType]types.UserAuthProvider{
				types.AuthTypeAuth0Username: {
					ProviderUserID: "auth0|chrisjen_avasarala",
					Type:           types.AuthTypeAuth0Username,
					Provider:       "auth0",
					Federated:      false,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"user_6": {
			ID:       "user_6",
			Email:    "amos.burton789@test.com",
			SignedUp: true,
			AuthProviders: map[types.AuthType]types.UserAuthProvider{
				types.AuthTypeAuth0Username: {
					ProviderUserID: "auth0|amos_burton",
					Type:           types.AuthTypeAuth0Username,
					Provider:       "auth0",
					Federated:      false,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"user_7": {
			ID:       "user_7",
			Email:    "frodo.baggins123@test.com",
			SignedUp: true,
			AuthProviders: map[types.AuthType]types.UserAuthProvider{
				types.AuthTypeAuth0Username: {
					ProviderUserID: "auth0|frodo_baggins",
					Type:           types.AuthTypeAuth0Username,
					Provider:       "auth0",
					Federated:      false,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"user_8": {
			ID:       "user_8",
			Email:    "rick.deckard456@test.com",
			SignedUp: true,
			AuthProviders: map[types.AuthType]types.UserAuthProvider{
				types.AuthTypeAuth0Username: {
					ProviderUserID: "auth0|rick_deckard",
					Type:           types.AuthTypeAuth0Username,
					Provider:       "auth0",
					Federated:      false,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"user_9": {
			ID:       "user_9",
			Email:    "tyrion.lannister789@test.com",
			SignedUp: true,
			AuthProviders: map[types.AuthType]types.UserAuthProvider{
				types.AuthTypeAuth0Username: {
					ProviderUserID: "auth0|tyrion_lannister",
					Type:           types.AuthTypeAuth0Username,
					Provider:       "auth0",
					Federated:      false,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"user_10": {
			ID:            "user_10",
			Email:         "daenerys.targaryen123@test.com",
			SignedUp:      false,
			AuthProviders: map[types.AuthType]types.UserAuthProvider{},
			CreatedAt:     MockTimestamp,
			UpdatedAt:     MockTimestamp,
		},
	}

	TestCreateUser = types.CreateUser{
		Email:            "commander.data@example.com",
		AuthProviderType: types.AuthTypeAuth0Username,
		ProviderUserID:   "auth0|commander_data",
	}

	UserPermissions = map[types.UserID]*types.UserPermissions{
		"user_1": {
			UserID: "user_1",
			Accounts: map[types.AccountID]types.AccountPermissions{
				"account_1": types.AccountPermissions{
					RoleName: types.RoleOwner,
					Permissions: []types.Permissions{
						types.PermReadEndpoint,
						types.PermWriteEndpoint,
						types.PermDeleteEndpoint,
						types.PermTransferEndpoint,
					},
				},
			},
		},
		"user_2": {
			UserID: "user_2",
			Accounts: map[types.AccountID]types.AccountPermissions{
				"account_1": types.AccountPermissions{
					RoleName: types.RoleAdmin,
					Permissions: []types.Permissions{
						types.PermReadEndpoint,
						types.PermWriteEndpoint,
					},
				},
				"account_2": types.AccountPermissions{
					RoleName: types.RoleMember,
					Permissions: []types.Permissions{
						types.PermReadEndpoint,
					},
				},
			},
		},
		"user_3": {
			UserID: "user_3",
			Accounts: map[types.AccountID]types.AccountPermissions{
				"account_2": types.AccountPermissions{
					RoleName: types.RoleOwner,
					Permissions: []types.Permissions{
						types.PermReadEndpoint,
						types.PermWriteEndpoint,
						types.PermDeleteEndpoint,
						types.PermTransferEndpoint,
					},
				},
			},
		},
		"user_4": {
			UserID: "user_4",
			Accounts: map[types.AccountID]types.AccountPermissions{
				"account_2": types.AccountPermissions{
					RoleName: types.RoleMember,
					Permissions: []types.Permissions{
						types.PermReadEndpoint,
					},
				},
				"account_4": types.AccountPermissions{
					RoleName: types.RoleOwner,
					Permissions: []types.Permissions{
						types.PermReadEndpoint,
						types.PermWriteEndpoint,
						types.PermDeleteEndpoint,
						types.PermTransferEndpoint,
					},
				},
				"account_5": types.AccountPermissions{
					RoleName: types.RoleOwner,
					Permissions: []types.Permissions{
						types.PermReadEndpoint,
						types.PermWriteEndpoint,
						types.PermDeleteEndpoint,
						types.PermTransferEndpoint,
					},
				},
			},
		},
		"user_5": {
			UserID: "user_5",
			Accounts: map[types.AccountID]types.AccountPermissions{
				"account_3": types.AccountPermissions{
					RoleName: types.RoleOwner,
					Permissions: []types.Permissions{
						types.PermReadEndpoint,
						types.PermWriteEndpoint,
						types.PermDeleteEndpoint,
						types.PermTransferEndpoint,
					},
				},
			},
		},
		"user_6": {
			UserID: "user_6",
			Accounts: map[types.AccountID]types.AccountPermissions{
				"account_3": types.AccountPermissions{
					RoleName: types.RoleAdmin,
					Permissions: []types.Permissions{
						types.PermReadEndpoint,
						types.PermWriteEndpoint,
					},
				},
			},
		},
		"user_7": {
			UserID: "user_7",
			Accounts: map[types.AccountID]types.AccountPermissions{
				"account_3": types.AccountPermissions{
					RoleName: types.RoleMember,
					Permissions: []types.Permissions{
						types.PermReadEndpoint,
					},
				},
			},
		},
	}

	UserIDs = map[types.ProviderUserID]types.UserID{
		"auth0|amos_burton":        "user_6",
		"auth0|bernard_marx":       "user_11",
		"auth0|chrisjen_avasarala": "user_5",
		"auth0|ellen_ripley":       "user_3",
		"auth0|frodo_baggins":      "user_7",
		"auth0|james_holden":       "user_1",
		"auth0|paul_atreides":      "user_2",
		"auth0|rick_deckard":       "user_8",
		"auth0|tyrion_lannister":   "user_9",
		"auth0|ulfric_stormcloak":  "user_4",
		"github|james_holden":      "user_1",
		"github|paul_atreides":     "user_2",
	}

	PortalApps = map[types.PortalAppID]*types.PortalApp{
		"test_app_1": {
			ID:        "test_app_1",
			AccountID: "account_1",
			Name:      "pokt_app_123",
			Gigastake: true,
			Staked:    false,
			AATs: map[types.ProtocolAppID]types.AAT{
				"test_protocol_app_1": {
					ID:              "test_protocol_app_1",
					Address:         "test_34715cae753e67c75fbb340442e7de8e",
					PublicKey:       "test_34715cae753e67c75fbb340442e7de8e",
					ClientPublicKey: "test_89a3af6a587aec02cfade6f5000424c2",
					PrivateKey:      "test_11b8d394ca331d7c7a71ca1896d630f6",
					Signature:       "test_1dc39a2e5a84a35bf030969a0b3231f7",
					Version:         "0.0.1",
				},
			},
			Settings: types.Settings{
				Environment:       types.EnvironmentProduction,
				SecretKey:         "test_40f482d91a5ef2300ebb4e2308c",
				SecretKeyRequired: true,
			},
			Whitelists: types.Whitelists{
				Origins:     map[types.Origin]struct{}{"https://test.com": {}},
				UserAgents:  map[types.UserAgent]struct{}{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)": {}},
				Blockchains: map[types.RelayChainID]struct{}{"0053": {}},
				Contracts: map[types.RelayChainID]map[types.Contract]struct{}{
					"0001": {"0x1234567890abcdef": {}},
				},
				Methods: map[types.RelayChainID]map[types.Method]struct{}{
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
		},
		"test_app_2": {
			ID:        "test_app_2",
			AccountID: "account_2",
			Name:      "pokt_app_456",
			Gigastake: false,
			Staked:    true,
			AATs: map[types.ProtocolAppID]types.AAT{
				"test_protocol_app_2": {
					ID:              "test_protocol_app_2",
					Address:         "test_8237c72345f12d1b1a8b64a1a7f66fa4",
					PublicKey:       "test_8237c72345f12d1b1a8b64a1a7f66fa4",
					ClientPublicKey: "test_04c71d90a92f40416b6f1d7d8af17e02",
					PrivateKey:      "test_2e83c836a29b423a47d8e18c779fd422",
					Signature:       "test_f48d33b30ddaf60a1e5bb50d2ba8da5a",
					Version:         "0.0.1",
				},
			},
			Settings: types.Settings{
				Environment:       types.EnvironmentProduction,
				SecretKey:         "test_9c9e3b193cfba5348f93bb2f3e3fb794",
				SecretKeyRequired: false,
			},
			Whitelists: types.Whitelists{
				Origins:     map[types.Origin]struct{}{"https://example.com": {}},
				UserAgents:  map[types.UserAgent]struct{}{"Mozilla/5.0 (Linux; Android 10; SM-A205U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36": {}},
				Blockchains: map[types.RelayChainID]struct{}{"0021": {}},
				Contracts: map[types.RelayChainID]map[types.Contract]struct{}{
					"0064": {"0x0987654321abcdef": {}},
				},
				Methods: map[types.RelayChainID]map[types.Method]struct{}{
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
		},
		"test_app_3": {
			ID:        "test_app_3",
			AccountID: "account_3",
			Name:      "pokt_app_789",
			Gigastake: false,
			Staked:    true,
			AATs: map[types.ProtocolAppID]types.AAT{
				"test_protocol_app_3": {
					ID:              "test_protocol_app_3",
					Address:         "test_b5e07928fc80083c13ad0201b81bae9b",
					PublicKey:       "test_f608500e4fe3e09014fe2411b4a560b5",
					ClientPublicKey: "test_328a9cf1b35085eeaa669aa858f6fba9",
					PrivateKey:      "test_8663e187c19f3c6e27317eab4ed6d7d5",
					Signature:       "test_c3cd8be16ba32e24dd49fdb0247fc9b8",
					Version:         "0.0.1",
				},
				"test_protocol_app_4": {
					ID:              "test_protocol_app_4",
					Address:         "test_eb2e5bcba557cfe8fa76fd7fff54f9d1",
					PublicKey:       "test_f6a5d8690ecb669865bd752b7796a920",
					ClientPublicKey: "test_6ee5ea553408f0895923fd1569dc5072",
					PrivateKey:      "test_838d29d61a65401f7d56d084cb6e4783",
					Signature:       "test_cf05cf9bb26111c548e88fb6157af708",
					Version:         "0.0.1",
				},
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
				CustomLimit:        0,
				RequestTimeout:     10_000,
				GigastakeRedirect:  false,
				FirstDateSurpassed: MockTimestamp,
			},
		},
	}

	// TestCreatePortalApp app used to test creation of PortalApps
	TestCreatePortalApp = &types.PortalApp{
		ID:        "test_app_create_208r23r",
		AccountID: "account_4",
		Name:      "create_pokt_app_1",
		Gigastake: true,
		Staked:    false,

		Settings: types.Settings{
			Environment:       types.EnvironmentProduction,
			SecretKey:         "test_3e3fb7949c9e3b193cfba5348f93bb2f",
			SecretKeyRequired: true,
		},
		Notifications: map[types.NotificationType]types.AppNotification{
			types.NotificationTypeEmail: types.AppNotification{
				Active:      true,
				Destination: "ulfric.stormcloak123@test.com",
				Events: map[types.NotificationEvent]bool{
					types.NotificationEventSignedUp:      true,
					types.NotificationEventThreeQuarters: true,
					types.NotificationEventFull:          true,
				},
			},
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
		// TODO remove legacy fields when migration to V2 schema complete
		LegacyFields: types.LegacyFields{
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

	TestCreatePortalAppAAT = types.AAT{
		Address:         "test_1a8b64a1a7f66fa48237c72345f12dgr",
		PublicKey:       "test_8237c72345f1a7f66fa41b1b8b644g2f",
		ClientPublicKey: "test_d4222e83c836a29b423a47d8e18c779f",
		PrivateKey:      "test_a92f40416b6f1d7d8af17e0204c71d90",
		Signature:       "test_da5af48d33b30ddaf60a1e5bb50d2b8f",
		Version:         "0.0.1",
	}

	// TestUpdatePortalApp app used to test updates of PortalApps
	TestUpdatePortalApp = &types.PortalApp{
		ID:        "test_app_update_b03ca84c",
		AccountID: "account_1",
		Name:      "", // name set in test
		Gigastake: true,
		Staked:    false,
		AATs: map[types.ProtocolAppID]types.AAT{
			"test_protocol_app_1": {
				ID:              "test_protocol_app_1",
				Address:         "test_7d0cd2743543a6200e41224594954b06",
				PublicKey:       "test_7d0cd2743543a6200e41224594954b06",
				ClientPublicKey: "test_3d2b1cf05bd9b479b6fd65b9ffdf1976",
				PrivateKey:      "test_9c59143368436aeee593c2e6cdbda57b",
				Signature:       "test_a8546957653d23e3b2e76bb718099e7a",
				Version:         "0.0.1",
			},
		},
		Settings: types.Settings{
			Environment:       types.EnvironmentProduction,
			SecretKey:         "test_849c1397586f9fb6f902576120d0d10f",
			SecretKeyRequired: true,
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
		LegacyFields: types.LegacyFields{
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

	Chains = map[types.RelayChainID]*types.Chain{
		"0001": {
			ID:            "0001",
			Blockchain:    "mainnet",
			Description:   "Pocket Network Mainnet",
			EnforceResult: "JSON",
			Path:          "/v1/query/height",
			Ticker:        "POKT",
			ChainAliases:  []string{"mainnet"},
			Active:        true,
			Altruists: []types.Altruist{
				{
					URL:      "https://altruist-0001.com:1234",
					AuthType: types.ChainAuthTypeBasicAuth,
					Auth:     "test_pocket:auth123456",
				},
			},
			Redirects: []types.GigastakeRedirect{
				{PortalApplicationID: "test_app_1", Alias: "altruist-0001", Domain: "pokt-rpc.gateway.pokt.network"},
			},
			Checks: map[types.ChainCheckType]types.Check{
				types.ChainCheckTypeSync: {
					Type:      types.ChainCheckTypeSync,
					Payload:   `{"id":1,"jsonrpc":"2.0","method":"query"}`,
					ResultKey: "result.sync_info",
					Allowance: 1,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"0053": {
			ID:             "0053",
			Blockchain:     "optimism-mainnet",
			Description:    "Optimism Mainnet",
			EnforceResult:  "JSON",
			Ticker:         "OP",
			ChainAliases:   []string{"optimism-mainnet"},
			LogLimitBlocks: 100000,
			RequestTimeout: 0,
			Active:         true,
			Altruists: []types.Altruist{
				{
					URL:      "https://altruist-0053.com:1234",
					AuthType: types.ChainAuthTypeBasicAuth,
					Auth:     "test_pocket:auth123456",
				},
			},
			Redirects: []types.GigastakeRedirect{
				{PortalApplicationID: "test_app_2", Alias: "altruist-0053", Domain: "op-rpc.gateway.pokt.network"},
			},
			Checks: map[types.ChainCheckType]types.Check{
				types.ChainCheckTypeSync: {
					Type:      types.ChainCheckTypeSync,
					Payload:   `{"id":1,"jsonrpc":"2.0","method":"eth_blockNumber","params":[]}`,
					ResultKey: "result",
					Allowance: 2,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"0021": {
			ID:             "0021",
			Blockchain:     "eth-mainnet",
			Description:    "Ethereum Mainnet",
			EnforceResult:  "JSON",
			Ticker:         "ETH",
			ChainAliases:   []string{"eth-mainnet"},
			LogLimitBlocks: 100000,
			Active:         true,
			Altruists: []types.Altruist{
				{
					URL:      "https://altruist-0021.com:1234",
					AuthType: types.ChainAuthTypeBasicAuth,
					Auth:     "test_pocket:auth123456",
				},
			},
			Redirects: []types.GigastakeRedirect{
				{PortalApplicationID: "test_app_3", Alias: "altruist-0021", Domain: "eth-rpc.gateway.pokt.network"},
			},
			Checks: map[types.ChainCheckType]types.Check{
				types.ChainCheckTypeSync: {
					Type:      types.ChainCheckTypeSync,
					Payload:   `{"id":1,"jsonrpc":"2.0","method":"eth_blockNumber","params":[]}`,
					ResultKey: "result",
					Allowance: 5,
				},
				types.ChainCheckTypeChain: {
					Type:       types.ChainCheckTypeChain,
					Payload:    `{"method":"eth_chainId","id":1,"jsonrpc":"2.0"}`,
					ResultKey:  "id",
					EVMChainID: 1,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"0064": {
			ID:             "0064",
			Blockchain:     "sui-testnet",
			Description:    "Sui Testnet",
			EnforceResult:  "JSON",
			Ticker:         "SUI-TESTNET",
			ChainAliases:   []string{"sui-testnet"},
			LogLimitBlocks: 100000,
			RequestTimeout: 60000,
			Active:         false,
			Altruists: []types.Altruist{
				{
					URL:      "https://altruist-0064.com:1234",
					AuthType: types.ChainAuthTypeBasicAuth,
					Auth:     "test_pocket:auth123456",
				},
			},
			Redirects: []types.GigastakeRedirect{},
			Checks: map[types.ChainCheckType]types.Check{
				types.ChainCheckTypeSync: {
					Type:      types.ChainCheckTypeSync,
					Payload:   `{"id":1,"jsonrpc":"2.0","method":"sui_blockNumber","params":[]}`,
					ResultKey: "result",
					Allowance: 7,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
		"0040": {
			ID:            "0040",
			Blockchain:    "harmony-0",
			Description:   "Harmony Shard 0",
			EnforceResult: "JSON",
			Ticker:        "HMY",
			ChainAliases:  []string{"harmony-0"},
			Active:        true,
			Altruists: []types.Altruist{
				{
					URL:      "https://altruist-0040.com:1234",
					AuthType: types.ChainAuthTypeBasicAuth,
					Auth:     "test_pocket:auth123456",
				},
			},
			Redirects: []types.GigastakeRedirect{
				{PortalApplicationID: "test_app_3", Alias: "altruist-0040", Domain: "hmy-rpc.gateway.pokt.network"},
			},
			Checks: map[types.ChainCheckType]types.Check{
				types.ChainCheckTypeSync: {
					Type:      types.ChainCheckTypeSync,
					Payload:   `{"id":1,"jsonrpc":"2.0","method":"hmy_blockNumber","params":[]}`,
					ResultKey: "result",
					Allowance: 8,
				},
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
	}

	// TestCreateChain used to test creation of Chains
	TestCreateChain = &types.Chain{
		ID:            "0006",
		Blockchain:    "solana-mainnet",
		Description:   "Solana",
		EnforceResult: "JSON",
		Ticker:        "SOL",
		ChainAliases:  []string{"solana-mainnet"},
		Altruists: []types.Altruist{
			{
				URL:      "https://test-rpc.solana-1.io:1234",
				AuthType: types.ChainAuthTypeBasicAuth,
				Auth:     "test_pocket:auth123456",
			},
			{
				URL:      "https://test-rpc.solana-2.io:1234",
				AuthType: types.ChainAuthTypeBasicAuth,
				Auth:     "test_pocket:auth123456",
			},
		},
		Redirects: []types.GigastakeRedirect{
			{PortalApplicationID: "test_app_1", Alias: "solana-mainnet", Domain: "sol-rpc.gateway.pokt.network"},
		},
		Checks: map[types.ChainCheckType]types.Check{
			types.ChainCheckTypeSync: {
				Type:      types.ChainCheckTypeSync,
				Payload:   `{"id":1,"jsonrpc":"2.0","method":"getSync"}`,
				ResultKey: "sync",
				Allowance: 2,
			},
			types.ChainCheckTypeChain: {
				Type:       types.ChainCheckTypeChain,
				Payload:    `{"id":1,"jsonrpc":"2.0","method":"getChain"}`,
				ResultKey:  "chain",
				EVMChainID: 5,
			},
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
	}

	GlobalBlockedContracts = types.GlobalBlockedContracts{
		BlockedAddresses: map[types.BlockedAddress]struct{}{
			"0xtest_6789abcdef0123456789abcdef01234567":   {},
			"0xtest_f0123456789abcdef0123456789abcdef01":  {},
			"0xtest_cdef0123456789abcdef0123456789abcdef": {},
			"0xtest_56789abcdef0123456789abcdef01234567":  {},
			"0xtest_789abcdef0123456789abcdef0123456789":  {},
		},
	}

	/* ----- Update Data ----- */

	UpdatePortalAppName     = "portal-app-updated"
	UpdatePortalAppSettings = &types.UpdateAppSettings{
		Environment:       types.EnvironmentProduction,
		SecretKey:         "test_9d07c8a96ad53e7c288b0e86f37c5680",
		SecretKeyRequired: true,
		MonthlyRelayLimit: 2_500_000,
		FavoritedChainIDs: []string{"0003", "0009", "00H3"},
	}
	UpdatePortalAppNotifications = []types.UpdateAppNotifications{
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
	UpdatePortalAppWhitelists = &types.WhitelistsObject{
		AppWhitelists: [3]types.ApplicationWhitelists{
			{Type: "origins", Values: []string{"https://portalgun.io", "https://subdomain.example.com", "https://www.example.com"}},
			{Type: "userAgents", Values: []string{"Brave", "Google Chrome", "Mozilla Firefox", "Netscape Navigator", "Safari"}},
			{Type: "blockchains", Values: []string{"0001", "0002", "003E", "0056"}},
		},
		ChainWhitelists: [2]types.ChainWhitelists{
			{Type: "contracts", Values: []types.ChainIDWhitelists{
				{ChainID: "0001", Values: []string{"0xtest_2f78db6436527729929aaf6c616361de0f7", "0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be"}},
				{ChainID: "0002", Values: []string{"0xtest_1111117dc0aa78b770fa6a738034120c302", "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2"}},
				{ChainID: "003E", Values: []string{"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984", "0xtest_f958d2ee523a2206206994597c13d831ec7"}},
				{ChainID: "0056", Values: []string{"0xtest_00000f279d81a1d3cc75430faa017fa5a2e", "0xtest_5068778dd592e39a122f4f5a5cf09c90fe2"}},
			},
			},
			{Type: "methods", Values: []types.ChainIDWhitelists{
				{ChainID: "0001", Values: []string{"GET", "POST", "PUT"}},
				{ChainID: "0002", Values: []string{"DELETE", "GET", "POST", "PUT"}},
				{ChainID: "003E", Values: []string{"GET"}},
				{ChainID: "0056", Values: []string{"GET", "POST"}},
			},
			},
		},
	}

	UpdateChainOne = types.Chain{
		ID:            "0001",
		Blockchain:    "mainnet-NEW",
		Description:   "Pocket Network Mainnet Update",
		EnforceResult: "JSON",
		Path:          "/v1/query/height/wow",
		Ticker:        "POKT-123",
		ChainAliases:  []string{"mainnet"},
		Active:        true,
		Altruists: []types.Altruist{
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
			{
				URL:      "https://altruist-0001-3.com:1234",
				AuthType: types.ChainAuthTypeBasicAuth,
				Auth:     "test_pocket:auth123456",
			},
		},
		Redirects: []types.GigastakeRedirect{
			{PortalApplicationID: "test_app_1", Alias: "altruist-0001", Domain: "pokt-rpc.gateway.pokt.network"},
		},
		Checks: map[types.ChainCheckType]types.Check{
			types.ChainCheckTypeSync: {
				Type:      types.ChainCheckTypeSync,
				Payload:   `{"id":1,"jsonrpc":"2.0","method":"query"}`,
				ResultKey: "result.sync_info",
				Allowance: 1,
			},
			types.ChainCheckTypeChain: {
				Type:       types.ChainCheckTypeChain,
				Payload:    `{"id":1,"jsonrpc":"2.0","method":"chain"}`,
				ResultKey:  "result.sync_info",
				EVMChainID: 3,
			},
			types.ChainCheckTypeMerge: {
				Type:      types.ChainCheckTypeMerge,
				Payload:   `{"id":1,"jsonrpc":"2.0","method":"merge"}`,
				ResultKey: "result.sync_info",
			},
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
	}
	UpdateChainTwo = types.Chain{
		ID:            "0001",
		Blockchain:    "mainnet-ULTRA",
		Description:   "Pocket Network Mainnet Original",
		EnforceResult: "JSON",
		Path:          "/v1/query/height/wow",
		Ticker:        "POKT-456",
		ChainAliases:  []string{"mainnet-again"},
		Active:        true,
		Altruists: []types.Altruist{
			{
				URL:      "https://altruist-0001-3.com:1234",
				AuthType: types.ChainAuthTypeBasicAuth,
				Auth:     "test_pocket:auth123456",
			},
		},
		Redirects: []types.GigastakeRedirect{
			{PortalApplicationID: "test_app_1", Alias: "altruist-0001", Domain: "pokt-rpc.gateway.pokt.network"},
		},
		Checks: map[types.ChainCheckType]types.Check{
			types.ChainCheckTypeSync: {
				Type:      types.ChainCheckTypeSync,
				Payload:   `{"id":1,"jsonrpc":"2.0","method":"query"}`,
				ResultKey: "result.sync_info",
				Allowance: 1,
			},
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
	}
	UpdateChainThree = types.Chain{
		ID:            "0001",
		Blockchain:    "mainnet-ULTRA",
		Description:   "Pocket Network Mainnet Original",
		EnforceResult: "JSON",
		Path:          "/v1/query/height/wow",
		Ticker:        "POKT-456",
		ChainAliases:  []string{"mainnet-again"},
		Active:        true,
		CreatedAt:     MockTimestamp,
		UpdatedAt:     MockTimestamp,
	}

	// UpdateChainNotExists used to test updating a Chain that doesn't exist
	UpdateChainNotExists = types.Chain{ID: "0073"}

	/* ----- Legacy Data ----- */

	V2Account = &types.Account{
		ID:       "account_1",
		PlanType: types.PayPlanType("basic_plan"),
		Users: map[types.UserID]types.AccountUserAccess{
			"user_1": AccountUserAccess[1],
			"user_2": AccountUserAccess[2],
			"user_8": AccountUserAccess[8],
		},
		PartnerChainIDs:        map[types.RelayChainID]struct{}{"0001": {}, "0053": {}},
		PartnerThroughputLimit: 2_000,
		PartnerAppLimit:        1,
		CreatedAt:              MockTimestamp,
		UpdatedAt:              MockTimestamp,

		// Assume Plan and PortalApps have been set in PHD
		Plan: PayPlans["basic_plan"],
		PortalApps: map[types.PortalAppID]*types.PortalApp{
			"test_app_1": PortalApps["test_app_1"],
		},
	}

	LegacyLoadBalancer = types.LoadBalancer{
		ID:                "test_app_1",
		Name:              "pokt_app_123",
		UserID:            "auth0|james_holden",
		ApplicationIDs:    []string(nil),
		RequestTimeout:    5000,
		Gigastake:         true,
		GigastakeRedirect: true,
		StickyOptions: types.StickyOptions{
			Duration:      "60",
			StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
			StickyMax:     300,
			Stickiness:    true,
		},
		Applications: LegacyApplications,
		Users: []types.UserAccess{
			{
				UserID:   "user_1",
				RoleName: types.RoleOwner,
				Email:    "james.holden123@test.com",
				Accepted: true,
			},
			{
				UserID:   "user_2",
				RoleName: types.RoleAdmin,
				Email:    "paul.atreides456@test.com",
				Accepted: true,
			},
			{
				UserID:   "user_8",
				RoleName: types.RoleAdmin,
				Email:    "rick.deckard456@test.com",
				Accepted: false,
			},
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
	}

	LegacyApplications = []*types.Application{
		{
			ID:                 "test_protocol_app_1",
			UserID:             "auth0|james_holden",
			Name:               "pokt_app_123",
			FirstDateSurpassed: MockTimestamp,
			GatewayAAT: types.GatewayAAT{
				Address:              "test_34715cae753e67c75fbb340442e7de8e",
				ApplicationPublicKey: "test_34715cae753e67c75fbb340442e7de8e",
				ApplicationSignature: "test_1dc39a2e5a84a35bf030969a0b3231f7",
				ClientPublicKey:      "test_89a3af6a587aec02cfade6f5000424c2",
				PrivateKey:           "test_11b8d394ca331d7c7a71ca1896d630f6",
				Version:              "0.0.1",
			},
			GatewaySettings: types.GatewaySettings{
				SecretKey:            "test_40f482d91a5ef2300ebb4e2308c",
				SecretKeyRequired:    true,
				WhitelistOrigins:     []string{"https://test.com"},
				WhitelistUserAgents:  []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
				WhitelistBlockchains: []string{"0053"},
				WhitelistContracts:   []types.WhitelistContracts{{BlockchainID: "0001", Contracts: []string{"0x1234567890abcdef"}}},
				WhitelistMethods:     []types.WhitelistMethods{{BlockchainID: "0001", Methods: []string{"GET"}}},
			},
			Limit: types.AppLimit{
				PayPlan:     types.PayPlan{Type: "basic_plan", Limit: 1_000},
				CustomLimit: 0,
			},
			NotificationSettings: types.NotificationSettings{
				SignedUp:      false,
				Quarter:       true,
				Half:          false,
				ThreeQuarters: true,
				Full:          true,
			},
			CreatedAt: MockTimestamp,
			UpdatedAt: MockTimestamp,
		},
	}

	LegacyBlockchain = types.Blockchain{
		ID:                "0001",
		Altruist:          "https://test_pocket:auth123456@altruist-0001.com:1234", // pragma: allowlist secret
		Blockchain:        "mainnet",
		Description:       "Pocket Network Mainnet",
		EnforceResult:     "JSON",
		Path:              "/v1/query/height",
		Ticker:            "POKT",
		BlockchainAliases: []string{"mainnet"},
		Active:            true,
		Redirects: []types.Redirect{
			{
				LoadBalancerID: "test_app_1",
				Alias:          "altruist-0001",
				Domain:         "pokt-rpc.gateway.pokt.network",
			},
		},
		SyncCheckOptions: types.SyncCheckOptions{
			Body:      "{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"query\"}",
			ResultKey: "result.sync_info",
			Allowance: 1,
		},
		CreatedAt: MockTimestamp,
		UpdatedAt: MockTimestamp,
	}

	LegacyUpdateApplication = types.UpdateApplication{
		Name: "test_portal_app_123",
		GatewaySettings: &types.UpdateGatewaySettings{
			SecretKey:            "test_90210ac4bdd3423e24877d1ff92",
			SecretKeyRequired:    boolToPointer(false),
			WhitelistOrigins:     []string{"https://portalgun.io", "https://subdomain.example.com", "https://www.example.com"},
			WhitelistBlockchains: []string{"0001", "0002", "003E", "0056"},
			WhitelistUserAgents:  []string{"Brave", "Google Chrome", "Mozilla Firefox", "Netscape Navigator", "Safari"},
			WhitelistContracts: []types.WhitelistContracts{
				{BlockchainID: "0001", Contracts: []string{"0xtest_2f78db6436527729929aaf6c616361de0f7", "0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be"}},
				{BlockchainID: "0002", Contracts: []string{"0xtest_1111117dc0aa78b770fa6a738034120c302", "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2"}},
				{BlockchainID: "003E", Contracts: []string{"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984", "0xtest_f958d2ee523a2206206994597c13d831ec7"}},
				{BlockchainID: "0056", Contracts: []string{"0xtest_00000f279d81a1d3cc75430faa017fa5a2e", "0xtest_5068778dd592e39a122f4f5a5cf09c90fe2"}},
			},
			WhitelistMethods: []types.WhitelistMethods{
				{BlockchainID: "0001", Methods: []string{"GET", "POST", "PUT"}},
				{BlockchainID: "0002", Methods: []string{"DELETE", "GET", "POST", "PUT"}},
				{BlockchainID: "003E", Methods: []string{"GET"}},
				{BlockchainID: "0056", Methods: []string{"GET", "POST"}},
			},
		},
		Limit:                &types.AppLimit{PayPlan: types.PayPlan{Type: types.FreetierV0, Limit: 250_000}, CustomLimit: 0},
		NotificationSettings: &types.UpdateNotificationSettings{SignedUp: boolToPointer(true), Quarter: boolToPointer(true), Half: boolToPointer(false), ThreeQuarters: boolToPointer(true), Full: boolToPointer(false)},
	}

	V2CreateAccount = types.Account{
		PlanType: types.PayPlanType("basic_plan"),
	}

	V2CreatePortalApp = types.PortalApp{
		Name:      "pokt_app_123",
		Gigastake: true,
		Staked:    false,
		AccountID: "",
		Settings: types.Settings{
			Environment:       "production",
			SecretKey:         "test_40f482d91a5ef2300ebb4e2308c",
			SecretKeyRequired: false,
		},
		Notifications: map[types.NotificationType]types.AppNotification{
			types.NotificationTypeEmail: {
				Active:      true,
				Destination: "james.holden123@test.com",
				Events: map[types.NotificationEvent]bool{
					types.NotificationEventFull:          true,
					types.NotificationEventHalf:          false,
					types.NotificationEventQuarter:       true,
					types.NotificationEventSignedUp:      false,
					types.NotificationEventThreeQuarters: true,
				},
			},
		},
		LegacyFields: types.LegacyFields{
			RequestTimeout:    5000,
			GigastakeRedirect: true,
			StickyOptions: types.StickyOptions{
				Duration:      "60",
				StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
				StickyMax:     300,
				Stickiness:    true,
			},
		},
	}

	V2CreatePortalAppAAT = types.AAT{
		Address:         "test_34715cae753e67c75fbb340442e7de8e",
		PublicKey:       "test_34715cae753e67c75fbb340442e7de8e",
		ClientPublicKey: "test_89a3af6a587aec02cfade6f5000424c2",
		PrivateKey:      "test_11b8d394ca331d7c7a71ca1896d630f6",
		Signature:       "test_1dc39a2e5a84a35bf030969a0b3231f7",
		Version:         "0.0.1",
	}

	V2UpdatePortalApp = types.UpdatePortalApp{
		AppID:    "test_app_1",
		Name:     "test_portal_app_123",
		Settings: &types.UpdateAppSettings{SecretKey: "test_90210ac4bdd3423e24877d1ff92"},
		Notifications: []types.UpdateAppNotifications{
			{
				NotificationType: types.NotificationTypeEmail,
				Events:           []types.NotificationEvent{"signedUp", "quarter", "threeQuarters"},
			},
		},
		Whitelists: &types.WhitelistsObject{
			AppWhitelists: [3]types.ApplicationWhitelists{
				{Type: "origins", Values: []string{"https://portalgun.io", "https://subdomain.example.com", "https://www.example.com"}},
				{Type: "userAgents", Values: []string{"Brave", "Google Chrome", "Mozilla Firefox", "Netscape Navigator", "Safari"}},
				{Type: "blockchains", Values: []string{"0001", "0002", "003E", "0056"}},
			},
			ChainWhitelists: [2]types.ChainWhitelists{
				{
					Type: "contracts",
					Values: []types.ChainIDWhitelists{
						{ChainID: "0001", Values: []string{"0xtest_2f78db6436527729929aaf6c616361de0f7", "0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be"}},
						{ChainID: "0002", Values: []string{"0xtest_1111117dc0aa78b770fa6a738034120c302", "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2"}},
						{ChainID: "003E", Values: []string{"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984", "0xtest_f958d2ee523a2206206994597c13d831ec7"}},
						{ChainID: "0056", Values: []string{"0xtest_00000f279d81a1d3cc75430faa017fa5a2e", "0xtest_5068778dd592e39a122f4f5a5cf09c90fe2"}},
					},
				},
				{
					Type: "methods",
					Values: []types.ChainIDWhitelists{
						{ChainID: "0001", Values: []string{"GET", "POST", "PUT"}},
						{ChainID: "0002", Values: []string{"DELETE", "GET", "POST", "PUT"}},
						{ChainID: "003E", Values: []string{"GET"}},
						{ChainID: "0056", Values: []string{"GET", "POST"}},
					},
				},
			},
		},
	}

	LegacyUpdateBlockchain = types.UpdateBlockchain{
		Blockchain:        "mainnet-ULTRA",
		Description:       "Pocket Network Mainnet Original",
		EnforceResult:     "JSON",
		Path:              "/v1/query/height/wow",
		Ticker:            "POKT-456",
		BlockchainAliases: []string{"mainnet-again"},
		Body:              `{"id":1,"jsonrpc":"2.0","method":"query"}`,
		ResultKey:         "result.sync_info",
		Allowance:         intToPointer(1),
		Altruist:          "https://test_pocket:auth123456@altruist-0001-3.com:1234", // pragma: allowlist secret
	}

	LegacyRedirect = types.Redirect{
		BlockchainID:   "0001",
		LoadBalancerID: "test_lb_5c6f50bc30b530a8",
		Domain:         "pokt-rpc.gateway.pokt.network",
		Alias:          "altruist-0001",
	}

	LegacyUserAccess = types.UserAccess{
		Email:    "james.holden123@test.com",
		RoleName: types.RoleOwner,
		Accepted: true,
		UserID:   "james_holden_push_button",
	}

	LegacyUpdateUserAccess = types.UpdateUserAccess{
		UserID:   "test_user_c66f399519ba23",
		RoleName: types.RoleAdmin,
		Email:    "test_admin@user.com",
	}
)

func boolToPointer(boolVar bool) *bool {
	return &boolVar
}

func intToPointer(intVar int) *int {
	return &intVar
}
