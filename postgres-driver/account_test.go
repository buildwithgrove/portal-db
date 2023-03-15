package postgresdriver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pokt-foundation/portal-db/testdata"
	"github.com/pokt-foundation/portal-db/types"
)

func PrettyString(label string, thing interface{}) {
	jsonThing, _ := json.Marshal(thing)
	str := string(jsonThing)

	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, []byte(str), "", "    ")
	output := prettyJSON.String()

	fmt.Println(label, output)
}

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
			account:         *testdata.Accounts[types.AccountID(5)],
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
			err:             errAccountMustHavePlanTypeSet,
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
				ts.Equal(test.err, err)
				ts.Equal(&test.account, accounts[createdAccount.ID])
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_WriteAccountUser() {
	tests := []struct {
		name                    string
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
			name: "Should create a new AccountUserAccess row in the database for a user that hasn't signed up yet",
			createAccountUser: types.CreateAccountUserAccess{
				AccountID: 4,
				Email:     "winston.smith@test.com",
				RoleName:  types.RoleAdmin,
			},
			accountUser: testdata.AccountUserAccess[14],
			accountUsersAfterCreate: map[types.UserID]types.AccountUserAccess{
				4:  testdata.AccountUserAccess[11],
				13: testdata.AccountUserAccess[14],
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			accountUser, err := ts.driver.WriteAccountUser(context.Background(), test.createAccountUser, testdata.MockTimestamp)
			ts.Equal(test.err, err)
			ts.Equal(&test.accountUser, accountUser)

			if test.err == nil {
				accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{})
				ts.Equal(test.err, err)
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
				ts.Equal(test.err, err)
				ts.Equal(test.accountUsersAfterUpdate, accounts[test.updateAccountUser.AccountID].Users)
			}
		})
	}
}

// func (ts *PGDriverTestSuite) Test_SetAccountDeleted() {
// 	tests := []struct {
// 		name                                      string
// 		deleteParams                              DeleteAccountParams
// 		accountsBeforeDelete, accountsAfterDelete map[types.AccountID]*types.Account
// 		err                                       error
// 	}{
// 		{
// 			name: "Should set a Account's deleted field to true, causing it to not appear in the ReadAccounts query",
// 			deleteParams: DeleteAccountParams{
// 				ID: testdata.Accounts[4].ID, DeletedAt: newSQLNullTime(testdata.MockTimestamp),
// 			},
// 			accountsBeforeDelete: map[types.AccountID]*types.Account{
// 				types.AccountID(1): testdata.Accounts[types.AccountID(1)],
// 				types.AccountID(2): testdata.Accounts[types.AccountID(2)],
// 				types.AccountID(3): testdata.Accounts[types.AccountID(3)],
// 				types.AccountID(4): testdata.Accounts[types.AccountID(4)],
// 			},
// 			accountsAfterDelete: map[types.AccountID]*types.Account{
// 				types.AccountID(1): testdata.Accounts[types.AccountID(1)],
// 				types.AccountID(2): testdata.Accounts[types.AccountID(2)],
// 				types.AccountID(3): testdata.Accounts[types.AccountID(3)],
// 			},
// 			err: nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		ts.Run(test.name, func() {
// 			// Check all Accounts exist before delete
// 			accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{IncludeDeleted: false})
// 			ts.Equal(test.err, err)
// 			ts.Equal(test.accountsBeforeDelete, accounts)

// 			// Delete Account
// 			err = ts.driver.SetAccountDeleted(context.Background(), test.deleteParams.ID, test.deleteParams.DeletedAt.Time)
// 			ts.Equal(test.err, err)

// 			// Check Account was deleted
// 			accounts, err = ts.driver.ReadAccounts(context.Background(), types.DriverOptions{IncludeDeleted: false})
// 			ts.Equal(test.err, err)
// 			ts.Equal(test.accountsAfterDelete, accounts)

// 			// Check Account still appears if IncludeDeleted: true
// 			accounts, err = ts.driver.ReadAccounts(context.Background(), types.DriverOptions{IncludeDeleted: true})
// 			ts.Equal(test.err, err)
// 			testDeletedApp, ok := test.accountsBeforeDelete[test.deleteParams.ID]
// 			ts.True(ok)
// 			testDeletedApp.Deleted = true
// 			ts.Equal(test.accountsBeforeDelete, accounts)
// 		})
// 	}
// }
