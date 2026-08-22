package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/models"
)

const textTimestampDriverName = "alice-text-timestamp-test"

func init() {
	sql.Register(textTimestampDriverName, textTimestampDriver{})
}

func TestCreateUserUsesPostgresCompatibleColumnValues(t *testing.T) {
	testDB, err := sql.Open(textTimestampDriverName, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { testDB.Close() })

	originalDB, originalDriver := DB, DriverName
	DB, DriverName = testDB, "postgres"
	t.Cleanup(func() { DB, DriverName = originalDB, originalDriver })

	user := &models.User{
		Email:        "postgres-registration@example.test",
		PasswordHash: "hash",
		FirstName:    "Postgres",
		LastName:     "Registration",
		Role:         "reader",
	}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser passed an incompatible PostgreSQL value: %v", err)
	}
}

type textTimestampDriver struct{}

func (textTimestampDriver) Open(string) (driver.Conn, error) {
	return textTimestampConn{}, nil
}

type textTimestampConn struct{}

func (textTimestampConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("Prepare is not supported")
}

func (textTimestampConn) Close() error { return nil }

func (textTimestampConn) Begin() (driver.Tx, error) {
	return nil, errors.New("Begin is not supported")
}

func (textTimestampConn) CheckNamedValue(value *driver.NamedValue) error {
	if _, ok := value.Value.(time.Time); ok {
		return errors.New("time.Time cannot be written to a TEXT column")
	}
	if _, ok := value.Value.(bool); ok {
		return errors.New("bool cannot be written to an INTEGER column")
	}
	return nil
}

func (textTimestampConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	return driver.RowsAffected(1), nil
}
