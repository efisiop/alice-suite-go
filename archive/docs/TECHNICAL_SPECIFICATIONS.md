# Alice Suite - Technical Specifications

**Created:** 2025-01-20  
**Status:** Technical specifications for implementation  
**Based on:** REQUIREMENTS.md and current database structure

---

## 🏗️ System Architecture

### High-Level Architecture

```
┌─────────────────┐
│   Web Frontend  │  (HTML/CSS/JavaScript)
│   (Browser)     │
└────────┬────────┘
         │ HTTP/REST API
┌────────▼────────┐
│  Go HTTP Server │  (cmd/reader/main.go)
│  Port 8080      │
└────────┬────────┘
         │
┌────────▼────────┐
│   Handlers      │  (internal/handlers/)
│   Layer         │
└────────┬────────┘
         │
┌────────▼────────┐
│   Services      │  (internal/services/)
│   Layer         │
└────────┬────────┘
         │
┌────────▼────────┐
│   Database      │  (internal/database/)
│   Layer         │
└────────┬────────┘
         │
┌────────▼────────┐
│   SQLite DB     │  (data/alice-suite.db)
└─────────────────┘
```

### Technology Stack

**Backend:**
- **Language:** Go 1.21+
- **Database:** SQLite 3 (file-based)
- **HTTP Server:** Standard library `net/http`
- **Password Hashing:** `golang.org/x/crypto/bcrypt`
- **UUID Generation:** `github.com/google/uuid`
- **SQLite Driver:** `github.com/mattn/go-sqlite3`

**Frontend:**
- **Approach:** Web-based (HTML/CSS/JavaScript)
- **Deployment:** Served as static files or Go templates
- **Mobile:** Responsive design (mobile-first)

**External Services:**
- **AI Service:** Moonshot AI (Kimi) or Anthropic API
- **Configuration:** Environment variables

---

## 📊 Database Schema

### Current Structure (Must Be Maintained)

#### 1. Pages Table
```sql
CREATE TABLE pages (
  id TEXT PRIMARY KEY,
  book_id TEXT NOT NULL,
  page_number INTEGER NOT NULL,
  chapter_id TEXT,                    -- Nullable, only if chapter starts on this page
  chapter_title TEXT,                 -- Nullable, only if chapter starts on this page
  content TEXT NOT NULL,               -- Full page content (concatenated sections)
  word_count INTEGER NOT NULL,         -- Total words on page (180-200 average)
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (book_id) REFERENCES books(id),
  FOREIGN KEY (chapter_id) REFERENCES chapters(id),
  UNIQUE(book_id, page_number)
);
```

**Key Constraints:**
- Page numbers: 1-100 (for full book), 1-17 (for test chapters)
- Average 180-200 words per page
- 4-5 sections per page

#### 2. Sections Table
```sql
CREATE TABLE sections (
  id TEXT PRIMARY KEY,
  page_id TEXT NOT NULL,
  page_number INTEGER NOT NULL,       -- Denormalized for quick lookup
  section_number INTEGER NOT NULL,    -- 1, 2, 3, etc. within page
  content TEXT NOT NULL,
  word_count INTEGER NOT NULL,        -- Average 40 words (range: 35-45)
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE,
  UNIQUE(page_id, section_number)
);
```

**Key Constraints:**
- Average 40 words per section (target range: 35-45 words)
- Sections numbered sequentially within each page
- Must maintain word count distribution

#### 3. Glossary Tables

**alice_glossary:**
```sql
CREATE TABLE alice_glossary (
  id TEXT PRIMARY KEY,
  book_id TEXT NOT NULL,
  term TEXT NOT NULL,
  definition TEXT NOT NULL,
  source_sentence TEXT,
  example TEXT,
  chapter_reference TEXT,
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (book_id) REFERENCES books(id),
  UNIQUE(book_id, term)
);
```

