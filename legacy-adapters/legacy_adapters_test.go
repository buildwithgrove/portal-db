package legacyadapters

import (
	"testing"

	v1Types "github.com/pokt-foundation/portal-db/types"
	"github.com/pokt-foundation/portal-db/v2/testdata"
	v2Types "github.com/pokt-foundation/portal-db/v2/types"
	"github.com/stretchr/testify/require"
)

func Test_LegacyAdapators_ConvertToLegacyLoadBalancer(t *testing.T) {
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
			legacyLoadBalancer := ConvertToLegacyLoadBalancer(test.portalApp, test.account, test.account.Users)
			c.Equal(test.expectedLegacyLoadBalancer, legacyLoadBalancer)
		})
	}
}

func Test_LegacyAdapators_ConvertToLegacyApplication(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                       string
		portalApp                  v2Types.PortalApp
		userID                     string
		planType                   v2Types.PayPlanType
		dailyLimit                 int32
		expectedLegacyApplications []*v1Types.Application
	}{
		{
			name:                       "Should convert a V2 PortalApp struct to a legacy Application struct",
			portalApp:                  *testdata.PortalApps["test_app_1"],
			userID:                     "auth0|james_holden",
			planType:                   "basic_plan",
			dailyLimit:                 1_000,
			expectedLegacyApplications: testdata.LegacyApplications,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			legacyApplications := ConvertToLegacyApplications(test.portalApp, test.userID, test.planType, test.dailyLimit)
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
	}{
		{
			name:                     "Should convert a V2 Chain struct to a legacy Blockchain struct",
			chain:                    *testdata.Chains["0001"],
			expectedLegacyBlockchain: testdata.LegacyBlockchain,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
			plan:                  *testdata.PayPlans["enterprise_plan"],
			expectedLegacyPayPlan: v1Types.PayPlan{Type: "enterprise_plan", Limit: 10_000},
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
	}{
		{
			name:                   "Should convert a legacy LoadBalancer struct to V2 PortalApp & AAT structs",
			loadBalancer:           testdata.LegacyLoadBalancer,
			expectedV2PortalApp:    testdata.V2CreatePortalApp,
			expectedV2PortalAppAAT: testdata.V2CreatePortalAppAAT,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

func Test_LegacyAdapators_ConvertToV2Chain(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name            string
		blockchain      v1Types.Blockchain
		expectedV2Chain v2Types.Chain
	}{
		{
			name:            "Should convert a legacy Blockchain struct to a V2 Chain struct",
			blockchain:      testdata.LegacyBlockchain,
			expectedV2Chain: *testdata.Chains["0001"],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v2Chain := ConvertToV2Chain(test.blockchain)
			v2Chain.Redirects[0].PortalApplicationID = "test_app_1"
			c.Equal(test.expectedV2Chain, v2Chain)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2UpdateChain(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                  string
		blockchainID          string
		updateBlockchain      v1Types.UpdateBlockchain
		expectedV2UpdateChain v2Types.Chain
	}{
		{
			name:                  "Should convert a legacy UpdateBlockchain struct to a V2 Chain struct for use in the update method",
			blockchainID:          "0001",
			updateBlockchain:      testdata.LegacyUpdateBlockchain,
			expectedV2UpdateChain: testdata.UpdateChainTwo,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v2UpdateChain := ConvertToV2UpdateChain(test.blockchainID, test.updateBlockchain)

			// Set fields not used in the legacy update to reuse test struct
			v2UpdateChain.ID = test.expectedV2UpdateChain.ID
			v2UpdateChain.Active = test.expectedV2UpdateChain.Active
			v2UpdateChain.Redirects = test.expectedV2UpdateChain.Redirects
			v2UpdateChain.GigastakeRedirectDomains = test.expectedV2UpdateChain.GigastakeRedirectDomains
			v2UpdateChain.CreatedAt = test.expectedV2UpdateChain.CreatedAt
			v2UpdateChain.UpdatedAt = test.expectedV2UpdateChain.UpdatedAt

			c.Equal(test.expectedV2UpdateChain, v2UpdateChain)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2Redirect(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name               string
		redirect           v1Types.Redirect
		expectedV2Redirect v2Types.GigastakeRedirect
	}{
		{
			name:               "Should convert a legacy Redirect struct to a V2 ChainGigastakesRedirect struct",
			redirect:           testdata.LegacyRedirect,
			expectedV2Redirect: testdata.Chains["0001"].Redirects[0],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v2Redirect := ConvertToV2Redirect(test.redirect)
			v2Redirect.PortalApplicationID = "test_app_1"
			c.Equal(test.expectedV2Redirect, v2Redirect)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2AccountUserAccess(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                        string
		userAccess                  v1Types.UserAccess
		expectedV2AccountUserAccess v2Types.AccountUserAccess
	}{
		{
			name:                        "Should convert a legacy Redirect struct to a V2 ChainGigastakesRedirect struct",
			userAccess:                  testdata.LegacyUserAccess,
			expectedV2AccountUserAccess: testdata.AccountUserAccess[1],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v2AccountUserAccess := ConvertToV2AccountUserAccess(test.userAccess)

			// Set fields not used in the legacy update to reuse test struct
			test.expectedV2AccountUserAccess.UserID = ""
			delete(test.expectedV2AccountUserAccess.ProviderUserIDs, v2Types.AuthTypeAuth0Github)
			test.expectedV2AccountUserAccess.ProviderUserIDs[v2Types.AuthTypeAuth0Username] = test.userAccess.UserID

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
				RoleName:             v2Types.RoleAdmin,
				UserID:               "user_123",
				LegacyLoadBalancerID: "test_lb_3127flsdhfoi323f",
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
