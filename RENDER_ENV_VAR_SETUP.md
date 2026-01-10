# How to Set Environment Variables in Render.com

## Simple Step-by-Step Guide

### Method 1: Using the Render Dashboard (Recommended)

1. **Login to Render:**
   - Go to: https://dashboard.render.com
   - Sign in to your account

2. **Find Your Service:**
   - You'll see a list of your services
   - Click on **`alice-suite-go`** (or whatever your service is named)

3. **Open Environment Settings:**
   - On the service page, look for tabs at the top (like "Events", "Metrics", "Logs", **"Environment"**)
   - Click on the **"Environment"** tab
   - OR look in the left sidebar for "Environment"

4. **Add New Variable:**
   - Scroll down until you see "Environment Variables" section
   - You'll see a list of existing variables
   - Click the button that says **"+ Add Environment Variable"** or **"Add"** (usually at the top right or bottom of the list)

5. **Enter the Variable:**
   - A form/popup will appear with two fields:
     - **Key**: Type exactly: `GEMINI_API_KEY`
     - **Value**: Paste your new Gemini API key
   - Click **"Save"** or **"Add"**

6. **Wait for Restart:**
   - Render will automatically restart your service (2-3 minutes)
   - You'll see a notification or the status will change to "Deploying"

### Method 2: If You Don't See "Environment" Tab

Sometimes the interface varies. Try these locations:

1. **In Settings:**
   - Click on your service
   - Look for "Settings" in the left sidebar
   - Find "Environment Variables" section there

2. **In the Service Overview:**
   - Some Render layouts show environment variables directly on the main service page
   - Scroll down to find "Environment Variables" section

### What You Need to Add

- **Variable Name**: `GEMINI_API_KEY`
- **Variable Value**: `your-new-api-key-from-google`

### Verify It's Set

After adding:
1. Refresh the page
2. Look for `GEMINI_API_KEY` in the environment variables list
3. Check that it shows as set (value might be hidden with dots for security)

### Need Screenshots?

If you're still stuck, tell me what you see on your Render dashboard and I can guide you more specifically!
