# Administrator Side — Proposal

This document describes the **Administrator** side of the Alice Suite: what it is for, what it will do, and how it fits with the existing Reader and Consultant apps.

---

## 1. The Three Sides of Alice Suite

| Side           | Scale        | Who uses it                    | Purpose                                      |
|----------------|-------------|---------------------------------|----------------------------------------------|
| **Reader**     | Up to millions | End readers (students, etc.)   | Read the book, get help, use the companion app |
| **Consultant** | Up to thousands | Support staff / consultants    | Manage assigned readers, send prompts, handle help requests |
| **Administrator** | **Max ~10–20** | HQ / board / ops managers      | Monitor the whole operation and consultant effectiveness |

Administrators are few, with different levels (e.g. **global** at publisher HQ, only 1–2 users).

---

## 2. What the Administrator Interface Will Do

### 2.1 Monitor “App open” events

- **Reader app open**  
  Know when a reader has the Reader app open (e.g. from existing sessions/activity or explicit “app open” events).
- **Consultant app open**  
  Know when a consultant has the Consultant app open (same idea: sessions or explicit events).

So the admin view can show “readers online now” and “consultants online now” (and optionally simple history).

### 2.2 Supervise consultant effectiveness

- Use a **stats matrix** to see if consultants are doing a good job.
- Examples of metrics (to be refined with you):
  - Response time to help requests.
  - How many readers each consultant is responsible for vs. how many are “followed up” (e.g. had recent contact or activity).
  - Prompts sent vs. accepted/dismissed (already in your DB).
  - Resolution rate of help requests.

This gives a board-level view: “Are consultants on top of their readers?”

### 2.3 Reader follow-up overview

- See whether readers are being followed up (e.g. had consultant interaction, prompt, or help request resolution in a given period).
- High-level counts and trends (e.g. “X% of active readers had follow-up in the last 7 days”).

### 2.4 Board of directors overview

- One place to see:
  - How many readers and consultants are active (e.g. “open app” or “active in last N minutes”).
  - High-level consultant performance (from the stats matrix).
  - Reader follow-up health (from the metrics above).

No need for admins to do day-to-day consultant work; they only **monitor and oversee**.

---

## 3. What We Can Build (Technical Summary)

The following fits your current codebase (Go server, SQLite, JWT, sessions, activity_logs, reader_states, consultant-related tables).

### 3.1 Authentication and access

- Add an **administrator** role (and optionally **admin_level**, e.g. `global` vs `regional`, with only 1–2 global).
- Reuse existing login flow: admin users log in at e.g. `/admin/login`; JWT and session store role so only `administrator` can access admin routes and APIs.
- Middleware: `RequireAdmin` (and optionally `RequireAdminLevel("global")` for the most sensitive views).

### 3.2 “App open” and presence

- **Readers:** You already have sessions and activity (e.g. LOGIN, page views). We can treat “has active session in last N minutes” as “Reader app open.” If you want an explicit “app open” event, we can add a small activity type (e.g. `APP_OPEN`) when the Reader app loads.
- **Consultants:** Same idea: use existing consultant sessions for “Consultant app open,” and optionally add an explicit event when the Consultant dashboard loads.
- Admin APIs can then expose:
  - Readers online now (and optionally last 24h).
  - Consultants online now (and optionally last 24h).

### 3.3 Stats matrix (consultant effectiveness)

- New **admin-only APIs** that aggregate existing data:
  - From `sessions`, `activity_logs`, `reader_states`, `help_requests`, `consultant_prompts` (and any reader–consultant assignment table).
- Example metrics we can derive:
  - Per consultant: readers assigned, help requests received, resolved, avg time to resolve.
  - Per consultant: prompts sent, accepted vs dismissed.
  - Per consultant: share of assigned readers with recent activity/follow-up.
- We can add a **stats matrix** table later if you want to cache daily/hourly aggregates; for “what we can do now,” we can compute from existing tables and then optimize with caching if needed.

### 3.4 Reader follow-up overview

- APIs that answer:
  - How many readers were “active” in a period (e.g. had at least one session or activity).
  - How many of those had “follow-up” (e.g. help request resolved, or consultant prompt accepted, or similar).
- Again, computed from existing tables; no new app-side concept, just definitions of “active” and “follow-up” we agree on.

### 3.5 Administrator UI (dashboard)

- **Admin app** under `/admin`:
  - `/admin/login` — login (administrator role only).
  - `/admin` or `/admin/dashboard` — main dashboard.
- Dashboard content (we can build step by step):
  1. **Live presence:** Readers online now, Consultants online now (and optionally simple charts over last 24h).
  2. **Consultant stats matrix:** Table or cards per consultant with the metrics above.
  3. **Reader follow-up:** e.g. “X% of active readers had follow-up in last 7 days” and a simple trend.
  4. **Board overview:** Summary tiles (total readers, total consultants, key health metrics).

All of this is **read-only** for the board; no editing of readers or consultants from this view unless you later ask for it.

---

## 4. What Exists Already (What We Reuse)

- **Sessions:** `sessions` table with `user_id`, `last_active_at`, `expires_at` — we can derive “reader/consultant app open” from active sessions per role.
- **Activity:** `activity_logs` (e.g. LOGIN, PAGE_VIEW, etc.) and `reader_states` — we can derive reader activity and “last seen.”
- **Consultant–reader link:** `help_requests.assigned_to`, consultant prompts, and any assignment table — we can attribute follow-up and workload to consultants.
- **Auth:** Same JWT + cookie + DB sessions; we only add role `administrator` and admin-level checks.

So the Administrator side is mostly **new routes, new role, new dashboard and APIs** that **read** from your current schema.

---

## 5. Next Steps (We Can Get Into Details Along the Way)

1. **Confirm scope**  
   Agree on: “app open” definition (session-based vs explicit event), exact consultant metrics, and definition of “reader follow-up.”
2. **Database**  
   One migration: add `administrator` (and optionally `admin_level`) to `users` and create 1–2 admin users.
3. **Backend**  
   - Add `RequireAdmin` (and optionally `RequireAdminLevel`) middleware.  
   - Expose admin-only APIs: presence, consultant stats matrix, reader follow-up, board summary.
4. **Frontend**  
   - Admin login page and dashboard (read-only).  
   - Widgets for: live presence, consultant matrix, reader follow-up, board overview.

If you tell me your priority (e.g. “first just reader/consultant online counts and one summary page”), we can do that first and then add the stats matrix and follow-up metrics in the next steps.

---

## 6. Creating the first administrator user

After running migration `013_add_administrator_role.sql`, you need at least one user with `role = 'administrator'`. Options:

- **Manual (SQL):**  
  `INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at) VALUES (..., 'admin@example.com', '<bcrypt-hash>', 'Admin', 'User', 'administrator', 1, datetime('now'), datetime('now'));`
- **Extend `cmd/init-users`:**  
  Add a step that creates an admin user (e.g. `admin@example.com` / `admin123`) if missing, same way as the existing reader and consultant seeds.
- **One-off script:**  
  A small program that hashes a password and inserts one administrator row.

You can add an optional `admin_level` (e.g. `'global'` or `'regional'`) on that row later if you use that column.
