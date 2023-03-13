package types

import "time"

/* Account Struct Definition and Methods */
type (
	// Account represents a single account for a single application in the Portal
	Account struct {
		ID              AccountID                    `json:"id"`
		Plan            Plan                         `json:"payPlan"`
		Users           map[UserID]AccountUserAccess `json:"users"`
		PartnerChainIDs map[ChainID]struct{}         `json:"partnerBlockchainIDs"`
		// PartnerThroughputLimit is the number of relays per second for an accounts partners
		PartnerThroughputLimit int32 `json:"partnerThroughputLimit"`
		// PartnerAppLimit is the number of apps for an accounts partners
		PartnerAppLimit int32     `json:"partnerAppLimit"`
		CreatedAt       time.Time `json:"createdAt"`
		UpdatedAt       time.Time `json:"updatedAt"`
		Deleted         bool      `json:"deleted"`
	}

	// AccountUserAccess represents a single Portal user's role for a single Portal application
	AccountUserAccess struct {
		ID       AccountID `json:"accountID,omitempty"`
		User     User      `json:"userID"`
		RoleName RoleName  `json:"roleName"`
		Accepted bool      `json:"accepted"`
	}
)

func (a *Account) Table() Table {
	return TableAccounts
}

func (a *AccountUserAccess) Table() Table {
	return TableAccountUserAccess
}
