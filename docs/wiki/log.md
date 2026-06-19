# Alice Suite wiki log

Chronological timeline of evolution decisions and changes.

Use this header format for every new entry:
`## [YYYY-MM-DD] <type> | <short title>`

---

## [2026-06-19] feature | Italian reader interface demo

### what changed

- Added a Reader UI translation layer for Italian that uses the saved reader language preference.
- Translated the main reader navigation, dashboard, My Page/account settings, and visible reading/AI/dictionary controls for the demo.
- Kept the translator away from book/page text containers so the physical book source text is not rewritten.
- Added preferred language to login responses and browser session storage so the UI can switch after login.

### why

- User confirmed AI output was Italian but the overall Reader interface still looked English.

### files touched

- `internal/static/js/app.js`
- `internal/handlers/auth.go`
- `internal/templates/reader/login.html`
- `internal/templates/reader/my-page.html`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `node --check internal/static/js/app.js` passed.

## [2026-06-19] feature | reader help language preference

### what changed

- Added `reader_preferences` storage for a reader's selected help/output language.
- Added a `Help Language` selector during reader registration.
- Added Account Settings on My Page so readers can change their help language after signup.
- Added `/api/reader/preferences` for authenticated preference load/save.
- Updated AI assistant prompt construction so answers use the saved help language while preserving book quotes, titles, and character names unless translation is requested.

### why

- User requested language selection at signup and under account settings, while keeping a clear boundary between the physical book edition language and the app's help/output language.

### files touched

- `migrations/017_reader_preferences.sql`
- `internal/models/models.go`
- `internal/database/database.go`
- `internal/database/reader_preferences.go`
- `pkg/auth/auth.go`
- `internal/handlers/auth.go`
- `internal/handlers/helpers.go`
- `internal/handlers/api.go`
- `internal/services/ai_service.go`
- `internal/templates/reader/register.html`
- `internal/templates/reader/my-page.html`
- `archive/old-code/handlers-stub.go`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `go build ./cmd/server` passed.
- `go build ./cmd/migrate` passed.
- `go test ./internal/database ./pkg/auth ./cmd/server ./cmd/migrate` passed.
- `go test ./...` was run; it is still blocked by existing failures in `internal/services/book_service_test.go:300`, `internal/handlers/api_test.go`, and `internal/middleware/auth_test.go`.

## [2026-06-17] feature | dictionary picture action

### what changed

- Added a `Picture` action to the word-click dictionary popup.
- Added a picture panel that first tries a fast Wikipedia summary thumbnail for the looked-up term.
- Added fallback generation through Alice's existing `/api/ai/generate-image` and `/api/ai/image-status` endpoints.
- Added evaluator markers for the picture button, Wikipedia lookup, generated fallback, and picture panel container.
- Marked Reader fix list item `R-021` done.

### why

- User requested a graphical representation of the looked-up word or definition inside the dictionary popup.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.

## [2026-06-19] infra | publish dictionary Picture action for Render

### what changed

- Published the dictionary Picture action to `main` in commit `8a15f01`.
- Render is configured to auto-deploy `main` through `render.yaml`.

### why

- The Picture control existed only in the working tree, while the deployed Render revision could not include it.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- The GitHub push to `main` completed.
- Render's public endpoint timed out from this environment, so final live-page verification remains pending.

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

## [2026-06-16] fix | tighten dictionary derivation and examples

### what changed

- Added a concise derivation cleanup/summarizer before showing Wiktionary etymology in the dictionary popup.
- Replaced generated example fallbacks that framed the looked-up term as a vocabulary word with natural usage-context sentences.
- Added evaluator markers for the derivation summarizer and regression checks for `I used...` / `The word...` example phrasing.
- Marked Reader fix list item `R-020` done.

### why

- User requested derivation show only the main origin, not a long history tree, and examples should be direct usage sentences instead of repeated word-awareness phrases.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.

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

## [2026-06-12] docs | add Reader autoresearch workflow

### what changed

- Added an Alice-specific Reader autoresearch workflow inspired by Karpathy's `autoresearch` loop.
- Added a fixed evaluator script for Reader layout/build checks with optional live localhost/Render checks.
- Added the workflow to the wiki index.

### why

- User asked to adapt Karpathy's autoresearch idea as a source of testing and improvement hints for the Reader app.

### files touched

