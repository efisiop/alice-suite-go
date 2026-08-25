# Consultant signal autoresearch

Status: active. Owner: user + agent. Created: 2026-06-26.

## Goal

Decide which Reader inputs should be visible to Consultants, at what urgency, and in what form.

The dashboard should not show every click. It should surface events that help a Consultant answer:

- Is this reader actively reading?
- Is this reader stuck or asking for help?
- Did a Consultant prompt produce a useful reader action?
- Is the Reader app itself blocking progress?

## Operating rule

Every Reader input gets one of four signal levels:

| Level | Dashboard behavior | Use when |
| --- | --- | --- |
| L0 audit-only | Store for history/analytics, never live-card by default | Low-value mechanical events |
| L1 quiet timeline | Show in reader inspector or collapsed history | Useful context, not urgent |
| L2 live card | Add to Consultant dashboard reader card | Current activity or learning friction |
| L3 alert/action | Highlight and route to Help Requests or follow-up queue | Human decision may be needed |

## Current tracked live path

Source: `internal/templates/reader/interaction.html`, `internal/handlers/activity.go`, `internal/templates/consultant/dashboard.html`.

| Reader input | Current event | Current signal level | Notes |
| --- | --- | --- | --- |
| Login | `LOGIN` | L2 live card | Also updates presence. |
| Logout | `LOGOUT` | L1/L2 | Useful for presence; live card may be less important than state. |
| Page or section navigation | `PAGE_SYNC` | L2 live card | Confirms current location. High volume, should be throttled/deduped by page/section. |
| Dictionary lookup | `DEFINITION_LOOKUP` | L2 live card | Useful as vocabulary friction, especially repeated lookups. |
| AI help completed | `AI_HELP` | L2 live card | Useful, but should distinguish scope/mode and failure. |
| Human consultant request | `HELP_REQUEST` | L3 alert/action | Must remain visible and actionable. |

## Dashboard visibility rule

Consultant dashboard reader cards must keep recent activity visible even when the reader is not currently online. Presence is a separate annotation (`Active now` / `Recent activity`), not a filter.

This prevents a real reader history from disappearing when the session is older than the active-session window.

Dashboard labels must distinguish live sessions from recent feedback:

- `Logged In Now`: readers with an active session.
- `Recent Readers`: reader cards with recent activity loaded on the dashboard.
- `Recent Reader Activity`: the card section containing live and historical feedback.

## Gaps found in the first audit

Source: `docs/AGENT_LOOP_INVENTORY.md`, `docs/AI_ASSISTANT_FLOW_INVENTORY.md`, `internal/templates/reader/interaction.html`, `internal/handlers/api.go`.

| Reader input | Proposed event | Proposed level | Why it matters |
| --- | --- | --- | --- |
| Consultant prompt shown to reader | `CONSULTANT_PROMPT_SHOWN` | L1 quiet timeline | Confirms delivery of Consultant intervention. |
| Consultant prompt dismissed | `CONSULTANT_PROMPT_DISMISSED` | L2 live card | Consultant should know the prompt did not help or was not wanted. |
| Consultant prompt accepted / Open AI Help | `CONSULTANT_PROMPT_ACCEPTED` | L2 live card | Shows Consultant guidance produced a reader action. |
| AI help opened but no question sent | `AI_HELP_OPENED` then timeout/close | L1, promote to L2 if repeated | Indicates hesitation or UI confusion. |
| Text selection mode started | `READING_CONTEXT_SELECTION_STARTED` | L0/L1 | Useful for UX diagnosis, not live by default. |
| Text selection abandoned | `READING_CONTEXT_SELECTION_ABANDONED` | L2 if repeated | Strong signal that selecting context is hard. |
| AI request failed | `AI_HELP_FAILED` | L2, L3 if repeated | Consultant may need to step in when Tier 2 fails. |
| Quiz generated | `QUIZ_STARTED` | L1 quiet timeline | Useful learning context. |
| Quiz answer checked | `QUIZ_ANSWERED` | L0 aggregate | Too noisy live; aggregate score is more useful. |
| Quiz finished | `QUIZ_COMPLETED` | L1/L2 | Shows comprehension/checkpoint; include score and scope. |
| Quiz abandoned | `QUIZ_ABANDONED` | L2 if repeated | Indicates friction or uncertainty. |
| Scan to locate started | `SCAN_STARTED` | L1 quiet timeline | Reader is trying to find location in physical book. |
| Scan matched page | `SCAN_SUCCEEDED` | L2 live card | Current reading location improved. |
| Scan failed | `SCAN_FAILED` | L2, L3 if repeated | App may be blocking the reader; Consultant can help locate page. |
| Ah Ah moment created | `AHA_MOMENT_CREATED` | L1 quiet timeline | Positive learning signal; not urgent. |
| Reader preference changed | `READER_PREFERENCE_CHANGED` | L0 audit-only | Useful support context only when debugging. |
| Verification completed | `BOOK_VERIFIED` | L1 quiet timeline | Reader is ready to use the book companion. |

