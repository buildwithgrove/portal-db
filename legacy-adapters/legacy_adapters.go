package legacyadapters

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	v1Types "github.com/pokt-foundation/portal-db/types"
	v2Types "github.com/pokt-foundation/portal-db/v2/types"
)

/*
This package exists to convert to and from between the legacy (v1) and v2 types.
It is used in PHD to enable the use of the new V2 schema in PHD without affecting upstream services.

Once the V2 migration is completed this package may be removed from this repo.
*/

/* V2 Struct to Legacy Struct Adaptors */
func ConvertPortalAppToLegacyLoadBalancer(a v2Types.PortalApp, account v2Types.Account) v1Types.LoadBalancer {
	var userID string

	var users []v1Types.UserAccess
	for _, accountUser := range account.Users {
		for portalAppID, roleName := range accountUser.PortalAppRoles {
			if portalAppID == a.ID {
				users = append(users, v1Types.UserAccess{
					UserID:   string(accountUser.UserID),
					Email:    string(accountUser.Email),
					RoleName: v1Types.RoleName(roleName),
					Accepted: accountUser.Accepted,
				})

				if roleName == v2Types.RoleOwner {
					userID = string(accountUser.UserID)
				}
			}

		}
	}

	sortUsersByRole(users)

	return v1Types.LoadBalancer{
		ID:           string(a.ID),
		AccountID:    string(account.ID),
		Name:         a.Name,
		UserID:       userID,
		Applications: ConvertPortalAppToLegacyApplications(a, userID),
		Users:        users,
		Integrations: v1Types.AccountIntegrations{
			CovalentAPIKeyFree: account.Integrations.CovalentAPIKeyFree,
			CovalentAPIKeyPaid: account.Integrations.CovalentAPIKeyPaid,
		},
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
		Gigastake:         false,
		RequestTimeout:    int(a.LegacyFields.RequestTimeout),
		GigastakeRedirect: true,
	}
}

func ConvertPortalAppToLegacyApplications(a v2Types.PortalApp, userID string) []*v1Types.Application {
	baseApp := v1Types.Application{
		UserID:          userID,
		Name:            a.Name,
		GatewaySettings: ConvertToLegacyGatewaySettings(a),
		Limit: v1Types.AppLimit{
			PayPlan:     v1Types.PayPlan{Type: v1Types.PayPlanType(a.LegacyFields.PlanType), Limit: int(a.LegacyFields.DailyLimit)},
			CustomLimit: int(a.LegacyFields.CustomLimit),
		},
		NotificationSettings: v1Types.NotificationSettings{
			SignedUp:      a.Notifications[v2Types.NotificationTypeEmail].Events[v2Types.NotificationEventSignedUp],
			Quarter:       a.Notifications[v2Types.NotificationTypeEmail].Events[v2Types.NotificationEventQuarter],
			Half:          a.Notifications[v2Types.NotificationTypeEmail].Events[v2Types.NotificationEventHalf],
			ThreeQuarters: a.Notifications[v2Types.NotificationTypeEmail].Events[v2Types.NotificationEventThreeQuarters],
			Full:          a.Notifications[v2Types.NotificationTypeEmail].Events[v2Types.NotificationEventFull],
		},
		FirstDateSurpassed: a.LegacyFields.FirstDateSurpassed,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}

	var legacyApps []*v1Types.Application

	for _, aat := range a.AATs {
		app := baseApp

		app.ID = string(aat.LegacyAppID)

		app.GatewayAAT = v1Types.GatewayAAT{
			Address:              aat.Address,
			ApplicationPublicKey: aat.PublicKey,
			ApplicationSignature: aat.Signature,
			ClientPublicKey:      aat.ClientPublicKey,
			Version:              aat.Version,
		}

		legacyApps = append(legacyApps, &app)
	}

	return legacyApps
}

