package types

/* String ID types */
type (
	AccountID   int64
	ChainID     string
	Email       string
	PortalAppID string
	UserID      string
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

	TableUsers     Table = "users"
	TableUserRoles Table = "user_roles"

	TablePayPlans Table = "pay_plans"

	TableBlockchains              Table = "chains"
	TableChainAltruists           Table = "chain_altruists"
	TableChainGigastakesRedirects Table = "chain_gigastake_redirects"
	TableChainChecks              Table = "chain_checks"

	TableGlobalBlockedContracts Table = "global_blocked_contracts"

	ActionInsert Action = "INSERT"
	ActionUpdate Action = "UPDATE"
	ActionDelete Action = "DELETE"
)

type SavedOnDB interface {
	Table() Table
}