## Autoresearch tasks

## #1: Reader Input Inventory

Blocked by: none
Type: Research
Status: active

### Question

Which Reader inputs exist today, and which already call `trackActivity()`?

### Current answer

Tracked today: `LOGIN`, `LOGOUT`, `PAGE_SYNC`, `DEFINITION_LOOKUP`, `AI_HELP`, `HELP_REQUEST`.

Known untracked inputs: consultant prompt shown/dismissed/accepted, AI help opened/abandoned/failed, context selection started/abandoned, quiz started/answered/completed/abandoned, scan started/succeeded/failed, Ah Ah moment created, reader preference changes, verification.

### Repeatable check

Run:

```bash
scripts/consultant_signal_autoresearch_check.sh
```

## #2: Signal Priority Rules

Blocked by: #1
Type: Discuss
Status: active

### Question

Which signals should be live-card events versus quiet history versus aggregate-only?

### Current answer

Use L2/L3 only for present-tense reading state, explicit help, repeated friction, failed Tier 2/scan flows, and Consultant-prompt feedback. Keep routine quiz answers, settings, and raw selection mechanics out of the live feed unless repeated failure promotes them.

## #3: Event Payload Contract

Blocked by: #2
Type: Research
Status: proposed

### Question

What fields must every event carry so Consultants get context without exposing excessive reader content?

### Initial shape

Recommended common fields: `event_type`, `book_id`, `page_number`, `section_id`, `scope`, `source`, `status`, `summary`, `severity`, and `correlation_id` for multi-step flows. Avoid raw long selected passages in live cards; store them only when needed for AI/help context.

## #4: Dashboard Presentation Model

Blocked by: #2, #3
Type: Prototype
Status: proposed

### Question

How should the Consultant dashboard group signals so it predicts reader need without overwhelming the Consultant?

### Initial shape

Group by reader card with three lanes: `Reading now`, `Friction`, `Needs human`. Keep `Learning/positive` signals in inspector history unless specifically requested.

## #5: Implementation Slices

Blocked by: #3, #4
Type: Research
Status: proposed

### Question

What is the safest order to add events?

### Initial order

1. Consultant prompt feedback: shown, accepted, dismissed.
2. AI failure/opened/abandoned signals.
3. Scan started/succeeded/failed.
4. Quiz started/completed/abandoned, with answer details aggregate-only.
5. Positive/low-urgency signals: Ah Ah moments and preference changes.

## Source files

- `docs/AGENT_LOOP_INVENTORY.md`
- `docs/AI_ASSISTANT_FLOW_INVENTORY.md`
- `docs/UPGRADE_ROADMAP.md`
- `internal/templates/reader/interaction.html`
- `internal/handlers/api.go`
- `internal/handlers/activity.go`
- `internal/templates/consultant/dashboard.html`
