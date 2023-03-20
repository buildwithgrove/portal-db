package types

import "regexp"

/* String ID types */
type (
	AccountID      int32
	UserID         int32
	ChainID        string
	PortalAppID    string
	Email          string
	BlockedAddress string
)

// Validates that an Email fits a valid email format eg. test@example.com
func (e Email) IsValid() bool {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(emailPattern)
	return regex.MatchString(string(e))
}

/* Config Structs */

// Provides options for passing to Driver interface methods
type DriverOptions struct {
	IncludeDeleted bool
}

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
