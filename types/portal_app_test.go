package types

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testPortalApplication = PortalApp{
	ID:        "test_5416bb8d696386455b8",
	Name:      "test_portal_app_123",
	Gigastake: true,
	Account: Account{
		Plan: Plan{
			Type:              FreetierV0,
			MonthlyRelayLimit: 2_500_000,
			ThroughputLimit:   2_000,
			AppLimit:          2,
			BlockchainIDs:     map[BlockchainID]struct{}{"0001": {}, "0056": {}, "0002": {}, "003E": {}},
			LegacyDailyLimit:  250_000,
		},
		Users: map[UserID]AccountUserAccess{
			"user_id_123": {
				User:     User{ID: "user_id_123", Email: "test_owner@user.com", AuthProvider: AuthProviderAuth0},
				RoleName: RoleOwner,
				Accepted: true,
			},
			"user_id_456": {
				User:     User{ID: "user_id_456", Email: "test_member@user.com", AuthProvider: AuthProviderAuth0},
				RoleName: RoleMember,
				Accepted: true,
			},
		},
		PartnerBlockchainIDs:   map[BlockchainID]struct{}{"0001": {}, "0056": {}, "0002": {}, "003E": {}},
		PartnerThroughputLimit: 2_000,
		PartnerAppLimit:        2,
	},
	AAT: AAT{
		Address:         "test_34715cae753e67c75fbb340442e7de8e",
		PublicKey:       "test_11b8d394ca331d7c7a71ca1896d630f6",
		ClientPublicKey: "test_9e9ca4fe13725d412003f4bc518f6974",
		PrivateKey:      "test_89a3af6a587aec02cfade6f5000424c2",
		Signature:       "test_1dc39a2e5a84a35bf030969a0b3231f7",
		Version:         "0.0.1",
	},
	Settings: Settings{
		Environment:            EnvironmentProduction,
		SecretKey:              "test_90210ac4bdd3423e24877d1ff92",
		SecretKeyRequired:      true,
		MonthlyRelayLimit:      250_000,
		FavoritedBlockchainIDs: map[BlockchainID]struct{}{"0001": {}, "0056": {}, "0002": {}, "003E": {}},
	},
	Whitelists: Whitelists{
		Origins:     map[Origin]struct{}{"https://www.example.com": {}, "https://subdomain.example.com": {}, "https://portalgun.io": {}},
		UserAgents:  map[UserAgent]struct{}{"Mozilla Firefox": {}, "Brave": {}, "Google Chrome": {}, "Safari": {}, "Netscape Navigator": {}},
		Blockchains: map[BlockchainID]struct{}{"0001": {}, "0056": {}, "0002": {}, "003E": {}},
		Contracts: map[BlockchainID]map[Contract]struct{}{
			"0001": {"0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be": {}, "0xtest_2f78db6436527729929aaf6c616361de0f7": {}},
			"0056": {"0xtest_5068778dd592e39a122f4f5a5cf09c90fe2": {}, "0xtest_00000f279d81a1d3cc75430faa017fa5a2e": {}},
			"0002": {"0xtest_1111117dc0aa78b770fa6a738034120c302": {}, "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2": {}},
			"003E": {"0xtest_f958d2ee523a2206206994597c13d831ec7": {}, "0xtest_0a85d5af5bf1d1762f925bdaddc4201f984": {}},
		},
		Methods: map[BlockchainID]map[Method]struct{}{
			"0001": {"GET": {}, "POST": {}, "PUT": {}},
			"0056": {"GET": {}, "POST": {}},
			"0002": {"GET": {}, "POST": {}, "PUT": {}, "DELETE": {}},
			"003E": {"GET": {}},
		},
	},
	Notifications: map[NotificationType]AppNotification{
		NotificationTypeEmail: {
			Active: true, Destination: "test@user.com", Trigger: "what_am_i", // TODO what is trigger?
			Events: map[NotificationEvent]bool{
				"signedUp":      true,
				"quarter":       true,
				"threeQuarters": true,
			},
		},
		NotificationTypeWebhook: {
			Active: false, Destination: "https://wh.destination.io", Trigger: "what_am_i",
			Events: map[NotificationEvent]bool{
				"signedUp": true,
				"full":     true,
			},
		},
	},
	CreatedAt: time.Date(2023, time.February, 14, 11, 11, 11, 0, time.UTC),
	UpdatedAt: time.Date(2023, time.February, 27, 13, 13, 13, 0, time.UTC),

	LegacyFields: LegacyFields{
		ApplicationID:      "test_475jf893f9j2f30jd230e",
		CustomLimit:        0,
		RequestTimeout:     5_000,
		GigastakeRedirect:  true,
		FirstDateSurpassed: time.Date(2023, time.February, 28, 15, 15, 15, 0, time.UTC),
		StickyOptions: StickyOptions{
			Duration:      "4000",
			StickyOrigins: []string{"origin123"},
			StickyMax:     4_000,
			Stickiness:    true,
		},
	},
}

