package postgresdriver

import (
	"context"
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

func (ts *PGDriverTestSuite) Test_SetAccountDeleted() {
	tests := []struct {
		name                                      string
		deleteParams                              DeleteAccountParams
		accountsBeforeDelete, accountsAfterDelete map[types.AccountID]*types.Account
		err                                       error
	}{
		{
			name: "Should set a Account's deleted field to true, causing it to not appear in the ReadAccounts query",
			deleteParams: DeleteAccountParams{
				ID: testdata.Accounts[4].ID, DeletedAt: newSQLNullTime(testdata.MockTimestamp),
			},
			accountsBeforeDelete: map[types.AccountID]*types.Account{
				types.AccountID(1): testdata.Accounts[types.AccountID(1)],
				types.AccountID(2): testdata.Accounts[types.AccountID(2)],
				types.AccountID(3): testdata.Accounts[types.AccountID(3)],
				types.AccountID(4): testdata.Accounts[types.AccountID(4)],
			},
			accountsAfterDelete: map[types.AccountID]*types.Account{
				types.AccountID(1): testdata.Accounts[types.AccountID(1)],
				types.AccountID(2): testdata.Accounts[types.AccountID(2)],
				types.AccountID(3): testdata.Accounts[types.AccountID(3)],
			},
			err: nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			// Check all Accounts exist before delete
			accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{IncludeDeleted: false})
			ts.Equal(test.err, err)
			ts.Equal(test.accountsBeforeDelete, accounts)

			// Delete Account
			err = ts.driver.SetAccountDeleted(context.Background(), test.deleteParams.ID, test.deleteParams.DeletedAt.Time)
			ts.Equal(test.err, err)

			// Check Account was deleted
			accounts, err = ts.driver.ReadAccounts(context.Background(), types.DriverOptions{IncludeDeleted: false})
			ts.Equal(test.err, err)
			ts.Equal(test.accountsAfterDelete, accounts)

			// Check Account still appears if IncludeDeleted: true
			accounts, err = ts.driver.ReadAccounts(context.Background(), types.DriverOptions{IncludeDeleted: true})
			ts.Equal(test.err, err)
			testDeletedApp, ok := test.accountsBeforeDelete[test.deleteParams.ID]
			ts.True(ok)
			testDeletedApp.Deleted = true
			ts.Equal(test.accountsBeforeDelete, accounts)
		})
	}
}

func (ts *PGDriverTestSuite) Test_WriteAccount() {
	tests := []struct {
		name            string
		ownerID         types.UserID
		account         types.Account
		testCreatedTime time.Time
		err             error
	}{
		{
			name:            "Should create a new Account in the database",
			ownerID:         "test_user_a06ab0cf00a714",
			account:         *testdata.Accounts[types.AccountID(5)],
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name:            "Should fail if input Account does not have a PayPlanType set",
			ownerID:         "test_user_a06ab0cf00a714",
			account:         types.Account{Plan: types.Plan{Type: ""}},
			testCreatedTime: testdata.MockTimestamp,
			err:             errAccountMustHavePlanTypeSet,
		},
		{
			name:            "Should fail if input User does not exist in the db",
			ownerID:         "sir_not_appearing_in_this_film",
			account:         *testdata.Accounts[types.AccountID(5)],
			testCreatedTime: testdata.MockTimestamp,
			err:             errUserDoesNotExist,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			createdAccount, err := ts.driver.WriteAccount(context.Background(), test.ownerID, test.account, test.testCreatedTime)
			ts.Equal(test.err, err)

			if test.err == nil {
				testOwner := testdata.Users[test.ownerID]
				test.account.ID = createdAccount.ID
				test.account.Users = map[types.Email]types.AccountUserAccess{
					testOwner.Email: {
						UserID:   test.ownerID,
						Email:    testOwner.Email,
						RoleName: types.RoleOwner,
						Accepted: true,
					},
				}
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
		accountID               types.AccountID
		accountUser             types.AccountUserAccess
		accountUsersAfterCreate map[types.Email]types.AccountUserAccess
		testCreatedTime         time.Time
		err                     error
	}{
		{
			name:        "Should create a new AccountUserAccess row in the database for an existing User",
			accountID:   1,
			accountUser: testdata.UserAccess["new_user@example.com"],
			accountUsersAfterCreate: map[types.Email]types.AccountUserAccess{
				"user1@example.com":    testdata.UserAccess["user1@example.com"],
				"user2@example.com":    testdata.UserAccess["user2@example.com"],
				"user8@example.com":    testdata.UserAccess["user8@example.com"],
				"new_user@example.com": testdata.UserAccess["new_user@example.com"],
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
		{
			name:        "Should create a new AccountUserAccess row in the database for a user that hasn't signed up yet",
			accountID:   2,
			accountUser: testdata.UserAccess["not_signed_up@example.com"],
			accountUsersAfterCreate: map[types.Email]types.AccountUserAccess{
				"user3@example.com":         testdata.UserAccess["user3@example.com"],
				"user4@example.com":         testdata.UserAccess["user4@example.com"],
				"user9@example.com":         testdata.UserAccess["user9@example.com"],
				"not_signed_up@example.com": testdata.UserAccess["not_signed_up@example.com"],
			},
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			test.accountUser.AccountID = test.accountID
			accountUser, err := ts.driver.WriteAccountUser(context.Background(), test.accountUser, test.testCreatedTime)
			ts.Equal(test.err, err)
			accountUser.AccountID = test.accountID
			ts.Equal(&test.accountUser, accountUser)

			if test.err == nil {
				accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{})
				ts.Equal(test.err, err)
				ts.Equal(test.accountUsersAfterCreate, accounts[test.accountUser.AccountID].Users)
			}
		})
	}
}
