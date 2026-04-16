# Alice Suite — Codebase summary (for QMD / project memory)

This summary is derived from the GitHub dump (`efisiop-alice-suite-go-8a5edab282632443-2.txt`) so that QMD and agents have a single place for high-level codebase facts. Re-index with `qmd update` and `qmd embed` after edits.

---

## What Alice Suite is

- **Physical book companion** for classic literature (Alice in Wonderland first). It works **alongside** the physical book; it does not replace it.
- **Three-tier assistance:** (1) Instant dictionary, (2) AI assistance, (3) Human consultant.
- **Three apps:** Reader (end users), Consultant (support staff), Administrator (board/ops monitoring).
- **Content scope:** First 3 chapters as test ground; full book to be added later.

---

## Technology stack

- **Backend:** Go.
- **Database:** SQLite locally; PostgreSQL when `DATABASE_URL` is set (e.g. Render).
- **Frontend:** Server-rendered HTML templates in `internal/templates/` (base, reader/, consultant/, admin/), plus shared CSS/JS in `internal/static/`.
- **Auth:** JWT + cookie; roles: `reader`, `consultant`, `administrator`. Sessions and heartbeat for “online” presence.
- **Real-time:** Server-Sent Events (SSE) for consultant dashboard; optional WebSocket.
- **AI:** Configurable (e.g. Gemini, Moonshot); see AI_SETUP_* and LOGIN_CREDENTIALS.

---

## Repo layout (high level)

```
alice-suite-go/
├── cmd/                    # Entry points
│   ├── server/            # Main server (Reader + Consultant + Admin)
│   ├── reader/            # Reader-only API server (alternative)
│   ├── migrate/           # Run DB migrations
│   ├── init-users/        # Seed reader, consultant, admin test users
│   ├── verify-deployment/
│   ├── fix-render/        # Fix sections on Render
│   └── ...
├── internal/
│   ├── handlers/          # HTTP handlers (auth, api, reader, consultant, admin, rest, rpc, sse)
│   ├── middleware/        # Auth, rate limit, heartbeat, hostname
│   ├── database/          # DB layer (sessions, activity, consultant, settings, etc.)
│   ├── services/          # Book, AI, dictionary, help, image
│   ├── config/            # Load config (DB path, port, DATABASE_URL)
│   ├── templates/         # base.html + reader/, consultant/, admin/
│   ├── static/            # css/app.css, js/app.js
│   ├── email/             # SMTP + simulated login emails
│   ├── realtime/          # Broadcaster for SSE
│   └── models/, errors/, query/, useragent/
├── pkg/auth/              # JWT, verification (public package)
├── migrations/             # SQL migrations (001–015)
├── docs/                   # Docs and project memory (QMD, APP_UPGRADES, etc.)
├── scripts/                # Sync with GitHub dump, generate-dump, etc.
├── archive/                # Deprecated and reference docs
├── go.mod, go.sum
├── Makefile                # start, stop, build, test, clean, check
└── render.yaml, Dockerfile # Deploy
```

---

## Three apps — URLs and purpose


| App               | Login URL           | Main URL / dashboard                                              | Purpose                                                                |
| ----------------- | ------------------- | ----------------------------------------------------------------- | ---------------------------------------------------------------------- |
| **Reader**        | `/reader/login`     | `/reader`, `/reader/interaction`                                  | Read book, dictionary, AI help, help requests, progress.               |
| **Consultant**    | `/consultant/login` | `/consultant`, `/consultant/readers`, `/consultant/help-requests` | See readers, activity, help requests, send prompts.                    |
| **Administrator** | `/admin/login`      | `/admin`                                                          | Monitor readers/consultants online, email notifications, future stats. |


- **Health:** `/health` (no auth).
- **Reader API** (when using `cmd/reader`): `/api/health`, `/api/auth/login`, `/api/books`, `/api/chapters`, `/api/sections`, `/api/pages`, `/api/dictionary/lookup`, `/api/ai/ask`, `/api/help/request`, etc.

---

## Key commands

```bash
go run ./cmd/migrate      # Run migrations
go run ./cmd/init-users   # Create reader, consultant, admin test users
make start                # Build and start server (default port 8080)
make stop / make restart / make test / make check
```

- **Compare DB (localhost vs Render):** `./bin/compare-db-structure`
- **Fix sections on Render:** `./bin/fix-render`
- **Simulate login emails (no SMTP):** `SIMULATE_LOGIN_EMAILS=true make start`

---

## Important docs (in repo)

- **AGENTS.MD** — Project rules and scope for agents; read first.
- **docs/AI_AGENTS_CHECKLIST.md** — One-page list of Cursor, MCP, Paperclip, `qmd`, and CLI helpers.
- **LOGIN_CREDENTIALS.md** — Test accounts (reader, consultant, admin).
- **docs/QMD_PROJECT_MEMORY.md** — How to use qmd collection `alice`.
- **docs/APP_UPGRADES.md** — Log of app upgrades (noscript, canonical URLs, nav, logout).
- **ADMIN_SETUP.md** — How to create and use the administrator account.
- **DEPLOYMENT.md**, **RENDER_*.md** — Deploy and env vars.
- **APPLICATION_STATE_*.md**, **FEATURE_INVENTORY.md** — State and feature list.

---

## One-line reminder (for qmd search)

**Alice Suite** = Go app, physical book companion (dictionary + AI + human consultant). Three apps: Reader, Consultant, Admin. Stack: Go, SQLite/Postgres, HTML templates in `internal/templates`, `cmd/server` + `cmd/reader`, migrations in `migrations/`. Use `qmd query "…" -c alice` for deeper lookup.