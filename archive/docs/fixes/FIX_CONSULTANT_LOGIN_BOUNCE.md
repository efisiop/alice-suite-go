# Fix: Consultant Login Bounce-Back Issue

## Problem
When logging into the consultant portal, users were being "bounced back" to the login page immediately after successful login. This happened in a fraction of a second.

## Root Cause
The login process stored the authentication token only in `localStorage`, but when the browser navigated to `/consultant` after login, it was a regular GET request that didn't include the token in the Authorization header. The server-side handler checked for the token in:
1. Authorization header (not present in browser navigation)
2. Cookie (not set during login)

Since neither was present, the server redirected back to login before the dashboard could load.

## Solution
Modified the consultant login JavaScript to set **both** localStorage (for JavaScript API calls) **and** a cookie (for server-side page navigation):

1. **Login (`internal/templates/consultant/login.html`)**:
   - Sets token in localStorage: `localStorage.setItem('auth_token', token)`
   - Sets token in cookie: `document.cookie = 'auth_token=${token}; expires=${expires}; path=/; SameSite=Lax'`
   - Added 100ms delay before redirect to ensure cookie is set

2. **Logout (`internal/templates/consultant/dashboard.html`)**:
   - Clears localStorage: `localStorage.removeItem('auth_token')`
   - Clears cookie: `document.cookie = 'auth_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; SameSite=Lax'`

## How It Works Now

1. User enters credentials and clicks Login
2. Server validates credentials and returns JWT token
3. JavaScript stores token in **both** localStorage and cookie
4. Browser redirects to `/consultant`
5. Server reads token from cookie
6. Server validates token and serves dashboard
7. Dashboard JavaScript uses token from localStorage for API calls

## Testing
After this fix, the consultant login should:
- ✅ Successfully log in without bouncing back
- ✅ Load the dashboard immediately after login
- ✅ Maintain authentication across page refreshes
- ✅ Properly log out and clear both storage locations

## Files Changed
- `internal/templates/consultant/login.html` - Added cookie setting on login
- `internal/templates/consultant/dashboard.html` - Added cookie clearing on logout

