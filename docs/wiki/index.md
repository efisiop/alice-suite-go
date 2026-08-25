# Alice Suite wiki index

State-oriented map of active evolution topics.

Update this file when the current state changes.

---

## Active topics


| Topic                          | Status      | Owner        | Main files                                                                                                                                                                             |
| ------------------------------ | ----------- | ------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Product direction (north star) | in-progress | user         | `docs/wiki/DIRECTION.md`                                                                                                                                                               |
| Obsidian as editor for repo docs | done      | user + agent | `docs/wiki/OBSIDIAN_STEPS.md`                                                                                                                                                          |
| Obsidian vault file import workflow | done   | user + agent | `Alice Suite/IMPORT_INSTRUCTIONS.md`, `scripts/import_to_alice_vault.sh`                                                                                                              |
| Evolution control workflow     | done        | user + agent | `docs/EVOLUTION_CONTROL.md`, `AGENTS.MD`, `docs/wiki/log.md`                                                                                                                           |
| Default knowledge base routine | done        | user + agent | `docs/wiki/START_HERE.md`, `docs/EVOLUTION_CONTROL.md`, `docs/wiki/log.md`                                                                                                             |
| Reader autoresearch workflow   | done        | user + agent | `docs/READER_AUTORESEARCH.md`, `scripts/reader_autoresearch_check.sh`, `internal/templates/reader/interaction.html`                                                                     |
| Reader fix list                | in-progress | user + agent | `docs/READER_FIX_LIST.md`, `docs/READER_AUTORESEARCH.md`, `internal/templates/reader/interaction.html` (concise dictionary derivation, natural generated examples, dictionary Picture action; published to `main` for Render deployment) |
| Reader UX upgrades             | in-progress | user + agent | `docs/UPGRADE_ROADMAP.md` item 1.1, `docs/APP_UPGRADES.md`, `docs/AI_ASSISTANT_FLOW_INVENTORY.md`, `internal/templates/reader/interaction.html` (dictionary + reading layout: prominent center reader, hideable Sections/Services panels, Navigation tools under Services, concise dictionary derivation/examples, dictionary Picture action, responsive full-width workspace), `internal/templates/reader/login.html`, `internal/templates/reader/register.html`, `internal/templates/base.html` (auth-page navbar format restored; per-page main container override) |
| Reader help language preference | done       | user + agent | `migrations/017_reader_preferences.sql`, `internal/database/reader_preferences.go`, `internal/handlers/api.go`, `internal/handlers/auth.go`, `internal/services/ai_service.go`, `internal/static/js/reader-i18n.js`, `internal/static/js/app.js`, `internal/templates/reader/login.html`, `internal/templates/reader/register.html`, `internal/templates/reader/my-page.html` |
| Reader localization performance | done       | user + agent | `internal/static/js/reader-i18n.js` (single initial translation pass; incremental mutation handling only) |
| Reader-to-Consultant activity handoff | done | user + agent | `CONTEXT.md`, `docs/HERMES_HANDOVER_BRIEFING.md`, `internal/handlers/activity.go`, `internal/handlers/auth.go`, `internal/templates/reader/interaction.html`, `internal/templates/consultant/dashboard.html`, `internal/handlers/activity_test.go` |
| Consultant signal autoresearch | active | user + agent | `docs/CONSULTANT_SIGNAL_AUTORESEARCH.md`, `scripts/consultant_signal_autoresearch_check.sh`, `docs/AGENT_LOOP_INVENTORY.md`, `docs/AI_ASSISTANT_FLOW_INVENTORY.md`, `internal/templates/reader/interaction.html`, `internal/templates/consultant/dashboard.html` |
| Consultant UX upgrades         | in-progress | user + agent | `docs/UPGRADE_ROADMAP.md`, `docs/AGENT_LOOP_INVENTORY.md`, `cmd/server/main.go`, `migrations/018_consultant_request_shares.sql`, `internal/database/consultant_assignments.go`, `internal/handlers/api.go`, `internal/services/help_service.go`, `internal/templates/consultant/help-requests.html`, `internal/templates/consultant/dashboard.html`, `internal/templates/consultant/readers.html`, `internal/templates/consultant/reader-inspector.html` |
| AI assistance UX/readability   | in-progress | user + agent | `docs/UPGRADE_ROADMAP.md`, `docs/AI_ASSISTANT_FLOW_INVENTORY.md`, `docs/READER_FIX_LIST.md`, `internal/templates/reader/interaction.html` (clean AI context popup: This section, This page, Selected text, Ask about this, linked-text starter actions)                                             |
| Reader AI context module design | done | user + agent | `CONTEXT.md`, `docs/READER_AI_CONTEXT_MODULE_DESIGN.md`, `docs/AI_ASSISTANT_FLOW_INVENTORY.md`, `internal/templates/reader/interaction.html`, `internal/handlers/api.go`, `internal/services/ai_service.go` |
| Agent-loop inventory | done | user + agent | `docs/AGENT_LOOP_INVENTORY.md`, `docs/HERMES_HANDOVER_BRIEFING.md`, `README.md`, `docs/UPGRADE_ROADMAP.md`, `docs/AI_ASSISTANT_FLOW_INVENTORY.md`, `docs/READER_AI_CONTEXT_MODULE_DESIGN.md` |
| GBrain reference integration | done | user + agent | `docs/GBRAIN_ALICE_INTEGRATION.md`, `.gitmodules`, `archive/reference/gbrain` |
| Handover to Hermes (VPS)       | in-progress | user + agent | `docs/HERMES_HANDOVER_BRIEFING.md`, `ADMINISTRATOR_PROPOSAL.md`, `docs/UPGRADE_ROADMAP.md`                                                                                            |


---

## How to use

- check this file first for current state
- then open the linked docs for details
- add/remove topics as scope changes
