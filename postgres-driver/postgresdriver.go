package postgresdriver

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	// PQ import is required
	_ "github.com/lib/pq"
	"github.com/pokt-foundation/portal-db/v2/types"
)

const idLength = 24

var (
	ErrMissingID     = errors.New("missing id")
	ErrMissingLBID   = errors.New("missing load balancer id")
	ErrMissingUserID = errors.New("missing user id")
	ErrMissingEmail  = errors.New("missing user email")
	ErrMissingRole   = errors.New("missing user role")
)

// The PostgresDriver struct satisfies the Driver interface which defines all database driver methods
type PostgresDriver struct {
	*Queries
	DB           *sql.DB
	notification chan *types.Notification
	listener     Listener
}

/* NewPostgresDriver returns PostgresDriver instance from Postgres connection string */
func NewPostgresDriver(connectionString string, listener Listener) (*PostgresDriver, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	driver := &PostgresDriver{
		Queries:      New(db),
		DB:           db,
		notification: make(chan *types.Notification, 32),
		listener:     listener,
	}

	err = driver.listener.Listen("events")
	if err != nil {
		return nil, err
	}

	go Listen(driver.listener.NotificationChannel(), driver.notification)

	return driver, nil
}

/* NewPostgresDriverFromDBInstance returns PostgresDriver instance from sdl.DB instance */
// mostly used for mocking tests
func NewPostgresDriverFromDBInstance(db *sql.DB, listener Listener) *PostgresDriver {
	driver := &PostgresDriver{
		Queries:      New(db),
		DB:           db,
		notification: make(chan *types.Notification, 32),
		listener:     listener,
	}

	err := driver.listener.Listen("events")
	if err != nil {
		panic(err)
	}

	go Listen(driver.listener.NotificationChannel(), driver.notification)

	return driver
}

/* NotificationChannel returns receiver Notification channel  */
func (d *PostgresDriver) NotificationChannel() <-chan *types.Notification {
	return d.notification
}

func generatePortalAppID() string {
	bytes := make([]byte, idLength/2)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func newSQLNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func newSQLNullInt32(value int32, allowZero bool) sql.NullInt32 {
	if !allowZero && value == 0 {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Int32: value,
		Valid: true,
	}
}

func newSQLNullBool(value *bool) sql.NullBool {
	if value == nil {
		return sql.NullBool{Valid: false}
	}

	return sql.NullBool{
		Bool:  *value,
		Valid: true,
	}
}

func newSQLNullTime(value time.Time) sql.NullTime {
	if value.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  value,
		Valid: true,
	}
}
