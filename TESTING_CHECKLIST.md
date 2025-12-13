# Testing Checklist

**Date:** 2025-01-23  
**Purpose:** Comprehensive testing checklist for Alice Suite Go migration

---

## Pre-Deployment Testing

### ✅ Build & Compilation
- [x] Project compiles without errors
- [x] All dependencies resolved
- [x] Binary created successfully
- [ ] Binary runs without errors
- [ ] Database initializes correctly

### ✅ Authentication System
- [ ] User registration works
- [ ] User login works
- [ ] JWT token generation works
- [ ] Token validation works
- [ ] Logout works
- [ ] Session management works
- [ ] Book verification works
- [ ] Role-based access control works

### ✅ REST API Endpoints
- [ ] GET `/rest/v1/books` - Returns books list
- [ ] GET `/rest/v1/books?id=eq.{id}` - Returns specific book
- [ ] GET `/rest/v1/chapters?book_id=eq.{id}` - Returns chapters
- [ ] GET `/rest/v1/sections?chapter_id=eq.{id}` - Returns sections
- [ ] POST `/rest/v1/help_requests` - Creates help request
- [ ] GET `/rest/v1/help_requests?status=eq.pending` - Filters help requests
- [ ] PATCH `/rest/v1/help_requests?id=eq.{id}` - Updates help request
- [ ] GET `/rest/v1/interactions` - Returns interactions
- [ ] POST `/rest/v1/interactions` - Creates interaction
- [ ] GET `/rest/v1/profiles` - Returns profiles
- [ ] GET `/rest/v1/reading_progress` - Returns reading progress
- [ ] POST `/rest/v1/reading_progress` - Saves reading progress

### ✅ RPC Functions
- [ ] POST `/rest/v1/rpc/get_definition_with_context` - Dictionary lookup
- [ ] POST `/rest/v1/rpc/get_sections_for_page` - Get page sections
- [ ] POST `/rest/v1/rpc/verify-book-code` - Verify book code
- [ ] GET `/rest/v1/rpc/check-book-verified` - Check verification status
- [ ] POST `/rest/v1/rpc/check_table_exists` - Check table exists

### ✅ Reader App Pages
- [ ] `/` - Landing page loads
- [ ] `/login` - Login page loads and works
- [ ] `/register` - Registration page loads and works
- [ ] `/verify` - Verification page loads and works
- [ ] `/reader` - Dashboard loads (requires auth + verification)
- [ ] `/reader/interaction` - Reading interface loads
- [ ] `/reader/statistics` - Statistics page loads

### ✅ Consultant Dashboard Pages
- [ ] `/consultant/login` - Consultant login works
- [ ] `/consultant` - Dashboard loads (requires consultant auth)
- [ ] `/consultant/help-requests` - Help requests page loads
- [ ] `/consultant/readers` - Readers page loads

### ✅ Real-time Features
- [ ] SSE connection establishes (`/api/realtime/events`)
- [ ] Events broadcast correctly
- [ ] Help request events received by consultants
- [ ] Login/logout events broadcast
- [ ] Activity events broadcast
- [ ] WebSocket connection works (optional)

### ✅ Activity Tracking
- [ ] POST `/api/activity/track` - Tracks activity
- [ ] Activities saved to database
- [ ] Activities broadcast to consultants
- [ ] Event types: LOGIN, LOGOUT, PAGE_SYNC, DEFINITION_LOOKUP, HELP_REQUEST

### ✅ Database Operations
- [ ] Foreign key constraints work
- [ ] UUID generation works
- [ ] Timestamp handling works
- [ ] Boolean conversion works
- [ ] Query parsing works (select, filters, order, limit, offset)
- [ ] Join syntax works

---

## Integration Testing

### Authentication Flow
1. Register new user → Verify redirects to verification
2. Verify book code → Verify redirects to dashboard
3. Login → Verify redirects to dashboard
4. Logout → Verify redirects to landing page

### Reading Flow
1. Navigate to reading interface
2. Load page content
3. Look up word in dictionary
4. Navigate to next page
5. Submit help request
6. Verify activity tracking

### Consultant Flow
1. Login as consultant
2. View dashboard
3. See real-time help requests
4. View online readers
5. Update help request status
6. Verify real-time updates

---

## Performance Testing

- [ ] Server starts quickly (< 2 seconds)
- [ ] Database queries execute quickly (< 100ms)
- [ ] SSE connections handle multiple clients
- [ ] Concurrent requests handled correctly
- [ ] Memory usage reasonable

---

## Security Testing

- [ ] SQL injection prevention (table name validation)
- [ ] XSS prevention (template escaping)
- [ ] CSRF protection (token validation)
- [ ] Authentication required for protected routes
- [ ] Role-based access control works
- [ ] Token expiration works

---

## Browser Compatibility

- [ ] Chrome (latest)
- [ ] Firefox (latest)
- [ ] Safari (latest)
- [ ] Edge (latest)
- [ ] Mobile browsers

---

## Deployment Checklist

- [ ] Database file exists and is accessible
- [ ] Environment variables set (if needed)
- [ ] Port configuration correct
- [ ] Static files served correctly
- [ ] Templates render correctly
- [ ] Logging works
- [ ] Error handling works

---

## Known Issues / TODO

- [ ] Add more comprehensive error handling
- [ ] Add request rate limiting
- [ ] Add CORS configuration for production
- [ ] Add HTTPS support
- [ ] Add database backup functionality
- [ ] Add monitoring/logging

---

**Status:** Ready for testing

