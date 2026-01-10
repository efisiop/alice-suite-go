# Quick Fix: AI 500 Error

## 🐛 Problem Found

Your local environment has:
- ❌ `GEMINI_API_KEY` **NOT SET** (that's why it's falling back to Moonshot)
- ❌ `ANTHROPIC_BASE_URL` set to **WRONG URL**: `https://api.moonshot.ai/anthropic` (incorrect)

## ✅ Solution: Fix Your Local Environment

### Option 1: Use Gemini (Recommended - You Already Have the Key)

```bash
# Set Gemini API key (from your render.yaml)
export GEMINI_API_KEY="your-new-gemini-api-key-here"

# Unset the incorrect Moonshot URL
unset ANTHROPIC_BASE_URL

# Restart your server
./start_dev_server.sh
```

**Result:** Gemini will be used (faster, free tier) ✅

### Option 2: Fix Moonshot URL (If You Prefer Moonshot)

```bash
# Set correct Moonshot URL
export MOONSHOT_BASE_URL="https://api.moonshot.cn/v1"

# Or unset the incorrect one (code will use default)
unset ANTHROPIC_BASE_URL

# Make sure you have MOONSHOT_API_KEY set
export MOONSHOT_API_KEY="your-moonshot-key-here"

# Restart your server
./start_dev_server.sh
```

### Option 3: Make It Permanent

Add to your `~/.zshrc` (or `~/.bashrc`):

```bash
# Add these lines
export GEMINI_API_KEY="your-new-gemini-api-key-here"
unset ANTHROPIC_BASE_URL  # Remove incorrect URL

# Then reload
source ~/.zshrc
```

---

## 🔍 What Was Wrong?

1. **GEMINI_API_KEY not set** → Service tried Gemini first, failed, then fell back to Moonshot
2. **ANTHROPIC_BASE_URL set incorrectly** → Pointed to wrong API endpoint (`api.moonshot.ai/anthropic` instead of `api.moonshot.cn/v1`)
3. **TLS certificate error** → Wrong endpoint had invalid certificate

---

## ✅ After Fixing

1. **Restart your server** (stop with Ctrl+C, then `./start_dev_server.sh`)
2. **Check startup message** - should say "✅ AI API key configured"
3. **Test AI Help** - should work with Gemini now!

---

## 📝 Verification

Check your environment:
```bash
echo "GEMINI_API_KEY: ${GEMINI_API_KEY:+✅ set}"
echo "ANTHROPIC_BASE_URL: ${ANTHROPIC_BASE_URL:-✅ not set (good)}"
```

If GEMINI_API_KEY shows "✅ set", you're good!

---

## 🔧 Code Fix

I've also updated the code to:
- ✅ Detect and fix incorrect Moonshot URLs automatically
- ✅ Use `MOONSHOT_BASE_URL` instead of `ANTHROPIC_BASE_URL` (clearer name)
- ✅ Warn when incorrect URLs are detected

This fix is already in the code (committed). Just set your `GEMINI_API_KEY` and restart!
