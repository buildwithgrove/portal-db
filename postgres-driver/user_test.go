package postgresdriver

import (
	"context"
	"fmt"

	"github.com/pokt-foundation/portal-db/testdata"
	"github.com/pokt-foundation/portal-db/types"
)

func (ts *PGDriverTestSuite) Test_GetPortalUserIDFromProviderID() {
	tests := []struct {
		name           string
		providerUserID string
		portalUserID   types.UserID
		err            error
	}{
		{
			name:           "Should return the Portal UserID when passed the auth provider user ID",
			providerUserID: "auth0|paul_atreides",
			portalUserID:   2,
			err:            nil,
		},
		{
			name:           "Should return the Portal UserID when passed the auth provider user ID",
			providerUserID: "auth0|tyrion_lannister",
			portalUserID:   9,
			err:            nil,
		},
		{
			name:           "Should fail if the user does not exist in the database",
			providerUserID: "auth0|deckard_cain",
			err:            fmt.Errorf(errUserIDDoesntExist.Error(), "auth0|deckard_cain"),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			portalUserID, err := ts.driver.GetPortalUserIDFromProviderID(context.Background(), test.providerUserID)
			ts.Equal(test.err, err)
			ts.Equal(test.portalUserID, portalUserID)
		})
	}
}

func (ts *PGDriverTestSuite) Test_ReadUserByUserID() {
	tests := []struct {
		name   string
		userID types.UserID
		user   *types.User
		err    error
	}{
		{
			name:   "Should return a Portal User when passed the portal UserID",
			userID: 1,
			user:   testdata.Users[1],
			err:    nil,
		},
		{
			name:   "Should return a Portal User when passed the portal UserID",
			userID: 7,
			user:   testdata.Users[7],
			err:    nil,
		},
		{
			name:   "Should fail if the user does not exist in the database",
			userID: 42,
			err:    fmt.Errorf(errUserDoesntExist.Error(), 42),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			user, err := ts.driver.ReadUserByUserID(context.Background(), test.userID)
			ts.Equal(test.err, err)
			ts.Equal(test.user, user)
		})
	}
}
