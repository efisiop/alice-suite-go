# AI Assistant Troubleshooting Guide

If the AI assistant is not working, follow these steps to diagnose and fix the issue:

## ✅ Step 1: Verify Render Deployment Completed

1. Go to your Render dashboard: https://dashboard.render.com
2. Click on `alice-suite-go` service
3. Go to **Logs** tab
4. Check for:
   - ✅ Build completed successfully
   - ✅ Deployment completed
   - ✅ Server started without errors
   - ❌ Any errors mentioning "GEMINI_API_KEY" or "MOONSHOT_API_KEY"

**If deployment is still running, wait for it to complete first!**

---

## ✅ Step 2: Check Environment Variables on Render

1. In Render dashboard, go to **Environment** tab
2. Verify these variables exist:
   - `GEMINI_API_KEY` - should have a value (not empty)
   - `MOONSHOT_API_KEY` (optional) - should have a value if using Moonshot
   - `AI_PROVIDER` (optional) - should be `auto`, `gemini`, or `moonshot`

**If variables are missing:**
- The `render.yaml` file should have set them automatically
- If not, manually add them via the Environment tab → Edit button
- After adding, restart the service (Manual Deploy)

---

## ✅ Step 3: Check Render Logs for API Errors

1. In Render dashboard, go to **Logs** tab
2. Look for errors when clicking "AI Help" button:
   - ❌ "GEMINI_API_KEY not set" - API key not configured
   - ❌ "AI service unavailable" - API call failed
   - ❌ "503 Service Unavailable" - API service is down
   - ❌ "401 Unauthorized" - API key is invalid

**Common errors and fixes:**

| Error | Cause | Fix |
|-------|-------|-----|
| "GEMINI_API_KEY not set" | Environment variable not set | Add `GEMINI_API_KEY` in Render Environment tab |
| "AI service unavailable" | API endpoint unreachable | Check internet connection, API service status |
| "401 Unauthorized" | Invalid API key | Verify API key is correct from Google AI Studio |
| "503 Service Unavailable" | API service down | Wait and retry, or check Google AI status |

---

## ✅ Step 4: Test the API Endpoint Directly

Test if the backend is working by calling the API directly:

```bash
# Replace YOUR_JWT_TOKEN with a valid reader token
# Replace YOUR_QUESTION with your question

curl -X POST https://alice-suite-go.onrender.com/api/ai/ask \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "book_id": "alice-in-wonderland",
    "interaction_type": "chat",
    "question": "What is this section about?",
    "context": "Test context"
  }'
```

**Expected response:**
- ✅ `200 OK` with JSON response containing `response` field
- ❌ `503 Service Unavailable` - API key issue
- ❌ `401 Unauthorized` - Authentication issue
- ❌ `500 Internal Server Error` - Backend error

---

## ✅ Step 5: Verify API Keys Are Valid

### Test Gemini API Key:

```bash
# Replace YOUR_GEMINI_API_KEY with your actual key
curl "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=YOUR_GEMINI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "contents": [{
      "parts": [{
        "text": "Say hello"
      }]
    }]
  }'
```

**Expected:** `200 OK` with a response
**If 401:** API key is invalid - get a new one from https://aistudio.google.com/

### Test Moonshot API Key:

```bash
# Replace YOUR_MOONSHOT_KEY with your actual key
curl https://api.moonshot.cn/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_MOONSHOT_KEY" \
  -d '{
    "model": "moonshot-v1-8k",
    "messages": [{
      "role": "user",
      "content": "Say hello"
    }]
  }'
```

**Expected:** `200 OK` with a response
**If 401:** API key is invalid - check your Moonshot account

---

## ✅ Step 6: Check Browser Console for Frontend Errors

1. Open your app: https://alice-suite-go.onrender.com
2. Log in as a reader
3. Open Browser Developer Tools (F12 or Cmd+Option+I)
4. Go to **Console** tab
5. Click "AI Help" button and ask a question
6. Look for JavaScript errors:
   - ❌ `404 Not Found` - API endpoint not found
   - ❌ `401 Unauthorized` - Not logged in
   - ❌ `503 Service Unavailable` - Backend error
   - ❌ `Network Error` - Connection issue

---

## ✅ Step 7: Verify Backend Code is Deployed

The latest code should include:
- ✅ `internal/services/ai_service.go` - Supports both Gemini and Moonshot
- ✅ `internal/handlers/api.go` - Has `/api/ai/ask` endpoint
- ✅ `internal/templates/reader/interaction.html` - Has AI Help modal and functions

**Check if deployed:**
1. Look at GitHub commits - should have commits for "Tier 2: AI service"
2. Check Render deployment logs - should show the latest commit
3. If not deployed, trigger Manual Deploy in Render dashboard

---

## 🔧 Common Fixes

### Fix 1: API Key Not Set on Render
**Problem:** Environment variable not in Render
**Solution:**
1. Render dashboard → Environment tab
2. Click Edit button
3. Add `GEMINI_API_KEY` with your key
4. Save and restart service

### Fix 2: Invalid API Key
**Problem:** API key is wrong or expired
**Solution:**
1. Go to https://aistudio.google.com/
2. Generate a new API key
3. Update `render.yaml` or Render Environment tab
4. Redeploy

### Fix 3: Deployment Didn't Include Latest Code
**Problem:** Old code is still running
**Solution:**
1. Check Render deployment logs
2. If needed, trigger "Manual Deploy"
3. Wait for deployment to complete (2-5 minutes)

### Fix 4: Frontend Not Updated
**Problem:** Browser has cached old JavaScript
**Solution:**
1. Hard refresh browser: `Ctrl+Shift+R` (Windows) or `Cmd+Shift+R` (Mac)
2. Or clear browser cache

---

## 📞 Still Not Working?

If none of these steps fix the issue:

1. **Check Render Logs** for specific error messages
2. **Check Browser Console** for JavaScript errors
3. **Verify API keys** are correct and active
4. **Test API endpoints** directly with curl
5. **Check deployment status** - ensure latest code is deployed

Share the specific error message you see, and I can help troubleshoot further!
