# Browser Compatibility Testing Checklist

Use this checklist to verify the application works consistently across different browsers.

## Test URL
**Production:** https://alice-suite-go.onrender.com/reader/login

## Browsers to Test
- ✅ Chrome (latest)
- ✅ Firefox (latest)
- ✅ Safari (latest)
- ✅ Edge (latest)
- ✅ Opera (optional)

---

## Test Scenarios

### 1. Login Flow
- [ ] **Chrome:** Login works, redirects correctly
- [ ] **Firefox:** Login works, redirects correctly
- [ ] **Safari:** Login works, redirects correctly
- [ ] **Edge:** Login works, redirects correctly

**Credentials to test:**
- Email: `efisio@efisio.com`
- Password: `efisio123`

**What to verify:**
- Login form submits correctly
- Token is stored in sessionStorage
- Cookie is set properly
- Redirect to dashboard works
- User name displays in navbar

---

### 2. Book Page Loading
- [ ] **Chrome:** Pages load correctly when clicking "Open Book"
- [ ] **Firefox:** Pages load correctly when clicking "Open Book"
- [ ] **Safari:** Pages load correctly when clicking "Open Book"
- [ ] **Edge:** Pages load correctly when clicking "Open Book"

**What to verify:**
- Click "Open Book" button on dashboard
- Page content displays correctly
- Sections are visible
- Navigation (Previous/Next page) works
- No console errors (check F12 → Console)

---

### 3. Page Navigation
- [ ] **Chrome:** Previous/Next page buttons work
- [ ] **Firefox:** Previous/Next page buttons work
- [ ] **Safari:** Previous/Next page buttons work
- [ ] **Edge:** Previous/Next page buttons work

**What to verify:**
- Use "Previous Page" button
- Use "Next Page" button
- Use "Go to Page" input field
- Page number updates correctly
- Content changes correctly

---

### 4. Real-Time Updates (SSE)
- [ ] **Chrome:** No SSE connection errors in console
- [ ] **Firefox:** No SSE connection errors in console
- [ ] **Safari:** No SSE connection errors in console
- [ ] **Edge:** No SSE connection errors in console

**What to verify:**
- Open browser console (F12)
- Look for SSE connection errors
- Should see: `[SSE] Connected successfully` (or similar)
- No "Failed to load resource" errors for SSE endpoints

---

### 5. Logout
- [ ] **Chrome:** Logout works, redirects to login
- [ ] **Firefox:** Logout works, redirects to login
- [ ] **Safari:** Logout works, redirects to login
- [ ] **Edge:** Logout works, redirects to login

**What to verify:**
- Click logout link
- SessionStorage is cleared
- Cookie is cleared
- Redirects to login page
- Cannot access protected pages after logout

---

## Known Browser-Specific Features

### Safari Compatibility
The code includes Safari-specific fixes:
- Cookie handling with encoding fallback
- Hostname normalization (localhost vs 127.0.0.1)
- Proxy bypass configuration (for local development)

### Modern Browser Features Used
- **EventSource (SSE):** Supported in all modern browsers (IE not supported)
- **Fetch API:** Supported in all modern browsers (IE11+ with polyfill)
- **sessionStorage:** Supported in all modern browsers
- **JSON.parse/stringify:** Universal support

---

## What to Check in Browser Console

### Good Signs ✅
- No red errors
- SSE connection messages
- Successful API calls (200 status)
- Token stored in sessionStorage

### Bad Signs ❌
- 404 errors (endpoints not found)
- 500 errors (server errors)
- SSE connection errors
- JSON parsing errors
- CORS errors

---

## Quick Test Script

1. Open browser
2. Open Developer Tools (F12)
3. Go to Console tab
4. Navigate to: https://alice-suite-go.onrender.com/reader/login
5. Login with credentials
6. Check console for errors
7. Click "Open Book"
8. Check console for errors
9. Navigate between pages
10. Check console for errors
11. Logout
12. Check console for errors

Repeat for each browser.

---

## Reporting Issues

If you find browser-specific issues, note:
- Browser name and version
- Operating system
- Steps to reproduce
- Console error messages (if any)
- Screenshot (if helpful)

---

## Browser Support Matrix

| Feature | Chrome | Firefox | Safari | Edge |
|---------|--------|---------|--------|------|
| Login | ✅ | ✅ | ✅ | ✅ |
| SSE (Real-time) | ✅ | ✅ | ✅ | ✅ |
| Fetch API | ✅ | ✅ | ✅ | ✅ |
| sessionStorage | ✅ | ✅ | ✅ | ✅ |
| Cookies | ✅ | ✅ | ✅* | ✅ |

*Safari has stricter cookie handling (code includes workarounds)
