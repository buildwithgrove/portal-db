package postgresdriver

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	// PQ import is required
	_ "github.com/lib/pq"
	"github.com/pokt-foundation/portal-db/types"
)

const (
	psqlDateLayout = "2006-01-02T15:04:05.999999"
	idLength       = 24
)

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
	db           *sql.DB
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
		db:           db,
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

func generateRandomID() (string, error) {
	bytes := make([]byte, idLength/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
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

func psqlDateToTime(rawDate string) time.Time {
	date, _ := time.Parse(psqlDateLayout, rawDate)
	return date
}

func boolPointer(value bool) *bool {
	return &value
}

// Typeguard for a derived field from Postgres that must be either nil or a string
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return v.(string)
}
