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

	TableUsers     Table = "users"
	TableUserRoles Table = "user_roles"

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
