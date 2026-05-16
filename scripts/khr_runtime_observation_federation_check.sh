#!/usr/bin/env bash
# KHR-AU: aggregate read-only continuity/access observations across sibling repos.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RUN_ID="${KHR_FEDERATION_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
OUT_DIR="${ROOT}/docs/evidence/khr-runtime-observation-federation/${RUN_ID}"
CLUSTER_CONTEXT="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"

find_repo() {
  local var_name="$1"
  shift
  local candidate resolved=""
  if [[ -n "${!var_name:-}" && -d "${!var_name}" ]]; then
    echo "$(cd "${!var_name}" && pwd)"
    return 0
  fi
  for candidate in "$@"; do
    [[ -n "${candidate}" && -d "${candidate}" ]] || continue
    resolved="$(cd "${candidate}" && pwd)"
    echo "${resolved}"
    return 0
  done
  return 1
}

RDP_GW="$(find_repo KHR_RDP_GW_PATH \
  "${ROOT}/../rdp-GW" "${ROOT}/../../rdp-GW" "/home/m.simoncini/rdp-GW" 2>/dev/null || true)"
INVENTORY="$(find_repo KHR_INVENTORY_PATH \
  "${ROOT}/../Karl-Inventory" "${ROOT}/../../Karl-Inventory" "/home/m.simoncini/GitHub/Karl-Inventory" 2>/dev/null || true)"
DASHBOARD="$(find_repo KHR_DASHBOARD_PATH \
  "${ROOT}/../Karl-Dashboard" "${ROOT}/../../Karl-Dashboard" "/home/m.simoncini/GitHub/Karl-Dashboard" 2>/dev/null || true)"
KARL_APP="$(find_repo KHR_KARL_APP_PATH \
  "${ROOT}/../KARL-APP" "${ROOT}/../../KARL-APP" "/home/m.simoncini/GitHub/KARL-APP" 2>/dev/null || true)"
DIRECTORY="$(find_repo KHR_DIRECTORYSERVICE_PATH \
  "${ROOT}/../karl-directoryservice" "${ROOT}/../../karl-directoryservice" "/home/m.simoncini/GitHub/karl-directoryservice" 2>/dev/null || true)"

mkdir -p "${OUT_DIR}"

log() { echo "[khr_runtime_observation_federation_check] $*" | tee -a "${OUT_DIR}/run.log"; }

log "runId=${RUN_ID} cluster=${CLUSTER_CONTEXT}"
log "rdp-GW=${RDP_GW:-<missing>} inventory=${INVENTORY:-<missing>} dashboard=${DASHBOARD:-<missing>}"

export ROOT RUN_ID OUT_DIR CLUSTER_CONTEXT
export RDP_GW="${RDP_GW:-}" INVENTORY="${INVENTORY:-}" DASHBOARD="${DASHBOARD:-}"
export KARL_APP="${KARL_APP:-}" DIRECTORY="${DIRECTORY:-}"

python3 <<'PY'
from __future__ import annotations

import json
import os
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

ROOT = Path(os.environ["ROOT"])
OUT_DIR = Path(os.environ["OUT_DIR"])
RUN_ID = os.environ["RUN_ID"]
CLUSTER = os.environ["CLUSTER_CONTEXT"]

RDP_GW = os.environ.get("RDP_GW", "")
INVENTORY = os.environ.get("INVENTORY", "")
DASHBOARD = os.environ.get("DASHBOARD", "")
KARL_APP = os.environ.get("KARL_APP", "")
DIRECTORY = os.environ.get("DIRECTORY", "")

TRUST_RANK = {
    "live-readonly": 4,
    "inventory-observed": 3,
    "projected-readonly": 2,
    "fixture-readonly": 1,
    "fixture": 1,
}


def load_json(path: Path) -> dict[str, Any] | None:
    if not path.is_file():
        return None
    with path.open() as f:
        return json.load(f)


def pick_rdpgw_summary() -> tuple[dict[str, Any] | None, Path | None, dict[str, Any] | None]:
    if not RDP_GW:
        return None, None, None
    evidence = Path(RDP_GW) / "docs/evidence/khr-accessgraph-continuity"
    if not evidence.is_dir():
        return None, None, None
    best: tuple[int, str, Path, dict[str, Any]] | None = None
    for name in sorted(evidence.iterdir()):
        if not name.is_dir():
            continue
        summary_path = name / "summary.json"
        if not summary_path.is_file():
            continue
        try:
            s = json.loads(summary_path.read_text())
        except (OSError, json.JSONDecodeError):
            continue
        if s.get("status") != "PASS":
            continue
        src = s.get("source", "")
        trust = s.get("trustLevel", "")
        rank = TRUST_RANK.get(trust, 0)
        if src == "live-readonly":
            rank = max(rank, 4)
        elif src == "fixture-readonly":
            rank = max(rank, 1)
        if best is None or rank > best[0] or (rank == best[0] and name.name > best[1]):
            best = (rank, name.name, summary_path, s)
    if not best:
        return None, None, None
    _, run_name, summary_path, summary = best
    graph_path = summary_path.parent / "accessgraph-session.json"
    graph = load_json(graph_path)
    return summary, summary_path, graph


