#!/usr/bin/env python3
"""Validate Hyperdensity JSON schemas and examples with stdlib only."""

from __future__ import annotations

import json
import sys
from pathlib import Path


def fail(message: str) -> None:
    print(f"[validate_json] ERROR: {message}", file=sys.stderr)
    raise SystemExit(1)


def load_json(path: Path) -> dict:
    try:
        data = json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:  # pragma: no cover
        fail(f"{path} is not valid JSON: {exc}")
    if not isinstance(data, dict):
        fail(f"{path} must contain a top-level JSON object")
    return data


def main() -> int:
    repo_root = Path(__file__).resolve().parents[1]
    schema_paths = sorted((repo_root / "schemas").glob("*.json"))
    example_paths = sorted((repo_root / "examples").glob("*.json"))

    if not schema_paths:
        fail("no schema files found under schemas/*.json")
    if not example_paths:
        fail("no example files found under examples/*.json")

    schema_count = 0
    for path in schema_paths:
        data = load_json(path)
        schema_count += 1

        if "$schema" not in data:
            fail(f"{path} is missing required '$schema'")
        if "$id" not in data and "title" not in data:
            fail(f"{path} must contain either '$id' or 'title'")
        if "type" not in data:
            fail(f"{path} is missing required 'type'")

    example_count = 0
    for path in example_paths:
        _ = load_json(path)
        example_count += 1

    print(
        f"[validate_json] OK: parsed {schema_count} schema files and {example_count} example files"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
