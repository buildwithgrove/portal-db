package postgresdriver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pokt-foundation/portal-db/types"
	"github.com/pokt-foundation/utils-go/id"
)

var (
	ErrInvalidUsersJSON            = errors.New("error: users JSON is invalid")
	ErrUserInputIsMissingField     = errors.New("error: user access input is missing a required field")
	ErrLBMustHaveUser              = errors.New("error: a new load balancer must have at least one user")
	ErrCannotSetToOwner            = errors.New("error: load balancers may only have one owner and the owner role is already set")
	ErrCannotSetToOwnerNotAccepted = errors.New("error: cannot set a user to owner if they have not yet accepted their invitation")
)

/* ReadLoadBalancers returns all LoadBalancers in the database */
func (p *PostgresDriver) ReadLoadBalancers(ctx context.Context) ([]*types.LoadBalancer, error) {
	dbLoadBalancers, err := p.SelectLoadBalancers(ctx)
	if err != nil {
		return nil, err
	}

	var loadbalancers []*types.LoadBalancer
	for _, dbLoadBalancer := range dbLoadBalancers {
		loadBalancer, err := dbLoadBalancer.toLoadBalancer()
		if err != nil {
			return nil, err
		}

		loadbalancers = append(loadbalancers, loadBalancer)
	}

	return loadbalancers, nil
}

func (lb *SelectLoadBalancersRow) toLoadBalancer() (*types.LoadBalancer, error) {
	loadBalancer := types.LoadBalancer{
		ID:                lb.LbID,
		Name:              lb.Name.String,
		UserID:            lb.UserID.String,
		ApplicationIDs:    strings.Split(string(lb.AppIds), ","),
		RequestTimeout:    int(lb.RequestTimeout.Int32),
		Gigastake:         lb.Gigastake.Bool,
		GigastakeRedirect: lb.GigastakeRedirect.Bool,

		StickyOptions: types.StickyOptions{
			Duration:      lb.SDuration.String,
			StickyOrigins: lb.SOrigins,
			StickyMax:     int(lb.SStickyMax.Int32),
			Stickiness:    lb.SStickiness.Bool,
		},

		Integrations: types.AccountIntegrations{
			CovalentAPIKeyFree: lb.CovalentApiKeyFree.String,
			CovalentAPIKeyPaid: lb.CovalentApiKeyPaid.String,
		},

		CreatedAt: lb.CreatedAt.Time,
		UpdatedAt: lb.UpdatedAt.Time,

		// Note: here in prep for the V2 migration, needed for Covalent API keys
		AccountID: lb.AccountID.String,
	}

	// Unmarshal LoadBalancer Users JSON into []types.UserAccess
	err := json.Unmarshal(lb.Users, &loadBalancer.Users)
	if err != nil {
		return &types.LoadBalancer{}, fmt.Errorf("%w: %s", ErrInvalidUsersJSON, err)
	}

	return &loadBalancer, nil
}

/* ReadUserPermissions returns all UserPermissions in the database as a map that takes the form map[types.UserID]*types.UserPermissions */
func (p *PostgresDriver) ReadUserPermissions(ctx context.Context) (map[types.UserID]*types.UserPermissions, error) {
	userRoles, err := p.SelectUserRoles(ctx)
	if err != nil {
		return nil, err
	}

	userPermissionsMap := make(map[types.UserID]*types.UserPermissions)

	for _, userRoleRow := range userRoles {
		userID := types.UserID(userRoleRow.UserID.String)
		lbID := types.LoadBalancerID(userRoleRow.LbID)

		if userPermissions, ok := userPermissionsMap[userID]; ok {
			_, err := userPermissions.UpsertPermissions(lbID, userRoleRow.RoleName)
			if err != nil {
				return nil, err
			}
		} else {
			emptyPermissions := types.UserPermissions{
				UserID:        userID,
				LoadBalancers: map[types.LoadBalancerID]types.LoadBalancerPermissions{},
			}

			permissions, err := emptyPermissions.UpsertPermissions(lbID, userRoleRow.RoleName)
			if err != nil {
				return nil, err
			}

			userPermissionsMap[userID] = permissions
		}
	}

	return userPermissionsMap, nil
}

// generateID generates a new random account_id
// Note: here in prep for the V2 migration, needed for Covalent API keys
func (p *PostgresDriver) generateAccountID(ctx context.Context) (string, error) {
	var generatedID string
	var err error
	idExists := true

	for idExists {
		generatedID = id.GenerateID(8)
		idExists, err = p.CheckAccountIDExists(ctx, newSQLNullString(generatedID))
		if err != nil {
			return "", fmt.Errorf("error checking ID %s", generatedID)
		}
	}

	return generatedID, nil
}

