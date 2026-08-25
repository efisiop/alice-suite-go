# Hermes Handover Briefing — Alice Suite

**Purpose:** Hand over everything done so far on the Alice Suite codebase to **Hermes** (the agent runtime on the VPS), and list the **skills** needed to take the three apps — **Reader, Consultant, Admin** — to a much higher standard and better optimization.

**Date:** 2026-06-03
**Repo:** `alice-suite-go` (Go)
**Read first:** `AGENTS.MD`, then this file.

---

## 1. What Alice Suite is (in one paragraph)

Alice Suite is a **companion app for a physical book** (starting with *Alice's Adventures in Wonderland*). It does **not** replace the book — it helps the reader **alongside** it, in three tiers: **(1)** instant dictionary, **(2)** AI assistance (explain / simplify / chat / quiz), **(3)** a human consultant for live help. There are **three apps** in one Go server: a **Reader** app for end users, a **Consultant** app for support staff, and an **Administrator** app for oversight. Current content scope is the **first 3 chapters** as a test ground; the full book comes later.

---

## 2. Current state snapshot

### 2.1 Technology stack
- **Backend:** Go 1.24, standard `net/http` with `http.ServeMux`.
- **Database:** SQLite locally (`mattn/go-sqlite3`); PostgreSQL in production via `pgx` when `DATABASE_URL` is set.
- **Auth:** JWT (`golang-jwt/jwt/v5`) + cookie + DB-backed sessions. Roles: `reader`, `consultant`, `administrator`.
- **Real-time:** Server-Sent Events (SSE) for the consultant dashboard; optional WebSocket (`gorilla/websocket`). A `realtime.Broadcaster` fans out events.
- **Frontend:** Server-rendered Go `html/template` files in `internal/templates/` + shared `internal/static/css/app.css` and `internal/static/js/app.js`. (No SPA / framework.)
- **AI:** Pluggable provider service (`internal/services/ai_service.go`) — Gemini / Moonshot / others, selected by env vars.
- **Deploy today:** Render (`render.yaml`, `Dockerfile`) plus shell scripts. **Target now: the VPS where Hermes lives.**

### 2.2 The three apps and what already works

| App | Login | Main pages | State today |
|-----|-------|-----------|-------------|
| **Reader** | `/reader/login` | `/reader` dashboard, `/reader/interaction` (reading + dictionary + AI), `my-page`, `statistics`, `verify` | Most mature. Dictionary popup, AI help, sections, progress tracking, activity logging all work. |
| **Consultant** | `/consultant/login` | `dashboard`, `readers`, `reader-inspector`, `help-requests`, `send-prompt`, `assign-readers`, `feedback`, `reading-insights`, `reports` | Functional. Live presence, reader activity charts/summaries, help-request handling, prompts. |
| **Admin** | `/admin/login` | `/admin` dashboard | **Thinnest.** Only live presence (readers/consultants online) + login-email notification settings. The bigger vision (consultant effectiveness matrix, reader follow-up, board view) is proposed but **not built** — see `ADMINISTRATOR_PROPOSAL.md`. |

### 2.3 Where the code lives (map)
```
cmd/            entry points: server (main), reader, migrate, init-users, seed, verify-deployment, fix-* 
internal/
  handlers/     HTTP handlers: auth, reader, consultant, admin, api, rest, rpc, sse, websocket, routes
  middleware/   auth, rate_limit, heartbeat, hostname
  database/     sessions, activity, consultant, settings, verification, queries, pg_transform
  services/     ai_service, book_service, dictionary_service, help_service, image_service
  templates/    base.html + reader/ + consultant/ + admin/
  static/       css/app.css, js/app.js, images/
  realtime/     broadcaster (SSE/WebSocket fan-out)
pkg/auth/       JWT, verification (public package)
migrations/     001–016 SQL migrations
docs/           live docs + evolution wiki (docs/wiki/)
archive/        old code + old docs (do NOT treat as current)
```

### 2.4 House rules already in place (Hermes must follow these)
- **Evolution control:** For each non-trivial change, add/update a row in `docs/wiki/index.md` and append a dated entry in `docs/wiki/log.md`. Schema is in `docs/EVOLUTION_CONTROL.md`.
- **Roadmap:** Approved upgrade sequence is in `docs/UPGRADE_ROADMAP.md` (Reader 1.x, AI 2.x, Consultant 3.x). Work item-by-item.
- **Scope:** Agentic tooling rules apply to this repo only (`AGENTS.MD`).
- **Test creds:** `LOGIN_CREDENTIALS.md`.

---

## 3. Known gaps and technical debt (be honest with Hermes)

These are the concrete things standing between "works" and "high standard". Hermes should treat them as the backlog.

1. **Templates are re-parsed on every request.** ~23 `template.ParseFiles(...)` calls live inside handlers (e.g. `admin.go`, `reader.go`, `consultant.go`). This is slow and fragile, and it breaks if the working directory changes on the VPS. → Parse once at startup and/or use `embed.FS`.
2. **Static files & templates read from disk by relative path.** On a VPS this depends on the launch directory. → Embed templates/static with `//go:embed` so the binary is self-contained.
3. **Thin test coverage.** Only 6 test files (`pkg/auth`, `middleware`, `handlers/api`, `services/book`). No end-to-end coverage of reader/consultant/admin flows. → Add handler + integration tests.
4. **No CI pipeline.** No automated build/test/lint on push. → Add CI (build, `go vet`, `go test`, `golangci-lint`).
5. **Admin app underbuilt** vs. `ADMINISTRATOR_PROPOSAL.md` (presence only; no stats matrix or follow-up metrics).
6. **AI readability/UX not finished** (`UPGRADE_ROADMAP.md` 2.1–2.3: response formatting, quiz feedback, single AI panel).
7. **Documentation sprawl.** 50+ markdown files at repo root, many overlapping "fix" notes. → Consolidate into `docs/`, archive the rest.
8. **North star not filled in.** `docs/wiki/DIRECTION.md` is still a template. → Pin down product direction so optimization has a target.
9. **Secrets / config on the VPS.** A `.env` exists at repo root; confirm it is git-ignored and move to VPS environment/secret management. See `URGENT_API_KEY_SECURITY.md`.
10. **No structured logging / metrics / health depth.** `/health` is basic. → Add structured logs, request metrics, and readiness checks for VPS monitoring.

---

## 4. Skills Hermes needs — to reach a higher standard & optimization

Organized as **cross-cutting skills** (apply to all three apps) and **per-app skills**. Each skill says *why it matters here*.

### 4.1 Cross-cutting (foundation — do these first)

**A. Go backend engineering**
- Idiomatic Go, `net/http`, `context`, error handling, goroutines/channels (the SSE broadcaster uses them).
- **Template management:** compile `html/template` once at boot; use `embed.FS` for templates + static assets.
- Dependency hygiene: `go mod tidy`, pin versions, `go vet`.

**B. Database & migrations**
- SQL for **both SQLite and PostgreSQL** (the code targets both). Understand `internal/database/pg_transform.go`.
- Writing safe forward migrations in `migrations/` (next is 017).
- Indexing and query optimization for `sessions`, `activity_logs`, `reader_states`, `help_requests`, `consultant_prompts` (these drive presence and the future admin stats).
- Connection pooling and avoiding N+1 queries.

**C. Testing & quality**
- Go unit tests (`testing`), table-driven tests, `httptest` for handlers.
- Integration tests against a temp SQLite DB (`internal/database/testutils.go` exists).
- `golangci-lint` / `go vet` discipline.

**D. DevOps on the VPS (this is the handover target)**
- Linux service management: run the Go binary under **systemd** (auto-restart, logs via journald), or Docker (a `Dockerfile` exists).
- **Reverse proxy + TLS:** Nginx/Caddy in front, HTTPS certificates (Let's Encrypt/Caddy auto-TLS).
- **Environment/secrets management:** move API keys and DB URL out of files into env/secret store.
- **Backups:** scheduled DB backups (SQLite file or Postgres dump).
- **CI/CD:** build → test → deploy pipeline (GitHub Actions or a VPS-side deploy script + webhook).
- **Observability:** structured logging, basic metrics, uptime/health monitoring, log rotation.

**E. Security**
- Auth hardening: JWT/cookie flags (HttpOnly, Secure, SameSite), session expiry, CSRF on form posts.
- Rate limiting (already started in `middleware/rate_limit.go`) and input validation.
- Secrets never committed; review `.gitignore`.
- Dependency vulnerability scanning (`govulncheck`).

**F. Performance & optimization**
- Caching (templates, dictionary lookups — `dictionary_service` + `dictionary_cache` table already exist).
- HTTP: gzip/compression, cache headers, fingerprinted static assets.
- Profiling with `pprof` to find hot paths before optimizing.

**G. Frontend (server-rendered) & UX**
- HTML/CSS, responsive/mobile-first (readers use phones next to the book).
- Accessibility (semantic HTML, contrast, keyboard).
- Vanilla JS in `app.js` (no framework); progressive enhancement (there are `<noscript>` blocks).
- Light markdown→HTML rendering for AI answers (roadmap 2.1).

**H. AI integration**
- Prompt design and provider abstraction (`ai_service.go`), streaming responses, timeouts/fallbacks, cost control, caching AI results.

**I. Agent-loop design and governance**
- Design closed, bounded loops with explicit trigger, action, feedback signal, stop condition, and human escalation path.
- Apply the Alice loop inventory before adding agentic behavior: prioritize reader-context help, AI grounding/recovery, and the consultant handoff; defer autonomous optimization until consent, evaluation, cost, and privacy controls exist.
- Treat this as a reusable planning skill for Reader, Consultant, and Admin changes. Reference `docs/AGENT_LOOP_INVENTORY.md` for the current catalogue and guardrails.

### 4.2 Per-app skills and target improvements

**READER (Tier 1 + 2) — most mature, polish & speed**
- Skills: frontend/UX, CSS responsive design, vanilla JS, dictionary caching, AI readability.
- Targets (from `UPGRADE_ROADMAP.md` 1.x & 2.x): clearer dictionary popup (source/part-of-speech/examples), better reading layout on phones, dashboard entry points, readable AI responses, working quiz with feedback, one unified "AI help" panel.

**CONSULTANT (Tier 3) — functional, organize & scale**
- Skills: data aggregation/SQL, SSE/real-time, dashboard data-viz, table UX (search/filter/sort).
- Targets (roadmap 3.x): organized Readers list with search/filter, clearer help-request statuses and actions, consistent reader-activity cards, reliable live presence.

**ADMIN — thinnest, biggest build-out**
- Skills: SQL aggregation & caching, role-based access (`RequireAdmin` exists), dashboard/reporting UI, data-viz.
- Targets (from `ADMINISTRATOR_PROPOSAL.md`): consultant effectiveness matrix (response time, resolution rate, prompts sent/accepted, readers followed up), reader follow-up overview, board-level summary tiles, optional `admin_level` (global vs regional). All **read-only** aggregations over existing tables, with caching for speed.

---

## 5. Recommended sequence for Hermes

1. **Stabilize for the VPS:** embed templates/static + parse-once; run under systemd/Docker behind Nginx/Caddy with TLS; move secrets to env. (Removes gaps 1, 2, 9.)
2. **Safety net:** add CI (build/vet/test/lint) and handler/integration tests. (Gaps 3, 4.)
3. **Pin direction:** fill `docs/wiki/DIRECTION.md` with the user so optimization has a clear target. (Gap 8.)
4. **Reader polish:** roadmap 1.x then 2.x.
5. **Consultant organize:** roadmap 3.x.
6. **Admin build-out:** implement `ADMINISTRATOR_PROPOSAL.md` step by step.
7. **Optimize & observe:** profiling, caching, structured logging/metrics. (Gaps 6, 10.)
8. **Docs cleanup:** consolidate root markdown into `docs/`, archive the rest. (Gap 7.)

---

## 6. How to verify success (definition of done)

- `go build ./...` and `go test ./...` pass; lint clean.
- App runs from the binary alone on the VPS (no dependence on launch directory), served over HTTPS, auto-restarts on crash.
- Each app's roadmap items pass a manual check on desktop **and** phone using `LOGIN_CREDENTIALS.md`.
- Every non-trivial change is recorded in `docs/wiki/index.md` + `docs/wiki/log.md`.
- DB backups run on a schedule; secrets are not in the repo.

---

## 7. Quick start for Hermes on the VPS

```bash
go mod download
go run ./cmd/migrate        # apply migrations (001–016)
go run ./cmd/init-users     # seed reader, consultant, admin test users
make start                  # build + run (default port 8080)
# Reader:     http://<host>:8080/reader/login
# Consultant: http://<host>:8080/consultant/login
# Admin:      http://<host>:8080/admin/login
# Health:     http://<host>:8080/health
```

For production: put it behind Nginx/Caddy with TLS, set `DATABASE_URL` for Postgres, set AI provider keys, and run as a systemd service or Docker container.

---

*Source docs referenced: `AGENTS.MD`, `docs/CODEBASE_SUMMARY.md`, `docs/UPGRADE_ROADMAP.md`, `docs/APP_UPGRADES.md`, `ADMINISTRATOR_PROPOSAL.md`, `docs/EVOLUTION_CONTROL.md`, `docs/wiki/index.md`, `docs/wiki/log.md`, `LOGIN_CREDENTIALS.md`.*
