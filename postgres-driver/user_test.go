package postgresdriver

import (
	"context"
	"fmt"

	"github.com/pokt-foundation/portal-db/v2/testdata"
	"github.com/pokt-foundation/portal-db/v2/types"
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
				Email:            "geralt.of.rivia623@example.com",
				AuthProviderType: types.AuthTypeAuth0Username,
				ProviderUserID:   "auth0|geralt_of_rivia",
			},
			user: &types.User{
				ID:       0, // user ID set in test case
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
				Email:            "jar.jar.binks3@example.com",
				AuthProviderType: types.AuthType("wrong_type"),
			},
			err: fmt.Errorf(errInvalidAuthProviderType.Error(), types.AuthType("wrong_type")),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			userID, err := ts.driver.WriteUserNewSignUp(context.Background(), test.createUser, testdata.MockTimestamp)
			ts.Equal(test.err, err)

			if test.err == nil {
				user, err := ts.driver.ReadUserByUserID(context.Background(), userID)
				ts.NoError(err)
				test.user.ID = userID
				ts.Equal(test.user, user)
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_DeletePortalUser() {
	tests := []struct {
		name                    string
		createUser              types.CreateUser
		accountID               types.AccountID
		numUsersBeforeDelete    int
		accountUsersAfterDelete map[types.UserID]types.AccountUserAccess
		userID                  types.UserID
		expectedErr             error
		err                     error
	}{
		{
			name:        "Should delete a portal User from the DB if they are not part of any accounts",
			createUser:  testdata.TestCreateUser,
			expectedErr: errUserDoesntExist,
			err:         nil,
		},
		{
			name:                 "Should delete a portal User from the DB including its account user and auth provider rows",
			accountID:            2,
			numUsersBeforeDelete: 5,
			accountUsersAfterDelete: map[types.UserID]types.AccountUserAccess{
				3: testdata.AccountUserAccess[3],
				4: testdata.AccountUserAccess[4],
				9: testdata.AccountUserAccess[9],
				2: testdata.AccountUserAccess[10],
			},
			createUser:  testdata.TestCreateUser,
			expectedErr: errUserDoesntExist,
			err:         nil,
		},
		{
			name:   "Should fail if the user does not exist in the database",
			userID: 42,
			err:    fmt.Errorf(errUserDoesntExist.Error(), 42),
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func() {
			createdUserID := test.userID

			if test.createUser.Email != "" { // if createUser set in test case use user ID from created User
				userID, err := ts.driver.WriteUserNewSignUp(context.Background(), test.createUser, testdata.MockTimestamp)
				ts.Equal(test.err, err)
				createdUserID = userID

				if test.accountID != types.AccountID(0) {
					// if accountID set in test case then test deleting a user with an account
					_, err = ts.driver.WriteAccountUser(context.Background(), types.CreateAccountUserAccess{
						AccountID: test.accountID, Email: test.createUser.Email, RoleName: types.RoleMember,
					}, testdata.MockTimestamp)
					ts.Equal(test.err, err)
					accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{})
					ts.NoError(err)
					ts.Len(accounts[test.accountID].Users, test.numUsersBeforeDelete)
				}
			}

			userID, err := ts.driver.DeletePortalUser(context.Background(), createdUserID)
			ts.Equal(test.err, err)

			if test.err == nil {
				_, err := ts.driver.ReadUserByUserID(context.Background(), userID)
				ts.Equal(fmt.Errorf(test.expectedErr.Error(), userID), err)

				if test.accountID != types.AccountID(0) {
					accounts, err := ts.driver.ReadAccounts(context.Background(), types.DriverOptions{})
					ts.NoError(err)
					ts.Equal(test.accountUsersAfterDelete, accounts[test.accountID].Users)
				}
			}
		})
	}
}

func (ts *PGDriverTestSuite) Test_ReadUserPermissions() {
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
