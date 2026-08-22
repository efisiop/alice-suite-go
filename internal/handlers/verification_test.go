package handlers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

const verificationPostgresDriverName = "alice-verification-postgres-test"

func init() {
	sql.Register(verificationPostgresDriverName, verificationPostgresDriver{})
}

func TestHandleVerifyBookCodeReturnsBadRequestForInvalidCodeOnPostgres(t *testing.T) {
	withVerificationPostgresDB(t, "invalid")

	recorder := performVerificationRequest(t, "not-a-real-code")
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%q", recorder.Code, http.StatusBadRequest, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "Invalid verification code") {
		t.Fatalf("body = %q, want a clear invalid-code message", recorder.Body.String())
	}
}

func TestHandleVerifyBookCodeCompletesValidCodeOnPostgres(t *testing.T) {
	withVerificationPostgresDB(t, "valid")

	recorder := performVerificationRequest(t, "valid-code")
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%q", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"valid":true`) {
		t.Fatalf("body = %q, want successful verification", recorder.Body.String())
	}
}

func TestHandleVerifyBookCodeRejectsConsultantRole(t *testing.T) {
	withVerificationPostgresDB(t, "invalid")

	recorder := performVerificationRequestForRole(t, "not-a-real-code", "consultant")
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d; body=%q", recorder.Code, http.StatusForbidden, recorder.Body.String())
	}
}

func withVerificationPostgresDB(t *testing.T, mode string) {
	t.Helper()
	testDB, err := sql.Open(verificationPostgresDriverName, mode)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { testDB.Close() })

	originalDB, originalDriver := database.DB, database.DriverName
	database.DB, database.DriverName = testDB, "postgres"
	t.Cleanup(func() {
		database.DB, database.DriverName = originalDB, originalDriver
	})
}

func performVerificationRequest(t *testing.T, code string) *httptest.ResponseRecorder {
	return performVerificationRequestForRole(t, code, "reader")
}

func performVerificationRequestForRole(t *testing.T, code, role string) *httptest.ResponseRecorder {
	t.Helper()
	token, err := auth.GenerateJWT(role+"-1", role+"@example.test", role)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/rest/v1/rpc/verify-book-code", strings.NewReader(`{"code":"`+code+`"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	HandleVerifyBookCode(recorder, req)
	return recorder
}

type verificationPostgresDriver struct{}

func (verificationPostgresDriver) Open(mode string) (driver.Conn, error) {
	return &verificationPostgresConn{mode: mode}, nil
}

type verificationPostgresConn struct {
	mode     string
	execStep int
}

func (*verificationPostgresConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("Prepare is not supported")
}

func (*verificationPostgresConn) Close() error { return nil }

func (*verificationPostgresConn) Begin() (driver.Tx, error) {
	return verificationPostgresTx{}, nil
}

func (*verificationPostgresConn) CheckNamedValue(value *driver.NamedValue) error {
	if _, ok := value.Value.(time.Time); ok {
		return errors.New("time.Time cannot be written to a TEXT column")
	}
	return nil
}

func (conn *verificationPostgresConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	if strings.Contains(query, "?") {
		return nil, errors.New(`syntax error at or near "?"`)
	}
	if !strings.Contains(query, "FROM verification_codes WHERE code = $1") {
		return nil, errors.New("unexpected verification lookup query")
	}
	if conn.mode == "invalid" {
		return &verificationRows{columns: []string{"book_id", "is_used"}}, nil
	}
	return &verificationRows{
		columns: []string{"book_id", "is_used"},
		values:  [][]driver.Value{{"book-1", false}},
	}, nil
}

func (conn *verificationPostgresConn) ExecContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Result, error) {
	if strings.Contains(query, "?") {
		return nil, errors.New(`syntax error at or near "?"`)
	}
	if strings.Contains(query, "used_at") {
		return nil, errors.New(`column "used_at" does not exist`)
	}
	conn.execStep++
	switch conn.execStep {
	case 1:
		if !strings.Contains(query, "UPDATE verification_codes") || !strings.Contains(query, "AND is_used = 0") {
			return nil, errors.New("verification code was not claimed safely")
		}
	case 2:
		if !strings.Contains(query, "UPDATE users SET is_verified = 1") {
			return nil, errors.New("user was not marked verified")
		}
	default:
		return nil, errors.New("unexpected verification update")
	}
	return driver.RowsAffected(1), nil
}

type verificationPostgresTx struct{}

func (verificationPostgresTx) Commit() error   { return nil }
func (verificationPostgresTx) Rollback() error { return nil }

type verificationRows struct {
	columns []string
	values  [][]driver.Value
	index   int
}

func (rows *verificationRows) Columns() []string { return rows.columns }
func (rows *verificationRows) Close() error      { return nil }
func (rows *verificationRows) Next(dest []driver.Value) error {
	if rows.index >= len(rows.values) {
		return io.EOF
	}
	copy(dest, rows.values[rows.index])
	rows.index++
	return nil
}
