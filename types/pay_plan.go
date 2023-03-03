package types

type (
	PayPlan struct {
		ID            string                    `json:"id"`
		Type          PayPlanType               `json:"planType"`
		BlockchainIDs map[BlockchainID]struct{} `json:"blockchainIDs"`
		// MonthlyRelayLimit is the number of relays-per-month for a pay plan
		MonthlyRelayLimit int `json:"monthlyRelayLimit"`
		// ThroughputLimit is the number of relays-per-second for a pay plan
		ThroughputLimit int `json:"throughputLimit"`
		// AppLimit is the number of apps permitted for a pay plan
		AppLimit int `json:"appLimit"`
	}

	PayPlanType string
)

const (
	// TODO will be updating plan types
	TestPlanV0   PayPlanType = "TEST_PLAN_V0"
	TestPlan10K  PayPlanType = "TEST_PLAN_10K"
	TestPlan90k  PayPlanType = "TEST_PLAN_90K"
	FreetierV0   PayPlanType = "FREETIER_V0"
	PayAsYouGoV0 PayPlanType = "PAY_AS_YOU_GO_V0"
	Enterprise   PayPlanType = "ENTERPRISE"
)

var (
	ValidPayPlanTypes = map[PayPlanType]bool{
		"":           true, // needs to be allowed while the change for all apps to have plans is done
		TestPlanV0:   true,
		TestPlan10K:  true,
		TestPlan90k:  true,
		FreetierV0:   true,
		PayAsYouGoV0: true,
		Enterprise:   true,
	}
)

func (t *Plan) Table() Table {
	return TablePayPlans
}
