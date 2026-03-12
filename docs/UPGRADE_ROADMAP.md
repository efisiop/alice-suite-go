# Alice Suite – Upgrade Roadmap (Sequence List)

Nice upgrades for the three tiers, in order. Each item has 4–5 steps; **approve a step (or item) and we proceed**.

---

## 1. READER (Tier 1 – Instant dictionary & reading)

### 1.1 Optimize the definition output
**Goal:** Make the dictionary popup clearer and more useful (source, part of speech, examples).

**Steps:**
1. **Backend:** In the RPC/API that returns definitions, include `source` (glossary / cache / external), `part_of_speech`, and `phonetic` when available (e.g. from `DictionaryCache` / external API).
2. **Template:** In `internal/templates/reader/interaction.html`, update the dictionary popup HTML to show a small “Source: glossary” or “Source: dictionary” label and optional part of speech (e.g. “noun”).
3. **Template:** Show “Example” in the popup in a consistent way (e.g. always in a dedicated line or expandable section), and ensure the “See example” / “More examples” buttons work with the new data.
4. **Styling:** Add light styling so definition text, examples, and source are visually distinct (e.g. definition bold or larger, examples italic, source muted).
5. **Optional:** If the external API returns multiple meanings, add a simple “Meaning 1”, “Meaning 2” (or “noun”, “verb”) toggle or list in the popup.

---

### 1.2 Improve reading layout (section list + content)
**Goal:** Make the reader section list and main content easier to scan and use on small screens.

**Steps:**
1. **Template:** In `interaction.html`, give the section-snippets list a clear heading (e.g. “Sections”) and ensure the collapse/expand control is obvious (icon + label).
2. **Template/CSS:** Make section snippets more scannable: slightly larger tap targets, clearer active state, and optional short labels (e.g. “Ch 1”, “Ch 2”) if the data supports it.
3. **Template:** Ensure the main reading content area has a comfortable max-width and line-height so long paragraphs don’t feel cramped.
4. **Optional:** Add a “Jump to section” or “Table of contents” link at the top that scrolls to the section list or opens it if collapsed.
5. **Test:** Check on a real phone or narrow viewport that the section list and content don’t overlap and that taps work reliably.

---

### 1.3 Reader dashboard: clearer entry points
**Goal:** Make the Reader dashboard (Start Reading, Stats, My Page) clearer and a bit more inviting.

**Steps:**
1. **Template:** In `internal/templates/reader/dashboard.html`, add one short sentence under each card (e.g. “Pick up where you left off” for Open Book, “See your progress” for Stats).
2. **Template:** Ensure the “Recent Activity” block has a clear empty state (e.g. “No recent activity yet – start reading to see it here”).
3. **Optional:** Add a small “Last read” or “Continue reading” hint on the Open Book card if the backend can expose last section/page.
4. **Styling:** Align card styles with the rest of the app (same borders/shadows as consultant or reader interaction if you have a design system).
5. **Links:** Double-check all dashboard links (Reading, Stats, My Page, Logout) and that they work for logged-in readers.

---

## 2. AI ASSISTANCE (Tier 2)

### 2.1 Make AI answers easier to read (explain / simplify / chat)
**Goal:** Improve readability of AI responses in the reading interface (explain, simplify, chat).

**Steps:**
1. **Template:** In `interaction.html`, find where AI responses are rendered (e.g. in a panel or modal) and wrap the response text in a container with a class like `ai-response-content`.
2. **CSS:** Add styles for that container: comfortable line-height (e.g. 1.6), max-width, and optional subtle background so the AI block is visually distinct from the book text.
3. **Optional:** If the backend can return markdown, add a small markdown-to-HTML step (client-side or server-side) so lists and bold/italic in AI answers render correctly.
4. **UX:** Ensure “Copy” or “Close” actions (if any) are clearly visible and that loading state (“Thinking…”) is obvious.
5. **Test:** Run explain, simplify, and one chat flow and confirm responses look good on desktop and mobile.

---

### 2.2 Quiz (Test your knowledge) – clearer questions and feedback
**Goal:** Make quiz questions and feedback easier to understand and use.

**Steps:**
1. **Backend/Template:** Ensure each quiz item shows the question text clearly and that the multiple-choice options are labeled (e.g. A, B, C or 1, 2, 3) and clickable.
2. **Template/JS:** After the user picks an answer, show immediate feedback: “Correct” or “Incorrect – the answer is …” with a short explanation if the backend provides it.
3. **Template:** Add a simple progress indicator (e.g. “Question 2 of 5”) and a “Next” or “Finish” button so the flow is clear.
4. **Optional:** At the end, show a short summary (e.g. “You got 4 out of 5 correct”) and optionally a “Try again” or “Back to reading” button.
5. **Test:** Generate a quiz from a section and complete it once; fix any layout or logic bugs (e.g. double-submit, wrong answer highlighted).

