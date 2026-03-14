# QMD project memory (Alice Suite)

**qmd** ([@tobilu/qmd](https://github.com/tobi/qmd)) is installed globally and this project is indexed as a **collection** so you can search it like a memory.

## What’s set up

- **Collection name:** `alice`
- **Path:** this repo (`alice-suite-go`)
- **Indexed:** all `**/*.md` files (113 docs)
- **Context:** a short project summary is attached so search results are more useful for LLMs/agents

## Useful commands

```bash
# Hybrid search (best quality, uses keywords + vectors + reranking)
qmd query "how does reader login work" -c alice

# Keyword-only search
qmd search "consultant dashboard" -c alice

# Semantic search
qmd vsearch "database migrations" -c alice

# After adding or changing .md files: re-index and refresh embeddings
qmd update
qmd embed
```

## Index location

- Default index: `~/.cache/qmd/index.sqlite`
- If you see “directory does not exist”, run: `mkdir -p ~/.cache/qmd` then run `qmd` again.

## One-line reminder

**Alice Suite** = Go app, physical book companion (dictionary + AI + human consultant), SQLite/Postgres, `cmd/` + `internal/` + `migrations/` + `docs/`.
