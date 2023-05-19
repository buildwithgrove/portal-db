package types

import (
	"errors"
	"time"
)

var errNoOwner error = errors.New("account does not have an owner")

/* Account Struct Definition and Methods */
type (
	// Account represents a single account for a single application in the Portal
	Account struct {
		ID                     AccountID                    `json:"id"`
		PlanType               PayPlanType                  `json:"planType"`
		Users                  map[UserID]AccountUserAccess `json:"users"`
		PartnerChainIDs        map[RelayChainID]struct{}    `json:"partnerBlockchainIDs"`
		PartnerThroughputLimit int32                        `json:"partnerThroughputLimit"`
		PartnerAppLimit        int32                        `json:"partnerAppLimit"`
		Integrations           AccountIntegrations          `json:"integrations"`
		CreatedAt              time.Time                    `json:"createdAt"`
		UpdatedAt              time.Time                    `json:"updatedAt"`
		Deleted                bool                         `json:"deleted"`

		// PortalApps and Plan are set inside PHD
		PortalApps map[PortalAppID]*PortalApp `json:"portalApps"`
		Plan       *Plan                      `json:"payPlan"`
	}

	// AccountUserAccess represents a single Portal user for a single Account
	AccountUserAccess struct {
		AccountID              AccountID                   `json:"accountID,omitempty"`
		UserID                 UserID                      `json:"userID"`
		Email                  Email                       `json:"email"`
		Owner                  bool                        `json:"owner"`
		Accepted               bool                        `json:"accepted"`
		ProviderUserIDs        map[AuthType]ProviderUserID `json:"providerUserIDs"`
		PortalApplicationRoles map[PortalAppID]RoleName    `json:"portalApplicationRoles"`
	}

	// AccountUserAccess represents fields used for integrations with other platforms
	AccountIntegrations struct {
		AccountID          AccountID `json:"id,omitempty"`
		CovalentAPIKeyFree string    `json:"covalentAPIKeyFree"`
		CovalentAPIKeyPaid string    `json:"covalentAPIKeyPaid"`
	}

	// CreateAccountUserAccess contains all fields required to create a new Account User
	CreateAccountUserAccess struct {
		AccountID AccountID `json:"accountID"`
		Email     Email     `json:"email"`
		RoleName  RoleName  `json:"roleName"`
	}

	// UpdateAccount contains all fields required to update an Account
	UpdateAccount struct {
		AccountID AccountID   `json:"accountID"`
		PlanType  PayPlanType `json:"planType"`
	}

	// UpdateAccountUserRole contains all fields required to update an Account User's Role
	UpdateAccountUserRole struct {
		UserID    UserID    `json:"userID"`
		AccountID AccountID `json:"accountID"`
		RoleName  RoleName  `json:"roleName"`
		// TODO - remove when v2 migration finished
		// LegacyLoadBalancerID is the load balancer ID that the account was migrated from
		LegacyLoadBalancerID string `json:"legacyLoadBalancerID"`
	}

	// UpdateAccountUserRole contains all fields required to update an Account User's Role
	UpdateAcceptAccountUser struct {
		AccountID        AccountID      `json:"accountID"`
		UserID           UserID         `json:"userID"`
		AuthProviderType AuthType       `json:"type"`
		ProviderUserID   ProviderUserID `json:"providerUserID"`
	}
)

// LegacyUserID returns the legacy user ID for an Account
func (a *Account) LegacyUserID() string {
	for _, user := range a.Users {
		if user.Owner {
			var userID ProviderUserID
			// in case user has two auth providers, default to using their Auth0 Username/PW ID
			switch {
			case user.ProviderUserIDs[AuthTypeAuth0Username] != "":
				userID = user.ProviderUserIDs[AuthTypeAuth0Username]
			default:
				userID = user.ProviderUserIDs[AuthTypeAuth0Github]
			}
			return string(userID)
		}
	}
	return ""
}

// GetOwnerEmail returns the Email of the Application OWNER
func (a *Account) GetOwner() (AccountUserAccess, error) {
	for _, userAccess := range a.Users {
		if userAccess.Owner {
			return userAccess, nil
		}
	}

	return AccountUserAccess{}, errNoOwner
}

func (a *Account) Table() Table {
	return TableAccounts
}

func (a *AccountUserAccess) Table() Table {
	return TableAccountUserAccess
}

func (s *AccountIntegrations) Table() Table {
	return TableAccountIntegrations
}
