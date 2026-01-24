# Refresher Protocol Execution - 2026-01-24

## Summary
Files moved from main codebase during refresher protocol execution.

## Archived Items

### Completion Documentation
- Step completion docs (STEP_*.md, PHASE_*.md)
- Old implementation summaries

### Old Code
- Empty cmd directories
- Unused services/models
- Old static files
- Placeholder handlers

### Reference Documentation
- Old documentation files
- Prompt specifications

## Active Codebase Structure

After cleanup:
- cmd/server/ - Main server entry point
- cmd/init-users/ - User initialization tool
- cmd/migrate/ - Database migration tool
- cmd/seed/ - Seed data tool
- internal/handlers/ - Active HTTP handlers
- internal/templates/ - Go HTML templates
- internal/static/ - Static assets
- internal/database/ - Database layer
- internal/realtime/ - Real-time features
- pkg/auth/ - Authentication package
- migrations/ - Database migrations
- data/ - Database files

## Essential Documentation

Active docs:
- README.md - Main readme
- DEPLOYMENT.md - Deployment guide
- LOGIN_CREDENTIALS.md - Login credentials
- TESTING_CHECKLIST.md - Testing checklist
- MIGRATION_TO_GO_COMPLETE.md - Migration guide
- FEATURE_INVENTORY.md - Feature inventory
