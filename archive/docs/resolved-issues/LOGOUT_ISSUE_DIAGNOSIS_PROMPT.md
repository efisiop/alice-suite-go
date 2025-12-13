# COMPREHENSIVE PROMPT: Real-Time Logout Issue Diagnosis and Fix

## Problem Statement

I have a consultant dashboard application built with Go backend and JavaScript frontend using Server-Sent Events (SSE) for real-time updates. The login functionality works PERFECTLY in real-time, but logout does NOT work as expected.

### What Works (Login) ✅
1. When a reader logs in, the "Logged In" count increases immediately
2. The reader's window/card appears immediately in the dashboard
3. Everything updates in real-time via SSE

### What Doesn't Work (Logout) ❌
1. When a reader logs out, the "Logged In" count SHOULD decrease immediately but does not
2. The reader's window/card SHOULD disappear immediately but does not
3. Real-time updates are not happening on logout

## System Architecture

**Technology Stack:**
- Backend: Go (Golang)
- Frontend: JavaScript (vanilla, no frameworks)
- Real-time: Server-Sent Events (SSE)
- Database: SQLite

**Key Files:**
- `internal/handlers/auth.go` - Handles login/logout API endpoints
- `internal/handlers/sse.go` - Manages SSE connections and event broadcasting
- `internal/realtime/broadcaster.go` - Event broadcasting system
- `internal/templates/consultant/dashboard.html` - Frontend dashboard with JavaScript

## Current Implementation Details

### Login Flow (WORKING)
1. Reader POSTs to `/auth/v1/login`
2. Server creates session, logs activity, and calls `BroadcastLogin(userID, email, firstName, lastName)`
3. SSE event is sent to consultant dashboard with event type `"login"`
4. Frontend `eventSource.addEventListener('login', ...)` handler:
   - Immediately increments "Logged In" count
   - Creates reader card using `addActivityToReaderCard()` with LOGIN activity
   - Card appears instantly

**Login Event Format (Backend):**
```go
// internal/handlers/sse.go
func BroadcastLogin(userID, email, firstName, lastName string) {
    broadcaster := realtime.GetBroadcaster()
    event := realtime.CreateEvent(realtime.EventTypeLogin, map[string]interface{}{
        "user_id":    userID,
        "email":      email,
        "first_name": firstName,
        "last_name":  lastName,
    })
    broadcaster.BroadcastToRole(event, "consultant")
}
```

**Login Event Handler (Frontend) - WORKING:**
```javascript
eventSource.addEventListener('login', function(e) {
    const eventData = JSON.parse(e.data);
    const loginData = eventData.data || eventData;
    
    // Update count immediately
    const currentCount = parseInt(document.getElementById('logged-in-readers').textContent) || 0;
    document.getElementById('logged-in-readers').textContent = currentCount + 1;
    
    // Create card immediately
    const userId = loginData.user_id || loginData.userId || loginData.id;
    if (userId) {
        const loginActivity = {
            id: 'login-' + Date.now(),
            user_id: userId,
            first_name: loginData.first_name || loginData.firstName,
            last_name: loginData.last_name || loginData.lastName,
            email: loginData.email,
            event_type: 'LOGIN',
            activity_type: 'LOGIN',
            created_at: new Date().toISOString()
        };
        addActivityToReaderCard(loginActivity);  // This creates the card
    }
});
```

### Logout Flow (NOT WORKING)
1. Reader POSTs to `/auth/v1/logout`
2. Server logs activity, calls `BroadcastLogout(userID)`, deletes sessions
3. SSE event is sent with event type `"logout"`
4. Frontend handler should:
   - Immediately decrement "Logged In" count
   - Remove reader card/window from dashboard
   - BUT THIS IS NOT HAPPENING

**Logout Event Format (Backend):**
```go
// internal/handlers/sse.go
func BroadcastLogout(userID string) {
    broadcaster := realtime.GetBroadcaster()
    event := realtime.CreateEvent(realtime.EventTypeLogout, map[string]interface{}{
        "user_id": userID,
    })
    broadcaster.BroadcastToRole(event, "consultant")
}
```

**Logout Event Handler (Frontend) - NOT WORKING:**
```javascript
eventSource.addEventListener('logout', function(e) {
    const eventData = JSON.parse(e.data);
    const logoutData = eventData.data || eventData;
    const userId = logoutData.user_id || logoutData.userId || logoutData.id || eventData.user_id;
    
    // This SHOULD work but doesn't
    const currentCount = parseInt(document.getElementById('logged-in-readers').textContent) || 0;
    document.getElementById('logged-in-readers').textContent = Math.max(0, currentCount - 1);
    
    // Multiple methods tried to remove card, but none work
    const readerCard = readerCards.get(userId);
    if (readerCard && readerCard.cardElement) {
        readerCard.cardElement.remove();
        readerCards.delete(userId);
    }
    // ... more removal attempts ...
});
```

## Data Structures

**Reader Cards Storage:**
- Cards are stored in a JavaScript Map: `const readerCards = new Map()`
- Key: `user_id` (string)
- Value: Object with `{cardElement, activities, activitiesList, unreadCount, ...}`

**Card DOM Structure:**
- Card element ID: `reader-card-${userId}`
- Card has attribute: `data-user-id="${userId}"`
- Card class: `card reader-card`

## Questions to Investigate

1. **Is the logout event being received?** Check if `eventSource.addEventListener('logout', ...)` is actually being triggered
2. **Is the event data format correct?** Compare logout event structure with login event structure
3. **Is the user_id matching?** The user_id used to create the card during login might be different from the user_id in the logout event
4. **Timing issue?** Could the card be created/removed in a way that breaks the reference?
5. **SSE connection issue?** Is the logout event being broadcast correctly to all consultant connections?

## Expected Behavior

When a reader logs out:
1. ✅ "Logged In" count decreases by 1 immediately (e.g., from 3 to 2)
2. ✅ Reader's card/window disappears from dashboard immediately
3. ✅ No delay - instant real-time update just like login

## Current Behavior

When a reader logs out:
1. ❌ Count does NOT decrease (stays the same)
2. ❌ Reader's card/window does NOT disappear
3. ❌ Must manually refresh page to see changes

## Debugging Information Needed

Please help me:
1. Add comprehensive logging to trace the logout event from backend to frontend
2. Compare login event flow with logout event flow to find discrepancies
3. Verify the SSE event format for logout matches what frontend expects
4. Check if readerCards Map is properly maintained and accessible
5. Ensure the user_id format/type is consistent between login and logout
6. Test if the logout event is even reaching the frontend

## Code Locations

- **Backend Logout Handler:** `internal/handlers/auth.go` lines 208-280
- **Backend Logout Broadcast:** `internal/handlers/sse.go` lines 145-152
- **Frontend Logout Handler:** `internal/templates/consultant/dashboard.html` lines 671-760
- **Frontend Login Handler:** `internal/templates/consultant/dashboard.html` lines 630-669
- **Card Creation Function:** `internal/templates/consultant/dashboard.html` lines 938-1060
- **SSE Connection Setup:** `internal/templates/consultant/dashboard.html` lines 555-600

## Task

Please:
1. **Diagnose** why logout events are not working in real-time
2. **Compare** login vs logout implementation to find differences
3. **Fix** the logout handler to work exactly like login handler
4. **Verify** the fix ensures real-time count updates and card removal
5. **Test** that logout events are properly received and processed

The goal is to make logout work EXACTLY like login works - instant, real-time updates with no page refresh needed.

