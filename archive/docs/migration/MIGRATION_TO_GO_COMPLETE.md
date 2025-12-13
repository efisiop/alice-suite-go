# Complete Migration Guide: React/TypeScript to Go + SQLite

**Created:** 2025-11-23  
**Purpose:** Complete migration of Alice Reader and Consultant Dashboard from React/TypeScript to Go  
**Target:** Single, self-contained Go application with embedded SQLite database

---

## 📋 Executive Summary

### Current Architecture (To Be Migrated)

**Location:** `/Users/efisiopittau/Project_1/alice-suite/`

**Technology Stack:**
- **Frontend:** React 18 + TypeScript + Material-UI (MUI)
- **Build Tool:** Vite 4.0
- **Backend:** Node.js Express (SQLite adapter on port 54321)
- **Database:** SQLite (`alice_suite.db`)
- **Real-time:** Socket.io/WebSocket
- **Apps:** Two separate React applications
  - Alice Reader (port 5173)
  - Consultant Dashboard (port 5174)

**Key Dependencies:**
- `@supabase/supabase-js` (using local SQLite adapter)
- `@mui/material`, `@mui/icons-material`
- `react-router-dom`
- `socket.io-client`
- `notistack` (notifications)

### Target Architecture (Go Implementation)

**Location:** `/Users/efisiopittau/Project_1/alice-suite-go/`

**Technology Stack:**
- **Backend:** Go 1.21+ (standard library `net/http`)
- **Database:** SQLite 3 (direct, no adapter)
- **Frontend:** Go HTML templates + HTMX (or embedded React static build)
- **Real-time:** Server-Sent Events (SSE) or WebSocket (Go native)
- **Single Binary:** One executable serving both apps

**Key Libraries:**
- `github.com/mattn/go-sqlite3` (SQLite driver)
- `golang.org/x/crypto/bcrypt` (password hashing)
- `github.com/google/uuid` (UUID generation)
- `html/template` (Go standard library)
- Optional: `github.com/gofiber/fiber` (if using Fiber framework)

---

## 🎯 Migration Goals

1. **Single Codebase:** Everything in `/Users/efisiopittau/Project_1/alice-suite-go/`
2. **Self-Contained:** No external dependencies, no adapters, no Node.js
3. **Streamlined:** Direct SQLite access, no API layers
4. **Same Functionality:** All features from React apps preserved
5. **Better Performance:** Native Go performance, compiled binary

---

## 📊 Current Application Analysis

### Alice Reader App (`/APPS/alice-reader/`)

**Key Features:**
- User authentication (login/signup)
- Book verification code system
- Reading interface with page/section navigation
- Dictionary/glossary lookup
- AI assistance (chat interface)
- Help request submission
- Activity tracking
- Reading progress tracking
- Reading statistics

**Main Components:**
- `src/pages/` - Page components
- `src/components/` - Reusable UI components
- `src/services/` - API service layer
- `src/contexts/` - React contexts (Auth, etc.)
- `src/hooks/` - Custom React hooks

**Key Routes:**
- `/` - Home/Login
- `/reader` - Main reading interface
- `/verify` - Book verification
- `/admin` - Admin dashboard

### Consultant Dashboard App (`/APPS/alice-consultant-dashboard/`)

**Key Features:**
- Consultant authentication
- Reader activity monitoring
- Help request management
- Reader assignment system
- Analytics dashboard
- Real-time updates

**Main Components:**
- `src/pages/Consultant/` - Consultant-specific pages
- `src/services/` - Data services
- `src/hooks/` - Real-time hooks

**Key Routes:**
- `/` - Login
- `/dashboard` - Main dashboard
- `/readers` - Reader management
- `/help-requests` - Help request queue

---

## 🗄️ Database Schema

**Location:** `/Users/efisiopittau/Project_1/alice-suite/local-database/schema/create_tables.sql`

**Key Tables:**
- `profiles` - Users (readers and consultants)
- `books` - Book metadata
- `chapters` - Chapter information
- `sections` - Section content
- `alice_glossary` - Glossary terms
- `help_requests` - Help request queue
- `interactions` - Activity tracking
- `reading_progress` - User reading progress
- `reading_stats` - Reading statistics
- `consultant_assignments` - Reader-consultant assignments
- `verification_codes` - Book access codes

**Note:** The schema already exists and should be migrated as-is to the Go codebase.

---

