#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VAULT_DIR="${ROOT_DIR}/Alice Suite"
INBOX_DIR="${VAULT_DIR}/Inbox"
NOTES_DIR="${VAULT_DIR}/Knowledge/Imported"
ATTACHMENTS_DIR="${VAULT_DIR}/Knowledge/Attachments"

mkdir -p "${INBOX_DIR}" "${NOTES_DIR}" "${ATTACHMENTS_DIR}"

if ! command -v textutil >/dev/null 2>&1; then
  echo "Error: textutil is required (macOS built-in), but was not found."
  exit 1
fi

processed_any=false

while IFS= read -r path; do
  [[ -f "${path}" ]] || continue

  processed_any=true

  filename="$(basename "${path}")"
  name="${filename%.*}"
  ext="${filename##*.}"
  ext="$(printf '%s' "${ext}" | tr '[:upper:]' '[:lower:]')"

  note_path="${NOTES_DIR}/${name}.md"

  case "${ext}" in
    md)
      cp "${path}" "${note_path}"
      echo "Imported markdown: ${filename}"
      ;;

    txt)
      {
        echo "# ${name}"
        echo
        cat "${path}"
      } > "${note_path}"
      echo "Imported text: ${filename}"
      ;;

    docx|doc)
      {
        echo "# ${name}"
        echo
        textutil -convert txt -stdout "${path}"
      } > "${note_path}"
      echo "Imported document: ${filename}"
      ;;

    pdf|png|jpg|jpeg|webp|gif|svg)
      cp "${path}" "${ATTACHMENTS_DIR}/${filename}"
      {
        echo "# ${name}"
        echo
        echo "Imported attachment:"
        echo
        echo "![[${filename}]]"
      } > "${note_path}"
      echo "Imported attachment + note: ${filename}"
      ;;

    *)
      echo "Skipped unsupported file type: ${filename}"
      ;;
  esac
done < <(rg --files "${INBOX_DIR}")

if [[ "${processed_any}" == false ]]; then
  echo "No files found in ${INBOX_DIR}"
  echo "Drop files there, then run this script again."
  exit 0
fi

echo
echo "Done. Open Obsidian vault at: ${VAULT_DIR}"
echo "See imported notes in: Knowledge/Imported"
echo "See attachments in: Knowledge/Attachments"
