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
	}

	// UpdateAccountUserRole contains all fields required to update an Account User's Role
	UpdateAcceptAccountUser struct {
		AccountID        AccountID `json:"accountID"`
		UserID           UserID    `json:"userID"`
		AuthProviderType AuthType  `json:"type"`
		ProviderUserID   string    `json:"providerUserID"`
	}
)

func (a *Account) Table() Table {
	return TableAccounts
}

func (a *AccountUserAccess) Table() Table {
	return TableAccountUserAccess
}
