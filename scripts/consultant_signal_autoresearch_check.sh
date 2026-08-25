#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

reader_template="internal/templates/reader/interaction.html"
consultant_dashboard="internal/templates/consultant/dashboard.html"
activity_handler="internal/handlers/activity.go"

echo "== Reader activity emitters =="
rg -n "trackActivity\\(" "$reader_template" "$activity_handler" || true

echo
echo "== Consultant live-card labels =="
rg -n "'(LOGIN|LOGOUT|PAGE_SYNC|SECTION_SYNC|DEFINITION_LOOKUP|AI_QUERY|AI_HELP|HELP_REQUEST|FEEDBACK_SUBMISSION)'|function formatEventType|function getEventIcon" "$consultant_dashboard" || true

echo
echo "== Dashboard keeps historical activity visible =="
if rg -n "Filtered activities \\(active readers only\\)|filterReaderCardsToActiveOnly|Removing inactive reader card|Reader windows will appear here as readers become active" "$consultant_dashboard"; then
  echo "Dashboard still hides reader activity behind active-session filtering." >&2
  exit 1
fi
if rg -n ">Active Readers<|card-title\">Active Readers|mb-0\">Active Readers" "$consultant_dashboard"; then
  echo "Dashboard still labels recent reader activity as active readers." >&2
  exit 1
fi
rg -n "Logged In Now|Recent Readers|Recent Reader Activity" "$consultant_dashboard"
echo "OK: recent reader activity is not filtered out by active-session presence."

echo
echo "== Untracked candidate Reader flows =="
rg -n "function (dismissConsultantPromptBubble|acceptConsultantPromptAndOpenAIHelp|showAIHelp|activateTextSelectionMode|generateQuiz|submitQuizCurrentAnswer|showQuizSummary|quizBackToScope|submitAhAhMoment|startScanToLocate|scanUploadedImage|captureAndProcessImage|showScanResult|showScanError)|/api/reader/prompt-(dismiss|accept)|/api/reader/quiz|/api/reader/ah-ah-moments|find_page_by_text" "$reader_template" internal/handlers/api.go || true

echo
echo "== Signal design doc =="
test -f docs/CONSULTANT_SIGNAL_AUTORESEARCH.md
rg -n "CONSULTANT_PROMPT|AI_HELP|QUIZ_|SCAN_|AHA_|READER_PREFERENCE|BOOK_VERIFIED|Level" docs/CONSULTANT_SIGNAL_AUTORESEARCH.md
