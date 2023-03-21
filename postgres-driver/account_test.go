package postgresdriver

import (
	"context"
	"fmt"
	"time"

	"github.com/pokt-foundation/portal-db/testdata"
	"github.com/pokt-foundation/portal-db/types"
)

func (ts *PGDriverTestSuite) Test_ReadAccounts() {
	tests := []struct {
		name     string
		accounts map[types.AccountID]*types.Account
		options  types.DriverOptions
		err      error
	}{
		{
			name: "Should return all non-deleted Accounts from the database",
			accounts: map[types.AccountID]*types.Account{
				types.AccountID(1): testdata.Accounts[types.AccountID(1)],
				types.AccountID(2): testdata.Accounts[types.AccountID(2)],
				types.AccountID(3): testdata.Accounts[types.AccountID(3)],
				types.AccountID(4): testdata.Accounts[types.AccountID(4)],
				types.AccountID(5): testdata.Accounts[types.AccountID(5)],
			},
			options: types.DriverOptions{},
			err:     nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			accounts, err := ts.driver.ReadAccounts(context.Background(), test.options)
			ts.Equal(test.err, err)
			ts.Equal(test.accounts, accounts)
		})
	}
}

func (ts *PGDriverTestSuite) Test_WriteAccount() {
	tests := []struct {
		name            string
		ownerID         types.UserID
		account         types.Account
		testCreatedTime time.Time
		users           map[types.UserID]types.AccountUserAccess
		err             error
	}{
		{
			name:            "Should create a new Account in the database",
			ownerID:         1,
			account:         *testdata.Accounts[types.AccountID(6)],
			testCreatedTime: testdata.MockTimestamp,
			users: map[types.UserID]types.AccountUserAccess{
				1: {
					UserID:   testdata.Users[1].ID,
					Email:    testdata.Users[1].Email,
					RoleName: types.RoleOwner,
					Accepted: true,
					ProviderUserIDs: map[types.AuthType]string{
						types.AuthTypeAuth0Username: "auth0|james_holden",
						types.AuthTypeAuth0Github:   "github|james_holden",
					},
				},
			},
			err: nil,
		},
		{
			name:            "Should fail if input Account does not have a PayPlanType set",
			ownerID:         1,
			account:         types.Account{Plan: types.Plan{Type: ""}},
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errPayPlanDoesntExist.Error(), ""),
		},
		{
			name:            "Should fail if input Account has an invalid plan type",
			ownerID:         1,
			account:         types.Account{Plan: types.Plan{Type: types.PayPlanType("turbo_ultra_mega_plan")}},
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errPayPlanDoesntExist.Error(), types.PayPlanType("turbo_ultra_mega_plan")),
		},
		{
			name:            "Should fail if input User does not exist in the db",
			ownerID:         451,
			account:         *testdata.Accounts[types.AccountID(5)],
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errUserDoesntExist.Error(), 451),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			createdAccount, err := ts.driver.WriteAccount(context.Background(), test.ownerID, test.account, test.testCreatedTime)
			ts.Equal(test.err, err)

			if test.err == nil {
				test.account.ID = createdAccount.ID
				test.account.Users = test.users
				ts.Equal(&test.account, createdAccount)

				accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				ts.Equal(&test.account, accounts[createdAccount.ID])
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_WriteAccountUser() {
	tests := []struct {
		name                    string
		notSignedUp             bool
		createAccountUser       types.CreateAccountUserAccess
		accountUser             types.AccountUserAccess
		accountUsersAfterCreate map[types.UserID]types.AccountUserAccess
		testCreatedTime         time.Time
		err                     error
	}{
		{
			name: "Should create a new AccountUserAccess row in the database for an existing User",
			createAccountUser: types.CreateAccountUserAccess{
				AccountID: 1,
				Email:     "bernard.marx@test.com",
				RoleName:  types.RoleMember,
			},
			accountUser: testdata.AccountUserAccess[13],
			accountUsersAfterCreate: map[types.UserID]types.AccountUserAccess{
				1:  testdata.AccountUserAccess[1],
				2:  testdata.AccountUserAccess[2],
				8:  testdata.AccountUserAccess[8],
				11: testdata.AccountUserAccess[13],
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name:        "Should create a new AccountUserAccess row in the database for a user that hasn't signed up yet",
			notSignedUp: true,
			createAccountUser: types.CreateAccountUserAccess{
				AccountID: 4,
				Email:     "winston.smith@test.com",
				RoleName:  types.RoleAdmin,
			},
			accountUser: testdata.AccountUserAccess[14],
			accountUsersAfterCreate: map[types.UserID]types.AccountUserAccess{
				4: testdata.AccountUserAccess[11],
				// Winston assigned in test case after creation
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name: "Should fail if an invalid email provided",
			createAccountUser: types.CreateAccountUserAccess{
				AccountID: 4,
				Email:     "winston.smith",
				RoleName:  types.RoleAdmin,
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errInvalidEmail.Error(), types.Email("winston.smith")),
		},
		{
			name: "Should fail if account does not exist",
			createAccountUser: types.CreateAccountUserAccess{
				AccountID: 674,
				Email:     "winston.smith@test.com",
				RoleName:  types.RoleAdmin,
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errAccountDoesntExist.Error(), 674),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			accountUser, err := ts.driver.WriteAccountUser(context.Background(), test.createAccountUser, testdata.MockTimestamp)
			if test.notSignedUp {
				test.accountUser.UserID = accountUser.UserID
				test.accountUsersAfterCreate[accountUser.UserID] = test.accountUser
			}
			ts.Equal(test.err, err)

			if test.err == nil {
				ts.Equal(&test.accountUser, accountUser)

				accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				ts.Equal(test.accountUsersAfterCreate, accounts[test.createAccountUser.AccountID].Users)
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_SetAccountUserRole() {
	tests := []struct {
		name                    string
		updateAccountUser       types.UpdateAccountUserRole
		accountUsersAfterUpdate map[types.UserID]types.AccountUserAccess
		testCreatedTime         time.Time
		err                     error
	}{
		{
			name: "Should update an existing AccountUserAccess row's role to non-OWNER role",
			updateAccountUser: types.UpdateAccountUserRole{
				AccountID: 3,
				UserID:    7,
				RoleName:  types.RoleAdmin,
			},
			accountUsersAfterUpdate: map[types.UserID]types.AccountUserAccess{
				5:  testdata.AccountUserAccess[5],
				6:  testdata.AccountUserAccess[6],
				10: testdata.AccountUserAccess[12],
				7: {
					UserID:          7,
					Email:           "frodo.baggins123@test.com",
					RoleName:        types.RoleAdmin,
					Accepted:        true,
					ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|frodo_baggins"},
				},
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name: "Should update an existing AccountUserAccess row's role back to original role",
			updateAccountUser: types.UpdateAccountUserRole{
				AccountID: 3,
				UserID:    7,
				RoleName:  types.RoleMember,
			},
			accountUsersAfterUpdate: map[types.UserID]types.AccountUserAccess{
				5:  testdata.AccountUserAccess[5],
				6:  testdata.AccountUserAccess[6],
				10: testdata.AccountUserAccess[12],
				7: {
					UserID:          7,
					Email:           "frodo.baggins123@test.com",
					RoleName:        types.RoleMember,
					Accepted:        true,
					ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|frodo_baggins"},
				},
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name: "Should transfer the OWNER of an Account",
			updateAccountUser: types.UpdateAccountUserRole{
				AccountID: 2,
				UserID:    4,
				RoleName:  types.RoleOwner,
			},
			accountUsersAfterUpdate: map[types.UserID]types.AccountUserAccess{
				9: testdata.AccountUserAccess[9],
				2: testdata.AccountUserAccess[10],
				3: {
					UserID:          3,
					Email:           "ellen.ripley789@test.com",
					RoleName:        types.RoleAdmin,
					Accepted:        true,
					ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|ellen_ripley"},
				},
				4: {
					UserID:          4,
					Email:           "ulfric.stormcloak123@test.com",
					RoleName:        types.RoleOwner,
					Accepted:        true,
					ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|ulfric_stormcloak"},
				},
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name: "Should work for users that have not accepted their invite yet",
			updateAccountUser: types.UpdateAccountUserRole{
				AccountID: 3,
				UserID:    10,
				RoleName:  types.RoleAdmin,
			},
			accountUsersAfterUpdate: map[types.UserID]types.AccountUserAccess{
				5:  testdata.AccountUserAccess[5],
				6:  testdata.AccountUserAccess[6],
				7:  testdata.AccountUserAccess[7],
				10: {UserID: 10, Email: "daenerys.targaryen123@test.com", RoleName: types.RoleAdmin, Accepted: false},
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name: "Should fail is attemptong to transfer ownership to user that has not accepted their invite yet",
			updateAccountUser: types.UpdateAccountUserRole{
				AccountID: 3,
				UserID:    10,
				RoleName:  types.RoleOwner,
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errCannotTransferNotAccepted.Error(), 10, 3),
		},
		{
			name: "Should fail if User is not a member of an Account",
			updateAccountUser: types.UpdateAccountUserRole{
				AccountID: 2,
				UserID:    512,
				RoleName:  types.RoleMember,
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             fmt.Errorf(errAccountUserDoesntExist.Error(), 512, 2),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			err := ts.driver.SetAccountUserRole(context.Background(), test.updateAccountUser, test.testCreatedTime)
			ts.Equal(test.err, err)

			if test.err == nil {
				accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				ts.Equal(test.accountUsersAfterUpdate, accounts[test.updateAccountUser.AccountID].Users)
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_UpdateAcceptAccountUser() {
	tests := []struct {
		name              string
		userID            types.UserID
		acceptAccountUser types.UpdateAcceptAccountUser
		user              *types.User
		accountUsers      map[types.UserID]types.AccountUserAccess
		err               error
	}{
		{
			name:   "Should create a new UserAuthProvider for an existing user in the DB",
			userID: 10,
			acceptAccountUser: types.UpdateAcceptAccountUser{
				AccountID:        3,
				UserID:           10,
				AuthProviderType: types.AuthTypeAuth0Username,
				ProviderUserID:   "auth0|daenerys_targaryen",
			},
			user: &types.User{
				ID:       10,
				Email:    "daenerys.targaryen123@test.com",
				SignedUp: true,
				AuthProviders: map[types.AuthType]types.UserAuthProvider{
					types.AuthTypeAuth0Username: {
						ProviderUserID: "auth0|daenerys_targaryen", Provider: types.AuthProviderAuth0,
						Type: types.AuthTypeAuth0Username, Federated: false,
					},
				},
				CreatedAt: testdata.MockTimestamp,
				UpdatedAt: testdata.MockTimestamp,
			},
			accountUsers: map[types.UserID]types.AccountUserAccess{
				5: testdata.AccountUserAccess[5],
				6: testdata.AccountUserAccess[6],
				7: testdata.AccountUserAccess[7],
				10: {
					UserID: 10, Email: "daenerys.targaryen123@test.com", RoleName: types.RoleAdmin, Accepted: true,
					ProviderUserIDs: map[types.AuthType]string{types.AuthTypeAuth0Username: "auth0|daenerys_targaryen"},
				},
			},
			err: nil,
		},
		{
			name: "Should fail if an invalid auth provider type provided",
			acceptAccountUser: types.UpdateAcceptAccountUser{
				AccountID:        3,
				UserID:           10,
				AuthProviderType: types.AuthType("ask_jeeves"),
			},
			err: fmt.Errorf(errInvalidAuthProviderType.Error(), types.AuthType("ask_jeeves")),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			err := ts.driver.UpdateAcceptAccountUser(context.Background(), test.acceptAccountUser, testdata.MockTimestamp)
			ts.Equal(test.err, err)

			if test.err == nil {
				user, err := ts.driver.ReadUserByUserID(context.Background(), test.userID)
				ts.NoError(err)
				ts.Equal(test.user, user)

				accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				ts.Equal(test.accountUsers, accounts[test.acceptAccountUser.AccountID].Users)
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_SetAccountDeleted() {
	tests := []struct {
		name                    string
		accountID               types.AccountID
		numAccountsBeforeDelete int
		accountsAfterDelete     map[types.AccountID]*types.Account
		err                     error
	}{
		{
			name:                    "Should set a Account's deleted field to true, causing it to not appear in the ReadAccounts query",
			numAccountsBeforeDelete: 6,
			accountsAfterDelete: map[types.AccountID]*types.Account{
				types.AccountID(1): testdata.Accounts[types.AccountID(1)],
				types.AccountID(2): testdata.Accounts[types.AccountID(2)],
				types.AccountID(3): testdata.Accounts[types.AccountID(3)],
				types.AccountID(4): testdata.Accounts[types.AccountID(4)],
				types.AccountID(5): testdata.Accounts[types.AccountID(5)],
			},
			err: nil,
		},
		{
			name:                    "Should fail if account does not exist",
			accountID:               347,
			numAccountsBeforeDelete: 5,
			err:                     fmt.Errorf(errAccountDoesntExist.Error(), 347),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			if test.accountID == types.AccountID(0) {
				// Create test Account to delete
				createdAccount, err := ts.driver.WriteAccount(context.Background(), 1, types.Account{Plan: testdata.PayPlans["developer_plan"]}, testdata.MockTimestamp)
				ts.NoError(err)
				test.accountID = createdAccount.ID
			}

			// Check all Accounts exist before delete
			accountsBeforeDelete, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{IncludeDeleted: false})
			ts.NoError(err)
			ts.Len(accountsBeforeDelete, test.numAccountsBeforeDelete)

			// Delete Account
			err = ts.driver.SetAccountDeleted(context.Background(), test.accountID, testdata.MockTimestamp)
			ts.Equal(test.err, err)

			if test.err == nil {
				// Check Account was deleted
				accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{IncludeDeleted: false})
				ts.NoError(err)
				ts.Len(accounts, test.numAccountsBeforeDelete-1)
				ts.Equal(test.accountsAfterDelete, accounts)

				// Check Account still appears if IncludeDeleted: true
				accounts, err = ts.driver.ReadAccounts(context.Background(), types.DriverOptions{IncludeDeleted: true})
				ts.NoError(err)
				testDeletedApp, ok := accountsBeforeDelete[test.accountID]
				ts.True(ok)
				testDeletedApp.Deleted = true
				ts.Equal(accountsBeforeDelete, accounts)
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_RemoveAccountUser() {
	tests := []struct {
		name                    string
		userID                  types.UserID
		accountID               types.AccountID
		numAccountsBeforeDelete int
		accountUsersAfterDelete map[types.UserID]types.AccountUserAccess
		err                     error
	}{
		{
			name:                    "Should delete a single AccountUserAccess row",
			accountID:               1,
			numAccountsBeforeDelete: 5,
			accountUsersAfterDelete: map[types.UserID]types.AccountUserAccess{
				1: testdata.AccountUserAccess[1],
				2: testdata.AccountUserAccess[2],
				8: testdata.AccountUserAccess[8],
			},
			err: nil,
		},
		{
			name:                    "Should fail if provided a UserID that doesn't exist for the Account",
			userID:                  789,
			accountID:               3,
			numAccountsBeforeDelete: 5,
			err:                     fmt.Errorf(errAccountUserDoesntExist.Error(), 789, 3),
		},
		{
			name:      "Should fail if attempting to delete the current Account OWNER",
			userID:    1,
			accountID: 1,
			err:       fmt.Errorf(errCannotDeleteOwner.Error(), 1, 1),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			if test.err == nil {
				if test.userID == types.UserID(0) {
					// create test User to delete
					createdUser, err := ts.driver.WriteAccountUser(context.Background(), types.CreateAccountUserAccess{
						AccountID: test.accountID, Email: "hermaeus.mora@example.com", RoleName: types.RoleMember,
					}, testdata.MockTimestamp)
					ts.NoError(err)
					test.userID = createdUser.UserID
				}

				// check all Accounts exist before delete
				accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				ts.Len(accounts, test.numAccountsBeforeDelete)
			}

			err := ts.driver.RemoveAccountUser(context.Background(), test.userID, test.accountID)
			ts.Equal(test.err, err)

			if test.err == nil {
				accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				ts.Equal(test.accountUsersAfterDelete, accounts[test.accountID].Users)
			}
		})
	}
}
