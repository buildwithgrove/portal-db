package types

import "time"

/* Account Struct Definition and Methods */
type (
	// Account represents a single account for a single application in the Portal
	Account struct {
		ID   AccountID `json:"id"`
		Plan Plan      `json:"payPlan"`
		// TODO: would this be needed?
		Apps                 map[PortalAppID]*PortalApp   `json:"apps"`
		Users                map[UserID]AccountUserAccess `json:"users"`
		PartnerBlockchainIDs map[ChainID]struct{}         `json:"partnerBlockchainIDs"`
		// PartnerThroughputLimit is the number of relays per second for an accounts partners
		PartnerThroughputLimit int `json:"partnerThroughputLimit"`
		// PartnerAppLimit is the number of apps for an accounts partners
		PartnerAppLimit int `json:"partnerAppLimit"`
	}

	// AccountUserAccess represents a single Portal user's role for a single Portal application
	AccountUserAccess struct {
		AppID     string    `json:"appID,omitempty"`
		ID        string    `json:"id"`
		User      User      `json:"userID"`
		RoleName  RoleName  `json:"roleName"`
		Accepted  bool      `json:"accepted"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)

func (a *Account) Table() Table {
	return TableAccounts
}

func (a *AccountUserAccess) Table() Table {
	return TableAccountUserAccess
}
