#!/usr/bin/env python3
"""
Sync Local Codebase with GitHub Dump

Compares the local alice-suite-go codebase with a GitHub repository dump file
and reports differences. Can optionally apply the dump content to make local match GitHub.

Usage:
  python sync-with-github-dump.py <dump_file.txt> [--report] [--sync] [--dry-run]

Examples:
  python sync-with-github-dump.py ~/Downloads/efisiop-alice-suite-go-8a5edab282632443.txt --report
  python sync-with-github-dump.py ~/Downloads/efisiop-alice-suite-go-8a5edab282632443.txt --sync --dry-run
  python sync-with-github-dump.py ~/Downloads/efisiop-alice-suite-go-8a5edab282632443.txt --sync
"""

import argparse
import os
import sys
from pathlib import Path


# Paths to exclude from comparison (local-only or generated)
EXCLUDE_LOCAL = {
    ".git",
    ".DS_Store",
    "Alice Suite",
    "archive/deprecated/venv",
    ".claude",
    "__pycache__",
    ".pyc",
    "bin/",  # compiled binaries
}
# Paths/files we never overwrite from dump (sensitive or local config)
EXCLUDE_SYNC = {
    ".env",
    ".gitignore",
    "AGENTS.md",
    "data/",  # databases and binary data
}


def should_exclude_local(path: str) -> bool:
    """Check if path should be excluded from local file listing."""
    path_lower = path.replace("\\", "/").lower()
    for exc in EXCLUDE_LOCAL:
        if exc in path_lower or path_lower.startswith(exc + "/"):
            return True
    return False


def should_exclude_sync(path: str) -> bool:
    """Check if we should not overwrite this file during sync."""
    path_norm = path.replace("\\", "/")
    for exc in EXCLUDE_SYNC:
        if exc.endswith("/"):
            if path_norm == exc.rstrip("/") or path_norm.startswith(exc):
                return True
        elif path_norm == exc or os.path.basename(path) == exc:
            return True
    return False


def parse_dump(dump_path: str) -> dict[str, str]:
    """
    Parse the GitHub dump file. Format:
    ================================================
    FILE: <relative_path>
    ================================================
    <content>
    """
    result = {}
    content = Path(dump_path).read_text(encoding="utf-8", errors="replace")
    separator = "\n================================================\nFILE: "
    parts = content.split(separator)
    # First part is directory structure; rest are file entries
    for part in parts[1:]:
        first_newline = part.find("\n")
        if first_newline < 0:
            continue
        filepath = part[:first_newline].strip()
        rest = part[first_newline + 1 :]
        # Rest starts with ================ line; content is after that
        sep2 = "================================================\n"
        if rest.startswith(sep2):
            body = rest[len(sep2) :]
        else:
            body = rest
        filepath = filepath.replace("\\", "/").strip()
        if filepath.startswith("efisiop-alice-suite-go/"):
            filepath = filepath[len("efisiop-alice-suite-go/") :]
        if filepath:
            result[filepath] = body.rstrip("\n") if body else ""
    return result


def get_local_files(root: Path) -> dict[str, str]:
    """Get all relevant files from local repo with their contents."""
    files = {}
    for entry in root.rglob("*"):
        if not entry.is_file():
            continue
        rel = str(entry.relative_to(root)).replace("\\", "/")
        if should_exclude_local(rel):
            continue
        try:
            content = entry.read_text(encoding="utf-8", errors="replace")
            files[rel] = content
        except Exception:
            pass  # Skip binary or unreadable files
    return files


def compare(dump_files: dict[str, str], local_files: dict[str, str]) -> dict:
    """Compare and return structured diff report."""
    dump_set = set(dump_files)
    local_set = set(local_files)

    only_in_github = sorted(dump_set - local_set)
    only_local = sorted(local_set - dump_set)
    both = sorted(dump_set & local_set)

    different = []
    identical = []
    for path in both:
        d = (dump_files.get(path) or "").rstrip("\n")
        l = (local_files.get(path) or "").rstrip("\n")
        if d != l:
            different.append(path)
        else:
            identical.append(path)

    return {
        "only_in_github": only_in_github,
        "only_local": only_local,
        "different": different,
        "identical": identical,
        "total_github": len(dump_set),
        "total_local_tracked": len(local_set),
    }


