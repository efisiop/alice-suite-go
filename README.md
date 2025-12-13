# Alice Suite Go

**Last Updated:** December 6, 2025  
**Status:** Production Ready ✅  
**Version:** 1.0

---

## 📚 Documentation

**📖 [Complete Application State Documentation](APPLICATION_STATE_2025_12_06.md)** - Comprehensive guide covering architecture, features, setup, and recent changes.

**Quick Links:**
- [Login Credentials](LOGIN_CREDENTIALS.md) - Test user accounts
- [Database Architecture](DATABASE_ARCHITECTURE_PLAN_CURSOR.md) - Database design and implementation
- [Feature Inventory](FEATURE_INVENTORY.md) - Complete feature list
- [Requirements](REQUIREMENTS.md) - Project requirements

---

## 🎯 Project Overview

This is the **production-ready codebase** for Alice Suite, built in **Go language** with **SQLite database**.

### Key Differences from Original Codebase

- **Language:** Go (instead of React/TypeScript)
- **Database:** SQLite (instead of Supabase/PostgreSQL)
- **Architecture:** Fresh start with streamlined design
- **Content Scope:** First 3 chapters only (test ground)

---

## 📋 Project Vision

**Alice Suite** is a **physical book companion app** that enhances the reading experience of classic literature by providing intelligent assistance **alongside** the physical book.

### Core Concept

- **Physical Book Companion** - Designed to work **alongside** the physical book, not replace it
- **Three-Tier Assistance System:**
  1. Tier 1: Instant Dictionary - Look up words from physical book
  2. Tier 2: AI Assistance - Ask questions about passages
  3. Tier 3: Human Consultant - Escalate to real-time human support

### Initial Scope

- **First 3 chapters only** - Test ground to ensure seamless operation
- **Full book** - To be added later once system works perfectly

---

## 🏗️ Planned Architecture

### Technology Stack

- **Backend:** Go language
- **Database:** SQLite
- **Frontend:** To be determined (may remain web-based)
- **Real-time:** To be determined based on Go implementation
- **AI Integration:** Moonshot AI (Kimi K2) or alternative

### Directory Structure (To Be Created)

```
alice-suite-go/
├── cmd/                    # Application entry points
│   ├── reader/            # Reader app server
│   └── consultant/        # Consultant dashboard server
├── internal/              # Private application code
│   ├── handlers/         # HTTP handlers
│   ├── services/         # Business logic
│   ├── models/           # Data models
│   └── database/         # Database layer
├── pkg/                   # Public library code
│   ├── auth/             # Authentication
│   ├── dictionary/       # Dictionary/glossary service
│   └── ai/               # AI integration
├── migrations/            # Database migrations
├── static/                # Static files (if web-based)
├── config/                # Configuration files
├── tests/                 # Test files
├── docs/                  # Documentation
├── go.mod                 # Go module file
└── README.md             # This file
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+ installed
- SQLite3
- Make (optional, for convenience commands)

### Initial Setup

```bash
# Install dependencies
go mod download

# Run database migrations
go run ./cmd/migrate

# Initialize test users
go run ./cmd/init-users

# Start development server
make start
```

### Access the Application

- **Reader App:** http://127.0.0.1:8080/reader/login
- **Consultant Dashboard:** http://127.0.0.1:8080/consultant/login
- **Health Check:** http://127.0.0.1:8080/health

### Test Credentials

See [LOGIN_CREDENTIALS.md](LOGIN_CREDENTIALS.md) for test user accounts.

### Development Commands

```bash
make help          # Show all available commands
make start         # Build and start server
make stop          # Stop server
make restart       # Restart server
make check         # Check if server is running
make test          # Run tests
make clean         # Clean build artifacts
```

---

## 📝 Development Priorities

### Phase 1: Foundation
1. Set up Go project structure
2. Initialize SQLite database
3. Create database schema
4. Implement basic authentication
5. Load first 3 chapters of Alice in Wonderland

### Phase 2: Core Features
1. Word definition lookup system
2. AI assistance integration
3. Help request system
4. Progress tracking
5. Consultant dashboard

### Phase 3: Refinement
1. Bug fixes and testing
2. Streamlined authentication flow
3. Performance optimization
4. Documentation

### Phase 4: Expansion
1. Add remaining chapters (full book)
2. Enhanced features
3. Mobile optimization

---

## 📚 Reference Documents

- **Recovered Brief:** `/Users/efisiopittau/Project_1/alice-suite/ALICE_SUITE_RECOVERED_BRIEF.md`
- **Features Report:** `/Users/efisiopittau/Project_1/alice-suite/ALICE_SUITE_FEATURES_REPORT.md`
- **Requirements:** `/Users/efisiopittau/Project_1/alice-suite/ALICE_SUITE_REQUIREMENTS.md`

---

## 🎯 Key Principles

1. **Physical Book Companion** - Always emphasize companion nature, not replacement
2. **Simplicity** - Streamlined, easy-to-use interface
3. **Performance** - Fast, responsive Go backend
4. **Reliability** - Bug-free, stable operation
5. **Test Ground** - Start with 3 chapters, expand when seamless

---

**This is a fresh start. Build it right from the beginning!**