/* WriteLoadBalancer saves input LoadBalancer to the database */
func (p *PostgresDriver) WriteLoadBalancer(ctx context.Context, loadBalancer *types.LoadBalancer) (*types.LoadBalancer, error) {
	if len(loadBalancer.Users) < 1 {
		return nil, ErrLBMustHaveUser
	}

	id, err := generateRandomID()
	if err != nil {
		return nil, err
	}
	loadBalancer.ID = id

	accountID, err := p.generateAccountID(ctx)
	if err != nil {
		return nil, err
	}
	loadBalancer.AccountID = accountID

	time := time.Now()
	loadBalancer.CreatedAt = time
	loadBalancer.UpdatedAt = time

	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := p.WithTx(tx)

	err = qtx.InsertLoadBalancer(ctx, extractInsertLoadBalancer(loadBalancer))
	if err != nil {
		return nil, err
	}

	stickinessParams := extractInsertStickinessOptions(loadBalancer)
	if stickinessParams.isNotNull() {
		err = qtx.InsertStickinessOptions(ctx, stickinessParams)
		if err != nil {
			return nil, err
		}
	}

	loadBalancer.Users[0].RoleName = types.RoleOwner                                                 // The first User will be the OWNER of the LoadBalancer
	err = qtx.InsertUserAccessOwner(ctx, extractInsertUserAccessOwner(id, loadBalancer, true, time)) // New LB owners always start with accepted = true
	if err != nil {
		return nil, err
	}

	lbAppParams := InsertLbAppsParams{LbID: loadBalancer.ID}
	lbAppParams.AppIds = append(lbAppParams.AppIds, loadBalancer.ApplicationIDs...)

	err = qtx.InsertLbApps(ctx, lbAppParams)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func extractInsertLoadBalancer(loadBalancer *types.LoadBalancer) InsertLoadBalancerParams {
	return InsertLoadBalancerParams{
		LbID:              loadBalancer.ID,
		Name:              newSQLNullString(loadBalancer.Name),
		UserID:            newSQLNullString(loadBalancer.UserID),
		RequestTimeout:    newSQLNullInt32(int32(loadBalancer.RequestTimeout), false),
		Gigastake:         newSQLNullBool(&loadBalancer.Gigastake),
		GigastakeRedirect: newSQLNullBool(&loadBalancer.GigastakeRedirect),
		CreatedAt:         newSQLNullTime(loadBalancer.CreatedAt),
		UpdatedAt:         newSQLNullTime(loadBalancer.UpdatedAt),
	}
}

func extractInsertStickinessOptions(loadBalancer *types.LoadBalancer) InsertStickinessOptionsParams {
	return InsertStickinessOptionsParams{
		LbID:       loadBalancer.ID,
		Duration:   newSQLNullString(loadBalancer.StickyOptions.Duration),
		Origins:    loadBalancer.StickyOptions.StickyOrigins,
		StickyMax:  newSQLNullInt32(int32(loadBalancer.StickyOptions.StickyMax), false),
		Stickiness: newSQLNullBool(&loadBalancer.StickyOptions.Stickiness),
	}
}
func (i *InsertStickinessOptionsParams) isNotNull() bool {
	return i.Duration.Valid || len(i.Origins) > 0 || i.StickyMax.Valid
}

func extractInsertUserAccessOwner(lbID string, loadBalancer *types.LoadBalancer, accepted bool, createdAt time.Time) InsertUserAccessOwnerParams {
	userAccess := loadBalancer.Users[0]

	return InsertUserAccessOwnerParams{
		LbID:      lbID,
		UserID:    newSQLNullString(userAccess.UserID),
		RoleName:  userAccess.RoleName,
		Email:     userAccess.Email,
		Accepted:  accepted,
		CreatedAt: newSQLNullTime(createdAt),
		UpdatedAt: newSQLNullTime(createdAt),
		// Note: here in prep for the V2 migration, needed for Covalent API keys
		AccountID: newSQLNullString(loadBalancer.AccountID),
	}
}

/* UpsertLoadBalancerIntegrations saves or updates input AccountIntegrations in the database */
func (p *PostgresDriver) UpsertLoadBalancerIntegrations(ctx context.Context, integrations types.AccountIntegrations) (*types.AccountIntegrations, error) {
	time := time.Now()

	accountIntegrations, err := p.UpsertAccountIntegrations(ctx, UpsertAccountIntegrationsParams{
		AccountID:          integrations.AccountID,
		CovalentApiKeyFree: newSQLNullString(integrations.CovalentAPIKeyFree),
		CovalentApiKeyPaid: newSQLNullString(integrations.CovalentAPIKeyPaid),
		CreatedAt:          newSQLNullTime(time),
		UpdatedAt:          newSQLNullTime(time),
	})
	if err != nil {
		return nil, err
	}

	return &types.AccountIntegrations{
		AccountID:          accountIntegrations.AccountID,
		CovalentAPIKeyFree: accountIntegrations.CovalentApiKeyFree.String,
		CovalentAPIKeyPaid: accountIntegrations.CovalentApiKeyPaid.String,
	}, nil
}

/* WriteLoadBalancerUser saves input UserAccess to the database */
func (p *PostgresDriver) WriteLoadBalancerUser(ctx context.Context, lbID string, userAccess types.UserAccess) error {
	if lbID == "" {
		return ErrMissingLBID
	}
	if userAccess.Email == "" {
		return ErrMissingEmail
	}
	if userAccess.RoleName == types.RoleName("") {
		return ErrMissingRole
	}
	if userAccess.RoleName == types.RoleOwner {
		return ErrCannotSetToOwner
	}

	userAccess.UserID = ""                                                           // New LB users do not start with their user ID set
	userAccessParams := extractInsertUserAccess(lbID, userAccess, false, time.Now()) // New LB users always start with accepted = false

	err := p.InsertUserAccess(ctx, userAccessParams)
	if err != nil {
		return err
	}

	return nil
}

func extractInsertUserAccess(lbID string, userAccess types.UserAccess, accepted bool, createdAt time.Time) InsertUserAccessParams {
	return InsertUserAccessParams{
		LbID:      lbID,
		UserID:    newSQLNullString(userAccess.UserID),
		RoleName:  userAccess.RoleName,
		Email:     userAccess.Email,
		Accepted:  accepted,
		CreatedAt: newSQLNullTime(createdAt),
		UpdatedAt: newSQLNullTime(createdAt),
	}
}

/* UpdateLoadBalancer updates LoadBalancer and related table rows */
func (p *PostgresDriver) UpdateLoadBalancer(ctx context.Context, id string, update *types.UpdateLoadBalancer) error {
	if id == "" {
		return ErrMissingID
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := p.WithTx(tx)

	err = qtx.UpdateLB(ctx, extractUpsertLoadBalancer(id, update))
	if err != nil {
		return err
	}

	stickinessOptionsParams := extractUpsertStickinessOptions(id, update)
	if stickinessOptionsParams.isNotNull() {
		err = qtx.UpsertStickinessOptions(ctx, *stickinessOptionsParams)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func extractUpsertLoadBalancer(id string, update *types.UpdateLoadBalancer) UpdateLBParams {
	return UpdateLBParams{
		LbID:      id,
		Name:      newSQLNullString(update.Name),
		UpdatedAt: newSQLNullTime(time.Now()),
	}
}

func extractUpsertStickinessOptions(id string, update *types.UpdateLoadBalancer) *UpsertStickinessOptionsParams {
	if update.StickyOptions == nil {
		return nil
	}

	return &UpsertStickinessOptionsParams{
		LbID:       id,
		Duration:   newSQLNullString(update.StickyOptions.Duration),
		StickyMax:  newSQLNullInt32(int32(update.StickyOptions.StickyMax), false),
		Stickiness: newSQLNullBool(update.StickyOptions.Stickiness),
		Origins:    update.StickyOptions.StickyOrigins,
	}
}
func (u *UpsertStickinessOptionsParams) isNotNull() bool {
	return u != nil && (u.Duration.Valid || u.StickyMax.Valid || u.Stickiness.Valid || len(u.Origins) != 0)
}

/* UpdateUserAccessRole updates the RoleName for a UserAccess row */
func (p *PostgresDriver) UpdateUserAccessRole(ctx context.Context, email, lbID string, roleName types.RoleName) error {
	if email == "" {
		return ErrMissingEmail
	}
	if lbID == "" {
		return ErrMissingLBID
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := p.WithTx(tx)

	// Block setting a user's role to owner if they have not yet accepted their invitation
	if roleName == types.RoleOwner {
		accepted, err := p.GetUserAccessAccepted(ctx, GetUserAccessAcceptedParams{Email: email, LbID: lbID})
		if err != nil {
			return err
		}
		if !accepted {
			return ErrCannotSetToOwnerNotAccepted
		}
	}

	params := UpdateUserAccessParams{
		Email:     email,
		LbID:      lbID,
		RoleName:  roleName,
		UpdatedAt: newSQLNullTime(time.Now()),
	}

	if roleName == types.RoleOwner {
		// Covalant Update: If the new role is owner, we need to transfer the account_id column for the
		// new owner then set the old owner_id columns account_id to null.
		ownerRow, err := qtx.GetPreviousOwner(ctx, lbID)
		if err != nil {
			return err
		}

		params.AccountID = ownerRow.AccountID

		err = qtx.UpdateAccountIDNullOldOwner(ctx,
			UpdateAccountIDNullOldOwnerParams{LbID: lbID, Email: ownerRow.Email},
		)
		if err != nil {
			return err
		}
	}

	err = qtx.UpdateUserAccess(ctx, params)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

/* AcceptUserAccess sets the User ID and the Accepted field to true for a UserAccess row */
func (p *PostgresDriver) AcceptUserAccess(ctx context.Context, email, userID, lbID string) error {
	if email == "" {
		return ErrMissingEmail
	}
	if userID == "" {
		return ErrMissingUserID
	}
	if lbID == "" {
		return ErrMissingLBID
	}

	params := SetUserAccessAcceptedParams{
		Email:     email,
		UserID:    newSQLNullString(userID),
		LbID:      lbID,
		Accepted:  true,
		UpdatedAt: newSQLNullTime(time.Now()),
	}

	err := p.SetUserAccessAccepted(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

/* RemoveLoadBalancer sets the user ID to an empty string (will not appear in Portal API or UI) */
func (p *PostgresDriver) RemoveLoadBalancer(ctx context.Context, id string) error {
	if id == "" {
		return ErrMissingID
	}

	err := p.RemoveLB(ctx, RemoveLBParams{LbID: id, UpdatedAt: newSQLNullTime(time.Now())})
	if err != nil {
		return err
	}

	return nil
}

/* RemoveUserAccess deletes a UserAccess row */
func (p *PostgresDriver) RemoveUserAccess(ctx context.Context, email, lbID string) error {
	if email == "" {
		return ErrMissingEmail
	}
	if lbID == "" {
		return ErrMissingLBID
	}

	params := DeleteUserAccessParams{Email: email, LbID: lbID}

	err := p.DeleteUserAccess(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

/* Used by Listener */
type (
	dbLoadBalancerJSON struct {
		LbID              string `json:"lb_id"`
		Name              string `json:"name"`
		UserID            string `json:"user_id"`
		RequestTimeout    int    `json:"request_timeout"`
		Gigastake         bool   `json:"gigastake"`
		GigastakeRedirect bool   `json:"gigastake_redirect"`
		AccountID         string `json:"account_id"`
		CreatedAt         string `json:"created_at"`
		UpdatedAt         string `json:"updated_at"`
	}
	dbStickinessOptionsJSON struct {
		LbID       string   `json:"lb_id"`
		Duration   string   `json:"duration"`
		Origins    []string `json:"origins"`
		StickyMax  int      `json:"sticky_max"`
		Stickiness bool     `json:"stickiness"`
	}
	dbUserAccessJSON struct {
		LbID     string `json:"lb_id"`
		UserID   string `json:"user_id"`
		RoleName string `json:"role_name"`
		Email    string `json:"email"`
		Accepted bool   `json:"accepted"`
	}
	dbAccountIntegrationsJSON struct {
		AccountID          string `json:"account_id"`
		CovalentAPIKeyFree string `json:"covalent_api_key_free"`
		CovalentAPIKeyPaid string `json:"covalent_api_key_paid"`
	}
)

func (j dbLoadBalancerJSON) toOutput() *types.LoadBalancer {
	return &types.LoadBalancer{
		ID:                j.LbID,
		Name:              j.Name,
		UserID:            j.UserID,
		RequestTimeout:    j.RequestTimeout,
		Gigastake:         j.Gigastake,
		GigastakeRedirect: j.GigastakeRedirect,
		AccountID:         j.AccountID,
		CreatedAt:         psqlDateToTime(j.CreatedAt),
		UpdatedAt:         psqlDateToTime(j.UpdatedAt),
	}
}
func (j dbStickinessOptionsJSON) toOutput() *types.StickyOptions {
	return &types.StickyOptions{
		ID:            j.LbID,
		Duration:      j.Duration,
		StickyOrigins: j.Origins,
		StickyMax:     j.StickyMax,
		Stickiness:    j.Stickiness,
	}
}
func (j dbUserAccessJSON) toOutput() *types.UserAccess {
	return &types.UserAccess{
		ID:       j.LbID,
		UserID:   j.UserID,
		RoleName: types.RoleName(j.RoleName),
		Email:    j.Email,
		Accepted: j.Accepted,
	}
}
func (j dbAccountIntegrationsJSON) toOutput() *types.AccountIntegrations {
	return &types.AccountIntegrations{
		AccountID:          j.AccountID,
		CovalentAPIKeyFree: j.CovalentAPIKeyFree,
		CovalentAPIKeyPaid: j.CovalentAPIKeyPaid,
	}
}
