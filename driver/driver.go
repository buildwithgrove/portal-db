package driver

import (
	"context"

	"github.com/pokt-foundation/portal-db/types"
)

type (
	// The Driver interface represents all database operations required by the Pocket HTTP DB
	Driver interface {
		Reader
		Writer
	}

	Reader interface {
		ReadPayPlans(ctx context.Context) ([]*types.PayPlan, error)
		ReadApplications(ctx context.Context) ([]*types.Application, error)
		ReadLoadBalancers(ctx context.Context) ([]*types.LoadBalancer, error)
		ReadUserPermissions(ctx context.Context) (map[types.UserID]types.UserPermissions, error)
		ReadBlockchains(ctx context.Context) ([]*types.Blockchain, error)

		NotificationChannel() <-chan *types.Notification
	}

	Writer interface {
		WriteLoadBalancer(ctx context.Context, loadBalancer *types.LoadBalancer) (*types.LoadBalancer, error)
		WriteLoadBalancerUser(ctx context.Context, lbID string, userAccess types.UserAccess) error
		UpdateLoadBalancer(ctx context.Context, id string, options *types.UpdateLoadBalancer) error
		UpdateUserAccessRole(ctx context.Context, userID, lbID string, roleName types.RoleName) error
		RemoveLoadBalancer(ctx context.Context, id string) error
		RemoveUserAccess(ctx context.Context, userID, lbID string) error

		WriteApplication(ctx context.Context, app *types.Application) (*types.Application, error)
		UpdateApplication(ctx context.Context, id string, update *types.UpdateApplication) error
		UpdateAppFirstDateSurpassed(ctx context.Context, update *types.UpdateFirstDateSurpassed) error
		RemoveApplication(ctx context.Context, id string) error

		WriteBlockchain(ctx context.Context, blockchain *types.Blockchain) (*types.Blockchain, error)
		WriteRedirect(ctx context.Context, redirect *types.Redirect) (*types.Redirect, error)
		ActivateChain(ctx context.Context, id string, active bool) error
	}
)