## 🗂️ Code Organization & Cleanup Strategy

### Archive Old/Redundant Files

**Goal:** Keep the main codebase clean while preserving old code for reference

**Actions:**

1. **Create Archive Directory Structure**
   ```bash
   cd /Users/efisiopittau/Project_1/alice-suite-go
   mkdir -p archive/old-code
   mkdir -p archive/reference
   mkdir -p archive/deprecated
   ```

2. **Label and Archive Old Files**
   - All old Node.js/React code → `archive/old-code/`
   - Reference documentation → `archive/reference/`
   - Deprecated scripts/tools → `archive/deprecated/`
   - Add `_OLD_` or `_DEPRECATED_` prefix to file names
   - Add clear comments at top of archived files explaining why they're archived

3. **File Naming Convention for Archived Files**
   ```
   _OLD_<original-filename>
   _DEPRECATED_<original-filename>
   _REFERENCE_<original-filename>
   ```

4. **Create Archive Index**
   Create `archive/README.md` documenting:
   - What's archived and why
   - When it was archived
   - If/when it can be deleted
   - Reference links to new implementations

5. **Keep Main Codebase Clean**
   - Only active Go code in main directories
   - No Node.js files in main codebase
   - No React source files in main codebase
   - No adapter code in main codebase
   - Only current, working code

**Example Archive Structure:**
```
alice-suite-go/
├── archive/
│   ├── README.md                    # Archive index
│   ├── old-code/
│   │   ├── _OLD_nodejs-adapter/     # Old SQLite adapter
│   │   ├── _OLD_react-reader/       # Old React reader app
│   │   └── _OLD_react-consultant/   # Old React consultant app
│   ├── reference/
│   │   ├── _REFERENCE_schema.sql    # Original schema for reference
│   │   └── _REFERENCE_api-docs.md   # API documentation
│   └── deprecated/
│       ├── _DEPRECATED_old-migrations/
│       └── _DEPRECATED_test-scripts/
├── cmd/                              # Active Go code only
├── internal/                         # Active Go code only
└── data/                             # Active database only
```

---

## 🏗️ Migration Strategy

### Phase 0: Archive & Cleanup (BEFORE Migration)

**Goal:** Organize existing codebase before starting migration

**Steps:**

1. **Identify Redundant Files**
   - Old test scripts
   - Deprecated migration files
   - Old configuration files
   - Unused utilities
   - Old documentation that's no longer relevant

2. **Move to Archive**
   ```bash
   # Example: Archive old test scripts
   mv test_*.sh archive/deprecated/
   mv test_*.py archive/deprecated/
   
   # Example: Archive old migrations
   mv migrations/old_*.sql archive/deprecated/
   
   # Example: Archive old documentation
   mv docs/old_*.md archive/reference/
   ```

3. **Label Files Clearly**
   - Add header comments to archived files
   - Update README.md in archive directory
   - Document what replaces each archived file

4. **Clean Main Codebase**
   - Remove node_modules if present
   - Remove old build artifacts
   - Remove temporary files
   - Keep only active, relevant code

### Phase 1: Backend API Migration

**Goal:** Replace Node.js SQLite adapter with native Go handlers

**Steps:**

1. **Analyze Current API Endpoints**
   - Document all REST endpoints from React apps
   - Map Supabase client calls to actual HTTP requests
   - List all RPC functions used

2. **Create Go Handlers**
   - Migrate authentication endpoints (`/auth/v1/token`, `/auth/v1/user`)
   - Migrate REST endpoints (`/rest/v1/:table`)
   - Migrate RPC functions (`/rest/v1/rpc/:function`)
   - Ensure same response format as Supabase adapter

3. **Database Layer**
   - Use existing `internal/database/` package
   - Add missing query functions
   - Ensure foreign key constraints are enabled
   - Add proper error handling

### Phase 2: Frontend Migration Options

**Option A: Go Templates + HTMX (Recommended)**
- Server-rendered HTML using Go templates
- HTMX for dynamic interactions
- No JavaScript build step
- Fully self-contained

**Option B: Embedded React Static Build**
- Build React apps to static files
- Serve via Go `http.FileServer`
- Keep React for complex UI
- Still self-contained (no Node.js runtime)

**Option C: Hybrid Approach**
- Go templates for simple pages
- Embedded React for complex components
- Best of both worlds

**Recommendation:** Start with Option A, fall back to Option B for complex components.

