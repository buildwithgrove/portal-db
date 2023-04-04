package types

import (
	"time"
)

/* Legacy Enums */
type (
	AppStatus string
)

const (
	AwaitingFreetierFunds   AppStatus = "AWAITING_FREETIER_FUNDS"
	AwaitingFreetierStaking AppStatus = "AWAITING_FREETIER_STAKING"
	AwaitingFunds           AppStatus = "AWAITING_FUNDS"
	AwaitingFundsRemoval    AppStatus = "AWAITING_FUNDS_REMOVAL"
	AwaitingGracePeriod     AppStatus = "AWAITING_GRACE_PERIOD"
	AwaitingSlotFunds       AppStatus = "AWAITING_SLOT_FUNDS"
	AwaitingSlotStaking     AppStatus = "AWAITING_SLOT_STAKING"
	AwaitingStaking         AppStatus = "AWAITING_STAKING"
	AwaitingUnstaking       AppStatus = "AWAITING_UNSTAKING"
	Decomissioned           AppStatus = "DECOMISSIONED"
	InService               AppStatus = "IN_SERVICE"
	Orphaned                AppStatus = "ORPHANED"
	Ready                   AppStatus = "READY"
	Swappable               AppStatus = "SWAPPABLE"

	TestPlanV0   PayPlanType = "TEST_PLAN_V0"
	TestPlan10K  PayPlanType = "TEST_PLAN_10K"
	TestPlan90k  PayPlanType = "TEST_PLAN_90K"
	FreetierV0   PayPlanType = "FREETIER_V0"
	PayAsYouGoV0 PayPlanType = "PAY_AS_YOU_GO_V0"
	Enterprise   PayPlanType = "ENTERPRISE"
)

var (
	ValidAppStatuses = map[AppStatus]bool{
		"":                      true, // needed since it can be empty too
		AwaitingFreetierFunds:   true,
		AwaitingFreetierStaking: true,
		AwaitingFunds:           true,
		AwaitingFundsRemoval:    true,
		AwaitingGracePeriod:     true,
		AwaitingSlotFunds:       true,
		AwaitingSlotStaking:     true,
		AwaitingStaking:         true,
		AwaitingUnstaking:       true,
		Decomissioned:           true,
		InService:               true,
		Orphaned:                true,
		Ready:                   true,
		Swappable:               true,
	}
)