func ConvertPortalAppToLegacyApplication(a v2Types.PortalApp, userID string, aat v2Types.AAT) v1Types.Application {
	return v1Types.Application{
		ID:              string(aat.LegacyAppID),
		UserID:          userID,
		Name:            a.Name,
		GatewaySettings: ConvertToLegacyGatewaySettings(a),
		GatewayAAT: v1Types.GatewayAAT{
			Address:              aat.Address,
			ApplicationPublicKey: aat.PublicKey,
			ApplicationSignature: aat.Signature,
			ClientPublicKey:      aat.ClientPublicKey,
			Version:              aat.Version,
		},
		Limit: v1Types.AppLimit{
			PayPlan:     v1Types.PayPlan{Type: v1Types.PayPlanType(a.LegacyFields.PlanType), Limit: int(a.LegacyFields.DailyLimit)},
			CustomLimit: int(a.LegacyFields.CustomLimit),
		},
		NotificationSettings: v1Types.NotificationSettings{
			SignedUp:      a.Notifications[v2Types.NotificationTypeEmail].Events[v2Types.NotificationEventSignedUp],
			Quarter:       a.Notifications[v2Types.NotificationTypeEmail].Events[v2Types.NotificationEventQuarter],
			Half:          a.Notifications[v2Types.NotificationTypeEmail].Events[v2Types.NotificationEventHalf],
			ThreeQuarters: a.Notifications[v2Types.NotificationTypeEmail].Events[v2Types.NotificationEventThreeQuarters],
			Full:          a.Notifications[v2Types.NotificationTypeEmail].Events[v2Types.NotificationEventFull],
		},
		FirstDateSurpassed: a.LegacyFields.FirstDateSurpassed,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
}

func ConvertToLegacyGatewaySettings(a v2Types.PortalApp) v1Types.GatewaySettings {
	gatewaySettings := v1Types.GatewaySettings{
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
		gatewaySettings.WhitelistContracts = append(gatewaySettings.WhitelistContracts, v1Types.WhitelistContracts{
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
		gatewaySettings.WhitelistMethods = append(gatewaySettings.WhitelistMethods, v1Types.WhitelistMethods{
			BlockchainID: string(chainID), Methods: methodList},
		)
	}
	sort.Slice(gatewaySettings.WhitelistMethods, func(i, j int) bool {
		return gatewaySettings.WhitelistMethods[i].BlockchainID < gatewaySettings.WhitelistMethods[j].BlockchainID
	})

	return gatewaySettings
}

func ConvertChainToLegacyGigastakeLoadBalancer(c v2Types.Chain) v1Types.LoadBalancer {
	var chainLBID string
	var apps []*v1Types.Application

	for _, app := range c.GigastakeApps {
		if chainLBID == "" {
			chainLBID = string(app.LegacyLBID)
		}

		apps = append(apps, ConvertGigastakeAppToLegacyApplication(app))
	}

	return v1Types.LoadBalancer{
		ID:           chainLBID,
		Applications: apps,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		Gigastake:    true,
	}
}

func ConvertGigastakeAppToLegacyApplication(a *v2Types.GigastakeApp) *v1Types.Application {
	return &v1Types.Application{
		ID:   string(a.ID),
		Name: a.Name,
		GatewayAAT: v1Types.GatewayAAT{
			Address:              a.Address,
			ApplicationPublicKey: a.PublicKey,
			ApplicationSignature: a.Signature,
			ClientPublicKey:      a.ClientPublicKey,
			Version:              a.Version,
		},
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func ConvertToLegacyBlockchain(c v2Types.Chain) v1Types.Blockchain {
	var redirectLBID string
	for _, gigastakeApp := range c.GigastakeApps {
		if gigastakeApp.LegacyLBID != "" {
			redirectLBID = string(gigastakeApp.LegacyLBID)
		}
	}

	var redirects []v1Types.Redirect
	var aliases []string
	for alias, domains := range c.AliasDomains {
		for _, domain := range domains {
			redirects = append(redirects, v1Types.Redirect{
				Alias:          string(alias),
				Domain:         string(domain),
				LoadBalancerID: string(redirectLBID),
			})
		}

		aliases = append(aliases, string(alias))
	}

	altruistURL := ""
	// for now we can assume each chain has only one altruist
	if len(c.Altruists) > 0 {
		var altruist v2Types.Altruist
		for _, alt := range c.Altruists {
			altruist = alt
			break
		}
		altruistURL = formatAltruistURL(altruist)
	}

	var chainID string
	if c.Checks[v2Types.ChainCheckTypeChain].EVMChainID != 0 {
		chainID = strconv.Itoa(int(c.Checks[v2Types.ChainCheckTypeChain].EVMChainID))
	}

	return v1Types.Blockchain{
		ID:                string(c.ID),
		Blockchain:        c.Blockchain,
		ChainID:           chainID,
		ChainIDCheck:      c.Checks[v2Types.ChainCheckTypeChain].Payload,
		Description:       c.Description,
		EnforceResult:     c.EnforceResult,
		Path:              c.Path,
		Ticker:            c.Ticker,
		BlockchainAliases: aliases,
		LogLimitBlocks:    int(c.LogLimitBlocks),
		RequestTimeout:    int(c.RequestTimeout),
		Active:            c.Active,
		Altruist:          altruistURL,
		SyncCheckOptions: v1Types.SyncCheckOptions{
			Body:      c.Checks[v2Types.ChainCheckTypeSync].Payload,
			ResultKey: c.Checks[v2Types.ChainCheckTypeSync].ResultKey,
			Allowance: int(c.Checks[v2Types.ChainCheckTypeSync].Allowance),
		},
		Redirects: redirects,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func formatAltruistURL(altruist v2Types.Altruist) string {
	formattedURL := string(altruist.URL)

	// insert the basic auth into the altruist URL
	if altruist.AuthType == v2Types.ChainAuthTypeBasicAuth {
		formattedURL = strings.Replace(formattedURL, "https://", "", 1)
		formattedURL = fmt.Sprintf("https://%s@%s", altruist.Auth, formattedURL)
		return formattedURL
	}

	return formattedURL
}

func ConvertToLegacyPayPlan(c v2Types.Plan) v1Types.PayPlan {
	return v1Types.PayPlan{
		Type:  v1Types.PayPlanType(c.Type),
		Limit: int(c.LegacyDailyLimit),
	}
}

/* Legacy Struct to V2 Struct Adaptors */

// Creates the struct with all fields needed to create a new PortalApp & AAT
// LoadBalancer must be sent to PHD containing its Application already defined inside PUB
// This way the Account and PortalApp can be created in only one operation (no PHD client changes needed)
func ConvertToV2PortalAppAndAAT(lb v1Types.LoadBalancer) (v2Types.PortalApp, v2Types.AAT) {
	app := lb.Applications[0]
	owner := lb.Users[0]

	portalApp := v2Types.PortalApp{ // Portal App ID is created inside postgresdriver
		Name: lb.Name,
		Settings: v2Types.Settings{
			Environment: v2Types.EnvironmentProduction,
			SecretKey:   app.GatewaySettings.SecretKey,
		},
		Notifications: map[v2Types.NotificationType]v2Types.AppNotification{
			v2Types.NotificationTypeEmail: {
				Type:        v2Types.NotificationTypeEmail,
				Active:      true,
				Destination: owner.Email,
				Events: map[v2Types.NotificationEvent]bool{
					v2Types.NotificationEventSignedUp:      app.NotificationSettings.SignedUp,
					v2Types.NotificationEventQuarter:       app.NotificationSettings.Quarter,
					v2Types.NotificationEventHalf:          app.NotificationSettings.Half,
					v2Types.NotificationEventThreeQuarters: app.NotificationSettings.ThreeQuarters,
					v2Types.NotificationEventFull:          app.NotificationSettings.Full,
				},
			},
		},

		LegacyFields: v2Types.LegacyFields{
			PlanType:       v2Types.PayPlanType(app.Limit.PayPlan.Type),
			DailyLimit:     int32(app.Limit.PayPlan.Limit),
			CustomLimit:    int32(app.Limit.CustomLimit),
			RequestTimeout: int32(lb.RequestTimeout),
		},
	}

	aat := v2Types.AAT{ // AAT ID (ProtocolAppID) is created inside postgresdriver
		Address:         app.GatewayAAT.Address,
		PublicKey:       app.GatewayAAT.ApplicationPublicKey,
		ClientPublicKey: app.GatewayAAT.ClientPublicKey,
		PrivateKey:      app.GatewayAAT.PrivateKey,
		Signature:       app.GatewayAAT.ApplicationSignature,
		Version:         app.GatewayAAT.Version,
	}

	return portalApp, aat
}

// Converts the existing UpdateApplication struct to a new one that updates all relevant fields in the PortalApp
func ConvertToV2UpdatePortalApp(u v1Types.UpdateApplication, lbID string) v2Types.UpdatePortalApp {
	var (
		settings                             *v2Types.UpdateAppSettings
		notifications                        []v2Types.UpdateAppNotifications
		whitelists                           *v2Types.WhitelistsObject
		contractWhitelists, methodWhitelists []v2Types.ChainIDWhitelists
	)

	if u.NotificationSettings != nil {
		notificationEvents := []v2Types.NotificationEvent{}

		if u.NotificationSettings.SignedUp != nil && *u.NotificationSettings.SignedUp {
			notificationEvents = append(notificationEvents, v2Types.NotificationEventSignedUp)
		}
		if u.NotificationSettings.Quarter != nil && *u.NotificationSettings.Quarter {
			notificationEvents = append(notificationEvents, v2Types.NotificationEventQuarter)
		}
		if u.NotificationSettings.Half != nil && *u.NotificationSettings.Half {
			notificationEvents = append(notificationEvents, v2Types.NotificationEventHalf)
		}
		if u.NotificationSettings.ThreeQuarters != nil && *u.NotificationSettings.ThreeQuarters {
			notificationEvents = append(notificationEvents, v2Types.NotificationEventThreeQuarters)
		}
		if u.NotificationSettings.Full != nil && *u.NotificationSettings.Full {
			notificationEvents = append(notificationEvents, v2Types.NotificationEventFull)
		}

		notifications = []v2Types.UpdateAppNotifications{
			{Active: true, NotificationType: v2Types.NotificationTypeEmail, Events: notificationEvents},
		}
	}

	if u.GatewaySettings != nil {
		settings = &v2Types.UpdateAppSettings{
			SecretKey:   u.GatewaySettings.SecretKey,
			Environment: v2Types.EnvironmentProduction,
		}
		if u.GatewaySettings.SecretKeyRequired != nil {
			settings.SecretKeyRequired = *u.GatewaySettings.SecretKeyRequired
		}

		for _, chainContracts := range u.GatewaySettings.WhitelistContracts {
			contracts := []string{}
			contracts = append(contracts, chainContracts.Contracts...)
			contractWhitelists = append(contractWhitelists, v2Types.ChainIDWhitelists{
				ChainID: chainContracts.BlockchainID, Values: contracts,
			})
		}
		for _, chainMethods := range u.GatewaySettings.WhitelistMethods {
			methods := []string{}
			methods = append(methods, chainMethods.Methods...)
			methodWhitelists = append(methodWhitelists, v2Types.ChainIDWhitelists{
				ChainID: chainMethods.BlockchainID, Values: methods,
			})
		}

		whitelists = &v2Types.WhitelistsObject{
			AppWhitelists: [3]v2Types.ApplicationWhitelists{
				{Type: v2Types.WhitelistTypeOrigins, Values: u.GatewaySettings.WhitelistOrigins},
				{Type: v2Types.WhitelistTypeUserAgents, Values: u.GatewaySettings.WhitelistUserAgents},
				{Type: v2Types.WhitelistTypeBlockchains, Values: u.GatewaySettings.WhitelistBlockchains},
			},
			ChainWhitelists: [2]v2Types.ChainWhitelists{
				{Type: v2Types.WhitelistTypeContracts, Values: contractWhitelists},
				{Type: v2Types.WhitelistTypeMethods, Values: methodWhitelists},
			},
		}
	}

	update := v2Types.UpdatePortalApp{
		AppID:         v2Types.PortalAppID(lbID),
		Name:          u.Name,
		Settings:      settings,
		Notifications: notifications,
		Whitelists:    whitelists,
	}

	if u.Limit != nil {
		update.PlanType = v2Types.PayPlanType(u.Limit.PayPlan.Type)
		if u.Limit.PayPlan.Type == v1Types.Enterprise {
			update.CustomLimit = int32(u.Limit.CustomLimit)
		} else {
			update.DailyLimit = int32(u.Limit.PayPlan.Limit)
		}
	}

	return update
}

func ConvertToV2AccountUserAccess(u v1Types.UserAccess, lbID, accountID string) v2Types.AccountUserAccess {
	portalAppID := v2Types.PortalAppID(lbID)

	userAccess := v2Types.AccountUserAccess{
		AccountID: v2Types.AccountID(accountID),
		UserID:    v2Types.UserID(u.UserID),
		Email:     v2Types.Email(u.Email),
		Owner:     u.RoleName == v1Types.RoleOwner,
		Accepted:  u.Accepted,
		PortalAppRoles: map[v2Types.PortalAppID]v2Types.RoleName{
			portalAppID: v2Types.RoleName(u.RoleName),
		},
	}

	return userAccess
}

func ConvertToV2UpdateAcceptAccountUserAccess(u v1Types.UpdateUserAccess, portalAppID v2Types.PortalAppID, userID v2Types.UserID) v2Types.UpdateAcceptAccountUser {
	authType := getAuthType(u.UserID)
	return v2Types.UpdateAcceptAccountUser{
		PortalAppID:      portalAppID,
		UserID:           userID,
		AuthProviderType: authType,
		ProviderUserID:   v2Types.ProviderUserID(u.UserID),
	}
}

func ConvertToV2UpdateAccountUserAccess(u v1Types.UpdateUserAccess, lbID string, userID string) v2Types.UpdateAccountUserRole {
	return v2Types.UpdateAccountUserRole{
		PortalAppID: v2Types.PortalAppID(lbID),
		UserID:      v2Types.UserID(userID),
		RoleName:    v2Types.RoleName(u.RoleName),
	}
}

func getAuthType(userID string) v2Types.AuthType {
	switch {
	case len(userID) == 24:
		return v2Types.AuthTypeAuth0Username
	default:
		return v2Types.AuthTypeAuth0Github
	}
}

/* Utils */
func sortUsersByRole(users []v1Types.UserAccess) {
	roleWeight := map[v1Types.RoleName]int{v1Types.RoleOwner: 0, v1Types.RoleAdmin: 1, v1Types.RoleMember: 2}

	sort.Slice(users, func(i, j int) bool {
		if roleWeight[users[i].RoleName] != roleWeight[users[j].RoleName] {
			return roleWeight[users[i].RoleName] < roleWeight[users[j].RoleName]
		}
		return users[i].UserID < users[j].UserID
	})
}
