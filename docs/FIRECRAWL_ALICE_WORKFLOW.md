# Firecrawl workflow for Alice Suite (development & UI/UX)

This document tells **humans and coding agents** how to use **Firecrawl** to improve Alice Suite: the **Reader** app (physical-book companion), **Consultant** tools, and shared **Go** backend. Use it whenever research, benchmarking, or optional “fetch URL → text” features are relevant.

---

## What Firecrawl is (in one line)

A hosted API that turns **web pages into clean markdown** (and related formats), so you can **read, compare, summarize, or feed text into AI** without copying messy HTML by hand.

---

## Prerequisites

| Item | Purpose |
|------|---------|
| **Firecrawl API key** | Required for HTTP API calls (`FIRECRAWL_API_KEY`). Store in env or local MCP config—**never commit keys**. |
| **Cursor MCP** (`firecrawl-mcp`) | Optional: lets the **IDE agent** scrape/search from chat when MCP is enabled and restarted. |
| **Network** | API calls need outbound HTTPS; some sites block bots (reCAPTCHA)—expect failures and retry with another URL. |

Project context: stack is **Go**, templates under `internal/templates/`, config via env (see `internal/config`). This doc does **not** require Firecrawl inside the live app unless you implement a feature that calls the API.

---

## Three ways to use Firecrawl

### 1. Cursor MCP (best for “while building”)

**When:** Researching competitors, award pages, UX articles, or summarizing a list of URLs **during design or implementation**.

**How:**

1. Ensure `firecrawl-mcp` is configured in Cursor MCP (user or project `.cursor/mcp.json`) with `FIRECRAWL_API_KEY` in `env`.
2. Restart MCP / Cursor after changes.
3. In chat, ask explicitly: e.g. “Scrape these URLs and summarize onboarding patterns” or “Extract the feature list from this page.”

**Agent note:** If MCP tools are not visible in a session, use **HTTP API** (section 2) from a terminal script or document the URLs for the user to open manually.

### 2. Firecrawl HTTP API (best for repeatable scripts)

**When:** You want **reproducible** research, CI-adjacent checks, or a one-off `curl`/`curl`+`jq` workflow.

**Endpoints (typical):**

