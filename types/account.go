package types

import (
	"errors"
	"strings"
	"time"
)

var errNoOwner error = errors.New("account does not have an owner")

/* Account Struct Definition and Methods */
type (
	// Account represents a single account for a single application in the Portal
	Account struct {
		ID                     AccountID                    `json:"id"`
		Name                   string                       `json:"name"`
		PlanType               PayPlanType                  `json:"planType"`
		Users                  map[UserID]AccountUserAccess `json:"users"`
		PartnerChainIDs        map[RelayChainID]struct{}    `json:"partnerBlockchainIDs"`
		PartnerThroughputLimit int32                        `json:"partnerThroughputLimit"`
		PartnerAppLimit        int32                        `json:"partnerAppLimit"`
		CreatedAt              time.Time                    `json:"createdAt"`
		UpdatedAt              time.Time                    `json:"updatedAt"`
		Deleted                bool                         `json:"deleted"`

		// PortalApps and Plan are set inside PHD
		PortalApps map[PortalAppID]*PortalApp `json:"portalApps"`
		Plan       *Plan                      `json:"payPlan"`

		// TODO - remove when v2 migration finished
		// LegacyLoadBalancerID is the load balancer ID that the account was migrated from
		LegacyLoadBalancerID string `json:"legacyLoadBalancerID"`
	}

	// AccountUserAccess represents a single Portal user's role for a single Account
	AccountUserAccess struct {
		AccountID       AccountID           `json:"accountID,omitempty"`
		UserID          UserID              `json:"userID"`
		Email           Email               `json:"email"`
		RoleName        RoleName            `json:"roleName"`
		Accepted        bool                `json:"accepted"`
		ProviderUserIDs map[AuthType]string `json:"providerUserIDs"`
	}

	// CreateAccountUserAccess contains all fields required to create a new Account User
	CreateAccountUserAccess struct {
		AccountID AccountID `json:"accountID"`
		Email     Email     `json:"email"`
		RoleName  RoleName  `json:"roleName"`
	}

	// UpdateAccount contains all fields required to update an Account
	UpdateAccount struct {
		AccountID AccountID `json:"accountID"`
		Name      string    `json:"name"`
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
		if user.RoleName == RoleOwner {
			var userID string
			// in case user has two auth providers, default to using their Auth0 Username/PW ID
			switch {
			case user.ProviderUserIDs[AuthTypeAuth0Username] != "":
				userID = user.ProviderUserIDs[AuthTypeAuth0Username]
			default:
				userID = user.ProviderUserIDs[AuthTypeAuth0Github]
			}
			// user ID will be stored as `auth0|userid123` or `github|userid123`
			return strings.Split(userID, "|")[1]
		}
	}
	return ""
}

// GetOwnerEmail returns the Email of the Application OWNER
func (a *Account) GetOwner() (AccountUserAccess, error) {
	for _, userAccess := range a.Users {
		if userAccess.RoleName == RoleOwner {
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
