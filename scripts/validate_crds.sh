#!/usr/bin/env bash
# Optional CRD + example YAML validation (requires kubectl).
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if ! command -v kubectl >/dev/null 2>&1; then
  echo "[validate_crds] SKIP: kubectl not in PATH"
  exit 0
fi

echo "[validate_crds] Validating CRD manifests (client dry-run)..."
find "${ROOT_DIR}/api/crds" -type f -name '*.yaml' -print0 | sort -z | while IFS= read -r -d '' f; do
  echo "  - $f"
  kubectl apply --dry-run=client -f "$f" >/dev/null
done

if [[ -d "${ROOT_DIR}/examples/crds" ]]; then
  echo "[validate_crds] Parsing example YAML (syntax only; CRDs need not be installed)..."
  python3 - <<'PY'
import pathlib
try:
    import yaml
except ImportError:
    print("  SKIP: PyYAML not installed (pip install pyyaml)")
    raise SystemExit(0)
root = pathlib.Path("examples/crds")
for p in sorted(root.glob("*.yaml")):
    with p.open() as f:
        list(yaml.safe_load_all(f))
    print(f"  - parsed {p}")
PY
fi

echo "[validate_crds] OK"
