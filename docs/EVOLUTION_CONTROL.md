# Alice Suite evolution control (simple guide)

This file defines how Alice Suite should control app evolution and documentation, based on the Karpathy `llm-wiki` pattern.

Use this as the source of truth for "what changed", "why", and "what is next".

---

## Goal

Keep all important product/technical decisions in repo files (not only in chat), so progress compounds over time.

---

## Default mode (always on)

By default, for any non-trivial task, Alice Suite should run in **knowledge-base-first mode**:

1. read `docs/wiki/index.md`
2. do the task (code/docs)
3. update `docs/wiki/index.md`
4. append `docs/wiki/log.md`

If the task is very small, keep documentation short (1 log entry is enough).

Use `docs/wiki/START_HERE.md` for the exact quick routine and copy-paste templates.

---

## Folder structure

- `docs/wiki/index.md` -> current map of active topics and where details live
- `docs/wiki/log.md` -> chronological timeline of changes
- existing docs in `docs/` -> source material (roadmap, upgrades, architecture, setup)

Raw source of truth remains code + migrations + existing docs.

---

## Required workflow for any meaningful change

1. **Plan briefly**
   - define problem and success criteria in 1-3 lines
2. **Implement**
   - make the code/doc change
3. **Document**
   - update one relevant topic in `docs/wiki/index.md`
   - append one new dated entry in `docs/wiki/log.md`
4. **Verify**
   - run the smallest useful checks (tests/lint/manual check)
5. **Close**
   - in final message, link the updated files

---

## Entry format rules

### `docs/wiki/index.md` (state-oriented)

For each active topic, keep:
- topic name
- current status (`planned`, `in-progress`, `done`, `blocked`)
- owner (`user`, `agent`, or both)
- links to source docs/files

### `docs/wiki/log.md` (time-oriented)

Append entries in this exact header format:

`## [YYYY-MM-DD] <type> | <short title>`

Allowed `<type>` values:
- `decision`
- `feature`
- `fix`
- `docs`
- `infra`
- `review`

Each entry should include:
- what changed
- why
- files touched
- verification done
- next step (optional)

---

## How to answer questions using this system

When user asks:
- **"What changed?"** -> read latest log entries
- **"What is the current state?"** -> read index first
- **"What should we do next?"** -> use roadmap docs + index statuses

Always cite concrete files, not chat memory.

---

## Keep it simple

- Do not create extra abstractions unless needed
- Prefer updating existing docs over creating many new docs
- If a change is tiny/trivial, one short log entry is enough