func Test_PortalApp_IsOriginWhitelisted(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name           string
		portalApp      PortalApp
		origin         Origin
		expectedResult bool
	}{
		{
			name:           "Should return true if a given origin is whitelisted for a given app",
			portalApp:      testPortalApplication,
			origin:         Origin("https://portalgun.io"),
			expectedResult: true,
		},
		{
			name:           "Should return true if a given origin is whitelisted for a given app",
			portalApp:      testPortalApplication,
			origin:         Origin("https://www.example.com"),
			expectedResult: true,
		},
		{
			name:           "Should return false if a given origin is not whitelisted for a given app",
			portalApp:      testPortalApplication,
			origin:         Origin("https://ricksanchez.io"),
			expectedResult: false,
		},
	}

	for _, test := range tests {
		isOriginWhitelisted := test.portalApp.IsOriginWhitelisted(test.origin)
		c.Equal(test.expectedResult, isOriginWhitelisted)
	}
}

func Test_PortalApp_IsUserAgentWhitelisted(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name           string
		portalApp      PortalApp
		userAgent      UserAgent
		expectedResult bool
	}{
		{
			name:           "Should return true if a given user agent is whitelisted for a given app",
			portalApp:      testPortalApplication,
			userAgent:      UserAgent("Brave"),
			expectedResult: true,
		},
		{
			name:           "Should return true if a given user agent is whitelisted for a given app",
			portalApp:      testPortalApplication,
			userAgent:      UserAgent("Safari"),
			expectedResult: true,
		},
		{
			name:           "Should return false if a given user agent is not whitelisted for a given app",
			portalApp:      testPortalApplication,
			userAgent:      UserAgent("Bird Person"),
			expectedResult: false,
		},
	}

	for _, test := range tests {
		isUserAgentWhitelisted := test.portalApp.IsUserAgentWhitelisted(test.userAgent)
		c.Equal(test.expectedResult, isUserAgentWhitelisted)
	}
}

func Test_PortalApp_IsBlockchainWhitelisted(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name           string
		portalApp      PortalApp
		blockchain     BlockchainID
		expectedResult bool
	}{
		{
			name:           "Should return true if a given blockchain is whitelisted for a given app",
			portalApp:      testPortalApplication,
			blockchain:     BlockchainID("0001"),
			expectedResult: true,
		},
		{
			name:           "Should return true if a given blockchain is whitelisted for a given app",
			portalApp:      testPortalApplication,
			blockchain:     BlockchainID("003E"),
			expectedResult: true,
		},
		{
			name:           "Should return false if a given blockchain is not whitelisted for a given app",
			portalApp:      testPortalApplication,
			blockchain:     BlockchainID("7009"),
			expectedResult: false,
		},
	}

	for _, test := range tests {
		isBlockchainWhitelisted := test.portalApp.IsBlockchainWhitelisted(test.blockchain)
		c.Equal(test.expectedResult, isBlockchainWhitelisted)
	}
}

func Test_PortalApp_IsContractWhitelisted(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name           string
		portalApp      PortalApp
		blockchain     BlockchainID
		contract       Contract
		expectedResult bool
	}{
		{
			name:           "Should return true if a given contract is whitelisted for a given app and blockchain",
			portalApp:      testPortalApplication,
			blockchain:     BlockchainID("0001"),
			contract:       Contract("0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be"),
			expectedResult: true,
		},
		{
			name:           "Should return true if a given contract is whitelisted for a given app and blockchain",
			portalApp:      testPortalApplication,
			blockchain:     BlockchainID("003E"),
			contract:       Contract("0xtest_0a85d5af5bf1d1762f925bdaddc4201f984"),
			expectedResult: true,
		},
		{
			name:           "Should return false if a given contract is not whitelisted for a given app and blockchain",
			portalApp:      testPortalApplication,
			blockchain:     BlockchainID("0056"),
			contract:       Contract("0xtest_04938rfj439fj3409jf0439fjf4304f4444"),
			expectedResult: false,
		},
		{
			name:           "Should return false if a given contract is not whitelisted for a given app and blockchain",
			portalApp:      testPortalApplication,
			blockchain:     BlockchainID("7009"),
			contract:       Contract("0xtest_439834fnin3f2032f03re3j2f30fj33f3f3"),
			expectedResult: false,
		},
	}

	for _, test := range tests {
		isBlockchainWhitelisted := test.portalApp.IsContractWhitelisted(test.blockchain, test.contract)
		c.Equal(test.expectedResult, isBlockchainWhitelisted)
	}
}

func Test_PortalApp_IsMethodWhitelisted(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name           string
		portalApp      PortalApp
		blockchain     BlockchainID
		method         Method
		expectedResult bool
	}{
		{
			name:           "Should return true if a given method is whitelisted for a given app and blockchain",
			portalApp:      testPortalApplication,
			blockchain:     BlockchainID("0001"),
			method:         Method("POST"),
			expectedResult: true,
		},
		{
			name:           "Should return true if a given method is whitelisted for a given app and blockchain",
			portalApp:      testPortalApplication,
			blockchain:     BlockchainID("003E"),
			method:         Method("GET"),
			expectedResult: true,
		},
		{
			name:           "Should return false if a given method is not whitelisted for a given app and blockchain",
			portalApp:      testPortalApplication,
			blockchain:     BlockchainID("0056"),
			method:         Method("PUT"),
			expectedResult: false,
		},
		{
			name:           "Should return false if a given method is not whitelisted for a given app and blockchain",
			portalApp:      testPortalApplication,
			blockchain:     BlockchainID("7009"),
			method:         Method("GET"),
			expectedResult: false,
		},
	}

	for _, test := range tests {
		isBlockchainWhitelisted := test.portalApp.IsMethodWhitelisted(test.blockchain, test.method)
		c.Equal(test.expectedResult, isBlockchainWhitelisted)
	}
}

