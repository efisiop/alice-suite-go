# Deploy Sections Fix – Steps 1, 2, 3

Follow these in order.

---

## Step 1: Deploy the updated code

The fix is in your repo (updated `internal/handlers/rpc.go` and docs). Deploy by pushing to the branch Render uses.

1. **Commit and push** (from the project root):

   ```bash
   git add internal/handlers/rpc.go FIX_RENDER_SECTIONS_ISSUE.md DEPLOY_SECTIONS_FIX.md
   git status
   git commit -m "Fix sections on Render: Rebind in RPC fallback, single-section fallback when no sections in DB"
   git push origin main
   ```

   If your Render service uses another branch (e.g. `master`), push that branch instead of `main`.

2. **Trigger deploy on Render**
   - If Render is connected to this repo, it will deploy automatically after the push.
   - Otherwise: Render Dashboard → your service (alice-suite-go) → **Manual Deploy** → **Deploy latest commit**.

3. Wait until the deploy finishes (build + start).

---

## Step 2: Ensure sections are loaded on Render

This is already set up:

- **`render.yaml`** builds `bin/fix-render`.
- **`start.sh`** runs `./bin/fix-render` on every start (after migrations and init-users).

You only need to confirm it runs:

1. After the deploy from Step 1, open **Render Dashboard** → **alice-suite-go** → **Logs**.
2. Right after a deploy, look for:
   - `Verifying and fixing sections data...`
   - Then either:
     - `✅ Sections fix completed successfully`, or
     - `✅ Page 1 now has 5 sections (expected 5+)` (from fix-render output)

If you see `⚠️  Warning: fix-render binary not found`, the build did not produce `bin/fix-render`; check the build command in Render (should match `render.yaml`).

---

## Step 3: Check Render logs after deploy

After the service has started, check logs to see whether sections are coming from the DB or from the fallback.

1. **Render Dashboard** → **alice-suite-go** → **Logs**.
2. Load the reader and open a page (e.g. Page 1) so the app calls `get_sections_for_page`.
3. In the logs, look for one of these:

**Success (sections from database):**

```text
Found N sections in database for page 1
Using new structure (page_number)
Successfully found N sections for page 1
```

**Fallback (no sections in DB; full page as one section):**

```text
No sections in DB for page 1; returning full page content as single section. Run fix-render for proper section division.
```

- If you see **Success**: Section division on Render matches local; you’re done.
- If you see **Fallback**: The reader still shows the full chapter as one section. Run fix-render (see below) or re-check that `start.sh` runs `./bin/fix-render` and that the build includes `bin/fix-render`.

**Optional – run fix-render manually on Render**

- Render Dashboard → **alice-suite-go** → **Shell** (if available).
- Run:
  ```bash
  ./bin/fix-render
  ```
- Restart the service or open the reader again and repeat Step 3.

---

## Quick checklist

- [ ] Step 1: Code committed and pushed; Render deploy finished.
- [ ] Step 2: Logs show “Verifying and fixing sections data...” and “Sections fix completed” (or similar).
- [ ] Step 3: Logs show “Successfully found N sections” for page 1 (or run fix-render and re-check).