def extract_graph_ids(graph: dict[str, Any] | None) -> dict[str, str]:
    if not graph:
        return {}
    ids: dict[str, str] = {}
    if graph.get("continuityLineageId"):
        ids["continuityLineageId"] = graph["continuityLineageId"]
    if graph.get("sessionCorrelationId"):
        ids["sessionCorrelationId"] = graph["sessionCorrelationId"]
    for node in graph.get("nodes") or []:
        ref = node.get("ref") or {}
        if ref.get("continuityLineageId") and "continuityLineageId" not in ids:
            ids["continuityLineageId"] = ref["continuityLineageId"]
        if ref.get("sessionCorrelationId") and "sessionCorrelationId" not in ids:
            ids["sessionCorrelationId"] = ref["sessionCorrelationId"]
    return ids


def trust_from_rdpgw(summary: dict[str, Any]) -> str:
    src = summary.get("source", "")
    trust = summary.get("trustLevel", "")
    if src == "live-readonly" or trust == "live-readonly":
        return "live-readonly"
    return "fixture-readonly"


def observation_source_rdpgw(summary: dict[str, Any]) -> str:
    return summary.get("observationSource") or (
        "rdpgw-live" if trust_from_rdpgw(summary) == "live-readonly" else "rdpgw-fixture"
    )


def add_obs(
    observations: list[dict[str, Any]],
    *,
    repo: str,
    observation_source: str,
    trust_level: str,
    path: str,
    summary: dict[str, Any],
    extra_ids: dict[str, str] | None = None,
) -> None:
    ids = dict(extra_ids or {})
    for key in ("continuityLineageId", "sessionCorrelationId", "federationCorrelationId"):
        if summary.get(key) and key not in ids:
            ids[key] = summary[key]
    if "federationCorrelationId" not in ids and ids.get("sessionCorrelationId"):
        ids["federationCorrelationId"] = ids["sessionCorrelationId"]
    observations.append(
        {
            "repo": repo,
            "observationSource": observation_source,
            "trustLevel": trust_level,
            "artifactPath": path,
            "status": summary.get("status", "PASS"),
            "readOnly": summary.get("readOnly", True),
            "mutating": summary.get("mutating", False),
            "noDisconnect": summary.get("noDisconnect", True),
            "noRevoke": summary.get("noRevoke", True),
            "continuityObserved": summary.get("continuityObserved", True),
            "productionReady": summary.get("productionReady", False),
            **ids,
        }
    )


observations: list[dict[str, Any]] = []
errors: list[str] = []

# rdp-GW
rdp_summary, rdp_path, rdp_graph = pick_rdpgw_summary()
if rdp_summary and rdp_path:
    graph_ids = extract_graph_ids(rdp_graph)
    add_obs(
        observations,
        repo="rdp-GW",
        observation_source=observation_source_rdpgw(rdp_summary),
        trust_level=trust_from_rdpgw(rdp_summary),
        path=str(rdp_path.relative_to(RDP_GW) if RDP_GW else rdp_path),
        summary=rdp_summary,
        extra_ids=graph_ids,
    )
else:
    errors.append("missing rdp-GW PASS accessgraph continuity summary")

# Inventory stub
if INVENTORY:
    inv_stub = Path(INVENTORY) / "examples/khr/federation-observation-stub.json"
    inv = load_json(inv_stub)
    if inv:
        add_obs(
            observations,
            repo="Karl-Inventory",
            observation_source=inv.get("observationSource", "inventory-observed"),
            trust_level=inv.get("trustLevel", "inventory-observed"),
            path=str(inv_stub.relative_to(INVENTORY)),
            summary=inv,
        )
    else:
        errors.append(f"missing {inv_stub}")

# Dashboard projection fixture
if DASHBOARD:
    dash_path = Path(DASHBOARD) / "examples/khr-dashboard/access-graph-continuity-summary.json"
    dash = load_json(dash_path)
    if dash:
        dash_summary = {
            **dash,
            "trustLevel": dash.get("trustLevel", "projected-readonly"),
            "observationSource": dash.get("observationSource", "dashboard-projected"),
        }
        add_obs(
            observations,
            repo="Karl-Dashboard",
            observation_source=dash_summary["observationSource"],
            trust_level=dash_summary["trustLevel"],
            path=str(dash_path.relative_to(DASHBOARD)),
            summary=dash_summary,
        )
    else:
        errors.append(f"missing {dash_path}")

