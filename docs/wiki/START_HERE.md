# Start here (default workflow)

If you are not sure what to do, follow this exact flow.

---

## 5-step default routine

1. **Read state**
  - open `docs/wiki/index.md`
2. **Pick one task**
  - keep scope small (one bug, one feature, or one doc update)
3. **Do the work**
  - code and/or docs
4. **Record state**
  - update `docs/wiki/index.md` (topic status)
5. **Record timeline**
  - append `docs/wiki/log.md` using the template below

---

## Template: index row

Add or edit one row in `docs/wiki/index.md`:

`| <Topic> | <planned|in-progress|done|blocked> | <user|agent|user + agent> | <main files> |`

Example:

`| Reader glossary improvements | in-progress | user + agent | docs/UPGRADE_ROADMAP.md, internal/templates/reader/interaction.html |`

---

## Template: log entry

Append to `docs/wiki/log.md`:

`## [YYYY-MM-DD] <decision|feature|fix|docs|infra|review> | <short title>`

Then:

- what changed
- why
- files touched
- verification
- next step (optional)

---

## Simple rule

If it is not written in index/log, it is not part of the project memory yet.