/* Legacy Structs */
type (
	LoadBalancerID        string
	LegacyUserPermissions struct {
		UserID        UserID                                     `json:"userID"`
		LoadBalancers map[LoadBalancerID]LoadBalancerPermissions `json:"loadBalancers"`
	}
	LoadBalancerPermissions struct {
		RoleName    RoleName      `json:"roleName"`
		Permissions []Permissions `json:"permissions"`
	}
	// LoadBalancer
	LoadBalancer struct {
		ID                string              `json:"id"`
		Name              string              `json:"name"`
		UserID            string              `json:"userID"`
		ApplicationIDs    []string            `json:"applicationIDs,omitempty"`
		RequestTimeout    int                 `json:"requestTimeout"`
		Gigastake         bool                `json:"gigastake"`
		GigastakeRedirect bool                `json:"gigastakeRedirect"`
		StickyOptions     StickyOptions       `json:"stickinessOptions"`
		Applications      []*Application      `json:"applications"`
		Integrations      AccountIntegrations `json:"integrations"`
		Users             []UserAccess        `json:"users"`
		CreatedAt         time.Time           `json:"createdAt"`
		UpdatedAt         time.Time           `json:"updatedAt"`
		AccountID         string              `json:"accountID"`
	}

	UserAccess struct {
		ID       string   `json:"id,omitempty"`
		UserID   string   `json:"userID"`
		RoleName RoleName `json:"roleName"`
		Email    string   `json:"email"`
		Accepted bool     `json:"accepted"`
	}
	// Application
	Application struct {
		ID                   string               `json:"id"`
		UserID               string               `json:"userID"`
		Name                 string               `json:"name"`
		ContactEmail         string               `json:"contactEmail"`
		Description          string               `json:"description"`
		Owner                string               `json:"owner"`
		URL                  string               `json:"url"`
		Dummy                bool                 `json:"dummy"`
		Status               AppStatus            `json:"status"`
		FirstDateSurpassed   time.Time            `json:"firstDateSurpassed"`
		GatewayAAT           GatewayAAT           `json:"gatewayAAT"`
		GatewaySettings      GatewaySettings      `json:"gatewaySettings"`
		Limit                AppLimit             `json:"limit"`
		NotificationSettings NotificationSettings `json:"notificationSettings"`
		CreatedAt            time.Time            `json:"createdAt"`
		UpdatedAt            time.Time            `json:"updatedAt"`
	}
	GatewayAAT struct {
		ID                   string `json:"id,omitempty"`
		Address              string `json:"address"`
		ApplicationPublicKey string `json:"applicationPublicKey"`
		ApplicationSignature string `json:"applicationSignature"`
		ClientPublicKey      string `json:"clientPublicKey"`
		PrivateKey           string `json:"privateKey"`
		Version              string `json:"version"`
	}
	GatewaySettings struct {
		ID                   string               `json:"id,omitempty"`
		SecretKey            string               `json:"secretKey"`
		SecretKeyRequired    bool                 `json:"secretKeyRequired"`
		WhitelistOrigins     []string             `json:"whitelistOrigins,omitempty"`
		WhitelistUserAgents  []string             `json:"whitelistUserAgents,omitempty"`
		WhitelistContracts   []WhitelistContracts `json:"whitelistContracts,omitempty"`
		WhitelistMethods     []WhitelistMethods   `json:"whitelistMethods,omitempty"`
		WhitelistBlockchains []string             `json:"whitelistBlockchains,omitempty"`
	}
	WhitelistContracts struct {
		ID           string   `json:"id,omitempty"`
		BlockchainID string   `json:"chainID"`
		Contracts    []string `json:"contracts"`
	}
	WhitelistMethods struct {
		ID           string   `json:"id,omitempty"`
		BlockchainID string   `json:"chainID"`
		Methods      []string `json:"methods"`
	}
	AppLimit struct {
		ID          string  `json:"id,omitempty"`
		PayPlan     PayPlan `json:"payPlan"`
		CustomLimit int     `json:"customLimit"`
	}
	PayPlan struct {
		Type  PayPlanType `json:"planType"`
		Limit int         `json:"dailyLimit"`
	}
	NotificationSettings struct {
		ID            string `json:"id,omitempty"`
		SignedUp      bool   `json:"signedUp"`
		Quarter       bool   `json:"quarter"`
		Half          bool   `json:"half"`
		ThreeQuarters bool   `json:"threeQuarters"`
		Full          bool   `json:"full"`
	}
	// Blockchain
	Blockchain struct {
		ID                string           `json:"id"`
		Altruist          string           `json:"altruist"`
		Blockchain        string           `json:"blockchain"`
		ChainID           string           `json:"chainID"`
		ChainIDCheck      string           `json:"chainIDCheck"`
		Description       string           `json:"description"`
		EnforceResult     string           `json:"enforceResult"`
		Path              string           `json:"path"`
		Ticker            string           `json:"ticker"`
		BlockchainAliases []string         `json:"blockchainAliases"`
		LogLimitBlocks    int              `json:"logLimitBlocks"`
		RequestTimeout    int              `json:"requestTimeout"`
		Active            bool             `json:"active"`
		Redirects         []Redirect       `json:"redirects"`
		SyncCheckOptions  SyncCheckOptions `json:"syncCheckOptions"`
		CreatedAt         time.Time        `json:"createdAt"`
		UpdatedAt         time.Time        `json:"updatedAt"`
	}
	Redirect struct {
		ChainID        string    `json:"chainID,omitempty"`
		Alias          string    `json:"alias"`
		Domain         string    `json:"domain"`
		LoadBalancerID string    `json:"loadBalancerID"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}
	SyncCheckOptions struct {
		BlockchainID string `json:"chainID,omitempty"`
		Body         string `json:"body"`
		ResultKey    string `json:"resultKey"`
		Allowance    int    `json:"allowance"`
	}

	/* Update structs */
	UpdateLoadBalancer struct {
		Name          string               `json:"name,omitempty"`
		StickyOptions *UpdateStickyOptions `json:"stickinessOptions,omitempty"`
		Remove        bool                 `json:"remove,omitempty"`
	}
	UpdateStickyOptions struct {
		ID            string   `json:"id,omitempty"`
		Duration      string   `json:"duration"`
		StickyOrigins []string `json:"stickyOrigins"`
		StickyMax     int      `json:"stickyMax"`
		Stickiness    *bool    `json:"stickiness"`
	}
	UpdateUserAccess struct {
		ID       string   `json:"id,omitempty"`
		UserID   string   `json:"userID"`
		Email    string   `json:"email"`
		RoleName RoleName `json:"roleName"`
	}

	UpdateApplication struct {
		Name                 string                      `json:"name,omitempty"`
		Status               AppStatus                   `json:"status,omitempty"`
		FirstDateSurpassed   time.Time                   `json:"firstDateSurpassed,omitempty"`
		GatewaySettings      *UpdateGatewaySettings      `json:"gatewaySettings,omitempty"`
		NotificationSettings *UpdateNotificationSettings `json:"notificationSettings,omitempty"`
		Limit                *AppLimit                   `json:"appLimit,omitempty"`
		Remove               bool                        `json:"remove,omitempty"`
	}
	UpdateGatewaySettings struct {
		ID                   string               `json:"id,omitempty"`
		SecretKey            string               `json:"secretKey"`
		SecretKeyRequired    *bool                `json:"secretKeyRequired"`
		WhitelistOrigins     []string             `json:"whitelistOrigins,omitempty"`
		WhitelistUserAgents  []string             `json:"whitelistUserAgents,omitempty"`
		WhitelistContracts   []WhitelistContracts `json:"whitelistContracts,omitempty"`
		WhitelistMethods     []WhitelistMethods   `json:"whitelistMethods,omitempty"`
		WhitelistBlockchains []string             `json:"whitelistBlockchains,omitempty"`
	}
	UpdateNotificationSettings struct {
		ID            string `json:"id,omitempty"`
		SignedUp      *bool  `json:"signedUp"`
		Quarter       *bool  `json:"quarter"`
		Half          *bool  `json:"half"`
		ThreeQuarters *bool  `json:"threeQuarters"`
		Full          *bool  `json:"full"`
	}

	UpdateBlockchain struct {
		Altruist          string   `json:"altruist,omitempty"`
		Blockchain        string   `json:"blockchain,omitempty"`
		ChainIDCheck      string   `json:"chainIDCheck,omitempty"`
		Description       string   `json:"description,omitempty"`
		EnforceResult     string   `json:"enforceResult,omitempty"`
		Path              string   `json:"path,omitempty"`
		Ticker            string   `json:"ticker,omitempty"`
		BlockchainAliases []string `json:"blockchainAliases,omitempty"`
		LogLimitBlocks    int      `json:"logLimitBlocks,omitempty"`
		RequestTimeout    int      `json:"requestTimeout,omitempty"`

		Body      string `json:"body,omitempty"`
		ResultKey string `json:"resultKey,omitempty"`
		Allowance *int   `json:"allowance,omitempty"`

		UpdatedAt time.Time `json:"updatedAt"`
	}
)

/* Legacy Methods */
func (a *Application) DailyLimit() int {
	if a.Limit.PayPlan.Type == Enterprise {
		return a.Limit.CustomLimit
	}

	return a.Limit.PayPlan.Limit
}
