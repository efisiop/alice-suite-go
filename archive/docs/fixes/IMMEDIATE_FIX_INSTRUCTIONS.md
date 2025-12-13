# ⚡ Immediate Fix Instructions

## 🚨 **CRITICAL: Build Failures** - Must Fix Now

The codebase currently has **build failures** due to undefined handler functions. This is a **blocking issue** that prevents compilation.

### ❌ Current Error:
```
cmd/reader/main.go:31:45: undefined: handlers.Login
cmd/reader/main.go:32:48: undefined: handlers.Register
cmd/reader/main.go:33:40: undefined: handlers.GetBooks
// ... 8 total undefined functions
```

---

## 🛠️ **FIX 1: Create Missing Handler Functions**

### Step 1: Create helper functions in `internal/handlers/helpers.go`

Create a new file to bridge the reader module with existing functionality:

```go
package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/efisiopittau/alice-suite-go/pkg/auth"
)

// Login wraps the auth Login function to match API standards
func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Delegate to existing HandleLogin
    // We need to ensure the Login function accepts the right type
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := auth.Login(req.Email, req.Password)
    if err != nil {
        if err == auth.ErrInvalidCredentials {
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }
        HandleError(w, err)
        return
    }

    // Generate JWT token
    token, err := auth.GenerateJWT(user.ID, user.Email, user.Role)
    if err != nil {
        HandleError(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "access_token": token,
        "message": "Login successful",
    })
}

// Register wraps the auth Register function to match API standards
func Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        Email     string `json:"email"`
        Password  string `json:"password"`
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := auth.Register(req.Email, req.Password, req.FirstName, req.LastName)
    if err != nil {
        if err == auth.ErrUserExists {
            http.Error(w, "User already exists", http.StatusConflict)
            return
        }
        HandleError(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

### Step 2: Create shortcut functions for book-related endpoints

Add to `internal/handlers/api.go`:

```go
// GetBooks provides book listing for reader API
func GetBooks(w http.ResponseWriter, r *http.Request) {
    HandleBooks(w, r) // Reuse existing implementation
}

// GetChapters provides chapter listing for reader API
func GetChapters(w http.ResponseWriter, r *http.Request) {
    HandleChapters(w, r) // Reuse existing implementation
}

// GetSections provides section listing for reader API
func GetSections(w http.ResponseWriter, r *http.Request) {
    HandleSections(w, r) // Reuse existing implementation
}

// GetPage provides page content for reader API
func GetPage(w http.ResponseWriter, r *http.Request) {
    HandlePages(w, r) // Reuse existing implementation
}
```

### Step 3: Create dictionary function shortcuts

Add to `internal/handlers/api.go`:

```go
// LookupWord provides dictionary lookup for reader API
func LookupWord(w http.ResponseWriter, r *http.Request) {
    HandleLookupWord(w, r) // Reuse existing implementation
}

// GetSectionGlossaryTerms provides glossary for reader API
func GetSectionGlossaryTerms(w http.ResponseWriter, r *http.Request) {
    HandleGetSectionGlossaryTerms(w, r) // Reuse existing implementation
}
```

### Step 4: Create AI and help function shortcuts

Add to `internal/handlers/api.go`:

```go
// AskAI provides AI service for reader API
func AskAI(w http.ResponseWriter, r *http.Request) {
    HandleAskAI(w, r) // Reuse existing implementation
}

