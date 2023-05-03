package postgresdriver

import (
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pokt-foundation/portal-db/v2/types"
)

type (
	// ListenerMock simulates a listener that receives notifications from a PostgreSQL database.
	ListenerMock struct {
		Notify chan *pq.Notification
	}
	inputStruct struct {
		action types.Action
		table  types.Table
		input  any
	}
)

func NewListenerMock() *ListenerMock {
	return &ListenerMock{Notify: make(chan *pq.Notification, 32)}
}

// NotificationChannel returns a channel that receives pq.Notification instances.
func (l *ListenerMock) NotificationChannel() <-chan *pq.Notification {
	return l.Notify
}

func (l *ListenerMock) Listen(channel string) error {
	return nil
}

// MockEvent simulates a database event by sending mock notifications to the ListenerMock's Notify channel.
// To exclude main or side table actions simply pass an empty string as the parameter.
func (l *ListenerMock) MockEvent(mainTableAction, sideTablesAction types.Action, content types.SavedOnDB) {
	notifications := mockContent(mainTableAction, sideTablesAction, content)

	for _, notification := range notifications {
		l.Notify <- notification
	}
}

// mockInput creates a mock pq.Notification based on the inputStruct provided.
func mockInput(inStruct inputStruct) *pq.Notification {
	notification, _ := json.Marshal(notification{
		Table:  inStruct.table,
		Action: inStruct.action,
		Data:   inStruct.input,
	})

	return &pq.Notification{
		Extra: string(notification),
	}
}

// mockContent creates a list of mock notifications for a given mainTableAction, sideTablesAction, and content.
func mockContent(mainTableAction, sideTablesAction types.Action, content types.SavedOnDB) []*pq.Notification {
	var inputs []inputStruct

	switch content.(type) {
	case *types.Account:
		inputs = accountInputs(mainTableAction, sideTablesAction, content)
	case *types.PortalApp:
		inputs = portalAppInputs(mainTableAction, sideTablesAction, content)
	case *types.Chain:
		inputs = chainInputs(mainTableAction, sideTablesAction, content)
	default:
		panic("invalid content type")
	}

	var notifications []*pq.Notification

	for _, input := range inputs {
		notifications = append(notifications, mockInput(input))
	}

	return notifications
}

// accountInputs generates the mock data for a listener notification for an Account struct
func accountInputs(mainTableAction, sideTablesAction types.Action, content types.SavedOnDB) []inputStruct {
	account := content.(*types.Account)

	var inputs []inputStruct

	partnerChainIDs := make([]string, 0, len(account.PartnerChainIDs))
	for chainID := range account.PartnerChainIDs {
		partnerChainIDs = append(partnerChainIDs, string(chainID))
	}

	if mainTableAction != "" {
		inputs = append(inputs, inputStruct{
			action: mainTableAction,
			table:  types.TableAccounts,
			input: dbAccount{
				ID:                      account.ID,
				PlanType:                account.PlanType,
				PartnerChainIDs:         partnerChainIDs,
				PartnerThroughputLimit:  account.PartnerThroughputLimit,
				PartnerApplicationLimit: account.PartnerAppLimit,
				CreatedAt:               account.CreatedAt,
				UpdatedAt:               account.UpdatedAt,
				Deleted:                 account.Deleted,
			},
		})
	}

	if sideTablesAction != "" {
		for _, userAccess := range account.Users {
			inputs = append(inputs, inputStruct{
				action: sideTablesAction,
				table:  types.TableAccountUserAccess,
				input: dbAccountUserAccess{
					AccountID: account.ID,
					UserID:    userAccess.UserID,
					RoleName:  userAccess.RoleName,
					Accepted:  userAccess.Accepted,
				},
			})
		}
	}

	return inputs
}

