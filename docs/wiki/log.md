# Alice Suite wiki log

Chronological timeline of evolution decisions and changes.

Use this header format for every new entry:
`## [YYYY-MM-DD] <type> | <short title>`

---

## [2026-06-06] fix | Restore reader auth-page navbar format

### what changed

- Restored page-specific `nav` blocks on reader login and register pages.
- Removed the base-template rule that hid the navbar collapse on public reader entry/auth pages.
- Updated the login/register cross-link to use canonical `/reader/...` routes.

### why

- The register page no longer matched the previous visible format because the auth-page nav had been removed and then hidden globally from `base.html`.

### files touched

- `internal/templates/reader/login.html`
- `internal/templates/reader/register.html`
- `internal/templates/base.html`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `go build -o bin/server ./cmd/server` succeeded.
- Restarted local server and confirmed `/reader/register` renders `Home`, `Login`, and `Create Account`.

## [2026-06-03] docs | Hermes handover briefing (VPS) with required skills

### what changed

- Added `docs/HERMES_HANDOVER_BRIEFING.md`: a handover document for the Hermes agent runtime on the VPS.
- It captures current state of the three apps (Reader, Consultant, Admin), the codebase map, known gaps/debt, and a full skills inventory (cross-cutting + per-app) needed to raise quality and optimization.
- Added a tracking row in `docs/wiki/index.md`.

### why

- User wants to turn over all work done so far to Hermes (installed on the VPS) and have a clear briefing of the skills needed to bring the apps to a higher standard.

### files touched

- `docs/HERMES_HANDOVER_BRIEFING.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- Briefing cross-references existing repo docs (CODEBASE_SUMMARY, UPGRADE_ROADMAP, ADMINISTRATOR_PROPOSAL, EVOLUTION_CONTROL) and matches observed code (e.g. per-request `template.ParseFiles`, test coverage, admin handler scope).

### next step

- User fills `docs/wiki/DIRECTION.md` north star; Hermes starts with VPS stabilization (embed templates/static, systemd/Docker + TLS, secrets) then the roadmap.

---

## [2026-04-28] docs | add Obsidian vault import script for source files

### what changed

- Added `scripts/import_to_alice_vault.sh` to bulk-import source files into the vault.
- Added `Alice Suite/IMPORT_INSTRUCTIONS.md` with simple copy/paste + one-command steps.
- Import flow now supports `.docx`, `.pdf`, images, `.md`, and `.txt`.

### why

- User asked for an easy way to add Alice Suite creation files into the overall Obsidian vault as notes.

### files touched

- `scripts/import_to_alice_vault.sh`
- `Alice Suite/IMPORT_INSTRUCTIONS.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- Ran `bash scripts/import_to_alice_vault.sh` successfully.
- Script created/checked `Alice Suite/Inbox`, `Alice Suite/Knowledge/Imported`, and `Alice Suite/Knowledge/Attachments`.

### next step

- User copies files into `Alice Suite/Inbox/` and runs the script to generate notes and attachments in the vault.

---

## [2026-04-28] docs | Direction note + Obsidian guide for repo wiki

### what changed

- Added `docs/wiki/DIRECTION.md` — a short, fillable template for north star, non-goals, and sharp priorities.
- Added `docs/wiki/OBSIDIAN_STEPS.md` — numbered steps to open this repo (or `docs/`) in Obsidian and keep Git as the source of truth.
- Updated `docs/wiki/index.md` with new topics for direction and Obsidian.

### why

- User wanted a clear place to steer the product in-repo and to use Obsidian to read and edit the same files as the evolution wiki.

### files touched

- `docs/wiki/DIRECTION.md`, `docs/wiki/OBSIDIAN_STEPS.md`, `docs/wiki/index.md`, `docs/wiki/log.md`

### verification

- Confirmed new files exist and index links to them.

### next step

- Open `DIRECTION.md` in Obsidian (or here) and replace the empty sections in order.