---

### 2.3 AI panel organization (explain / simplify / chat in one place)
**Goal:** Group explain, simplify, and chat in a single, organized “AI help” area so the reader isn’t confused.

**Steps:**
1. **Template:** In the reading interface, introduce one “AI help” or “Ask AI” entry point (button or tab) that opens a single panel or drawer.
2. **Template/JS:** Inside that panel, use tabs or buttons: “Explain”, “Simplify”, “Chat” (and optionally “Quiz” if you want it there). Selecting one shows the relevant input and history.
3. **State:** Keep the last-used mode (explain/simplify/chat) so when the user reopens the panel, the same mode is selected.
4. **Styling:** Use the same `ai-response-content` styling from 2.1 so all AI output looks consistent.
5. **Optional:** Add a short hint under each mode (e.g. “Explain: get a short explanation of the selected text”) to guide new users.

---

## 3. CONSULTANT (Tier 3)

### 3.1 Make the Reader section more organized
**Goal:** Consultant “Readers” view easier to scan and use (list, filters, reader detail).

**Steps:**
1. **Template:** In `internal/templates/consultant/readers.html`, add a clear page title and a short subtitle (e.g. “Readers and their activity”).
2. **Template:** Ensure the readers table has a visible header row (Name, Last name, Purchased, Activated, Calendar/Activity) and that columns align; consider making the “active now” indicator (e.g. green dot) and “Last seen” or “Last activity” prominent.
3. **Template/JS:** If the list can be long, add a simple filter or search (e.g. by name or email) so consultants can find a reader quickly; start with a single search box that filters the visible rows.
4. **Template:** When clicking a reader, ensure the reader detail/inspector view (e.g. `reader-inspector.html`) opens or expands with a clear “Back to list” or “Close” so navigation is obvious.
5. **Optional:** Add a “Last 24h” or “Last 7 days” toggle for activity so consultants can narrow the time range without changing code.

---

### 3.2 Help requests – clearer list and status
**Goal:** Help requests page easier to scan and act on (pending / assigned / resolved).

**Steps:**
1. **Template:** In `internal/templates/consultant/help-requests.html`, ensure the filter buttons (e.g. All, Pending, Assigned, Resolved) are visible and that the active filter is highlighted.
2. **Template:** For each request card, show clearly: reader name (or ID), book/section if available, short message preview, status badge, and date.
3. **Template:** Add an empty state when there are no requests (e.g. “No help requests” or “No pending requests” when filtered).
4. **Template/Backend:** Ensure “Assign to me” and “Mark resolved” (or equivalent) actions are visible and that the list refreshes or updates after an action so the consultant sees the new status.
5. **Optional:** Add a “Sort by date” (newest first) default so the most recent requests appear at the top.

---

### 3.3 Consultant dashboard – reader cards and activity
**Goal:** Dashboard “reader cards” (or activity feed) more informative and consistent.

**Steps:**
1. **Template:** In `internal/templates/consultant/dashboard.html`, ensure each reader card has a consistent header (e.g. reader name + “Active now” or “Last seen 5 min ago”).
2. **Template:** In the activity list inside each card, use the existing event types (e.g. LOGIN, PAGE_SYNC, DEFINITION_LOOKUP, AI_QUERY, HELP_REQUEST) and show a short, human-readable label and timestamp for each.
3. **Template:** If a card is collapsible, make the expand/collapse icon and click area obvious; consider “Expand all” / “Collapse all” if there are many cards.
4. **Optional:** Add a link from each card to the full reader inspector (readers list or reader-inspector) so the consultant can jump to full detail.
5. **Test:** Log in as consultant, open dashboard with at least one active or recent reader, and confirm activity and links work.

---

## How to use this roadmap

- **Sequence:** Work through the list in order (1.1 → 1.2 → … → 3.3), or skip to a tier you care about most.
- **Approval:** For each item (or each step), you say “approved” or “ok to proceed” and we implement only that. You can approve one step at a time or a full item.
- **Scope:** Steps are kept small so each approval is a clear, achievable chunk.

When you’re ready, reply with something like:  
“Approve 1.1” or “Approve step 1 and 2 of 1.1” or “Approve full 3.1”, and we’ll proceed accordingly.
