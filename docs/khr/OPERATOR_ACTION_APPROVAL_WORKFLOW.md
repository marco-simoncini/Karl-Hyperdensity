# Operator action approval workflow (KHR-W)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-W |
| **Mode** | Local evidence only |
| **Apply** | **Never** — approval does not trigger apply |

---

## ActionApproval model

| Field | Description |
|-------|-------------|
| `actionId` | Stable id for the proposed action |
| `resourceFutureRef` | ResourceFuture simulation ref |
| `resourceLeaseRef` | Related ResourceLease ref (not applied) |
| `laneId` | Lane (e.g. `native-live`) |
| `certificationRef` | Certification evidence path |
| `policyGateResult` | KHR-V gate evaluation snapshot |
| `approvalState` | `pending` \| `approved` \| `rejected` \| `expired` |
| `approvedBy` | Operator id |
| `approvedAt` | RFC3339 timestamp |
| `expiresAt` | Pending TTL expiry |
| `reason` | Optional operator or system reason |

---

## CLI `khr-action-approval`

```bash
# Generate pending from gated ResourceFuture simulation
go run ./cmd/khr-action-approval -cmd=generate \
  -simulation=simulation-gated.json -registry=registry.json \
  -cert-ref=docs/evidence/... -out=pending-bundle.json

# Approve (local evidence only)
go run ./cmd/khr-action-approval -cmd=approve \
  -approval=pending.json -registry=registry.json -by=operator-a -out=approved.json

# Reject
go run ./cmd/khr-action-approval -cmd=reject \
  -approval=pending.json -by=operator-b -reason="not now" -out=rejected.json

# Expire (simulation)
go run ./cmd/khr-action-approval -cmd=expire -approval=pending.json -out=expired.json
```

---

## Approval gates

| Gate | Blocks when |
|------|-------------|
| Expired TTL | `now > expiresAt` on pending |
| Stale certification | Registry evidence outside `validForSeconds` |
| Failed policy gate | Attestation or certification state fails KHR-V gates |
| Non-pending state | Approve/reject only valid on `pending` |

---

## Evidence

```bash
./scripts/khr_action_approval_evidence.sh
```

Artifacts under `docs/evidence/khr-action-approval/`: pending, approved, rejected, expired, stale-blocked.

Prerequisites: native-live certification + `khr_cert_registry_policy_gates.sh`.
