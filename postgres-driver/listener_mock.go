package postgresdriver

import (
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pokt-foundation/portal-db/types"
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

func gatewaySettingsIsNull(settings types.GatewaySettings) bool {
	return settings.SecretKey == "" &&
		len(settings.WhitelistOrigins) == 0 &&
		len(settings.WhitelistUserAgents) == 0 &&
		len(settings.WhitelistContracts) == 0 &&
		len(settings.WhitelistMethods) == 0 &&
		len(settings.WhitelistBlockchains) == 0
}

func applicationInputs(mainTableAction, sideTablesAction types.Action, content types.SavedOnDB) []inputStruct {
	app := content.(*types.Application)

	var inputs []inputStruct

	inputs = append(inputs, inputStruct{
		action: mainTableAction,
		table:  types.TableApplications,
		input: dbAppJSON{
			ApplicationID:      app.ID,
			UserID:             app.UserID,
			Name:               app.Name,
			ContactEmail:       app.ContactEmail,
			Description:        app.Description,
			Owner:              app.Owner,
			URL:                app.URL,
			Status:             string(app.Status),
			CreatedAt:          app.CreatedAt.Format(psqlDateLayout),
			UpdatedAt:          app.UpdatedAt.Format(psqlDateLayout),
			FirstDateSurpassed: app.FirstDateSurpassed.Format(psqlDateLayout),
			Dummy:              app.Dummy,
		},
	})

	inputs = append(inputs, inputStruct{
		action: sideTablesAction,
		table:  types.TableAppLimits,
		input: dbAppLimitJSON{
			ApplicationID: app.ID,
			PlanType:      app.Limit.PayPlan.Type,
			CustomLimit:   app.Limit.CustomLimit,
		},
	})

	if app.GatewayAAT != (types.GatewayAAT{}) {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  types.TableGatewayAAT,
			input: dbGatewayAATJSON{
				ApplicationID:   app.ID,
				Address:         app.GatewayAAT.Address,
				ClientPublicKey: app.GatewayAAT.ClientPublicKey,
				PrivateKey:      app.GatewayAAT.PrivateKey,
				PublicKey:       app.GatewayAAT.ApplicationPublicKey,
				Signature:       app.GatewayAAT.ApplicationSignature,
				Version:         app.GatewayAAT.Version,
			},
		})
	}

	if !gatewaySettingsIsNull(app.GatewaySettings) {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  types.TableGatewaySettings,
			input: dbGatewaySettingsJSON{
				ApplicationID:        app.ID,
				SecretKey:            app.GatewaySettings.SecretKey,
				SecretKeyRequired:    app.GatewaySettings.SecretKeyRequired,
				WhitelistOrigins:     app.GatewaySettings.WhitelistOrigins,
				WhitelistUserAgents:  app.GatewaySettings.WhitelistUserAgents,
				WhitelistBlockchains: app.GatewaySettings.WhitelistBlockchains,
			},
		})
		for _, contract := range app.GatewaySettings.WhitelistContracts {
			inputs = append(inputs, inputStruct{
				action: sideTablesAction,
				table:  types.TableWhitelistContracts,
				input: dbWhitelistContractJSON{
					ApplicationID: app.ID,
					BlockchainID:  contract.BlockchainID,
					Contracts:     contract.Contracts,
				},
			})
		}
		for _, method := range app.GatewaySettings.WhitelistMethods {
			inputs = append(inputs, inputStruct{
				action: sideTablesAction,
				table:  types.TableWhitelistMethods,
				input: dbWhitelistMethodJSON{
					ApplicationID: app.ID,
					BlockchainID:  method.BlockchainID,
					Methods:       method.Methods,
				},
			})
		}
	}

	if app.NotificationSettings != (types.NotificationSettings{}) {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  types.TableNotificationSettings,
			input: dbNotificationSettingsJSON{
				ApplicationID: app.ID,
				SignedUp:      app.NotificationSettings.SignedUp,
				Quarter:       app.NotificationSettings.Quarter,
				Half:          app.NotificationSettings.Half,
				ThreeQuarters: app.NotificationSettings.ThreeQuarters,
				Full:          app.NotificationSettings.Full,
			},
		})
	}

	return inputs
}

func appLimitInputs(sideTablesAction types.Action, content types.SavedOnDB) []inputStruct {
	appLimit := content.(*types.AppLimit)

	var inputs []inputStruct

	inputs = append(inputs, inputStruct{
		action: sideTablesAction,
		table:  types.TableAppLimits,
		input: dbAppLimitJSON{
			ApplicationID: appLimit.ID,
			PlanType:      appLimit.PayPlan.Type,
			CustomLimit:   0,
		},
	})

	return inputs
}

