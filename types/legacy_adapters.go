package types

import (
	"sort"
	"time"
)

/* V2 Struct to Legacy Struct Adaptors */
func (a *PortalApp) ConvertToLegacyLoadBalancer() LoadBalancer {
	var users []UserAccess
	for _, accountUser := range a.Account.Users {
		var userID string
		// For now we can assume each ProviderUserIDs map only has one entry
		for _, id := range accountUser.ProviderUserIDs {
			userID = id
		}
		users = append(users, UserAccess{
			UserID:   userID,
			RoleName: accountUser.RoleName,
			Email:    string(accountUser.Email),
			Accepted: accountUser.Accepted,
		})
	}
	sortUsersByRole(users)

	return LoadBalancer{
		ID:           string(a.ID),
		Name:         a.Name,
		UserID:       a.LegacyUserID(),
		Gigastake:    a.Gigastake,
		Applications: a.ConvertToLegacyApplications(),
		Users:        users,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,

		RequestTimeout:    int(a.LegacyFields.RequestTimeout),
		GigastakeRedirect: a.LegacyFields.GigastakeRedirect,
		StickyOptions:     a.LegacyFields.StickyOptions,
	}
}

func (a *PortalApp) ConvertToLegacyApplications() []*Application {
	appDetails := &Application{
		UserID: a.LegacyUserID(),
		Name:   a.Name,
		GatewayAAT: GatewayAAT{
			Address:              a.AAT.Address,
			ApplicationPublicKey: a.AAT.PublicKey,
			ApplicationSignature: a.AAT.Signature,
			ClientPublicKey:      a.AAT.ClientPublicKey,
			PrivateKey:           a.AAT.PrivateKey,
			Version:              a.AAT.Version,
		},
		GatewaySettings: a.ConvertToLegacyGatewaySettings(),
		Limit: AppLimit{
			Plan: PayPlan{
				Type:  a.Account.Plan.Type,
				Limit: int(a.LegacyDailyLimit()),
			},
			CustomLimit: int(a.LegacyFields.CustomLimit),
		},
		NotificationSettings: NotificationSettings{
			SignedUp:      a.Notifications[NotificationTypeEmail].Events[NotificationEventSignedUp],
			Quarter:       a.Notifications[NotificationTypeEmail].Events[NotificationEventQuarter],
			Half:          a.Notifications[NotificationTypeEmail].Events[NotificationEventHalf],
			ThreeQuarters: a.Notifications[NotificationTypeEmail].Events[NotificationEventThreeQuarters],
			Full:          a.Notifications[NotificationTypeEmail].Events[NotificationEventFull],
		},
		FirstDateSurpassed: a.LegacyFields.FirstDateSurpassed,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}

	var applications []*Application
	for _, id := range a.LegacyFields.ApplicationIDs {
		application := appDetails
		application.ID = id

		applications = append(applications, application)
	}

	return applications
}

func (a *PortalApp) ConvertToLegacyGatewaySettings() GatewaySettings {
	gatewaySettings := GatewaySettings{
		SecretKey:         a.Settings.SecretKey,
		SecretKeyRequired: a.Settings.SecretKeyRequired,
	}

	for origin := range a.Whitelists.Origins {
		gatewaySettings.WhitelistOrigins = append(gatewaySettings.WhitelistOrigins, string(origin))
	}
	sort.Strings(gatewaySettings.WhitelistOrigins)

	for userAgent := range a.Whitelists.UserAgents {
		gatewaySettings.WhitelistUserAgents = append(gatewaySettings.WhitelistUserAgents, string(userAgent))
	}
	sort.Strings(gatewaySettings.WhitelistUserAgents)

	for chainID := range a.Whitelists.Blockchains {
		gatewaySettings.WhitelistBlockchains = append(gatewaySettings.WhitelistBlockchains, string(chainID))
	}
	sort.Strings(gatewaySettings.WhitelistBlockchains)

	for chainID, contracts := range a.Whitelists.Contracts {
		var contractList []string
		for contract := range contracts {
			contractList = append(contractList, string(contract))
		}
		sort.Strings(contractList)
		gatewaySettings.WhitelistContracts = append(gatewaySettings.WhitelistContracts, WhitelistContracts{
			ChainID: string(chainID), Contracts: contractList},
		)
	}
	sort.Slice(gatewaySettings.WhitelistContracts, func(i, j int) bool {
		return gatewaySettings.WhitelistContracts[i].ChainID < gatewaySettings.WhitelistContracts[j].ChainID
	})

	for chainID, methods := range a.Whitelists.Methods {
		var methodList []string
		for method := range methods {
			methodList = append(methodList, string(method))
		}
		sort.Strings(methodList)
		gatewaySettings.WhitelistMethods = append(gatewaySettings.WhitelistMethods, WhitelistMethods{
			ChainID: string(chainID), Methods: methodList},
		)
	}
	sort.Slice(gatewaySettings.WhitelistMethods, func(i, j int) bool {
		return gatewaySettings.WhitelistMethods[i].ChainID < gatewaySettings.WhitelistMethods[j].ChainID
	})

	return gatewaySettings
}

