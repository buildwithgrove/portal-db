package postgresdriver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

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

func (ts *PGDriverTestSuite) Test_ReadPortalApps() {
	tests := []struct {
		name       string
		portalApps []*types.PortalApp
		err        error
	}{
		{
			name: "Should return all PortalApps from the database ordered by application_id",
			portalApps: []*types.PortalApp{
				{
					ID: "test_app_47hfnths73j2se",
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		portalApps, err := ts.driver.ReadPortalApps(context.Background())
		ts.Equal(test.err, err)
		ts.Equal(test.portalApps, portalApps)

		PrettyString("APP", portalApps)
	}
}
