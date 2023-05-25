package postgresdriver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

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

func (ts *PGDriverTestSuite) Test_WriteGigastakeApp() {
	tests := []struct {
		name            string
		gigastakeApp    types.GigastakeApp
		testCreatedTime time.Time
		err             error
	}{
		{
			name:            "Should create a new GigastakeApp in the database",
			gigastakeApp:    testdata.TestCreateGigastakeApp,
			testCreatedTime: testdata.MockTimestamp,
			err:             nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			createdGigastakeApp, err := ts.driver.WriteGigastakeApp(context.Background(), test.gigastakeApp, test.testCreatedTime)
			ts.Equal(test.err, err)
			ts.NotEmpty(createdGigastakeApp.AATID)
			ts.NotEmpty(createdGigastakeApp.AAT.ID)

			// Ensure the ID and timestamps are the same as the created app
			test.gigastakeApp.AATID = createdGigastakeApp.AATID
			test.gigastakeApp.AAT.ID = createdGigastakeApp.AATID
			test.gigastakeApp.CreatedAt = createdGigastakeApp.CreatedAt
			test.gigastakeApp.UpdatedAt = createdGigastakeApp.UpdatedAt
			Plog(&test.gigastakeApp, createdGigastakeApp)
			ts.Equal(&test.gigastakeApp, createdGigastakeApp)

			gigastakeApps, err := ts.driver.ReadGigastakeApps(context.Background(), types.DriverOptions{})
			ts.NoError(err)
			ts.Equal(&test.gigastakeApp, gigastakeApps[test.gigastakeApp.AATID])
		})
	}
}

func Plog(args ...interface{}) {
	for _, arg := range args {
		var prettyJSON bytes.Buffer
		jsonArg, _ := json.Marshal(arg)
		str := string(jsonArg)
		_ = json.Indent(&prettyJSON, []byte(str), "", "    ")
		output := prettyJSON.String()

		fmt.Println(output)
	}
}