func (c *Chain) ConvertToLegacyBlockchain() Blockchain {
	var redirects []Redirect
	for _, chainRedirect := range c.Redirects {
		redirects = append(redirects, Redirect{
			Alias:          chainRedirect.Alias,
			Domain:         chainRedirect.Domain,
			LoadBalancerID: chainRedirect.ProtocolAppID,
		})
	}

	return Blockchain{
		ID:                string(c.ID),
		Blockchain:        c.Blockchain,
		ChainID:           c.ChainID,
		ChainIDCheck:      c.Checks[ChainCheckTypeChain].Payload,
		Description:       c.Description,
		EnforceResult:     c.EnforceResult,
		Path:              c.Path,
		Ticker:            c.Ticker,
		BlockchainAliases: c.BlockchainAliases,
		LogLimitBlocks:    c.LogLimitBlocks,
		RequestTimeout:    c.RequestTimeout,
		Active:            c.Active,
		Altruist:          c.Altruists[0].URL,
		SyncCheckOptions: SyncCheckOptions{
			Body:      c.Checks[ChainCheckTypeSync].Payload,
			ResultKey: c.Checks[ChainCheckTypeSync].ResultKey,
			Allowance: c.Checks[ChainCheckTypeSync].Allowance,
		},
		Redirects: redirects,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (c *Plan) ConvertToLegacyPayPlan() PayPlan {
	return PayPlan{
		Type:  c.Type,
		Limit: int(c.LegacyDailyLimit),
	}
}

/* Legacy Struct to V2 Struct Adaptors */

// Creates the struct with all fields needed to create a new PortalApp
// LoadBalancer must be sent to PHD containing its Application already defined inside PUB
// This way the PortalApp can be created in only one operation (no PHD client changes needed)
func (lb *LoadBalancer) ConvertToV2PortalApp(accountID AccountID) PortalApp {
	app := lb.Applications[0]
	owner := lb.Users[0]

	return PortalApp{
		ID:        "", // generate ID inside postgresdriver (same ID as LoadBalancer)
		Name:      lb.Name,
		Gigastake: lb.Gigastake,
		AccountID: accountID,
		AAT: AAT{
			Address:         app.GatewayAAT.Address,
			PublicKey:       app.GatewayAAT.ApplicationPublicKey,
			ClientPublicKey: app.GatewayAAT.ClientPublicKey,
			PrivateKey:      app.GatewayAAT.PrivateKey,
			Signature:       app.GatewayAAT.ApplicationSignature,
			Version:         app.GatewayAAT.Version,
		},
		Settings: Settings{
			Environment: EnvironmentProduction,
			SecretKey:   app.GatewaySettings.SecretKey,
		},
		Notifications: map[NotificationType]AppNotification{
			NotificationTypeEmail: {
				Active:      true,
				Destination: owner.Email,
				Events: map[NotificationEvent]bool{
					NotificationEventSignedUp:      app.NotificationSettings.SignedUp,
					NotificationEventQuarter:       app.NotificationSettings.Quarter,
					NotificationEventHalf:          app.NotificationSettings.Half,
					NotificationEventThreeQuarters: app.NotificationSettings.ThreeQuarters,
					NotificationEventFull:          app.NotificationSettings.Full,
				},
			},
		},

		LegacyFields: LegacyFields{
			ApplicationIDs:    []string{}, // generate ID inside postgres driver (same ID as Application)
			RequestTimeout:    int32(lb.RequestTimeout),
			GigastakeRedirect: lb.GigastakeRedirect,
			StickyOptions:     lb.StickyOptions,
		},
	}
}

// Converts the existing UpdateApplication struct to a new one that updates all relevant fields
// UpdateLoadBalancer struct is redundant since it previously only updated the Name field
// Must use /load_balancer update endpoint as the ID for PortalApps is the former LB ID
func (u *UpdateApplication) ConvertToV2UpdatePortalApp(loadBalancerID string) UpdatePortalApp {
	var (
		settings                             *UpdateAppSettings
		notifications                        []UpdateAppNotifications
		whitelists                           *WhitelistsObject
		contractWhitelists, methodWhitelists []BlockchainIDWhitelists
	)

	if u.NotificationSettings != nil {
		notificationEvents := []NotificationEvent{}

		if u.NotificationSettings.SignedUp != nil && *u.NotificationSettings.SignedUp {
			notificationEvents = append(notificationEvents, NotificationEventSignedUp)
		}
		if u.NotificationSettings.Quarter != nil && *u.NotificationSettings.Quarter {
			notificationEvents = append(notificationEvents, NotificationEventQuarter)
		}
		if u.NotificationSettings.Half != nil && *u.NotificationSettings.Half {
			notificationEvents = append(notificationEvents, NotificationEventHalf)
		}
		if u.NotificationSettings.ThreeQuarters != nil && *u.NotificationSettings.ThreeQuarters {
			notificationEvents = append(notificationEvents, NotificationEventThreeQuarters)
		}
		if u.NotificationSettings.Full != nil && *u.NotificationSettings.Full {
			notificationEvents = append(notificationEvents, NotificationEventFull)
		}

		notifications = []UpdateAppNotifications{
			{NotificationType: NotificationTypeEmail, Events: notificationEvents},
		}
	}

	if u.GatewaySettings != nil {
		settings = &UpdateAppSettings{SecretKey: u.GatewaySettings.SecretKey}
		if u.GatewaySettings.SecretKeyRequired != nil {
			settings.SecretKeyRequired = *u.GatewaySettings.SecretKeyRequired
		}

		for _, chainContracts := range u.GatewaySettings.WhitelistContracts {
			contracts := []string{}
			contracts = append(contracts, chainContracts.Contracts...)
			contractWhitelists = append(contractWhitelists, BlockchainIDWhitelists{
				ChainID: chainContracts.ChainID, Values: contracts,
			})
		}
		for _, chainMethods := range u.GatewaySettings.WhitelistMethods {
			methods := []string{}
			methods = append(methods, chainMethods.Methods...)
			methodWhitelists = append(methodWhitelists, BlockchainIDWhitelists{
				ChainID: chainMethods.ChainID, Values: methods,
			})
		}

		whitelists = &WhitelistsObject{
			AppWhitelists: [3]ApplicationWhitelists{
				{Type: WhitelistTypeOrigins, Values: u.GatewaySettings.WhitelistOrigins},
				{Type: WhitelistTypeUserAgents, Values: u.GatewaySettings.WhitelistUserAgents},
				{Type: WhitelistTypeBlockchains, Values: u.GatewaySettings.WhitelistBlockchains},
			},
			ChainWhitelists: [2]ChainWhitelists{
				{Type: WhitelistTypeContracts, Values: contractWhitelists},
				{Type: WhitelistTypeMethods, Values: methodWhitelists},
			},
		}
	}

	return UpdatePortalApp{
		AppID: PortalAppID(loadBalancerID), Name: u.Name,
		Settings: settings, Notifications: notifications, Whitelists: whitelists,
	}
}

func (u *UserAccess) ConvertToV2AccountUserAccess() AccountUserAccess {
	authType := getAuthType(u.UserID)
	return AccountUserAccess{
		UserID:          UserID(0), // autogenerated by DB
		Email:           Email(u.Email),
		RoleName:        u.RoleName,
		Accepted:        u.Accepted,
		ProviderUserIDs: map[AuthType]string{authType: u.UserID},
	}
}

func (u *UpdateUserAccess) ConvertToV2UpdateAccountUserAccess(userID int32, accepted bool) AccountUserAccess {
	authType := getAuthType(u.UserID)
	return AccountUserAccess{
		UserID:          UserID(userID),
		Email:           Email(u.Email),
		RoleName:        u.RoleName,
		Accepted:        accepted,
		ProviderUserIDs: map[AuthType]string{authType: u.UserID},
	}
}

func getAuthType(userID string) AuthType {
	switch {
	case len(userID) == 24:
		return AuthTypeAuth0Username
	default:
		return AuthTypeAuth0Github
	}
}

func (b *Blockchain) ConvertToV2Chain() Chain {
	checks := map[ChainCheckType]Check{
		ChainCheckTypeSync: {
			Payload:   b.SyncCheckOptions.Body,
			ResultKey: b.SyncCheckOptions.ResultKey,
			Allowance: b.SyncCheckOptions.Allowance,
		},
	}
	if b.ChainIDCheck != "" {
		checks[ChainCheckTypeChain] = Check{Payload: b.ChainIDCheck}
	}

	return Chain{
		ID:                ChainID(b.ID),
		Blockchain:        b.Blockchain,
		ChainID:           b.ChainID,
		Description:       b.Description,
		EnforceResult:     b.EnforceResult,
		Path:              b.Path,
		Ticker:            b.Ticker,
		BlockchainAliases: b.BlockchainAliases,
		LogLimitBlocks:    b.LogLimitBlocks,
		RequestTimeout:    b.RequestTimeout,
		Active:            b.Active,
		Altruists:         []Altruist{{URL: b.Altruist}},
		Checks:            checks,
	}
}

func (u *UpdateBlockchain) ConvertToV2UpdateChain() UpdateChain {
	checks := map[ChainCheckType]UpdateCheck{
		ChainCheckTypeSync: {Payload: u.Body, ResultKey: u.ResultKey, Allowance: u.Allowance},
	}
	if u.ChainIDCheck != "" {
		checks[ChainCheckTypeChain] = UpdateCheck{Payload: u.ChainIDCheck}
	}

	return UpdateChain{
		Blockchain:        u.Blockchain,
		Description:       u.Description,
		EnforceResult:     u.EnforceResult,
		Path:              u.Path,
		Ticker:            u.Ticker,
		BlockchainAliases: u.BlockchainAliases,
		LogLimitBlocks:    u.LogLimitBlocks,
		RequestTimeout:    u.RequestTimeout,
		Altruists:         []Altruist{{URL: u.Altruist}},
		Checks:            checks,
	}
}

func (r *Redirect) ConvertToV2Redirect() GigastakeRedirect {
	return GigastakeRedirect{
		Alias:         r.Alias,
		Domain:        r.Domain,
		ProtocolAppID: r.LoadBalancerID,
	}
}

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
	// LoadBalancer
	LoadBalancer struct {
		ID                string         `json:"id"`
		Name              string         `json:"name"`
		UserID            string         `json:"userID"`
		ApplicationIDs    []string       `json:"applicationIDs,omitempty"`
		RequestTimeout    int            `json:"requestTimeout"`
		Gigastake         bool           `json:"gigastake"`
		GigastakeRedirect bool           `json:"gigastakeRedirect"`
		StickyOptions     StickyOptions  `json:"stickinessOptions"`
		Applications      []*Application `json:"applications"`
		Users             []UserAccess   `json:"users"`
		CreatedAt         time.Time      `json:"createdAt"`
		UpdatedAt         time.Time      `json:"updatedAt"`
	}
	StickyOptions struct {
		ID            string   `json:"id,omitempty"`
		Duration      string   `json:"duration"`
		StickyOrigins []string `json:"stickyOrigins"`
		StickyMax     int      `json:"stickyMax"`
		Stickiness    bool     `json:"stickiness"`
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
		ID        string   `json:"id,omitempty"`
		ChainID   string   `json:"chainID"`
		Contracts []string `json:"contracts"`
	}
	WhitelistMethods struct {
		ID      string   `json:"id,omitempty"`
		ChainID string   `json:"chainID"`
		Methods []string `json:"methods"`
	}
	AppLimit struct {
		ID          string  `json:"id,omitempty"`
		Plan        PayPlan `json:"payPlan"`
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
		ChainID   string `json:"chainID,omitempty"`
		Body      string `json:"body"`
		ResultKey string `json:"resultKey"`
		Allowance int    `json:"allowance"`
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
	if a.Limit.Plan.Type == Enterprise {
		return a.Limit.CustomLimit
	}

	return a.Limit.Plan.Limit
}

/* Utils */
func sortUsersByRole(users []UserAccess) {
	roleWeight := map[RoleName]int{RoleOwner: 0, RoleAdmin: 1, RoleMember: 2}
	sort.Slice(users, func(i, j int) bool {
		return roleWeight[users[i].RoleName] < roleWeight[users[j].RoleName]
	})
}