**glossary_section_links (Junction Table):**
```sql
CREATE TABLE glossary_section_links (
  id TEXT PRIMARY KEY,
  glossary_id TEXT NOT NULL,
  section_id TEXT NOT NULL,
  page_number INTEGER NOT NULL,       -- Denormalized
  section_number INTEGER NOT NULL,   -- Denormalized
  term TEXT NOT NULL,                 -- Denormalized for quick lookup
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (glossary_id) REFERENCES alice_glossary(id),
  FOREIGN KEY (section_id) REFERENCES sections(id),
  UNIQUE(glossary_id, section_id)
);
```

**Key Features:**
- 1,209 glossary terms currently loaded
- Terms linked to sections where they appear
- Efficient lookup via indexes

#### 4. Other Core Tables

**users:**
```sql
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  first_name TEXT,
  last_name TEXT,
  role TEXT CHECK (role IN ('reader', 'consultant')) DEFAULT 'reader',
  is_verified INTEGER DEFAULT 0,
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now'))
);
```

**books:**
```sql
CREATE TABLE books (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  description TEXT,
  total_pages INTEGER NOT NULL,
  created_at TEXT DEFAULT (datetime('now'))
);
```

**chapters:**
```sql
CREATE TABLE chapters (
  id TEXT PRIMARY KEY,
  book_id TEXT NOT NULL,
  title TEXT NOT NULL,
  number INTEGER NOT NULL,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (book_id) REFERENCES books(id),
  UNIQUE(book_id, number)
);
```

**reading_progress:**
```sql
CREATE TABLE reading_progress (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  book_id TEXT NOT NULL,
  page_id TEXT,
  page_number INTEGER,
  section_number INTEGER,
  last_read_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (book_id) REFERENCES books(id),
  FOREIGN KEY (page_id) REFERENCES pages(id)
);
```

**vocabulary_lookups:**
```sql
CREATE TABLE vocabulary_lookups (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  book_id TEXT NOT NULL,
  word TEXT NOT NULL,
  definition TEXT,
  page_id TEXT,
  section_number INTEGER,
  context TEXT,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (book_id) REFERENCES books(id)
);
```

**verification_codes:**
```sql
CREATE TABLE verification_codes (
  code TEXT PRIMARY KEY,
  book_id TEXT NOT NULL,
  is_used INTEGER DEFAULT 0,
  used_by TEXT,
  used_at TEXT,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (book_id) REFERENCES books(id)
);
```

### Database Indexes

```sql
-- Performance indexes
CREATE INDEX idx_pages_book_id ON pages(book_id);
CREATE INDEX idx_sections_page_id ON sections(page_id);
CREATE INDEX idx_sections_page_number ON sections(page_number, section_number);
CREATE INDEX idx_glossary_links_term ON glossary_section_links(term);
CREATE INDEX idx_glossary_links_section ON glossary_section_links(section_id);
CREATE INDEX idx_glossary_links_page ON glossary_section_links(page_number, section_number);
CREATE INDEX idx_glossary_term ON alice_glossary(term);
CREATE INDEX idx_reading_progress_user_book ON reading_progress(user_id, book_id);
CREATE INDEX idx_vocabulary_user_book ON vocabulary_lookups(user_id, book_id);
```

---

## 🔌 API Endpoint Specifications

### Base URL
```
http://localhost:8080/api
```

### Authentication Endpoints

#### POST /api/auth/register
Register a new user.

**Request:**
```json
{
  "email": "student@example.com",
  "password": "password123",
  "first_name": "Alice",
  "last_name": "Student"
}
```

**Response (201 Created):**
```json
{
  "id": "user-uuid",
  "email": "student@example.com",
  "first_name": "Alice",
  "last_name": "Student",
  "role": "reader",
  "created_at": "2025-01-20T10:00:00Z"
}
```

**Errors:**
- `400 Bad Request` - Invalid input
- `409 Conflict` - Email already exists

---

#### POST /api/auth/login
Login user.

**Request:**
```json
{
  "email": "student@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "user-uuid",
    "email": "student@example.com",
    "first_name": "Alice",
    "last_name": "Student"
  },
  "token": "jwt-token-here"
}
```

