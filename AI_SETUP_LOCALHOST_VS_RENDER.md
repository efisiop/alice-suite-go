# AI Setup: Localhost vs Render.com

## 🔍 Current Status

### ✅ Render.com (Production) - Already Configured
- **GEMINI_API_KEY**: ✅ Set in `render.yaml` (line 19-20)
- **MOONSHOT_API_KEY**: ✅ Set in `render.yaml` (line 21-22)
- **AI_PROVIDER**: ✅ Set to `auto` in `render.yaml` (line 22-23)
- **Status**: Should work automatically after deployment! ✅

### ❌ Localhost (Development) - Needs Setup
- **GEMINI_API_KEY**: ❌ NOT SET (that's why you got the 500 error)
- **MOONSHOT_API_KEY**: ✅ SET (but has TLS certificate issue)
- **ANTHROPIC_BASE_URL**: ❌ Set to wrong URL (`https://api.moonshot.ai/anthropic`)
- **Status**: Needs manual setup ❌

---

## 🛠️ Fix for Localhost (Development)

You need to set the environment variables in your **local terminal**:

```bash
# Set Gemini API key (same one from render.yaml)
export GEMINI_API_KEY="your-new-gemini-api-key-here"

# Remove the incorrect Moonshot URL
unset ANTHROPIC_BASE_URL

# Optional: Force Gemini only
export AI_PROVIDER="gemini"

# Restart your local server
./start_dev_server.sh
```

**Result:** AI Help will work on localhost! ✅

---

## ✅ Render.com (Production) - Should Already Work

The API keys are **already configured** in `render.yaml`, so Render should automatically:
1. Read the keys from `render.yaml` during deployment
2. Set them as environment variables
3. Make AI Help work automatically

**To verify Render is working:**
1. Wait for latest deployment to complete (from our code commits)
2. Go to: https://alice-suite-go.onrender.com
3. Log in as a reader
4. Click "AI Help" button
5. Ask a question - should work! ✅

**If Render still doesn't work:**
- Check Render dashboard → Environment tab
- Verify `GEMINI_API_KEY` is listed there
- Check Render logs for errors
- See `AI_TROUBLESHOOTING.md` for detailed steps

---

## 📋 Summary

| Location | GEMINI_API_KEY | Status | Action Needed |
|----------|---------------|--------|---------------|
| **Localhost** | ❌ Not set | ❌ Not working | Set with `export GEMINI_API_KEY=...` |
| **Render.com** | ✅ Set in render.yaml | ✅ Should work | Wait for deployment, then test |

---

## 🚀 Quick Test Both

### Test Localhost:
```bash
# Set the key
export GEMINI_API_KEY="your-new-gemini-api-key-here"

# Start server
./start_dev_server.sh

# Test in browser: http://localhost:8080
# Click "AI Help" → Should work!
```

### Test Render:
```bash
# No setup needed - keys are in render.yaml
# Just wait for deployment to complete

# Test in browser: https://alice-suite-go.onrender.com
# Click "AI Help" → Should work!
```

---

## 🔧 Make Localhost Setup Permanent

To avoid setting the key every time you open a new terminal:

**Add to `~/.zshrc` (Mac) or `~/.bashrc` (Linux):**

```bash
echo 'export GEMINI_API_KEY="your-new-gemini-api-key-here"' >> ~/.zshrc
echo 'export AI_PROVIDER="gemini"' >> ~/.zshrc
source ~/.zshrc
```

**Or create a `.env` file** (if your server supports it):
```bash
# .env file in project root
GEMINI_API_KEY=AIzaSyB2RpUVMID-JJTLL4PakpZHDCqZuI_OZio
AI_PROVIDER=gemini
```

---

## ✅ Recommended Setup

**For localhost development:**
- Set `GEMINI_API_KEY` in your shell config (`~/.zshrc`)
- Use Gemini only (`AI_PROVIDER=gemini`)
- Skip Moonshot (to avoid TLS issues)

**For Render.com production:**
- Already configured in `render.yaml` ✅
- Uses `auto` mode (tries Gemini first, falls back to Moonshot)
- Should work automatically after deployment ✅

---

## 🐛 If Render Doesn't Work

Even though keys are in `render.yaml`, sometimes Render needs them manually set:

1. **Check Render Dashboard:**
   - Go to your service → Environment tab
   - Verify `GEMINI_API_KEY` is listed
   - If missing, add it manually

2. **Check Deployment Logs:**
   - Look for errors mentioning "GEMINI_API_KEY"
   - Look for successful API calls

3. **Manual Setup (if needed):**
   - Render Dashboard → Environment → Edit
   - Add: `GEMINI_API_KEY` = `your-key-here`
   - Restart service

---

## 📞 Need Help?

- **Localhost issues:** See `LOCAL_SETUP_AI.md`
- **Render issues:** See `AI_TROUBLESHOOTING.md`
- **TLS errors:** See `FIX_TLS_ERROR.md`
