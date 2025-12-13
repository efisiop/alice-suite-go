# Cursor App — SQLite Execution Plan (from **GPT-5 Thinking mini**)

> This document is an execution-ready guide for structuring the Cursor application's SQLite schema and operational patterns. It covers schema DDL, PRAGMAs, indexes, Go usage notes, real-time patterns for consultant monitoring, migrations, backups, scaling advice and an operational checklist.

---

## 1. Quick Rationale & design goals (short)

* **Isolation & integrity:** Every reader has private account data; consultants can view all readers but not modify reader-restricted fields without permission.
* **Auditability & append-only history:** Keep `interactions` (events) append-only for auditing and analytics.
* **SQLite-first, migration-ready:** Use portable SQL and keep patterns compatible with server DBs to allow later migration to Postgres/MySQL.
* **Concurrency-aware:** Use WAL mode + tuned pragmas, serialized writers or writer pool in Go, and denormalized `reader_states` for fast consultant queries.
* **Real-time UX:** Use an external pub/sub (preferred) or an efficient event-notification pattern to stream events to consultants.

---

## 2. High-level schema overview

Tables included:

* `users` — reader accounts
* `consultants` — consultant accounts
* `sessions` — active sessions (readers & consultants share schema via `actor_type`)
* `interactions` — append-only events (primary reader actions)
* `reader_states` — denormalized quick-access state per reader (kept updated by app logic)
* `messages` — messages/notes between consultant and reader
* `audit_logs` — admin/audit records
* `schema_migrations` — migration history

All tables use UUIDs (TEXT) as primary keys to keep distributed-unique identifiers.

---

## 3. Recommended SQLite PRAGMAs (set on DB open)

Run these as the first statements upon opening the DB connection (in Go, right after opening):

```sql
PRAGMA journal_mode = WAL; -- enables concurrent readers + writers
PRAGMA synchronous = NORMAL; -- good tradeoff: durability vs speed
PRAGMA foreign_keys = ON; -- enforce FK constraints
PRAGMA temp_store = MEMORY; -- speed up temp operations
PRAGMA cache_size = -2000; -- ~2000 pages (adjust for host RAM)
PRAGMA locking_mode = NORMAL; -- default; keep unless tuning for extreme use
```

Notes:

* WAL mode allows many concurrent readers and a single writer without heavy blocking.
* `synchronous = NORMAL` provides good throughput; set to `FULL` if you need max durability.

---

## 4. Full SQL schema (SQLite-compatible)

> Save these statements in a `schema.sql` file and apply via a migration tool or as part of app startup migration logic.

