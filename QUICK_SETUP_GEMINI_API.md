# Quick Setup: Adding Gemini API Key to Render.com

Based on your Render dashboard screenshot, here are the **sure-fire methods** to add your Gemini API key:

## ✅ Method 1: Using render.yaml (MOST RELIABLE)

This method **will definitely work** because Render reads this file directly from your Git repository.

### Steps:

1. **Edit the `render.yaml` file** in your project:
   ```bash
   # In your project directory
   code render.yaml  # or use your preferred editor
   ```

2. **Add the GEMINI_API_KEY** to the `envVars` section. Your file should look like this:

   ```yaml
   services:
     - type: web
       name: alice-suite-go
       runtime: go
       plan: free
       buildCommand: go mod download && CGO_ENABLED=1 go build -o bin/server ./cmd/server && CGO_ENABLED=1 go build -o bin/migrate ./cmd/migrate && CGO_ENABLED=1 go build -o bin/init-users ./cmd/init-users && CGO_ENABLED=1 go build -o bin/fix-render ./cmd/fix-render && CGO_ENABLED=1 go build -o bin/verify-deployment ./cmd/verify-deployment
       startCommand: ./start.sh
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
           value: YOUR_GEMINI_API_KEY_HERE  # ← Add this line (replace with your actual key)
         - key: AI_PROVIDER
           value: auto  # ← Optional: Add this line
   ```

3. **Commit and push to GitHub:**
   ```bash
   git add render.yaml
   git commit -m "Add Gemini API key configuration"
   git push origin main
   ```

4. **Render will automatically deploy** with the new environment variable!

⚠️ **IMPORTANT SECURITY WARNING:**
- This adds your API key to your Git repository
- **ONLY use this method if your GitHub repository is PRIVATE**
- If your repo is public, use Method 2 instead

---

## ✅ Method 2: Using Render Dashboard UI (If Method 1 is too risky)

If your repository is public, use this method instead:

### Steps:

1. **Click the "Edit" button** (pencil icon ✏️) next to "Environment Variables" in your Render dashboard

2. **In the modal that opens:**
   - Look for a **table** with columns "KEY" and "VALUE"
   - You should see `JWT_SECRET` already listed
   - Look for one of these:
     - A **"+" button** or **"Add Variable"** button (often at the top or bottom of the table)
     - An **empty row** at the bottom of the table where you can type
     - A **"New Variable"** link or button

3. **Add the new variable:**
   - **KEY column:** Type `GEMINI_API_KEY`
   - **VALUE column:** Type your actual API key (e.g., `AIza...`)
   - Click **Save**, **Apply**, or **Update** (button name varies)

4. **If you still can't find "Add":**
   - The modal might have a **"Add Environment Variable"** button at the top
   - OR look for a **"+ Create environment"** button at the top of the Environment section
   - OR try right-clicking in the table area

5. **After adding:**
   - Click **Save Changes** or **Apply**
   - Your service should automatically restart, or click **Manual Deploy**

---

## 🔍 Alternative: Check Render Documentation

If neither method works, Render's UI might have changed. Try:

1. Go to: https://render.com/docs/environment-variables
2. Check the latest instructions for adding environment variables
3. Or contact Render support

---

## ✅ Quick Test After Setup

Once you've added the key (using either method), test it:

1. **Wait for Render to deploy** (should happen automatically if using Method 1)
2. **Or manually deploy:** Click "Manual Deploy" in Render dashboard
3. **Test in your app:**
   - Go to: https://alice-suite-go.onrender.com
   - Log in as a reader
   - Click "AI Help" button
   - Ask a question - you should get an AI response!

---

## 🆘 Still Having Issues?

If you can't find the "Edit" button or add variables in the UI:

1. **Screenshot the Environment tab** and I can help identify the exact button
2. **Or use Method 1** (render.yaml) - it's the most reliable method
3. **Just remember:** Only use render.yaml if your GitHub repo is private!