**Errors:**
- `401 Unauthorized` - Invalid credentials

---

### Book & Content Endpoints

#### GET /api/books
Get all available books.

**Response (200 OK):**
```json
[
  {
    "id": "alice-in-wonderland",
    "title": "Alice's Adventures in Wonderland",
    "author": "Lewis Carroll",
    "description": "Physical book companion",
    "total_pages": 100
  }
]
```

---

#### GET /api/pages?book_id={book_id}&page_number={page_number}
Get a specific page with sections.

**Query Parameters:**
- `book_id` (required) - Book identifier
- `page_number` (required) - Page number (1-17 for test, 1-100 for full)

**Response (200 OK):**
```json
{
  "id": "page-1",
  "page_number": 1,
  "chapter_id": "chapter-1",
  "chapter_title": "Chapter 1: Down the Rabbit-Hole",
  "word_count": 148,
  "sections": [
    {
      "id": "page-1-section-1",
      "section_number": 1,
      "content": "Alice was beginning to get very tired...",
      "word_count": 38,
      "glossary_terms": [
        {
          "term": "beginning",
          "definition": "the act of starting something"
        }
      ]
    }
  ]
}
```

**Errors:**
- `400 Bad Request` - Missing parameters
- `404 Not Found` - Page not found

---

#### GET /api/sections?page_id={page_id}
Get sections for a page.

**Query Parameters:**
- `page_id` (required) - Page identifier

**Response (200 OK):**
```json
[
  {
    "id": "page-1-section-1",
    "page_number": 1,
    "section_number": 1,
    "content": "Alice was beginning...",
    "word_count": 38
  }
]
```

---

### Dictionary/Glossary Endpoints

#### POST /api/dictionary/lookup
Lookup word definition (glossary priority).

**Request:**
```json
{
  "book_id": "alice-in-wonderland",
  "word": "rabbit",
  "page_number": 1,
  "section_number": 1,
  "user_id": "user-uuid"
}
```

**Response (200 OK):**
```json
{
  "id": "glossary-29",
  "term": "rabbit",
  "definition": "any of various burrowing animals...",
  "example": "CHAPTER I. Down the Rabbit-Hole...",
  "chapter_reference": "I"
}
```

**If not found:**
```json
{
  "word": "rabbit",
  "definition": "Word not found in glossary"
}
```

---

#### GET /api/dictionary/section/{section_id}/terms
Get all glossary terms for a section.

**Response (200 OK):**
```json
[
  {
    "term": "beginning",
    "definition": "the act of starting something"
  },
  {
    "term": "book",
    "definition": "an object consisting of a number of pages bound together"
  }
]
```

---

### AI Assistance Endpoints

#### POST /api/ai/ask
Ask AI a question about the content.

**Request:**
```json
{
  "user_id": "user-uuid",
  "book_id": "alice-in-wonderland",
  "interaction_type": "explain",
  "question": "What does this passage mean?",
  "page_number": 1,
  "section_number": 1,
  "context": "Selected text from section"
}
```

**Response (200 OK):**
```json
{
  "id": "interaction-uuid",
  "user_id": "user-uuid",
  "interaction_type": "explain",
  "question": "What does this passage mean?",
  "response": "This passage describes...",
  "created_at": "2025-01-20T10:00:00Z"
}
```

**Interaction Types:**
- `explain` - Explain passages/concepts
- `quiz` - Generate quiz questions
- `simplify` - Simplify/rephrase text
- `definition` - Get definitions
- `chat` - General chat assistance

---

### Help Request Endpoints

#### POST /api/help/request
Create a help request (Tier 3: Human consultant).

**Request:**
```json
{
  "user_id": "user-uuid",
  "book_id": "alice-in-wonderland",
  "content": "I don't understand this passage",
  "page_number": 1,
  "section_number": 1,
  "context": "Selected text"
}
```

**Response (201 Created):**
```json
{
  "id": "request-uuid",
  "user_id": "user-uuid",
  "status": "pending",
  "content": "I don't understand this passage",
  "created_at": "2025-01-20T10:00:00Z"
}
```

