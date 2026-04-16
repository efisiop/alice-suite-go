# AI agents & helpers — one-page checklist (Alice Suite)

**Scope:** Alice Suite (`alice-suite-go`) only. Use this to pick a tool without re-explaining the whole stack.

---

## Daily coding (start here)

| Tool | Use when |
|------|----------|
| **Cursor — Chat / Agent / Composer** | You want code edits, refactors, or commands run in the repo. Give a short goal; use `@file` or a folder to point at context. |

---

## Cursor MCP (optional)

| Area | Use when |
|------|----------|
| **Web research MCP** (e.g. Tavily) | You need live web answers without pasting long URLs by hand. |
| **Firecrawl MCP** | You need to scrape or benchmark sites; see `docs/FIRECRAWL_ALICE_WORKFLOW.md`. |

Restart Cursor after MCP config changes. Never commit API keys.

---

## Repo rules & CLI tools (`AGENTS.MD`)

Listed for **this codebase**; not every tool is needed every day.

| Name | Rough job |
|------|-----------|
| **Oracle** | Codebase / analysis CLI (when installed and on PATH). |
| **Summarize** | Summarization (when on PATH). |
| **MCPorter** | MCP bridge / `mcporter list` vs your MCP config. |
| **Poltergeist** | Automation — run only when you approve it. |
| **Peekaboo** | GUI automation — needs macOS permissions. |
| **Trimmy**, **CodexBar**, **RepoBar** | Menu bar helpers — approve before use. |

Full install and permissions: `AGENTIC_SETUP_NOTES.md` (repo root).

---

## Paperclip (orchestration)

| Tool | Use when |
|------|----------|
| **Paperclip** | You want **tasks, assignees, heartbeats, and a dashboard** for multiple agent runtimes (e.g. Cursor, Claude, HTTP adapters), not a single Cursor session. |

Repo clone / run: `~/Projects/oss/paperclip` (or upstream docs). Does not replace Alice’s Go server.

---

## Project memory (`qmd`)

| Command | Use when |
|---------|----------|
| `qmd query "…" -c alice` | You want answers from **this repo’s `.md` docs** without typing long context. |
| `qmd update` + `qmd embed` | After you change markdown docs. |

Details: `docs/QMD_PROJECT_MEMORY.md`.

---

## API keys (unlock models / APIs)

Set in your shell profile as needed (see `AGENTIC_SETUP_NOTES.md`): e.g. `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, `GEMINI_API_KEY`, `X_AI_API_KEY`. Some CLIs stay idle without keys.

---

## Pick-one guide

1. **Single feature or bug in Alice** → Cursor Agent + short instruction.  
2. **“What did we document?”** → `qmd query "…" -c alice`.  
3. **External URLs / benchmarks** → Firecrawl or web MCP from Cursor (if configured).  
4. **Many agents, tasks, tracking** → Paperclip + adapters.  
5. **Special steipete CLIs** → Only if the task fits; many need explicit approval to run.

---

*Last updated: 2026-04. Re-index docs with `qmd` after large changes to this file.*
