# Reader fix list

Use this as the working list for Reader interaction improvements. Keep items short, testable, and mark them when done.

## Status legend

- `todo`: not started
- `doing`: in progress
- `done`: implemented and verified

## Items

| ID | Status | Item | Success check |
| --- | --- | --- | --- |
| R-001 | done | Dictionary examples are clumsy and should not show Alice book excerpts. The Example action should show short, simple life-usage examples for the word. | Clicking Example in the dictionary popup/modal shows concise generated everyday examples, not long book excerpts. |
| R-003 | done | Dictionary popup hierarchy needs to make the selected entry and applicable definition the dominant content. Buttons and dictionary heading should be smaller, source should move to the bottom, and popup should open prominently near the center. | Popup opens centered/prominent, word and definition are visually dominant, Derivation/Example buttons are secondary, and source appears at the bottom. |
| R-007 | done | Dictionary popup needs a light focus background behind the entry and main definition. | Entry word and main definition sit together in a subtle light background block. |
| R-004 | done | Expanded Derivation and Example dictionary panels need an obvious small close control. | Each expanded Derivation/Example panel includes a small `x` that closes the panel and resets the related button. |
| R-005 | done | Left and right Reader sidebars need the same direct close affordance. | Sections and Services sidebars each include a small `x` in the header that hides the panel while preserving the bottom show controls. |
| R-006 | done | Previous/Next page navigation should be visible under Sections, not buried in Services navigation tools. | Neat backward/forward arrow controls appear under the Sections list; Go to Page and Scan to Locate remain in Services. |
| R-008 | done | AI chat text selection should not rely on an Enter-key instruction that conflicts with sending messages. | The AI chat selection flow says `Select text from the book` and exposes an explicit `Add selected text` button after text is selected. |
| R-009 | done | Opening AI help from the text-selection spark should use the same floating chat as the Services AI Help button. | Selected text opens the floating AI chat and pre-fills `Help me understand this: "..."` without relying on the removed legacy Bootstrap modal. |
| R-010 | done | AI Assistant popup should visually feel close to the dictionary popup design type. | Floating AI assistant uses the same 8px card shape, pale header, restrained border/shadow, quieter controls, and light focus surfaces as the dictionary popup. |
| R-011 | done | Services sidebar should place Navigation tools below Info Center. | Services accordion order is Help Center, Info Center, then Navigation tools. |
| R-012 | done | Dictionary examples should use the actual meaning of the defined word, not generic word-awareness sentences. | Generated examples receive the definition context, avoid `I heard the word...`, and produce meaning-aware examples such as chair/sitting examples. |
| R-013 | done | Reader service tools should visually align with the dictionary popup design type. | Scan to Locate, Quiz, Ah Ah Moments, Human Consultant, and manual Dictionary modals share a dictionary-style popup shell, centered placement, pale header, restrained border/shadow, and quieter controls. |
| R-014 | done | Dictionary examples need more variation, and repeated Example clicks should show a different set. | Generated examples are grouped into varied pairs; each Example click advances to the next pair while the close `x` hides the panel. |
| R-015 | done | Sections sidebar should stay visually stable and show more section entries without collapsing when there are only a few. | On desktop, the Sections card has a stable minimum height and the snippet list has its own scrollable area sized for roughly 4-5 section entries. |
| R-020 | done | Dictionary derivation should be concise, and generated examples should not frame the term as a vocabulary word. | Derivation is summarized to the main origin line, and generated examples use the looked-up word in natural context instead of `I used...` / `The word...` phrasing. |
| R-021 | done | Dictionary popup should offer a picture for the looked-up word or definition. | The popup includes a `Picture` action that first tries a Wikipedia thumbnail, then falls back to Alice's configured image-generation endpoint. |
| R-016 | done | AI Assistant needs a visible, streamlined way to link help to what the reader is currently reading. | AI Assistant shows a `Linked to` context bar with This section, This page, and Selected text scopes, plus an `Ask about this` action that prefills the chat from the chosen context. |
| R-017 | done | AI Assistant popup should stop showing multiple older tools when the reader is trying to link AI to reading context. | Popup now visibly shows only the context-linking tools, chat history, and one input/send row; old tabs and quick-tool buttons are hidden from the reader. |
| R-018 | done | AI custom selection should visibly highlight what the reader is selecting in the book. | Choosing `Selected text` activates clear selection mode, outlines the reading area, and leaves the chosen passage highlighted while the context preview updates. |
| R-019 | done | AI Assistant popup should feel more ergonomic on first open. | Empty state offers direct linked-text starter actions, the context card is visually calmer, and the input/send row is easier to target. |
| R-002 | doing | AI Assistant interaction is convoluted. Selection, chat window focus, and book navigation conflict. Create a complete inventory of AI assistant entry points, order them sequentially, and test each flow. | A documented AI assistant flow map exists, each entry point has been tested, and the UI has a simpler sequence for selecting text and asking for help. |

## Notes

- Prioritize the reader's focus and confidence over feature density.
- Keep dictionary and AI assistance available but secondary to reading.
- Use `docs/READER_AUTORESEARCH.md` and `scripts/reader_autoresearch_check.sh` for kept Reader changes.
