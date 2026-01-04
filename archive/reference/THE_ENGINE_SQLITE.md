# The Engine: Technical Architecture (SQLite Version)

**Updated:** 2025-01-XX  
**Purpose:** Updated technical architecture documentation replacing Supabase/PostgreSQL with SQLite  
**Based on:** PROJECT_BRIEFING.pdf (Section 6: Technical Architecture)

---

## Overview

This document describes the technical architecture of **Alice AI Companion** (Alice Suite), updated to reflect the current SQLite-based implementation. The original briefing mentioned Supabase with PostgreSQL, but the actual implementation uses **SQLite 3 with WAL mode** for a self-contained, efficient database solution.

---

## Technical Architecture

The application is built on a **modern Go-based backend** with **SQLite database** to ensure security, performance, and robust functionality.

### Backend Platform

- **Backend:** Go language server (`cmd/server/main.go`)
- **Database:** **SQLite 3 with WAL mode** (`internal/database/database.go`)
- **HTTP Server:** Standard library `net/http`
- **Authentication:** JWT-based authentication (`pkg/auth/`)
- **Real-time:** Server-Sent Events (SSE) (`internal/handlers/sse.go`)
- **AI Integration:** Go service layer with external LLM API (`internal/services/ai_service.go`)

### Frontend Technology

- **Reader App:** Go HTML templates with Bootstrap 5 (`internal/templates/reader/`)
- **Publisher Dashboard:** Go HTML templates with Bootstrap 5 (`internal/templates/consultant/`)
- **UI Framework:** Bootstrap 5.3.0 for responsive design
- **JavaScript:** Vanilla JavaScript (no frameworks) for interactivity
- **Real-time Updates:** Server-Sent Events (SSE) for live activity feed

### Database Structure

The core of the backend is a **SQLite 3 database** with WAL (Write-Ahead Logging) mode for optimal concurrent access. The database includes several key tables designed to manage the application's logic:

---

## Database Tables

| Table Name | Primary Function |
|------------|------------------|
| `users` | Stores user information (name, email, role), verification status, and links their identity to authenticated sessions. Replaces the `profiles` table mentioned in briefing. |
| `books` | Contains the segmented text of "Alice in Wonderland," organized by page and section for easy retrieval. |
| `definitions` / `alice_glossary` | A curated database of pre-defined content (explanations, examples) that powers the Tier 1 instant pop-ups. |
| `verification_codes` | Manages the unique verification codes included with the physical books, tracking their usage and linking them to the user who verified them. |
| `interactions` | Detailed log of all user activity, including Tier 1 lookups, AI queries, help requests, and records of interactions. |
| `activity_logs` | Comprehensive activity tracking for consultant dashboards (additional to `interactions`). |
| `consultant_triggers` | Stores requests initiated by consultants from the Publisher Dashboard to send specific, subtle AI prompts to individual readers, tracking the status of each trigger. |
| `sessions` | Database-backed user sessions for authentication persistence. |
| `reading_progress` | Tracks reader progress (current page, section, book purchase date). |
| `help_requests` | Tier 3 help requests from readers to consultants. |
| `ai_interactions` | Records of AI assistant interactions (Tier 2). |
| `consultant_assignments` | Links consultants to assigned readers. |

---

## Database Configuration

### SQLite with WAL Mode

The database uses **SQLite 3 with Write-Ahead Logging (WAL mode)** for optimal concurrent access. This allows multiple readers and writers to access the database simultaneously without locking errors.

**Configuration PRAGMAs:**
- `journal_mode = WAL` - Enables Write-Ahead Logging (allows concurrent reads/writes)
- `synchronous = NORMAL` - Good balance between speed and safety
- `foreign_keys = ON` - Enforces referential integrity
- `busy_timeout = 5000` - Waits 5 seconds before "database locked" error
- `wal_autocheckpoint = 1000` - Auto-checkpoint WAL every 1000 pages
- `cache_size = -20000` - 20 MB page cache
- `temp_store = MEMORY` - Store temp tables in memory for speed

**Connection Pooling:**
- Max Open Connections: 25
- Max Idle Connections: 5
- Connection Lifetime: Unlimited (reuse connections)

**Database Location:**
- File: `data/alice-suite.db`
- WAL file: `data/alice-suite.db-wal`
- Shared memory: `data/alice-suite.db-shm`

---

## Key Differences from Briefing

### Original Briefing (Outdated)
- **Backend Platform:** Supabase (BaaS)
- **Database:** PostgreSQL
- **Frontend:** React (TypeScript) with Material-UI
- **Edge Functions:** Serverless Edge Functions for AI

