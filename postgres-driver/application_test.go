package postgresdriver

import (
	"database/sql"
	"time"

	"github.com/pokt-foundation/portal-db/types"
)

func (ts *PGDriverTestSuite) Test_ReadApplications() {
	tests := []struct {
		name         string
		applications []*types.Application
		err          error
	}{
		{
			name: "Should return all Applications from the database ordered by application_id",
			applications: []*types.Application{
				{
					ID:     "test_app_47hfnths73j2se",
					UserID: "test_user_1dbffbdfeeb225",
					Name:   "pokt_app_123",
					URL:    "https://test.app123.io",
					Dummy:  true,
					Status: types.InService,
					GatewayAAT: types.GatewayAAT{
						Address:              "test_34715cae753e67c75fbb340442e7de8e",
						ApplicationPublicKey: "test_11b8d394ca331d7c7a71ca1896d630f6",
						ApplicationSignature: "test_89a3af6a587aec02cfade6f5000424c2",
						ClientPublicKey:      "test_1dc39a2e5a84a35bf030969a0b3231f7",
						PrivateKey:           "test_d2ce53f115f4ecb2208e9188800a85cf",
					},
					GatewaySettings: types.GatewaySettings{
						SecretKey:            "test_40f482d91a5ef2300ebb4e2308c",
						SecretKeyRequired:    true,
						WhitelistBlockchains: []string{"chain_1"},
						WhitelistOrigins:     []string{"origin_1", "origin_2"},
						WhitelistUserAgents:  []string{"user_agent_1", "user_agent_2", "user_agent_3"},
						WhitelistContracts: []types.WhitelistContracts{
							{BlockchainID: "TST2", Contracts: []string{"test_address_67890"}},
							{BlockchainID: "TST1", Contracts: []string{"test_address_12345", "test_address_44444"}},
						},
						WhitelistMethods: []types.WhitelistMethods{
							{BlockchainID: "TST1", Methods: []string{"GET"}},
							{BlockchainID: "TST2", Methods: []string{"PUT", "POST"}},
						},
					},
					Limit: types.AppLimit{
						PayPlan: types.PayPlan{Type: types.FreetierV0, Limit: 250_000},
					},
					NotificationSettings: types.NotificationSettings{
						SignedUp:      true,
						Quarter:       false,
						Half:          false,
						ThreeQuarters: true,
						Full:          true,
					},
				},
				{
					ID:     "test_app_5hdf7sh23jd828",
					UserID: "test_user_04228205bd261a",
					Name:   "pokt_app_456",
					URL:    "https://test.app456.io",
					Dummy:  true,
					Status: types.InService,
					GatewayAAT: types.GatewayAAT{
						Address:              "test_558c0225c7019e14ccf2e7379ad3eb50",
						ApplicationPublicKey: "test_96c981db344ab6920b7e87853838e285",
						ApplicationSignature: "test_1272a8ab4cbbf636f09bf4fa5395b885",
						ClientPublicKey:      "test_d709871777b89ed3051190f229ea3f01",
						PrivateKey:           "test_53e50765d8bc1fb41b3b0065dd8094de",
					},
					GatewaySettings: types.GatewaySettings{
						SecretKey:         "test_90210ac4bdd3423e24877d1ff92",
						SecretKeyRequired: false,
					},
					Limit: types.AppLimit{
						PayPlan:     types.PayPlan{Type: types.Enterprise},
						CustomLimit: 2_000_000,
					},
					NotificationSettings: types.NotificationSettings{
						SignedUp:      true,
						Quarter:       false,
						Half:          false,
						ThreeQuarters: true,
						Full:          true,
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		applications, err := ts.driver.ReadApplications(testCtx)
		ts.Equal(test.err, err)
		for i, app := range applications {
			ts.Equal(test.applications[i].ID, app.ID)
			ts.Equal(test.applications[i].UserID, app.UserID)
			ts.Equal(test.applications[i].Name, app.Name)
			ts.Equal(test.applications[i].URL, app.URL)
			ts.Equal(test.applications[i].Dummy, app.Dummy)
			ts.Equal(test.applications[i].Status, app.Status)
			ts.Equal(test.applications[i].GatewayAAT, app.GatewayAAT)
			ts.Equal(test.applications[i].GatewaySettings, app.GatewaySettings)
			ts.Equal(test.applications[i].Limit, app.Limit)
			ts.Equal(test.applications[i].NotificationSettings, app.NotificationSettings)
			ts.NotEmpty(app.CreatedAt)
			ts.NotEmpty(app.UpdatedAt)
		}
	}
}

func (ts *PGDriverTestSuite) Test_ReadPayPlans() {
	tests := []struct {
		name     string
		payPlans []*types.PayPlan
		err      error
	}{
		{
			name: "Should return all PayPlans from the database ordered by plan_type",
			payPlans: []*types.PayPlan{
				{Type: types.Enterprise, Limit: 0},
				{Type: types.FreetierV0, Limit: 250000},
				{Type: types.PayAsYouGoV0, Limit: 0},
				{Type: types.TestPlan10K, Limit: 10000},
				{Type: types.TestPlan90k, Limit: 90000},
				{Type: types.TestPlanV0, Limit: 100},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		payPlans, err := ts.driver.ReadPayPlans(testCtx)
		ts.Equal(test.payPlans, payPlans)
		ts.Equal(test.err, err)
	}
}

func (ts *PGDriverTestSuite) Test_WriteApplication() {
	tests := []struct {
		name              string
		appInputs         []*types.Application
		expectedNumOfApps int
		expectedApp       SelectOneApplicationRow
		err               error
	}{
		{
			name: "Should create a single load balancer successfully with correct input",
			appInputs: []*types.Application{
				{
					Name:   "pokt_app_789",
					UserID: "test_user_47fhsd75jd756sh",
					Dummy:  true,
					Status: types.InService,
					GatewayAAT: types.GatewayAAT{
						Address:              "test_e209a2d1f3454ddc69cb9333d547bbcf",
						ApplicationPublicKey: "test_b95c35affacf6df4a5585388490542f0",
						ApplicationSignature: "test_e59760339d9ce02972d1080d73446c90",
						ClientPublicKey:      "test_d591178ab3f48f45b243303fe77dc8c3",
						PrivateKey:           "test_f403700aed7e039c0a8fc2dd22da6fd9",
					},
					GatewaySettings: types.GatewaySettings{
						SecretKey:         "test_489574398f34uhf4uhjf9328jf23f98j",
						SecretKeyRequired: true,
					},
					Limit: types.AppLimit{
						PayPlan: types.PayPlan{Type: types.FreetierV0},
					},
					NotificationSettings: types.NotificationSettings{
						SignedUp:      true,
						Quarter:       false,
						Half:          false,
						ThreeQuarters: true,
						Full:          true,
					},
				},
			},
			expectedNumOfApps: 3,
			expectedApp: SelectOneApplicationRow{
				Name:              sql.NullString{Valid: true, String: "pokt_app_789"},
				UserID:            sql.NullString{Valid: true, String: "test_user_47fhsd75jd756sh"},
				Dummy:             sql.NullBool{Valid: true, Bool: true},
				Status:            sql.NullString{Valid: true, String: "IN_SERVICE"},
				GaAddress:         sql.NullString{Valid: true, String: "test_e209a2d1f3454ddc69cb9333d547bbcf"},
				GaClientPublicKey: sql.NullString{Valid: true, String: "test_d591178ab3f48f45b243303fe77dc8c3"},
				GaPrivateKey:      sql.NullString{Valid: true, String: "test_f403700aed7e039c0a8fc2dd22da6fd9"},
				GaPublicKey:       sql.NullString{Valid: true, String: "test_b95c35affacf6df4a5585388490542f0"},
				GaSignature:       sql.NullString{Valid: true, String: "test_e59760339d9ce02972d1080d73446c90"},
				SecretKey:         sql.NullString{Valid: true, String: "test_489574398f34uhf4uhjf9328jf23f98j"},
				SecretKeyRequired: sql.NullBool{Valid: true, Bool: true},
				SignedUp:          sql.NullBool{Valid: true, Bool: true},
				OnQuarter:         sql.NullBool{Valid: true, Bool: false},
				OnHalf:            sql.NullBool{Valid: true, Bool: false},
				OnThreeQuarters:   sql.NullBool{Valid: true, Bool: true},
				OnFull:            sql.NullBool{Valid: true, Bool: true},
				PayPlan:           sql.NullString{Valid: true, String: "FREETIER_V0"},
			},
			err: nil,
		},
		{
			name: "Should fail if passing an invalid status",
			appInputs: []*types.Application{
				{Status: types.AppStatus("INVALID_STATUS")},
			},
			err: types.ErrInvalidAppStatus,
		},
		{
			name: "Should fail if passing an invalid pay plan",
			appInputs: []*types.Application{
				{
					Status: types.InService,
					Limit: types.AppLimit{
						PayPlan: types.PayPlan{Type: types.PayPlanType("INVALID_PAY_PLAN")},
					},
				},
			},
			err: types.ErrInvalidPayPlanType,
		},
		{
			name: "Should fail when trying to update to a non-enterprise plan with a custom limit",
			appInputs: []*types.Application{
				{
					Status: types.InService,
					Limit: types.AppLimit{
						PayPlan:     types.PayPlan{Type: types.PayAsYouGoV0},
						CustomLimit: 123,
					},
				},
			},
			err: types.ErrNotEnterprisePlan,
		},
	}

	for _, test := range tests {
		for _, input := range test.appInputs {
			createdApp, err := ts.driver.WriteApplication(testCtx, input)
			ts.Equal(test.err, err)
			if err == nil {
				ts.Len(createdApp.ID, 24)
				ts.Equal(input.Name, createdApp.Name)
				ts.NotEmpty(createdApp.CreatedAt)
				ts.NotEmpty(createdApp.UpdatedAt)

				apps, err := ts.driver.ReadApplications(testCtx)
				ts.Equal(test.err, err)
				ts.Len(apps, test.expectedNumOfApps)

				app, err := ts.driver.SelectOneApplication(testCtx, createdApp.ID)
				ts.Equal(test.err, err)
				for _, testInput := range test.appInputs {
					if testInput.Name == app.Name.String {
						ts.Equal(createdApp.ID, app.ApplicationID)
						ts.Equal(test.expectedApp.Dummy, app.Dummy)
						ts.Equal(test.expectedApp.Status, app.Status)
						ts.Equal(test.expectedApp.GaAddress, app.GaAddress)
						ts.Equal(test.expectedApp.GaClientPublicKey, app.GaClientPublicKey)
						ts.Equal(test.expectedApp.GaPrivateKey, app.GaPrivateKey)
						ts.Equal(test.expectedApp.GaPublicKey, app.GaPublicKey)
						ts.Equal(test.expectedApp.GaSignature, app.GaSignature)
						ts.Equal(test.expectedApp.SecretKey, app.SecretKey)
						ts.Equal(test.expectedApp.SecretKeyRequired, app.SecretKeyRequired)
						ts.Equal(test.expectedApp.SignedUp, app.SignedUp)
						ts.Equal(test.expectedApp.OnQuarter, app.OnQuarter)
						ts.Equal(test.expectedApp.OnHalf, app.OnHalf)
						ts.Equal(test.expectedApp.OnThreeQuarters, app.OnThreeQuarters)
						ts.Equal(test.expectedApp.OnFull, app.OnFull)
						ts.Equal(test.expectedApp.PayPlan, app.PayPlan)
						ts.NotEmpty(app.CreatedAt)
						ts.NotEmpty(app.UpdatedAt)
					}

				}
			}
		}
	}
}

func (ts *PGDriverTestSuite) Test_UpdateApplication() {
	tests := []struct {
		name                string
		appID               string
		appUpdate           *types.UpdateApplication
		expectedAfterUpdate SelectOneApplicationRow
		err                 error
	}{
		{
			name:  "Should update a single application successfully with all fields",
			appID: "test_app_47hfnths73j2se",
			appUpdate: &types.UpdateApplication{
				Name: "pokt_app_updated_lb",
				GatewaySettings: &types.UpdateGatewaySettings{
					WhitelistOrigins:    []string{"test-origin1", "test-origin2"},
					WhitelistUserAgents: []string{"test-agent1"},
					WhitelistContracts: []types.WhitelistContracts{
						{
							BlockchainID: "01",
							Contracts:    []string{"test-contract1"},
						},
					},
					WhitelistMethods: []types.WhitelistMethods{
						{
							BlockchainID: "01",
							Methods:      []string{"test-method1"},
						},
					},
					WhitelistBlockchains: []string{"test-chain1"},
				},
				NotificationSettings: &types.UpdateNotificationSettings{
					SignedUp:      boolPointer(false),
					Quarter:       boolPointer(true),
					Half:          boolPointer(true),
					ThreeQuarters: boolPointer(false),
					Full:          boolPointer(false),
				},
				Limit: &types.AppLimit{
					PayPlan: types.PayPlan{
						Type: types.Enterprise,
					},
					CustomLimit: 4_200_000,
				},
			},
			expectedAfterUpdate: SelectOneApplicationRow{
				Name:                 sql.NullString{Valid: true, String: "pokt_app_updated_lb"},
				WhitelistBlockchains: []string{"test-chain1"},
				WhitelistContracts:   "[{\"contracts\": [\"test-contract1\"], \"blockchainID\": \"01\"}]",
				WhitelistMethods:     "[{\"methods\": [\"test-method1\"], \"blockchainID\": \"01\"}]",
				WhitelistOrigins:     []string{"test-origin1", "test-origin2"},
				WhitelistUserAgents:  []string{"test-agent1"},
				SignedUp:             sql.NullBool{Valid: true, Bool: false},
				OnQuarter:            sql.NullBool{Valid: true, Bool: true},
				OnHalf:               sql.NullBool{Valid: true, Bool: true},
				OnThreeQuarters:      sql.NullBool{Valid: true, Bool: false},
				OnFull:               sql.NullBool{Valid: true, Bool: false},
				CustomLimit:          sql.NullInt32{Valid: true, Int32: 4_200_000},
				PayPlan:              sql.NullString{Valid: true, String: "ENTERPRISE"},
			},
			err: nil,
		},
		{
			name:  "Should update a single application successfully with only some fields",
			appID: "test_app_5hdf7sh23jd828",
			appUpdate: &types.UpdateApplication{
				GatewaySettings: &types.UpdateGatewaySettings{
					WhitelistOrigins:    []string{"test-origin1", "test-origin2"},
					WhitelistUserAgents: []string{"test-agent1"},
				},
				NotificationSettings: &types.UpdateNotificationSettings{
					Full: boolPointer(false),
				},
				Limit: &types.AppLimit{
					PayPlan: types.PayPlan{Type: types.PayAsYouGoV0},
				},
			},
			expectedAfterUpdate: SelectOneApplicationRow{
				Name:                 sql.NullString{Valid: true, String: "pokt_app_456"},
				WhitelistBlockchains: []string(nil),
				WhitelistOrigins:     []string{"test-origin1", "test-origin2"},
				WhitelistUserAgents:  []string{"test-agent1"},
				SignedUp:             sql.NullBool{Valid: true, Bool: true},
				OnQuarter:            sql.NullBool{Valid: true, Bool: false},
				OnHalf:               sql.NullBool{Valid: true, Bool: false},
				OnThreeQuarters:      sql.NullBool{Valid: true, Bool: true},
				OnFull:               sql.NullBool{Valid: true, Bool: false},
				CustomLimit:          sql.NullInt32{Valid: true, Int32: 0},
				PayPlan:              sql.NullString{Valid: true, String: "PAY_AS_YOU_GO_V0"},
			},
			err: nil,
		},
		{
			name:  "Should add whitelist contracts & methods",
			appID: "test_app_5hdf7sh23jd828",
			appUpdate: &types.UpdateApplication{
				GatewaySettings: &types.UpdateGatewaySettings{
					WhitelistContracts: []types.WhitelistContracts{
						{
							BlockchainID: "01",
							Contracts:    []string{"test-contract1", "test-contract2"},
						},
						{
							BlockchainID: "02",
							Contracts:    []string{"test-contract3", "test-contract4"},
						},
					},
					WhitelistMethods: []types.WhitelistMethods{
						{
							BlockchainID: "01",
							Methods:      []string{"test-method1", "test-method2", "test-method3"},
						},
						{
							BlockchainID: "03",
							Methods:      []string{"test-method4", "test-method5", "test-method6"},
						},
					},
				},
			},
			expectedAfterUpdate: SelectOneApplicationRow{
				Name:                 sql.NullString{Valid: true, String: "pokt_app_456"},
				WhitelistBlockchains: []string(nil),
				WhitelistOrigins:     []string{"test-origin1", "test-origin2"},
				WhitelistUserAgents:  []string{"test-agent1"},
				WhitelistContracts:   "[{\"contracts\": [\"test-contract1\", \"test-contract2\"], \"blockchainID\": \"01\"}, {\"contracts\": [\"test-contract3\", \"test-contract4\"], \"blockchainID\": \"02\"}]",
				WhitelistMethods:     "[{\"methods\": [\"test-method1\", \"test-method2\", \"test-method3\"], \"blockchainID\": \"01\"}, {\"methods\": [\"test-method4\", \"test-method5\", \"test-method6\"], \"blockchainID\": \"03\"}]",
				SignedUp:             sql.NullBool{Valid: true, Bool: true},
				OnQuarter:            sql.NullBool{Valid: true, Bool: false},
				OnHalf:               sql.NullBool{Valid: true, Bool: false},
				OnThreeQuarters:      sql.NullBool{Valid: true, Bool: true},
				OnFull:               sql.NullBool{Valid: true, Bool: false},
				CustomLimit:          sql.NullInt32{Valid: true, Int32: 0},
				PayPlan:              sql.NullString{Valid: true, String: "PAY_AS_YOU_GO_V0"},
			},
			err: nil,
		},
		{
			name:  "Should remove whitelist contracts & methods",
			appID: "test_app_5hdf7sh23jd828",
			appUpdate: &types.UpdateApplication{
				GatewaySettings: &types.UpdateGatewaySettings{
					WhitelistContracts: []types.WhitelistContracts{
						{
							BlockchainID: "01",
							Contracts:    []string{"test-contract2"},
						},
						{
							BlockchainID: "02",
							Contracts:    []string{},
						},
					},
					WhitelistMethods: []types.WhitelistMethods{
						{
							BlockchainID: "01",
							Methods:      []string{"test-method1"},
						},
						{
							BlockchainID: "03",
							Methods:      []string{"test-method5", "test-method6"},
						},
					},
				},
			},
			expectedAfterUpdate: SelectOneApplicationRow{
				Name:                 sql.NullString{Valid: true, String: "pokt_app_456"},
				WhitelistBlockchains: []string(nil),
				WhitelistOrigins:     []string{"test-origin1", "test-origin2"},
				WhitelistUserAgents:  []string{"test-agent1"},
				WhitelistContracts:   "[{\"contracts\": [\"test-contract2\"], \"blockchainID\": \"01\"}]",
				WhitelistMethods:     "[{\"methods\": [\"test-method1\"], \"blockchainID\": \"01\"}, {\"methods\": [\"test-method5\", \"test-method6\"], \"blockchainID\": \"03\"}]",
				SignedUp:             sql.NullBool{Valid: true, Bool: true},
				OnQuarter:            sql.NullBool{Valid: true, Bool: false},
				OnHalf:               sql.NullBool{Valid: true, Bool: false},
				OnThreeQuarters:      sql.NullBool{Valid: true, Bool: true},
				OnFull:               sql.NullBool{Valid: true, Bool: false},
				CustomLimit:          sql.NullInt32{Valid: true, Int32: 0},
				PayPlan:              sql.NullString{Valid: true, String: "PAY_AS_YOU_GO_V0"},
			},
			err: nil,
		},
		{
			name:  "Should fail if passing an invalid status",
			appID: "test_app_5hdf7sh23jd828",
			appUpdate: &types.UpdateApplication{
				Status: types.AppStatus("INVALID_STATUS"),
			},
			err: types.ErrInvalidAppStatus,
		},
		{
			name:  "Should fail if passing an invalid pay plan",
			appID: "test_app_5hdf7sh23jd828",
			appUpdate: &types.UpdateApplication{
				Status: types.InService,
				Limit: &types.AppLimit{
					PayPlan: types.PayPlan{Type: types.PayPlanType("INVALID_PAY_PLAN")},
				},
			},
			err: types.ErrInvalidPayPlanType,
		},
		{
			name:  "Should fail when trying to update to a non-enterprise plan with a custom limit",
			appID: "test_app_5hdf7sh23jd828",
			appUpdate: &types.UpdateApplication{
				Limit: &types.AppLimit{
					PayPlan:     types.PayPlan{Type: types.PayAsYouGoV0},
					CustomLimit: 123,
				},
			},
			err: types.ErrNotEnterprisePlan,
		},
		{
			name:  "Should fail when trying to update to an enterprise plan without a custom limit",
			appID: "test_app_5hdf7sh23jd828",
			appUpdate: &types.UpdateApplication{
				Limit: &types.AppLimit{
					PayPlan: types.PayPlan{Type: types.Enterprise},
				},
			},
			err: types.ErrEnterprisePlanNeedsCustomLimit,
		},
	}

	for _, test := range tests {
		_, err := ts.driver.SelectOneApplication(testCtx, test.appID)
		ts.NoError(err)

		err = ts.driver.UpdateApplication(testCtx, test.appID, test.appUpdate)
		ts.Equal(test.err, err)
		if err == nil {
			appAfterUpdate, err := ts.driver.SelectOneApplication(testCtx, test.appID)
			ts.NoError(err)
			ts.Equal(test.expectedAfterUpdate.Name, appAfterUpdate.Name)
			ts.Equal(test.expectedAfterUpdate.WhitelistBlockchains, appAfterUpdate.WhitelistBlockchains)
			ts.Equal(test.expectedAfterUpdate.WhitelistContracts, appAfterUpdate.WhitelistContracts)
			ts.Equal(test.expectedAfterUpdate.WhitelistMethods, appAfterUpdate.WhitelistMethods)
			ts.Equal(test.expectedAfterUpdate.WhitelistOrigins, appAfterUpdate.WhitelistOrigins)
			ts.Equal(test.expectedAfterUpdate.WhitelistUserAgents, appAfterUpdate.WhitelistUserAgents)
			ts.Equal(test.expectedAfterUpdate.SignedUp, appAfterUpdate.SignedUp)
			ts.Equal(test.expectedAfterUpdate.OnQuarter, appAfterUpdate.OnQuarter)
			ts.Equal(test.expectedAfterUpdate.OnHalf, appAfterUpdate.OnHalf)
			ts.Equal(test.expectedAfterUpdate.OnThreeQuarters, appAfterUpdate.OnThreeQuarters)
			ts.Equal(test.expectedAfterUpdate.OnFull, appAfterUpdate.OnFull)
			ts.Equal(test.expectedAfterUpdate.CustomLimit, appAfterUpdate.CustomLimit)
			ts.Equal(test.expectedAfterUpdate.PayPlan, appAfterUpdate.PayPlan)
		}
	}
}

func (ts *PGDriverTestSuite) Test_UpdateAppFirstDateSurpassed() {
	tests := []struct {
		name         string
		update       *types.UpdateFirstDateSurpassed
		expectedDate sql.NullTime
		err          error
	}{
		{
			name: "Should succeed without any errors",
			update: &types.UpdateFirstDateSurpassed{
				ApplicationIDs:     []string{"test_app_47hfnths73j2se", "test_app_5hdf7sh23jd828"},
				FirstDateSurpassed: time.Date(2022, time.December, 13, 5, 15, 0, 0, time.UTC),
			},
			expectedDate: sql.NullTime{Valid: true, Time: time.Date(2022, time.December, 13, 5, 15, 0, 0, time.UTC)},
			err:          nil,
		},
	}

	for _, test := range tests {
		err := ts.driver.UpdateAppFirstDateSurpassed(testCtx, test.update)
		ts.Equal(test.err, err)

		for _, appID := range test.update.ApplicationIDs {
			app, err := ts.driver.SelectOneApplication(testCtx, appID)
			ts.NoError(err)
			ts.Equal(test.expectedDate.Time, app.FirstDateSurpassed.Time.UTC()) // SQL time comes back without location
		}
	}
}

func (ts *PGDriverTestSuite) Test_RemoveApplication() {
	tests := []struct {
		name           string
		appID          string
		expectedStatus string
		err            error
	}{
		{
			name:           "Should remove a single application successfully with correct input",
			appID:          "test_app_47hfnths73j2se",
			expectedStatus: "AWAITING_GRACE_PERIOD",
			err:            nil,
		},
	}

	for _, test := range tests {
		err := ts.driver.RemoveApplication(testCtx, test.appID)
		ts.Equal(test.err, err)

		appAfterRemove, err := ts.driver.SelectOneApplication(testCtx, test.appID)
		ts.Equal(test.err, err)
		ts.Equal(test.expectedStatus, appAfterRemove.Status.String)
	}
}
