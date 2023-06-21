package postgresdriver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pokt-foundation/portal-db/v2/types"
	"github.com/pokt-foundation/utils-go/id"
)

// The PostgresDriver struct satisfies the Driver interface which defines all database driver methods
type PostgresDriver struct {
	*Queries
	DB           *pgx.Conn
	notification chan *types.Notification
	listener     Listener
}

type CloudSQLConfig struct {
	DBUser                 string
	DBPassword             string
	DBName                 string
	InstanceConnectionName string
	PrivateIP              string
}

/* ---------- Postgres Connection Funcs ---------- */

/*
NewCloudSQLPGXConnection
- Creates a connection to a Cloud SQL instance using the provided CloudSQLConfig.
- Uses the Cloud SQL Connector for Go and pgx for the connection.
- Establishes a dialer with the desired options (like using private IP).
- Returns the established connection.
*/
func NewCloudSQLPGXConnection(options CloudSQLConfig) (*pgx.Conn, error) {
	ctx := context.Background()

	dsn := fmt.Sprintf("user=%s password=%s database=%s", options.DBUser, options.DBPassword, options.DBName)

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	var opts []cloudsqlconn.Option
	if options.PrivateIP != "" {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}
	dialer, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return dialer.Dial(ctx, options.InstanceConnectionName)
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("pgx.ConnectConfig: %v", err)
	}

	err = registerCustomEnumTypes(ctx, conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

/*
NewPGXConnection (mainly used for testing)
- Establishes a connection to a PostgreSQL database using the provided connection string.
- Parses the connection string into a pgx configuration object.
- Connects to the database using the created configuration.
- Registers custom enum types to the connection if any.
- Returns the established connection.
*/
func NewPGXConnection(connectionString string) (*pgx.Conn, error) {
	ctx := context.Background()

	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = registerCustomEnumTypes(ctx, conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

/*
NewPostgresDriver
- Creates an instance of PostgresDriver using the provided pgx connection, Listener, and notification channel.
- Initiates the listening of notifications in a separate goroutine.
- Returns the created PostgresDriver instance.
*/
func NewPostgresDriver(conn *pgx.Conn, listener Listener, notificationChannel chan *types.Notification) (*PostgresDriver, error) {
	driver := &PostgresDriver{
		Queries:      New(conn),
		DB:           conn,
		listener:     listener,
		notification: notificationChannel,
	}

	go func() {
		err := driver.listener.Listen(context.Background())
		if err != nil {
			fmt.Printf("Error listening for notifications: %v\n", err)
		}
	}()

	return driver, nil
}

// NotificationChannel returns receiver Notification channel
func (pg *PostgresDriver) NotificationChannel() <-chan *types.Notification {
	return pg.notification
}

// registerCustomEnumTypes registers the custom enum types used in the DB
func registerCustomEnumTypes(ctx context.Context, conn *pgx.Conn) error {
	// Add new custom enum types to this slice, including _ for array types
	dataTypeNames := []string{
		"auth_provider",
		"_auth_provider",
		"auth_type",
		"_auth_type",
		"chain_auth_type",
		"_chain_auth_type",
		"chain_check_type",
		"_chain_check_type",
		"environment",
		"_environment",
		"notification_event",
		"_notification_event",
		"notification_type",
		"_notification_type",
		"permissions",
		"_permissions",
		"whitelist_type",
		"_whitelist_type",
	}

	for _, typeName := range dataTypeNames {
		dataType, err := conn.LoadType(ctx, typeName)
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(dataType)
	}

	return nil
}

/* ---------- Postgres Driver Util Methods ---------- */
const idLength = 8

var (
	errCheckingID = errors.New("error checking if ID %s exists")
)

// generateID generates a new random ID with the specified length and checks that it doesn't already exist in the DB.
// If it already exists in the DB a new random ID is generated until one is created that does not already exist.
// Checks the following tables for existing IDs: `portal_applications`, `gigastake_applications`, `accounts` & `users`
func (pg *PostgresDriver) generateID(ctx context.Context) (string, error) {
	var generatedID string
	var err error
	idExists := true

	for idExists {
		generatedID = id.GenerateID(idLength)
		idExists, err = pg.CheckIDExists(ctx, generatedID)
		if err != nil {
			return "", fmt.Errorf(errCheckingID.Error(), generatedID)
		}
	}

	return generatedID, nil
}

func newText(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}

	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}

func newInt4(value int32, allowZero bool) pgtype.Int4 {
	if !allowZero && value == 0 {
		return pgtype.Int4{}
	}

	return pgtype.Int4{
		Int32: value,
		Valid: true,
	}
}

func newBool(value *bool) pgtype.Bool {
	if value == nil {
		return pgtype.Bool{Valid: false}
	}

	return pgtype.Bool{
		Bool:  *value,
		Valid: true,
	}
}

func newTimestamptz(value time.Time) pgtype.Timestamptz {
	if value.IsZero() {
		return pgtype.Timestamptz{}
	}

	return pgtype.Timestamptz{
		Time:  value,
		Valid: true,
	}
}

func errNoRows(err error) bool {
	return strings.Contains(err.Error(), "no rows in result set")
}