### Phase 3: Real-time Features

**Current:** Socket.io/WebSocket

**Migration:**
- Use Go's native WebSocket support (`golang.org/x/net/websocket`)
- Or Server-Sent Events (SSE) for simpler use cases
- Migrate activity tracking, real-time updates

### Phase 4: Consolidation

**Goal:** Single binary serving both apps

**Structure:**
```
alice-suite-go/
├── cmd/
│   └── server/
│       └── main.go          # Single entry point
├── internal/
│   ├── handlers/
│   │   ├── reader.go        # Reader app handlers
│   │   ├── consultant.go    # Consultant app handlers
│   │   ├── api.go           # REST API handlers
│   │   └── auth.go          # Authentication handlers
│   ├── templates/
│   │   ├── reader/          # Reader app templates
│   │   └── consultant/      # Consultant app templates
│   ├── static/              # Static assets (CSS, JS, images)
│   └── database/            # Database layer (already exists)
└── data/
    └── alice-suite.db       # SQLite database
```

---

## 📝 Detailed Migration Instructions

### Step 1: Analyze Current React Applications

**Task:** Document all functionality, routes, and API calls

**Actions:**
1. List all React components and their purposes
2. Document all API service calls (`src/services/`)
3. List all routes and their handlers
4. Document authentication flow
5. List all database queries (via Supabase client)
6. Document real-time features

**Deliverable:** Complete feature inventory document

### Step 2: Set Up Go Project Structure

**Task:** Create proper Go project structure

**Actions:**
```bash
cd /Users/efisiopittau/Project_1/alice-suite-go

# Create directory structure
mkdir -p cmd/server
mkdir -p internal/handlers
mkdir -p internal/templates/reader
mkdir -p internal/templates/consultant
mkdir -p internal/static
mkdir -p internal/middleware
```

**Update `cmd/server/main.go`:**
```go
package main

import (
    "log"
    "net/http"
    "os"
    "path/filepath"
    
    "github.com/efisiopittau/alice-suite-go/internal/handlers"
    "github.com/efisiopittau/alice-suite-go/internal/database"
)

func main() {
    // Initialize database
    if err := database.Initialize("data/alice-suite.db"); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    
    // Setup routes
    mux := http.NewServeMux()
    
    // Static files
    staticDir := filepath.Join("internal", "static")
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
    
    // API routes
    handlers.SetupAPIRoutes(mux)
    
    // Reader app routes
    handlers.SetupReaderRoutes(mux)
    
    // Consultant app routes
    handlers.SetupConsultantRoutes(mux)
    
    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}
```

### Step 3: Migrate Authentication System

**Task:** Replace Supabase auth with Go-native auth

**Current Flow:**
1. User submits email/password
2. Supabase client calls `/auth/v1/token`
3. SQLite adapter validates and returns token
4. Token stored in browser, used for subsequent requests

**Go Implementation:**

**File:** `internal/handlers/auth.go`
```go
package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/efisiopittau/alice-suite-go/internal/database"
    "github.com/efisiopittau/alice-suite-go/pkg/auth"
    "github.com/google/uuid"
)

// HandleLogin handles POST /auth/v1/token
func HandleLogin(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    // Authenticate user
    user, err := auth.Login(req.Email, req.Password)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
    
    // Generate token (simple for now, upgrade to JWT later)
    token := uuid.New().String()
    expiresAt := time.Now().Add(24 * time.Hour)
    
    // Return Supabase-compatible response
    response := map[string]interface{}{
        "access_token": token,
        "token_type": "bearer",
        "expires_in": 86400,
        "expires_at": expiresAt.Unix(),
        "user": map[string]interface{}{
            "id": user.ID,
            "email": user.Email,
            "aud": "authenticated",
            "role": "authenticated",
            "user_metadata": map[string]string{
                "first_name": user.FirstName,
                "last_name": user.LastName,
            },
        },
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// HandleGetUser handles GET /auth/v1/user
func HandleGetUser(w http.ResponseWriter, r *http.Request) {
    // Extract token from Authorization header
    token := r.Header.Get("Authorization")
    if token == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Validate token and get user
    // (Implement token validation logic)
    
    // Return user info
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

### Step 4: Migrate REST API Endpoints

**Task:** Replace `/rest/v1/:table` endpoints

**Current:** Node.js adapter handles Supabase-style queries

**Go Implementation:**

**File:** `internal/handlers/api.go`
```go
package handlers

