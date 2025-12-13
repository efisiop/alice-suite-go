# Active Codebase Structure

**Last Updated:** 2025-01-23  
**Status:** Clean and organized after Refresher Protocol

---

## 📁 Directory Structure

```
alice-suite-go/
├── cmd/                          # Go application entry points
│   ├── init-users/              # User initialization tool
│   ├── migrate/                 # Database migration tool
│   ├── seed/                    # Seed data tool
│   └── server/                  # Main server (single entry point)
│
├── internal/                     # Internal application code
│   ├── database/                # Database layer
│   │   ├── database.go         # DB connection
│   │   ├── queries.go          # Query functions
│   │   └── verification.go      # Verification helpers
│   │
│   ├── handlers/                # HTTP handlers
│   │   ├── activity.go         # Activity tracking
│   │   ├── api.go              # API route setup
│   │   ├── auth.go             # Authentication
│   │   ├── consultant.go      # Consultant routes
│   │   ├── reader.go           # Reader routes
│   │   ├── rest.go             # REST API handlers
│   │   ├── routes.go           # Route setup
│   │   ├── rpc.go              # RPC handlers
│   │   ├── sse.go              # Server-Sent Events
│   │   ├── verification.go     # Verification handlers
│   │   └── websocket.go        # WebSocket handler
│   │
│   ├── middleware/              # HTTP middleware
│   │   ├── auth.go             # Auth middleware
│   │   └── middleware.go       # General middleware
│   │
│   ├── query/                   # Query parsing
│   │   ├── builder.go          # SQL query builder
│   │   └── parser.go           # Query parameter parser
│   │
│   ├── realtime/                # Real-time features
│   │   └── broadcaster.go     # Event broadcaster
│   │
│   ├── static/                  # Static assets
│   │   ├── css/
│   │   │   └── app.css
│   │   └── js/
│   │       └── app.js
│   │
│   └── templates/               # HTML templates
│       ├── base.html           # Base template
│       ├── consultant/
│       │   ├── dashboard.html
│       │   └── login.html
│       └── reader/
│           ├── dashboard.html
│           ├── interaction.html
│           ├── landing.html
│           ├── login.html
│           ├── register.html
│           ├── statistics.html
│           └── verify.html
│
├── pkg/                          # Reusable packages
│   └── auth/                    # Authentication package
│       ├── auth.go             # Auth functions
│       ├── jwt.go              # JWT handling
│       ├── session.go          # Session management
│       └── verification.go     # Book verification
│
├── migrations/                   # Database migrations
│   ├── 001_initial_schema.sql
│   ├── 002_seed_first_3_chapters.sql
│   ├── 003_restructure_pages_and_sections.sql
│   └── 004_link_glossary_to_sections.sql
│
├── data/                         # Database files
│   └── alice-suite.db           # SQLite database
│
├── archive/                      # Archived files (not part of active codebase)
│
├── alice-suite-server           # Compiled binary
│
└── Essential Documentation:
    ├── README.md                # Main readme
    ├── DEPLOYMENT.md            # Deployment guide
    ├── LOGIN_CREDENTIALS.md    # Login credentials
    ├── TESTING_CHECKLIST.md     # Testing checklist
    ├── MIGRATION_TO_GO_COMPLETE.md  # Migration guide
    ├── FEATURE_INVENTORY.md     # Feature inventory
    ├── REQUIREMENTS.md          # Requirements
    └── TECHNICAL_SPECIFICATIONS.md  # Technical specs
```

---

## ✅ Active Components

### Core Application
- **Server:** `cmd/server/main.go` - Single entry point
- **Database:** `internal/database/` - SQLite database layer
- **Handlers:** `internal/handlers/` - All HTTP handlers
- **Templates:** `internal/templates/` - Go HTML templates
- **Static Assets:** `internal/static/` - CSS and JavaScript

### Tools
- **Init Users:** `cmd/init-users/` - Create test users
- **Migrate:** `cmd/migrate/` - Database migrations
- **Seed:** `cmd/seed/` - Seed data

### Features
- **Authentication:** `pkg/auth/` - JWT-based auth
- **Real-time:** `internal/realtime/` - SSE/WebSocket
- **Query Parsing:** `internal/query/` - Supabase-compatible queries

---

## 📦 Archived Items

All archived items are in `archive/` directory:
- Completion documentation (STEP_*.md, PHASE_*.md)
- Old documentation files
- Unused code (services, models, empty directories)
- Old static files
- Reference documentation

See `archive/README.md` for details.

---

## 🚀 Quick Start

```bash
# Initialize test users
go run ./cmd/init-users

# Start server
./start.sh

# Or build and run
go build -o alice-suite-server ./cmd/server
./alice-suite-server
```

---

## 📝 Notes

- All code is Go-based (no Node.js/React in active codebase)
- Single binary deployment (`alice-suite-server`)
- Self-contained (no external runtime dependencies)
- Clean and organized structure

---

**Status:** ✅ Clean and ready for development


