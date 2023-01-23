package postgresdriver

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

const (
	connectionString = "postgres://postgres:pgpassword@localhost:5432/postgres?sslmode=disable" // pragma: allowlist secret
)

var (
	testCtx = context.Background()
)

type (
	PGDriverTestSuite struct {
		suite.Suite
		connectionString string
		driver           *PostgresDriver
	}
)

func Test_RunPGDriverSuite(t *testing.T) {
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
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Printf("Problem with listener, error: %s, event type: %d", err.Error(), ev)
		}
	}
	listener := pq.NewListener(ts.connectionString, 10*time.Second, time.Minute, reportProblem)

	driver, err := NewPostgresDriver(ts.connectionString, listener)
	if err != nil {
		return err
	}
	ts.driver = driver

	return nil
}