```sql
-- schema_migrations: simple migration tracker
CREATE TABLE IF NOT EXISTS schema_migrations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  applied_at TEXT NOT NULL
);

-- users (readers)
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY, -- UUID v4 or ULID
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  display_name TEXT,
  metadata JSON DEFAULT ('{}'), -- SQLite has no native JSON type but stores as TEXT; keep JSON functions in mind
  created_at TEXT NOT NULL,
  updated_at TEXT,
  is_active INTEGER NOT NULL DEFAULT 1
);

-- consultants
CREATE TABLE IF NOT EXISTS consultants (
  id TEXT PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  display_name TEXT,
  role TEXT NOT NULL DEFAULT 'consultant', -- e.g., admin, supervisor, consultant
  metadata JSON DEFAULT ('{}'),
  created_at TEXT NOT NULL,
  updated_at TEXT,
  is_active INTEGER NOT NULL DEFAULT 1
);

-- sessions (generic for both readers and consultants)
CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY,
  actor_id TEXT NOT NULL,
  actor_type TEXT NOT NULL CHECK(actor_type IN ('user','consultant')),
  token_hash TEXT NOT NULL,
  created_at TEXT NOT NULL,
  last_active_at TEXT,
  expires_at TEXT,
  FOREIGN KEY(actor_id) REFERENCES users(id) ON DELETE CASCADE
);

-- append-only interactions/events
CREATE TABLE IF NOT EXISTS interactions (
  id TEXT PRIMARY KEY,
  reader_id TEXT NOT NULL,
  type TEXT NOT NULL, -- e.g., 'progress', 'answer', 'heartbeat', 'open_page'
  payload JSON NOT NULL, -- store details, e.g., {"step": 3, "answer":"..."}
  created_at TEXT NOT NULL,
  processed INTEGER NOT NULL DEFAULT 0, -- app can mark as processed for background jobs
  FOREIGN KEY(reader_id) REFERENCES users(id) ON DELETE CASCADE
);

-- denormalized current state per reader (for fast consultant queries)
CREATE TABLE IF NOT EXISTS reader_states (
  reader_id TEXT PRIMARY KEY,
  last_interaction_id TEXT,
  last_seen_at TEXT,
  progress REAL DEFAULT 0, -- 0.0 - 1.0
  current_step INTEGER DEFAULT 0,
  status TEXT DEFAULT 'idle', -- 'idle', 'active', 'paused', 'completed'
  metadata JSON DEFAULT ('{}'),
  FOREIGN KEY(reader_id) REFERENCES users(id) ON DELETE CASCADE
);

-- messages between consultant and reader
CREATE TABLE IF NOT EXISTS messages (
  id TEXT PRIMARY KEY,
  reader_id TEXT NOT NULL,
  consultant_id TEXT, -- nullable if system message
  direction TEXT NOT NULL CHECK(direction IN ('user_to_consultant','consultant_to_user','system')),
  body TEXT NOT NULL,
  metadata JSON DEFAULT ('{}'),
  created_at TEXT NOT NULL,
  read_at TEXT,
  FOREIGN KEY(reader_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(consultant_id) REFERENCES consultants(id) ON DELETE SET NULL
);

-- audit logs
CREATE TABLE IF NOT EXISTS audit_logs (
  id TEXT PRIMARY KEY,
  actor_id TEXT,
  actor_type TEXT,
  action TEXT NOT NULL,
  target_table TEXT,
  target_id TEXT,
  details JSON DEFAULT ('{}'),
  created_at TEXT NOT NULL
);
```

Notes:

* All `TEXT` timestamps should be stored in ISO8601 UTC format: `YYYY-MM-DDTHH:MM:SSZ` (or with fractional seconds if needed).
* `metadata`/`payload` fields use JSON (store as TEXT). Keep JSON structure and keys stable.

---

## 5. Indexes (for performance)

Create indexes for the most common queries (consultant dashboards, recent activity, reader-specific fetches):

```sql
CREATE INDEX IF NOT EXISTS idx_interactions_reader_created_at ON interactions(reader_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_interactions_created_at ON interactions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_reader_states_last_seen ON reader_states(last_seen_at);
CREATE INDEX IF NOT EXISTS idx_sessions_actor ON sessions(actor_id, actor_type, last_active_at);
CREATE INDEX IF NOT EXISTS idx_messages_reader_created_at ON messages(reader_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_created_at ON audit_logs(created_at DESC);
```

Index rationale:

* `interactions(reader_id, created_at DESC)` is the primary index consultants will use to show recent activity per reader.
* `interactions(created_at DESC)` helps global recent activity streams.
* `reader_states(last_seen_at)` supports queries like "who's active now?" when combined with `last_seen_at > now - 2m`.

---

## 6. Go (golang) operational patterns

### 6.1 DB driver & connection settings

* Use `github.com/mattn/go-sqlite3` or `modernc.org/sqlite` (pure-Go) drivers. `mattn/go-sqlite3` is widely used but cgo-based.
* Open the DB with a reasonable connection pool size. **Important:** SQLite does not benefit from many concurrent writer connections. Instead, use:

  * `MaxOpenConns = N` (set to a low number for writers) and `SetMaxIdleConns`.
  * A typical pattern: `db.SetMaxOpenConns(1)` if your app handles writes via a single writer goroutine, or `db.SetMaxOpenConns(10)` if you have a buffered writer pool with retries.

Example:

```go
import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "time"
)

func openDB(path string) (*sql.DB, error) {
  db, err := sql.Open("sqlite3", path+"?_foreign_keys=1")
  if err != nil { return nil, err }
  db.SetMaxOpenConns(3) // conservative; tune based on app load
  db.SetConnMaxLifetime(time.Minute * 10)
  return db, nil
}
```

### 6.2 Transactions & retry logic

