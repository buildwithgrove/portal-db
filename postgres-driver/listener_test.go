package postgresdriver

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pokt-foundation/portal-db/types"
)

func TestListen(t *testing.T) {
	testCases := []struct {
		name                  string
		content               types.SavedOnDB
		expectedNotifications map[types.Table]*types.Notification
		wantPanic             bool
	}{
		{
			name: "application",
			content: &types.Application{
				ID: "321",
				GatewayAAT: types.GatewayAAT{
					Address: "123",
				},
				GatewaySettings: types.GatewaySettings{
					SecretKey:            "123",
					WhitelistBlockchains: []string{"test-chain-1", "test-chain-2"},
					WhitelistContracts: []types.WhitelistContracts{
						{BlockchainID: "001", Contracts: []string{"test123"}},
					},
					WhitelistMethods: []types.WhitelistMethods{
						{BlockchainID: "001", Methods: []string{"POST"}},
					},
				},
				Limit: types.AppLimit{
					PayPlan:     types.PayPlan{Type: types.Enterprise},
					CustomLimit: 2000000,
				},
				NotificationSettings: types.NotificationSettings{
					Full: true,
				},
			},
			expectedNotifications: map[types.Table]*types.Notification{
				types.TableApplications: {
					Table:  types.TableApplications,
					Action: types.ActionInsert,
					Data: &types.Application{
						ID: "321",
					},
				},
				types.TableGatewayAAT: {
					Table:  types.TableGatewayAAT,
					Action: types.ActionUpdate,
					Data: &types.GatewayAAT{
						ID:      "321",
						Address: "123",
					},
				},
				types.TableAppLimits: {
					Table:  types.TableAppLimits,
					Action: types.ActionUpdate,
					Data: &types.AppLimit{
						ID:          "321",
						PayPlan:     types.PayPlan{Type: types.Enterprise, Limit: 0},
						CustomLimit: 2000000,
					},
				},
				types.TableGatewaySettings: {
					Table:  types.TableGatewaySettings,
					Action: types.ActionUpdate,
					Data: &types.GatewaySettings{
						ID:                   "321",
						SecretKey:            "123",
						WhitelistBlockchains: []string{"test-chain-1", "test-chain-2"},
					},
				},
				types.TableWhitelistContracts: {
					Table:  types.TableWhitelistContracts,
					Action: types.ActionUpdate,
					Data: &types.WhitelistContract{
						ID: "321", BlockchainID: "001", Contract: "test123",
					},
				},
				types.TableWhitelistMethods: {
					Table:  types.TableWhitelistMethods,
					Action: types.ActionUpdate,
					Data: &types.WhitelistMethod{
						ID: "321", BlockchainID: "001", Method: "POST",
					},
				},
				types.TableNotificationSettings: {
					Table:  types.TableNotificationSettings,
					Action: types.ActionUpdate,
					Data: &types.NotificationSettings{
						ID:   "321",
						Full: true,
					},
				},
			},
		},
		{
			name: "blockchain",
			content: &types.Blockchain{
				ID: "0021",
				SyncCheckOptions: types.SyncCheckOptions{
					BlockchainID: "0021",
					Body:         "yeh",
				},
			},
			expectedNotifications: map[types.Table]*types.Notification{
				types.TableBlockchains: {
					Table:  types.TableBlockchains,
					Action: types.ActionInsert,
					Data: &types.Blockchain{
						ID: "0021",
					},
				},
				types.TableSyncCheckOptions: {
					Table:  types.TableSyncCheckOptions,
					Action: types.ActionUpdate,
					Data: &types.SyncCheckOptions{
						BlockchainID: "0021",
						Body:         "yeh",
					},
				},
			},
		},
		{
			name: "load balancer",
			content: &types.LoadBalancer{
				ID: "123",
				StickyOptions: types.StickyOptions{
					StickyOrigins: []string{"oahu"},
					Stickiness:    true,
				},
				ApplicationIDs: []string{"a123"},
				Users: []types.UserAccess{
					{RoleName: "ADMIN", UserID: "test_user_admin1234", Email: "admin1@test.com", Accepted: true},
				},
			},
			expectedNotifications: map[types.Table]*types.Notification{
				types.TableLoadBalancers: {
					Table:  types.TableLoadBalancers,
					Action: types.ActionInsert,
					Data: &types.LoadBalancer{
						ID: "123",
					},
				},
				types.TableStickinessOptions: {
					Table:  types.TableStickinessOptions,
					Action: types.ActionUpdate,
					Data: &types.StickyOptions{
						ID:            "123",
						StickyOrigins: []string{"oahu"},
						Stickiness:    true,
					},
				},
				types.TableUserAccess: {
					Table:  types.TableUserAccess,
					Action: types.ActionUpdate,
					Data: &types.UserAccess{
						ID:       "123",
						RoleName: "ADMIN",
						UserID:   "test_user_admin1234",
						Email:    "admin1@test.com",
						Accepted: true,
					},
				},
				types.TableLbApps: {
					Table:  types.TableLbApps,
					Action: types.ActionUpdate,
					Data: &types.LbApp{
						LbID:  "123",
						AppID: "a123",
					},
				},
			},
		},
		{
			name: "redirect",
			content: &types.Redirect{
				BlockchainID: "0021",
			},
			expectedNotifications: map[types.Table]*types.Notification{
				types.TableRedirects: {
					Table:  types.TableRedirects,
					Action: types.ActionInsert,
					Data: &types.Redirect{
						BlockchainID: "0021",
					},
				},
			},
		},
		{
			name:      "panic",
			content:   &types.GatewayAAT{},
			wantPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tc.wantPanic {
					t.Errorf("recover = %v, wantPanic = %v", r, tc.wantPanic)
				}
			}()

			listenerMock := NewListenerMock()
			driver := NewPostgresDriverFromDBInstance(nil, listenerMock)

			listenerMock.MockEvent(types.ActionInsert, types.ActionUpdate, tc.content)

			time.Sleep(1 * time.Second)
			driver.CloseListener()

			nMap := make(map[types.Table]*types.Notification)

			for n := range driver.NotificationChannel() {
				nMap[n.Table] = n
			}

			if diff := cmp.Diff(tc.expectedNotifications, nMap); diff != "" {
				t.Errorf("unexpected value (-want +got):\n%s", diff)
			}
		})
	}
}