# KARL-APP projected stub
if KARL_APP:
    app_stub = Path(KARL_APP) / "examples/khr/app-federation-observation-stub.json"
    app = load_json(app_stub)
    if app:
        add_obs(
            observations,
            repo="KARL-APP",
            observation_source=app.get("observationSource", "karl-app-projected"),
            trust_level=app.get("trustLevel", "projected-readonly"),
            path=str(app_stub.relative_to(KARL_APP)),
            summary=app,
        )

# Directory service stub
if DIRECTORY:
    ds_stub = Path(DIRECTORY) / "examples/khr/identity-federation-observation-stub.json"
    ds = load_json(ds_stub)
    if ds:
        add_obs(
            observations,
            repo="karl-directoryservice",
            observation_source=ds.get("observationSource", "directoryservice-projected"),
            trust_level=ds.get("trustLevel", "projected-readonly"),
            path=str(ds_stub.relative_to(DIRECTORY)),
            summary=ds,
        )

# Correlation merge
lineage_ids = {o["continuityLineageId"] for o in observations if o.get("continuityLineageId")}
session_ids = {o.get("federationCorrelationId") or o.get("sessionCorrelationId") for o in observations}
session_ids = {x for x in session_ids if x}

lineage_match = len(lineage_ids) <= 1
session_match = len(session_ids) <= 1

primary_trust = "fixture-readonly"
for o in observations:
    t = o.get("trustLevel", "fixture-readonly")
    if TRUST_RANK.get(t, 0) >= TRUST_RANK.get(primary_trust, 0):
        primary_trust = t

federation_correlation_id = ""
continuity_lineage_id = ""
session_correlation_id = ""
if session_ids:
    federation_correlation_id = sorted(session_ids)[0]
    session_correlation_id = federation_correlation_id
if lineage_ids:
    continuity_lineage_id = sorted(lineage_ids)[0]

continuity_flags = all(
    o.get("continuityObserved") is True
    and o.get("mutating") is False
    and o.get("noRevoke") is True
    and o.get("noDisconnect") is True
    and o.get("productionReady") is False
    for o in observations
)

status = "PASS"
if errors or not lineage_match or not session_match or not continuity_flags or not observations:
    status = "FAIL"

trust_map = {o["repo"]: o["trustLevel"] for o in observations}

summary: dict[str, Any] = {
    "phase": "khr-runtime-observation-federation",
    "sprint": "KHR-AU",
    "runId": RUN_ID,
    "clusterContext": CLUSTER,
    "contractSetId": "khr-tp-contract-v1",
    "status": status,
    "readOnly": True,
    "mutating": False,
    "noDisconnect": True,
    "noRevoke": True,
    "noSessionMutation": True,
    "continuityObserved": all(o.get("continuityObserved") for o in observations),
    "productionReady": False,
    "noAutonomousOrchestration": True,
    "primaryTrustLevel": primary_trust,
    "trustLevels": trust_map,
    "federationCorrelationId": federation_correlation_id,
    "continuityLineageId": continuity_lineage_id,
    "sessionCorrelationId": session_correlation_id,
    "observationSources": observations,
    "federationBundle": {
        "kind": "FederationBundle",
        "federationCorrelation": {
            "kind": "FederationCorrelation",
            "federationCorrelationId": federation_correlation_id,
            "continuityLineageId": continuity_lineage_id,
            "sessionCorrelationId": session_correlation_id,
        },
        "mergeRulesApplied": [
            "lineage-correlation",
            "session-correlation",
            "continuity-precedence",
        ],
        "continuityPrecedence": primary_trust,
    },
    "consistency": {
        "lineageCorrelationMatch": lineage_match,
        "sessionCorrelationMatch": session_match,
        "continuityObservedConsistent": continuity_flags,
        "errors": errors,
    },
    "evidencePath": f"docs/evidence/khr-runtime-observation-federation/{RUN_ID}",
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}

out_path = OUT_DIR / "federation-summary.json"
out_path.write_text(json.dumps(summary, indent=2) + "\n")
print(f"[khr_runtime_observation_federation_check] federation-summary={out_path}")
print(f"[khr_runtime_observation_federation_check] status={status} primaryTrustLevel={primary_trust}")
print(f"[khr_runtime_observation_federation_check] sources={len(observations)} lineageMatch={lineage_match} sessionMatch={session_match}")
if status != "PASS":
    for e in errors:
        print(f"[khr_runtime_observation_federation_check] FAIL: {e}", file=sys.stderr)
    if not lineage_match:
        print(f"[khr_runtime_observation_federation_check] FAIL: lineage ids mismatch {lineage_ids}", file=sys.stderr)
    if not session_match:
        print(f"[khr_runtime_observation_federation_check] FAIL: session ids mismatch {session_ids}", file=sys.stderr)
    sys.exit(1)
PY

log "PASS federation-summary=${OUT_DIR}/federation-summary.json"
