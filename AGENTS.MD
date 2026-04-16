READ ~/Projects/agent-scripts/AGENTS.MD BEFORE ANYTHING.

Scope: All Agentic Engineering tools (Poltergeist, Oracle, Summarize, MCPorter, Peekaboo, Trimmy, CodexBar, RepoBar, etc.) and this agent workflow apply only to this codebase — Alice Suite (alice-suite-go). Do not use them for other projects or repos unless explicitly asked.

Tool usage:
- Use tools the agent can run directly (Oracle, Summarize when in PATH and configured; MCP tools e.g. Render) as much as needed without asking.
- When a task requires the user to run a program (e.g. Poltergeist, Peekaboo, Trimmy, CodexBar, RepoBar), ask for approval by stating clearly: the program name in bold, and why it is needed.

---

**Project memory (use this for every task):**  
Alice Suite = Go app, physical book companion for classic literature (not a replacement). Three tiers: (1) instant dictionary, (2) AI assistance, (3) human consultant. Stack: Go backend, SQLite/Postgres, HTML templates in internal/templates. Layout: `cmd/` (migrate, init-users, verify-deployment), `internal/` (handlers, services, database, middleware, auth), `migrations/`, `docs/`. Reader app and Consultant dashboard; test creds in LOGIN_CREDENTIALS.md. For deeper lookup without reading many files: `qmd query "your question" -c alice` (qmd collection `alice` indexes this repo’s .md; see docs/QMD_PROJECT_MEMORY.md). Firecrawl (MCP/API) for benchmarking and UI/UX research: see docs/FIRECRAWL_ALICE_WORKFLOW.md.
