# Why localhost broke with your proxy (and what fixed it)

## The change that caused the problem

**Commit:** `f5841d4`  
**Date:** March 4, 2026  
**Title:** `feat: migrate from SQLite to PostgreSQL for Render deployment`

That commit added **PostgreSQL support** so the app can run on Render.com with a real database. To do that, it added two new Go dependencies:

| Dependency | Purpose |
|------------|--------|
| `github.com/jackc/pgx/v5` | PostgreSQL driver (used when `DATABASE_URL` is set, e.g. on Render) |
| `github.com/jmoiron/sqlx` | Query helpers for PostgreSQL (e.g. `?` → `$1` placeholders) |

They were added to **go.mod** and imported in:

- `internal/database/database.go` (main DB layer)
- `cmd/migrate/main.go`
- `cmd/init-users/main.go`
- `cmd/fix-render/main.go`

## Why that broke localhost behind your proxy

**Before that commit**

- The app only used **SQLite** for local development.
- The only extra dependency was `go-sqlite3` (and a few others). Those were likely already in your Go module cache from earlier, so no new download was needed when you ran `go build` or `go run`.

**After that commit**

- Every time you run `go build`, `go run ./cmd/migrate`, or `go run ./cmd/init-users`, Go tries to **download** `pgx` and `sqlx` from **proxy.golang.org** (if they’re not in the cache).
- Your network uses a **TLS-inspecting proxy** (e.g. corporate/school). That proxy intercepts HTTPS and presents a certificate that Go 1.23+ rejects: **"x509: negative serial number"**.
- So the **same proxy** you had before is still there; the **new** part is that the project now **must** download these two modules, and that download goes through the proxy and fails.

So: **the change that made it a problem to start the server** was adding PostgreSQL support (commit `f5841d4`), which introduced the need to download `pgx` and `sqlx` through your proxy.

## What we did to fix it (so you can start localhost again)

We did **not** remove PostgreSQL support (you still need it for Render). We did two things:

1. **`GODEBUG=x509negativeserial=1`**  
   This tells Go to accept the certificate your proxy uses (the one with the “negative serial number”). So Go can successfully download from proxy.golang.org again.  
   - You set this in your terminal before running `go run` / `go build`, or  
   - It’s set inside `start_dev_server.sh` so running the script also uses it.

2. **Baby-steps guide**  
   `BABY_STEPS_START.md` in the project root has step-by-step commands (including the `export GODEBUG=...` step) so you can start the server on localhost even with the same proxy.

## Summary

| What | Before commit f5841d4 | After commit f5841d4 |
|------|------------------------|----------------------|
| DB for localhost | SQLite only | Still SQLite (unchanged) |
| New dependencies | — | pgx, sqlx (for Render/PostgreSQL) |
| Download needed on build/run | Only if cache empty; often already cached | pgx + sqlx must be fetched if not cached |
| Behind TLS proxy | Often worked (no or few new downloads) | Fails when Go fetches from proxy.golang.org |
| Fix | — | `GODEBUG=x509negativeserial=1` (and use it in `start_dev_server.sh` / baby steps) |

So: **the change that made starting the server a problem** was the PostgreSQL migration commit that added `pgx` and `sqlx`. The fix is to keep that change but allow Go to work with your proxy by setting `GODEBUG=x509negativeserial=1` when you build and run.