- `docs/READER_AUTORESEARCH.md`
- `scripts/reader_autoresearch_check.sh`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally: Go build and Reader template marker checks.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.
- `BASE_URL=https://alice-suite-go.onrender.com scripts/reader_autoresearch_check.sh` passed against Render: live health and Reader interaction markers.

## [2026-06-12] docs | start Reader fix list

### what changed

- Added a Reader fix list with the first two user-reported issues.
- Marked dictionary examples as the active item.
- Added AI Assistant flow simplification/testing as the next item.

### why

- User wants unsequenced Reader interaction issues turned into a tracked list, then handled one by one and marked done.

### files touched

- `docs/READER_FIX_LIST.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- Pending implementation and Reader evaluator check.

## [2026-06-12] fix | simplify dictionary examples

### what changed

- Removed inline display of stored dictionary/glossary examples from the word-click dictionary popup.
- Changed the dictionary action label from `Examples` to `Example`.
- Added generated everyday examples for the popup and manual Dictionary modal.
- Marked Reader fix list item `R-001` done.

### why

- User reported the dictionary output felt clumsy and did not want Alice book excerpts as examples. The desired behavior is a short, simple life-usage example.

### files touched

- `internal/templates/reader/interaction.html`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.

## [2026-06-13] fix | stabilize Sections sidebar height

### what changed

- Added a stable desktop minimum height to the left Sections card.
- Made the Sections card body a vertical flex layout.
- Gave the section snippets list its own minimum height and scroll area so it can show more entries without shrinking the whole box when only a few sections are present.
- Kept mobile layout flexible by resetting the snippet minimum height at narrow widths.
- Added evaluator checks for the stable Sections sidebar sizing.
- Marked Reader fix list item `R-015` done.

### why

- User requested the left sidebar stay static and cover a larger amount of sections, while not collapsing when only two or three sections are shown.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.

## [2026-06-13] fix | vary dictionary examples on repeat clicks

### what changed

- Changed generated dictionary examples from one fixed pair into multiple varied example pairs.
- Updated the word-click dictionary popup so each `Example` click advances to a different pair instead of only hiding the panel.
- Kept the small `x` close control as the way to hide examples.
- Added more varied meaning-aware chair examples, including direct-use and question-style sentences.
- Added evaluator checks for varied example sets and cycling state.
- Marked Reader fix list item `R-014` done.

### why

- User said the dictionary examples are much better but the two examples should not be too similar, and clicking the Example button again should present a totally different set.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.

## [2026-06-12] fix | align Reader service tools with dictionary popup style

### what changed

- Added a shared `reader-service-popup` visual shell for Reader service modals.
- Applied it to Scan to Locate, Test your knowledge, Ah Ah Moments, Human Consultant, and manual Dictionary.
- Centered the service modals so they read more like the prominent dictionary popup.
- Matched the dictionary popup's restrained border, 8px radius, strong shadow, pale header, Lora title, quieter controls, and light focus surfaces.
- Added evaluator checks that all service modals carry the shared dictionary-style popup class.
- Marked Reader fix list item `R-013` done.

### why

- User requested all service tools such as Human Consultant, Quiz, Dictionary, and Ah Ah Moments visually move closer to the dictionary popup design type.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.

## [2026-06-12] fix | reorder Services and improve dictionary examples

### what changed

- Reordered the Reader Services sidebar so Navigation tools appears below Info Center.
- Passed dictionary definitions into the generated example helper.
- Added a definition-context helper so examples can refer to the meaning of the word, not just the word string.
- Added a specific seating-context path for entries such as `chair`.
- Removed the generic `I heard the word...` fallback from generated examples.
- Added evaluator checks for Services order and definition-aware dictionary examples.
- Marked Reader fix list items `R-011` and `R-012` done.

### why

- User requested Navigation tools move below Info Center and asked that dictionary examples include a little context about the actual meaning of the word.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.

## [2026-06-12] fix | align AI assistant popup with dictionary popup

### what changed

- Restyled the floating AI assistant to use the dictionary popup's visual language: 8px radius, pale draggable header, restrained border, stronger shadow, white content surface, and quieter controls.
- Removed the WhatsApp-style green header, patterned message area, and green send button styling.
- Made AI mode tabs, quick actions, selected-text controls, and chat bubbles visually closer to the dictionary popup's secondary action style.
- Added evaluator markers for the dictionary-style AI assistant design.
- Marked Reader fix list item `R-010` done.

### why

- User requested the AI assistant popup look very close to the dictionary popup as a design type.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.
- Confirmed `/reader/interaction` serves `getGeneratedEverydayExamples` and `Generated everyday examples`.

## [2026-06-12] fix | strengthen dictionary popup hierarchy

### what changed

- Made the word entry and definition the dominant dictionary popup content.
- Reduced the visual weight of the popup chrome and Derivation/Example buttons.
- Moved source attribution to the bottom of both dictionary popup and manual Dictionary modal output.
- Changed the word-click dictionary popup to open more centrally and prominently.
- Added evaluator markers for the more prominent dictionary popup.
- Marked Reader fix list item `R-003` done.

### why

- User requested dictionary output that reads as the main thing when called, not side information: clear entry, applicable definition, smaller controls, source at the bottom, and a more central popup.

### files touched

- `internal/templates/reader/interaction.html`
- `docs/READER_FIX_LIST.md`
- `scripts/reader_autoresearch_check.sh`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.
- Confirmed the served Reader page includes `dictionary-popup-actions`, centered popup placement, and prominent popup width markers.

### update

- Added a subtle light focus background behind the dictionary entry word and main definition.

## [2026-06-12] fix | add close control to dictionary panels

### what changed

- Added a small close `x` to expanded Derivation and Example dictionary panels.
- Added a shared helper to close the panel and reset the related button styling.
- Added evaluator markers for the dictionary close control.
- Marked Reader fix list item `R-004` done.

### why

- User requested a direct, reader-visible way to toggle out expanded Derivation or Example content.

### files touched

- `internal/templates/reader/interaction.html`
- `docs/READER_FIX_LIST.md`
- `scripts/reader_autoresearch_check.sh`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the restarted local app.
- Confirmed evaluator checks for `dictionary-popup-panel-close` and `closeDictionaryPanel`.

## [2026-06-12] fix | add close controls to Reader sidebars

### what changed

- Added small close `x` controls to the Sections and Services side panel headers.
- Wired the header close controls to the existing Reader panel hide/show state.
- Added evaluator markers for the sidebar close controls.
- Marked Reader fix list item `R-005` done.

### why

- User requested the same direct close affordance on the left and right sidebars.

### files touched

- `internal/templates/reader/interaction.html`
- `docs/READER_FIX_LIST.md`
- `scripts/reader_autoresearch_check.sh`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the restarted local app.
- Confirmed evaluator checks for `reader-panel-close` and `hideReaderPanelFromHeader`.

## [2026-06-12] fix | move page arrows under Sections

### what changed

- Moved Previous/Next page navigation from Services → Navigation tools to a compact arrow row under the Sections list.
- Left Go to Page and Scan to Locate inside Services → Navigation tools.
- Added evaluator markers for the section page navigation arrow row.
- Marked Reader fix list item `R-006` done.

### why

- User requested only the page backward/forward arrows be visibly available under Sections, while other navigation tools remain in Services.

### files touched

- `internal/templates/reader/interaction.html`
- `docs/READER_FIX_LIST.md`
- `scripts/reader_autoresearch_check.sh`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the restarted local app.
- Confirmed evaluator checks for `section-page-nav`.

### update

- Added a small `Page` label beside the section-page arrow controls.
- Added the current page number to the section-page arrow label and synced it from `loadPage()`.

## [2026-06-12] docs | inventory AI Assistant flows

### what changed

- Added an AI Assistant flow inventory for the Reader app.
- Listed current AI Assistant entry points, selection paths, conflicts, intended simplification direction, and a test checklist.
- Marked Reader fix list item `R-002` as `doing`.

### why

- User reported the AI Assistant interaction is convoluted, especially moving between text selection, the small AI window, and the book.

### files touched

- `docs/AI_ASSISTANT_FLOW_INVENTORY.md`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- Documentation-only inventory; implementation checks pending with the first AI flow cleanup.

## [2026-06-12] fix | make AI selected-text transfer explicit

### what changed

- Changed the AI chat selection button copy to `Select text from the book`.
- Removed the visible `Press Enter` instruction from the AI selected-text status.
- Added an explicit `Add selected text` button that appears when selected book text is ready for the assistant input.
- Added evaluator markers for the explicit AI selected-text action.
- Marked Reader fix list item `R-008` done while keeping broader AI Assistant flow item `R-002` in progress.

### why

- User reported that AI Assistant interaction is convoluted when moving between text selection, the floating chat window, and the book. The first simplification removes an Enter-key instruction that competes with sending chat messages.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/AI_ASSISTANT_FLOW_INVENTORY.md`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.

