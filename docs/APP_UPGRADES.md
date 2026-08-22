# App upgrades log (Alice Suite)

Short record of improvements made to the three apps (Reader, Consultant, Admin) so future work and agents know what was done.

---

## Reader onboarding and book verification (2026-08-22)

- **Physical-book purpose:** The landing page now says clearly that Alice Suite accompanies a physical copy and does not replace it.
- **Registration hand-off:** A successful registration now establishes the reader session before opening the book-code page.
- **Book-code verification:** PostgreSQL-compatible queries, an atomic code claim, and the real schema are used so an invalid code returns a clear response and a valid code cannot be consumed by a partial failure.
- **Access gate:** Signed-out visitors are guided to Reader Login. Signed-in readers must verify their physical book before reader pages and reader-only APIs open.
- **Regression coverage:** Production-style verification, onboarding session, physical-book messaging, login redirects, and the verification gate are covered by automated tests.

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

*Last updated: 2026-08-22. Re-index with `qmd update` and `qmd embed` if using qmd.*
