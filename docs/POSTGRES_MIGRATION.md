# PostgreSQL Migration Guide

Alice Suite now supports **PostgreSQL** in addition to SQLite. On Render.com, using PostgreSQL prevents data loss caused by the ephemeral filesystem.

## Why Migrate?

- **Render.com**: SQLite data is lost on every deploy/restart
- **PostgreSQL**: Data persists across deploys
- **Reader progress, help requests, new accounts** are never lost again

## Render.com Setup

### Option A: Blueprint (render.yaml) – Recommended

The `render.yaml` in the repo defines a PostgreSQL database and links it automatically:

1. Push your code to GitHub and connect the repo to Render.
2. Render will create:
   - **alice-suite-go-db** (PostgreSQL, free plan)
   - **alice-suite-go** (web service)
3. `DATABASE_URL` is set from the database via `fromDatabase`, so no manual env setup is needed.

On each deploy, `start.sh` runs migrations → init-users → fix-render → server.

### Option B: Manual Setup

1. Go to [Render Dashboard](https://dashboard.render.com)
2. Create **New +** → **PostgreSQL** (e.g. `alice-suite-go-db`)
3. Open your **alice-suite-go** web service → **Environment**
4. Add `DATABASE_URL` = Internal Database URL from the PostgreSQL service
5. Redeploy

### Production Requirement

In production (`ENV=production`), `DATABASE_URL` must be set. If it is missing, startup fails with a clear error. This prevents accidentally running with ephemeral SQLite on Render.

## Local Development

### Option A: SQLite (default)

No changes. Run as usual:

```bash
go run ./cmd/server
```

Uses `data/alice-suite.db` by default.

### Option B: PostgreSQL (matches production)

1. Start PostgreSQL (Docker):

   ```bash
   docker run -d -p 5432:5432 \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=alice \
     --name alice-pg postgres:16
   ```

2. Set `DATABASE_URL`:

   ```bash
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/alice?sslmode=disable"
   go run ./cmd/server
   ```

## Environment Variables

| Variable      | Required | Description                                      |
|--------------|----------|--------------------------------------------------|
| `DATABASE_URL` | No       | PostgreSQL connection URL. When set, uses PostgreSQL. |
| `DB_PATH`    | No       | SQLite file path (default: `data/alice-suite.db`). Used when `DATABASE_URL` is not set. |

## Troubleshooting

**Migrations fail on PostgreSQL**

- Ensure the database user has CREATE rights
- Check `DATABASE_URL` format: `postgres://user:pass@host:5432/dbname?sslmode=disable`

**init-users fails**

- Migrations must run first (start.sh does this automatically)
- If manually testing: run `./bin/migrate` before `./bin/init-users`
