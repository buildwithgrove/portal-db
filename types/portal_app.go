package types

import (
	"sort"
	"time"
)

/* Enums */
type (
	Environment string

	NotificationType  string
	NotificationEvent string

	Origin    string
	UserAgent string
	Method    string
	Contract  string

	AppWhitelistType   string
	ChainWhitelistType string
)

const (
	EnvProduction Environment = "production"
	EnvTest       Environment = "test"

	NotificationEmail   NotificationType = "email"
	NotificationWebhook NotificationType = "webhook"
	NotificationPortal  NotificationType = "portal"

	EventSignedUp      NotificationEvent = "signedUp"
	EventQuarter       NotificationEvent = "quarter"
	EventHalf          NotificationEvent = "half"
	EventThreeQuarters NotificationEvent = "threeQuarters"
	EventFull          NotificationEvent = "full"

	WLOrigins     AppWhitelistType = "origins"
	WLBlockchains AppWhitelistType = "blockchains"
	WLUserAgents  AppWhitelistType = "userAgents"

	WLMethods   ChainWhitelistType = "methods"
	WLContracts ChainWhitelistType = "contracts"
)

/* PortalApp Struct Definition and Methods */
type (
	// PortalApp represents a single application in the Portal
	PortalApp struct {
		ID            string                               `json:"id"`
		Name          string                               `json:"name"`
		Gigastake     bool                                 `json:"gigastake"`
		Staked        bool                                 `json:"staked"`
		Account       Account                              `json:"account"`
		AAT           AAT                                  `json:"aat"`
		Settings      Settings                             `json:"settings"`
		Whitelists    Whitelists                           `json:"whitelists"`
		Notifications map[NotificationType]AppNotification `json:"notifications"`
		CreatedAt     time.Time                            `json:"createdAt"`
		UpdatedAt     time.Time                            `json:"updatedAt"`
		// TODO - remove when v2 migration finished
		// Fields required for compatibility with the old Portal API and Services (temporary)
		LegacyFields LegacyFields `json:"legacyFields"`
	}

	// TODO - remove when v2 migration finished
	// Fields required for compatibility with the old Portal API and Services (temporary)
	LegacyFields struct {
		ApplicationID      string        `json:"applicationID"`
		CustomLimit        int           `json:"customLimit"`
		RequestTimeout     int           `json:"requestTimeout"`
		GigastakeRedirect  bool          `json:"gigastakeRedirect"`
		FirstDateSurpassed time.Time     `json:"firstDateSurpassed"`
		StickyOptions      StickyOptions `json:"stickyOptions"`
	}

	AAT struct {
		AppID           ApplicationID `json:"appID,omitempty"`
		Address         string        `json:"address"`
		PublicKey       string        `json:"publicKey"`
		ClientPublicKey string        `json:"clientPublicKey"`
		PrivateKey      string        `json:"privateKey"`
		Signature       string        `json:"signature"`
		Version         string        `json:"version"`
	}

	Settings struct {
		AppID                  ApplicationID             `json:"appID,omitempty"`
		Environment            Environment               `json:"environment"`
		SecretKey              string                    `json:"secretKey"`
		SecretKeyRequired      bool                      `json:"secretKeyRequired"`
		FavoritedBlockchainIDs map[BlockchainID]struct{} `json:"favoritedBlockchainIDs"`
		// MonthlyRelayLimit sets the monthly limit per-application
		// Sum of an Account's Apps MonthlyRelayLimits cannot exceed the Account's MonthlyRelayLimit
		MonthlyRelayLimit int `json:"monthlyRelayLimit"`
	}

	Whitelists struct {
		AppID       ApplicationID                          `json:"appID,omitempty"`
		Origins     map[Origin]struct{}                    `json:"origins"`
		UserAgents  map[UserAgent]struct{}                 `json:"userAgents"`
		Blockchains map[BlockchainID]struct{}              `json:"blockchains"`
		Contracts   map[BlockchainID]map[Contract]struct{} `json:"contracts"`
		Methods     map[BlockchainID]map[Method]struct{}   `json:"methods"`
	}

	AppNotification struct {
		AppID       string                     `json:"appID,omitempty"`
		Active      bool                       `json:"active"`
		Destination string                     `json:"destination"`
		Trigger     string                     `json:"trigger"`
		Events      map[NotificationEvent]bool `json:"events"`
	}

	// WhitelistsObject is a GraphQL-compatible representation of all the whitelists for a given application (used for the Portal UI)
	WhitelistsObject struct {
		AppWhitelists   [3]ApplicationWhitelists `json:"appWhitelists"`
		ChainWhitelists [2]ChainWhitelists       `json:"chainWhitelists"`
	}
	ApplicationWhitelists struct {
		Type   AppWhitelistType `json:"type"`
		Values []string         `json:"values"`
	}
	ChainWhitelists struct {
		Type   ChainWhitelistType       `json:"type"`
		Values []BlockchainIDWhitelists `json:"values"`
	}
	BlockchainIDWhitelists struct {
		BlockchainID string   `json:"blockchainID"`
		Values       []string `json:"values"`
	}

	// UpdatePortalApp Struct Definition and Methods
	UpdatePortalApp struct {
		AppID         ApplicationID           `json:"appID,omitempty"`
		Name          string                  `json:"name,omitempty"`
		Settings      *UpdateAppSettings      `json:"appSettings,omitempty"`
		Notifications *UpdateAppNotifications `json:"notificationSettings,omitempty"`
		Whitelists    *WhitelistsObject       `json:"whitelists,omitempty"`
	}

	UpdateAppSettings struct {
		AppID             ApplicationID `json:"appID,omitempty"`
		Environment       Environment   `json:"environment"`
		SecretKey         string        `json:"secretKey"`
		SecretKeyRequired bool          `json:"secretKeyRequired"`
		MonthlyRelayLimit int           `json:"monthlyRelayLimit"`
		FavoritedChainIDs []string      `json:"favoritedChainIDs"`
	}

	UpdateAppNotifications struct {
		AppID            ApplicationID              `json:"appID,omitempty"`
		NotificationType NotificationType           `json:"notificationType"`
		Active           *bool                      `json:"active"`
		Destination      string                     `json:"destination"`
		Trigger          string                     `json:"trigger"`
		Events           map[NotificationEvent]bool `json:"events"`
	}

	UpdateFirstDateSurpassed struct {
		ApplicationIDs     []string  `json:"applicationIDs"`
		FirstDateSurpassed time.Time `json:"firstDateSurpassed"`
	}
)

