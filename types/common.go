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

	TableBlockchains              Table = "blockchains"
	TableChainAltruists           Table = "blockchain_altruists"
	TableChainGigastakesRedirects Table = "blockchain_gigastake_redirects"
	TableChainChecks              Table = "blockchain_checks"

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

func (t *AAT) Table() Table {
	return TableAppAATs
}

func (t *Settings) Table() Table {
	return TableAppSettings
}

func (t *Whitelists) Table() Table {
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

func (t *Plan) Table() Table {
	return TablePayPlans
}

func (b *Chain) Table() Table {
	return TableBlockchains
}

func (r *Altruist) Table() Table {
	return TableChainAltruists
}

func (r *GigastakeRedirect) Table() Table {
	return TableChainGigastakesRedirects
}

func (t *Check) Table() Table {
	return TableChainChecks
}

func (t *GlobalBlockedContracts) Table() Table {
	return TableGlobalBlockedContracts
}
