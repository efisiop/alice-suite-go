#!/bin/bash
set -e

echo "🧹 Starting Refresher Protocol..."

# Create archive directories
mkdir -p archive/reference/docs/completion-docs
mkdir -p archive/reference/docs/old-docs
mkdir -p archive/old-code/static-files
mkdir -p archive/old-code/cmd-stubs
mkdir -p archive/old-code/services
mkdir -p archive/old-code/models
mkdir -p archive/old-code/empty-pkgs
mkdir -p archive/reference/prompt-specs

# Archive completion documentation
echo "📦 Archiving completion documentation..."
[ -f "STEP_2_COMPLETE.md" ] && mv STEP_2_COMPLETE.md archive/reference/docs/completion-docs/
[ -f "STEP_3_COMPLETE.md" ] && mv STEP_3_COMPLETE.md archive/reference/docs/completion-docs/
[ -f "STEP_4_COMPLETE.md" ] && mv STEP_4_COMPLETE.md archive/reference/docs/completion-docs/
[ -f "STEP_5_COMPLETE.md" ] && mv STEP_5_COMPLETE.md archive/reference/docs/completion-docs/
[ -f "STEP_6_COMPLETE.md" ] && mv STEP_6_COMPLETE.md archive/reference/docs/completion-docs/
[ -f "STEP_7_COMPLETE.md" ] && mv STEP_7_COMPLETE.md archive/reference/docs/completion-docs/
[ -f "PHASE_0_COMPLETE.md" ] && mv PHASE_0_COMPLETE.md archive/reference/docs/completion-docs/

# Archive old documentation
echo "📦 Archiving old documentation..."
[ -f "ALICE_SUITE_RECOVERED_BRIEF.md" ] && mv ALICE_SUITE_RECOVERED_BRIEF.md archive/reference/docs/old-docs/
[ -f "GETTING_STARTED.md" ] && mv GETTING_STARTED.md archive/reference/docs/old-docs/
[ -f "IMPLEMENTATION_SUMMARY.md" ] && mv IMPLEMENTATION_SUMMARY.md archive/reference/docs/old-docs/
[ -f "GLOSSARY_LINKING_SUMMARY.md" ] && mv GLOSSARY_LINKING_SUMMARY.md archive/reference/docs/old-docs/
[ -f "OPEN_VIEWER.md" ] && mv OPEN_VIEWER.md archive/reference/docs/old-docs/
[ -f "README_SERVER.md" ] && mv README_SERVER.md archive/reference/docs/old-docs/
[ -f "FORWARD_PIPELINE_GUIDE.md" ] && mv FORWARD_PIPELINE_GUIDE.md archive/reference/docs/old-docs/

# Archive old static files
echo "📦 Archiving old static files..."
[ -d "static" ] && [ "$(ls -A static 2>/dev/null)" ] && mv static archive/old-code/static-files/ || true

# Archive empty cmd directories
echo "📦 Checking cmd directories..."
[ -d "cmd/consultant" ] && [ -z "$(ls -A cmd/consultant 2>/dev/null)" ] && mv cmd/consultant archive/old-code/cmd-stubs/ 2>/dev/null || true
[ -d "cmd/reader" ] && [ -z "$(ls -A cmd/reader 2>/dev/null)" ] && mv cmd/reader archive/old-code/cmd-stubs/ 2>/dev/null || true

# Archive unused services
if [ -d "internal/services" ] && [ "$(ls -A internal/services 2>/dev/null)" ]; then
    if ! grep -r "internal/services" internal/handlers/*.go 2>/dev/null | grep -q "import\|services\."; then
        echo "  → Archiving unused services/"
        mv internal/services archive/old-code/services/ 2>/dev/null || true
    fi
fi

# Archive unused models
if [ -d "internal/models" ] && [ "$(ls -A internal/models 2>/dev/null)" ]; then
    if ! grep -r "internal/models" internal/*/*.go 2>/dev/null | grep -q "import\|models\."; then
        echo "  → Archiving unused models/"
        mv internal/models archive/old-code/models/ 2>/dev/null || true
    fi
fi

# Archive placeholder handlers
if [ -f "internal/handlers/handlers.go" ]; then
    if grep -q "TODO\|stub\|placeholder" internal/handlers/handlers.go 2>/dev/null; then
        echo "  → Archiving placeholder handlers.go"
        mv internal/handlers/handlers.go archive/old-code/handlers-stub.go 2>/dev/null || true
    fi
fi

# Archive empty directories
[ -d "config" ] && [ -z "$(ls -A config 2>/dev/null)" ] && rmdir config 2>/dev/null || true
[ -d "docs" ] && [ -z "$(ls -A docs 2>/dev/null)" ] && rmdir docs 2>/dev/null || true
[ -d "tests" ] && [ -z "$(ls -A tests 2>/dev/null)" ] && rmdir tests 2>/dev/null || true

# Archive empty pkg directories
[ -d "pkg/ai" ] && [ -z "$(ls -A pkg/ai 2>/dev/null)" ] && mv pkg/ai archive/old-code/empty-pkgs/ 2>/dev/null || true
[ -d "pkg/dictionary" ] && [ -z "$(ls -A pkg/dictionary 2>/dev/null)" ] && mv pkg/dictionary archive/old-code/empty-pkgs/ 2>/dev/null || true

# Archive old files
echo "📦 Archiving old files..."
[ -f "alice-suite" ] && mv alice-suite archive/old-code/ 2>/dev/null || true
[ -f "server" ] && mv server archive/old-code/ 2>/dev/null || true
[ -f "START_SERVER.sh" ] && mv START_SERVER.sh archive/deprecated/ 2>/dev/null || true

# Archive prompt specifications
[ -d "prompt_specifications" ] && mv prompt_specifications archive/reference/prompt-specs/ 2>/dev/null || true

# Clean up logs
echo "🧹 Cleaning up logs..."
[ -d "logs" ] && find logs -name "*.log" -type f -delete 2>/dev/null || true

# Create archive index
DATE=$(date +%Y-%m-%d)
cat > archive/REFRESHER_${DATE//-/_}.md << EOF
# Refresher Protocol Execution - $DATE

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
EOF

echo ""
echo "✅ Refresher Protocol Complete!"
echo ""
echo "📊 Summary:"
echo "  - Archived completion docs"
echo "  - Archived old documentation"
echo "  - Archived unused code"
echo "  - Cleaned up logs"
echo ""
echo "📁 Active codebase is now clean and organized!"