// CreateHelpRequest provides help system for reader API
func CreateHelpRequest(w http.ResponseWriter, r *http.Request) {
    HandleCreateHelpRequest(w, r) // Reuse existing implementation
}
```

---

## 🛠️ **FIX 2: Update Route Registration**

Update `cmd/reader/main.go:30-45` to match the function names:

```go
// Reader API routes (register BEFORE static files to avoid conflicts)
mux.HandleFunc("/api/health", handlers.HealthCheck)
mux.HandleFunc("/api/auth/login", handlers.Login)  // Now defined
mux.HandleFunc("/api/auth/register", handlers.Register)  // Now defined
mux.HandleFunc("/api/books", handlers.GetBooks)  // Now defined
mux.HandleFunc("/api/chapters", handlers.GetChapters)  // Now defined
mux.HandleFunc("/api/sections", handlers.GetSections)  // Now defined
mux.HandleFunc("/api/pages", handlers.GetPage)  // Now defined
mux.HandleFunc("/api/dictionary/lookup", handlers.LookupWord)  // Now defined
mux.HandleFunc("/api/dictionary/section/", handlers.GetSectionGlossaryTerms)  // Now defined
mux.HandleFunc("/api/ai/ask", handlers.AskAI)  // Now defined
mux.HandleFunc("/api/help/request", handlers.CreateHelpRequest)  // Now defined
mux.HandleFunc("/api/progress", handlers.GetProgress)
mux.HandleFunc("/api/user", handlers.GetUserProfile)
mux.HandleFunc("/api/interactions", handlers.GetInteractions)
mux.HandleFunc("/api/track", handlers.TrackEvent)
```

---

## 🧪 **Verify Fix**

### Step 1: Check compilation
```bash
go build ./...
```

### Step 2: Run tests
```bash
go test ./... -v
```

---

## 🚨 **Test Database Fix** (High Priority)

### Create `internal/database/testutils.go`:

```go
package database

import (
    "database/sql"
    "log"
    "os"
    "testing"
    "_/github.com/mattn/go-sqlite3"
)

// TestDB provides test database access
type TestDB struct {
    *sql.DB
}

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *TestDB {
    // Close any existing connection
    if DB != nil {
        DB.Close()
    }

    // Create in-memory SQLite database for tests
    testDB, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
    if err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }

    // Save as global DB temporarily for tests
    DB = testDB

    // Run migrations on test database
    testMigrations()

    return &TestDB{testDB}
}

// testMigrations runs database migrations for test database
func testMigrations() {
    // Run migrations similar to cmd/migrate/main.go but simplified
    migrationFiles := []string{
        "migrations/001_initial_schema.sql",
        "migrations/002_seed_first_3_chapters.sql",
    }

    for _, file := range migrationFiles {
        sql, err := os.ReadFile(file)
        if err != nil {
            log.Printf("Warning: Could not read migration %s: %v", file, err)
            continue
        }

        _, err = DB.Exec(string(sql))
        if err != nil {
            log.Printf("Warning: Migration %s failed: %v", file, err)
        }
    }
}

// Cleanup closes the test database
func (t *TestDB) Cleanup() {
    if t.DB != nil {
        t.DB.Close()
    }
}

```

### Update test files to use test database:

Add to each failing test file (book_service_test.go, api_test.go):

```go
package handlers

import (
    "testing"
    "internal/database"
)

func TestMain(m *testing.M) {
    // Set up test database
    tdb := database.SetupTestDB(nil)
    defer tdb.Cleanup()

    // Run tests
    code := m.Run()
    os.Exit(code)
}
```

---

## 🎯 **Success Criteria**

### ✅ Fix Complete When:
- [ ] `go build ./...` returns **no errors**
- [ ] All external handler functions are properly defined
- [ ] Test database setup works correctly
- [ ] **17+ tests passing** (currently 17/22 total)

### Expected Output:
```bash
$ go build ./...
# ✅ Success - no build errors

$ go test ./... -v
pkg/auth: passing (7/7) ✅
internal/middleware: passing (7/7) ✅
internal/services: passing (3/3) ✅
internal/handlers: passing (5/6) ✅
```

---

## ⏱️ **Estimated Completion Time: 1-2 hours**

The fixes are straightforward and involve:
1. Creating wrapper functions (30 minutes)
2. Fixing test database setup (45 minutes)
3. Testing and validation (30 minutes)

**This will resolve the critical blocking issues and make the codebase fully operational.**

---

*Execute these fixes in order to unblock the build and test infrastructure.* ❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋❋"**file_path:"/Users/efisiopittau/Project_1/alice-suite-go/IMMEDIATE_FIX_INSTRUCTIONS.md"}