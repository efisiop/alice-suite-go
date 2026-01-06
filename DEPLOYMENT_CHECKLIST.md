# Deployment Checklist: Localhost vs Render.com

## ⚠️ Important: Render.com Database Persistence

**Render.com Free Tier uses ephemeral filesystem** - the database may be reset on each deploy or service restart. This means:
- Data is lost when the service restarts or redeploys
- Migrations must run on every start
- Seed data must be re-inserted on every start

**Solutions:**
1. **Use Render Persistent Disk** (if available on your plan)
2. **Use External Database** (PostgreSQL, etc.)
3. **Accept ephemeral nature** and ensure all data is seeded on startup

## ✅ Pre-Deployment Checklist

### 1. Database Migrations
- [x] All migration files are in `migrations/` directory
- [x] Migrations are numbered sequentially (001, 002, 003, etc.)
- [x] Migration 008_seed_glossary_terms.sql includes all 1,818 glossary terms
- [x] All migrations use `INSERT OR IGNORE` to be idempotent

### 2. Build Process
- [x] `render.yaml` builds all required binaries:
  - `bin/server`
  - `bin/migrate`
  - `bin/init-users`
  - `bin/fix-render`
- [x] `start.sh` runs migrations, init-users, and fix-render on every start
- [x] Environment variables are set correctly

### 3. Data Seeding
- [x] Migration 002_seed_first_3_chapters.sql seeds chapters and sections
- [x] Migration 008_seed_glossary_terms.sql seeds all glossary terms
- [x] `init-users` creates required user accounts
- [x] `fix-render` ensures sections data is correct

### 4. Configuration
- [x] `DB_PATH` is set consistently (data/alice-suite.db)
- [x] `PORT` is set correctly (10000 for Render)
- [x] `ENV=production` is set for Render
- [x] `JWT_SECRET` is generated/configured

## 🔍 Verification Steps

### After Deployment to Render.com

1. **Check Database State:**
   ```bash
   # In Render Shell
   sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM alice_glossary;"
   # Should return ~1,818 (or 1,821 including the 3 from migration 002)
   ```

2. **Check Sections:**
   ```bash
   sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM sections WHERE page_number = 1;"
   # Should return 5 or more
   ```

3. **Check Users:**
   ```bash
   sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM users WHERE role = 'reader';"
   # Should return at least the seeded users
   ```

4. **Check Migrations:**
   ```bash
   sqlite3 data/alice-suite.db "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;"
   # Should show all expected tables
   ```

## 🐛 Common Discrepancies & Fixes

### Issue 1: Missing Glossary Terms
**Symptom:** Glossary definitions missing on Render but present on localhost

**Fix:** 
- Ensure migration 008_seed_glossary_terms.sql exists and runs
- Check that `start.sh` runs migrations on every start
- Verify migration file is included in deployment

### Issue 2: Missing Sections
**Symptom:** Only 1 section per page on Render

**Fix:**
- Ensure `fix-render` runs on every start (in `start.sh`)
- Check that `scripts/sections-data.sql` is accessible
- Verify `fix-render` binary is built and included

### Issue 3: Missing Users
**Symptom:** Users don't exist on Render

**Fix:**
- Ensure `init-users` runs on every start (in `start.sh`)
- Check that user creation logic is idempotent (uses INSERT OR IGNORE)

### Issue 4: Database Resets on Deploy
**Symptom:** All data lost after deployment

**Cause:** Render.com ephemeral filesystem

**Solutions:**
1. Use Render Persistent Disk (upgrade plan)
2. Use external database (PostgreSQL)
3. Accept ephemeral nature and ensure all data seeds on startup

## 📋 Deployment Process

### Current Setup (Render.com)

1. **Build Phase:**
   - Builds all binaries
   - Does NOT run migrations (moved to start phase)

2. **Start Phase (`start.sh`):**
   - Creates data directory
   - Runs migrations (creates schema + seeds data)
   - Runs init-users (creates user accounts)
   - Runs fix-render (ensures sections data)
   - Starts server

### Why This Approach?

- **Migrations in start phase:** Ensures database is always up-to-date, even if it was reset
- **Idempotent operations:** All INSERTs use `INSERT OR IGNORE` so they're safe to run multiple times
- **Automatic recovery:** If database is reset, everything is automatically recreated

## 🔄 Making Changes

When adding new data or migrations:

1. **Create new migration file** (e.g., `009_new_feature.sql`)
2. **Use `INSERT OR IGNORE`** for all data inserts
3. **Test locally** first
4. **Deploy to Render** - migration will run automatically on next start
5. **Verify** using the checklist above

## 📊 Monitoring

Check Render.com logs for:
- Migration execution messages
- "All migrations completed successfully!"
- "Seeding completed!"
- Any errors during startup

## 🎯 Goal: Perfect Match

To achieve perfect match between localhost and Render:

1. ✅ All migrations run on every start
2. ✅ All seed data is included in migrations
3. ✅ All initialization scripts run on every start
4. ✅ Use idempotent operations (INSERT OR IGNORE)
5. ✅ Verify data after each deployment