---

### Progress Endpoints

#### GET /api/progress?user_id={user_id}&book_id={book_id}
Get reading progress.

**Response (200 OK):**
```json
{
  "user_id": "user-uuid",
  "book_id": "alice-in-wonderland",
  "last_page_number": 5,
  "last_section_number": 2,
  "pages_read": 5,
  "sections_read": 22,
  "total_pages": 17,
  "total_sections": 77,
  "percentage_complete": 29.4
}
```

---

#### PUT /api/progress
Update reading progress.

**Request:**
```json
{
  "user_id": "user-uuid",
  "book_id": "alice-in-wonderland",
  "page_number": 5,
  "section_number": 2
}
```

**Response (200 OK):**
```json
{
  "message": "Progress updated",
  "page_number": 5,
  "section_number": 2
}
```

---

## 🔐 Authentication & Authorization

### Authentication Flow

1. **Registration:**
   - User provides: email, password, first_name, last_name
   - Password hashed with bcrypt (cost: 10)
   - User created in database
   - Return user object (no token yet)

2. **Login:**
   - User provides: email, password
   - Verify password hash
   - Generate JWT token (or session token)
   - Return user + token

3. **Authorization:**
   - Token included in `Authorization: Bearer {token}` header
   - Middleware validates token
   - Extract user_id from token
   - Attach to request context

### JWT Token Structure (Recommended)

```json
{
  "user_id": "user-uuid",
  "email": "student@example.com",
  "role": "reader",
  "exp": 1234567890,
  "iat": 1234567890
}
```

### Book Access Authorization

- User must enter verification code
- Code validated against `verification_codes` table
- Code marked as used (or tracked)
- Access granted to book content

---

## 🎨 Frontend Specifications

### Page Structure

#### 1. Login/Register Page
- Simple form: Email, Password, Name fields
- Toggle between login/register
- Error messages displayed clearly

#### 2. Book Code Entry Page
- Input field for verification code
- Submit button
- Error message if invalid

#### 3. Welcome Screen (First Time)
- Modal or dedicated page
- Explains how to use app
- "Got it" button to dismiss
- Accessible from menu

#### 4. Main Reader Page
- **Page Selection:**
  - Input field: "Which page are you reading?"
  - Page number validation (1-17 for test)
  - "Go" button

- **Section Display:**
  - Show all sections for selected page
  - Each section: number + word count
  - Clickable sections
  - Selected section highlighted

- **Content Display:**
  - Section text displayed clearly
  - Words with glossary definitions highlighted (yellow)
  - Hover shows tooltip with definition
  - Click word shows full definition modal

- **Glossary Terms List:**
  - Below section content
  - List of all glossary terms in section
  - Clickable to see definitions

#### 5. AI Help Interface
- Text selection from section
- Question input field
- "Ask AI" button
- Response displayed below
- Follow-up questions supported

#### 6. Help Request Interface
- "Ask for Help" button
- Question input field
- Context automatically included
- Submit button
- Confirmation message

#### 7. Progress Dashboard
- Pages read / Total pages
- Sections read / Total sections
- Percentage complete
- Last page viewed
- Vocabulary list

### UI/UX Guidelines

- **Mobile-First:** Responsive design, works on phones/tablets
- **Simple:** Clean, uncluttered interface
- **Fast:** Quick page loads, instant word lookups
- **Accessible:** WCAG 2.1 AA compliance
- **Clear:** Obvious navigation, clear error messages

---

## 🔄 Implementation Phases

### Phase 1: Foundation (Current Focus)
**Goal:** Basic reader functionality with first 3 chapters

**Tasks:**
1. ✅ Database schema created
2. ✅ Pages/sections structured (40 words average)
3. ✅ Glossary loaded and linked
4. ✅ Basic API endpoints implemented
5. ⏳ Frontend: Page/section selection
6. ⏳ Frontend: Word lookup with glossary
7. ⏳ Frontend: Last page memory
8. ⏳ Testing and bug fixes

