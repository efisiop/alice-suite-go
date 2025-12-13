# Archive Directory

This directory contains old, deprecated, or reference code that is no longer part of the active codebase.

**Created:** 2025-01-23  
**Purpose:** Organize old/redundant files before migration to single Go application

---

## Structure

- `old-code/` - Old implementations and one-time tools kept for reference
  - `migration-tools/` - One-time migration/restructuring tools
  - `test-tools/` - Test utilities and scripts
  - `scripts/` - One-time shell scripts
- `reference/` - Reference documentation and schemas
  - `docs/` - Old documentation files
  - `data/` - Reference data files
- `deprecated/` - Deprecated scripts and utilities (already contains Python scripts)

---

## Archived Items

### Migration Tools (`old-code/migration-tools/`)

#### `_OLD_export/` - Data Export Tool
- **Archived:** 2025-01-23
- **Reason:** One-time tool used to export data during initial setup
- **Replacement:** N/A (one-time use completed)
- **Can Delete:** After 3 months if no longer needed

#### `_OLD_export_glossary/` - Glossary Export Tool
- **Archived:** 2025-01-23
- **Reason:** One-time tool used to export glossary data
- **Replacement:** N/A (one-time use completed)
- **Can Delete:** After 3 months if no longer needed

#### `_OLD_full_restructure/` - Full Database Restructuring Tool
- **Archived:** 2025-01-23
- **Reason:** One-time tool used to restructure database schema
- **Replacement:** N/A (one-time use completed)
- **Can Delete:** After 3 months if no longer needed

#### `_OLD_restructure/` - Database Restructuring Tool
- **Archived:** 2025-01-23
- **Reason:** One-time tool used to restructure database to page-based system
- **Replacement:** N/A (one-time use completed)
- **Can Delete:** After 3 months if no longer needed

#### `_OLD_link_glossary/` - Glossary Linking Tool
- **Archived:** 2025-01-23
- **Reason:** One-time tool used to link glossary terms to sections
- **Replacement:** N/A (one-time use completed)
- **Can Delete:** After 3 months if no longer needed

### Test Tools (`old-code/test-tools/`)

#### `_OLD_test_help_requests/` - Help Request Test Tool
- **Archived:** 2025-01-23
- **Reason:** Test utility, not part of production codebase
- **Replacement:** N/A (testing tool)
- **Can Delete:** After 1 month if no longer needed

#### `_OLD_test_page/` - Page Test Tool
- **Archived:** 2025-01-23
- **Reason:** Test utility, not part of production codebase
- **Replacement:** N/A (testing tool)
- **Can Delete:** After 1 month if no longer needed

#### `_OLD_viewer/` - Viewer Tool
- **Archived:** 2025-01-23
- **Reason:** One-time viewer utility, functionality integrated into main app
- **Replacement:** `cmd/reader/` or `cmd/consultant/`
- **Can Delete:** After 3 months if functionality confirmed in main app

### Scripts (`old-code/scripts/`)

#### `_OLD_RESTRUCTURE_DB.sh` - Database Restructuring Script
- **Archived:** 2025-01-23
- **Reason:** One-time script used during database migration
- **Replacement:** N/A (one-time use completed)
- **Can Delete:** After 3 months if no longer needed

### Reference Documentation (`reference/docs/`)

#### `_REFERENCE_MIGRATION_PLAN.md` - Original Migration Plan
- **Archived:** 2025-01-23
- **Reason:** Superseded by `MIGRATION_TO_GO_COMPLETE.md`
- **Replacement:** `MIGRATION_TO_GO_COMPLETE.md` (complete migration guide)
- **Can Delete:** After migration is complete and stable

### Reference Data (`reference/data/`)

#### `_REFERENCE_alice_glossary.json` - Glossary JSON Export
- **Archived:** 2025-01-23
- **Reason:** Reference data file, data already in database
- **Replacement:** Database table `alice_glossary`
- **Can Delete:** After confirming data is in database

#### `_REFERENCE_alice_glossary.sql` - Glossary SQL Export
- **Archived:** 2025-01-23
- **Reason:** Reference data file, data already in database
- **Replacement:** Database table `alice_glossary`
- **Can Delete:** After confirming data is in database

### Backup Files (`deprecated/`)

#### `_BACKUP_queries.go.backup` - Database Queries Backup
- **Archived:** 2025-01-23
- **Reason:** Backup file, current version is `queries.go`
- **Replacement:** `internal/database/queries.go`
- **Can Delete:** After confirming current version is stable

---

## Cleanup Policy

- **Review archived items quarterly**
- **Delete after confirmation that replacement is stable**
- **Keep reference documentation indefinitely** (unless explicitly replaced)
- **Test tools:** Can delete after 1 month if not needed
- **Migration tools:** Can delete after 3 months if migration is stable
- **Backup files:** Can delete after 1 month if current version is stable

---

## File Naming Convention

All archived files follow this naming convention:
- `_OLD_` prefix for old implementations
- `_REFERENCE_` prefix for reference documentation/data
- `_BACKUP_` prefix for backup files
- `_DEPRECATED_` prefix for deprecated utilities (already in use)

---

## Notes

- All archived code is kept for reference only
- Do not modify archived files
- If you need functionality from archived code, check if it exists in the active codebase first
- Archive directory is excluded from git tracking (see `.gitignore`)

---

**Last Updated:** 2025-01-23


