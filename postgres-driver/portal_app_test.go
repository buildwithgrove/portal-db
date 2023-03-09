package postgresdriver

import (
	"context"

	"github.com/pokt-foundation/portal-db/testdata"
	"github.com/pokt-foundation/portal-db/types"
)

func (ts *PGDriverTestSuite) Test_ReadPortalApps() {
	tests := []struct {
		name       string
		portalApps map[types.PortalAppID]*types.PortalApp
		err        error
	}{
		{
			name: "Should return all PortalApps from the database",
			portalApps: map[types.PortalAppID]*types.PortalApp{
				testdata.TestPortalAppOne.ID: &testdata.TestPortalAppOne,
				testdata.TestPortalAppTwo.ID: &testdata.TestPortalAppTwo,
			},
			err: nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			portalApps, err := ts.driver.ReadPortalApps(context.Background())
			ts.Equal(test.err, err)
			ts.Equal(test.portalApps, portalApps)
		})
	}
}
