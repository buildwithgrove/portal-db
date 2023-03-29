package postgresdriver

import (
	"testing"
	"time"

	"github.com/pokt-foundation/portal-db/v2/testdata"
	"github.com/pokt-foundation/portal-db/v2/types"
	"github.com/stretchr/testify/require"
)

func TestListen(t *testing.T) {
	tests := []struct {
		name                  string
		content               types.SavedOnDB
		expectedNotifications map[types.Table]*types.Notification
		wantPanic             bool
	}{
		{
			name:    "Should process Account notifications",
			content: testdata.Accounts[5],
			expectedNotifications: map[types.Table]*types.Notification{
				types.TableAccounts: {
					Table:  types.TableAccounts,
					Action: types.ActionInsert,
					Data: &types.Account{
						ID:                     5,
						Name:                   "test_account_5",
						PlanType:               "basic_plan",
						PartnerChainIDs:        map[types.RelayChainID]struct{}{"0006": {}, "0040": {}},
						PartnerThroughputLimit: 6_000,
						PartnerAppLimit:        1,
						CreatedAt:              testdata.MockTimestamp,
						UpdatedAt:              testdata.MockTimestamp,
						LegacyLoadBalancerID:   "test_lb_f5ee77c7c58025231",
					},
				},
				types.TableAccountUserAccess: {
					Table:  types.TableAccountUserAccess,
					Action: types.ActionUpdate,
					Data: &types.AccountUserAccess{
						AccountID: 5,
						UserID:    4,
						RoleName:  types.RoleOwner,
						Accepted:  true,
					},
				},
			},
		},
		{
			name:    "Should process PortalApp notifications",
			content: testdata.PortalApps["test_app_3487u329rfn23f9"],
			expectedNotifications: map[types.Table]*types.Notification{
				types.TablePortalApps: {
					Table:  types.TablePortalApps,
					Action: types.ActionInsert,
					Data: &types.PortalApp{
						ID:        "test_app_3487u329rfn23f9",
						AccountID: 1,
						Name:      "pokt_app_123",
						Gigastake: true,
						Staked:    false,
						CreatedAt: testdata.MockTimestamp,
						UpdatedAt: testdata.MockTimestamp,
						Deleted:   false,
						LegacyFields: types.LegacyFields{
							CustomLimit:        0,
							RequestTimeout:     5000,
							GigastakeRedirect:  true,
							FirstDateSurpassed: testdata.MockTimestamp,
						},
					},
				},
				types.TableAppAATs: {
					Table:  types.TableAppAATs,
					Action: types.ActionUpdate,
					Data: &types.AAT{
						AppID:           "test_app_3487u329rfn23f9",
						Address:         "test_34715cae753e67c75fbb340442e7de8e",
						PublicKey:       "test_34715cae753e67c75fbb340442e7de8e",
						ClientPublicKey: "test_89a3af6a587aec02cfade6f5000424c2",
						PrivateKey:      "test_11b8d394ca331d7c7a71ca1896d630f6",
						Signature:       "test_1dc39a2e5a84a35bf030969a0b3231f7",
						Version:         "0.0.1",
					},
				},
				types.TableAppSettings: {
					Table:  types.TableAppSettings,
					Action: types.ActionUpdate,
					Data: &types.Settings{
						AppID:             "test_app_3487u329rfn23f9",
						Environment:       "production",
						SecretKey:         "test_40f482d91a5ef2300ebb4e2308c",
						SecretKeyRequired: true,
						MonthlyRelayLimit: 0,
					},
				},
				types.TableAppWhitelists: {
					Table:  types.TableAppWhitelists,
					Action: types.ActionUpdate,
					Data: &types.Whitelist{
						AppID: "test_app_3487u329rfn23f9",
						Type:  "blockchains",
						Value: "0053",
					},
				},
				types.TableAppNotifications: {
					Table:  types.TableAppNotifications,
					Action: types.ActionUpdate,
					Data: &types.AppNotification{
						AppID:       "test_app_3487u329rfn23f9",
						Active:      true,
						Destination: "test@test.com",
						Trigger:     "trigger123",
						Events:      map[types.NotificationEvent]bool{"full": true, "quarter": true, "threeQuarters": true},
					},
				},
				types.TableStickinessOptions: {
					Table:  types.TableStickinessOptions,
					Action: types.ActionUpdate,
					Data: &types.StickyOptions{
						ID:            "test_app_3487u329rfn23f9",
						Duration:      "60",
						StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
						StickyMax:     300,
						Stickiness:    true,
					},
				},
			},
		},
		{
			name:    "Should process Chain notifications",
			content: testdata.Chains["0001"],
			expectedNotifications: map[types.Table]*types.Notification{
				types.TableChains: {
					Table:  types.TableChains,
					Action: types.ActionInsert,
					Data: &types.Chain{
						ID:            "0001",
						Blockchain:    "mainnet",
						Description:   "Pocket Network Mainnet",
						EnforceResult: "JSON",
						Path:          "/v1/query/height",
						Ticker:        "POKT",
						ChainAliases:  []string{"mainnet"},
						Active:        true,
						CreatedAt:     time.Date(2022, time.November, 11, 11, 11, 11, 0, time.UTC),
						UpdatedAt:     time.Date(2022, time.November, 11, 11, 11, 11, 0, time.UTC),
					},
				},
				types.TableChainAltruists: {
					Table:  types.TableChainAltruists,
					Action: types.ActionUpdate,
					Data: &types.Altruist{
						ChainID:  "0001",
						URL:      "https://altruist-0001.com:1234",
						Auth:     "test_pocket:auth123456",
						AuthType: types.ChainAuthTypeBasicAuth,
					},
				},
				types.TableChainGigastakeRedirects: {
					Table:  types.TableChainGigastakeRedirects,
					Action: types.ActionUpdate,
					Data: &types.GigastakeRedirect{
						ChainID:              "0001",
						AccountID:            1,
						Domain:               "pokt-rpc.gateway.pokt.network",
						Alias:                "altruist-0001",
						LegacyLoadBalancerID: "test_lb_5c6f50bc30b530a8",
					},
				},
				types.TableChainChecks: {
					Table:  types.TableChainChecks,
					Action: types.ActionUpdate,
					Data: &types.Check{
						Type:      types.ChainCheckTypeSync,
						ChainID:   "0001",
						Payload:   "{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"query\"}",
						ResultKey: "result.sync_info",
						Allowance: 1,
					},
				},
			},
		},
		{
			name:                  "Should return error for unsupported content",
			content:               &types.User{},
			expectedNotifications: map[types.Table]*types.Notification{},
			wantPanic:             true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != test.wantPanic {
					t.Errorf("recover = %v, wantPanic = %v", r, test.wantPanic)
				}
			}()

			c := require.New(t)

			listenerMock := NewListenerMock()
			driver := NewPostgresDriverFromDBInstance(nil, listenerMock)

			listenerMock.MockEvent(types.ActionInsert, types.ActionUpdate, test.content)

			<-time.After(1 * time.Second)
			driver.CloseListener()

			notificationsMap := make(map[types.Table]*types.Notification)

			for n := range driver.NotificationChannel() {
				notificationsMap[n.Table] = n
			}

			c.Equal(test.expectedNotifications, notificationsMap)
		})
	}
}
