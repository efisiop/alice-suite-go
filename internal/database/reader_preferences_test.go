package database

import (
	"os"
	"strings"
	"testing"
)

func TestReaderPreferencesPostgresTimestampDefaultsAreText(t *testing.T) {
	if got := readerPreferencesTimestampDefault("postgres"); got != "CAST(CURRENT_TIMESTAMP AS TEXT)" {
		t.Fatalf("PostgreSQL timestamp default = %q", got)
	}

	migration, err := os.ReadFile("../../migrations/017_reader_preferences.sql")
	if err != nil {
		t.Fatal(err)
	}
	postgresSQL := TransformSQLiteToPostgres(string(migration))
	if count := strings.Count(postgresSQL, "DEFAULT (CAST(CURRENT_TIMESTAMP AS TEXT))"); count != 2 {
		t.Fatalf("migration has %d PostgreSQL-safe text timestamp defaults, want 2:\n%s", count, postgresSQL)
	}
}
