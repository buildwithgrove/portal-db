package types

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
	TableLoadBalancers     Table = "loadbalancers"
	TableStickinessOptions Table = "stickiness_options"
	TableUserAccess        Table = "user_access"

	TableLbApps Table = "lb_apps"

	TableApplications         Table = "applications"
	TableAppLimits            Table = "app_limits"
	TableGatewayAAT           Table = "gateway_aat"
	TableGatewaySettings      Table = "gateway_settings"
	TableWhitelistContracts   Table = "whitelist_contracts"
	TableWhitelistMethods     Table = "whitelist_methods"
	TableNotificationSettings Table = "notification_settings"

	TableBlockchains      Table = "blockchains"
	TableRedirects        Table = "redirects"
	TableSyncCheckOptions Table = "sync_check_options"

	ActionInsert Action = "INSERT"
	ActionUpdate Action = "UPDATE"
	ActionDelete Action = "DELETE"
)

type SavedOnDB interface {
	Table() Table
}

func (l *LoadBalancer) Table() Table {
	return TableLoadBalancers
}
func (s *StickyOptions) Table() Table {
	return TableStickinessOptions
}
func (s *UserAccess) Table() Table {
	return TableUserAccess
}

func (l *LbApp) Table() Table {
	return TableLbApps
}

func (a *Application) Table() Table {
	return TableApplications
}
func (a *GatewayAAT) Table() Table {
	return TableGatewayAAT
}
func (s *GatewaySettings) Table() Table {
	return TableGatewaySettings
}
func (s *WhitelistContract) Table() Table {
	return TableWhitelistContracts
}
func (s *WhitelistMethod) Table() Table {
	return TableWhitelistMethods
}
func (a *AppLimit) Table() Table {
	return TableAppLimits
}
func (s *NotificationSettings) Table() Table {
	return TableNotificationSettings
}

func (b *Blockchain) Table() Table {
	return TableBlockchains
}
func (r *Redirect) Table() Table {
	return TableRedirects
}
func (o *SyncCheckOptions) Table() Table {
	return TableSyncCheckOptions
}
