package types

import (
	"strings"
	"time"
)

/* Account Struct Definition and Methods */
type (
	// Account represents a single account for a single application in the Portal
	Account struct {
		ID                     AccountID                    `json:"id"`
		Name                   string                       `json:"name"`
		Plan                   Plan                         `json:"payPlan"`
		Users                  map[UserID]AccountUserAccess `json:"users"`
		PortalApps             map[PortalAppID]*PortalApp   `json:"portalApps"`
		PartnerChainIDs        map[ChainID]struct{}         `json:"partnerBlockchainIDs"`
		PartnerThroughputLimit int32                        `json:"partnerThroughputLimit"`
		PartnerAppLimit        int32                        `json:"partnerAppLimit"`
		CreatedAt              time.Time                    `json:"createdAt"`
		UpdatedAt              time.Time                    `json:"updatedAt"`
		Deleted                bool                         `json:"deleted"`

		// TODO - remove when v2 migration finished
		// LegacyLoadBalancerID is the load balancer ID that the account was migrated from
		LegacyLoadBalancerID string `json:"legacyLoadBalancerID"`
	}

	// AccountUserAccess represents a single Portal user's role for a single Account
	AccountUserAccess struct {
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
		AccountID        AccountID `json:"accountID"`
		UserID           UserID    `json:"userID"`
		AuthProviderType AuthType  `json:"type"`
		ProviderUserID   string    `json:"providerUserID"`
	}
)

// LegacyDailyLimit returns the legacy daily relay limit for a given application (temporary)
func (a *Account) LegacyUserID() string {
	for _, user := range a.Users {
		if user.RoleName == RoleOwner {
			for _, userID := range user.ProviderUserIDs {
				return strings.Split(userID, "|")[1]
			}
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

	return AccountUserAccess{}, ErrNoOwner
}

func (a *Account) Table() Table {
	return TableAccounts
}

func (a *AccountUserAccess) Table() Table {
	return TableAccountUserAccess
}
