package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	testLegacyLoadBalancer = LoadBalancer{
		ID:                "test_5416bb8d696386455b8",
		Name:              "test_portal_app_123",
		UserID:            "user_id_123",
		RequestTimeout:    5_000,
		Gigastake:         true,
		GigastakeRedirect: true,
		StickyOptions: StickyOptions{
			Duration:      "4_000",
			StickyOrigins: []string{"origin123"},
			StickyMax:     4_000,
			Stickiness:    true,
		},
		Applications: []*Application{&testLegacyApplication},
		Users: []UserAccess{
			{UserID: "user_id_123", RoleName: RoleOwner, Email: "test_owner@user.com", Accepted: true},
			{UserID: "user_id_456", RoleName: RoleMember, Email: "test_member@user.com", Accepted: true},
		},
		CreatedAt: time.Date(2023, time.February, 14, 11, 11, 11, 0, time.UTC),
		UpdatedAt: time.Date(2023, time.February, 27, 13, 13, 13, 0, time.UTC),
	}

	testLegacyApplication = Application{
		ID:                 "test_475jf893f9j2f30jd230e",
		UserID:             "user_id_123",
		Name:               "test_portal_app_123",
		FirstDateSurpassed: time.Date(2023, time.February, 28, 15, 15, 15, 0, time.UTC),
		GatewayAAT: GatewayAAT{
			Address:              "test_34715cae753e67c75fbb340442e7de8e",
			ApplicationPublicKey: "test_11b8d394ca331d7c7a71ca1896d630f6",
			ApplicationSignature: "test_1dc39a2e5a84a35bf030969a0b3231f7",
			ClientPublicKey:      "test_9e9ca4fe13725d412003f4bc518f6974",
			PrivateKey:           "test_89a3af6a587aec02cfade6f5000424c2",
			Version:              "0.0.1",
		},
		GatewaySettings: GatewaySettings{
			SecretKey:            "test_90210ac4bdd3423e24877d1ff92",
			SecretKeyRequired:    true,
			WhitelistOrigins:     []string{"https://portalgun.io", "https://subdomain.example.com", "https://www.example.com"},
			WhitelistBlockchains: []string{"0001", "0002", "003E", "0056"},
			WhitelistUserAgents:  []string{"Brave", "Google Chrome", "Mozilla Firefox", "Netscape Navigator", "Safari"},
			WhitelistContracts: []WhitelistContracts{
				{BlockchainID: "0001", Contracts: []string{"0xtest_2f78db6436527729929aaf6c616361de0f7", "0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be"}},
				{BlockchainID: "0002", Contracts: []string{"0xtest_1111117dc0aa78b770fa6a738034120c302", "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2"}},
				{BlockchainID: "003E", Contracts: []string{"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984", "0xtest_f958d2ee523a2206206994597c13d831ec7"}},
				{BlockchainID: "0056", Contracts: []string{"0xtest_00000f279d81a1d3cc75430faa017fa5a2e", "0xtest_5068778dd592e39a122f4f5a5cf09c90fe2"}},
			},
			WhitelistMethods: []WhitelistMethods{
				{BlockchainID: "0001", Methods: []string{"GET", "POST", "PUT"}},
				{BlockchainID: "0002", Methods: []string{"DELETE", "GET", "POST", "PUT"}},
				{BlockchainID: "003E", Methods: []string{"GET"}},
				{BlockchainID: "0056", Methods: []string{"GET", "POST"}},
			},
		},
		Limit: AppLimit{
			PayPlan:     LegacyPayPlan{Type: FreetierV0, Limit: 250_000},
			CustomLimit: 0,
		},
		NotificationSettings: NotificationSettings{SignedUp: true, Quarter: true, Half: false, ThreeQuarters: true, Full: false},
		CreatedAt:            time.Date(2023, time.February, 14, 11, 11, 11, 0, time.UTC),
		UpdatedAt:            time.Date(2023, time.February, 27, 13, 13, 13, 0, time.UTC),
	}

	testLegacyBlockchain = Blockchain{
		ID:                "0001",
		Altruist:          "https://user:test_123@pokt-test.us-1.pokt.network:1234",
		Blockchain:        "pokt-mainnet",
		Description:       "POKT Network Mainnet",
		EnforceResult:     "JSON",
		Path:              "/wow/test",
		Ticker:            "POKT",
		BlockchainAliases: []string{"pokt-mainnet"},
		LogLimitBlocks:    100_000,
		Active:            true,
		Redirects: []Redirect{
			{
				Alias:          "mainnet",
				Domain:         "pokt.test.com",
				LoadBalancerID: "test_5416bb8d696386455b8",
			},
		},
		SyncCheckOptions: SyncCheckOptions{
			Body:      `{}`,
			ResultKey: "testing",
			Allowance: 1,
		},
		CreatedAt: time.Date(2023, time.February, 14, 11, 11, 11, 0, time.UTC),
		UpdatedAt: time.Date(2023, time.February, 27, 13, 13, 13, 0, time.UTC),
	}

	testLegacyUpdateBlockchain = UpdateBlockchain{
		Altruist:          "https://user:test_123@pokt-test.us-1.pokt.network:1234",
		Blockchain:        "pokt-mainnet",
		Description:       "POKT Network Mainnet",
		EnforceResult:     "JSON",
		Path:              "/wow/test",
		Ticker:            "POKT",
		BlockchainAliases: []string{"pokt-mainnet"},
		LogLimitBlocks:    100_000,
		Body:              `{}`,
		ResultKey:         "testing",
		Allowance:         nil,
	}

	testLegacyUpdateRedirect = Redirect{
		Alias:          "mainnet",
		Domain:         "pokt.test.com",
		LoadBalancerID: "test_5416bb8d696386455b8",
	}

	testV2Chain = Chain{
		ID:                "0001",
		Blockchain:        "pokt-mainnet",
		Description:       "POKT Network Mainnet",
		EnforceResult:     "JSON",
		Path:              "/wow/test",
		Ticker:            "POKT",
		BlockchainAliases: []string{"pokt-mainnet"},
		LogLimitBlocks:    100_000,
		Active:            true,
		Altruists:         []ChainAltruist{{URL: "https://user:test_123@pokt-test.us-1.pokt.network:1234"}},
		SyncCheckOptions: ChainSyncCheckOptions{
			Body:      "{}",
			ResultKey: "testing",
			Allowance: 1,
		},
	}

	testV2UpdateChain = UpdateChain{
		Blockchain:        "pokt-mainnet",
		Description:       "POKT Network Mainnet",
		EnforceResult:     "JSON",
		Path:              "/wow/test",
		Ticker:            "POKT",
		BlockchainAliases: []string{"pokt-mainnet"},
		LogLimitBlocks:    100_000,
		Altruists:         []ChainAltruist{{URL: "https://user:test_123@pokt-test.us-1.pokt.network:1234"}},
		Body:              "{}",
		ResultKey:         "testing",
		Allowance:         nil,
	}
)

