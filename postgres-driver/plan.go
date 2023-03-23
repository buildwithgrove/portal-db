package postgresdriver

import "github.com/pokt-foundation/portal-db/v2/types"

/* Used by Listener */
func (json PayPlan) toOutput() *types.Plan {
	chainIDs := make(map[types.ChainID]struct{})
	for _, id := range json.ChainIDs {
		chainIDs[types.ChainID(id)] = struct{}{}
	}

	return &types.Plan{
		Type:              json.PlanType,
		ChainIDs:          chainIDs,
		MonthlyRelayLimit: json.MonthlyRelayLimit,
		ThroughputLimit:   json.ThroughputLimit,
		AppLimit:          json.ApplicationLimit,
		LegacyDailyLimit:  json.DailyLimit.Int32,
	}
}
