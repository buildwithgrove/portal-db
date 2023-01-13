package types

import (
	"errors"
	"time"
)

var (
	ErrNoFieldsToUpdate               = errors.New("no fields to update")
	ErrInvalidAppStatus               = errors.New("invalid app status")
	ErrInvalidPayPlanType             = errors.New("invalid pay plan type")
	ErrNotEnterprisePlan              = errors.New("custom limits may only be set on enterprise plans")
	ErrEnterprisePlanNeedsCustomLimit = errors.New("enterprise plans must have a custom limit set")
)

type (
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
		ID                   string              `json:"id,omitempty"`
		SecretKey            string              `json:"secretKey"`
		SecretKeyRequired    bool                `json:"secretKeyRequired"`
		WhitelistOrigins     []string            `json:"whitelistOrigins,omitempty"`
		WhitelistUserAgents  []string            `json:"whitelistUserAgents,omitempty"`
		WhitelistContracts   []WhitelistContract `json:"whitelistContracts,omitempty"`
		WhitelistMethods     []WhitelistMethod   `json:"whitelistMethods,omitempty"`
		WhitelistBlockchains []string            `json:"whitelistBlockchains,omitempty"`
	}
	WhitelistContract struct {
		ID           string   `json:"id,omitempty"`
		BlockchainID string   `json:"blockchainID"`
		Contracts    []string `json:"contracts"`
	}
	WhitelistMethod struct {
		ID           string   `json:"id,omitempty"`
		BlockchainID string   `json:"blockchainID"`
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
	/* Update structs */
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
		ID                   string              `json:"id,omitempty"`
		SecretKey            string              `json:"secretKey"`
		SecretKeyRequired    *bool               `json:"secretKeyRequired"`
		WhitelistOrigins     []string            `json:"whitelistOrigins,omitempty"`
		WhitelistUserAgents  []string            `json:"whitelistUserAgents,omitempty"`
		WhitelistContracts   []WhitelistContract `json:"whitelistContracts,omitempty"`
		WhitelistMethods     []WhitelistMethod   `json:"whitelistMethods,omitempty"`
		WhitelistBlockchains []string            `json:"whitelistBlockchains,omitempty"`
	}
	UpdateFirstDateSurpassed struct {
		ApplicationIDs     []string  `json:"applicationIDs"`
		FirstDateSurpassed time.Time `json:"firstDateSurpassed"`
	}
	UpdateNotificationSettings struct {
		ID            string `json:"id,omitempty"`
		SignedUp      *bool  `json:"signedUp"`
		Quarter       *bool  `json:"quarter"`
		Half          *bool  `json:"half"`
		ThreeQuarters *bool  `json:"threeQuarters"`
		Full          *bool  `json:"full"`
	}

	AppStatus   string
	PayPlanType string
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

	ValidPayPlanTypes = map[PayPlanType]bool{
		"":           true, // needs to be allowed while the change for all apps to have plans is done
		TestPlanV0:   true,
		TestPlan10K:  true,
		TestPlan90k:  true,
		FreetierV0:   true,
		PayAsYouGoV0: true,
		Enterprise:   true,
	}
)

func (a *Application) DailyLimit() int {
	if a.Limit.PayPlan.Type == Enterprise {
		return a.Limit.CustomLimit
	}

	return a.Limit.PayPlan.Limit
}

func (a *Application) Validate() error {
	if !ValidAppStatuses[a.Status] {
		return ErrInvalidAppStatus
	}

	if !ValidPayPlanTypes[a.Limit.PayPlan.Type] {
		return ErrInvalidPayPlanType
	}

	if a.Limit.PayPlan.Type != Enterprise && a.Limit.CustomLimit != 0 {
		return ErrNotEnterprisePlan
	}
	return nil
}

func (u *UpdateApplication) Validate() error {
	if u == nil {
		return ErrNoFieldsToUpdate
	}
	if !ValidAppStatuses[u.Status] {
		return ErrInvalidAppStatus
	}
	if u.Limit != nil && !ValidPayPlanTypes[u.Limit.PayPlan.Type] {
		return ErrInvalidPayPlanType
	}
	if u.Limit != nil && u.Limit.PayPlan.Type != Enterprise && u.Limit.CustomLimit != 0 {
		return ErrNotEnterprisePlan
	}
	if u.Limit != nil && u.Limit.PayPlan.Type == Enterprise && u.Limit.CustomLimit == 0 {
		return ErrEnterprisePlanNeedsCustomLimit
	}
	return nil
}

func (p *PayPlan) Validate() error {
	if !ValidPayPlanTypes[p.Type] {
		return ErrInvalidPayPlanType
	}

	return nil
}