import (
    "net/http"
    "github.com/efisiopittau/alice-suite-go/internal/database"
)

// HandleRESTTable handles GET/POST /rest/v1/:table
func HandleRESTTable(w http.ResponseWriter, r *http.Request) {
    table := r.PathValue("table")
    
    switch r.Method {
    case http.MethodGet:
        handleGETTable(w, r, table)
    case http.MethodPost:
        handlePOSTTable(w, r, table)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func handleGETTable(w http.ResponseWriter, r *http.Request, table string) {
    // Parse query parameters (select, filter, order, limit, offset)
    // Build SQL query
    // Execute query
    // Return JSON response
}

func handlePOSTTable(w http.ResponseWriter, r *http.Request, table string) {
    // Parse request body
    // Validate foreign keys
    // Insert into database
    // Return inserted row
}
```

**Key Features to Implement:**
- Query parameter parsing (`select`, `eq`, `gte`, `like`, etc.)
- Join syntax support (`profiles:user_id(first_name,last_name)`)
- Foreign key validation
- Proper error handling

### Step 5: Migrate Frontend - Option A (Go Templates)

**Task:** Convert React components to Go HTML templates

**Example: Reader Login Page**

**React Component:** `src/pages/Login.tsx`
```tsx
// Current React component
```

**Go Template:** `internal/templates/reader/login.html`
```html
<!DOCTYPE html>
<html>
<head>
    <title>Alice Reader - Login</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-5">
        <h1>Alice Reader</h1>
        <form hx-post="/auth/v1/token" hx-target="#result">
            <input type="email" name="email" placeholder="Email" required>
            <input type="password" name="password" placeholder="Password" required>
            <button type="submit">Login</button>
        </form>
        <div id="result"></div>
    </div>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</body>
</html>
```

**Go Handler:**
```go
func HandleReaderLogin(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("internal/templates/reader/login.html"))
    tmpl.Execute(w, nil)
}
```

### Step 6: Migrate Frontend - Option B (Embedded React)

**Task:** Build React apps and serve as static files

**Actions:**
1. Build React apps:
   ```bash
   cd /Users/efisiopittau/Project_1/alice-suite/APPS/alice-reader
   npm run build
   
   cd ../alice-consultant-dashboard
   npm run build
   ```

2. Copy build outputs to Go project:
   ```bash
   cp -r dist/* /Users/efisiopittau/Project_1/alice-suite-go/internal/static/reader/
   cp -r dist/* /Users/efisiopittau/Project_1/alice-suite-go/internal/static/consultant/
   ```

3. Update API endpoints in React build to point to Go server

4. Serve static files:
   ```go
   mux.Handle("/reader/", http.StripPrefix("/reader/", http.FileServer(http.Dir("internal/static/reader"))))
   mux.Handle("/consultant/", http.StripPrefix("/consultant/", http.FileServer(http.Dir("internal/static/consultant"))))
   ```

### Step 7: Migrate Real-time Features

**Task:** Replace Socket.io with Go WebSocket

**Current:** Socket.io client connects to Node.js server

**Go Implementation:**

**File:** `internal/handlers/websocket.go`
```go
package handlers

import (
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins in development
    },
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()
    
    // Handle WebSocket messages
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            break
        }
        
        // Process message and send response
        conn.WriteMessage(messageType, message)
    }
}
```

### Step 8: Database Migration

**Task:** Ensure database schema matches

**Actions:**
1. Copy schema from `/Users/efisiopittau/Project_1/alice-suite/local-database/schema/create_tables.sql`
2. Create Go migration script
3. Run migrations on first startup
4. Ensure foreign keys are enabled

**File:** `internal/database/migrations.go`
```go
package database

import (
    "database/sql"
    "embed"
    "io/fs"
    "path/filepath"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(db *sql.DB) error {
    // Read all migration files
    // Execute them in order
    // Track which migrations have been run
    return nil
}
```

### Step 9: Testing & Validation

**Task:** Ensure all functionality works

**Checklist:**
- [ ] User authentication (login/signup)
- [ ] Book verification
- [ ] Reading interface
- [ ] Dictionary lookup
- [ ] AI assistance
- [ ] Help request submission
- [ ] Consultant dashboard
- [ ] Reader activity tracking
- [ ] Real-time updates
- [ ] Reading progress
- [ ] Statistics

### Step 10: Final Cleanup & Archive

**Task:** Archive old code and clean up codebase

**Actions:**

1. **Archive React/Node.js Source Code**
   ```bash
   cd /Users/efisiopittau/Project_1/alice-suite-go
   
   # Create archive structure
   mkdir -p archive/old-code/{react-reader,react-consultant,nodejs-adapter}
   
   # Copy old React reader app (for reference)
   cp -r /Users/efisiopittau/Project_1/alice-suite/APPS/alice-reader archive/old-code/react-reader/
   echo "# Archived: $(date)" > archive/old-code/react-reader/_ARCHIVED.txt
   echo "# Reason: Migrated to Go implementation" >> archive/old-code/react-reader/_ARCHIVED.txt
   
   # Copy old React consultant app (for reference)
   cp -r /Users/efisiopittau/Project_1/alice-suite/APPS/alice-consultant-dashboard archive/old-code/react-consultant/
   echo "# Archived: $(date)" > archive/old-code/react-consultant/_ARCHIVED.txt
   echo "# Reason: Migrated to Go implementation" >> archive/old-code/react-consultant/_ARCHIVED.txt
   
   # Copy old Node.js adapter (for reference)
   cp -r /Users/efisiopittau/Project_1/alice-suite/local-database/sqlite-adapter archive/old-code/nodejs-adapter/
   echo "# Archived: $(date)" > archive/old-code/nodejs-adapter/_ARCHIVED.txt
   echo "# Reason: Replaced with native Go handlers" >> archive/old-code/nodejs-adapter/_ARCHIVED.txt
   ```

2. **Label All Archived Files**
   ```bash
   # Add prefix to all archived files
   find archive/old-code -type f -name "*.js" -exec sh -c 'mv "$1" "$(dirname "$1")/_OLD_$(basename "$1")"' _ {} \;
   find archive/old-code -type f -name "*.ts" -exec sh -c 'mv "$1" "$(dirname "$1")/_OLD_$(basename "$1")"' _ {} \;
   find archive/old-code -type f -name "*.tsx" -exec sh -c 'mv "$1" "$(dirname "$1")/_OLD_$(basename "$1")"' _ {} \;
   ```

3. **Create Archive Documentation**
   ```bash
   cat > archive/README.md << 'EOF'
   # Archive Directory
   
   This directory contains old, deprecated, or reference code that is no longer part of the active codebase.
   
   ## Structure
   
   - `old-code/` - Old implementations (React, Node.js) kept for reference
   - `reference/` - Reference documentation and schemas
   - `deprecated/` - Deprecated scripts and utilities
   
   ## Archived Items
   
   ### React Reader App (`old-code/react-reader/`)
   - **Archived:** [DATE]
   - **Reason:** Migrated to Go implementation
   - **Replacement:** `cmd/server/main.go` + `internal/handlers/reader.go`
   - **Can Delete:** After 6 months if migration is stable
   
   ### React Consultant Dashboard (`old-code/react-consultant/`)
   - **Archived:** [DATE]
   - **Reason:** Migrated to Go implementation
   - **Replacement:** `cmd/server/main.go` + `internal/handlers/consultant.go`
   - **Can Delete:** After 6 months if migration is stable
   
   ### Node.js SQLite Adapter (`old-code/nodejs-adapter/`)
   - **Archived:** [DATE]
   - **Reason:** Replaced with native Go handlers
   - **Replacement:** `internal/handlers/api.go`
   - **Can Delete:** After 3 months if Go handlers are stable
   
   ## Cleanup Policy
   
   - Review archived items quarterly
   - Delete after confirmation that replacement is stable
   - Keep reference documentation indefinitely
   EOF
   ```

4. **Remove Old Files from Main Codebase**
   ```bash
   # Remove any old test files
   rm -f test_*.sh test_*.py test_*.js
   
   # Remove old build artifacts
   find . -name "node_modules" -type d -exec rm -rf {} + 2>/dev/null
   find . -name "dist" -type d -exec rm -rf {} + 2>/dev/null
   find . -name "*.log" -type f -delete
   
   # Remove temporary files
   find . -name "*.tmp" -type f -delete
   find . -name ".DS_Store" -type f -delete
   ```

5. **Update .gitignore**
   ```bash
   cat >> .gitignore << 'EOF'
   
   # Archive directory (keep for reference but don't track changes)
   archive/old-code/
   archive/deprecated/
   
   # Old build artifacts
   node_modules/
   dist/
   *.log
   EOF
   ```

### Step 11: Build & Deploy

**Task:** Create single binary

**Actions:**
```bash
cd /Users/efisiopittau/Project_1/alice-suite-go

# Ensure codebase is clean
go mod tidy

# Build binary
go build -o alice-suite cmd/server/main.go

# Run
./alice-suite
```

**Result:** Single executable serving both apps on port 8080

**Verify Clean Codebase:**
```bash
# Check that no old files remain in main directories
find cmd internal pkg -name "*.js" -o -name "*.ts" -o -name "*.tsx" | grep -v archive
# Should return nothing

# Check that archive exists
ls -la archive/
# Should show archived directories
```

---

## 🔧 Technical Implementation Details

### API Compatibility Layer

**Goal:** Maintain compatibility with Supabase client calls

**Implementation:**
- Map Supabase client methods to Go handlers
- Return same JSON response format
- Support same query parameters
- Handle authentication tokens the same way

### Error Handling

**Current:** Supabase-style error responses

**Go Implementation:**
```go
type APIError struct {
    Message string `json:"message"`
    Code    string `json:"code"`
    Details string `json:"details,omitempty"`
}

func sendError(w http.ResponseWriter, status int, err APIError) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(err)
}
```

### Middleware

**Required Middleware:**
- CORS handling
- Authentication token validation
- Request logging
- Error recovery

**File:** `internal/middleware/middleware.go`
```go
package middleware

import (
    "log"
    "net/http"
)

func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate token
        // Set user context
        next.ServeHTTP(w, r)
    })
}
```

---

## 📦 Dependencies

### Go Modules

**File:** `go.mod`
```go
module github.com/efisiopittau/alice-suite-go

go 1.21

require (
    github.com/mattn/go-sqlite3 v1.14.18
    github.com/google/uuid v1.5.0
    golang.org/x/crypto v0.17.0
    github.com/gorilla/websocket v1.5.1
)
```

---

## ✅ Success Criteria

**Migration is complete when:**

1. ✅ Single Go binary runs both apps
2. ✅ No Node.js dependencies
3. ✅ Direct SQLite access (no adapter)
4. ✅ All React app features work
5. ✅ Authentication works
6. ✅ Real-time features work
7. ✅ Database schema matches
8. ✅ All API endpoints respond correctly
9. ✅ No external service dependencies
10. ✅ Self-contained deployment
11. ✅ **All old/redundant files archived and labeled**
12. ✅ **Main codebase contains only active Go code**
13. ✅ **Archive directory properly documented**
14. ✅ **No old files in main directories (cmd/, internal/, pkg/)**

---

## 🚀 Quick Start After Migration

```bash
# Build
cd /Users/efisiopittau/Project_1/alice-suite-go
go build -o alice-suite cmd/server/main.go

# Run
./alice-suite

# Access
# Reader: http://localhost:8080/reader
# Consultant: http://localhost:8080/consultant
```

---

## 📚 Additional Resources

- Go HTTP Server: https://pkg.go.dev/net/http
- SQLite Driver: https://github.com/mattn/go-sqlite3
- Go Templates: https://pkg.go.dev/html/template
- HTMX: https://htmx.org/ (if using Option A)
- WebSocket: https://pkg.go.dev/golang.org/x/net/websocket

---

## 🎯 Next Steps for Agent

1. **Read this document completely**
2. **Archive old/redundant files** (Phase 0 - BEFORE starting migration)
3. **Analyze current React applications** (Step 1)
4. **Set up Go project structure** (Step 2)
5. **Migrate authentication** (Step 3)
6. **Migrate API endpoints** (Step 4)
7. **Choose frontend migration approach** (Step 5 or 6)
8. **Migrate real-time features** (Step 7)
9. **Test thoroughly** (Step 9)
10. **Final cleanup and archive** (Step 10)
11. **Build and deploy** (Step 11)

**IMPORTANT:** Always archive old files BEFORE removing them. Keep them labeled and organized in the `archive/` directory for reference.

**Estimated Time:** 2-3 days for complete migration

**Priority:** High - This will simplify the entire codebase and remove all external dependencies.

**Code Organization:** All old/redundant files MUST be archived and labeled before removal. The main codebase should contain ONLY active Go code. Archive directory serves as reference and can be cleaned up after migration is confirmed stable.

---

**End of Migration Guide**

