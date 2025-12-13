# Feature Inventory: React Applications Analysis

**Created:** 2025-01-23  
**Purpose:** Complete documentation of all functionality, routes, API calls, and features from React/TypeScript applications  
**Source:** `/Users/efisiopittau/Project_1/alice-suite/APPS/`

---

## 📋 Table of Contents

1. [Application Overview](#application-overview)
2. [Alice Reader App - Routes & Pages](#alice-reader-app---routes--pages)
3. [Consultant Dashboard App - Routes & Pages](#consultant-dashboard-app---routes--pages)
4. [API Service Calls](#api-service-calls)
5. [Database Queries](#database-queries)
6. [Authentication Flow](#authentication-flow)
7. [Real-time Features](#real-time-features)
8. [Key Components](#key-components)
9. [Feature Summary](#feature-summary)

---

## Application Overview

### Alice Reader App
- **Location:** `/Users/efisiopittau/Project_1/alice-suite/APPS/alice-reader/`
- **Port:** 5173 (development)
- **Purpose:** Main reading interface for users to read books, look up words, get AI assistance, and request help
- **Technology:** React 18 + TypeScript + Material-UI + Vite

### Consultant Dashboard App
- **Location:** `/Users/efisiopittau/Project_1/alice-suite/APPS/alice-consultant-dashboard/`
- **Port:** 5174 (development)
- **Purpose:** Dashboard for consultants to monitor readers, manage help requests, and provide assistance
- **Technology:** React 18 + TypeScript + Material-UI + Vite

---

## Alice Reader App - Routes & Pages

### Public Routes (No Authentication Required)
- **`/`** - Landing page (public welcome page)
- **`/login`** - User login page
- **`/login-legacy`** - Legacy login route (redirects to `/login`)
- **`/register`** - User registration page
- **`/forgot-password`** - Password reset page
- **`/consultant-landing`** - Consultant landing page (public)

### Auth Routes (Require Authentication, Not Verification)
- **`/verify`** - Book verification page (user must enter verification code)
- **`/welcome`** - Welcome page after registration

### Reader Routes (Require Authentication + Verification)
- **`/reader`** - Main reader dashboard
- **`/reader/interaction`** - Main reading interaction page
- **`/reader/book/:bookId`** - Book-specific reading interface
- **`/reader/:bookId/page/:pageNumber`** - Specific page view
- **`/reader/statistics`** - Reading statistics and progress

### Admin Routes
- **`/admin`** - Admin dashboard
- **`/service-status`** - Service status check page
- **`/supabase-test`** - Supabase connection test page

### Development/Test Routes (DEV mode only)
- **`/test`** - Test page
- **`/test-reader-page/:pageNumber`** - Test reader page
- **`/test-main-interaction`** - Test main interaction page
- **`/test-simple`** - Simple test page
- **`/test-supabase`** - Supabase test page
- **`/dictionary-test`** - Dictionary functionality test
- **`/alice-glossary-demo`** - Alice glossary demo page

### Page Components
- `LandingPage.tsx` - Public landing page
- `LoginPage.tsx` - User login
- `RegisterPage.tsx` - User registration
- `VerifyPage.tsx` - Book verification
- `WelcomePage.tsx` - Post-registration welcome
- `ReaderDashboard.tsx` - Main reader dashboard
- `ReaderPage.tsx` - Individual page view
- `MainInteractionPage.tsx` - Main reading interface
- `ReaderStatistics.tsx` - Reading statistics
- `AdminDashboard.tsx` - Admin interface
- `ForgotPasswordPage.tsx` - Password reset

---

## Consultant Dashboard App - Routes & Pages

### Public Routes
- **`/consultant/login`** - Consultant login page

### Legacy Redirects (Redirect to `/consultant/login`)
- `/login` → `/consultant/login`
- `/register` → `/consultant/login`
- `/forgot-password` → `/consultant/login`
- `/bypass` → `/consultant/login`
- `/verify` → `/consultant/login`

### Protected Consultant Routes (Require Consultant Authentication)
- **`/`** - Main consultant dashboard
- **`/consultant/send-prompt`** - Send AI interaction prompt to readers
- **`/consultant/help-requests`** - Help request management
- **`/consultant/feedback`** - Feedback management
- **`/consultant/readers`** - Reader management
- **`/consultant/reports`** - Analytics and reports
- **`/consultant/reading-insights`** - Reader activity insights
- **`/consultant/assign-readers`** - Reader assignment interface

### Page Components
- `ConsultantLoginPage.tsx` - Consultant authentication
- `ConsultantDashboard.tsx` - Main dashboard
- `SendPromptPage.tsx` - Send prompts to readers
- `HelpRequestsPage.tsx` - Manage help requests
- `FeedbackManagementPage.tsx` - Manage user feedback
- `ReaderManagementPage.tsx` - Manage readers
- `AnalyticsReportsPage.tsx` - Analytics and reports
- `ReaderActivityInsightsPage.tsx` - Reader activity insights
- `AssignReadersPage.tsx` - Assign readers to consultants

---

## API Service Calls

### Authentication Services (`authService.ts`, `backendService.ts`)

#### Sign In
- **Method:** `signIn(email, password)`
- **Endpoint:** `/auth/v1/token` (POST)
- **Returns:** Session with access_token, user data
- **Database:** Updates `profiles` table, creates session

#### Sign Up
- **Method:** `signUp(email, password, firstName, lastName)`
- **Endpoint:** `/auth/v1/signup` (POST)
- **Returns:** User object
- **Database:** Creates entry in `profiles` table

#### Sign Out
- **Method:** `signOut()`
- **Endpoint:** `/auth/v1/logout` (POST)
- **Returns:** Success status
- **Database:** Clears session

#### Get Session
- **Method:** `getSession()`
- **Endpoint:** `/auth/v1/user` (GET)
- **Returns:** Current session and user data
- **Database:** Queries `profiles` table

#### Verify Book Code
- **Method:** `verifyBookCode(code, userId)`
- **Endpoint:** `/rest/v1/verification_codes` (GET with filter)
- **Returns:** Verification status, book access
- **Database:** Queries `verification_codes` table, updates `profiles.book_verified`

### Book Services (`bookService.ts`, `backendService.ts`)

#### Get Book Details
- **Method:** `getBook(bookId)`
- **Endpoint:** `/rest/v1/books` (GET)
- **Query:** `?id=eq.{bookId}`
- **Returns:** Book metadata (title, author, chapters, total_pages)

#### Get Page Content
- **Method:** `getPage(bookId, pageNumber)`
- **Endpoint:** `/rest/v1/rpc/get_sections_for_page` (RPC)
- **Parameters:** `{ book_id, page_number }`
- **Returns:** Page content with sections

#### Get Section Details
- **Method:** `getSectionDetails(sectionId)`
- **Endpoint:** `/rest/v1/sections` (GET)
- **Query:** `?id=eq.{sectionId}&select=*,chapters(*)`
- **Returns:** Section content with chapter info

#### Get Reading Progress
- **Method:** `getReadingProgress(userId, bookId)`
- **Endpoint:** `/rest/v1/reading_progress` (GET)
- **Query:** `?user_id=eq.{userId}&book_id=eq.{bookId}`
- **Returns:** Current reading position, last section

#### Save Reading Progress
- **Method:** `saveReadingProgress(userId, bookId, sectionId, position)`
- **Endpoint:** `/rest/v1/reading_progress` (POST/UPSERT)
- **Returns:** Success status
- **Database:** Upserts `reading_progress` table

#### Get Reading Statistics
- **Method:** `getReadingStats(userId, bookId)`
- **Endpoint:** `/rest/v1/reading_stats` (GET)
- **Query:** `?user_id=eq.{userId}&book_id=eq.{bookId}`
- **Returns:** Total reading time, pages read, progress percentage

#### Update Reading Statistics
- **Method:** `updateReadingStats(userId, bookId, currentPosition)`
- **Endpoint:** `/rest/v1/reading_stats` (POST/UPSERT)
- **Returns:** Updated stats
- **Database:** Upserts `reading_stats` table

### Dictionary Services (`dictionaryService.ts`, `glossaryService.ts`)

#### Get Definition
- **Method:** `getDefinition(term, bookId, sectionId, context)`
- **Endpoint:** `/rest/v1/rpc/get_definition_with_context` (RPC)
- **Parameters:** `{ term, book_id, section_id, context }`
- **Returns:** Dictionary entry with definition, examples, related terms
- **Fallback:** Local glossary cache, external dictionary API

#### Get Glossary Terms
- **Method:** `getAllGlossaryTerms(bookId)`
- **Endpoint:** `/rest/v1/alice_glossary` (GET)
- **Query:** `?book_id=eq.{bookId}`
- **Returns:** List of glossary terms for the book

### Help Request Services (`backendService.ts`, `consultantService.ts`)

#### Submit Help Request
- **Method:** `submitHelpRequest(userId, bookId, content, sectionId, context)`
- **Endpoint:** `/rest/v1/help_requests` (POST)
- **Body:** `{ user_id, book_id, content, section_id, context, status: 'PENDING' }`
- **Returns:** Created help request
- **Database:** Inserts into `help_requests` table

#### Get Help Requests
- **Method:** `getHelpRequests(statusFilter, userId, consultantId)`
- **Endpoint:** `/rest/v1/help_requests` (GET)
- **Query:** `?status=eq.{status}&user_id=eq.{userId}&consultant_id=eq.{consultantId}`
- **Returns:** List of help requests with user and book info

#### Update Help Request Status
- **Method:** `updateHelpRequestStatus(requestId, status, consultantId)`
- **Endpoint:** `/rest/v1/help_requests` (PATCH)
- **Query:** `?id=eq.{requestId}`
- **Body:** `{ status, consultant_id, updated_at }`
- **Returns:** Updated help request

### Consultant Services (`consultantService.ts`)

#### Get Verified Readers
- **Method:** `getVerifiedReaders()`
- **Endpoint:** `/rest/v1/profiles` (GET)
- **Query:** `?book_verified=eq.true&select=*,books(*)`
- **Returns:** List of verified reader profiles

#### Get Reader Interactions
- **Method:** `getReaderInteractions(userId, eventType)`
- **Endpoint:** `/rest/v1/interactions` (GET)
- **Query:** `?user_id=eq.{userId}&event_type=eq.{eventType}&order=created_at.desc&limit=50`
- **Returns:** Recent activity/interactions

#### Get Currently Logged In Users
- **Method:** `getCurrentlyLoggedInUsers()`
- **Endpoint:** `/rest/v1/interactions` (GET)
- **Query:** `?event_type=eq.LOGIN&select=*,profiles:user_id(*)&order=created_at.desc`
- **Returns:** Users who logged in recently

#### Get Reader Details
- **Method:** `getUserReadingDetails(userId, bookId)`
- **Endpoint:** Multiple queries:
  - `/rest/v1/profiles` - User profile
  - `/rest/v1/reading_progress` - Reading progress
  - `/rest/v1/reading_stats` - Reading statistics
  - `/rest/v1/interactions` - Activity history
- **Returns:** Complete reader profile with progress and activity

#### Get Dashboard Data
- **Method:** `getDashboardData(consultantId)`
- **Endpoint:** Multiple queries:
  - `/rest/v1/help_requests` - Pending requests
  - `/rest/v1/interactions` - Recent activity
  - `/rest/v1/profiles` - Reader count
- **Returns:** Dashboard summary data

### Activity Tracking Services (`activityTrackingService.ts`)

#### Track Event
- **Method:** `trackEvent(event)`
- **Endpoint:** `/rest/v1/interactions` (POST)
- **Body:** `{ user_id, event_type, book_id, section_id, page_number, content, context }`
- **Event Types:** `LOGIN`, `LOGOUT`, `PAGE_SYNC`, `SECTION_SYNC`, `DEFINITION_LOOKUP`, `AI_QUERY`, `HELP_REQUEST`, `FEEDBACK_SUBMISSION`
- **Database:** Inserts into `interactions` table

#### Track Login
- **Method:** `trackLogin(userId)`
- **Calls:** `trackEvent({ event_type: 'LOGIN' })`

#### Track Logout
- **Method:** `trackLogout(userId)`
- **Calls:** `trackEvent({ event_type: 'LOGOUT' })`

#### Track Page Sync
- **Method:** `trackPageSync(userId, bookId, pageNumber)`
- **Calls:** `trackEvent({ event_type: 'PAGE_SYNC', page_number })`

---

## Database Queries

### Tables Used

1. **`profiles`** - User profiles (readers and consultants)
   - Fields: `id`, `email`, `first_name`, `last_name`, `book_verified`, `is_consultant`, `is_verified`, `created_at`, `updated_at`
   - Queries: SELECT, INSERT, UPDATE

2. **`books`** - Book metadata
   - Fields: `id`, `title`, `author`, `total_pages`, `description`, `cover_image`
   - Queries: SELECT

3. **`chapters`** - Chapter information
   - Fields: `id`, `book_id`, `number`, `title`, `start_page`
   - Queries: SELECT

4. **`sections`** - Section content (subdivisions of chapters)
   - Fields: `id`, `chapter_id`, `book_id`, `title`, `content`, `start_page`, `word_count`
   - Queries: SELECT

5. **`pages`** - Page-level content (physical book pages)
   - Fields: `id`, `book_id`, `page_number`, `chapter_id`, `content`, `word_count`
   - Queries: SELECT

6. **`verification_codes`** - Book access codes
   - Fields: `id`, `code`, `book_id`, `is_used`, `used_by`, `created_at`
   - Queries: SELECT, UPDATE

7. **`reading_progress`** - User reading progress
   - Fields: `id`, `user_id`, `book_id`, `section_id`, `last_position`, `last_read_at`
   - Queries: SELECT, INSERT, UPDATE (UPSERT)

8. **`reading_stats`** - Reading statistics
   - Fields: `id`, `user_id`, `book_id`, `total_reading_time`, `pages_read`, `last_updated`
   - Queries: SELECT, INSERT, UPDATE (UPSERT)

9. **`help_requests`** - Help request queue
   - Fields: `id`, `user_id`, `book_id`, `section_id`, `content`, `status`, `consultant_id`, `created_at`, `updated_at`
   - Queries: SELECT, INSERT, UPDATE

10. **`interactions`** - Activity tracking
    - Fields: `id`, `user_id`, `event_type`, `book_id`, `section_id`, `page_number`, `content`, `context`, `created_at`
    - Queries: SELECT, INSERT

11. **`alice_glossary`** - Glossary terms
    - Fields: `id`, `term`, `definition`, `book_id`, `section_id`, `examples`, `related_terms`
    - Queries: SELECT

12. **`user_feedback`** - User feedback
    - Fields: `id`, `user_id`, `book_id`, `section_id`, `feedback_type`, `content`, `is_public`, `created_at`
    - Queries: SELECT, INSERT

13. **`consultant_assignments`** - Reader-consultant assignments
    - Fields: `id`, `consultant_id`, `reader_id`, `book_id`, `assigned_at`
    - Queries: SELECT, INSERT

### Common Query Patterns

#### Supabase Query Syntax
```typescript
// Basic SELECT
.from('table_name')
.select('*')
.eq('column', 'value')

// SELECT with joins
.from('interactions')
.select('*, profiles:user_id(first_name, last_name, email)')

// SELECT with ordering and limits
.from('interactions')
.select('*')
.order('created_at', { ascending: false })
.limit(50)

// INSERT
.from('table_name')
.insert({ field1: 'value1', field2: 'value2' })
.select()

// UPDATE
.from('table_name')
.update({ field: 'new_value' })
.eq('id', id)
.select()

// RPC (Remote Procedure Call)
.rpc('function_name', { param1: 'value1' })
```

#### Foreign Key Relationships
- `profiles.id` → `reading_progress.user_id`
- `profiles.id` → `help_requests.user_id`
- `profiles.id` → `interactions.user_id`
- `books.id` → `chapters.book_id`
- `books.id` → `sections.book_id`
- `books.id` → `pages.book_id`
- `chapters.id` → `sections.chapter_id`
- `sections.id` → `reading_progress.section_id`

---

## Authentication Flow

### Reader Authentication Flow

1. **Landing Page** (`/`)
   - User sees welcome page
   - Options: Login or Register

2. **Registration** (`/register`)
   - User enters: email, password, first name, last name
   - Calls `signUp(email, password, firstName, lastName)`
   - Creates profile in `profiles` table
   - Redirects to `/verify` (book verification required)

3. **Login** (`/login`)
   - User enters: email, password
   - Calls `signIn(email, password)`
   - Validates credentials via `/auth/v1/token`
   - Creates session, stores access_token
   - Checks `book_verified` status
   - If verified: redirects to `/reader`
   - If not verified: redirects to `/verify`

4. **Book Verification** (`/verify`)
   - User enters verification code
   - Calls `verifyBookCode(code, userId)`
   - Queries `verification_codes` table
   - Updates `profiles.book_verified = true`
   - Redirects to `/reader`

5. **Reader Dashboard** (`/reader`)
   - Requires: authenticated + verified
   - Shows book selection, reading progress
   - Can navigate to reading interface

6. **Session Management**
   - Session stored in browser (localStorage/sessionStorage)
   - Access token used for API requests (Bearer token)
   - Session validated on each protected route
   - Auto-logout on token expiry

### Consultant Authentication Flow

1. **Consultant Login** (`/consultant/login`)
   - Consultant enters: email, password
   - Calls `signIn(email, password)`
   - Validates credentials
   - Checks `profiles.is_consultant = true`
   - Creates session
   - Redirects to `/` (consultant dashboard)

2. **Protected Routes**
   - All consultant routes wrapped in `ConsultantProtectedRoute`
   - Validates consultant role
   - Redirects to login if not authenticated

3. **Session Management**
   - Same as reader authentication
   - Separate session context (`ConsultantAuthContext`)

---

## Real-time Features

### Current Implementation

#### WebSocket/Socket.IO (Currently Disabled by Default)
- **Reader App:** Socket.IO client (disabled unless `VITE_ENABLE_REALTIME=true`)
- **Consultant Dashboard:** Socket.IO client for real-time updates
- **Server:** WebSocket server on port 3001 (optional)

#### Activity Tracking (Database-Based)
- **Primary Method:** Database polling and inserts
- **Table:** `interactions` table
- **Events Tracked:**
  - `LOGIN` - User login events
  - `LOGOUT` - User logout events
  - `PAGE_SYNC` - Page navigation
  - `SECTION_SYNC` - Section navigation
  - `DEFINITION_LOOKUP` - Dictionary lookups
  - `AI_QUERY` - AI assistance requests
  - `HELP_REQUEST` - Help request submissions
  - `FEEDBACK_SUBMISSION` - Feedback submissions

#### Real-time Updates in Consultant Dashboard
- **Method:** Polling `interactions` table
- **Frequency:** Auto-refresh every 30 seconds (configurable)
- **Data:** Recent reader activity, online users, help requests

#### Online Users Detection
- **Method:** Query `interactions` table for recent `LOGIN` events
- **Query:** `SELECT * FROM interactions WHERE event_type = 'LOGIN' AND created_at > NOW() - INTERVAL '5 minutes'`
- **Display:** Shows currently logged-in readers

### Future Migration (Go Implementation)

**Recommended Approach:**
- **Server-Sent Events (SSE)** for one-way updates (simpler)
- **WebSocket** for bidirectional communication (if needed)
- **Go Native:** Use `golang.org/x/net/websocket` or standard library

**Events to Support:**
- Reader login/logout notifications
- Help request updates
- Reading progress updates
- Real-time activity feed

---

## Key Components

### Reader App Components

#### Authentication Components
- `AuthContext.tsx` - Authentication state management
- `RouteGuard.tsx` - Route protection based on auth status
- `ProtectedRoute.tsx` - Protected route wrapper
- `LoginPage.tsx` - Login form
- `RegisterPage.tsx` - Registration form
- `VerifyPage.tsx` - Book verification form

#### Reading Components
- `ReaderDashboard.tsx` - Main reader dashboard
- `ReaderPage.tsx` - Individual page view
- `MainInteractionPage.tsx` - Main reading interface
- `ReaderStatistics.tsx` - Reading statistics display
- `DictionaryDialog.tsx` - Dictionary lookup dialog
- `AIAssistanceDialog.tsx` - AI assistance interface
- `HelpRequestDialog.tsx` - Help request submission

#### Navigation Components
- `EnhancedAppHeader.tsx` - App header with navigation
- `NavigationListener.tsx` - Tracks route changes

#### UI Components
- `ErrorBoundary.tsx` - Error handling
- `SnackbarProvider.tsx` - Notification system
- `AccessibilityMenu.tsx` - Accessibility options
- `LogViewer.tsx` - Debug log viewer

### Consultant Dashboard Components

#### Authentication Components
- `ConsultantAuthContext.tsx` - Consultant authentication state
- `ConsultantProtectedRoute.tsx` - Protected route wrapper
- `ConsultantLoginPage.tsx` - Consultant login form

#### Dashboard Components
- `ConsultantDashboard.tsx` - Main dashboard
- `ReaderActivityDashboard.tsx` - Reader activity display
- `OnlineReadersWidget.tsx` - Online users widget
- `HelpRequestsPage.tsx` - Help request management
- `ReaderManagementPage.tsx` - Reader list and details

#### Analytics Components
- `AnalyticsReportsPage.tsx` - Analytics and reports
- `ReaderActivityInsightsPage.tsx` - Activity insights

---

## Feature Summary

### Core Features

#### Reader Features
1. **User Authentication**
   - Registration with email/password
   - Login/logout
   - Session management
   - Book verification system

2. **Reading Interface**
   - Page-by-page navigation
   - Section-based reading
   - Reading progress tracking
   - Last position bookmarking

3. **Dictionary/Glossary**
   - Word lookup from text
   - Context-aware definitions
   - Alice-specific glossary terms
   - Related terms suggestions

4. **AI Assistance**
   - Ask questions about passages
   - Context-aware responses
   - Integration with AI service (Moonshot/Kimi)

5. **Help Requests**
   - Submit help requests
   - Track request status
   - View consultant responses

6. **Reading Statistics**
   - Total reading time
   - Pages read
   - Progress percentage
   - Reading history

7. **Activity Tracking**
   - Automatic activity logging
   - Page navigation tracking
   - Dictionary lookup tracking
   - AI query tracking

#### Consultant Features
1. **Consultant Authentication**
   - Separate consultant login
   - Role-based access control

2. **Reader Monitoring**
   - View all verified readers
   - See reader activity in real-time
   - View online users
   - Reader profile details

3. **Help Request Management**
   - View pending requests
   - Assign requests to consultants
   - Update request status
   - Respond to requests

4. **Reader Assignment**
   - Assign readers to consultants
   - View assignment history
   - Manage assignments

5. **Analytics & Reports**
   - Reader engagement metrics
   - Help request trends
   - Reading activity insights
   - Dashboard summaries

6. **Feedback Management**
   - View user feedback
   - Public/private feedback
   - Feedback filtering

7. **AI Prompt Sending**
   - Send prompts to readers
   - Trigger AI interactions
   - Monitor responses

### Technical Features

1. **Service Registry Pattern**
   - Centralized service management
   - Lazy initialization
   - Dependency injection

2. **Error Handling**
   - Comprehensive error boundaries
   - Retry logic for API calls
   - Fallback to mock data

3. **Caching**
   - Dictionary cache
   - Service response cache
   - Local storage caching

4. **Mock Backend**
   - Fallback when backend unavailable
   - Development testing
   - Demo mode

5. **Accessibility**
   - ARIA labels
   - Keyboard navigation
   - Screen reader support
   - Accessibility menu

6. **Performance Optimization**
   - Lazy loading
   - Code splitting
   - Memoization
   - Debouncing

---

## Migration Notes

### Key Considerations for Go Migration

1. **API Compatibility**
   - Maintain Supabase-compatible API responses
   - Support same query parameters (`select`, `eq`, `order`, `limit`)
   - Handle join syntax (`profiles:user_id(...)`)
   - Support RPC functions

2. **Authentication**
   - JWT token generation/validation
   - Session management
   - Role-based access control (reader vs consultant)

3. **Real-time Updates**
   - Replace Socket.IO with Go WebSocket or SSE
   - Maintain same event types
   - Support same data structures

4. **Database Schema**
   - Use existing SQLite schema
   - Maintain foreign key relationships
   - Support all table operations

5. **Frontend Options**
   - Option A: Go templates + HTMX (recommended)
   - Option B: Embedded React static build
   - Option C: Hybrid approach

---

## Next Steps

1. **Review this inventory** with stakeholders
2. **Prioritize features** for initial migration
3. **Design Go API** matching Supabase patterns
4. **Implement authentication** system
5. **Migrate core features** (reading, dictionary, help requests)
6. **Add consultant features** (dashboard, monitoring)
7. **Implement real-time** updates
8. **Test thoroughly** before removing React apps

---

**End of Feature Inventory**