---

## [2026-04-28] feature | Reader reading area: sections left, navigation in Services

### what changed

- **Left column** shows only the Sections card: title "Sections", subheading "Select the section you're reading", and the section list.
- **Services (right column)** includes a new **Navigation tools** accordion: Go to Page, Previous/Next, and Scan to Locate (same behavior and element ids as before).
- `jumpToSections()` scrolls to the section subheading (collapse expand removed with the old left nav).

### why

- User-requested layout: section picker on the left; page/scan tools grouped with Services.

### files touched

- `internal/templates/reader/interaction.html`
- `docs/wiki/index.md`, `docs/wiki/log.md`

### verification

- `go build` on `cmd/server`. Manual: open reader interaction page, confirm left shows only sections; open Services → Navigation tools for page nav and scan.

---

## [2026-04-27] feature | Reader dictionary popup (roadmap 1.1)

### what changed

- Clarified **source** labels: Alice glossary vs dictionary (live) vs dictionary (saved cache).
- **Phonetic** display avoids duplicate slashes when the API already includes them; **part of speech** row uses a small flex meta line.
- **Example** uses a dedicated “Example” label and italic line; single-example modal shows inline; multiple examples still use expand.
- Dictionary requests can include `**section_id`** when the reader is on a section (matches RPC support).

### why

- Implements `docs/UPGRADE_ROADMAP.md` item **1.1** (clearer definition output). Backend in `internal/handlers/rpc.go` already returned `source`, `part_of_speech`, and `phonetic` where applicable; UI now matches.

### files touched

- `internal/templates/reader/interaction.html`
- `docs/wiki/index.md`, `docs/wiki/log.md`

### verification

- `go build -o /dev/null ./cmd/server/` succeeded. Manual: open reading page, tap a word, confirm source line + optional POS/phonetic + example label; use Dictionary button in UI for modal parity.

### next step

- Roadmap **1.2** (reading layout) or polish **1.1** optional step (multiple meanings) if desired.

## [2026-04-25] docs | bootstrap evolution control wiki

### what changed

- Added a simple operating schema for evolution/documentation.
- Added a state index for active topics.
- Added this chronological log file.
- Linked agent behavior to this workflow in `AGENTS.MD`.

### why

- To keep app evolution and documentation controlled in-repo, not only in chat history.
- To follow a lightweight Karpathy-style wiki pattern (`schema + index + log`).

### files touched

- `AGENTS.MD`
- `docs/EVOLUTION_CONTROL.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- Confirmed files exist and are readable in repo.

### next step

- For each future feature/fix/docs task, append one entry here and update `docs/wiki/index.md` status if needed.

## [2026-04-25] docs | set default knowledge-base workflow

### what changed

- Added a beginner-friendly default routine file.
- Added copy-paste templates for index rows and log entries.
- Updated evolution schema to mark knowledge-base-first mode as default.
- Updated index topic statuses.

### why

- User asked to make knowledge-base building the default behavior.
- This reduces friction and keeps documentation consistent.

### files touched

- `docs/EVOLUTION_CONTROL.md`
- `docs/wiki/START_HERE.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- Confirmed all files exist and reference each other correctly.

### next step

- Use this routine on the next approved roadmap item (for example `1.1` in `docs/UPGRADE_ROADMAP.md`).

## [2026-04-25] infra | run refresher protocol

### what changed

- Executed `bash refresher-protocol.sh`.
- Generated a new refresher report file.

### why

- User requested to run the refresher protocol.

### files touched

- `archive/REFRESHER_2026_04_25.md`
- `docs/wiki/log.md`

### verification

- Script finished with success message: "Refresher Protocol Complete!".

### next step

- Review `archive/REFRESHER_2026_04_25.md` if you want to inspect protocol summary details.

## [2026-04-25] feature | add concise reader activity summary panel

### what changed