### Phase 2: AI Integration
**Goal:** Add AI assistance

**Tasks:**
1. ✅ AI service implemented (backend)
2. ⏳ Frontend: AI question interface
3. ⏳ Context-aware prompts
4. ⏳ Interaction history
5. ⏳ Testing

### Phase 3: Human Help
**Goal:** Add consultant help system

**Tasks:**
1. ✅ Help request service (backend)
2. ⏳ Frontend: Help request interface
3. ⏳ Consultant dashboard (separate)
4. ⏳ Notification system
5. ⏳ Testing

### Phase 4: Progress Tracking
**Goal:** Track reading progress

**Tasks:**
1. ⏳ Progress service implementation
2. ⏳ Frontend: Progress dashboard
3. ⏳ Vocabulary tracking
4. ⏳ Statistics display
5. ⏳ Testing

### Phase 5: Full Book
**Goal:** Expand to all 12 chapters

**Tasks:**
1. ⏳ Load remaining chapters
2. ⏳ Update page numbers (1-100)
3. ⏳ Re-link glossary terms
4. ⏳ Test all features with full book
5. ⏳ Performance optimization

---

## 🧪 Testing Strategy

### Unit Tests
- Service layer functions
- Database query functions
- Authentication functions
- Word lookup logic

### Integration Tests
- API endpoint testing
- Database operations
- Glossary linking
- AI service integration

### End-to-End Tests
- User registration → Login → Page selection → Word lookup
- AI question flow
- Help request flow
- Progress tracking

### Performance Tests
- Page load times (< 2 seconds)
- Word lookup response (< 500ms)
- Database query performance
- Concurrent user handling

---

## 🔒 Security Considerations

1. **Password Security:**
   - Bcrypt hashing (cost: 10)
   - Minimum 8 characters
   - No password storage in plain text

2. **SQL Injection Prevention:**
   - Parameterized queries only
   - No string concatenation in SQL

3. **Input Validation:**
   - Validate all user inputs
   - Sanitize text inputs
   - Validate page numbers (range checks)

4. **Session Management:**
   - Secure token generation
   - Token expiration
   - Secure token storage (httpOnly cookies recommended)

5. **API Security:**
   - Rate limiting (future)
   - CORS configuration
   - Error messages don't leak sensitive info

---

## 📝 Environment Variables

```bash
# Database
DB_PATH=data/alice-suite.db

# Server
PORT=8080
ENV=development  # development, production

# AI Service
MOONSHOT_API_KEY=your_api_key_here
ANTHROPIC_AUTH_TOKEN=your_token_here  # Alternative
ANTHROPIC_BASE_URL=https://api.moonshot.cn/v1

# JWT (if implemented)
JWT_SECRET=your_secret_key_here
JWT_EXPIRY=24h
```

---

## 🚀 Deployment Considerations

### Development
- SQLite file-based database
- Local file storage
- Development server on localhost:8080

### Production (Future)
- Database backup strategy
- Static file serving (CDN or Go embed)
- HTTPS configuration
- Environment variable management
- Logging and monitoring

---

## 📊 Key Metrics to Track

1. **Performance:**
   - API response times
   - Database query times
   - Page load times

2. **Usage:**
   - User registrations
   - Pages viewed
   - Words looked up
   - AI questions asked
   - Help requests submitted

3. **Quality:**
   - Error rates
   - User feedback
   - Bug reports

---

## ✅ Success Criteria

### Technical Success
- [ ] All API endpoints working
- [ ] Database queries optimized
- [ ] Frontend responsive and fast
- [ ] Glossary linking accurate
- [ ] AI integration working
- [ ] Help system functional

### User Success
- [ ] Users can register/login easily
- [ ] Page/section selection intuitive
- [ ] Word lookup instant and accurate
- [ ] AI help useful and clear
- [ ] Progress tracking accurate
- [ ] Overall positive user experience

---

**Next Steps:** Begin Phase 1 frontend implementation based on these specifications.


