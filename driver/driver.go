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

	// TODO --> update all read methods to return maps

	Reader interface {
		/* ReadAccounts returns all PortalApps in the database */
		ReadAccounts(ctx context.Context) ([]*types.Account, error)
		/* ReadChains returns all blockchains in the database and marshals to types struct */
		ReadChains(ctx context.Context) ([]*types.Chain, error)
		/* ReadPlans returns all Plans in the database */
		ReadPlans(ctx context.Context) ([]*types.Plan, error)
		/* ReadPortalApps returns all PortalApps in the database */
		ReadPortalApps(ctx context.Context) (map[types.PortalAppID]*types.PortalApp, error)
		/* ReadUserPermissions returns all UserPermissions in the database as a map that takes the form map[types.UserID]*types.UserPermissions */
		ReadUserPermissions(ctx context.Context) (map[types.UserID]*types.UserPermissions, error)

		NotificationChannel() <-chan *types.Notification
	}

	Writer interface {
		/* WritePortalApp saves input LoadBalancer to the database */
		WritePortalApp(ctx context.Context, portalApp *types.PortalApp) (*types.PortalApp, error)
		/* UpdateLoadBalancer updates PortalApp and related table rows */
		UpdatePortalApp(ctx context.Context, portalAppID types.PortalAppID, options *types.UpdatePortalApp) error
		/* UpdatePortalAppsFirstDateSurpassed updates multiple Applications firstDateSurpassed field */
		UpdatePortalAppsFirstDateSurpassed(ctx context.Context, update *types.UpdateFirstDateSurpassed) error
		/* DeletePortalApp sets the  portal app Deleted field to true (will not appear in Portal API or UI) */
		DeletePortalApp(ctx context.Context, portalAppID types.PortalAppID) error

		/* WriteAccount saves input Account to the database */
		WriteAccount(ctx context.Context, account types.Account) error
		/* DeleteAccount saves input Account to the database */
		DeleteAccount(ctx context.Context, account types.Account) error

		/* WriteAccountUser saves input AccountUserAccess to the database */
		WriteAccountUser(ctx context.Context, portalAppID types.PortalAppID, accountUser types.AccountUserAccess) error
		/* UpdateUserAccessRole updates the RoleName for a UserAccess row */
		UpdateAccountUserRole(ctx context.Context, email types.Email, portalAppID types.PortalAppID, roleName types.RoleName) error
		/* AcceptAccountUser sets the User ID and the Accepted field to true for an AccountUserAccess row */
		AcceptAccountUser(ctx context.Context, email types.Email, userID types.UserID, portalAppID string) error
		/* DeleteAccountUser deletes a UserAccess row */
		DeleteAccountUser(ctx context.Context, email types.Email, portalAppID types.PortalAppID) error

		/* WriteChain saves input Chain struct to the database */
		WriteChain(ctx context.Context, blockchain *types.Chain) (*types.Chain, error)
		/* WriteGigastakeRedirect saves input Redirect struct to the database .*/
		WriteGigastakeRedirect(ctx context.Context, redirect *types.GigastakeRedirect) (*types.GigastakeRedirect, error)
		/* UpdateChain updates Blockchain and Sync Check Options */
		UpdateChain(ctx context.Context, blockchainID types.BlockchainID, update *types.UpdateChain) error
		/* ActivateChain toggles chain.active field on or off */
		ActivateChain(ctx context.Context, blockchainID types.BlockchainID, active bool) error
		/* DeleteGigastakeRedirect removes a single GigastakeRedirect for a chain */
		DeleteGigastakeRedirect(ctx context.Context, blockchainID types.BlockchainID, domain string) error
	}
)