* Wrap writes in short-lived transactions. If you get `SQLITE_BUSY`, retry with backoff (e.g., exponential backoff with jitter) a small number of times.
* Use `BEGIN IMMEDIATE` for writes that must avoid waiting for other writers; otherwise `BEGIN` is fine.

Example pseudocode for retrying a transaction:

```go
for attempt := 0; attempt < maxRetries; attempt++ {
  tx, err := db.Begin()
  if err != nil { /* handle */ }
  _, err = tx.Exec("INSERT INTO interactions (...) VALUES (...)", ...)
  if err == nil {
    err = tx.Commit()
    if err == nil { break }
  }
  tx.Rollback()
  if isSqliteBusy(err) { sleepWithBackoff(attempt) ; continue }
  return err
}
```

### 6.3 Single writer or writer pool

* **Option A (recommended for simplicity):** Single writer goroutine. All writes are sent to a channel and the goroutine serializes them into DB transactions. This eliminates concurrency write contention.
* **Option B:** Small pool (2-5) of writers with retry on `SQLITE_BUSY`.

### 6.4 Prepared statements & batching

* Use prepared statements for frequent writes (e.g., insert interactions).
* Batch writes where possible (e.g., insert multiple interactions in one transaction) for throughput.

### 6.5 JSON handling

* Store JSON as strings. Use Go structs + `encoding/json` to marshal/unmarshal. Keep schema of `payload` stable.

---

## 7. Real-time consultant monitoring

Consultants need near-real-time visibility of reader activity. Two recommended architectures:

### 7.1 Preferred: App + Pub/Sub + WebSockets

* Flow:

  1. Reader action arrives at app server -> app writes compact event into `interactions` table.
  2. App publishes a small event message to a pub/sub broker (Redis Pub/Sub, NATS, or Kafka for heavy scale).
  3. Consultant-facing service(s) subscribe to channels and push events to consultant clients over WebSockets (or server-sent events).

* Advantages:

  * Low latency, scalable, decoupled from DB write speed.
  * Consultants get events as soon as processed by app logic.

* Implementation notes:

  * Event message should include minimal fields: `{event_id, reader_id, event_type, created_at, summary}`.
  * Use Redis Streams or Redis Pub/Sub for simple setups. Redis Streams provide persistence.

### 7.2 SQLite-only fallback: event table + notification

* Flow:

  1. App writes into `interactions` table.
  2. The app also writes to an in-memory notification broker (within the same process) to push to consultant websockets.
  3. If multiple app instances exist without a shared broker, use short polling or a lightweight message queue (recommended to add Redis/NATS).

* Notes:

  * Polling `SELECT * FROM interactions WHERE created_at > ? ORDER BY created_at` every 500-2000ms may be acceptable at small scale.
  * Use `reader_states` to reduce data transferred: publish only changed `reader_states` rather than full interaction payloads.

### 7.3 Consultant queries examples

* Active readers (last seen within 2 minutes):

```sql
SELECT reader_id, progress, current_step, last_seen_at
FROM reader_states
WHERE last_seen_at >= datetime('now', '-2 minutes')
ORDER BY last_seen_at DESC;
```

* Recent interactions globally (most recent 100):

```sql
SELECT id, reader_id, type, json_extract(payload, '$.summary') as summary, created_at
FROM interactions
ORDER BY created_at DESC
LIMIT 100;
```

* Reader-specific recent interactions:

```sql
SELECT id, type, payload, created_at
FROM interactions
WHERE reader_id = ?
ORDER BY created_at DESC
LIMIT 200;
```

---

## 8. Security checklist

* **Passwords:** Hash with `bcrypt` (cost depending on environment). Prefer `argon2` if available.
* **Sessions:** use random opaque tokens; store `token_hash` in DB, not raw token. Expire sessions and rotate tokens on sensitive changes.
* **Transport:** Serve all traffic over TLS; WebSockets must be wss://.
* **RBAC:** Enforce at application layer. Consultants can query any `reader_states` or `interactions`, readers only their own `interactions` unless consented.
* **Input validation:** Validate JSON payload fields and length. Cap payload sizes.
* **Audit logs:** Write to `audit_logs` on privilege changes, consultant actions (e.g., messages, notes), and admin actions.

---

## 9. Backups, migrations & upgrades

### 9.1 Backups

