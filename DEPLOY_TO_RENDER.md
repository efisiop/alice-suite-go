# How to Deploy Latest Changes to Render.com

## Quick Steps to Deploy

### Option 1: Manual Deploy Trigger (Recommended)

1. **Go to Render Dashboard:**
   - Visit: https://dashboard.render.com
   - Log in to your account
   - Find your "alice-suite-go" service

2. **Trigger Manual Deploy:**
   - Click on your service
   - Look for "Manual Deploy" button (usually in the top right)
   - Click "Deploy latest commit"
   - This will pull the latest code from GitHub and rebuild

3. **Monitor the Deployment:**
   - Watch the build logs
   - You should see:
     - Code being pulled from GitHub
     - Go modules downloading
     - Binaries being built
     - Migrations running (including 011_add_alice_characters.sql)
     - Server starting

### Option 2: Push an Empty Commit (If Manual Deploy Not Available)

If you don't see a manual deploy button, you can trigger a deployment by pushing a small change:

```bash
git commit --allow-empty -m "Trigger Render deployment"
git push origin main
```

This will trigger Render's auto-deploy.

## What Should Happen During Deployment

### Build Phase:
1. ✅ Pulls latest code from GitHub (commit c7abeaa)
2. ✅ Downloads Go dependencies
3. ✅ Builds binaries:
   - `bin/server`
   - `bin/migrate`
   - `bin/init-users`
   - `bin/fix-render`
   - `bin/verify-deployment`

### Start Phase (start.sh):
1. ✅ Creates data directory
2. ✅ Runs migrations (including 011_add_alice_characters.sql)
3. ✅ Initializes users
4. ✅ Fixes sections data
5. ✅ Starts server

## Verify Deployment Worked

After deployment completes:

1. **Check the app is running:**
   - Visit your Render URL (e.g., https://alice-suite-go.onrender.com)
   - Should load without errors

2. **Check glossary terms are highlighted:**
   - Log in to the reader app
   - Navigate to a page
   - Look for words highlighted in **light blue** (#6BA3D6)
   - These are glossary terms

3. **Test character names:**
   - Click on "Alice", "White Rabbit", or "Cheshire Cat"
   - Should show character definitions from the glossary

4. **Check deployment logs:**
   - In Render dashboard → Your Service → Logs
   - Look for:
     - "✅ Migration 011_add_alice_characters.sql completed"
     - "Database initialized successfully"
     - No errors about NULL values

## Troubleshooting

### If glossary terms still not highlighted:

1. **Check if migrations ran:**
   - In Render Shell, run:
     ```bash
     sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM alice_glossary WHERE id LIKE 'character-%';"
     ```
   - Should return 11 (number of characters added)

2. **Check if code was deployed:**
   - In Render Shell, check the code:
     ```bash
     grep -n "glossary-term" internal/templates/reader/interaction.html | head -5
     ```
   - Should show the CSS with color #6BA3D6

3. **Clear browser cache:**
   - Hard refresh: Ctrl+Shift+R (Windows) or Cmd+Shift+R (Mac)
   - Or clear browser cache completely

### If deployment fails (or you get a "deploy didn't complete" email):

1. **Find the actual error (most important):**
   - Go to [Render Dashboard](https://dashboard.render.com) → your **alice-suite-go** service
   - Open the **Logs** tab (or the failed deploy’s log)
   - Check **Build logs** first: if the failure is there, the error message will show (e.g. Go compile error, missing binary, timeout)
   - Then check **Deploy / Runtime logs**: if the build succeeded but the service didn’t start, the error will be there (e.g. `DATABASE_URL` missing, migration failed, `./start.sh` failed)

2. **Typical causes:**
   - **Build:** Go build failed (syntax/import error), one of the extra binaries (`diagnose-sections`, `compare-db-structure`) failed to build, or build timed out
   - **Start:** `DATABASE_URL` not set or PostgreSQL not linked, or a step in `start.sh` failed (migrate, init-users, fix-render)

3. **Quick checks on your machine:**
   - Confirm the latest commit is on GitHub: `git log --oneline -1`
   - Optional: try building locally: `CGO_ENABLED=1 go build -o bin/server ./cmd/server` (needs C compiler for SQLite)

4. **Verify migration file exists (if logs mention migrations):**
   ```bash
   ls -la migrations/011_add_alice_characters.sql
   ```

## Expected Results After Deployment

✅ Glossary terms highlighted in light blue (#6BA3D6)  
✅ Character names (Alice, White Rabbit, etc.) in glossary  
✅ Clicking glossary terms shows definitions from database  
✅ No 500 errors when loading glossary API  
✅ All 1,212+ glossary terms available (including characters)
