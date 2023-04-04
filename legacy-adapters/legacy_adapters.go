package legacyadapters

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/pokt-foundation/portal-db/v2/types"
)

/* V2 Struct to Legacy Struct Adaptors */
func ConvertToLegacyLoadBalancer(a types.PortalApp, account types.Account) types.LoadBalancer {
	var users []types.UserAccess

	for _, accountUser := range account.Users {
		users = append(users, types.UserAccess{
			UserID:   string(accountUser.UserID),
			RoleName: accountUser.RoleName,
			Email:    string(accountUser.Email),
			Accepted: accountUser.Accepted,
		})
	}

	sortUsersByRole(users)

	userID := account.LegacyUserID()
	legacyDailyLimit := account.Plan.LegacyDailyLimit

	return types.LoadBalancer{
		ID:                string(a.ID),
		Name:              a.Name,
		UserID:            userID,
		Applications:      ConvertToLegacyApplications(a, userID, account.Plan.Type, legacyDailyLimit),
		Users:             users,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
		Gigastake:         a.Gigastake,
		RequestTimeout:    int(a.LegacyFields.RequestTimeout),
		GigastakeRedirect: a.LegacyFields.GigastakeRedirect,
		StickyOptions:     a.LegacyFields.StickyOptions,
	}
}

func ConvertToLegacyApplications(a types.PortalApp, userID string, planType types.PayPlanType, dailyLimit int32) []*types.Application {
	baseApp := types.Application{
		UserID:          userID,
		Name:            a.Name,
		GatewaySettings: ConvertToLegacyGatewaySettings(a),
		Limit: types.AppLimit{
			PayPlan:     types.PayPlan{Type: planType, Limit: int(dailyLimit)},
			CustomLimit: int(a.LegacyFields.CustomLimit),
		},
		NotificationSettings: types.NotificationSettings{
			SignedUp:      a.Notifications[types.NotificationTypeEmail].Events[types.NotificationEventSignedUp],
			Quarter:       a.Notifications[types.NotificationTypeEmail].Events[types.NotificationEventQuarter],
			Half:          a.Notifications[types.NotificationTypeEmail].Events[types.NotificationEventHalf],
			ThreeQuarters: a.Notifications[types.NotificationTypeEmail].Events[types.NotificationEventThreeQuarters],
			Full:          a.Notifications[types.NotificationTypeEmail].Events[types.NotificationEventFull],
		},
		FirstDateSurpassed: a.LegacyFields.FirstDateSurpassed,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}

	var legacyApps []*types.Application

	for appID, aat := range a.AATs {
		app := baseApp

		app.ID = string(appID)

		app.GatewayAAT = types.GatewayAAT{
			Address:              aat.Address,
			ApplicationPublicKey: aat.PublicKey,
			ApplicationSignature: aat.Signature,
			ClientPublicKey:      aat.ClientPublicKey,
			PrivateKey:           aat.PrivateKey,
			Version:              aat.Version,
		}

		legacyApps = append(legacyApps, &app)
	}

	return legacyApps
}

func ConvertToLegacyGatewaySettings(a types.PortalApp) types.GatewaySettings {
	gatewaySettings := types.GatewaySettings{
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
		gatewaySettings.WhitelistContracts = append(gatewaySettings.WhitelistContracts, types.WhitelistContracts{
			BlockchainID: string(chainID), Contracts: contractList},
		)
	}
	sort.Slice(gatewaySettings.WhitelistContracts, func(i, j int) bool {
		return gatewaySettings.WhitelistContracts[i].BlockchainID < gatewaySettings.WhitelistContracts[j].BlockchainID
	})

	for chainID, methods := range a.Whitelists.Methods {
		var methodList []string
		for method := range methods {
			methodList = append(methodList, string(method))
		}
		sort.Strings(methodList)
		gatewaySettings.WhitelistMethods = append(gatewaySettings.WhitelistMethods, types.WhitelistMethods{
			BlockchainID: string(chainID), Methods: methodList},
		)
	}
	sort.Slice(gatewaySettings.WhitelistMethods, func(i, j int) bool {
		return gatewaySettings.WhitelistMethods[i].BlockchainID < gatewaySettings.WhitelistMethods[j].BlockchainID
	})

	return gatewaySettings
}

