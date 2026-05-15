#!/usr/bin/env bash
# CRD + provider contract + example YAML validation.
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

REQUIRED_PROVIDER_DOCS=(
  "docs/providers/KUBEVIRT_LEGACY_PROVIDER_CONTRACT.md"
  "docs/providers/KUBEVIRT_SHELL_CELL_MAPPING.md"
  "docs/providers/KUBEVIRT_LABEL_HANDLE_CONTRACT.md"
  "docs/providers/KUBEVIRT_MIGRATION_SAFETY.md"
)

REQUIRED_PROVIDER_API=(
  "api/providers/kubevirt/kubevirt-legacy-provider.yaml"
  "api/providers/kubevirt/kubevirt-label-contract.yaml"
  "api/providers/kubevirt/kubevirt-handle-contract.yaml"
)

REQUIRED_PROVIDER_EXAMPLES=(
  "examples/providers/kubevirt/shell-windows-desktop-kubevirt-legacy.yaml"
  "examples/providers/kubevirt/shell-linux-vm-kubevirt-legacy.yaml"
  "examples/providers/kubevirt/cell-kubevirt-vm-handle.yaml"
  "examples/providers/kubevirt/runtimeprovider-kubevirt-legacy.yaml"
  "examples/providers/kubevirt/resourceport-kubevirt-legacy.yaml"
  "examples/providers/kubevirt/resourcelease-kubevirt-guarded-example.yaml"
)

echo "[validate_crds] Checking required provider contract files..."
for f in "${REQUIRED_PROVIDER_DOCS[@]}" "${REQUIRED_PROVIDER_API[@]}" "${REQUIRED_PROVIDER_EXAMPLES[@]}"; do
  if [[ ! -f "${ROOT_DIR}/${f}" ]]; then
    echo "[validate_crds] ERROR: missing required file: ${f}" >&2
    exit 1
  fi
  echo "  - ok ${f}"
done

echo "[validate_crds] Parsing YAML (syntax) for CRDs, provider contracts, and examples..."
python3 - <<'PY'
import pathlib
import sys

try:
    import yaml
except ImportError:
    print("[validate_crds] ERROR: PyYAML required (pip install pyyaml)", file=sys.stderr)
    sys.exit(1)

root = pathlib.Path(".")


def load_file(p: pathlib.Path) -> None:
    with p.open() as f:
        list(yaml.safe_load_all(f))
    print(f"  - parsed {p}")


for p in sorted((root / "api/crds").rglob("*.yaml")):
    load_file(p)
for p in sorted((root / "api/providers/kubevirt").glob("*.yaml")):
    load_file(p)
for p in sorted((root / "examples/crds").glob("*.yaml")):
    load_file(p)
for p in sorted((root / "examples/providers/kubevirt").glob("*.yaml")):
    load_file(p)
PY

if command -v kubectl >/dev/null 2>&1; then
  echo "[validate_crds] Validating Kubernetes CRD manifests (kubectl client dry-run)..."
  find "${ROOT_DIR}/api/crds" -type f -name '*.yaml' -print0 | sort -z | while IFS= read -r -d '' f; do
    echo "  - $f"
    kubectl apply --dry-run=client -f "$f" >/dev/null
  done
else
  echo "[validate_crds] SKIP: kubectl not in PATH (CRD kubectl validation only)"
fi

echo "[validate_crds] OK"
