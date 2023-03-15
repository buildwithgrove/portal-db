package postgresdriver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pokt-foundation/portal-db/types"
)

type (
	authProviderDBRow struct {
		UserID         types.UserID       `json:"user_id"`
		Type           types.AuthType     `json:"type"`
		Provider       types.AuthProvider `json:"provider"`
		ProviderUserID string             `json:"provider_user_id"`
		Federated      bool               `json:"federated"`
	}
)

var (
	errUserIDDoesntExist       = errors.New("error user ID does not exist for auth provider ID '%s'")
	errUserDoesntExist         = errors.New("error user does not exist for portal ID '%d'")
	errNotValidEmail           = errors.New("error email input is not a valid email address '%s'")
	errInvalidAuthProviderType = errors.New("error invalid auth provider type '%s'")
)

// /* ----- postgresdriver Account Read Methods ----- */

// GetPortalUserIDFromProviderID takes a user's auth provider ID and returns the Portal UserID
func (pg *PostgresDriver) GetPortalUserIDFromProviderID(ctx context.Context, providerUserID string) (types.UserID, error) {
	userID, err := pg.GetPortalUserID(ctx, providerUserID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return types.UserID(0), fmt.Errorf(errUserIDDoesntExist.Error(), providerUserID)
		default:
			return types.UserID(0), err
		}
	}

	return userID, nil
}

// ReadUserByUserID takes a portal UserID and returns a single user in the database as a User struct
func (pg *PostgresDriver) ReadUserByUserID(ctx context.Context, userID types.UserID) (*types.User, error) {
	userData, err := pg.GetUserDataFromPortalUserID(ctx, userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, fmt.Errorf(errUserDoesntExist.Error(), userID)
		default:
			return nil, err
		}
	}

	user, err := userData.toUser()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// toUser converts User SELECT output to User struct
func (u *GetUserDataFromPortalUserIDRow) toUser() (*types.User, error) {
	var providerRows []authProviderDBRow
	if err := json.Unmarshal(u.AuthProviders, &providerRows); err != nil {
		return nil, err
	}

	authProviders := make(map[types.AuthType]types.UserAuthProvider, len(providerRows))
	for _, provider := range providerRows {
		authProviders[provider.Type] = types.UserAuthProvider{
			ProviderUserID: provider.ProviderUserID,
			Type:           provider.Type,
			Provider:       provider.Provider,
			Federated:      provider.Federated,
		}
	}

	return &types.User{
		ID:            u.ID,
		Email:         u.Email,
		SignedUp:      u.SignedUp,
		AuthProviders: authProviders,
		CreatedAt:     u.CreatedAt.UTC(),
		UpdatedAt:     u.UpdatedAt.UTC(),
	}, nil
}

/* ----- postgresdriver User Create Methods ----- */

// WriteUserNewSignUp creates a new portal User in the DB from a CreateUser input when a new user signs up
// Includes the user's auth method
func (pg *PostgresDriver) WriteUserNewSignUp(ctx context.Context, user types.CreateUser, createdAt time.Time) (types.UserID, error) {
	if !user.Email.IsValid() {
		return types.UserID(0), fmt.Errorf(errNotValidEmail.Error(), user.Email)
	}
	if !user.AuthProviderType.IsValid() {
		return types.UserID(0), fmt.Errorf(errInvalidAuthProviderType.Error(), user.AuthProviderType)
	}

	params := CreateUserNewSignUpParams{
		Email:          user.Email,
		ProviderUserID: user.ProviderUserID,
		Type:           user.AuthProviderType,
		Provider:       user.AuthProviderType.Provider(),
		Federated:      user.AuthProviderType.IsFederated(),
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
	}

	userID, err := pg.CreateUserNewSignUp(ctx, params)
	if err != nil {
		return types.UserID(0), err
	}

	return types.UserID(userID), nil
}
