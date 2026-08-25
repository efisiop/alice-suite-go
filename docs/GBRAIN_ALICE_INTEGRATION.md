# GBrain integration for Alice Suite

Status: reference integration.

Alice Suite now includes `garrytan/gbrain` as a pinned external reference at
`archive/reference/gbrain`.

## Source

- Repository: https://github.com/garrytan/gbrain.git
- Local OSS clone for inspection: `~/Projects/oss/gbrain`
- Alice submodule path: `archive/reference/gbrain`
- Pinned commit: `814258dda67945ffec9457a1e73980e947b7e462`
- Upstream date: `2026-06-24`
- License: MIT

## Why this is a reference integration

Alice Suite is a Go application with SQLite/Postgres and server-rendered HTML
templates. GBrain is a Bun/TypeScript personal/company knowledge system with
CLI, MCP, schema packs, skills, PGLite/Postgres support, and retrieval/synthesis
workflows.

Do not import the GBrain runtime directly into Alice Suite server code without a
specific product slice. The current useful boundary is:

- use GBrain as a reference implementation for knowledge retrieval, citation,
  graph traversal, schema packs, and agent memory workflows
- keep Alice Suite runtime dependencies unchanged until a dedicated feature
  needs GBrain as a sidecar, CLI, or MCP service
- preserve Alice-specific memory in `docs/`, `docs/wiki/`, and the existing
  Alice vault/import workflow

## Alice-relevant areas to inspect

- `archive/reference/gbrain/README.md` for product shape and setup modes
- `archive/reference/gbrain/AGENTS.md` for agent operating protocol
- `archive/reference/gbrain/skills/query/SKILL.md` for cited retrieval behavior
- `archive/reference/gbrain/skills/strategic-reading/SKILL.md` for applied
  reading workflows relevant to classic literature support
- `archive/reference/gbrain/src/core/schema-pack/base/` for schema-pack design
- `archive/reference/gbrain/src/mcp/` for MCP server design

## Candidate future slices

1. Reader research memory sidecar: use GBrain concepts as inspiration for a
   cited "what do we know about this reader/book/context" workflow.
2. Consultant knowledge retrieval: adapt the citation and gap-flagging contract
   from `skills/query/SKILL.md` for consultant-facing evidence summaries.
3. Literature strategy workflow: adapt `strategic-reading` patterns for
   problem-focused book guidance without replacing the physical book.
4. MCP bridge: evaluate whether Alice should expose a small MCP surface before
   adding a GBrain dependency.

## Maintenance

Update the submodule deliberately:

```bash
git -C archive/reference/gbrain fetch origin
git -C archive/reference/gbrain checkout <commit-or-tag>
```

Then update this file and append `docs/wiki/log.md` with the new pinned commit.