* Use `sqlite3_backup` API (supported by drivers) for hot backups. Don't copy DB file while writers are active unless using `sqlite3_backup` or `VACUUM INTO` (SQLite >= 3.27).
* Example: periodically run a snapshot backup to S3/remote storage.

### 9.2 Schema migrations

* Use `schema_migrations` table to track applied migration files.
* Apply migrations in transactions (`BEGIN IMMEDIATE` then DDL then `COMMIT`).
* For heavy migrations (large table changes), create new tables, backfill in background, then swap.

### 9.3 Archival

* Archive old `interactions` > retention window (e.g., older than 90 days) into separate archive DB files (per month or per quarter). This reduces active DB size and improves query speed.
* Keep an archive index of where records live.

---

## 10. Scaling roadmap (0 → 1000+ readers and beyond)

### Phase A — Small (0–1k readers)

* Single app instance. SQLite WAL mode. Single-writer pattern in Go.
* Use in-memory pub/sub within process for WebSocket notifications.
* Daily backups to remote.

### Phase B — Medium (1k–10k readers, multiple app instances)

* Introduce Redis (or NATS) for pub/sub and ephemeral session store.
* Keep SQLite for persistence per instance only if reads/writes remain low — prefer central DB.
* Start planning migration to Postgres if write concurrency increases.

### Phase C — Large (10k+ readers, heavy concurrent writes)

* Migrate to a client-server DB (Postgres) and centralize storage.
* Use horizontal scaling for app servers, Redis/NATS for pub/sub, and a stateless WebSocket cluster or socket gateway.
* Partition or shard data logically if needed (e.g., multi-tenant DB per customer).

Migration note: keep schema normalized and SQL-compatible to ease transition.

---

## 11. Operational notes & monitoring

* Monitor `busy` and `lock` errors in app logs; high rates indicate write contention.
* Monitor DB file size; VACUUM occasionally (but not frequently on heavy write DBs as it can be expensive).
* Track slow queries via a query logging wrapper in Go, and add indexes where necessary.
* Observe memory and pool settings; tune `PRAGMA cache_size` per host.

---

## 12. Sample queries & useful utilities

* Get readers currently active in last X minutes:

```sql
SELECT u.id, u.email, rs.progress, rs.status, rs.last_seen_at
FROM users u
JOIN reader_states rs ON rs.reader_id = u.id
WHERE rs.last_seen_at >= datetime('now', '-5 minutes')
ORDER BY rs.last_seen_at DESC;
```

* Promote a user to consultant (admin action):

```sql
BEGIN;
INSERT INTO consultants (id, email, password_hash, display_name, created_at)
SELECT id, email, password_hash, display_name, datetime('now') FROM users WHERE id = ?;
COMMIT;
```

* Compact old interactions into archive DB (conceptual):

  1. Create `archive_YYYY_QX.db` and copy range of `interactions` into it.
  2. Delete the archived rows from production DB within a transaction.

---

## 13. Appendix — example JSON payload shapes

* `interactions.payload` examples:

```json
{ "type": "progress", "step": 5, "percent": 0.52 }
{ "type": "answer", "question_id": "q_123", "answer": "B", "time_ms": 4200 }
{ "type": "heartbeat", "client_version": "1.2.3" }
```

Design tip: reserve a `summary` short field inside payload for quick displays (e.g., `{"summary":"Answered q_123: B"}`) so consultant UI can show short lists without parsing large JSON bodies.

---

## 14. Checklist before rollout

* [ ] Implement PRAGMAs at DB open
* [ ] Create schema via migration tool
* [ ] Implement single-writer goroutine (or small writer pool) with retry/backoff
* [ ] Add pub/sub (Redis or NATS) for real-time notifications
* [ ] Secure password hashing and session token practices
* [ ] Implement audit logging for consultant actions
* [ ] Setup backup snapshots and retention policy
* [ ] Prepare migration plan to Postgres if needed

---

## 15. Final notes & next steps

If you want, I can:

* Produce a ready-to-run `schema.sql` file (this document's SQL is already runnable) and a `migrate.sh` script.
* Generate Go code templates for DB access (models, prepared statements, writer goroutine, websocket notification wiring).
* Draft a lightweight Redis-based event publisher/subscriber example in Go.

Say which of the above you want next and I'll include it as a follow-up file or code doc.

---

*End of document.*
