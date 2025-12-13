# Alice Suite - Quick Reference Guide

**Last Updated:** December 6, 2025

---

## 🚀 Quick Start

```bash
# Setup (first time only)
make setup          # Configure Safari proxy bypass
go run ./cmd/migrate
go run ./cmd/init-users

# Start server
make start

# Check status
make check
```

**Access:**
- Reader: http://127.0.0.1:8080/reader/login
- Consultant: http://127.0.0.1:8080/consultant/login

---

## 🔑 Test Credentials

**Consultant:**
- Email: `consultant@example.com`
- Password: `consultant123`

**Reader:**
- Email: `reader@example.com`
- Password: `reader123`
- Verification Code: `ALICE2024`

---

## 📁 Key Directories

```
cmd/server/         # Main server entry point
internal/handlers/  # HTTP handlers
internal/database/  # Database layer
internal/templates/ # HTML templates
internal/static/    # CSS and JavaScript
migrations/         # Database migrations
```

---

## 🔧 Common Tasks

### Add New Migration

```bash
# Create migration file
touch migrations/007_your_migration.sql

# Run migration
go run ./cmd/migrate
```

### Add New Handler

1. Create handler in `internal/handlers/`
2. Register route in `cmd/server/main.go`
3. Add template if needed in `internal/templates/`

### Debug Issues

```bash
# Check server logs
make start  # Watch console output

# Check database
sqlite3 data/alice-suite.db

# Browser console
F12 → Console tab
```

---

## 🐛 Troubleshooting

### Safari Not Connecting

```bash
make setup  # Configure proxy bypass
```

### JavaScript Not Updating

- Hard refresh: `Cmd+Shift+R` (Mac) or `Ctrl+Shift+R` (Windows)
- Clear browser cache

### Server Won't Start

```bash
# Check if port is in use
lsof -i :8080

# Kill process if needed
make stop
```

---

## 📊 Database

**Location:** `data/alice-suite.db`  
**Mode:** WAL (Write-Ahead Logging)  
**Connection Pool:** 25 max connections

**Key Tables:**
- `users` - User accounts
- `sessions` - Active sessions
- `activity_logs` - Activity tracking
- `reader_states` - Reader state (denormalized)
- `books`, `chapters`, `sections` - Book content

---

## 🔐 Authentication

- **Method:** JWT (JSON Web Tokens)
- **Storage:** `sessionStorage` (per-tab isolation)
- **Cookie Sync:** HTTP cookie for server-side navigation
- **Sessions:** Database-backed with expiration

---

## 📡 API Endpoints

### Auth
- `POST /auth/v1/token` - Login
- `POST /auth/v1/signup` - Register
- `GET /auth/v1/user` - Get current user
- `POST /auth/v1/logout` - Logout

### Reader
- `GET /reader` - Dashboard
- `GET /reader/interaction` - Reading interface
- `POST /api/activity/track` - Track activity

### Consultant
- `GET /consultant` - Dashboard
- `GET /api/consultant/logged-in-readers-count` - Active readers
- `GET /api/consultant/reader-activities` - Reader activities
- `GET /sse` - Real-time updates (SSE)

---

## 🎨 Recent Features (Dec 2025)

✅ Consultant dashboard with individual reader cards  
✅ Real-time activity monitoring  
✅ Active readers filtering (last hour)  
✅ User name in navbar (left side)  
✅ Safari compatibility fixes  
✅ Session isolation (per-tab)  

---

## 📚 Full Documentation

See [APPLICATION_STATE_2025_12_06.md](APPLICATION_STATE_2025_12_06.md) for complete documentation.

---

**Quick Reference Version:** 1.0  
**Last Updated:** December 6, 2025