## [2026-06-12] fix | route AI selection spark to floating chat

### what changed

- Changed `openAIHelpWithSelection()` to open the existing floating AI chat instead of the removed legacy Bootstrap modal.
- Preserved the selected-text prefill as `Help me understand this: "..."`.
- Removed the active modal/backdrop retry loop from the current chat selection mode.
- Made Enter in the AI input consistently send the chat message; selected text transfer now uses `Add selected text`.
- Added evaluator checks that the legacy selected-text modal path does not reappear.
- Marked Reader fix list item `R-009` done while keeping broader AI Assistant flow item `R-002` in progress.

### why

- User reported conflicts moving between the book, selected text, and the AI window. The text-selection spark and Services AI Help should land in one assistant surface rather than split between old modal code and the current floating chat.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/AI_ASSISTANT_FLOW_INVENTORY.md`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.

## [2026-06-13] feature | prototype visible AI reading context

### what changed

- Added a visible AI Assistant `Linked to` context bar.
- Added reader-controlled scopes for `This section`, `This page`, and `Selected text`.
- Added an `Ask about this` action that prefills the chat with the chosen reading context.
- Updated the AI request context so the selected scope is sent to the backend.
- Added evaluator markers and updated the Reader fix list and AI Assistant flow inventory.

### why

- The reader should be able to connect the assistant to what they are reading without fiddling with text selection for every question.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/AI_ASSISTANT_FLOW_INVENTORY.md`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the restarted local app.