def format_report(report: dict, root: Path) -> str:
    """Build full report as string."""
    lines = []
    g = report["only_in_github"]
    l = report["only_local"]
    d = report["different"]

    lines.append("=" * 70)
    lines.append("  SYNC REPORT: Local vs GitHub Dump")
    lines.append("=" * 70)
    lines.append("")
    lines.append(f"GitHub dump:      {report['total_github']} files")
    lines.append(f"Local (tracked):  {report['total_local_tracked']} files")
    lines.append(f"Identical:        {len(report['identical'])} files")
    lines.append(f"Different:        {len(d)} files")
    lines.append("")
    if g:
        lines.append(f"--- Only in GitHub (missing locally): {len(g)} ---")
        for p in g:
            lines.append(f"  + {p}")
        lines.append("")
    if l:
        lines.append(f"--- Only local (not in GitHub): {len(l)} ---")
        for p in l:
            lines.append(f"  - {p}")
        lines.append("")
    if d:
        lines.append(f"--- Different content (both exist): {len(d)} ---")
        for p in d:
            lines.append(f"  ~ {p}")
        lines.append("")
    lines.append("=" * 70)
    return "\n".join(lines)


def print_report(report: dict, root: Path) -> None:
    """Print a human-readable diff report."""
    print(format_report(report, root))


def apply_sync(
    dump_files: dict[str, str],
    root: Path,
    dry_run: bool = True,
) -> list[str]:
    """
    Apply dump content to local. Returns list of affected paths.
    - Adds missing files
    - Overwrites differing files (except excluded)
    """
    affected = []
    for path, content in dump_files.items():
        if should_exclude_sync(path):
            continue
        local_path = root / path
        existed = local_path.exists()
        existing = local_path.read_text(encoding="utf-8", errors="replace") if existed else None
        new_content = (content or "").rstrip("\n")
        if existed and existing is not None:
            existing_norm = (existing or "").rstrip("\n")
            if existing_norm == new_content:
                continue
        if dry_run:
            action = "would create" if not existed else "would update"
        else:
            local_path.parent.mkdir(parents=True, exist_ok=True)
            local_path.write_text(new_content + "\n" if new_content else "", encoding="utf-8")
            action = "created" if not existed else "updated"
        affected.append(f"  {action}: {path}")
    return affected


def main():
    parser = argparse.ArgumentParser(
        description="Compare and sync local codebase with GitHub dump"
    )
    parser.add_argument(
        "dump_file",
        help="Path to the GitHub dump .txt file",
    )
    parser.add_argument(
        "--report",
        action="store_true",
        default=True,
        help="Print diff report (default: True)",
    )
    parser.add_argument(
        "--sync",
        action="store_true",
        help="Apply dump content to local (add missing, overwrite different)",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="With --sync: show what would be done without writing",
    )
    parser.add_argument(
        "--repo",
        default=".",
        help="Local repo root (default: current directory)",
    )
    parser.add_argument(
        "--output",
        metavar="FILE",
        help="Write detailed report to file",
    )
    args = parser.parse_args()

    dump_path = Path(args.dump_file).expanduser().resolve()
    if not dump_path.exists():
        print(f"Error: Dump file not found: {dump_path}")
        sys.exit(1)

    root = Path(args.repo).resolve()
    if not (root / "go.mod").exists():
        print(f"Warning: {root} may not be alice-suite-go (no go.mod)")

    print("Parsing dump...")
    dump_files = parse_dump(str(dump_path))
    print(f"  Found {len(dump_files)} files in dump")

    print("Scanning local repo...")
    local_files = get_local_files(root)
    print(f"  Found {len(local_files)} tracked files locally")

    report = compare(dump_files, local_files)

    if args.report:
        report_text = format_report(report, root)
        print(report_text)
        if args.output:
            Path(args.output).write_text(report_text, encoding="utf-8")
            print(f"\nReport written to: {args.output}")

    if args.sync:
        affected = apply_sync(dump_files, root, dry_run=args.dry_run)
        if args.dry_run:
            print("\n--- DRY RUN: would apply these changes ---")
        else:
            print("\n--- Applied changes ---")
        for line in affected:
            print(line)
        print(f"\nTotal: {len(affected)} files {'would be' if args.dry_run else ''} affected")

    sys.exit(0)


if __name__ == "__main__":
    main()