func ConvertToLegacyBlockchain(c types.Chain) types.Blockchain {
	var redirects []types.Redirect
	for _, chainRedirect := range c.Redirects {
		redirects = append(redirects, types.Redirect{
			Alias:          chainRedirect.Alias,
			Domain:         string(chainRedirect.Domain),
			LoadBalancerID: string(chainRedirect.PortalApplicationID),
		})
	}

	altruistURL := ""
	// for now we can assume each chain has only one altruist
	if len(c.Altruists) > 0 {
		altruistURL = formatAltruistURL(c.Altruists[0])
	}

	var chainID string
	if c.Checks[types.ChainCheckTypeChain].EVMChainID != 0 {
		chainID = strconv.Itoa(int(c.Checks[types.ChainCheckTypeChain].EVMChainID))
	}

	return types.Blockchain{
		ID:                string(c.ID),
		Blockchain:        c.Blockchain,
		ChainID:           chainID,
		ChainIDCheck:      c.Checks[types.ChainCheckTypeChain].Payload,
		Description:       c.Description,
		EnforceResult:     c.EnforceResult,
		Path:              c.Path,
		Ticker:            c.Ticker,
		BlockchainAliases: c.ChainAliases,
		LogLimitBlocks:    int(c.LogLimitBlocks),
		RequestTimeout:    int(c.RequestTimeout),
		Active:            c.Active,
		Altruist:          altruistURL,
		SyncCheckOptions: types.SyncCheckOptions{
			Body:      c.Checks[types.ChainCheckTypeSync].Payload,
			ResultKey: c.Checks[types.ChainCheckTypeSync].ResultKey,
			Allowance: int(c.Checks[types.ChainCheckTypeSync].Allowance),
		},
		Redirects: redirects,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func formatAltruistURL(altruist types.Altruist) string {
	formattedURL := string(altruist.URL)

	// insert the basic auth into the altruist URL
	if altruist.AuthType == types.ChainAuthTypeBasicAuth {
		formattedURL = strings.Replace(formattedURL, "https://", "", 1)
		formattedURL = fmt.Sprintf("https://%s@%s", altruist.Auth, formattedURL)
		return formattedURL
	}

	return formattedURL
}

func ConvertToLegacyPayPlan(c types.Plan) types.PayPlan {
	return types.PayPlan{
		Type:  c.Type,
		Limit: int(c.LegacyDailyLimit),
	}
}

/* Legacy Struct to V2 Struct Adaptors */

// Creates the struct with all fields needed to create a new Account, PortalApp & AAT
// LoadBalancer must be sent to PHD containing its Application already defined inside PUB
// This way the Account and PortalApp can be created in only one operation (no PHD client changes needed)
func ConvertToV2AccountPortalAppAndAAT(lb types.LoadBalancer) (types.Account, types.PortalApp, types.AAT) {
	app := lb.Applications[0]
	owner := lb.Users[0]

	account := types.Account{ // Account ID is created inside postgresdriver
		PlanType: types.PayPlanType(lb.Applications[0].Limit.PayPlan.Type),
	}

	portalApp := types.PortalApp{ // Portal App ID is created inside postgresdriver
		Name:      lb.Name,
		Gigastake: lb.Gigastake,
		Settings: types.Settings{
			Environment: types.EnvironmentProduction,
			SecretKey:   app.GatewaySettings.SecretKey,
		},
		Notifications: map[types.NotificationType]types.AppNotification{
			types.NotificationTypeEmail: {
				Active:      true,
				Destination: owner.Email,
				Events: map[types.NotificationEvent]bool{
					types.NotificationEventSignedUp:      app.NotificationSettings.SignedUp,
					types.NotificationEventQuarter:       app.NotificationSettings.Quarter,
					types.NotificationEventHalf:          app.NotificationSettings.Half,
					types.NotificationEventThreeQuarters: app.NotificationSettings.ThreeQuarters,
					types.NotificationEventFull:          app.NotificationSettings.Full,
				},
			},
		},

		LegacyFields: types.LegacyFields{
			RequestTimeout:    int32(lb.RequestTimeout),
			GigastakeRedirect: lb.GigastakeRedirect,
			StickyOptions:     lb.StickyOptions,
		},
	}

	aat := types.AAT{ // AAT ID (ProtocolAppID) is created inside postgresdriver
		Address:         app.GatewayAAT.Address,
		PublicKey:       app.GatewayAAT.ApplicationPublicKey,
		ClientPublicKey: app.GatewayAAT.ClientPublicKey,
		PrivateKey:      app.GatewayAAT.PrivateKey,
		Signature:       app.GatewayAAT.ApplicationSignature,
		Version:         app.GatewayAAT.Version,
	}

	return account, portalApp, aat
}

// Converts the existing UpdateApplication struct to a new one that updates all relevant fields in the PortalApp
func ConvertToV2UpdatePortalApp(u types.UpdateApplication, appID string) types.UpdatePortalApp {
	var (
		settings                             *types.UpdateAppSettings
		notifications                        []types.UpdateAppNotifications
		whitelists                           *types.WhitelistsObject
		contractWhitelists, methodWhitelists []types.ChainIDWhitelists
	)

	if u.NotificationSettings != nil {
		notificationEvents := []types.NotificationEvent{}

		if u.NotificationSettings.SignedUp != nil && *u.NotificationSettings.SignedUp {
			notificationEvents = append(notificationEvents, types.NotificationEventSignedUp)
		}
		if u.NotificationSettings.Quarter != nil && *u.NotificationSettings.Quarter {
			notificationEvents = append(notificationEvents, types.NotificationEventQuarter)
		}
		if u.NotificationSettings.Half != nil && *u.NotificationSettings.Half {
			notificationEvents = append(notificationEvents, types.NotificationEventHalf)
		}
		if u.NotificationSettings.ThreeQuarters != nil && *u.NotificationSettings.ThreeQuarters {
			notificationEvents = append(notificationEvents, types.NotificationEventThreeQuarters)
		}
		if u.NotificationSettings.Full != nil && *u.NotificationSettings.Full {
			notificationEvents = append(notificationEvents, types.NotificationEventFull)
		}

		notifications = []types.UpdateAppNotifications{
			{NotificationType: types.NotificationTypeEmail, Events: notificationEvents},
		}
	}

	if u.GatewaySettings != nil {
		settings = &types.UpdateAppSettings{SecretKey: u.GatewaySettings.SecretKey}
		if u.GatewaySettings.SecretKeyRequired != nil {
			settings.SecretKeyRequired = *u.GatewaySettings.SecretKeyRequired
		}

		for _, chainContracts := range u.GatewaySettings.WhitelistContracts {
			contracts := []string{}
			contracts = append(contracts, chainContracts.Contracts...)
			contractWhitelists = append(contractWhitelists, types.ChainIDWhitelists{
				ChainID: chainContracts.BlockchainID, Values: contracts,
			})
		}
		for _, chainMethods := range u.GatewaySettings.WhitelistMethods {
			methods := []string{}
			methods = append(methods, chainMethods.Methods...)
			methodWhitelists = append(methodWhitelists, types.ChainIDWhitelists{
				ChainID: chainMethods.BlockchainID, Values: methods,
			})
		}

		whitelists = &types.WhitelistsObject{
			AppWhitelists: [3]types.ApplicationWhitelists{
				{Type: types.WhitelistTypeOrigins, Values: u.GatewaySettings.WhitelistOrigins},
				{Type: types.WhitelistTypeUserAgents, Values: u.GatewaySettings.WhitelistUserAgents},
				{Type: types.WhitelistTypeBlockchains, Values: u.GatewaySettings.WhitelistBlockchains},
			},
			ChainWhitelists: [2]types.ChainWhitelists{
				{Type: types.WhitelistTypeContracts, Values: contractWhitelists},
				{Type: types.WhitelistTypeMethods, Values: methodWhitelists},
			},
		}
	}

	return types.UpdatePortalApp{
		AppID: types.PortalAppID(appID), Name: u.Name,
		Settings: settings, Notifications: notifications, Whitelists: whitelists,
	}
}

func ConvertToV2AccountUserAccess(u types.UserAccess) types.AccountUserAccess {
	authType := getAuthType(u.UserID)
	return types.AccountUserAccess{
		Email:           types.Email(u.Email),
		RoleName:        u.RoleName,
		Accepted:        u.Accepted,
		ProviderUserIDs: map[types.AuthType]string{authType: u.UserID},
	}
}

func ConvertToV2UpdateAcceptAccountUserAccess(u types.UpdateUserAccess, accountID types.AccountID, userID types.UserID) types.UpdateAcceptAccountUser {
	authType := getAuthType(u.UserID)
	return types.UpdateAcceptAccountUser{
		AccountID:        accountID,
		UserID:           userID,
		AuthProviderType: authType,
		ProviderUserID:   types.ProviderUserID(u.UserID),
	}
}

func ConvertToV2UpdateAccountUserAccess(u types.UpdateUserAccess, lbID string, userID string) types.UpdateAccountUserRole {
	return types.UpdateAccountUserRole{
		LegacyLoadBalancerID: lbID,
		UserID:               types.UserID(userID),
		RoleName:             u.RoleName,
	}
}

func getAuthType(userID string) types.AuthType {
	switch {
	case len(userID) == 24:
		return types.AuthTypeAuth0Username
	default:
		return types.AuthTypeAuth0Github
	}
}

func ConvertToV2Chain(b types.Blockchain) types.Chain {
	checks := map[types.ChainCheckType]types.Check{
		types.ChainCheckTypeSync: {
			Type:      types.ChainCheckTypeSync,
			Payload:   b.SyncCheckOptions.Body,
			ResultKey: b.SyncCheckOptions.ResultKey,
			Allowance: int32(b.SyncCheckOptions.Allowance),
		},
	}
	if b.ChainIDCheck != "" {
		chainID, _ := strconv.Atoi(b.ChainID)
		evmChainID := int32(chainID)

		checks[types.ChainCheckTypeChain] = types.Check{
			Payload:    b.ChainIDCheck,
			EVMChainID: evmChainID,
		}
	}

	altruist := parseAltruistURL(b.Altruist)

	redirects := []types.GigastakeRedirect{}
	for _, redirect := range b.Redirects {
		redirects = append(redirects, types.GigastakeRedirect{
			PortalApplicationID: types.PortalAppID(redirect.LoadBalancerID),
			Domain:              types.RedirectDomain(redirect.Domain),
			Alias:               redirect.Alias,
		})
	}

	return types.Chain{
		ID:             types.RelayChainID(b.ID),
		Blockchain:     b.Blockchain,
		Description:    b.Description,
		EnforceResult:  b.EnforceResult,
		Path:           b.Path,
		Ticker:         b.Ticker,
		ChainAliases:   b.BlockchainAliases,
		LogLimitBlocks: int32(b.LogLimitBlocks),
		RequestTimeout: int32(b.RequestTimeout),
		Active:         b.Active,
		Altruists:      []types.Altruist{altruist},
		Redirects:      redirects,
		Checks:         checks,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}

func parseAltruistURL(rawURL string) types.Altruist {
	parsedURL, _ := url.Parse(rawURL)

	auth := ""
	authType := types.ChainAuthTypeNone

	if parsedURL.User != nil {
		auth = parsedURL.User.String()
		parsedURL.User = nil
		authType = types.ChainAuthTypeBasicAuth
	}

	altruistURL := types.AltruistURL(strings.TrimPrefix(parsedURL.String(), "//"))

	return types.Altruist{
		URL:      altruistURL,
		Auth:     auth,
		AuthType: authType,
	}
}

func ConvertToV2UpdateChain(u types.UpdateBlockchain) types.Chain {
	allowance := *u.Allowance
	checks := map[types.ChainCheckType]types.Check{
		types.ChainCheckTypeSync: {Type: types.ChainCheckTypeSync, Payload: u.Body, ResultKey: u.ResultKey, Allowance: int32(allowance)},
	}
	if u.ChainIDCheck != "" {
		checks[types.ChainCheckTypeChain] = types.Check{Payload: u.ChainIDCheck}
	}

	altruist := parseAltruistURL(u.Altruist)

	return types.Chain{
		Blockchain:     u.Blockchain,
		Description:    u.Description,
		EnforceResult:  u.EnforceResult,
		Path:           u.Path,
		Ticker:         u.Ticker,
		ChainAliases:   u.BlockchainAliases,
		LogLimitBlocks: int32(u.LogLimitBlocks),
		RequestTimeout: int32(u.RequestTimeout),
		Altruists:      []types.Altruist{altruist},
		Checks:         checks,
	}
}

func ConvertToV2Redirect(r types.Redirect) types.GigastakeRedirect {
	return types.GigastakeRedirect{
		PortalApplicationID: types.PortalAppID(r.LoadBalancerID),
		Alias:               r.Alias,
		Domain:              types.RedirectDomain(r.Domain),
	}
}

/* Utils */
func sortUsersByRole(users []types.UserAccess) {
	roleWeight := map[types.RoleName]int{types.RoleOwner: 0, types.RoleAdmin: 1, types.RoleMember: 2}

	sort.Slice(users, func(i, j int) bool {
		if roleWeight[users[i].RoleName] != roleWeight[users[j].RoleName] {
			return roleWeight[users[i].RoleName] < roleWeight[users[j].RoleName]
		}
		return users[i].UserID < users[j].UserID
	})
}
