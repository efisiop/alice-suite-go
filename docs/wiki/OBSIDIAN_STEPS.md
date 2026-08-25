# Obsidian + Alice Suite (step by step)

Use [Obsidian](https://obsidian.md/) to read and edit the same markdown files that live in this repo. Nothing special is required: Obsidian is a **viewer and editor** for the files on disk.

---

## 1) Install Obsidian

- Download Obsidian for macOS from the [official site](https://obsidian.md/).
- Open it and skip paid sync unless you want it. Local files are free.

## 2) Open this project as a vault

1. In Obsidian, choose **Open folder as vault** (or **Create new vault** → then pick *existing* folder — exact wording can vary by version).
2. Select the **root** of the Alice repo, for example:
   - `.../Project_1/alice-suite-go`
3. You should see folders like `docs/`, `cmd/`, `internal/` in the file list.

**Optional:** If you only want documentation in the sidebar, open **`docs`** as the vault instead of the whole repo. Then paths like `wiki/DIRECTION.md` are at the top of the list.

## 3) What to open first

- `docs/wiki/index.md` — what is *in progress* vs *done*
- `docs/wiki/DIRECTION.md` — *your* north star and sharp priorities
- `docs/wiki/START_HERE.md` — the 5-step routine
- `docs/EVOLUTION_CONTROL.md` — rules for agents and for you

## 4) When you change something important

1. Edit the file in Obsidian (or in Cursor — same files).
2. If it is a **non-trivial** product or process change, also:
   - update a row in `docs/wiki/index.md` if a topic’s status changes, and
   - add a dated block at the **top** of the entries in `docs/wiki/log.md` (newest first), using the format in `docs/EVOLUTION_CONTROL.md`.
3. Commit in Git when you are ready.

## 5) Optional Obsidian features

- **Wikilinks** `[[DIRECTION]]` work between notes. They are fine in Git; plain Markdown also works.
- **Graph** shows how notes connect — useful once you add links.
- Avoid relying on **plugins** as the *only* place for decisions: if it matters, put a sentence in `DIRECTION.md` or `log.md` so the repo stays the source of truth.

## 6) Old codebases and extra history

- Put long-term notes in this repo, e.g. `docs/legacy/` or `archive/`, as normal `.md` files.
- Then they show up in the same Obsidian vault and in Git history.

That is the full setup. The evolution source remains **the repo**; Obsidian is your comfortable desk.