func Test_PortalApp_GetWhitelistsObject(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name           string
		portalApp      PortalApp
		expectedResult *WhitelistsObject
		expectedJSON   string
	}{
		{
			name:      "Should return an applications whitelists in the form of a WhitelistContracts struct",
			portalApp: testPortalApplication,
			expectedResult: &WhitelistsObject{
				AppWhitelists: [3]ApplicationWhitelists{
					{Type: "origins", Values: []string{"https://portalgun.io", "https://subdomain.example.com", "https://www.example.com"}},
					{Type: "userAgents", Values: []string{"Brave", "Google Chrome", "Mozilla Firefox", "Netscape Navigator", "Safari"}},
					{Type: "blockchains", Values: []string{"0001", "0002", "003E", "0056"}},
				},
				ChainWhitelists: [2]ChainWhitelists{
					{Type: "contracts", Values: []BlockchainIDWhitelists{
						{BlockchainID: "0001", Values: []string{"0xtest_2f78db6436527729929aaf6c616361de0f7", "0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be"}},
						{BlockchainID: "0002", Values: []string{"0xtest_1111117dc0aa78b770fa6a738034120c302", "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2"}},
						{BlockchainID: "003E", Values: []string{"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984", "0xtest_f958d2ee523a2206206994597c13d831ec7"}},
						{BlockchainID: "0056", Values: []string{"0xtest_00000f279d81a1d3cc75430faa017fa5a2e", "0xtest_5068778dd592e39a122f4f5a5cf09c90fe2"}},
					},
					},
					{Type: "methods", Values: []BlockchainIDWhitelists{
						{BlockchainID: "0001", Values: []string{"GET", "POST", "PUT"}},
						{BlockchainID: "0002", Values: []string{"DELETE", "GET", "POST", "PUT"}},
						{BlockchainID: "003E", Values: []string{"GET"}},
						{BlockchainID: "0056", Values: []string{"GET", "POST"}},
					},
					},
				},
			},
			expectedJSON: `{
				"appWhitelists": [
					{
						"type": "origins",
						"values": [
							"https://portalgun.io",
							"https://subdomain.example.com",
							"https://www.example.com"
						]
					},
					{
						"type": "userAgents",
						"values": [
							"Brave",
							"Google Chrome",
							"Mozilla Firefox",
							"Netscape Navigator",
							"Safari"
						]
					},
					{
						"type": "blockchains",
						"values": [
							"0001",
							"0002",
							"003E",
							"0056"
						]
					}
				],
				"chainWhitelists": [
					{
						"type": "contracts",
						"values": [
							{
								"blockchainID": "0001",
								"values": [
									"0xtest_2f78db6436527729929aaf6c616361de0f7",
									"0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be"
								]
							},
							{
								"blockchainID": "0002",
								"values": [
									"0xtest_1111117dc0aa78b770fa6a738034120c302",
									"0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2"
								]
							},
							{
								"blockchainID": "003E",
								"values": [
									"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984",
									"0xtest_f958d2ee523a2206206994597c13d831ec7"
								]
							},
							{
								"blockchainID": "0056",
								"values": [
									"0xtest_00000f279d81a1d3cc75430faa017fa5a2e",
									"0xtest_5068778dd592e39a122f4f5a5cf09c90fe2"
								]
							}
						]
					},
					{
						"type": "methods",
						"values": [
							{
								"blockchainID": "0001",
								"values": [
									"GET",
									"POST",
									"PUT"
								]
							},
							{
								"blockchainID": "0002",
								"values": [
									"DELETE",
									"GET",
									"POST",
									"PUT"
								]
							},
							{
								"blockchainID": "003E",
								"values": [
									"GET"
								]
							},
							{
								"blockchainID": "0056",
								"values": [
									"GET",
									"POST"
								]
							}
						]
					}
				]
			}`,
		},
	}

	for _, test := range tests {
		whitelistsObject := test.portalApp.GetWhitelistsObject()
		c.Equal(test.expectedResult, whitelistsObject)

		// check that JSON is equal as well
		resultJSON, err := json.MarshalIndent(whitelistsObject, "", "  ")
		c.NoError(err)
		expectedJSON := strings.ReplaceAll(strings.ReplaceAll(test.expectedJSON, " ", ""), "\t", "")
		actualJSON := strings.ReplaceAll(strings.ReplaceAll(string(resultJSON), " ", ""), "\t", "")
		c.Equal(expectedJSON, actualJSON)
	}
}
