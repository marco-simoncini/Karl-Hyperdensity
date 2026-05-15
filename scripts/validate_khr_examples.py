#!/usr/bin/env python3
"""Parse examples/khr/*.{json,yaml} for Sprint 5 contract fixtures (stdlib only)."""
from __future__ import annotations

import json
import sys
from pathlib import Path


def main() -> int:
    root = Path(__file__).resolve().parents[1]
    khr = root / "examples" / "khr"
    if not khr.is_dir():
        print("[validate_khr_examples] ERROR: examples/khr missing", file=sys.stderr)
        return 1
    for path in sorted(khr.iterdir()):
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
    print("[validate_khr_examples] SUCCESS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
