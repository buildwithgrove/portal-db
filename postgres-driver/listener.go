package postgresdriver

import (
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pokt-foundation/portal-db/v2/types"
)

type Listener interface {
	NotificationChannel() <-chan *pq.Notification
	Listen(channel string) error
}

type notification struct {
	Table  types.Table  `json:"table"`
	Action types.Action `json:"action"`
	Data   any          `json:"data"`
}

// func (n notification) parseLoadBalancerNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbLoadBalancer dbLoadBalancerJSON
// 	_ = json.Unmarshal(rawData, &dbLoadBalancer)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbLoadBalancer.toOutput(),
// 	}
// }

// func (n notification) parseStickinessOptionsNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbStickinessOpts dbStickinessOptionsJSON
// 	_ = json.Unmarshal(rawData, &dbStickinessOpts)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbStickinessOpts.toOutput(),
// 	}
// }

// func (n notification) parseUserAccessNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbUserAccess dbUserAccessJSON
// 	_ = json.Unmarshal(rawData, &dbUserAccess)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbUserAccess.toOutput(),
// 	}
// }

// func (n notification) parseLbApps() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var lbApp types.LbApp
// 	_ = json.Unmarshal(rawData, &lbApp)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   &lbApp,
// 	}
// }

// func (n notification) parseApplicationNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbApp dbAppJSON
// 	_ = json.Unmarshal(rawData, &dbApp)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbApp.toOutput(),
// 	}
// }

// func (n notification) parseAppLimitNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbAppLimit dbAppLimitJSON
// 	_ = json.Unmarshal(rawData, &dbAppLimit)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbAppLimit.toOutput(),
// 	}
// }

// func (n notification) parseGatewayAATNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbGatewayAAT dbGatewayAATJSON
// 	_ = json.Unmarshal(rawData, &dbGatewayAAT)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbGatewayAAT.toOutput(),
// 	}
// }

// func (n notification) parseGatewaySettingsNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbGatewaySettings dbGatewaySettingsJSON
// 	_ = json.Unmarshal(rawData, &dbGatewaySettings)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbGatewaySettings.toOutput(),
// 	}
// }

// func (n notification) parseWhitelistContractNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbWhitelistContract dbWhitelistContractJSON
// 	_ = json.Unmarshal(rawData, &dbWhitelistContract)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbWhitelistContract.toOutput(),
// 	}
// }

// func (n notification) parseWhitelistMethodNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbWhitelistMethod dbWhitelistMethodJSON
// 	_ = json.Unmarshal(rawData, &dbWhitelistMethod)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbWhitelistMethod.toOutput(),
// 	}
// }

// func (n notification) parseNotificationSettingsNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbNotificationSettings dbNotificationSettingsJSON
// 	_ = json.Unmarshal(rawData, &dbNotificationSettings)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbNotificationSettings.toOutput(),
// 	}
// }

// func (n notification) parseBlockchainNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbBlockchain dbBlockchainJSON
// 	_ = json.Unmarshal(rawData, &dbBlockchain)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbBlockchain.toOutput(),
// 	}
// }

// func (n notification) parseRedirectNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbRedirect dbRedirectJSON
// 	_ = json.Unmarshal(rawData, &dbRedirect)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbRedirect.toOutput(),
// 	}
// }

// func (n notification) parseSyncOptionsNotification() *types.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbSyncOpts dbSyncCheckOptionsJSON
// 	_ = json.Unmarshal(rawData, &dbSyncOpts)

// 	return &types.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbSyncOpts.toOutput(),
// 	}
// }

func (n notification) parseNotification() *types.Notification {
	// switch n.Table {
	// case types.TableLoadBalancers:
	// 	return n.parseLoadBalancerNotification()
	// case types.TableStickinessOptions:
	// 	return n.parseStickinessOptionsNotification()
	// case types.TableUserAccess:
	// 	return n.parseUserAccessNotification()

	// case types.TableLbApps:
	// 	return n.parseLbApps()

	// case types.TableApplications:
	// 	return n.parseApplicationNotification()
	// case types.TableAppLimits:
	// 	return n.parseAppLimitNotification()
	// case types.TableGatewayAAT:
	// 	return n.parseGatewayAATNotification()
	// case types.TableGatewaySettings:
	// 	return n.parseGatewaySettingsNotification()
	// case types.TableWhitelistContracts:
	// 	return n.parseWhitelistContractNotification()
	// case types.TableWhitelistMethods:
	// 	return n.parseWhitelistMethodNotification()
	// case types.TableNotificationSettings:
	// 	return n.parseNotificationSettingsNotification()

	// case types.TableBlockchains:
	// 	return n.parseBlockchainNotification()
	// case types.TableRedirects:
	// 	return n.parseRedirectNotification()
	// case types.TableSyncCheckOptions:
	// 	return n.parseSyncOptionsNotification()
	// }

	return nil
}

func parsePQNotification(n *pq.Notification, outCh chan *types.Notification) {
	if n != nil {
		var notification notification
		_ = json.Unmarshal([]byte(n.Extra), &notification)
		outCh <- notification.parseNotification()
	}
}

func Listen(inCh <-chan *pq.Notification, outCh chan *types.Notification) {
	for {
		n := <-inCh
		go parsePQNotification(n, outCh)
	}
}

func (d *PostgresDriver) CloseListener() {
	close(d.notification)
}