### Current Implementation (Correct)
- **Backend Platform:** Go HTTP server
- **Database:** **SQLite 3 with WAL mode**
- **Frontend:** Go HTML templates with Bootstrap 5
- **Real-time:** Server-Sent Events (SSE)
- **AI Integration:** Go service layer with external LLM API

---

## Architecture Layers

```
┌─────────────────────────────────────────────┐
│         Web Browser (Frontend)              │
│   - HTML Templates (Go templates)           │
│   - Bootstrap 5 CSS Framework               │
│   - Vanilla JavaScript                      │
│   - SSE Client for Real-time Updates        │
└─────────────────┬───────────────────────────┘
                  │ HTTP/REST API + SSE
┌─────────────────▼───────────────────────────┐
│         Go HTTP Server                      │
│   - Port 8080 (configurable)                │
│   - Standard library net/http               │
│   - Middleware (auth, heartbeat, rate limit)│
└─────────────────┬───────────────────────────┘
                  │
┌─────────────────▼───────────────────────────┐
│         Handlers Layer                      │
│   - internal/handlers/                      │
│   - REST API endpoints                      │
│   - Template rendering                      │
│   - SSE event broadcasting                  │
└─────────────────┬───────────────────────────┘
                  │
┌─────────────────▼───────────────────────────┐
│         Services Layer                      │
│   - internal/services/                      │
│   - Business logic                          │
│   - AI service integration                  │
│   - Help request management                 │
└─────────────────┬───────────────────────────┘
                  │
┌─────────────────▼───────────────────────────┐
│         Database Layer                      │
│   - internal/database/                      │
│   - Query functions                         │
│   - Session management                      │
│   - Activity logging                        │
└─────────────────┬───────────────────────────┘
                  │
┌─────────────────▼───────────────────────────┐
│         SQLite 3 Database                   │
│   - data/alice-suite.db                     │
│   - WAL mode enabled                        │
│   - Connection pooling                      │
└─────────────────────────────────────────────┘
```

---

## Security & Access Control

### Authentication
- **JWT-based authentication** (`pkg/auth/auth.go`)
- **Session management** (database-backed sessions table)
- **Password hashing** (bcrypt)
- **Role-based access control** (reader/consultant roles)

### Data Isolation
- **Readers:** Can only access their own data (filtered by `user_id`)
- **Consultants:** Can view all reader data (filtered by `role = 'reader'`)
- **Middleware:** `RequireAuth` and `RequireConsultant` protect endpoints

### Privacy Controls
- **Purpose Limitation:** Data access strictly for support purposes
- **No Unsolicited Contact:** Consultants cannot initiate direct contact
- **Transparency:** Clear messaging in onboarding about data usage
- **Reader Control:** Reader always controls escalation to Tier 2/3

---

## Real-time Features

### Server-Sent Events (SSE)
- **Endpoint:** `/api/sse`
- **Purpose:** Real-time activity feed for consultant dashboard
- **Events:** Activity updates, login notifications
- **Implementation:** `internal/handlers/sse.go`, `internal/realtime/broadcaster.go`

### Heartbeat Mechanism
- **Middleware:** `internal/middleware/heartbeat.go`
- **Purpose:** Track "who's online" by updating `users.last_active_at`
- **Frequency:** Updates on every authenticated request

---

## AI Integration

### AI Service Layer
- **Location:** `internal/services/ai_service.go`
- **Purpose:** Handles Tier 2 AI Assistant interactions
- **External API:** Moonshot AI (Kimi) or Anthropic API
- **Context Awareness:** Includes book context and reader position
- **Storage:** AI interactions stored in `ai_interactions` table

---

## Deployment

### Local Development
- **Database:** SQLite file (`data/alice-suite.db`)
- **Server:** Go HTTP server (`cmd/server/main.go`)
- **Port:** 8080 (configurable via environment variable)

### Production
- **Database:** SQLite with WAL mode (same as development)
- **Server:** Go binary compiled for target platform
- **Deployment:** Render.com or similar platform
- **Configuration:** Environment variables for sensitive data

---

## Summary

The current implementation uses **SQLite 3 with WAL mode** as the database engine, providing:
- ✅ **Self-contained** database (no external dependencies)
- ✅ **Concurrent access** via WAL mode
- ✅ **High performance** with optimized PRAGMAs
- ✅ **Simplicity** - single file database
- ✅ **Reliability** - ACID-compliant transactions
- ✅ **Portability** - easy backup and migration

This architecture is production-ready and supports 100-1000+ concurrent readers with 10-20 consultants monitoring in real-time.

