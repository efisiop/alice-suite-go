# AI Assistant flow inventory

This documents the current Reader AI Assistant interaction surface before simplifying it.

## Current entry points

1. **Services -> Help Center -> AI Help**
   - Calls `showAIHelp()`.
   - Opens the floating chat container without selected text.

2. **Text selection spark**
   - Reader selects text in the book.
   - Floating spark toolbar appears.
   - Clicking the spark calls `openAIHelpWithSelection()`.
   - The floating AI chat opens and the selected text is placed into the AI input as `Help me understand this: "..."`

3. **AI chat -> Select text from the book**
   - Calls `activateTextSelectionMode()`.
   - Enables page selection while chat remains open.
   - After text is selected, the chat shows an explicit `Add selected text` action.

4. **AI chat reading context bar**
   - Shows what the assistant is linked to before the reader asks.
   - Supports `This section`, `This page`, and `Selected text`.
   - `Ask about this` sends the chosen context directly.
   - When `Selected text` is chosen, the reading area enters selection mode and the chosen passage stays highlighted after selection.

5. **AI chat starter actions**
   - The empty chat state offers direct linked-text actions: explain what is happening, make the text easier to read, and what to notice.
   - Starter actions use the currently linked scope and send through the same `askAI()` path as typed questions.

6. **Hidden legacy AI helpers**
   - Older Explain/Simplify/Chat tabs and quick actions remain in code paths where needed, but are no longer visible as primary popup tools.

## Main conflicts observed

- There are two selection systems:
  - always-on page selection spark
  - chat-initiated selection mode
- The chat-initiated flow previously asked the reader to select text and press Enter, while Enter also sends chat messages. The visible flow now uses an explicit `Add selected text` action.
- The selected-text spark previously tried to open an older modal path; it now uses the same floating chat as the Services AI Help button.
- The assistant now shows its linked reading scope visibly, reducing the need for the reader to manually copy or remember context.
- The popup no longer presents tabs, visual generation, misunderstood-word finding, and separate text-selection controls as equal choices. This keeps the first AI interaction focused on reading context.
- The first-open empty state now presents direct reading actions instead of only instructions, reducing the gap between opening the popup and getting help.
- Custom selection now needs visible feedback: the page shows a selection-mode cue and persistent highlight so the reader can trust what was captured.
- Some older selection helpers still modify pointer events and global key handlers, which makes focus feel fragile and should be reviewed after the primary reader-facing flows are stable.

## Intended simplification direction

1. Keep one primary visible sequence:
   - select text in book
   - click AI spark
   - choose Selected text in the AI window
   - click `Ask about this`
2. When AI chat is already open, make selection explicit:
   - click `Select text from book`
   - select passage
   - click an `Add selected text` control
3. Avoid implicit Enter-based selection transfer.
4. Keep the chat floating window stable; do not auto-minimize or alter book pointer events during normal reading unless absolutely needed.
5. Prefer visible context selection over requiring the reader to select text for every question.

## Test checklist

- Open AI Help from Services with no selected text.
- Use `This section` and `Ask about this`.
- Use the empty-state starter actions and confirm they ask about the linked scope.
- Switch to `This page` and confirm the prompt changes.
- Confirm old AI tabs and quick action buttons are not visible in the popup.
- Select text in the book and open AI Help from the spark.
- Switch to `Selected text` and confirm the context bar follows the selected passage.
- Confirm the chosen selected-text passage remains highlighted in the book after mouseup/touchend.
- Use Explain mode with selected text.
- Use Simplify mode with selected text.
- Use Chat mode with typed question.
- Use Find misunderstood word.
- Use Visual.
- Minimize, restore, resize, and close chat.
- Return to selecting words for dictionary lookup after closing AI.
