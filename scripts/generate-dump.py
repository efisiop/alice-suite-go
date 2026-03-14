#!/usr/bin/env python3
"""
Generate a GitHub-style dump file from the local alice-suite-go codebase.

Recreates the same format as the DeepWiki/GitHub dump:
  1. Directory tree at the top
  2. Each file's full content in FILE: blocks

Usage:
  python scripts/generate-dump.py [--output OUTPUT_PATH]
"""

import argparse
import os
import sys
from pathlib import Path

SKIP_DIRS = {
    ".git",
    ".DS_Store",
    "Alice Suite",
    "__pycache__",
    "node_modules",
    ".obsidian",
    ".claude",
    "venv",
    ".venv",
}

SKIP_FILES = {
    ".DS_Store",
    ".env",
}

BINARY_EXTENSIONS = {
    ".db", ".db-shm", ".db-wal",
    ".pdf",
    ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".svg",
    ".woff", ".woff2", ".ttf", ".eot",
    ".zip", ".tar", ".gz", ".bz2",
    ".exe", ".dll", ".so", ".dylib",
    ".pyc", ".pyo",
    ".class", ".jar",
}

REPO_NAME = "efisiop-alice-suite-go"


def should_skip_dir(name: str) -> bool:
    return name in SKIP_DIRS


def should_skip_file(path: Path, rel: str) -> bool:
    if path.name in SKIP_FILES:
        return True
    if path.suffix.lower() in BINARY_EXTENSIONS:
        return True
    return False


def is_binary(path: Path) -> bool:
    """Quick heuristic: read first 8KB and check for null bytes."""
    try:
        with open(path, "rb") as f:
            chunk = f.read(8192)
        return b"\x00" in chunk
    except Exception:
        return True


def collect_files(root: Path) -> list[str]:
    """Collect all includeable file paths relative to root, sorted."""
    results = []
    for dirpath, dirnames, filenames in os.walk(root):
        dirnames[:] = [
            d for d in sorted(dirnames) if not should_skip_dir(d)
        ]
        for fname in sorted(filenames):
            full = Path(dirpath) / fname
            rel = str(full.relative_to(root)).replace("\\", "/")
            if should_skip_file(full, rel):
                continue
            if is_binary(full):
                continue
            results.append(rel)
    return results


def build_tree(files: list[str]) -> str:
    """Build a directory tree string matching the dump format."""
    tree_dict: dict = {}

    for f in files:
        parts = f.split("/")
        node = tree_dict
        for part in parts:
            if part not in node:
                node[part] = {}
            node = node[part]

    lines = ["Directory structure:"]
    lines.append(f"└── {REPO_NAME}/")

    def render(node: dict, prefix: str, is_last_parent: bool = True):
        entries = sorted(node.keys(), key=lambda x: (not bool(node[x]), x))
        for i, name in enumerate(entries):
            is_last = i == len(entries) - 1
            connector = "└── " if is_last else "├── "
            children = node[name]
            if children:
                lines.append(f"{prefix}{connector}{name}/")
                ext = "    " if is_last else "│   "
                render(children, prefix + ext, is_last)
            else:
                lines.append(f"{prefix}{connector}{name}")

    render(tree_dict, "    ")
    return "\n".join(lines)


def generate_dump(root: Path, output: Path):
    """Generate the full dump file."""
    print(f"Scanning {root} ...")
    files = collect_files(root)
    print(f"  Found {len(files)} files to include")

    tree = build_tree(files)

    sep = "=" * 48

    with open(output, "w", encoding="utf-8") as out:
        out.write(tree)
        out.write("\n\n")

        for i, rel in enumerate(files):
            full = root / rel
            try:
                content = full.read_text(encoding="utf-8", errors="replace")
            except Exception as e:
                print(f"  Warning: could not read {rel}: {e}")
                continue

            out.write(f"{sep}\n")
            out.write(f"FILE: {rel}\n")
            out.write(f"{sep}\n")
            out.write(content)
            if content and not content.endswith("\n"):
                out.write("\n")

            if (i + 1) % 50 == 0:
                print(f"  Processed {i + 1}/{len(files)} files...")

    print(f"\nDone! {len(files)} files written to: {output}")
    size_mb = output.stat().st_size / (1024 * 1024)
    print(f"File size: {size_mb:.2f} MB")


def main():
    parser = argparse.ArgumentParser(
        description="Generate a dump file from local alice-suite-go codebase"
    )
    parser.add_argument(
        "--output", "-o",
        default=os.path.expanduser(
            "~/Downloads/efisiop-alice-suite-go-8a5edab282632443.txt"
        ),
        help="Output file path (default: overwrite existing dump in ~/Downloads)",
    )
    parser.add_argument(
        "--repo",
        default=os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
        help="Repo root directory",
    )
    args = parser.parse_args()

    root = Path(args.repo).resolve()
    if not (root / "go.mod").exists():
        print(f"Warning: {root} does not look like alice-suite-go (no go.mod)")

    output = Path(args.output).expanduser().resolve()
    generate_dump(root, output)


if __name__ == "__main__":
    main()
