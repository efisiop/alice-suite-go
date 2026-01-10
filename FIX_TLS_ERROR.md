# Fix: TLS Certificate Error with Moonshot API

## 🐛 Problem

You're getting: `tls: failed to parse certificate from server: x509: negative serial number`

**Root cause:**
- `GEMINI_API_KEY` is **NOT SET** locally
- `MOONSHOT_API_KEY` **IS SET**
- Service tries Gemini first, but fails (no key), then tries Moonshot
- Moonshot API server has a TLS certificate issue

## ✅ Solution 1: Use Gemini (Recommended - Fastest Fix)

**Set GEMINI_API_KEY so it uses Gemini instead of Moonshot:**

```bash
export GEMINI_API_KEY="AIzaSyB2RpUVMID-JJTLL4PakpZHDCqZuI_OZio"
```

**Restart your server:**
```bash
./start_dev_server.sh
```

**Result:** Gemini will be used first (free, fast, works). No more Moonshot TLS error! ✅

---

## ✅ Solution 2: Fix Moonshot TLS Issue (If You Must Use Moonshot)

If you really want to use Moonshot despite the certificate issue, you can bypass TLS verification:

```bash
# Allow skipping TLS verification for Moonshot (not recommended for production)
export MOONSHOT_SKIP_TLS_VERIFY="true"
```

**Restart your server:**
```bash
./start_dev_server.sh
```

**⚠️ Warning:** This reduces security by skipping TLS certificate verification. Only use for development/testing.

---

## ✅ Solution 3: Force Gemini Only (Best Practice)

Set AI_PROVIDER to only use Gemini:

```bash
export GEMINI_API_KEY="AIzaSyB2RpUVMID-JJTLL4PakpZHDCqZuI_OZio"
export AI_PROVIDER="gemini"  # Only use Gemini, never Moonshot
```

**Restart your server:**
```bash
./start_dev_server.sh
```

---

## 🔧 What I Fixed in the Code

1. ✅ **Only tries providers with configured keys** - Won't try Gemini if key not set
2. ✅ **Better logging** - Shows which provider is being tried
3. ✅ **TLS skip option** - Can bypass TLS verification for Moonshot if needed
4. ✅ **Better error messages** - Shows which provider failed and why

---

## 🧪 Test After Fix

1. **Set GEMINI_API_KEY** (see Solution 1 above)
2. **Restart server**
3. **Check startup logs** - Should see: "✅ AI API key configured"
4. **Click "AI Help"** and ask a question
5. **Check server logs** - Should see: "Trying Gemini API..." then "AI API call successful using gemini"

---

## 📝 Make It Permanent

Add to your `~/.zshrc` (or `~/.bashrc`):

```bash
export GEMINI_API_KEY="AIzaSyB2RpUVMID-JJTLL4PakpZHDCqZuI_OZio"
export AI_PROVIDER="gemini"  # Optional: force Gemini only
# Don't set MOONSHOT_SKIP_TLS_VERIFY unless you really need Moonshot
```

Then reload:
```bash
source ~/.zshrc
```

---

## ✅ Recommended Setup

**Best practice for local development:**

```bash
export GEMINI_API_KEY="your-gemini-key"
export AI_PROVIDER="gemini"  # Use Gemini only (free tier, fast, reliable)
# Don't set MOONSHOT_API_KEY or ANTHROPIC_BASE_URL unless you need Moonshot
```

This way:
- ✅ Uses Gemini (free, fast, works)
- ✅ No Moonshot TLS issues
- ✅ Simpler setup

**After setting, restart your server and it should work!**
