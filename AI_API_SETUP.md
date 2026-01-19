# AI API Setup Guide - Tier 2

This guide shows you where to paste your **Gemini API key**, optionally configure your **Moonshot API key**, and optionally configure your **Freepik API key** for Tier 2 AI assistance features including image generation.

## 🔑 Where to Set API Keys

### Option 1: Local Development (Mac/Linux)

**For your current shell session:**
```bash
export GEMINI_API_KEY="your-gemini-api-key-here"
export MOONSHOT_API_KEY="your-moonshot-api-key-here"  # Optional
export FREEPIK_API_KEY="your-freepik-api-key-here"  # Optional - for visual example image generation
export AI_PROVIDER="auto"  # Optional: "gemini", "moonshot", or "auto" (default)
```

**To make it permanent (recommended), add to your shell config file:**

**For Zsh (default on Mac):**
```bash
echo 'export GEMINI_API_KEY="your-gemini-api-key-here"' >> ~/.zshrc
echo 'export MOONSHOT_API_KEY="your-moonshot-api-key-here"' >> ~/.zshrc  # Optional
echo 'export FREEPIK_API_KEY="your-freepik-api-key-here"' >> ~/.zshrc  # Optional - for visual example image generation
echo 'export AI_PROVIDER="auto"' >> ~/.zshrc  # Optional
source ~/.zshrc
```

**For Bash:**
```bash
echo 'export GEMINI_API_KEY="your-gemini-api-key-here"' >> ~/.bashrc
echo 'export MOONSHOT_API_KEY="your-moonshot-api-key-here"' >> ~/.bashrc  # Optional
echo 'export FREEPIK_API_KEY="your-freepik-api-key-here"' >> ~/.bashrc  # Optional - for visual example image generation
source ~/.bashrc
```

### Option 2: Render.com Deployment

**Method 1: Via Render Dashboard - Edit Button (Recommended)**
Based on your Render dashboard:

1. Go to your Render.com dashboard: https://dashboard.render.com
2. Click on your `alice-suite-go` service
3. You should be on the **Environment** tab (you can see "Environment Variables" section with `JWT_SECRET` listed)
4. **Click the "Edit" button** (pencil icon ✏️) that appears next to "Environment Variables"
5. A modal/popup should appear where you can edit environment variables
6. In that modal, look for:
   - A **"+"** or **"Add Variable"** or **"Add Environment Variable"** button
   - OR you might see an empty row where you can type directly
7. Add the new variable:
   - In the **KEY** field: `GEMINI_API_KEY`
   - In the **VALUE** field: `your-gemini-api-key-here`
   - Click **Save**, **Add**, or **Apply** (depending on what buttons appear)
8. (Optional) Add Freepik API key for visual example image generation:
   - In the **KEY** field: `FREEPIK_API_KEY`
   - In the **VALUE** field: `your-freepik-api-key-here`
   - Click **Save**, **Add**, or **Apply**
9. If you want to add other optional variables, repeat step 7:
   - **KEY:** `MOONSHOT_API_KEY`, **VALUE:** `your-moonshot-api-key-here`
   - **KEY:** `AI_PROVIDER`, **VALUE:** `auto`
9. Click **Save Changes** or **Apply** to confirm
10. **Restart your service** (the service should auto-restart, or use "Manual Deploy")

**If Method 1 doesn't work (no Edit button or can't find Add option):**

**Method 2: Via "+ Create environment" button**
1. Look for the **"+ Create environment"** button next to the "Environment" heading
2. Click it - this might open a dialog to add environment variables
3. Follow the prompts to add `GEMINI_API_KEY`

**Method 3: Via render.yaml (Alternative - Less Secure)**
⚠️ **Warning:** This method will add the API key to your Git repository. Only use if you trust your repository is private.

