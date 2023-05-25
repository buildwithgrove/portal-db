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

func (n notification) parseAccountNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbAccount dbAccount
	_ = json.Unmarshal(rawData, &dbAccount)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbAccount.toOutput(),
	}
}

func (n notification) parseAccountUserAccessNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbAccountUser dbAccountUserAccess
	_ = json.Unmarshal(rawData, &dbAccountUser)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbAccountUser.toOutput(),
	}
}

func (n notification) parseAccountIntegrationsNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbAccountIntegrations dbAccountIntegration
	_ = json.Unmarshal(rawData, &dbAccountIntegrations)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbAccountIntegrations.toOutput(),
	}
}

func (n notification) parsePortalAppNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbPortalApp dbPortalApplication
	_ = json.Unmarshal(rawData, &dbPortalApp)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbPortalApp.toOutput(),
	}
}

func (n notification) parseAATNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbAAT dbPortalApplicationAAT
	_ = json.Unmarshal(rawData, &dbAAT)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbAAT.toOutput(),
	}
}

func (n notification) parseSettingsNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbAppSettings dbPortalApplicationSetting
	_ = json.Unmarshal(rawData, &dbAppSettings)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbAppSettings.toOutput(),
	}
}

func (n notification) parseWhitelistNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbAppWhitelist dbPortalApplicationWhitelist
	_ = json.Unmarshal(rawData, &dbAppWhitelist)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbAppWhitelist.toOutput(),
	}
}

func (n notification) parseAppNotificationNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbAppNotification dbPortalApplicationNotification
	_ = json.Unmarshal(rawData, &dbAppNotification)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbAppNotification.toOutput(),
	}
}

func (n notification) parseChainNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbChain dbChain
	_ = json.Unmarshal(rawData, &dbChain)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbChain.toOutput(),
	}
}

func (n notification) parseChainAltruistNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbChainAltruist dbChainAltruist
	_ = json.Unmarshal(rawData, &dbChainAltruist)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbChainAltruist.toOutput(),
	}
}

func (n notification) parseChainCheckNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbChainCheck dbChainCheck
	_ = json.Unmarshal(rawData, &dbChainCheck)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbChainCheck.toOutput(),
	}
}

func (n notification) parseChainGigastakeAliasNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbAlias dbChainAliasDomains
	_ = json.Unmarshal(rawData, &dbAlias)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbAlias.toOutput(),
	}
}

func (n notification) parseUsersNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbUser dbUser
	_ = json.Unmarshal(rawData, &dbUser)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbUser.toOutput(),
	}
}

func (n notification) parseUserAuthProviderNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbUserAuthProvider dbUserAuthProvider
	_ = json.Unmarshal(rawData, &dbUserAuthProvider)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbUserAuthProvider.toOutput(),
	}
}

func (n notification) parsePayPlanNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbPayPlan dbPayPlan
	_ = json.Unmarshal(rawData, &dbPayPlan)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbPayPlan.toOutput(),
	}
}

func (n notification) parseGlobalBlockedContractNotification() *types.Notification {
	rawData, _ := json.Marshal(n.Data)
	var dbBlockedContract dbGlobalBlockedContract
	_ = json.Unmarshal(rawData, &dbBlockedContract)

	return &types.Notification{
		Table:  n.Table,
		Action: n.Action,
		Data:   dbBlockedContract.toOutput(),
	}
}

func (n notification) parseNotification() *types.Notification {
	switch n.Table {
	case types.TableAccounts:
		return n.parseAccountNotification()
	case types.TableAccountUserAccess:
		return n.parseAccountUserAccessNotification()
	case types.TableAccountIntegrations:
		return n.parseAccountIntegrationsNotification()

	case types.TablePortalApps:
		return n.parsePortalAppNotification()
	case types.TableAATs:
		return n.parseAATNotification()
	case types.TableAppSettings:
		return n.parseSettingsNotification()
	case types.TableAppWhitelists:
		return n.parseWhitelistNotification()
	case types.TableAppNotifications:
		return n.parseAppNotificationNotification()

	case types.TableChains:
		return n.parseChainNotification()
	case types.TableChainAltruists:
		return n.parseChainAltruistNotification()
	case types.TableChainGigastakeAliases:
		return n.parseChainGigastakeAliasNotification()
	case types.TableChainChecks:
		return n.parseChainCheckNotification()

	case types.TableUsers:
		return n.parseUsersNotification()
	case types.TableUserAuthProviders:
		return n.parseUserAuthProviderNotification()

	case types.TablePayPlans:
		return n.parsePayPlanNotification()

	case types.TableGlobalBlockedContracts:
		return n.parseGlobalBlockedContractNotification()

	default:
		return nil
	}
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