## [2026-06-14] fix | simplify AI popup around context tools

### what changed

- Removed the visible AI mode tabs from the floating AI Assistant popup.
- Removed the visible old select-text button and quick action buttons from the popup.
- Made the context-linking block the dominant visible control: `This section`, `This page`, `Selected text`, and `Ask about this`.
- Changed `Ask about this` to send the linked context directly when readable text is available.
- Updated evaluator markers, Reader fix list, AI Assistant flow inventory, and wiki index.

### why

- User liked the context-linking concept but found the popup confusing because too many older tools competed for attention.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/AI_ASSISTANT_FLOW_INVENTORY.md`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the restarted local app.

## [2026-06-14] fix | make AI custom selection visible

### what changed

- Added a clearer selection-mode cue when the reader chooses `Selected text` in the AI popup.
- Added a subtle outline around the reading area while AI custom selection is active.
- Changed AI custom selection to persistently highlight the selected passage after mouseup/touchend.
- Updated the AI context preview from the selected passage so the reader can see what the AI is linked to.
- Added evaluator markers and updated the Reader fix list and AI Assistant flow inventory.

### why

- User reported that custom selection gave no visual response in the text, making it unclear what was actually selected.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/AI_ASSISTANT_FLOW_INVENTORY.md`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the restarted local app.

## [2026-06-14] fix | move reader panel toggles outside book card

### what changed

- Moved the Sections and Services hide/show buttons out of the central book card.
- Positioned the button row as fixed page chrome at the bottom-right of the reader background.
- Indented the fixed button row farther left from the right edge.
- Overrode Bootstrap row-child sizing so the fixed control shrink-wraps the buttons instead of stretching full width.

### why

- User requested the hide/show controls sit outside the book section box, fixed down on the bottom-right background area.

### files touched

- `internal/templates/reader/interaction.html`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally: Go build and Reader template marker checks.

## [2026-06-14] fix | make AI assistant popup more ergonomic

### what changed

- Polished the floating AI Assistant header, context card, message area, input, and send button sizing.
- Replaced the instructional empty state with direct linked-text starter actions: explain what is happening, make this easier to read, and what should I notice.
- Added `runAIPromptForCurrentContext()` so starter actions use the current section/page/selection context and send through the existing AI flow.
- Updated the Reader evaluator markers for the new starter-action empty state.
- Updated the Reader fix list, AI Assistant flow inventory, and wiki index.

### why

- User wants the AI Assistant popup to feel neat, ergonomic, and user-friendly rather than technical or cluttered.

### files touched

- `internal/templates/reader/interaction.html`
- `scripts/reader_autoresearch_check.sh`
- `docs/AI_ASSISTANT_FLOW_INVENTORY.md`
- `docs/READER_FIX_LIST.md`
- `docs/wiki/index.md`
- `docs/wiki/log.md`

### verification

- `scripts/reader_autoresearch_check.sh` passed locally: Go build and Reader template marker checks.
- `BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh` passed against the running local app.
