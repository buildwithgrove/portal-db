package types

import "regexp"

/* String ID types */
type (
	AccountID      string
	UserID         string
	RelayChainID   string
	PortalAppID    string
	ProtocolAppID  string
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
	TableAppSettings      Table = "portal_application_settings"
	TableAppWhitelists    Table = "portal_application_whitelists"
	TableAppNotifications Table = "portal_application_notifications"

	TableAATs Table = "aats"

	TableGigastakeApplications Table = "gigastake_applications"

	TableAccounts            Table = "accounts"
	TableAccountUserAccess   Table = "account_user_access"
	TableAccountIntegrations Table = "account_integrations"

	TableUsers             Table = "users"
	TableUserAuthProviders Table = "user_auth_providers"

	TablePayPlans Table = "pay_plans"

	TableChains            Table = "chains"
	TableChainAltruists    Table = "chain_altruists"
	TableChainChecks       Table = "chain_checks"
	TableChainAliasDomains Table = "chain_alias_domains"

	TableGlobalBlockedContracts Table = "global_blocked_contracts"

	ActionInsert Action = "INSERT"
	ActionUpdate Action = "UPDATE"
	ActionDelete Action = "DELETE"
)

type SavedOnDB interface {
	Table() Table
}
