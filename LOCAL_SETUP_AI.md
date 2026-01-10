# Local Setup: AI Assistant

To test the AI assistant on localhost, you need to set the `GEMINI_API_KEY` environment variable.

## ✅ Quick Setup

### Option 1: Set for Current Terminal Session

```bash
export GEMINI_API_KEY="your-gemini-api-key-here"
```

Then start your server:
```bash
./start_dev_server.sh
# or
go run cmd/server/main.go
```

### Option 2: Set in Shell Config (Permanent)

**For Zsh (default on Mac):**
```bash
echo 'export GEMINI_API_KEY="your-gemini-api-key-here"' >> ~/.zshrc
source ~/.zshrc
```

**For Bash:**
```bash
echo 'export GEMINI_API_KEY="your-gemini-api-key-here"' >> ~/.bashrc
source ~/.bashrc
```

### Option 3: Set in start_dev_server.sh

Edit `start_dev_server.sh` and add:
```bash
#!/bin/bash
export GEMINI_API_KEY="your-gemini-api-key-here"
# ... rest of the script
```

---

## 🔍 Verify Setup

Check if the environment variable is set:
```bash
echo $GEMINI_API_KEY
```

If it shows your API key, it's set correctly. If it's empty, the variable is not set.

---

## 🧪 Test the Setup

1. **Start the server** (make sure GEMINI_API_KEY is set):
   ```bash
   ./start_dev_server.sh
   ```

2. **Check server logs** for any errors:
   - Look for: "GEMINI_API_KEY not set" or similar errors
   - If you see this, the environment variable is not being read

3. **Test in browser**:
   - Go to: http://localhost:8080
   - Log in as a reader
   - Click "AI Help" button
   - Ask a question
   - Check browser console (F12) for errors

4. **Check server terminal**:
   - Look for error messages when you click "Ask AI"
   - You should now see detailed error messages (after the fix)

---

## 🐛 Troubleshooting

### Problem: "500 Internal Server Error"

**Check server logs for the actual error:**

After starting your server, when you ask a question, you should see error messages in the server terminal like:
- "Error in HandleAskAI: GEMINI_API_KEY not set"
- "Error in HandleAskAI: Gemini API error (status 401): ..."
- etc.

**Common fixes:**

1. **"GEMINI_API_KEY not set"**
   - Solution: Set the environment variable (see above)
   - Verify with: `echo $GEMINI_API_KEY`

2. **"Gemini API error (status 401)"**
   - Solution: Your API key is invalid
   - Get a new key from: https://aistudio.google.com/

3. **"Gemini API error (status 400)"**
   - Solution: Request format issue (check code)
   - This is a code bug - let me know and I'll fix it

4. **"Network error" or "Connection refused"**
   - Solution: Check your internet connection
   - The server needs internet access to call Gemini API

---

## 📝 Getting Your Gemini API Key

1. Go to: https://aistudio.google.com/
2. Sign in with your Google account
3. Click "Get API Key" or go to API Keys section
4. Create a new API key (or use existing one)
5. Copy the key (starts with `AIza...`)
6. Set it as environment variable (see above)

---

## ✅ After Setup

Once `GEMINI_API_KEY` is set:

1. **Restart your server** (if it was already running)
2. **Test the AI Help feature**:
   - Click "AI Help" button
   - Ask: "What is this section about?"
   - You should get an AI response!

3. **Check server logs** for success:
   - Should see successful API calls
   - No error messages

---

## 🔄 Optional: Set Moonshot Key Too

If you also want to use Moonshot (as a fallback):

```bash
export MOONSHOT_API_KEY="your-moonshot-api-key-here"
export AI_PROVIDER="auto"  # Try Gemini first, fallback to Moonshot
```

This will make the system try Gemini first, and if it fails, automatically try Moonshot.

---

## 📞 Still Not Working?

1. **Check server terminal** for error messages (should now show detailed errors)
2. **Check browser console** (F12) for frontend errors
3. **Verify API key** is set: `echo $GEMINI_API_KEY`
4. **Test API key directly**:
   ```bash
   curl "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=YOUR_GEMINI_API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"contents":[{"parts":[{"text":"Say hello"}]}]}'
   ```

If the curl works but the app doesn't, there's a code issue. Share the server error logs and I'll help fix it!
