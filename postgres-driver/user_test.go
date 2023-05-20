package postgresdriver

import (
	"context"
	"fmt"

	"github.com/pokt-foundation/portal-db/v2/testdata"
	"github.com/pokt-foundation/portal-db/v2/types"
)

func (ts *PGDriverTestSuite) Test_ReadUserIDsMap() {
	tests := []struct {
		name               string
		expectedUserIDsMap map[types.ProviderUserID]types.UserID
		err                error
	}{
		{
			name:               "Should return the Portal UserID when passed the auth provider user ID",
			expectedUserIDsMap: testdata.UserIDs,
			err:                nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			userIDsMap, err := ts.driver.ReadUserIDsMap(context.Background())
			ts.Equal(test.err, err)
			ts.Equal(test.expectedUserIDsMap, userIDsMap)
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
			userID: "user_1",
			user:   testdata.Users["user_1"],
			err:    nil,
		},
		{
			name:   "Should return a Portal User when passed the portal UserID",
			userID: "user_7",
			user:   testdata.Users["user_7"],
			err:    nil,
		},
		{
			name:   "Should fail if the user does not exist in the database",
			userID: "user_42",
			err:    fmt.Errorf(errUserDoesntExist.Error(), "user_42"),
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

func (ts *PGDriverTestSuite) Test_WriteNewUser() {
	tests := []struct {
		name       string
		createUser types.CreateUser
		user       *types.User
		err        error
	}{
		{
			name: "Should create a new portal User in the DB",
			createUser: types.CreateUser{
				Email:          "geralt.of.rivia623@example.com",
				ProviderUserID: "auth0|geralt_of_rivia",
			},
			user: &types.User{
				ID:       "", // user ID set in test case
				Email:    "geralt.of.rivia623@example.com",
				SignedUp: true,
				AuthProviders: map[types.AuthType]types.UserAuthProvider{
					types.AuthTypeAuth0Username: {
						ProviderUserID: "auth0|geralt_of_rivia",
						Type:           types.AuthTypeAuth0Username,
						Provider:       types.AuthProviderAuth0,
						Federated:      false,
					},
				},
				CreatedAt: testdata.MockTimestamp,
				UpdatedAt: testdata.MockTimestamp,
			},
			err: nil,
		},
		{
			name: "Should fail if an invalid email provided",
			createUser: types.CreateUser{
				Email: "jar.jar.binks3",
			},
			err: fmt.Errorf(errInvalidEmail.Error(), types.Email("jar.jar.binks3")),
		},
		{
			name: "Should fail if an invalid auth provider type provided",
			createUser: types.CreateUser{
				Email:          "jar.jar.binks3@example.com",
				ProviderUserID: "wtf|jar_jar_binks",
			},
			err: fmt.Errorf(errInvalidAuthProviderType.Error(), types.AuthType("wtf")),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			user, accountID, err := ts.driver.WriteUserNewSignUp(context.Background(), test.createUser, testdata.MockTimestamp)
			ts.Equal(test.err, err)

			if test.err == nil {
				test.user.ID = user.ID
				ts.Equal(test.user, user)

				ts.NotEmpty(accountID)

				accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{})
				ts.NoError(err)
				exists := false
				for _, account := range accounts {
					for _, user := range account.Users {
						if user.UserID == test.user.ID {
							exists = true
							ts.Equal(test.createUser.Email, user.Email)
							break
						}
					}
				}
				ts.True(exists)
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_Z_DeletePortalUser() {
	tests := []struct {
		name        string
		userID      types.UserID
		expectedErr error
		err         error
	}{
		{
			name:        "Should fail to delete a portal User from the DB if they are on the team of any accounts",
			userID:      "user_1",
			expectedErr: errUserDoesntExist,
			err:         errUserHasAccount,
		},
		{
			name:        "Should delete a portal User from the DB",
			userID:      "user_11",
			expectedErr: errUserDoesntExist,
			err:         nil,
		},
		{
			name:   "Should fail if the user does not exist in the database",
			userID: "user_42",
			err:    fmt.Errorf(errUserDoesntExist.Error(), "user_42"),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			userID, err := ts.driver.DeletePortalUser(context.Background(), test.userID)
			ts.Equal(test.err, err)

			if test.err == nil {
				_, err := ts.driver.ReadUserByUserID(context.Background(), userID)
				ts.Equal(fmt.Errorf(test.expectedErr.Error(), userID), err)
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_AllReadUserPermissions() {
	tests := []struct {
		name                    string
		expectedUserPermissions map[types.UserID]*types.UserPermissions
		err                     error
	}{
		{
			name:                    "Should read all UserPermissions for the DB as a map[types.UserID]*types.UserPermissions",
			expectedUserPermissions: testdata.UserPermissions,
			err:                     nil,
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			userPermissions, err := ts.driver.ReadUserPermissions(context.Background())
			ts.Equal(test.err, err)

			if test.err == nil {
				ts.Equal(test.expectedUserPermissions, userPermissions)
			}
		})
	}
}
