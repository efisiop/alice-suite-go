# Sync Local Codebase with GitHub Dump

This tool compares your local `alice-suite-go` folder with a GitHub repository dump file and can update local files to match GitHub.

## Quick Start

### 1. View the differences (report only)

```bash
cd /Users/efisiopittau/Project_1/alice-suite-go
python3 scripts/sync-with-github-dump.py ~/Downloads/efisiop-alice-suite-go-8a5edab282632443.txt --report
```

### 2. Save the report to a file

```bash
python3 scripts/sync-with-github-dump.py ~/Downloads/efisiop-alice-suite-go-8a5edab282632443.txt --report --output SYNC_REPORT.md
```

### 3. See what would change (dry run, no files modified)

```bash
python3 scripts/sync-with-github-dump.py ~/Downloads/efisiop-alice-suite-go-8a5edab282632443.txt --sync --dry-run
```

### 4. Actually sync: make local match GitHub

**Warning:** This overwrites local files with the GitHub dump content. Make a backup first (e.g. `git stash` or commit your work).

```bash
python3 scripts/sync-with-github-dump.py ~/Downloads/efisiop-alice-suite-go-8a5edab282632443.txt --sync
```

## What the report means

| Symbol | Meaning |
|--------|---------|
| `+` | File exists in GitHub but not locally |
| `-` | File exists locally but not in GitHub |
| `~` | File exists in both but has different content |

## What is NOT synced (protected)

- `.env` – your local environment/config
- `AGENTS.md` – project rules
- `data/` – database files (kept as-is)

## Current summary (from last run)

- **259** files in the GitHub dump
- **231** files identical in both
- **28** files differ between local and GitHub
- **38** files only in local (e.g. `pkg/auth/`, `migrations/008_seed_glossary_terms.sql`, this sync script)
