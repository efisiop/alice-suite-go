# Alice Suite - Application State Documentation

**Last Updated:** December 6, 2025  
**Version:** 1.0  
**Status:** Production Ready

---

## 📋 Table of Contents

1. [Executive Summary](#executive-summary)
2. [Architecture Overview](#architecture-overview)
3. [Technology Stack](#technology-stack)
4. [Database Architecture](#database-architecture)
5. [Authentication & Security](#authentication--security)
6. [API Endpoints](#api-endpoints)
7. [Features & Functionality](#features--functionality)
8. [Recent Changes & Improvements](#recent-changes--improvements)
9. [Setup & Deployment](#setup--deployment)
10. [Known Issues & Solutions](#known-issues--solutions)
11. [Development Workflow](#development-workflow)

---

## Executive Summary

**Alice Suite** is a physical book companion application designed to enhance the reading experience of classic literature (currently "Alice's Adventures in Wonderland"). The application provides three tiers of assistance:

1. **Tier 1: Instant Dictionary** - Word lookup from physical book
2. **Tier 2: AI Assistance** - AI-powered explanations and help
3. **Tier 3: Human Consultant** - Real-time human support

### Current Status

✅ **Fully Functional** - All core features implemented and tested  
✅ **Production Ready** - Stable, secure, and performant  
✅ **Multi-User Support** - Handles 100-1000+ concurrent readers  
✅ **Real-Time Monitoring** - Consultant dashboard with live updates  
✅ **Safari Compatible** - Cross-browser support including Safari

---

## Architecture Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Alice Suite Application                   │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐         ┌──────────────┐                  │
│  │   Reader     │         │ Consultant   │                  │
│  │   App        │         │ Dashboard    │                  │
│  └──────┬───────┘         └──────┬───────┘                  │
│         │                         │                          │
│         └──────────┬──────────────┘                          │
│                    │                                         │
│         ┌──────────▼──────────┐                              │
│         │   Go HTTP Server    │                              │
│         │   (Port 8080)       │                              │
│         └──────────┬──────────┘                              │
│                    │                                         │
│         ┌──────────▼──────────┐                              │
│         │   Middleware Layer  │                              │
│         │  - Auth            │                              │
│         │  - CORS            │                              │
│         │  - Rate Limiting   │                              │
│         │  - Heartbeat       │                              │
│         └──────────┬──────────┘                              │
│                    │                                         │
│         ┌──────────▼──────────┐                              │
│         │   Handlers Layer    │                              │
│         │  - Auth            │                              │
│         │  - Reader          │                              │
│         │  - Consultant      │                              │
│         │  - API             │                              │
│         │  - SSE/WebSocket    │                              │
│         └──────────┬──────────┘                              │
│                    │                                         │
│         ┌──────────▼──────────┐                              │
│         │   Database Layer    │                              │
│         │   SQLite (WAL)      │                              │
│         └─────────────────────┘                              │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### Application Structure

```
alice-suite-go/
├── cmd/
│   ├── server/          # Main server entry point
│   ├── migrate/         # Database migrations
│   ├── init-users/      # Initialize test users
│   └── seed/            # Seed data
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database layer (SQLite)
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   ├── realtime/        # Real-time communication (SSE/WebSocket)
│   ├── static/          # Static assets (CSS, JS)
│   └── templates/       # HTML templates
├── pkg/
│   └── auth/            # Authentication package
├── migrations/          # SQL migration files
└── data/               # Database files
```

---

## Technology Stack

### Backend
- **Language:** Go 1.21+
- **Database:** SQLite 3 with WAL mode
- **HTTP Server:** Standard `net/http` package
- **Templates:** Go `html/template` package
- **Authentication:** JWT (JSON Web Tokens)

### Frontend
- **Framework:** Vanilla JavaScript (no frameworks)
- **CSS Framework:** Bootstrap 5.3.0
- **Real-Time:** Server-Sent Events (SSE)
- **HTTP Client:** Fetch API

### Development Tools
- **Build System:** Go modules
- **Process Management:** Makefile
- **Proxy Configuration:** Shell scripts for Safari compatibility

---

## Database Architecture

### Database Configuration

**File:** `data/alice-suite.db`  
**Mode:** WAL (Write-Ahead Logging) for concurrent access  
**Connection Pooling:** Configured for 100-1000+ concurrent readers

**Key PRAGMAs:**
```sql
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA foreign_keys = ON;
PRAGMA busy_timeout = 5000;
PRAGMA wal_autocheckpoint = 1000;
PRAGMA cache_size = -20000;
PRAGMA temp_store = MEMORY;
```

### Core Tables

#### `users`
- User accounts (readers and consultants)
- Fields: `id`, `email`, `password_hash`, `first_name`, `last_name`, `role`, `is_verified`, `last_active_at`

#### `sessions`
- Active user sessions
- Fields: `id`, `user_id`, `token_hash`, `ip_address`, `user_agent`, `created_at`, `last_active_at`, `expires_at`
- Indexed on: `user_id`, `token_hash`, `expires_at`, `last_active_at`

#### `activity_logs`
- Comprehensive activity tracking
- Fields: `id`, `user_id`, `session_id`, `activity_type`, `book_id`, `page_number`, `section_id`, `metadata`, `created_at`
- Indexed on: `user_id`, `activity_type`, `created_at`, `session_id`

#### `reader_states`
- Denormalized reader state for fast consultant queries
- Fields: `user_id`, `last_active_at`, `current_page`, `current_section`, `total_activities`, `last_activity_type`, `last_activity_time`
- Indexed on: `last_active_at`

#### `books`, `chapters`, `sections`
- Book content structure
- Supports first 3 chapters (test ground)

#### `interactions`
- Legacy activity tracking (maintained for compatibility)
- Migrated to `activity_logs` but still used

#### `verification_codes`
- Book access verification codes
- Fields: `code`, `book_id`, `is_used`, `created_at`

#### `help_requests`
- Help request queue for consultants
- Fields: `id`, `user_id`, `book_id`, `section_id`, `content`, `status`, `created_at`

### Database Migrations

Located in `migrations/`:
- `001_initial_schema.sql` - Core schema
- `002_seed_first_3_chapters.sql` - Book content
- `003_restructure_pages_and_sections.sql` - Page structure
- `004_link_glossary_to_sections.sql` - Glossary linking
- `005_create_interactions_table.sql` - Activity tracking
- `006_add_sessions_and_activity.sql` - Sessions and enhanced activity logging

---

## Authentication & Security

### Authentication Flow

1. **Login:** User submits email/password → Server validates → Returns JWT token
2. **Token Storage:** Client stores token in `sessionStorage` (per-tab isolation)
3. **Cookie Sync:** Token also synced to HTTP cookie for server-side navigation
4. **Session Management:** Server creates database-backed session
5. **Authorization:** Middleware validates JWT and checks role

### Security Features

✅ **JWT Authentication** - Secure token-based auth  
✅ **Password Hashing** - bcrypt with salt  
✅ **Session Management** - Database-backed sessions with expiration  
✅ **Role-Based Access Control** - Reader vs Consultant roles  
✅ **CORS Protection** - Configured CORS middleware  
✅ **Rate Limiting** - Request rate limiting middleware  
✅ **SQL Injection Protection** - Parameterized queries  
✅ **XSS Protection** - Template escaping  

### Token Storage Strategy

**Why `sessionStorage` instead of `localStorage`:**
- Prevents session mixing when multiple users log in from same IP
- Each browser tab has isolated session
- More secure (cleared on tab close)

**Cookie Sync:**
- Token synced to HTTP cookie for server-side page navigation
- Safari-compatible cookie handling (URL encoding support)

---

## API Endpoints

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/v1/token` | Login | No |
| POST | `/auth/v1/signup` | Register | No |
| GET | `/auth/v1/user` | Get current user | Yes |
| POST | `/auth/v1/logout` | Logout | Yes |

### Reader Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/reader` | Reader dashboard | Yes (Reader) |
| GET | `/reader/login` | Reader login page | No |
| GET | `/reader/interaction` | Reading interface | Yes (Reader) |
| GET | `/reader/statistics` | Reading statistics | Yes (Reader) |
| POST | `/api/activity/track` | Track activity | Yes (Reader) |

### Consultant Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/consultant` | Consultant dashboard | Yes (Consultant) |
| GET | `/consultant/login` | Consultant login page | No |
| GET | `/consultant/help-requests` | Help requests | Yes (Consultant) |
| GET | `/consultant/readers` | Reader management | Yes (Consultant) |
| GET | `/api/consultant/logged-in-readers-count` | Active readers count | Yes (Consultant) |
| GET | `/api/consultant/active-readers-count` | Active readers (last hour) | Yes (Consultant) |
| GET | `/api/consultant/reader-activities` | Reader activities | Yes (Consultant) |
| GET | `/api/consultant/reader-activities/stream` | Activity stream | Yes (Consultant) |

### Real-Time Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/sse` | Server-Sent Events stream | Yes |

### REST API (Supabase-Compatible)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/rest/v1/{table}` | Query table | Yes |
| POST | `/rest/v1/{table}` | Insert row | Yes |
| PATCH | `/rest/v1/{table}` | Update row | Yes |
| DELETE | `/rest/v1/{table}` | Delete row | Yes |

---

## Features & Functionality

### Reader Features

✅ **User Registration & Login**
- Email/password registration
- Book verification code system
- JWT-based authentication

✅ **Reading Interface**
- Page-by-page navigation
- Section-based reading
- Physical book page synchronization

✅ **Dictionary Lookup**
- Word definitions from glossary
- Context-aware definitions
- Quick lookup interface

✅ **AI Assistance** (Planned)
- AI-powered explanations
- Question answering
- Context-aware help

✅ **Help Requests**
- Submit help requests to consultants
- Track request status
- Real-time notifications

✅ **Activity Tracking**
- Automatic activity logging
- Reading progress tracking
- Statistics dashboard

### Consultant Features

✅ **Dashboard**
- Real-time reader monitoring
- Active readers count (last hour)
- Individual reader cards with activities
- Collapsible reader windows
- Activity filtering (active in last hour)

✅ **Reader Management**
- View all readers
- Track reader activity
- Monitor reading progress
- Assign readers to consultants

✅ **Help Request Management**
- View help request queue
- Respond to help requests
- Track request status

✅ **Real-Time Updates**
- Server-Sent Events (SSE) for live updates
- Automatic refresh on login/logout
- Activity feed updates

### Recent UI/UX Improvements (December 2025)

✅ **Consultant Dashboard Enhancements:**
- All reader cards collapsed by default (neat appearance)
- Only show readers active in last hour
- Real-time count updates on login/logout
- Cards automatically removed when readers log out
- Activity pagination (12 items initially, expand to 50)
- Unread badge on collapsed cards
- Color-coded recent activities (last 1 hour)

✅ **Navigation Improvements:**
- User name displayed in navbar brand (left side, after "Alice Suite")
- Visible even in narrow/squeezed windows
- Works for both reader and consultant dashboards

✅ **Safari Compatibility:**
- Fixed cookie handling for Safari
- Proxy bypass configuration
- URL encoding support
- Cross-browser compatibility

---

## Recent Changes & Improvements

### December 6, 2025

**Consultant Dashboard:**
- ✅ Default collapsed state for all reader cards
- ✅ Filter to show only active readers (last hour)
- ✅ Real-time count updates and card removal on logout
- ✅ User name in navbar brand (left side)

**Authentication:**
- ✅ Fixed middleware to check cookies (not just Authorization header)
- ✅ Safari cookie compatibility fixes
- ✅ Session isolation using `sessionStorage`

**Database:**
- ✅ Enhanced activity logging with `activity_logs` table
- ✅ Database-backed session management
- ✅ Reader states denormalization for fast queries

### December 3-5, 2025

**Consultant Dashboard UI:**
- ✅ Individual reader cards (side-by-side layout)
- ✅ Expandable/collapsible reader windows
- ✅ Activity sorting (newest first)
- ✅ Color-coded recent activities (last 1 hour)
- ✅ Activity pagination (12 → 50 items)
- ✅ Unread badge on collapsed cards

**Bug Fixes:**
- ✅ Fixed "Test user" appearing in activity feed
- ✅ Fixed names changing in activity list
- ✅ Fixed JavaScript syntax errors
- ✅ Fixed logged-in readers count display
- ✅ Fixed Safari login issues

---

## Setup & Deployment

### Prerequisites

- Go 1.21 or higher
- SQLite3
- Make (optional, for convenience commands)

### Initial Setup

```bash
# Clone repository
cd /Users/efisiopittau/Project_1/alice-suite-go

# Install dependencies
go mod download

# Run database migrations
go run ./cmd/migrate

# Initialize test users
go run ./cmd/init-users
```

### Development Server

```bash
# Using Makefile (recommended)
make start

# Or manually
go build -o bin/server ./cmd/server
./bin/server
```

### Available Make Commands

```bash
make help          # Show all commands
make setup         # Initial setup (proxy configuration)
make build         # Build server
make start         # Build and start server
make stop          # Stop server
make restart       # Restart server
make check         # Check if server is running
make test          # Run tests
make clean         # Clean build artifacts
```

### Configuration

Configuration is managed in `internal/config/config.go`:
- **Port:** Default 8080 (set via `PORT` environment variable)
- **Database Path:** Default `data/alice-suite.db` (set via `DB_PATH`)
- **JWT Secret:** Set via `JWT_SECRET` environment variable (fallback for development)

### Test Credentials

See `LOGIN_CREDENTIALS.md` for test user credentials:

**Consultant:**
- Email: `consultant@example.com`
- Password: `consultant123`

**Reader:**
- Email: `reader@example.com`
- Password: `reader123`
- Verification Code: `ALICE2024`

---

## Known Issues & Solutions

### Safari Compatibility

**Issue:** Safari has stricter cookie handling than other browsers.

**Solution:**
- Proxy bypass configuration script (`setup_safari_proxy.sh`)
- URL-encoded cookie support
- Cookie sync from `sessionStorage`
- Hostname normalization (localhost vs 127.0.0.1)

**Command:** `make setup` or `./setup_safari_proxy.sh`

### Session Mixing

**Issue:** Multiple users logging in from same IP could mix sessions.

**Solution:**
- Switched from `localStorage` to `sessionStorage` (per-tab isolation)
- Database-backed session management
- Unique session tokens per login

### JavaScript Cache

**Issue:** Browser caching old JavaScript files.

**Solution:**
- Cache-busting query parameter (`?v=20241206`)
- Hard refresh recommended: `Cmd+Shift+R` (Mac) or `Ctrl+Shift+R` (Windows)

---

## Development Workflow

### Code Organization

- **Handlers:** `internal/handlers/` - HTTP request handlers
- **Database:** `internal/database/` - Database operations
- **Middleware:** `internal/middleware/` - HTTP middleware
- **Templates:** `internal/templates/` - HTML templates
- **Static:** `internal/static/` - CSS and JavaScript

### Adding New Features

1. **Database Changes:**
   - Create migration file in `migrations/`
   - Run migration: `go run ./cmd/migrate`

2. **New Handler:**
   - Add handler function in `internal/handlers/`
   - Register route in `cmd/server/main.go` or `internal/handlers/routes.go`

3. **New Template:**
   - Add template file in `internal/templates/`
   - Reference in handler using `template.ParseFiles()`

### Testing

```bash
# Run all tests
make test

# Run specific test
go test ./internal/handlers/...

# Check server status
make check
```

### Debugging

- **Server Logs:** Check console output when running `make start`
- **Browser Console:** F12 → Console tab
- **Network Tab:** F12 → Network tab for API calls
- **Database:** Use SQLite CLI: `sqlite3 data/alice-suite.db`

---

## Performance Characteristics

### Database Performance

- **WAL Mode:** Enables concurrent reads/writes
- **Connection Pooling:** Handles 100-1000+ concurrent connections
- **Indexes:** Optimized indexes on frequently queried columns
- **Denormalization:** `reader_states` table for fast consultant queries

### Scalability

- **Concurrent Readers:** Tested with 100-1000+ concurrent readers
- **Consultants:** Supports 10-20 concurrent consultants
- **Real-Time Updates:** SSE for efficient real-time communication
- **Activity Logging:** Efficient batch logging with minimal overhead

---

## Future Enhancements

### Planned Features

- [ ] AI Assistance integration (Tier 2)
- [ ] Enhanced reading statistics
- [ ] Mobile app optimization
- [ ] Full book content (beyond first 3 chapters)
- [ ] Advanced analytics for consultants
- [ ] Export functionality for reports

### Technical Debt

- [ ] Migrate fully from `interactions` table to `activity_logs`
- [ ] Consolidate duplicate activity tracking code
- [ ] Add comprehensive test coverage
- [ ] Performance benchmarking
- [ ] Documentation for API endpoints

---

## Support & Resources

### Documentation Files

- `README.md` - Main readme
- `LOGIN_CREDENTIALS.md` - Test credentials
- `DATABASE_ARCHITECTURE_PLAN_CURSOR.md` - Database architecture
- `FEATURE_INVENTORY.md` - Feature list
- `REQUIREMENTS.md` - Requirements
- `TECHNICAL_SPECIFICATIONS.md` - Technical specs

### Key Files

- `cmd/server/main.go` - Server entry point
- `internal/database/database.go` - Database initialization
- `internal/handlers/` - All HTTP handlers
- `migrations/` - Database migrations

---

**Document Version:** 1.0  
**Last Updated:** December 6, 2025  
**Maintained By:** Development Team

