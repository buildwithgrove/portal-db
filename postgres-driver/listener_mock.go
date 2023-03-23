package postgresdriver

import (
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pokt-foundation/portal-db/v2/types"
)

type ListenerMock struct {
	Notify chan *pq.Notification
}

func NewListenerMock() *ListenerMock {
	return &ListenerMock{
		Notify: make(chan *pq.Notification, 32),
	}
}

func (l *ListenerMock) NotificationChannel() <-chan *pq.Notification {
	return l.Notify
}

func (l *ListenerMock) Listen(channel string) error {
	return nil
}

type inputStruct struct {
	action types.Action
	table  types.Table
	input  any
}

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

func mockContent(mainTableAction, sideTablesAction types.Action, content types.SavedOnDB) []*pq.Notification {
	var inputs []inputStruct

	switch content.(type) {
	case *types.PortalApp:
		inputs = portalAppInputs(mainTableAction, sideTablesAction, content)
	// case *types.Chain:
	// 	inputs = blockchainInputs(mainTableAction, sideTablesAction, content)
	// case *types.Account:
	// 	inputs = loadBalancerInputs(mainTableAction, sideTablesAction, content)
	default:
		panic("type not supported")
	}

	var notifications []*pq.Notification

	for _, input := range inputs {
		notifications = append(notifications, mockInput(input))
	}

	return notifications
}

func (l *ListenerMock) MockEvent(mainTableAction, sideTablesAction types.Action, content types.SavedOnDB) {
	notifications := mockContent(mainTableAction, sideTablesAction, content)

	for _, notification := range notifications {
		l.Notify <- notification
	}
}

func portalAppInputs(mainTableAction, sideTablesAction types.Action, content types.SavedOnDB) []inputStruct {
	portalApp := content.(*types.PortalApp)

	var inputs []inputStruct

	inputs = append(inputs, inputStruct{
		action: mainTableAction,
		table:  types.TablePortalApps,
		input: PortalApplication{
			ID:                 portalApp.ID,
			AccountID:          int32(portalApp.AccountID),
			Name:               portalApp.Name,
			Gigastake:          portalApp.Gigastake,
			Staked:             portalApp.Staked,
			CreatedAt:          portalApp.CreatedAt,
			UpdatedAt:          portalApp.UpdatedAt,
			Deleted:            portalApp.Deleted,
			RequestTimeout:     newSQLNullInt32(portalApp.LegacyFields.RequestTimeout, true),
			GigastakeRedirect:  newSQLNullBool(&portalApp.LegacyFields.GigastakeRedirect),
			FirstDateSurpassed: newSQLNullTime(portalApp.LegacyFields.FirstDateSurpassed),
			CustomLimit:        newSQLNullInt32(portalApp.LegacyFields.CustomLimit, true),
		},
	})

	inputs = append(inputs, inputStruct{
		action: sideTablesAction,
		table:  types.TableAppAATs,
		input: PortalApplicationAat{
			ApplicationID:   portalApp.ID,
			Address:         portalApp.AAT.Address,
			PublicKey:       portalApp.AAT.PublicKey,
			ClientPublicKey: portalApp.AAT.ClientPublicKey,
			PrivateKey:      portalApp.AAT.PrivateKey,
			Signature:       portalApp.AAT.Signature,
			Version:         portalApp.AAT.Version,
		},
	})

	favoritedChainIDs := make([]string, 0, len(portalApp.Settings.FavoritedChainIDs))
	for chainID := range portalApp.Settings.FavoritedChainIDs {
		favoritedChainIDs = append(favoritedChainIDs, string(chainID))
	}

	inputs = append(inputs, inputStruct{
		action: sideTablesAction,
		table:  types.TableAppSettings,
		input: PortalApplicationSetting{
			ApplicationID:     portalApp.ID,
			MonthlyRelayLimit: portalApp.Settings.MonthlyRelayLimit,
			Environment:       portalApp.Settings.Environment,
			FavoritedChainIDs: favoritedChainIDs,
			SecretKey:         newSQLNullString(portalApp.Settings.SecretKey),
			SecretKeyRequired: newSQLNullBool(&portalApp.Settings.SecretKeyRequired),
		},
	})

	for origin := range portalApp.Whitelists.Origins {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  types.TableAppWhitelists,
			input: PortalApplicationWhitelist{
				ApplicationID: portalApp.ID,
				Type:          types.WhitelistTypeOrigins,
				Value:         string(origin),
			},
		})
	}

	for userAgent := range portalApp.Whitelists.UserAgents {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  types.TableAppWhitelists,
			input: PortalApplicationWhitelist{
				ApplicationID: portalApp.ID,
				Type:          types.WhitelistTypeUserAgents,
				Value:         string(userAgent),
			},
		})
	}

	for chainID := range portalApp.Whitelists.Blockchains {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  types.TableAppWhitelists,
			input: PortalApplicationWhitelist{
				ApplicationID: portalApp.ID,
				Type:          types.WhitelistTypeBlockchains,
				Value:         string(chainID),
			},
		})
	}

	for chainID, contracts := range portalApp.Whitelists.Contracts {
		for contract := range contracts {
			inputs = append(inputs, inputStruct{
				action: sideTablesAction,
				table:  types.TableAppWhitelists,
				input: PortalApplicationWhitelist{
					ApplicationID: portalApp.ID,
					Type:          types.WhitelistTypeContracts,
					Value:         string(contract),
					ChainID:       newSQLNullString(string(chainID)),
				},
			})
		}
	}

	for chainID, methods := range portalApp.Whitelists.Methods {
		for method := range methods {
			inputs = append(inputs, inputStruct{
				action: sideTablesAction,
				table:  types.TableAppWhitelists,
				input: PortalApplicationWhitelist{
					ApplicationID: portalApp.ID,
					Type:          types.WhitelistTypeMethods,
					Value:         string(method),
					ChainID:       newSQLNullString(string(chainID)),
				},
			})
		}
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
			input: PortalApplicationNotification{
				ApplicationID: portalApp.ID,
				Active:        notification.Active,
				Type:          notification.Type,
				Destination:   newSQLNullString(notification.Destination),
				Trigger:       newSQLNullString(notification.Trigger),
				Events:        events,
			},
		})
	}

	return inputs
}
