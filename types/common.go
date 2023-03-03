package types

/* Common string types */
type (
	ApplicationID string
	BlockchainID  string
	UserID        string
)

/* Listener types */
type (
	Table  string
	Action string

	Notification struct {
		Table  Table
		Action Action
		Data   SavedOnDB
	}
)

const (
	TablePortalApps       Table = "portal_applications"
	TableAppAATs          Table = "portal_application_aats"
	TableAppSettings      Table = "portal_application_settings"
	TableAppWhitelists    Table = "portal_application_whitelists"
	TableAppNotifications Table = "portal_application_notifications"

	TableAccounts          Table = "accounts"
	TableAccountUserAccess Table = "account_user_access"

	TableUsers Table = "users"

	TablePayPlans Table = "pay_plans"

	TableBlockchains               Table = "blockchains"
	TableChainAltruists            Table = "blockchain_altruists"
	TableChainGigastakesRedirect   Table = "blockchain_gigastakes_redirects"
	TableChainSyncCheckOptions     Table = "blockchain_sync_check_options"
	TableChainGlobalAllowedMethods Table = "blockchain_global_allowed_methods"

	TableGlobalBlockedContracts Table = "global_blocked_contracts"

	ActionInsert Action = "INSERT"
	ActionUpdate Action = "UPDATE"
	ActionDelete Action = "DELETE"
)

type SavedOnDB interface {
	Table() Table
}

func (t *PortalApp) Table() Table {
	return TablePortalApps
}

func (t *AppAAT) Table() Table {
	return TableAppAATs
}

func (t *AppSettings) Table() Table {
	return TableAppSettings
}

func (t *AppWhitelists) Table() Table {
	return TableAppWhitelists
}

func (t *AppNotification) Table() Table {
	return TableAppNotifications
}

func (t *Account) Table() Table {
	return TableAccounts
}

func (t *AccountUserAccess) Table() Table {
	return TableAccountUserAccess
}

func (t *User) Table() Table {
	return TableUsers
}

func (t *PayPlan) Table() Table {
	return TablePayPlans
}

func (t *Blockchain) Table() Table {
	return TableBlockchains
}

func (t *ChainAltruist) Table() Table {
	return TableChainAltruists
}

func (t *ChainGigastakesRedirect) Table() Table {
	return TableChainGigastakesRedirect
}

func (t *ChainSyncCheckOptions) Table() Table {
	return TableChainSyncCheckOptions
}

func (t *ChainGlobalAllowedMethods) Table() Table {
	return TableChainGlobalAllowedMethods
}

func (t *GlobalBlockedContracts) Table() Table {
	return TableGlobalBlockedContracts
}
