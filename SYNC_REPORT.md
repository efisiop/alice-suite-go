======================================================================
  SYNC REPORT: Local vs GitHub Dump
======================================================================

GitHub dump:      259 files
Local (tracked):  297 files
Identical:        231 files
Different:        28 files

--- Only local (not in GitHub): 38 ---
  - .DS_Store
  - .env
  - Alice Suite/.obsidian/app.json
  - Alice Suite/.obsidian/appearance.json
  - Alice Suite/.obsidian/core-plugins.json
  - Alice Suite/.obsidian/graph.json
  - Alice Suite/.obsidian/workspace.json
  - Alice Suite/Untitled.md
  - Alice Suite/Welcome.md
  - Alice Suite/create a link.md
  - archive/.DS_Store
  - archive/deprecated/Untitled 2.csv
  - archive/deprecated/server_test.log
  - archive/old-code/alice-suite
  - archive/old-code/server
  - archive/old-code/static-files/static/viewer-glossary.html
  - archive/reference/.DS_Store
  - archive/reference/alice_wonderland.pdf
  - archive/reference/alice_wonderland_by_pages.txt
  - archive/reference/data/_REFERENCE_alice_glossary.json
  - archive/reference/data/_REFERENCE_alice_glossary.sql
  - data/alice-suite.db
  - data/alice-suite.db-shm
  - data/alice-suite.db-wal
  - docs/POSTGRES_MIGRATION.md
  - internal/database/pg_transform.go
  - internal/database/testsetup.go.bak
  - internal/templates/consultant/dashboard.html
  - internal/templates/consultant/reader-inspector.html
  - internal/templates/reader/interaction.html
  - migrations/008_seed_glossary_terms.sql
  - pkg/auth/auth.go
  - pkg/auth/auth_test.go
  - pkg/auth/jwt.go
  - pkg/auth/jwt_test.go
  - pkg/auth/session.go
  - pkg/auth/verification.go
  - scripts/sync-with-github-dump.py

--- Different content (both exist): 28 ---
  ~ archive/docs/TECHNICAL_SPECIFICATIONS.md
  ~ archive/reference/docs/old-docs/GETTING_STARTED.md
  ~ archive/scripts/test_consultant_login_flow.sh
  ~ cmd/compare-db-structure/main.go
  ~ cmd/diagnose-sections/main.go
  ~ cmd/fix-render/main.go
  ~ cmd/fix-sections/main.go
  ~ cmd/init-users/main.go
  ~ cmd/migrate/main.go
  ~ cmd/reader/main.go
  ~ cmd/server/main.go
  ~ cmd/verify-deployment/main.go
  ~ data/alice-suite-go
  ~ data/alice-suite.db.backup
  ~ go.mod
  ~ internal/config/config.go
  ~ internal/database/activity.go
  ~ internal/database/consultant.go
  ~ internal/database/consultant_prompts.go
  ~ internal/database/database.go
  ~ internal/database/queries.go
  ~ internal/database/sessions.go
  ~ internal/database/settings.go
  ~ internal/database/verification.go
  ~ internal/handlers/reader_activity.go
  ~ internal/middleware/heartbeat.go
  ~ render.yaml
  ~ start.sh

======================================================================