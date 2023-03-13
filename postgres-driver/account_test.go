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
		err             error
	}{
		{
			name:            "Should create a new Account in the database",
			ownerID:         "test_user_a06ab0cf00a714",
			account:         *testdata.Accounts[types.AccountID(4)],
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
			account:         *testdata.Accounts[types.AccountID(4)],
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
				test.account.Users = map[types.UserID]types.AccountUserAccess{
					test.ownerID: {
						User:     types.User{ID: test.ownerID, Email: testOwner.Email, AuthProvider: testOwner.AuthProvider},
						RoleName: types.RoleOwner, Accepted: true,
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
