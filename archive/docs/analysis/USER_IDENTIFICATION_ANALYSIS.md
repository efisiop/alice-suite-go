# User Identification Analysis - Same IP Address

**Date:** 2025-12-03  
**Question:** Can two readers logging in from the same IP address cause identification conflicts?

---

## ✅ Answer: No Conflicts

Users are **uniquely identified** by their `user_id` and JWT token, **not** by IP address. Each user from the same IP will be correctly identified and tracked separately.

---

## How User Identification Works

### 1. Session Management ✅

**Identification Method:** `token_hash` (unique JWT token hash)

- Each user receives a **unique JWT token** on login
- Token is hashed (SHA-256) and stored as `token_hash` in `sessions` table
- Sessions are identified by `token_hash`, not IP address
- Each user from the same IP has a **different token** and **different session**

**Code Reference:**
```go
// internal/database/sessions.go
func CreateSession(userID, token, ipAddress, userAgent string, expiresIn time.Duration)
// Creates session with unique token_hash per user
```

### 2. Activity Tracking ✅

**Identification Method:** `user_id` + `session_id`

- All activities are logged with `user_id` (from authenticated token)
- Each activity includes `session_id` linking to the specific session
- IP address is stored in `metadata` JSON field for audit purposes only
- **Not used for identification**

**Code Reference:**
```go
// internal/database/activity.go
func LogActivity(activity *ActivityLog) error
// Logs with user_id and session_id, not IP
```

### 3. Real-Time Updates (SSE) ✅

**Identification Method:** `user_id`

- SSE connections are registered by `user_id`
- Each user gets their own SSE connection
- Events are broadcast with `user_id` included
- **IP address not used**

**Code Reference:**
```go
// internal/realtime/broadcaster.go
func (b *Broadcaster) RegisterClient(id, role string) *Client
// Registers by user_id, not IP
```

### 4. Consultant Dashboard ✅

**Identification Method:** `user_id`

- All queries use `user_id` to identify users
- Dashboard displays users by:
  - `user_id`
  - `first_name` + `last_name`
  - `email`
- **IP address not used for identification**

**Code Reference:**
```go
// internal/database/consultant.go
func GetActiveReaders(minutesThreshold int) ([]*ActiveReader, error)
// Queries by user_id, not IP
```

---

## ⚠️ One Consideration: Rate Limiting

### Current Implementation: Per-IP Rate Limiting

**Location:** `internal/middleware/rate_limit.go`

**How it works:**
- Rate limiting is applied **per IP address**
- Two users from the same IP share the same rate limit bucket
- Limits: 10 requests/second, burst of 20 requests

**Potential Issue:**
If User A makes many rapid requests and hits the rate limit, User B (from same IP) might also be rate-limited temporarily.

**Impact:** 
- **Low** - Rate limits are generous (10 req/sec)
- Only affects rapid-fire requests
- Normal usage won't trigger this

### Solution Options:

#### Option 1: Keep Per-IP (Current)
- **Pros:** Simple, prevents abuse from single IP
- **Cons:** Shared IP users share limits
- **Recommendation:** Keep this for now (low impact)

#### Option 2: Per-User Rate Limiting
- **Pros:** Each user has independent limits
- **Cons:** Requires extracting `user_id` from token (adds overhead)
- **Implementation:** Would need to modify `RateLimit` middleware to check authenticated user

---

## Verification Queries

### Check Sessions from Same IP
```sql
-- See all sessions from a specific IP
SELECT user_id, email, created_at, expires_at 
FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.ip_address = '192.168.1.100'
ORDER BY created_at DESC;
```

### Check Activities from Same IP
```sql
-- See activities from users on same IP
SELECT al.user_id, u.email, al.activity_type, al.created_at
FROM activity_logs al
JOIN users u ON al.user_id = u.id
JOIN sessions s ON al.session_id = s.id
WHERE s.ip_address = '192.168.1.100'
ORDER BY al.created_at DESC;
```

---

## Summary

| Component | Identification Method | IP Used? | Conflict Risk |
|-----------|---------------------|----------|---------------|
| **Sessions** | `token_hash` (unique per user) | No | ✅ None |
| **Activities** | `user_id` + `session_id` | No | ✅ None |
| **SSE Connections** | `user_id` | No | ✅ None |
| **Dashboard Queries** | `user_id` | No | ✅ None |
| **Rate Limiting** | IP address | Yes | ⚠️ Shared limits |

---

## Conclusion

✅ **No identification conflicts** - Each user is uniquely identified by their `user_id` and token, regardless of IP address.

⚠️ **Minor consideration** - Rate limiting is shared per IP, but this is unlikely to cause issues in normal usage.

**Recommendation:** Current implementation is correct and safe. No changes needed unless you want per-user rate limiting.

