package types

/* Enums */
type PayPlanType string

/* Pay Plan Type and Methods */
type (
	Plan struct {
		Type     PayPlanType          `json:"planType"`
		ChainIDs map[ChainID]struct{} `json:"blockchainIDs"`
		// MonthlyRelayLimit is the number of relays-per-month for a pay plan
		MonthlyRelayLimit int32 `json:"monthlyRelayLimit"`
		// ThroughputLimit is the number of relays-per-second for a pay plan
		ThroughputLimit int32 `json:"throughputLimit"`
		// AppLimit is the number of apps permitted for a pay plan
		AppLimit int32 `json:"appLimit"`

		// TODO - remove when v2 migration finished
		// LegacyDailyLimit is the daily limit (required for legacy apps to function)
		LegacyDailyLimit int32 `json:"legacyDailyLimit"`
	}
)

func (p *Plan) Table() Table {
	return TablePayPlans
}
