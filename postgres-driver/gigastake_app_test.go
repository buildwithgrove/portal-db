package postgresdriver

import (
	"context"

	"github.com/pokt-foundation/portal-db/v2/testdata"
	"github.com/pokt-foundation/portal-db/v2/types"
)

func (ts *PGDriverTestSuite) Test_ReadGigastakeApps() {
	tests := []struct {
		name          string
		gigastakeApps map[types.ProtocolAppID]*types.GigastakeApp
		options       types.DriverOptions
		err           error
	}{
		{
			name: "Should return all non-deleted GigastakeApps from the database",
			gigastakeApps: map[types.ProtocolAppID]*types.GigastakeApp{
				testdata.GigastakeApps["test_gigastake_app_1"].AATID: testdata.GigastakeApps["test_gigastake_app_1"],
				testdata.GigastakeApps["test_gigastake_app_2"].AATID: testdata.GigastakeApps["test_gigastake_app_2"],
				testdata.GigastakeApps["test_gigastake_app_3"].AATID: testdata.GigastakeApps["test_gigastake_app_3"],
			},
			options: types.DriverOptions{},
			err:     nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			gigastakeApps, err := ts.driver.ReadGigastakeApps(context.Background(), test.options)
			ts.Equal(test.err, err)
			ts.Equal(test.gigastakeApps, gigastakeApps)
		})
	}
}
