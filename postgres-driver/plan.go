package postgresdriver

import (
	"context"

	"github.com/pokt-foundation/portal-db/v2/types"
)

/* ----- postgresdriver Plan Read Methods ----- */

// ReadPlans returns all Chains in the database as Chains structs
func (pg *PostgresDriver) ReadPlans(ctx context.Context) (map[types.PayPlanType]*types.Plan, error) {
	dbPlans, err := pg.SelectPlans(ctx)
	if err != nil {
		return nil, err
	}

	plans := make(map[types.PayPlanType]*types.Plan, len(dbPlans))
	for _, dbPlan := range dbPlans {
		plan, err := dbPlan.toPlan()
		if err != nil {
			return nil, err
		}

		plans[dbPlan.PlanType] = plan
	}

	return plans, nil
}

// toPlan converts PayPlan to Plan struct
func (p *SelectPlansRow) toPlan() (*types.Plan, error) {
	plan := &types.Plan{
		Type:              p.PlanType,
		ChainIDs:          make(map[types.RelayChainID]struct{}, len(p.ChainIDs)),
		MonthlyRelayLimit: p.MonthlyRelayLimit,
		ThroughputLimit:   p.ThroughputLimit,
		AppLimit:          p.ApplicationLimit,
		LegacyDailyLimit:  p.DailyLimit.Int32,
		CreatedAt:         p.CreatedAt.UTC(),
	}

	for _, chainID := range p.ChainIDs {
		plan.ChainIDs[types.RelayChainID(chainID)] = struct{}{}
	}

	return plan, nil
}

/* ----- Used by Listener ----- */
func (json PayPlan) toOutput() *types.Plan {
	chainIDs := make(map[types.RelayChainID]struct{})
	for _, id := range json.ChainIDs {
		chainIDs[types.RelayChainID(id)] = struct{}{}
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