- Added a compact summary panel in Consultant -> Reader Inspector -> "Reader activity — stats & charts".
- Summary now lists key metrics in plain language: pages clicked, words looked up, AI requests, help requests, active days, events/day, and AI share.
- Wired the panel to existing activity-charts data load and added fallback text for load errors.

### why

- User asked for a quick consultant-facing overview next to graphs, so key behavior can be understood without reading every chart.

### files touched

- `internal/templates/consultant/reader-inspector.html`
- `docs/wiki/index.md`
- `docs/wiki/log.md`
- source docs referenced: `docs/UPGRADE_ROADMAP.md`, `docs/APP_UPGRADES.md`

### verification

- Reviewed template diff to confirm the new panel appears inside "Reader activity — stats & charts".
- Verified JS fills the concise summary when activity data loads and shows a readable fallback on error.

### next step

- Open localhost Reader Inspector and click "Show activity" to validate layout and wording with real reader data.

## [2026-04-25] feature | add cumulative summary on all readers charts

### what changed

- Added a compact "Cumulative summary" panel to Consultant -> Readers -> "All readers — activity stats & charts".
- Summary uses aggregate activity data and lists key cumulative metrics in plain language: pages clicked, words looked up, AI requests, active readers, help requests, events/day, and AI share.
- Connected this panel to the existing all-readers chart load flow with safe fallback text on errors.

### why

- User asked to have the same quick summary style used in Reader Inspector, but for all readers as a cumulative consultant view.

### files touched

- `internal/templates/consultant/readers.html`
- `docs/wiki/log.md`

### verification

- Reviewed template/JS diff to confirm the new panel is visible in the chart row.
- Verified summary is populated from the aggregate API response fields already used by existing stat cards.

### next step

- Open localhost Readers page and click "Show activity" to confirm the summary aligns with chart values.

## [2026-04-25] feature | refresh reader landing page with Alice-style design

### what changed

- Redesigned the reader landing page to feel more "Alice-like" with warmer visual style and clearer copy.
- Added a public-domain cover image for *Alice's Adventures in Wonderland* (1865 edition) to the landing hero section.
- Kept login and account creation actions prominent and simple.

### why

- User requested a more thematic welcome page where login/signup appears, including a copyright-safe book cover image.

### files touched

- `internal/templates/reader/landing.html`
- `internal/static/images/alice-cover-1865.jpg`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- Confirmed image file exists in static assets and is referenced by `/static/images/alice-cover-1865.jpg`.
- Reviewed landing template diff to ensure updated copy, styling, and buttons render in one cohesive layout.

### next step

- Open `/reader` on localhost and verify image load and text readability on desktop and mobile widths.

## [2026-04-25] feature | switch to more colorful landing cover

### what changed

- Replaced the landing hero image with a more colorful public-domain Alice cover.
- Updated the on-page attribution text to match the new source (1907 Charles Robinson edition).

### why

- User requested a cover that feels more alive and colorful.

### files touched

- `internal/static/images/alice-cover-1865.jpg` (replaced image content)
- `internal/templates/reader/landing.html`
- `docs/wiki/log.md`

### verification

- Confirmed replacement image downloaded successfully and landing page still points to the same static path.

## [2026-04-25] feature | use user-provided colorful landing cover

### what changed

- Replaced the landing cover image with the user-provided colorful Alice artwork.
- Updated landing image alt text and attribution note to match the new source.

### why

- User requested this specific, more alive image instead of the previous cover.

### files touched

- `internal/static/images/alice-cover-1865.jpg` (replaced image content)
- `internal/templates/reader/landing.html`
- `docs/wiki/log.md`

### verification

- Confirmed new image was copied into the static path used by landing page.

## [2026-04-26] fix | simplify reader entry header and align landing colors

### what changed