- `POST https://api.firecrawl.dev/v1/scrape` — one URL → markdown (and metadata).
- `POST https://api.firecrawl.dev/v1/search` — web search that returns URLs/snippets (verify exact request/response in [Firecrawl docs](https://docs.firecrawl.dev)).

**Pattern:**

```bash
# Replace YOUR_KEY and URL. Do not commit keys.
curl -sS -X POST "https://api.firecrawl.dev/v1/scrape" \
  -H "Authorization: Bearer YOUR_KEY" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","formats":["markdown"],"onlyMainContent":true}'
```

Save outputs under `docs/research/` **only if** they contain no secrets and licensing is respected (often better to save **your own notes**, not full third-party text).

### 3. Future: Go service (only if product needs it)

**When:** A **reader-facing or consultant-facing** feature must load **trusted reference text** from a URL (e.g. optional “Reference link” next to AI chat).

**How (high level):**

1. Add `FIRECRAWL_API_KEY` to config (same pattern as other API keys).
2. Implement a small client (e.g. `internal/services/`) that calls `/v1/scrape`, caps length, and passes excerpt into existing AI `context` strings.
3. Enforce **allowlisted domains**, rate limits, and user-visible errors.

Do **not** add this until there is a clear product requirement and privacy review.

---

## Development phase: how to improve Alice with Firecrawl

### A. Competitive & pattern research

**Goal:** Find **successful web/apps** that solve adjacent problems (reading help, language, calm utility, education) and extract **patterns**, not pixels.

**Steps:**

1. **Define a question** (e.g. “How do apps combine reading + AI + glossary?”).
2. **Build a URL list** (5–15): award hubs, product homepages, “how it works” pages, help centers.
3. **Scrape** each URL (MCP or API); **trim** to main content when possible (`onlyMainContent: true`).
4. **Synthesize** in a short note: *Adopt / Avoid / Later* for Alice Reader vs Consultant.
5. **Turn into tasks** (e.g. “Clear primary action on reader screen”, “Stronger empty state for AI chat”).

**Reliable public hubs (starting points):**

- [Awwwards — Annual Awards](https://www.awwwards.com/annual-awards/) (yearly “best of web” craft).
- [Awwwards — Culture & Education](https://www.awwwards.com/websites/culture-education/) (visual inspiration; often institutional sites).
- [Webby Awards — Learning & Education apps](https://winners.webbyawards.com/winners/apps-software/software-services-platforms/learning-education) (named products, awards by year).
- [Nielsen Norman Group — Top UX articles](https://www.nngroup.com/articles/) (e.g. annual “top articles” lists for **AI, navigation, buttons, prompts**).

**Limitations:**

- Some listing pages use **reCAPTCHA**; scrapes may fail—open the site in a browser or pick a different URL.
- Firecrawl returns **text**, not subjective “quality”—**open the live app** for feel, motion, and accessibility.

### B. UX research backlog (tie to NN/g themes)

When improving UI/UX, scrape NN/g (or similar) **article pages** and map recommendations to Alice:

| Theme | Example Alice application |
|-------|---------------------------|
| AI + search behavior | Where glossary, AI chat, and help should live so users don’t get lost. |
| Button / control states | Reader controls: disabled, loading, focus for keyboard users. |
| Hidden navigation | If the hamburger or tabs hide too much, reading flow suffers. |
| Prompt structure (CARE-style) | Align with how `AskAI` builds prompts in `internal/services/ai_service.go`. |
| Service design + AI | Consultant workflows and handoff between AI and humans. |

Agents: **quote repo behavior** from code when suggesting prompt or UI changes—read the relevant handler/service first.

---

## UI/UX phase: what to do with Firecrawl vs the browser

| Use Firecrawl for | Use the browser for |
|-------------------|---------------------|
| Feature lists, pricing copy, help docs, long articles | Layout, spacing, animation, real performance |
| Comparing **wording** and **flows** described on marketing sites | Touch targets, color contrast, accessibility |
| Building a **benchmark matrix** (copy into a table) | Emotional “feel” and first-time onboarding |

**For Alice specifically:**

- **Reader:** prioritize **calm**, **book-first** context, **obvious** “where am I in the book” — benchmark **utility** and **learning** apps, not flashy campaign sites.
- **Consultant dashboard:** benchmark **admin/analytics** patterns (density, filters, alerts)—often different from consumer reading apps.

---

## Optional product feature: “reference URL” for AI (design rules)

If implemented later:

1. **Allowlist** domains (e.g. your own docs, Project Gutenberg, partner publishers)—no arbitrary open web by default.
2. **Cap** scraped characters sent to the model; **log** domain + length only, not full text in production logs if avoidable.
3. **User-visible** label: “Excerpt from [domain] for context.”
4. **Fallback** when scrape fails: show the error and do not silently hallucinate sources.

---

## Agent checklist (before changing UI or flows)

1. Read this doc and the relevant template under `internal/templates/reader/` or consultant UI.
2. Decide whether the task needs **Firecrawl** (research) or **only** code/templates.
3. If research: collect URLs → scrape → summarize → list **concrete** UI tasks (max 3–5 per iteration).
4. If implementing Go + Firecrawl: confirm product approval, allowlist, and env vars—**do not** ship keys in the repo.
5. Prefer **small, testable** UI changes aligned with existing patterns in the codebase.

---

## Security & ethics (non-negotiable)

- **Never** commit API keys or paste them into tracked markdown.
- **Respect** robots.txt, terms of service, and copyright—use Firecrawl for **research and short excerpts**, not for republishing full copyrighted pages inside the app without permission.
- **Rotate** keys if they appear in chat logs or issues.

---

## Related project files

- `AGENTS.md` — project memory and tool scope for agents.
- `docs/QMD_PROJECT_MEMORY.md` — optional `qmd` lookup for repo docs.
- `internal/services/ai_service.go` — AI prompt construction for Reader.
- `internal/config/config.go` — pattern for new env-backed API keys.

---

## Revision

Update this file when you add a **production** Firecrawl integration or change **MCP** layout for the team.
