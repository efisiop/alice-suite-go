# Reader autoresearch workflow

This adapts the useful part of Karpathy's `autoresearch` idea to Alice Suite.
Do not import the ML training code. The value is the loop:

1. choose a narrow scope
2. run a fixed evaluator
3. try one change
4. rerun the evaluator
5. keep only improvements that pass the checks and reduce real friction

## Scope

Use this workflow for Reader app improvements only:

- `internal/templates/reader/interaction.html`
- `internal/templates/base.html` when needed for Reader layout shell behavior
- Reader-specific static assets or JavaScript only when directly required
- Reader handlers only when a UI change requires backend support

Avoid touching Consultant, Admin, migrations, deployment, or unrelated docs during a Reader autoresearch run.

## Fixed evaluator

Run the evaluator before and after each experiment:

```bash
scripts/reader_autoresearch_check.sh
```

For a live localhost or Render check:

```bash
BASE_URL=http://localhost:8080 scripts/reader_autoresearch_check.sh
BASE_URL=https://alice-suite-go.onrender.com scripts/reader_autoresearch_check.sh
```

The evaluator checks:

- Go server builds
- Reader interaction template contains the expected workspace markers
- Reader panel toggles and persistence code are present
- Central reader typography markers are present
- Optional live `/health` and `/reader/interaction` responses contain expected markers

## Results log

For longer runs, keep an untracked TSV:

```text
commit	status	score	description
```

Suggested statuses:

- `baseline`: first run before changes
- `keep`: evaluator passes and UX is clearly better
- `discard`: evaluator passes but UX is worse/no clearer
- `crash`: build/check/live page fails

Suggested score:

- `0`: fails required checks
- `1`: passes checks but weak UX value
- `2`: passes checks and has clear UX value
- `3`: passes checks, clear UX value, and simpler or more maintainable

## Keep/discard rule

Keep a change only when all are true:

- evaluator passes
- change is scoped to Reader
- UX improvement is visible or defensible
- complexity is proportional to the gain

Discard a change when it adds fragile layout code, global CSS side effects, hidden dependencies, or broad refactors for a narrow Reader issue.

## Autoresearch loop

1. Read `docs/wiki/index.md`, latest `docs/wiki/log.md`, and this file.
2. Run `scripts/reader_autoresearch_check.sh` for a baseline.
3. Pick one Reader friction point.
4. Make the smallest useful change.
5. Run the evaluator again.
6. Manually inspect localhost for UX quality when the change is visual.
7. Update `docs/wiki/index.md` and `docs/wiki/log.md` for kept non-trivial changes.
8. Commit only the relevant files.

## Good first experiments

- Reader layout clarity at mobile/tablet/desktop widths.
- Better wording and placement for Reader controls.
- Dictionary popup readability and hierarchy.
- AI Help entry points and selected-text workflow.
- Reducing console noise or brittle inline styling in the Reader template.

## What this is not

This is not a replacement for product judgment, accessibility testing, or user feedback. It is a disciplined loop for generating and filtering Reader improvements without drifting into random UI changes.