Add to your `render.yaml` file:
```yaml
services:
  - type: web
    name: alice-suite-go
    # ... existing config ...
    envVars:
      - key: PORT
        value: 10000
      - key: ENV
        value: production
      - key: JWT_SECRET
        generateValue: true
      - key: DB_PATH
        value: data/alice-suite.db
      - key: GEMINI_API_KEY
        value: your-gemini-api-key-here  # ⚠️ This will be in Git!
      - key: MOONSHOT_API_KEY  # Optional
        value: your-moonshot-api-key-here
      - key: FREEPIK_API_KEY  # Optional - for visual example image generation
        value: your-freepik-api-key-here
      - key: AI_PROVIDER  # Optional
        value: auto
```

Then commit and push to GitHub - Render will auto-deploy with the new variables.

**Method 4: Using Render API (Advanced)**
If the UI doesn't work, you can use Render's API:
```bash
curl -X PATCH "https://api.render.com/v1/services/srv-d4uunpm3jp1c73eahffg/env-vars" \
  -H "Authorization: Bearer YOUR_RENDER_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "envVars": {
      "GEMINI_API_KEY": "your-gemini-api-key-here"
    }
  }'
```

**Security Note:** Method 3 (render.yaml) exposes your API key in Git. **Prefer Method 1 or 2** (Dashboard UI) which keeps keys secure and secret.

## 🔧 How It Works

### Provider Selection (`AI_PROVIDER`)

- **`auto`** (default): Tries Gemini first, automatically falls back to Moonshot if Gemini fails or is not configured
- **`gemini`**: Uses only Gemini API
- **`moonshot`**: Uses only Moonshot API

### Recommended Setup

1. **For Development/Testing (Free):**
   ```
   GEMINI_API_KEY=your-key
   AI_PROVIDER=gemini
   ```
   - Uses Gemini's free tier (~60 requests/minute)
   - No cost for development

2. **For Production (Flexible):**
   ```
   GEMINI_API_KEY=your-key
   MOONSHOT_API_KEY=your-key
   AI_PROVIDER=auto
   ```
   - Automatically uses Gemini (free tier) when available
   - Falls back to Moonshot if Gemini fails
   - Best of both worlds

3. **For Production (Moonshot Only):**
   ```
   MOONSHOT_API_KEY=your-key
   AI_PROVIDER=moonshot
   ```
   - If you prefer Moonshot's performance or pricing
   - No fallback

## 📝 Getting Your API Keys

### Gemini API Key

1. Go to [Google AI Studio](https://aistudio.google.com/)
2. Sign in with your Google account
3. Click **Get API Key**
4. Create a new project or select existing one
5. Copy the API key
6. Paste it where indicated above

### Freepik API Key (Optional - for Visual Examples)

The Freepik API key enables image generation for the "Visual Example" feature. Without it, visual examples will only show text descriptions.

1. Go to [Freepik API Documentation](https://docs.freepik.com/)
2. Sign up or log in to your Freepik account
3. Navigate to API settings or developer dashboard
4. Generate or copy your API key
5. Add it as `FREEPIK_API_KEY` environment variable (see setup instructions above)

**Note:** Freepik uses a credit-based system. Check their [pricing page](https://www.freepik.com/api/pricing) for details. The free tier includes limited credits for testing.

## 🧪 Testing Your Setup

After setting the environment variables, restart your server:

```bash
# Stop current server (Ctrl+C)
# Then restart
./start.sh
```

Or if running manually:
```bash
go run cmd/server/main.go
```

The AI service will automatically detect which providers are configured and use them accordingly.

## 🔍 Troubleshooting

**Problem:** "AI service not configured"
- **Solution:** Make sure `GEMINI_API_KEY` or `MOONSHOT_API_KEY` is set as an environment variable

**Problem:** "Gemini API error: API key not valid"
- **Solution:** Double-check your Gemini API key from Google AI Studio

**Problem:** "Moonshot API error: 401 Unauthorized"
- **Solution:** Verify your Moonshot API key is correct

**Problem:** On Render.com, keys don't work
- **Solution:** Make sure you added them via the Environment tab in Render dashboard, then restart the service

**Problem:** Visual examples show text but no images
- **Solution:** Make sure `FREEPIK_API_KEY` is set. Without it, only text descriptions will be shown.

## ✅ Verification

You can verify your setup is working by:
1. Starting the server
2. Logging in as a reader
3. Clicking "AI Help" button
4. Asking a question - you should get an AI response!
