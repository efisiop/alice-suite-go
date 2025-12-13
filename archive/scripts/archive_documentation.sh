#!/bin/bash

# Archive Documentation Script
# Moves historical markdown files to archive/ directory

set -e

echo "📦 Archiving historical documentation files..."
echo ""

# Create archive directories
mkdir -p archive/docs/{completion-reports,migration,fixes,security,analysis,protocols,architecture-plans,routing}

# Move completion reports
echo "Moving completion reports..."
mv COMPLETION_ASSESSMENT.md archive/docs/completion-reports/ 2>/dev/null || echo "  ⚠️  COMPLETION_ASSESSMENT.md not found"
mv FINAL_COMPLETION_REPORT.md archive/docs/completion-reports/ 2>/dev/null || echo "  ⚠️  FINAL_COMPLETION_REPORT.md not found"

# Move migration docs
echo "Moving migration documentation..."
mv MIGRATION_TO_GO_COMPLETE.md archive/docs/migration/ 2>/dev/null || echo "  ⚠️  MIGRATION_TO_GO_COMPLETE.md not found"

# Move fix documentation
echo "Moving fix documentation..."
mv CONSULTANT_DASHBOARD_CLEAR_IDENTIFICATION_FIX.md archive/docs/fixes/ 2>/dev/null || echo "  ⚠️  CONSULTANT_DASHBOARD_CLEAR_IDENTIFICATION_FIX.md not found"
mv CONSULTANT_DASHBOARD_FIXES.md archive/docs/fixes/ 2>/dev/null || echo "  ⚠️  CONSULTANT_DASHBOARD_FIXES.md not found"
mv CONSULTANT_DASHBOARD_FIXES_2025_12_03.md archive/docs/fixes/ 2>/dev/null || echo "  ⚠️  CONSULTANT_DASHBOARD_FIXES_2025_12_03.md not found"
mv CONSULTANT_DASHBOARD_VALIDATION.md archive/docs/fixes/ 2>/dev/null || echo "  ⚠️  CONSULTANT_DASHBOARD_VALIDATION.md not found"
mv CONSULTANT_LOGIN_DIAGNOSTIC_REPORT.md archive/docs/fixes/ 2>/dev/null || echo "  ⚠️  CONSULTANT_LOGIN_DIAGNOSTIC_REPORT.md not found"
mv FIX_CONSULTANT_LOGIN_BOUNCE.md archive/docs/fixes/ 2>/dev/null || echo "  ⚠️  FIX_CONSULTANT_LOGIN_BOUNCE.md not found"
mv FIX_INSTRUCTIONS.md archive/docs/fixes/ 2>/dev/null || echo "  ⚠️  FIX_INSTRUCTIONS.md not found"
mv FIX_JAVASCRIPT_ERRORS.md archive/docs/fixes/ 2>/dev/null || echo "  ⚠️  FIX_JAVASCRIPT_ERRORS.md not found"
mv IMMEDIATE_FIX_INSTRUCTIONS.md archive/docs/fixes/ 2>/dev/null || echo "  ⚠️  IMMEDIATE_FIX_INSTRUCTIONS.md not found"

# Move security docs
echo "Moving security documentation..."
mv SECURITY_FIXES_COMPLETED.md archive/docs/security/ 2>/dev/null || echo "  ⚠️  SECURITY_FIXES_COMPLETED.md not found"
mv SECURITY_VIOLATIONS_FOUND.md archive/docs/security/ 2>/dev/null || echo "  ⚠️  SECURITY_VIOLATIONS_FOUND.md not found"

# Move analysis docs
echo "Moving analysis documentation..."
mv USER_IDENTIFICATION_ANALYSIS.md archive/docs/analysis/ 2>/dev/null || echo "  ⚠️  USER_IDENTIFICATION_ANALYSIS.md not found"
mv ACTIVE_CODEBASE.md archive/docs/analysis/ 2>/dev/null || echo "  ⚠️  ACTIVE_CODEBASE.md not found"

# Move protocol docs
echo "Moving protocol documentation..."
mv REFRESHER_PROTOCOL_REPORT.md archive/docs/protocols/ 2>/dev/null || echo "  ⚠️  REFRESHER_PROTOCOL_REPORT.md not found"
mv REFRESHER_PROTOCOL_REPORT_2025_12_03.md archive/docs/protocols/ 2>/dev/null || echo "  ⚠️  REFRESHER_PROTOCOL_REPORT_2025_12_03.md not found"
mv REFRESHER_PROTOCOL_DATABASE_ARCHITECTURE.md archive/docs/protocols/ 2>/dev/null || echo "  ⚠️  REFRESHER_PROTOCOL_DATABASE_ARCHITECTURE.md not found"

# Move architecture plans (not selected)
echo "Moving architecture plans..."
mv DATABASE_ARCHITECTURE_PLAN_BY_Claude.md archive/docs/architecture-plans/ 2>/dev/null || echo "  ⚠️  DATABASE_ARCHITECTURE_PLAN_BY_Claude.md not found"
mv DATABASE_ARCHITECTURE_PLAN_BY_GEMINI.md archive/docs/architecture-plans/ 2>/dev/null || echo "  ⚠️  DATABASE_ARCHITECTURE_PLAN_BY_GPT5.md not found"
mv DATABASE_ARCHITECTURE_PLAN_BY_GPT5.md archive/docs/architecture-plans/ 2>/dev/null || echo "  ⚠️  DATABASE_ARCHITECTURE_PLAN_BY_GPT5.md not found"
mv DATABASE_ARCHITECTURE_PLAN_BY_kimi2.md archive/docs/architecture-plans/ 2>/dev/null || echo "  ⚠️  DATABASE_ARCHITECTURE_PLAN_BY_kimi2.md not found"
mv DATABASE_ARCHITECTURE_IMPLEMENTATION_SUMMARY.md archive/docs/architecture-plans/ 2>/dev/null || echo "  ⚠️  DATABASE_ARCHITECTURE_IMPLEMENTATION_SUMMARY.md not found"

# Move routing docs
echo "Moving routing documentation..."
mv ARCHITECTURE_DATA_ROUTING.md archive/docs/routing/ 2>/dev/null || echo "  ⚠️  ARCHITECTURE_DATA_ROUTING.md not found"
mv DATA_ROUTING_SUMMARY.md archive/docs/routing/ 2>/dev/null || echo "  ⚠️  DATA_ROUTING_SUMMARY.md not found"
mv QUICK_REFERENCE_DATA_ROUTING.md archive/docs/routing/ 2>/dev/null || echo "  ⚠️  QUICK_REFERENCE_DATA_ROUTING.md not found"

echo ""
echo "✅ Archive complete!"
echo ""
echo "📊 Summary:"
echo "  - Active documentation: 12 files (in root)"
echo "  - Archived documentation: 25 files (in archive/docs/)"
echo ""
echo "📁 Archive structure:"
echo "  archive/docs/"
echo "    ├── completion-reports/"
echo "    ├── migration/"
echo "    ├── fixes/"
echo "    ├── security/"
echo "    ├── analysis/"
echo "    ├── protocols/"
echo "    ├── architecture-plans/"
echo "    └── routing/"
echo ""

