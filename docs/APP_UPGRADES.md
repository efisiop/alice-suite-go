# App upgrades log (Alice Suite)

Short record of improvements made to the three apps (Reader, Consultant, Admin) so future work and agents know what was done.

---

## Cross-app: Login pages (JavaScript required message)

- **Consultant** `internal/templates/consultant/login.html`: Added `<noscript>` block with message “JavaScript is disabled. Please enable JavaScript to use the login form.” (same as Reader).
- **Admin** `internal/templates/admin/login.html`: Same `<noscript>` block added.
- **Reader** already had it. All three login pages now show a consistent message when JavaScript is off.

---

## Reader app

- **Canonical URLs:** All Reader entry points and auth redirects now use `/reader/login` and `/reader/register` instead of `/login` and `/register`.
- **Files changed:** `landing.html`, `register.html`, `dashboard.html`, `interaction.html`, `my-page.html`, `statistics.html`, `verify.html`.
- **Effect:** One fewer redirect when not logged in; consistent Reader URLs across the app.

---

## Consultant app

- **Navigation:** “Readers” link added to the main nav on every Consultant page, and active state fixed on the Readers page.
- **Files changed:** `help-requests.html`, `readers.html`, `reader-inspector.html`.
- **Nav order (all pages):** Dashboard → Help Requests → Readers → Logout. On the Readers page, “Readers” is marked active (not Dashboard).

---

## Admin app

- **Logout behavior:** Dashboard “Logout” now performs a real logout: clears `auth_token` from `sessionStorage` and cookie, then redirects to `/admin/login`.
- **File changed:** `internal/templates/admin/dashboard.html` (logout link + script).
- **Before:** Logout only linked to `/admin/login`; token stayed set so admin remained logged in. **After:** Clicking Logout clears the session; next visit to `/admin` requires login again.

---

*Last updated: 2025-03 (session). Re-index with `qmd update` and `qmd embed` if using qmd.*
