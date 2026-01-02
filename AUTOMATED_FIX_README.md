# Automated Fix for Render.com Sections Issue

## ✅ What's Been Automated

The fix is now **automatically applied during every Render.com deployment**!

### How It Works

1. **Build Process** (in `render.yaml`):
   - Builds the `fix-render` binary alongside other binaries
   - Runs migrations first
   - Runs `init-users` 
   - **Automatically runs `fix-render`** to ensure sections data is correct

2. **The `fix-render` Script** (`cmd/fix-render/main.go`):
   - **Diagnoses** the current database state
   - **Checks** if sections data is correct (5+ sections per page)
   - **If incorrect**: Automatically imports the correct sections data (77 sections)
   - **Safe to run multiple times** - won't duplicate data if already correct
   - Provides detailed progress output

### What Gets Fixed Automatically

- ✅ Checks if Page 1 has 5+ sections (should be 5)
- ✅ If not, imports all 77 sections from `scripts/sections-data.sql`
- ✅ Verifies the import was successful
- ✅ Safe to re-run (detects if data already exists)

## 📋 Manual Run (Optional)

If you want to run the fix manually on Render.com:

```bash
# In Render.com Shell
DB_PATH=data/alice-suite.db ./bin/fix-render
```

Or if the binary isn't built yet:

```bash
DB_PATH=data/alice-suite.db go run cmd/fix-render/main.go
```

## 🔍 What to Expect

The script will show:

```
🔧 Render.com Sections Fix Script
============================================================

✅ Database connected: data/alice-suite.db

📊 Step 1: Diagnosing current database state...
------------------------------------------------------------
   Current sections in database: 1
   Sections for page 1: 1

⚠️  Issue detected: Page 1 has only 1 section(s) (expected 5+)

📊 Step 2: Checking pages table...
------------------------------------------------------------
   Pages in database: 17

📥 Step 3: Preparing to import sections data...
------------------------------------------------------------
   Found 1 sections (expected 70+), will replace with correct data
   ✓ Found sections data at: scripts/sections-data.sql
   Found 77 sections to import

💾 Step 4: Importing sections data...
------------------------------------------------------------
   Clearing existing sections...
   ✓ Cleared 1 existing sections
   ✓ Successfully imported 77 sections

✅ Step 5: Verifying import...
------------------------------------------------------------
   Total sections after import: 77
   Sections for page 1: 5

   Sample sections for page 1:
     Section 1: Alice was beginning to get very tired of sitting b...
     Section 2: the book her sister was reading, but it had no pic...
     Section 3: ' So she was considering in her own mind (as well ...
     Section 4: of making a daisy-chain would be worth the trouble...
     Section 5: There was nothing so very remarkable in that; nor ...

🎉 SUCCESS! Fix completed successfully!
   Page 1 now has 5 sections (expected 5+)
   You can now test the Render.com reader app.
```

## 🚀 Next Deployment

The fix will automatically run during the next Render.com deployment. After deployment:

1. The sections data will be automatically fixed
2. Check the Render.com build logs to see the fix script output
3. Test the reader app - Page 1 should show 5 section snippets

## 📁 Files Involved

- **`cmd/fix-render/main.go`** - The automated fix script
- **`scripts/sections-data.sql`** - The sections data to import (77 sections)
- **`render.yaml`** - Updated to run fix-render during build

## 🛠️ Troubleshooting

**If the fix doesn't run automatically:**
- Check Render.com build logs for errors
- Manually run: `DB_PATH=data/alice-suite.db ./bin/fix-render`

**If sections still show incorrectly:**
- Check that migrations ran successfully
- Verify the `pages` table has data (sections need pages to exist)
- Run diagnostic: `DB_PATH=data/alice-suite.db go run cmd/fix-sections/main.go`