// portalAppInputs generates the mock data for a listener notification for a PortalApp struct
func portalAppInputs(mainTableAction, sideTablesAction types.Action, content types.SavedOnDB) []inputStruct {
	portalApp := content.(*types.PortalApp)

	var inputs []inputStruct

	if mainTableAction != "" {
		inputs = append(inputs, inputStruct{
			action: mainTableAction,
			table:  types.TablePortalApps,
			input: dbPortalApplication{
				ID:                 portalApp.ID,
				AccountID:          string(portalApp.AccountID),
				Name:               portalApp.Name,
				Gigastake:          portalApp.Gigastake,
				Staked:             portalApp.Staked,
				CreatedAt:          portalApp.CreatedAt,
				UpdatedAt:          portalApp.UpdatedAt,
				Deleted:            portalApp.Deleted,
				RequestTimeout:     portalApp.LegacyFields.RequestTimeout,
				GigastakeRedirect:  portalApp.LegacyFields.GigastakeRedirect,
				FirstDateSurpassed: portalApp.LegacyFields.FirstDateSurpassed,
				CustomLimit:        portalApp.LegacyFields.CustomLimit,
			},
		})
	}

	if sideTablesAction != "" {
		aat := types.AAT{}
		for _, portalAAT := range portalApp.AATs {
			aat = portalAAT
		}

		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  types.TableAppAATs,
			input: dbPortalApplicationAAT{
				ApplicationID:   portalApp.ID,
				ID:              aat.ID,
				Address:         aat.Address,
				PublicKey:       aat.PublicKey,
				ClientPublicKey: aat.ClientPublicKey,
				PrivateKey:      aat.PrivateKey,
				Signature:       aat.Signature,
				Version:         aat.Version,
			},
		})

		favoritedChainIDs := make([]string, 0, len(portalApp.Settings.FavoritedChainIDs))
		for chainID := range portalApp.Settings.FavoritedChainIDs {
			favoritedChainIDs = append(favoritedChainIDs, string(chainID))
		}

		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  types.TableAppSettings,
			input: dbPortalApplicationSetting{
				ApplicationID:     portalApp.ID,
				MonthlyRelayLimit: portalApp.Settings.MonthlyRelayLimit,
				Environment:       portalApp.Settings.Environment,
				FavoritedChainIDs: favoritedChainIDs,
				SecretKey:         portalApp.Settings.SecretKey,
				SecretKeyRequired: portalApp.Settings.SecretKeyRequired,
			},
		})

		for chainID := range portalApp.Whitelists.Blockchains {
			inputs = append(inputs, inputStruct{
				action: sideTablesAction,
				table:  types.TableAppWhitelists,
				input: dbPortalApplicationWhitelist{
					ApplicationID: portalApp.ID,
					Type:          types.WhitelistTypeBlockchains,
					Value:         string(chainID),
				},
			})
		}

		for _, notification := range portalApp.Notifications {
			events := []types.NotificationEvent{}
			for event, active := range notification.Events {
				if active {
					events = append(events, event)
				}
			}

			inputs = append(inputs, inputStruct{
				action: sideTablesAction,
				table:  types.TableAppNotifications,
				input: dbPortalApplicationNotification{
					ApplicationID: portalApp.ID,
					Active:        notification.Active,
					Type:          notification.Type,
					Destination:   notification.Destination,
					Trigger:       notification.Trigger,
					Events:        events,
				},
			})
		}
	}

	return inputs
}

// chainInputs generates the mock data for a listener notification for a Chain struct
func chainInputs(mainTableAction, sideTablesAction types.Action, content types.SavedOnDB) []inputStruct {
	chain := content.(*types.Chain)

	var inputs []inputStruct

	if mainTableAction != "" {
		inputs = append(inputs, inputStruct{
			action: mainTableAction,
			table:  types.TableChains,
			input: dbChain{
				ID:             chain.ID,
				Blockchain:     chain.Blockchain,
				Description:    chain.Description,
				EnforceResult:  chain.EnforceResult,
				Ticker:         chain.Ticker,
				Path:           chain.Path,
				RequestTimeout: chain.RequestTimeout,
				LogLimitBlocks: chain.LogLimitBlocks,
				ChainAliases:   chain.ChainAliases,
				AllowedMethods: chain.AllowedMethods,
				Active:         chain.Active,
				CreatedAt:      chain.CreatedAt,
				UpdatedAt:      chain.UpdatedAt,
			},
		})
	}

	if sideTablesAction != "" {
		for _, altruist := range chain.Altruists {
			inputs = append(inputs, inputStruct{
				action: sideTablesAction,
				table:  types.TableChainAltruists,
				input: dbChainAltruist{
					ChainID:  chain.ID,
					URL:      altruist.URL,
					Auth:     altruist.Auth,
					AuthType: altruist.AuthType,
				},
			})
		}

		for _, redirect := range chain.Redirects {
			inputs = append(inputs, inputStruct{
				action: sideTablesAction,
				table:  types.TableChainGigastakeRedirects,
				input: dbChainGigastakeRedirect{
					ChainID:             chain.ID,
					PortalApplicationID: redirect.PortalApplicationID,
					Alias:               redirect.Alias,
					Domain:              redirect.Domain,
				},
			})
		}

		for _, check := range chain.Checks {
			inputs = append(inputs, inputStruct{
				action: sideTablesAction,
				table:  types.TableChainChecks,
				input: dbChainCheck{
					ChainID:    chain.ID,
					Type:       check.Type,
					Payload:    check.Payload,
					ResultKey:  check.ResultKey,
					Allowance:  check.Allowance,
					EVMChainID: check.EVMChainID,
				},
			})
		}
	}

	return inputs
}
