package types

import "time"

/* Enums */
type PayPlanType string

var (
	BasicPlan      PayPlanType = "basic_plan"
	DeveloperPlan  PayPlanType = "developer_plan"
	EnterprisePlan PayPlanType = "enterprise_plan"
	ProPlan        PayPlanType = "pro_plan"
	StartupPlan    PayPlanType = "startup_plan"
)

/* Pay Plan Type and Methods */
type (
	Plan struct {
		ID       string                    `json:"id"`
		Type     PayPlanType               `json:"planType"`
		ChainIDs map[RelayChainID]struct{} `json:"blockchainIDs"`
		// MonthlyRelayLimit is the number of relays-per-month for a pay plan
		MonthlyRelayLimit int32 `json:"monthlyRelayLimit"`
		// ThroughputLimit is the number of relays-per-second for a pay plan
		ThroughputLimit int32 `json:"throughputLimit"`
		// AppLimit is the number of apps permitted for a pay plan
		AppLimit int32 `json:"appLimit"`

		// TODO - remove when v2 migration finished
		// LegacyDailyLimit is the daily limit (required for legacy apps to function)
		LegacyDailyLimit int32     `json:"legacyDailyLimit"`
		CreatedAt        time.Time `json:"createdAt"`
	}
)

func (p *Plan) Table() Table {
	return TablePayPlans
}
