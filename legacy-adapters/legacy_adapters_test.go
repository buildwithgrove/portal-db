package legacyadapters

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/pokt-foundation/portal-db/v2/testdata"
	"github.com/pokt-foundation/portal-db/v2/types"
	"github.com/stretchr/testify/require"
)

func generateTestAccountID() string {
	bytes := make([]byte, 24/2)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func Test_LegacyAdapators_ConvertToLegacyLoadBalancer(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                       string
		account                    types.Account
		expectedLegacyLoadBalancer types.LoadBalancer
	}{
		{
			name:                       "Should convert a V2 Account struct to a legacy LoadBalancer struct",
			account:                    *testdata.V2Account,
			expectedLegacyLoadBalancer: testdata.LegacyLoadBalancer,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			legacyLoadBalancer := ConvertToLegacyLoadBalancer(test.account)
			c.Equal(test.expectedLegacyLoadBalancer, legacyLoadBalancer)
		})
	}
}

func Test_LegacyAdapators_ConvertToLegacyApplication(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                      string
		portalApp                 types.PortalApp
		userID                    string
		planType                  types.PayPlanType
		dailyLimit                int32
		expectedLegacyApplication types.Application
	}{
		{
			name:                      "Should convert a V2 PortalApp struct to a legacy Application struct",
			portalApp:                 *testdata.PortalApps["test_app_3487u329rfn23f9"],
			userID:                    "james_holden",
			planType:                  "basic_plan",
			dailyLimit:                1_000,
			expectedLegacyApplication: testdata.LegacyApplication,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			legacyApplication := ConvertToLegacyApplication(test.portalApp, test.userID, test.planType, test.dailyLimit)
			c.Equal(test.expectedLegacyApplication, legacyApplication)
		})
	}
}

func Test_LegacyAdapators_ConvertToLegacyBlockchain(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                     string
		chain                    types.Chain
		expectedLegacyBlockchain types.Blockchain
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
		plan                  types.Plan
		expectedLegacyPayPlan types.PayPlan
	}{
		{
			name:                  "Should convert a V2 Plan struct to a legacy PayPlan struct",
			plan:                  *testdata.PayPlans["enterprise_plan"],
			expectedLegacyPayPlan: types.PayPlan{Type: "enterprise_plan", Limit: 10_000},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			legacyPayPlan := ConvertToLegacyPayPlan(&test.plan)
			c.Equal(test.expectedLegacyPayPlan, legacyPayPlan)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2AccountAndPortalApp(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                string
		loadBalancer        types.LoadBalancer
		accountID           types.AccountID
		expectedV2Account   types.Account
		expectedV2PortalApp types.PortalApp
	}{
		{
			name:                "Should convert a legacy LoadBalancer struct to V2 Account & PortalApp structs",
			loadBalancer:        testdata.LegacyLoadBalancer,
			accountID:           1,
			expectedV2Account:   testdata.V2CreateAccount,
			expectedV2PortalApp: testdata.V2CreatePortalApp,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testID := generateTestAccountID()
			test.expectedV2Account.LegacyLoadBalancerID = testID

			v2Account, v2PortalApp := ConvertToV2AccountAndPortalApp(test.loadBalancer, testID)
			c.Equal(test.expectedV2Account, v2Account)
			c.Equal(test.expectedV2PortalApp, v2PortalApp)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2UpdatePortalApp(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                      string
		appID                     string
		updateApplication         types.UpdateApplication
		expectedV2UpdatePortalApp types.UpdatePortalApp
	}{
		{
			name:                      "Should convert a legacy UpdateApplication struct to a V2 UpdatePortalApp struct",
			appID:                     "test_app_3487u329rfn23f9",
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
		blockchain      types.Blockchain
		expectedV2Chain types.Chain
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
			v2Chain.Redirects[0].AccountID = 1
			c.Equal(test.expectedV2Chain, v2Chain)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2UpdateChain(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                  string
		updateBlockchain      types.UpdateBlockchain
		expectedV2UpdateChain types.Chain
	}{
		{
			name:                  "Should convert a legacy UpdateBlockchain struct to a V2 Chain struct for use in the update method",
			updateBlockchain:      testdata.LegacyUpdateBlockchain,
			expectedV2UpdateChain: testdata.UpdateChainTwo,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v2UpdateChain := ConvertToV2UpdateChain(test.updateBlockchain)

			// Set fields not used in the legacy update to reuse test struct
			v2UpdateChain.ID = test.expectedV2UpdateChain.ID
			v2UpdateChain.Active = test.expectedV2UpdateChain.Active
			v2UpdateChain.Redirects = test.expectedV2UpdateChain.Redirects
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
		redirect           types.Redirect
		expectedV2Redirect types.GigastakeRedirect
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
			v2Redirect.AccountID = 1
			c.Equal(test.expectedV2Redirect, v2Redirect)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2AccountUserAccess(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                        string
		userAccess                  types.UserAccess
		expectedV2AccountUserAccess types.AccountUserAccess
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
			test.expectedV2AccountUserAccess.UserID = 0
			delete(test.expectedV2AccountUserAccess.ProviderUserIDs, types.AuthTypeAuth0Github)
			test.expectedV2AccountUserAccess.ProviderUserIDs[types.AuthTypeAuth0Username] = test.userAccess.UserID

			c.Equal(test.expectedV2AccountUserAccess, v2AccountUserAccess)
		})
	}
}

func Test_LegacyAdapators_ConvertToV2UpdateAccountUserAccess(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                        string
		lbID                        string
		userID                      int32
		userAccess                  types.UpdateUserAccess
		expectedV2AccountUserAccess types.UpdateAccountUserRole
	}{
		{
			name:       "Should convert a legacy Redirect struct to a V2 ChainGigastakesRedirect struct",
			lbID:       "test_lb_3127flsdhfoi323f",
			userID:     123,
			userAccess: testdata.LegacyUpdateUserAccess,
			expectedV2AccountUserAccess: types.UpdateAccountUserRole{
				RoleName:             types.RoleAdmin,
				UserID:               123,
				LegacyLoadBalancerID: "test_lb_3127flsdhfoi323f",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v2AccountUserAccess := ConvertToV2UpdateAccountUserAccess(test.userAccess, test.lbID, test.userID)
			test.expectedV2AccountUserAccess.UserID = types.UserID(test.userID)
			c.Equal(test.expectedV2AccountUserAccess, v2AccountUserAccess)
		})
	}
}
