package driver

import (
	"context"
	"time"

	"github.com/pokt-foundation/portal-db/types"
)

type (
	// The Driver interface represents all database operations required by the Pocket HTTP DB
	Driver interface {
		Reader
		Writer
	}

	Reader interface {
		/* ReadAccounts returns all Accounts in the database. Can specify if deleted Accounts should be included. */
		ReadAccounts(ctx context.Context, options types.DriverOptions) (map[types.AccountID]*types.Account, error)
		/* ReadChains returns all Chains in the databas. Can specify if deleted Chains should be included. */
		ReadChains(ctx context.Context, options types.DriverOptions) (map[types.ChainID]*types.Chain, error)
		/* ReadPortalApps returns all PortalApps in the database. Can specify if deleted PortalApps should be included.  */
		ReadPortalApps(ctx context.Context, options types.DriverOptions) (map[types.PortalAppID]*types.PortalApp, error)

		/* ReadPlans returns all Plans in the database */
		ReadPlans(ctx context.Context) (map[types.PayPlanType]*types.Plan, error)
		/* ReadUserPermissions returns all UserPermissions in the database as a map that takes the form map[types.UserID]*types.UserPermissions */
		ReadUserPermissions(ctx context.Context) (map[types.UserID]*types.UserPermissions, error)

		/* GetPortalUserIDFromProviderID returns a user's portal UserID */
		GetPortalUserIDFromProviderID(ctx context.Context, providerUserID string) (types.UserID, error)
		/* ReadUserByUserID returns a single user from a portal UserID */
		ReadUserByUserID(ctx context.Context, userID types.UserID) (*types.User, error)

		NotificationChannel() <-chan *types.Notification
	}

	Writer interface {
		/* WritePortalApp saves input PortalApp to the database. */
		WritePortalApp(ctx context.Context, portalApp types.PortalApp, createdAt time.Time) (*types.PortalApp, error)
		/* UpdateLoadBalancer updates PortalApp and related table rows. */
		UpdatePortalApp(ctx context.Context, update types.UpdatePortalApp, updatedAt time.Time) error
		/* SetPortalAppDeleted sets the portal app Deleted field to true. */
		SetPortalAppDeleted(ctx context.Context, portalAppID types.PortalAppID, deletedAt time.Time) error

		/* UpdatePortalAppsFirstDateSurpassed updates multiple PortalApps firstDateSurpassed field. */
		// TODO legacy method - determine if still needed and remove if not when V2 migration completed
		UpdatePortalAppsFirstDateSurpassed(ctx context.Context, update *types.UpdateFirstDateSurpassed) error

		/* WriteUserNewSignUp creates a new portal User in the DB from a CreateUser input. */
		WriteUserNewSignUp(ctx context.Context, user types.CreateUser, createdAt time.Time) error
		/* WriteUserProviderSignedUp creates a new portal UserAuthProvider in the DB when a user accepts their team invite. */
		WriteUserProviderSignedUp(ctx context.Context, userID types.UserID, user types.CreateUser, createdAt time.Time) (types.UserID, error)
		/* DeletePortalUser deletes a portal User from the DB. WARNING will do a full delete in the case of users. */
		DeletePortalUser(ctx context.Context, userID types.UserID) (types.UserID, error)

		/* WriteAccount saves input Account to the database. */
		WriteAccount(ctx context.Context, creatorID types.UserID, account types.Account, createdAt time.Time) (*types.Account, error)
		/* DeleteAccount saves input Account to the database. */
		DeleteAccount(ctx context.Context, account types.Account, deletedAt time.Time) error

		/* WriteAccountUser saves input AccountUserAccess to the database. */
		WriteAccountUser(ctx context.Context, accountID types.AccountID, accountUser types.AccountUserAccess, createdAt time.Time) (*types.AccountUserAccess, error)
		/* UpdateUserAccessRole updates the RoleName for an AccountUserAccess row. */
		UpdateAccountUserRole(ctx context.Context, email types.Email, portalAppID types.PortalAppID, roleName types.RoleName) error
		/* AcceptAccountUser sets the User ID and the Accepted field to true for an AccountUserAccess row. */
		AcceptAccountUser(ctx context.Context, email types.Email, userID types.UserID, portalAppID string) error
		/* DeleteAccountUser deletes a UserAccess row. */
		DeleteAccountUser(ctx context.Context, email types.Email, portalAppID types.PortalAppID) error

		/* WriteChain saves input Chain struct to the database. */
		WriteChain(ctx context.Context, blockchain *types.Chain) (*types.Chain, error)
		/* WriteGigastakeRedirect saves input GigastakeRedirect struct to the database for a given Chain. */
		WriteGigastakeRedirect(ctx context.Context, redirect *types.GigastakeRedirect) (*types.GigastakeRedirect, error)
		/* UpdateChain updates Chain and ChainCheck for a given Chain. */
		UpdateChain(ctx context.Context, chainID types.ChainID, update *types.UpdateChain) error
		/* ActivateChain toggles Chain.Active field on or off. */
		ActivateChain(ctx context.Context, chainID types.ChainID, active bool) error
		/* DeleteGigastakeRedirect removes a single GigastakeRedirect for a given Chain. */
		DeleteGigastakeRedirect(ctx context.Context, chainID types.ChainID, domain string) error
	}
)
