package types

import (
	"sort"
	"time"
)

/* V2 Struct to Legacy Struct Adaptors */
func (a *PortalApp) ConvertToLegacyLoadBalancer() LoadBalancer {
	var users []UserAccess
	for userID, accountUser := range a.Account.Users {
		users = append(users, UserAccess{
			UserID:   string(userID),
			RoleName: accountUser.RoleName,
			Email:    string(accountUser.User.Email),
			Accepted: accountUser.Accepted,
		})
	}
	sortUsersByRole(users)

	application := a.ConvertToLegacyApplication()

	return LoadBalancer{
		ID:           a.ID,
		Name:         a.Name,
		UserID:       string(a.UserID()),
		Gigastake:    a.Gigastake,
		Applications: []*Application{&application},
		Users:        users,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,

		RequestTimeout:    a.LegacyFields.RequestTimeout,
		GigastakeRedirect: a.LegacyFields.GigastakeRedirect,
		StickyOptions:     a.LegacyFields.StickyOptions,
	}
}

func (a *PortalApp) ConvertToLegacyApplication() Application {
	return Application{
		ID:     a.LegacyFields.PortalAppID,
		UserID: string(a.UserID()),
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
				Limit: a.LegacyDailyLimit(),
			},
			CustomLimit: a.LegacyFields.CustomLimit,
		},
		NotificationSettings: NotificationSettings{
			SignedUp:      a.Notifications[NotificationEmail].Events[EventSignedUp],
			Quarter:       a.Notifications[NotificationEmail].Events[EventQuarter],
			Half:          a.Notifications[NotificationEmail].Events[EventHalf],
			ThreeQuarters: a.Notifications[NotificationEmail].Events[EventThreeQuarters],
			Full:          a.Notifications[NotificationEmail].Events[EventFull],
		},
		FirstDateSurpassed: a.LegacyFields.FirstDateSurpassed,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
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
		ChainIDCheck:      c.Checks[CheckChain].Payload,
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
			Body:      c.Checks[CheckSync].Payload,
			ResultKey: c.Checks[CheckSync].ResultKey,
			Allowance: c.Checks[CheckSync].Allowance,
		},
		Redirects: redirects,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (c *Plan) ConvertToLegacyPayPlan() PayPlan {
	return PayPlan{
		Type:  c.Type,
		Limit: c.LegacyDailyLimit,
	}
}

/* Legacy Struct to V2 Struct Adaptors */

// Creates the struct with all fields needed to create a new PortalApp
// LoadBalancer must be sent to PHD containing its Application already defined inside PUB
// This way the PortalApp can be created in only one operation (no PHD client changes needed)
func (lb *LoadBalancer) ConvertToV2PortalApp() PortalApp {
	app := lb.Applications[0]
	owner := lb.Users[0]

	return PortalApp{
		ID:        "", // generate ID inside postgresdriver (same ID as LoadBalancer)
		Name:      lb.Name,
		Gigastake: lb.Gigastake,
		Account: &Account{
			Plan: Plan{Type: app.Limit.Plan.Type},
			Users: map[UserID]AccountUserAccess{
				UserID(owner.UserID): {
					User:     User{ID: UserID(owner.UserID), Email: Email(owner.Email), AuthProvider: ProviderAuth0},
					RoleName: RoleOwner,
					Accepted: true,
				},
			},
		},
		AAT: AAT{
			Address:         app.GatewayAAT.Address,
			PublicKey:       app.GatewayAAT.ApplicationPublicKey,
			ClientPublicKey: app.GatewayAAT.ClientPublicKey,
			PrivateKey:      app.GatewayAAT.PrivateKey,
			Signature:       app.GatewayAAT.ApplicationSignature,
			Version:         app.GatewayAAT.Version,
		},
		Settings: Settings{
			Environment: EnvProduction,
			SecretKey:   app.GatewaySettings.SecretKey,
		},
		Notifications: map[NotificationType]AppNotification{
			NotificationEmail: {
				Active:      true,
				Destination: owner.Email,
				Events: map[NotificationEvent]bool{
					EventSignedUp:      app.NotificationSettings.SignedUp,
					EventQuarter:       app.NotificationSettings.Quarter,
					EventHalf:          app.NotificationSettings.Half,
					EventThreeQuarters: app.NotificationSettings.ThreeQuarters,
					EventFull:          app.NotificationSettings.Full,
				},
			},
		},

		LegacyFields: LegacyFields{
			PortalAppID:       "", // generate ID inside postgres driver (same ID as Application)
			RequestTimeout:    lb.RequestTimeout,
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
		notifications                        *UpdateAppNotifications
		whitelists                           *WhitelistsObject
		contractWhitelists, methodWhitelists []BlockchainIDWhitelists
	)

	if u.NotificationSettings != nil {
		notifications = &UpdateAppNotifications{
			NotificationType: NotificationEmail,
			Events: map[NotificationEvent]bool{
				EventSignedUp:      *u.NotificationSettings.SignedUp,
				EventQuarter:       *u.NotificationSettings.Quarter,
				EventHalf:          *u.NotificationSettings.Half,
				EventThreeQuarters: *u.NotificationSettings.ThreeQuarters,
				EventFull:          *u.NotificationSettings.Full,
			},
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
				{Type: WLOrigins, Values: u.GatewaySettings.WhitelistOrigins},
				{Type: WLUserAgents, Values: u.GatewaySettings.WhitelistUserAgents},
				{Type: WLBlockchains, Values: u.GatewaySettings.WhitelistBlockchains},
			},
			ChainWhitelists: [2]ChainWhitelists{
				{Type: WLContracts, Values: contractWhitelists},
				{Type: WLMethods, Values: methodWhitelists},
			},
		}
	}

	return UpdatePortalApp{
		AppID: PortalAppID(loadBalancerID), Name: u.Name,
		Settings: settings, Notifications: notifications, Whitelists: whitelists,
	}
}

func (u *UserAccess) ConvertToV2AccountUserAccess() AccountUserAccess {
	return AccountUserAccess{
		User:     User{ID: UserID(u.UserID), Email: Email(u.Email), AuthProvider: ProviderAuth0},
		RoleName: u.RoleName, Accepted: u.Accepted,
	}
}

func (u *UpdateUserAccess) ConvertToV2UpdateAccountUserAccess(accepted bool) AccountUserAccess {
	return AccountUserAccess{
		User:     User{ID: UserID(u.UserID), Email: Email(u.Email), AuthProvider: ProviderAuth0},
		RoleName: u.RoleName, Accepted: accepted,
	}
}

func (b *Blockchain) ConvertToV2Chain() Chain {
	checks := map[ChainCheckType]Check{
		CheckSync: {
			Payload:   b.SyncCheckOptions.Body,
			ResultKey: b.SyncCheckOptions.ResultKey,
			Allowance: b.SyncCheckOptions.Allowance,
		},
	}
	if b.ChainIDCheck != "" {
		checks[CheckChain] = Check{Payload: b.ChainIDCheck}
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
		CheckSync: {Payload: u.Body, ResultKey: u.ResultKey, Allowance: u.Allowance},
	}
	if u.ChainIDCheck != "" {
		checks[CheckChain] = UpdateCheck{Payload: u.ChainIDCheck}
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
