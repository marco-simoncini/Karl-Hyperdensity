#!/usr/bin/env python3
"""Parse examples/khr/**, examples/ingest/**, and examples/grandepadre/evidence-store/** for JSON/YAML fixtures (stdlib JSON; PyYAML for yaml)."""
from __future__ import annotations

import json
import sys
from pathlib import Path


def validate_tree(root: Path, rel: Path, required: bool) -> int:
    base = root / rel
    if not base.is_dir():
        if required:
            print(f"[validate_khr_examples] ERROR: {rel} missing", file=sys.stderr)
            return 1
        return 0
    for path in sorted(base.rglob("*")):
        if not path.is_file():
            continue
        if path.suffix.lower() not in (".json", ".yaml", ".yml"):
            continue
        if path.suffix.lower() == ".json":
            json.loads(path.read_text(encoding="utf-8"))
        else:
            try:
                import yaml  # type: ignore
            except ImportError:
                print("[validate_khr_examples] ERROR: PyYAML required for .yaml fixtures", file=sys.stderr)
                return 1
            list(yaml.safe_load_all(path.read_text(encoding="utf-8")))
        print(f"[validate_khr_examples] OK {path.relative_to(root)}")
    return 0


def main() -> int:
    root = Path(__file__).resolve().parents[1]
    if validate_tree(root, Path("examples") / "khr", required=True):
        return 1
    if validate_tree(root, Path("examples") / "ingest", required=True):
        return 1
    if validate_tree(root, Path("examples") / "grandepadre" / "evidence-store", required=True):
        return 1
    print("[validate_khr_examples] SUCCESS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
