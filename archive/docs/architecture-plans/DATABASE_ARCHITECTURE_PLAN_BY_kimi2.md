&lt;!--  ================================================================
      SQLite + Go scalability recipe
      Written by: Kimi (Moonshot AI) – 2025-12-03
      File purpose: Drop-in execution plan for Cursor IDE
================================================================= --&gt;

# cursor-sqlite-scalability.md
**Author:** Kimi (Moonshot AI) – 2025-12-03  
**Audience:** Cursor IDE user building a dual-interface Go app (Readers + Consultants) on SQLite.

--------------------------------------------------------------------
0.  One-sentence goal
--------------------------------------------------------------------
Keep 1 000+ concurrent readers and 20 consultants isolated, consistent,
and real-time while using a *single* SQLite file and zero external DB
processes.

--------------------------------------------------------------------
1.  Project layout (copy/paste into Cursor)
--------------------------------------------------------------------
db/
 ├─ migrate.go          // embeds all .sql files, runs on start-up
 ├─ 001_schema.up.sql   // this file
 ├─ 002_indexes.up.sql
 ├─ 003_wal_mode.sql
 └─ queries.go          // generated or hand-written SQL stubs

--------------------------------------------------------------------
2.  Pragmas – run once per connection
--------------------------------------------------------------------
PRAGMA journal_mode  = WAL;          -- allows concurrent readers
PRAGMA synchronous   = NORMAL;       -- good balance
PRAGMA foreign_keys  = ON;           -- mandatory integrity
PRAGMA wal_autocheckpoint = 1000;    -- keep WAL file small
PRAGMA cache_size    = -20000;       -- 20 MB page cache per conn

--------------------------------------------------------------------
3.  Logical ER diagram (ASCII)
--------------------------------------------------------------------
readers
-------
reader_id      TEXT PK  -- UUID v4
full_name      TEXT
password_hash  TEXT
created_at     INTEGER  -- unix millis
last_seen      INTEGER  -- updated every 30 s by writer goroutine

sessions
--------
session_id     TEXT PK  -- UUID v4
reader_id      FK → readers.reader_id  ON DELETE CASCADE
started_at     INTEGER
ended_at       INTEGER  -- NULL while active
ip             TEXT
user_agent     TEXT

events
------
event_id       INTEGER PK AUTOINCREMENT
reader_id      FK → readers.reader_id  ON DELETE CASCADE
session_id     FK → sessions.session_id ON DELETE CASCADE
seq            INTEGER  -- monotonic inside session
etype          TEXT     -- 'page','click','scroll','answer' …
payload        JSON     -- arbitrary small JSON
created_at     INTEGER  -- unix millis
 UNIQUE (reader_id, session_id, seq)  -- prevents dupes

reader_state   -- MATERIALISED, updated by writer goroutine
------------
reader_id  PK
session_id FK → sessions.session_id
last_event_seq INTEGER
last_page      TEXT
last_answer    JSON
updated_at     INTEGER

consultants
-----------
consultant_id  TEXT PK  -- UUID v4
email          TEXT UNIQUE
password_hash  TEXT
full_name      TEXT
role           TEXT CHECK(role IN ('admin','viewer'))
created_at     INTEGER

--------------------------------------------------------------------
4.  Indexes (002_indexes.up.sql)
--------------------------------------------------------------------
-- reader hot-path: everything by reader_id
CREATE INDEX idx_events_reader_created ON events(reader_id, created_at DESC);
CREATE INDEX idx_sessions_reader        ON sessions(reader_id, started_at DESC);

-- consultant dashboard: last state
CREATE INDEX idx_reader_state_updated   ON reader_state(updated_at DESC);

-- login look-ups
CREATE INDEX idx_consultant_email       ON consultants(email);
CREATE INDEX idx_readers_last_seen      ON readers(last_seen DESC);

--------------------------------------------------------------------
5.  Writer-goroutine contract (Go pseudo-code)
--------------------------------------------------------------------
// One single goroutine owns the *write* connection.
// All HTTP handlers send commands through a buffered channel:
type cmd struct {
    readerID string
    sql      string   -- prepared statement name
    args     []interface{}
    reply    chan sql.Result
}

// This goroutine also every 30 s:
//   1. UPDATE readers SET last_seen = NOW WHERE reader_id = ?
//   2. REPLACE INTO reader_state (…) SELECT … FROM events …
// so consultants always hit a tiny, hot-cache-friendly table.

--------------------------------------------------------------------
6.  Reader HTTP path (read-only connection)
--------------------------------------------------------------------
OpenDB("file:app.db?mode=ro&_query_only=1")  // prevents accidental writes
Queries:
  SELECT * FROM events WHERE reader_id = ? AND session_id = ?
  SELECT last_seen FROM readers WHERE reader_id = ?

--------------------------------------------------------------------
7.  Consultant HTTP path (read-only connection, but full scan allowed)
--------------------------------------------------------------------
OpenDB("file:app.db?mode=ro")
Queries:
  SELECT * FROM reader_state ORDER BY updated_at DESC LIMIT 50;
  SELECT COUNT(*) FROM readers WHERE last_seen &gt; strftime('%s','now','-5 minutes');

--------------------------------------------------------------------
8.  Migration runner (migrate.go snippet)
--------------------------------------------------------------------
//go:embed *.sql
var files embed.FS

func Migrate(db *sql.DB) error {
    if _, err := db.Exec("CREATE TABLE IF NOT EXISTS schema_version (version INTEGER PRIMARY KEY)"); err != nil {
        return err
    }
    var v int
    db.QueryRow("SELECT version FROM schema_version ORDER BY version DESC LIMIT 1").Scan(&v)

    for _, name := range []string{"001_schema.up.sql","002_indexes.up.sql","003_wal_mode.sql"} {
        if seq := extractSeq(name); seq &gt; v {
            b, _ := files.ReadFile(name)
            if _, err := db.Exec(string(b)); err != nil {
                return fmt.Errorf("migrate %s: %w", name, err)
            }
            db.Exec("INSERT INTO schema_version(version) VALUES(?)", seq)
        }
    }
    return nil
}

--------------------------------------------------------------------
9.  Capacity cheat-sheet
--------------------------------------------------------------------
- 1 000 active readers × 100 events/hour  =  2.4 M events/day  ≈  1 GB/month
- WAL file stays &lt; 50 MB with autocheckpoint every 1 000 pages
- Page cache 20 MB × 1 020 connections  ≈  20 GB RAM → still fine on a 32 GB box
- When DB &gt; 50 GB or write throughput &gt; 500 qps → split by reader_id range into second file

--------------------------------------------------------------------
10.  Daily backup one-liner
--------------------------------------------------------------------
sqlite3 app.db ".backup backup/app-$(date +%F).db"

--------------------------------------------------------------------
11.  Next steps for Cursor
--------------------------------------------------------------------
1. Copy this file into `db/001_schema.up.sql`.
2. Create `db/002_indexes.up.sql` with the indexes in §4.
3. Create `db/003_wal_mode.sql` with the PRAGMA lines in §2.
4. Paste the migrate.go snippet and run `go generate ./...` if needed.
5. Implement the single writer goroutine and the two HTTP middlewares
   (reader-ro, consultant-ro).
6. Ship a single binary; no external DB required.

End of file.