// LegacyDailyLimit returns the legacy daily relay limit for a given application (temporary)
func (a *PortalApp) LegacyDailyLimit() int {
	return a.Account.Plan.LegacyDailyLimit
}

// UserID returns the UserID of the Application OWNER
func (a *PortalApp) UserID() UserID {
	for userID, user := range a.Account.Users {
		if user.RoleName == RoleOwner {
			return userID
		}
	}
	return ""
}

// MonthlyLimit returns the monthly relay limit for a given application
func (a *PortalApp) MonthlyLimit() int {
	return a.Settings.MonthlyRelayLimit
}

// IsOriginWhitelisted returns a boolean indicating whether the given ORIGIN is whitelisted for an application
func (a *PortalApp) IsOriginWhitelisted(origin Origin) bool {
	_, ok := a.Whitelists.Origins[origin]
	return ok
}

// IsUserAgentWhitelisted returns a boolean indicating whether the given USER AGENT is whitelisted for an application
func (a *PortalApp) IsUserAgentWhitelisted(userAgent UserAgent) bool {
	_, ok := a.Whitelists.UserAgents[userAgent]
	return ok
}

// IsBlockchainWhitelisted returns a boolean indicating whether the given BLOCKCHAIN is whitelisted for an application
func (a *PortalApp) IsBlockchainWhitelisted(blockchain BlockchainID) bool {
	_, ok := a.Whitelists.Blockchains[blockchain]
	return ok
}

// IsContractWhitelisted returns a boolean indicating whether the given CONTRACT is whitelisted for a blockchain and application
func (a *PortalApp) IsContractWhitelisted(blockchainID BlockchainID, contract Contract) bool {
	if chainContracts, contractsOK := a.Whitelists.Contracts[blockchainID]; contractsOK {
		if _, contractOK := chainContracts[contract]; contractOK {
			return true
		}
	}
	return false
}

// IsMethodWhitelisted returns a boolean indicating whether the given METHOD is whitelisted for a blockchain and application
func (a *PortalApp) IsMethodWhitelisted(blockchainID BlockchainID, method Method) bool {
	if chainMethods, methodsOK := a.Whitelists.Methods[blockchainID]; methodsOK {
		if _, methodOK := chainMethods[method]; methodOK {
			return true
		}
	}
	return false
}

// GetWhitelistsObject returns a GraphQL-compatible WhitelistsObject struct that contains all whitelists for an application (used for Portal UI)
func (a *PortalApp) GetWhitelistsObject() *WhitelistsObject {
	var origins, userAgents, blockchains []string // App whitelists

	for origin := range a.Whitelists.Origins {
		origins = append(origins, string(origin))
	}
	sort.Strings(origins)

	for userAgent := range a.Whitelists.UserAgents {
		userAgents = append(userAgents, string(userAgent))
	}
	sort.Strings(userAgents)

	for blockchain := range a.Whitelists.Blockchains {
		blockchains = append(blockchains, string(blockchain))
	}
	sort.Strings(blockchains)

	var contractWhitelists, methodWhitelists []BlockchainIDWhitelists // Chain whitelists

	for blockchainID, chainContracts := range a.Whitelists.Contracts {
		contracts := []string{}
		for contract := range chainContracts {
			contracts = append(contracts, string(contract))
		}
		sort.Strings(contracts)
		contractWhitelists = append(contractWhitelists, BlockchainIDWhitelists{BlockchainID: string(blockchainID), Values: contracts})
	}
	sort.Slice(contractWhitelists, func(i, j int) bool {
		return contractWhitelists[i].BlockchainID < contractWhitelists[j].BlockchainID
	})

	for blockchainID, chainMethods := range a.Whitelists.Methods {
		methods := []string{}
		for method := range chainMethods {
			methods = append(methods, string(method))
		}
		sort.Strings(methods)
		methodWhitelists = append(methodWhitelists, BlockchainIDWhitelists{BlockchainID: string(blockchainID), Values: methods})
	}
	sort.Slice(methodWhitelists, func(i, j int) bool {
		return methodWhitelists[i].BlockchainID < methodWhitelists[j].BlockchainID
	})

	return &WhitelistsObject{
		AppWhitelists: [3]ApplicationWhitelists{
			{Type: WLOrigins, Values: origins},
			{Type: WLUserAgents, Values: userAgents},
			{Type: WLBlockchains, Values: blockchains},
		},
		ChainWhitelists: [2]ChainWhitelists{
			{Type: WLContracts, Values: contractWhitelists},
			{Type: WLMethods, Values: methodWhitelists},
		},
	}
}

func (a *PortalApp) Table() Table {
	return TablePortalApps
}

func (a *AAT) Table() Table {
	return TableAppAATs
}

func (a *Settings) Table() Table {
	return TableAppSettings
}

func (a *Whitelists) Table() Table {
	return TableAppWhitelists
}

func (a *AppNotification) Table() Table {
	return TableAppNotifications
}
