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
		/* ReadPayPlans returns all pay plans in the database and marshals to types struct */
		ReadPayPlans(ctx context.Context) ([]*types.PayPlan, error)
		/* ReadApplications returns all Applications in the database */
		ReadApplications(ctx context.Context) ([]*types.Application, error)
		/* ReadLoadBalancers returns all LoadBalancers in the database */
		ReadLoadBalancers(ctx context.Context) ([]*types.LoadBalancer, error)
		/* ReadUserPermissions returns all UserPermissions in the database as a map that takes the form map[types.UserID]*types.UserPermissions */
		ReadUserPermissions(ctx context.Context) (map[types.UserID]*types.UserPermissions, error)
		/* ReadBlockchains returns all blockchains in the database and marshals to types struct */
		ReadBlockchains(ctx context.Context) ([]*types.Blockchain, error)

		NotificationChannel() <-chan *types.Notification
	}

	Writer interface {
		/* WriteLoadBalancer saves input LoadBalancer to the database */
		WriteLoadBalancer(ctx context.Context, loadBalancer *types.LoadBalancer) (*types.LoadBalancer, error)
		/* WriteLoadBalancerUser saves input LoadBalancer to the database */
		WriteLoadBalancerUser(ctx context.Context, lbID string, userAccess types.UserAccess) error
		/* UpdateLoadBalancer updates LoadBalancer and related table rows */
		UpdateLoadBalancer(ctx context.Context, id string, options *types.UpdateLoadBalancer) error
		/* UpdateUserAccessRole updates the RoleName for a UserAccess row */
		UpdateUserAccessRole(ctx context.Context, email, lbID string, roleName types.RoleName) error
		/* AcceptUserAccess sets the User ID and the Accepted field to true for a UserAccess row */
		AcceptUserAccess(ctx context.Context, email, userID, lbID string) error
		/* RemoveLoadBalancer sets the user ID to an empty string (will not appear in Portal API or UI) */
		RemoveLoadBalancer(ctx context.Context, id string) error
		/* RemoveUserAccess deletes a UserAccess row */
		RemoveUserAccess(ctx context.Context, email, lbID string) error

		/* WriteApplication saves input Application to the database */
		WriteApplication(ctx context.Context, app *types.Application) (*types.Application, error)
		/* UpdateApplication updates Application and related table rows */
		UpdateApplication(ctx context.Context, id string, update *types.UpdateApplication) error
		/* UpdateAppFirstDateSurpassed updates Application's firstDateSurpassed field */
		UpdateAppFirstDateSurpassed(ctx context.Context, update *types.UpdateFirstDateSurpassed) error
		/* RemoveApplication updates Application's status field to AwaitingGracePeriod */
		RemoveApplication(ctx context.Context, id string) error

		/* WriteBlockchain saves input Blockchain struct to the database */
		WriteBlockchain(ctx context.Context, blockchain *types.Blockchain) (*types.Blockchain, error)

		/* WriteRedirect saves input Redirect struct to the database.
		   It must be called separately from WriteBlockchain due to how new chains are added to the dB */
		WriteRedirect(ctx context.Context, redirect *types.Redirect) (*types.Redirect, error)
		/* UpdateChain updates Blockchain and Sync Check Options */
		UpdateChain(ctx context.Context, blockchainID string, update *types.UpdateBlockchain) error
		/* ActivateChain toggles chain.active field on or off */
		ActivateChain(ctx context.Context, id string, active bool) error
		RemoveRedirect(ctx context.Context, blockchainID, domain string) error
	}
)
