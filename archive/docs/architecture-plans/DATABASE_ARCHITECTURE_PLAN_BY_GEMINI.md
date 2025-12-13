Here is the execution plan designed specifically for your Cursor environment. You can save the text below as a file named `DATABASE_ARCHITECTURE_PLAN_BY_GEMINI.md` in the root of your project folder.

Once saved, you can open this file in Cursor and press `Cmd+K` (or `Ctrl+K`) or use the Chat feature to reference it, telling Cursor: *"@DATABASE\_ARCHITECTURE\_PLAN\_BY\_GEMINI.md implement Step 1 and 2 now."*

-----

### `DATABASE_ARCHITECTURE_PLAN_BY_GEMINI.md`

````markdown
# Database Architecture & Execution Plan (By Gemini)

## Context & Objective
We are building a Go application with an SQLite backend serving two distinct user types:
1.  **Readers (100-1000+):** Write progress logs and read their own history.
2.  **Consultants (10-20):** Read real-time progress of all Readers.

**Critical Constraints:**
* **Concurrency:** SQLite must be configured to handle simultaneous writes (Readers) and reads (Consultants) without locking errors.
* **Isolation:** Readers must strictly only access their own data.
* **Visibility:** Consultants must access all data.

---

## Step 1: SQLite Configuration (Performance & Concurrency)

To support 1000+ users on SQLite, we must move away from the default rollback journal to **Write-Ahead Logging (WAL)**. This allows readers and writers to exist simultaneously.

**Directives for Cursor:**
Create a database initialization function in `internal/database/db.go`. It must execute the following PRAGMAs immediately upon connection:

```go
// internal/database/db.go

// Ensure the following Pragmas are executed on Open:
// 1. journal_mode = WAL (Enables concurrency)
// 2. synchronous = NORMAL (Faster writes, safe enough for WAL)
// 3. busy_timeout = 5000 (Waits 5s before throwing "database locked" error)
// 4. foreign_keys = ON (Enforce referential integrity)

const initPragmaSQL = `
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA busy_timeout = 5000;
PRAGMA foreign_keys = ON;
`
````

-----

## Step 2: Schema Design (The "Unified" Model)

We will use a single database with a Role-Based Access Control (RBAC) column.

**Directives for Cursor:**
Create a migration file `internal/database/schema.sql` with the following structure:

### 1\. Users Table (The Anchor)

Includes an indexed `last_active_at` column for high-speed "Who is online?" queries by Consultants.

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT CHECK(role IN ('reader', 'consultant')) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_active_at DATETIME -- Indexed for real-time dashboards
);

CREATE INDEX idx_users_last_active ON users(last_active_at);
```

### 2\. Sessions (Auth & Heartbeat)

```sql
CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### 3\. Progress Milestones (State Tracking)

Stores the *current status* of a reader. Used for Consultant Dashboards.

```sql
CREATE TABLE milestones (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    chapter_id TEXT NOT NULL,
    status TEXT CHECK(status IN ('started', 'completed', 'stuck')),
    score INTEGER,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, chapter_id) -- Prevent duplicate entries for same chapter
);
```

### 4\. Activity Logs (Event Stream)

Stores the high-volume interaction history (clicks, logins, pauses).

```sql
CREATE TABLE activity_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    action_type TEXT NOT NULL, -- e.g., 'LOGIN', 'PAGE_VIEW', 'SUBMIT'
    metadata JSON,             -- SQLite supports JSON operations
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

-----

## Step 3: Go Application Logic & Security Layers

We will use the **Repository Pattern** to separate the database queries from the HTTP handlers. This ensures we can enforce logic constraints.

**Directives for Cursor:**
Structure the Go application as follows:

### 1\. The Store Interface

Define an interface in `internal/store/store.go` that abstracts the DB operations.

```go
type Store interface {
    CreateUser(ctx context.Context, user *User) error
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    
    // Updates last_active_at
    UpdateHeartbeat(ctx context.Context, userID int64) error 

    // Reader Ops (Strictly scoped by userID)
    GetMyMilestones(ctx context.Context, userID int64) ([]Milestone, error)
    LogActivity(ctx context.Context, log *ActivityLog) error

    // Consultant Ops (Can query by any targetUserID)
    GetAllActiveUsers(ctx context.Context, timeWindow time.Duration) ([]User, error)
    GetReaderMilestones(ctx context.Context, targetUserID int64) ([]Milestone, error)
}
```

### 2\. Middleware Logic (Context Injection)

Create middleware that:

1.  Validates the Session Token.
2.  Extracts `user_id` and `role`.
3.  Injects them into the Request Context.
4.  **Crucial:** Triggers a "Fire and Forget" update to `last_active_at` so Consultants see the user is online.

### 3\. Service Layer Security (The Guardrail)

In the implementation of the handlers, enforce the isolation:

  * **If Role == Reader:** The handler *automatically* takes the `user_id` from the Context and passes it to `GetMyMilestones`. The Reader *cannot* supply a target user ID in the request body.
  * **If Role == Consultant:** The handler accepts a `target_user_id` URL parameter (e.g., `/api/readers/{id}/progress`) and passes it to `GetReaderMilestones`.

-----

## Step 4: Real-Time Polling Strategy

Since we are using SQLite (and not using complex WebSockets yet), the Consultant Frontend should implement **Short Polling** to the `last_active_at` endpoint.

**Directives for Cursor:**
Create a specific endpoint `GET /api/consultant/dashboard/live`.

  * **Query:** `SELECT * FROM users WHERE last_active_at > datetime('now', '-1 minute') AND role = 'reader'`
  * **Performance:** Because `last_active_at` is indexed (Step 2), this query is instantaneous even with 1000 users.

-----

## Execution Checklist

  - [ ] Initialize `go.mod`.
  - [ ] Install SQLite driver (`go get github.com/mattn/go-sqlite3` or `modernc.org/sqlite`).
  - [ ] Implement **Step 1** (DB Init with WAL Mode).
  - [ ] Implement **Step 2** (SQL Schema migration).
  - [ ] Implement **Step 3** (Go Structs and Store Interface).
  - [ ] Create Middleware for Auth & Heartbeat.
  - [ ] build `main.go` wiring it all together.

<!-- end list -->

```

---

### Next Step for You

You can now copy the content above, create the file in your Cursor project, and ask Cursor:

**"Read DATABASE_ARCHITECTURE_PLAN_BY_GEMINI.md and strictly follow Step 1 and Step 2 to generate the initial database setup code."**

Would you like me to generate the specific **Go code** for the "Heartbeat Middleware" right now, as that is the trickiest part to get right for performance?
```