- Removed top header navigation links on reader entry/auth pages (`/reader`, `/reader/login`, `/reader/register`, `/login`, `/register`, `/verify`) to avoid random navigation.
- Made navbar brand non-clickable on those entry/auth pages.
- Retuned landing palette (gold/red/blue accents) to better match the colorful Wonderland image.

### why

- User requested cleaner entry experience without extra header links and closer visual match between page colors and hero image context.

### files touched

- `internal/templates/base.html`
- `internal/templates/reader/landing.html`
- `internal/templates/reader/login.html`
- `internal/templates/reader/register.html`
- `docs/wiki/log.md`

### verification

- Checked base template logic hides collapse navigation for the specified reader entry/auth routes.
- Checked landing CSS now uses warmer/high-contrast colors consistent with the new image.

### update

- Extended the same header-hide rule to the root landing route (`/`) so top-right links are removed there too and brand title is non-clickable.

## [2026-06-12] fix | make reader workspace resize correctly

### what changed

- Added a per-page main container override in the base template.
- Changed `/reader/interaction` from a fixed Bootstrap container to a fluid workspace.
- Updated the reader panels to use the three-column layout only on large screens and stack earlier on medium/narrow windows.
- Made side panels sticky only when there is enough viewport width.

### why

- The reader interaction page is a workspace with three active panels, not a simple content page. The fixed-width shell made it feel incorrectly scaled when the app was opened outside the in-app browser or resized.

### files touched

- `internal/templates/base.html`
- `internal/templates/reader/interaction.html`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `go build -o bin/server ./cmd/server` succeeded.
- Restarted localhost server on port `8080`.
- Confirmed `/health` returns `status: ok`.
- Confirmed `/reader/interaction` serves `container-fluid reader-interaction-main`, `reader-workspace`, and `col-lg-*` responsive layout classes.

## [2026-06-12] feature | add hideable reader side panels

### what changed

- Made the center reading panel more prominent in the desktop layout.
- Added reader layout controls to hide/show the left Sections panel and right Services panel.
- Persisted each reader's panel visibility choice in browser local storage.
- Kept side-panel toggles functional on narrow windows too.

### why

- User requested stronger focus on the central reading window while keeping side menus available when needed.

### files touched

- `internal/templates/reader/interaction.html`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `go build -o bin/server ./cmd/server` succeeded.
- Restarted localhost server on port `8080`.
- Confirmed `/health` returns `status: ok`.
- Confirmed `/reader/interaction` serves the `Hide sections` / `Hide services` controls, the `reader-panel-center` layout, and the `alice-reader-layout` persistence key.

## [2026-06-12] fix | enlarge central reader text

### what changed

- Increased the central book text size again for stronger reader focus.
- Raised the reading line-height to keep the larger typography comfortable.
- Kept the increase scoped to the central reading content rather than side menus or modals.

### why

- User requested the central app characters/text to be even larger so the center gets more attention.

### files touched

- `internal/templates/reader/interaction.html`
- `docs/wiki/log.md`

### verification

- `go build -o bin/server ./cmd/server` succeeded.
- Restarted localhost server on port `8080`.
- Confirmed `/reader/interaction` serves `font-size: 1.375rem`, `line-height: 1.78`, and the mobile reader text override.
- Health check verification was not rerun because the approval request was declined.

## [2026-06-12] fix | move reader panel toggles below book content

### what changed

- Moved the Sections and Services hide/show buttons from the top of the reader workspace to the bottom of the central reading card.
- Added a subtle top border and spacing so the controls read as secondary layout actions below the book content.

### why

- User requested the toggle buttons for Sections and Services to sit below the book content.

### files touched

- `internal/templates/reader/interaction.html`
- `docs/wiki/log.md`

### verification

- `go build -o bin/server ./cmd/server` succeeded.
- Restarted localhost server on port `8080`.
- Confirmed `/reader/interaction` serves `#page-content` before `.reader-layout-toolbar`, with `Hide sections` and `Hide services` directly below the central book content.
