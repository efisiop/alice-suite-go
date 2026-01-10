# ⚠️ URGENT: API Key Security Issue

## What Happened
Your Gemini API key was exposed in your code repository and Google has revoked it for security reasons.

## What You Need to Do NOW:

### 1. Generate a New Gemini API Key
1. Go to: https://aistudio.google.com/app/apikey
2. Click "Create API Key"
3. Copy your new API key

### 2. Set the New Key in Render Dashboard (Production)

**Step-by-step instructions:**

1. **Go to Render Dashboard:**
   - Visit: https://dashboard.render.com
   - Log in to your account

2. **Find your service:**
   - Look for `alice-suite-go` in your services list
   - Click on the service name to open it

3. **Go to Environment settings:**
   - In the service page, look at the **top menu/tabs** (horizontal tabs)
   - Click on the **"Environment"** tab
   - (If you don't see tabs, look for "Environment" in the left sidebar)

4. **Add the environment variable:**
   - Scroll down to the "Environment Variables" section
   - Look for a button that says **"Add Environment Variable"** or **"Add"** or **"New Variable"**
   - Click that button
   - In the popup/form that appears:
     - **Key**: Type `GEMINI_API_KEY` (exactly like this, all caps)
     - **Value**: Paste your new API key
   - Click **"Save"** or **"Add"**

5. **Save and restart:**
   - Render should automatically save the variable
   - Your service will automatically restart/redeploy (this takes 2-3 minutes)
   - Wait for the deployment to finish

**Alternative method (if you can't find "Environment" tab):**
- Click on your service name
- Look for "Settings" or "Config" in the menu
- Find "Environment Variables" section there

### 3. Set the New Key for Local Development
For your local machine, set it as an environment variable:

**On Mac/Linux:**
```bash
export GEMINI_API_KEY="your-new-api-key-here"
```

**To make it permanent (Mac/Linux):**
```bash
echo 'export GEMINI_API_KEY="your-new-api-key-here"' >> ~/.zshrc
source ~/.zshrc
```

**On Windows:**
```powershell
setx GEMINI_API_KEY "your-new-api-key-here"
```

### 4. Verify It Works
After setting the new key:
- Restart your local server if running locally
- Wait for Render to redeploy (usually 2-3 minutes)
- Test the AI feature again

## What I Fixed:
✅ Removed exposed API keys from `render.yaml`
✅ Removed exposed keys from documentation files
✅ Created `render.yaml.example` as a safe template
✅ Updated `.gitignore` to protect sensitive files

## Important Security Rules:
🚫 **NEVER** commit API keys to git
✅ **ALWAYS** use environment variables for secrets
✅ Use Render Dashboard for production keys (not code files)

## Need Help?
- Local setup: Check `LOCAL_SETUP_AI.md`
- Render setup: Check `AI_SETUP_LOCALHOST_VS_RENDER.md`