func Test_LegacyAdapators_ConvertToLegacyLoadBalancer(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                       string
		portalApp                  PortalApp
		expectedLegacyLoadBalancer LoadBalancer
	}{
		{
			name:                       "Should convert a V2 PortalApp struct to a legacy LoadBalancer struct",
			portalApp:                  testPortalApplication,
			expectedLegacyLoadBalancer: testLegacyLoadBalancer,
		},
	}

	for _, test := range tests {
		legacyLoadBalancer := test.portalApp.ConvertToLegacyLoadBalancer()
		c.Equal(test.expectedLegacyLoadBalancer, legacyLoadBalancer)
	}
}

func Test_LegacyAdapators_ConvertToLegacyApplication(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                      string
		portalApp                 PortalApp
		expectedLegacyApplication Application
	}{
		{
			name:                      "Should convert a V2 PortalApp struct to a legacy Application struct",
			portalApp:                 testPortalApplication,
			expectedLegacyApplication: testLegacyApplication,
		},
	}

	for _, test := range tests {
		legacyApplication := test.portalApp.ConvertToLegacyApplication()
		c.Equal(test.expectedLegacyApplication, legacyApplication)
	}
}

func Test_LegacyAdapators_ConvertToLegacyBlockchain(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                     string
		chain                    Chain
		expectedLegacyBlockchain Blockchain
	}{
		{
			name:                     "Should convert a V2 Chain struct to a legacy Blockchain struct",
			chain:                    testChain,
			expectedLegacyBlockchain: testLegacyBlockchain,
		},
	}

	for _, test := range tests {
		legacyBlockchain := test.chain.ConvertToLegacyBlockchain()
		c.Equal(test.expectedLegacyBlockchain, legacyBlockchain)
	}
}

func Test_LegacyAdapators_ConvertToV2Chain(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name            string
		blockchain      Blockchain
		expectedV2Chain Chain
	}{
		{
			name:            "Should convert a legacy Blockchain struct to a V2 Chain struct",
			blockchain:      testLegacyBlockchain,
			expectedV2Chain: testV2Chain,
		},
	}

	for _, test := range tests {
		v2Chain := test.blockchain.ConvertToV2Chain()
		c.Equal(test.expectedV2Chain, v2Chain)
	}
}

func Test_LegacyAdapators_ConvertToV2UpdateChain(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                 string
		updateBlockchain     UpdateBlockchain
		expectedV2Blockchain UpdateChain
	}{
		{
			name:                 "Should convert a legacy UpdateBlockchain struct to a V2 UpdateChain struct",
			updateBlockchain:     testLegacyUpdateBlockchain,
			expectedV2Blockchain: testV2UpdateChain,
		},
	}

	for _, test := range tests {
		v2UpdateChain := test.updateBlockchain.ConvertToV2UpdateChain()
		c.Equal(test.expectedV2Blockchain, v2UpdateChain)
	}
}

func Test_LegacyAdapators_ConvertToV2Redirect(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name               string
		redirect           Redirect
		expectedV2Redirect ChainGigastakesRedirect
	}{
		{
			name:               "Should convert a legacy Redirect struct to a V2 ChainGigastakesRedirect struct",
			redirect:           testLegacyUpdateRedirect,
			expectedV2Redirect: testRedirect,
		},
	}

	for _, test := range tests {
		v2Redirect := test.redirect.ConvertToV2Redirect()
		c.Equal(test.expectedV2Redirect, v2Redirect)
	}
}
