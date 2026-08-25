# Reader AI context module design

## Purpose

Make one reliable Reader AI flow without changing its user-facing capabilities:

1. choose a Reading Context — section, page, or Selected Passage;
2. see exactly what is linked;
3. choose an intent or write a question;
4. submit one AI Assistance Request.

This is a design record for the next implementation slice. It follows the current Reader AI inventory and fix list; it does not replace either document.

Sources:

- `docs/AI_ASSISTANT_FLOW_INVENTORY.md`
- `docs/READER_FIX_LIST.md` (R-002)
- `internal/templates/reader/interaction.html`
- `internal/handlers/api.go` (`HandleAskAI`)
- `internal/services/ai_service.go` (`AskAI`)

## Current design pressure

`internal/templates/reader/interaction.html` currently owns five separate concerns in the same script:

- deriving section, page, and selected-passage context;
- maintaining overlapping globals (`selectedText`, `currentAISelectedContextText`, `textSelectionModeActive`, and `aiChatSelectionModeActive`);
- rendering and clearing selection highlights;
- capturing global keyboard events and forcing focus into the AI input;
- constructing an AI question and the HTTP request.

The overlap produces two selection paths: the always-available selection spark and chat-initiated selection. Either can update the AI context, but neither has a single owner for the resulting state. This explains why selection, focus, and book navigation remain coupled in R-002.

## Proposed modules and seams

### `ReadingContextController`

Owns the active Reading Context. Its interface is intentionally small:

```text
setScope(section | page | selection)
beginSelectedPassage()
captureSelectedPassage(text, location)
current()
clearSelectedPassage()
```

`current()` returns a normalized value:

```text
{
  scope: section | page | selection,
  text: string,
  bookID: string,
  pageNumber: number,
  sectionID: string | null,
  sectionNumber: number | null
}
```

It is the only module allowed to decide what context is active. It must not manipulate the DOM, submit network requests, or decide AI intent.

### `ReaderSelectionAdapter`

Owns browser selection and persistent visual highlighting. Its interface is:

```text
startCapture(onPassageCaptured)
stopCapture()
clearHighlight()
```

The callback passes only a Selected Passage and location. The adapter does not know about chat, prompts, or AI request types. There is one adapter now; no extra abstraction should be introduced until a second selection implementation exists.

### `AIChatController`

Owns the chat panel and submits AI Assistance Requests. Its interface is:

```text
open()
renderContext(readingContext)
submit(intent, question, readingContext)
```

It may create reader-facing wording such as “Explain this page,” but does not read browser selection or page DOM directly. The controller uses the normalized Reading Context from `ReadingContextController`.

### `AIRequestAdapter`

Encapsulates the request to `/api/ai/ask`. Initially it preserves the existing wire contract (`book_id`, `interaction_type`, `question`, `section_id`, `context`) so this refactor can be behavior-preserving. A later API change may add structured location metadata only when a concrete consumer needs it.

## Ownership rules

| Concern | Owner | Must not own |
| --- | --- | --- |
| Active scope and selected passage | `ReadingContextController` | DOM mutation or network calls |
| Browser-selection listeners and highlights | `ReaderSelectionAdapter` | Chat focus or AI request construction |
| Panel visibility, focus, intent wording | `AIChatController` | Raw page selection |
| HTTP payload and response handling | `AIRequestAdapter` | UI state |
| Prompt construction, provider fallback, persistence | Go `AIService` | Browser/UI state |

## Migration plan

1. Introduce `ReadingContextController` inside `interaction.html`; route the context bar through `current()` while retaining existing behavior.
2. Route both selection paths through `ReaderSelectionAdapter` and remove duplicate state writes. Keep the visible spark and explicit `Add selected text` behavior unchanged.
3. Move chat input, focus, and request construction behind `AIChatController` and `AIRequestAdapter`.
4. Remove obsolete selection globals, legacy Enter-transfer handling, and modal-only helper paths only after behavior tests pass.
5. Extract the modules to static JavaScript files only when their interfaces are stable; initial co-location minimizes deployment risk.

## Behavior checks before and after each slice

- Section, page, and selected-passage context show the correct preview and location.
- The selection spark opens the same chat with the captured passage linked.
- Chat-initiated selection exposes `Add selected text`; Enter sends chat and never transfers selection.
- Closing or minimizing chat does not disable normal dictionary word selection.
- The submitted request contains the intended question, section ID, and rendered context.

## Non-goals

- No redesign of the AI response UI.
- No provider or prompt-policy change.
- No change to stored historical `AIInteraction` records.
- No broad extraction of the entire Reader template.
