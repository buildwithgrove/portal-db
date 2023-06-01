package legacyadapters

import (
	"testing"

	v1Types "github.com/pokt-foundation/portal-db/types"
	"github.com/pokt-foundation/portal-db/v2/testdata"
	v2Types "github.com/pokt-foundation/portal-db/v2/types"
	"github.com/stretchr/testify/require"
)

func Test_LegacyAdapators_ConvertPortalAppToLegacyLoadBalancer(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                       string
		portalApp                  v2Types.PortalApp
		account                    v2Types.Account
		expectedLegacyLoadBalancer v1Types.LoadBalancer
	}{
		{
			name:                       "Should convert a V2 PortalApp & Account struct to a legacy LoadBalancer struct",
			portalApp:                  *testdata.PortalApps["test_app_1"],
			account:                    *testdata.V2Account,
			expectedLegacyLoadBalancer: testdata.LegacyLoadBalancer,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			legacyLoadBalancer := ConvertPortalAppToLegacyLoadBalancer(test.portalApp, test.account)
			c.Equal(test.expectedLegacyLoadBalancer, legacyLoadBalancer)
		})
	}
}

func Test_LegacyAdapators_ConvertPortalAppToLegacyApplication(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                       string
		portalApp                  v2Types.PortalApp
		userID                     string
		expectedLegacyApplications []*v1Types.Application
	}{
		{
			name:                       "Should convert a V2 PortalApp struct to a legacy Application struct",
			portalApp:                  *testdata.PortalApps["test_app_1"],
			userID:                     "user_1",
			expectedLegacyApplications: testdata.LegacyApplications,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			legacyApplications := ConvertPortalAppToLegacyApplications(test.portalApp, test.userID)
			c.Equal(test.expectedLegacyApplications, legacyApplications)
		})
	}
}

func Test_LegacyAdapators_ConvertToLegacyBlockchain(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                     string
		chain                    v2Types.Chain
		expectedLegacyBlockchain v1Types.Blockchain
		gigastakeApps            map[v2Types.ProtocolAppID]*v2Types.GigastakeApp
	}{
		{
			name:                     "Should convert a V2 Chain struct to a legacy Blockchain struct",
			chain:                    *testdata.Chains["0001"],
			expectedLegacyBlockchain: testdata.LegacyBlockchain,
			gigastakeApps: map[v2Types.ProtocolAppID]*v2Types.GigastakeApp{
				"test_gigastake_app_1": testdata.GigastakeApps["test_gigastake_app_1"],
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Gigastake apps will be set in PHD cache
			test.chain.GigastakeApps = test.gigastakeApps

			legacyBlockchain := ConvertToLegacyBlockchain(test.chain)
			c.Equal(test.expectedLegacyBlockchain, legacyBlockchain)
		})
	}
}

func Test_LegacyAdapators_ConvertToLegacyPayPlan(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                  string
		plan                  v2Types.Plan
		expectedLegacyPayPlan v1Types.PayPlan
	}{
		{
			name:                  "Should convert a V2 Plan struct to a legacy PayPlan struct",
			plan:                  *testdata.PayPlans["ENTERPRISE"],
			expectedLegacyPayPlan: v1Types.PayPlan{Type: "ENTERPRISE", Limit: 0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			legacyPayPlan := ConvertToLegacyPayPlan(test.plan)
			c.Equal(test.expectedLegacyPayPlan, legacyPayPlan)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2PortalAppAndAAT(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                   string
		loadBalancer           v1Types.LoadBalancer
		expectedV2PortalApp    v2Types.PortalApp
		expectedV2PortalAppAAT v2Types.AAT
		testPrivateKeyInput    string
	}{
		{
			name:                   "Should convert a legacy LoadBalancer struct to V2 PortalApp & AAT structs",
			loadBalancer:           testdata.LegacyLoadBalancer,
			expectedV2PortalApp:    testdata.V2CreatePortalApp,
			expectedV2PortalAppAAT: testdata.V2CreatePortalAppAAT,
			testPrivateKeyInput:    "test_11b8d394ca331d7c7a71ca1896d630f6",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.loadBalancer.Applications[0].GatewayAAT.PrivateKey = test.testPrivateKeyInput
			v2PortalApp, v2AAT := ConvertToV2PortalAppAndAAT(test.loadBalancer)
			c.Equal(test.expectedV2PortalApp, v2PortalApp)
			c.Equal(test.expectedV2PortalAppAAT, v2AAT)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2UpdatePortalApp(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                      string
		appID                     string
		updateApplication         v1Types.UpdateApplication
		expectedV2UpdatePortalApp v2Types.UpdatePortalApp
	}{
		{
			name:                      "Should convert a legacy UpdateApplication struct to a V2 UpdatePortalApp struct",
			appID:                     "test_app_1",
			updateApplication:         testdata.LegacyUpdateApplication,
			expectedV2UpdatePortalApp: testdata.V2UpdatePortalApp,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v2UpdatePortalApp := ConvertToV2UpdatePortalApp(test.updateApplication, test.appID)
			c.Equal(test.expectedV2UpdatePortalApp, v2UpdatePortalApp)
			c.Equal(test.expectedV2UpdatePortalApp.Settings, v2UpdatePortalApp.Settings)
			c.Equal(test.expectedV2UpdatePortalApp.Whitelists, v2UpdatePortalApp.Whitelists)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2AccountUserAccess(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                        string
		userAccess                  v1Types.UserAccess
		lbID, accountID             string
		expectedV2AccountUserAccess v2Types.AccountUserAccess
	}{
		{
			name:                        "Should convert a legacy Redirect struct to a V2 ChainGigastakesRedirect struct",
			userAccess:                  testdata.LegacyUserAccess,
			accountID:                   "account_1",
			lbID:                        "test_app_1",
			expectedV2AccountUserAccess: testdata.AccountUserAccess[1],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v2AccountUserAccess := ConvertToV2AccountUserAccess(test.userAccess, test.lbID, test.accountID)

			test.expectedV2AccountUserAccess.AccountID = v2Types.AccountID(test.accountID)

			c.Equal(test.expectedV2AccountUserAccess, v2AccountUserAccess)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2UpdateAccountUserAccess(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                        string
		lbID                        string
		userID                      string
		userAccess                  v1Types.UpdateUserAccess
		expectedV2AccountUserAccess v2Types.UpdateAccountUserRole
	}{
		{
			name:       "Should convert a legacy Redirect struct to a V2 ChainGigastakesRedirect struct",
			lbID:       "test_lb_3127flsdhfoi323f",
			userID:     "user_123",
			userAccess: testdata.LegacyUpdateUserAccess,
			expectedV2AccountUserAccess: v2Types.UpdateAccountUserRole{
				RoleName:    v2Types.RoleAdmin,
				UserID:      "user_123",
				PortalAppID: "test_lb_3127flsdhfoi323f",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v2AccountUserAccess := ConvertToV2UpdateAccountUserAccess(test.userAccess, test.lbID, test.userID)
			test.expectedV2AccountUserAccess.UserID = v2Types.UserID(test.userID)
			c.Equal(test.expectedV2AccountUserAccess, v2AccountUserAccess)
		})
	}
}
