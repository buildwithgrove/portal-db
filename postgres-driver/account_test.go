package postgresdriver

import (
	"context"

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
				types.AccountID(1): testdata.TestAccounts[types.AccountID(1)],
				types.AccountID(2): testdata.TestAccounts[types.AccountID(2)],
				types.AccountID(3): testdata.TestAccounts[types.AccountID(3)],
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