func blockchainInputs(mainTableAction, sideTablesAction types.Action, content types.SavedOnDB) []inputStruct {
	blockchain := content.(*types.Blockchain)

	var inputs []inputStruct

	inputs = append(inputs, inputStruct{
		action: mainTableAction,
		table:  types.TableBlockchains,
		input: dbBlockchainJSON{
			BlockchainID:      blockchain.ID,
			Altruist:          blockchain.Altruist,
			Blockchain:        blockchain.Blockchain,
			ChainID:           blockchain.ChainID,
			ChainIDCheck:      blockchain.ChainIDCheck,
			ChainPath:         blockchain.Path,
			Description:       blockchain.Description,
			EnforceResult:     blockchain.EnforceResult,
			Network:           blockchain.Network,
			Ticker:            blockchain.Ticker,
			BlockchainAliases: blockchain.BlockchainAliases,
			LogLimitBlocks:    blockchain.LogLimitBlocks,
			RequestTimeout:    blockchain.RequestTimeout,
			Active:            blockchain.Active,
			CreatedAt:         blockchain.CreatedAt.Format(psqlDateLayout),
			UpdatedAt:         blockchain.UpdatedAt.Format(psqlDateLayout),
		},
	})

	if blockchain.SyncCheckOptions != (types.SyncCheckOptions{}) {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  types.TableSyncCheckOptions,
			input: dbSyncCheckOptionsJSON{
				BlockchainID: blockchain.SyncCheckOptions.BlockchainID,
				Body:         blockchain.SyncCheckOptions.Body,
				Path:         blockchain.SyncCheckOptions.Path,
				ResultKey:    blockchain.SyncCheckOptions.ResultKey,
				Allowance:    blockchain.SyncCheckOptions.Allowance,
			},
		})
	}

	return inputs
}

func loadBalancerInputs(mainTableAction, sideTablesAction types.Action, content types.SavedOnDB) []inputStruct {
	lb := content.(*types.LoadBalancer)

	var inputs []inputStruct

	inputs = append(inputs, inputStruct{
		action: mainTableAction,
		table:  types.TableLoadBalancers,
		input: dbLoadBalancerJSON{
			LbID:              lb.ID,
			Name:              lb.Name,
			UserID:            lb.UserID,
			RequestTimeout:    lb.RequestTimeout,
			Gigastake:         lb.Gigastake,
			GigastakeRedirect: lb.GigastakeRedirect,
			CreatedAt:         lb.CreatedAt.Format(psqlDateLayout),
			UpdatedAt:         lb.UpdatedAt.Format(psqlDateLayout),
		},
	})

	if !lb.StickyOptions.IsEmpty() {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  types.TableStickinessOptions,
			input: dbStickinessOptionsJSON{
				LbID:       lb.ID,
				Duration:   lb.StickyOptions.Duration,
				Origins:    lb.StickyOptions.StickyOrigins,
				StickyMax:  lb.StickyOptions.StickyMax,
				Stickiness: lb.StickyOptions.Stickiness,
			},
		})
	}

	if len(lb.Users) != 0 {
		for _, user := range lb.Users {
			inputs = append(inputs, inputStruct{
				action: sideTablesAction,
				table:  types.TableUserAccess,
				input: dbUserAccessJSON{
					LbID:     lb.ID,
					UserID:   user.UserID,
					RoleName: string(user.RoleName),
					Email:    user.Email,
					Accepted: user.Accepted,
				},
			})
		}
	}

	for _, appID := range lb.ApplicationIDs {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  types.TableLbApps,
			input: types.LbApp{
				LbID:  lb.ID,
				AppID: appID,
			},
		})
	}

	return inputs
}

func redirectInput(action types.Action, content types.SavedOnDB) inputStruct {
	redirect := content.(*types.Redirect)

	return inputStruct{
		action: action,
		table:  types.TableRedirects,
		input: dbRedirectJSON{
			BlockchainID:   redirect.BlockchainID,
			Alias:          redirect.Alias,
			LoadBalancerID: redirect.LoadBalancerID,
			Domain:         redirect.Domain,
			CreatedAt:      redirect.CreatedAt.Format(psqlDateLayout),
			UpdatedAt:      redirect.UpdatedAt.Format(psqlDateLayout),
		},
	}
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
	case *types.Application:
		inputs = applicationInputs(mainTableAction, sideTablesAction, content)
	case *types.AppLimit:
		inputs = appLimitInputs(sideTablesAction, content)
	case *types.Blockchain:
		inputs = blockchainInputs(mainTableAction, sideTablesAction, content)
	case *types.LoadBalancer:
		inputs = loadBalancerInputs(mainTableAction, sideTablesAction, content)
	case *types.Redirect:
		inputs = []inputStruct{redirectInput(mainTableAction, content)}
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
