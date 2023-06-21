package postgresdriver

import (
	"testing"

	"github.com/pokt-foundation/portal-db/v2/types"
	"github.com/stretchr/testify/suite"
)

const (
	connectionString = "postgres://postgres:pgpassword@localhost:5432/postgres?sslmode=disable" // pragma: allowlist secret
)

type (
	PGDriverTestSuite struct {
		suite.Suite
		connectionString string
		driver           *PostgresDriver
	}
)

func Test_RunPGDriverSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end to end test")
	}

	testSuite := new(PGDriverTestSuite)
	testSuite.connectionString = connectionString

	suite.Run(t, testSuite)
}

// SetupSuite runs before each test suite run
func (ts *PGDriverTestSuite) SetupSuite() {
	err := ts.initPostgresDriver()
	ts.NoError(err)
}

// Initializes a real instance of the Postgres driver that connects to the test Postgres Docker container
func (ts *PGDriverTestSuite) initPostgresDriver() error {
	mainConn, err := NewPGXConnection(connectionString)
	if err != nil {
		panic(err)
	}
	listenerConn, err := NewPGXConnection(connectionString)
	if err != nil {
		panic(err)
	}

	notificationChannel := make(chan *types.Notification, 32)

	listener := NewPGXListener(listenerConn, notificationChannel)

	driver, err := NewPostgresDriver(mainConn, listener, notificationChannel)
	if err != nil {
		panic(err)
	}
	ts.driver = driver

	return nil
}
