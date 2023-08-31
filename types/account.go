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
		Name                   string                       `json:"name"`
		IconURL                string                       `json:"iconURL"`
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
		AccountID      AccountID                `json:"id,omitempty"` // used for listener
		Owner          bool                     `json:"owner"`
		UserID         UserID                   `json:"userID"`
		Email          Email                    `json:"email"`
		Accepted       bool                     `json:"accepted"`
		PortalAppRoles map[PortalAppID]RoleName `json:"portalApplicationRoles"`
	}

	// AccountUserAccess represents fields used for integrations with other platforms
	AccountIntegrations struct {
		AccountID          AccountID `json:"id,omitempty"` // used for listener
		CovalentAPIKeyFree string    `json:"covalentAPIKeyFree"`
		CovalentAPIKeyPaid string    `json:"covalentAPIKeyPaid"`
	}

	// CreateAccountUserAccess contains all fields required to create a new Account User
	CreateAccountUserAccess struct {
		AccountID   AccountID   `json:"accountID"`
		PortalAppID PortalAppID `json:"portalAppID"`
		Email       Email       `json:"email"`
		RoleName    RoleName    `json:"roleName"`
	}

	// UpdateAccount contains all fields required to update an Account
	UpdateAccount struct {
		AccountID AccountID   `json:"accountID"`
		PlanType  PayPlanType `json:"planType"`
	}

	// UpdateAccountUserRole contains all fields required to update an Account User's Role
	UpdateAccountUserRole struct {
		PortalAppID PortalAppID `json:"portalAppID"`
		UserID      UserID      `json:"userID"`
		AccountID   AccountID   `json:"accountID"`
		RoleName    RoleName    `json:"roleName"`
	}

	// UpdateAccountUserRole contains all fields required to accept an Account User
	UpdateAcceptAccountUser struct {
		PortalAppID      PortalAppID    `json:"portalAppID"`
		UserID           UserID         `json:"userID"`
		AuthProviderType AuthType       `json:"type"`
		ProviderUserID   ProviderUserID `json:"providerUserID"`
	}

	// UpdateRemoveAccountUser contains all fields required to remove an Account User's Role
	UpdateRemoveAccountUser struct {
		UserID      UserID      `json:"userID"`
		PortalAppID PortalAppID `json:"portalAppID"`
		AccountID   AccountID   `json:"accountID"`
	}
)

// GetOwner returns the Account OWNER
func (a *Account) GetOwner() (AccountUserAccess, error) {
	for _, user := range a.Users {
		if user.Owner {
			return user, nil
		}
	}
	return AccountUserAccess{}, errNoOwner
}

// GetOwnerID returns the Account OWNER's ID
func (a *Account) GetOwnerID() (UserID, error) {
	for _, user := range a.Users {
		if user.Owner {
			return user.UserID, nil
		}
	}
	return UserID(""), errNoOwner
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
