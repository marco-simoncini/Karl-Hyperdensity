# Identity / Access Lease contract (read-only — KHR-AQ)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AQ |
| **Contract set** | `khr-tp-contract-v1` |
| **Mode** | Read-only Technical Preview |
| **Production** | **NOT production ready** |

---

## Purpose

Introduce **IdentityBinding** and **AccessLease** semantics to connect directory identity, entitlements, and ShellLease in the Access Graph — without enforcement, revoke, disconnect, or auth mutation.

---

## Core semantics

| Concept | Description |
|---------|-------------|
| **IdentityBinding** | Binds a `UserIdentityRef` to tenant/subject context for access decisions (read-only projection) |
| **AccessLease** | Time-bounded access grant linking identity binding to gateway route/session path |
| **PrivilegeLease** | Observed admin/session capability lease (documentation only in KHR-AQ) |
| **ShellEntitlement** | Entitlement to access a Shell via ShellLease (compatibility with directory GPO/AD) |
| **PolicyBinding** | Policy bundle refs governing session/route (no enforcement in KHR-AQ) |

---

## Relationship chain

```
UserIdentityRef
    └── binds ──► IdentityBinding
                      └── entitles ──► ShellEntitlement
                                            └── grants ──► ShellLeaseRef
AccessLeaseRef
    └── routes ──► GatewayRouteRef
PrivilegeLeaseRef
    └── observes ──► session/admin capability (read-only stub)
PolicyBindingRef
    └── governs ──► AccessLease / GatewayRoute (no apply)
```

| From | Edge | To |
|------|------|-----|
| `UserIdentityRef` | binds | `IdentityBinding` |
| `IdentityBinding` | entitles | `ShellEntitlement` |
| `ShellEntitlement` | grants | `ShellLeaseRef` |
| `AccessLeaseRef` | routes | `GatewayRouteRef` |
| `PrivilegeLeaseRef` | observes | capability descriptor |
| `PolicyBindingRef` | governs | `AccessLease` / route |

---

## Ref shapes (compatibility)

| Ref | Example id |
|-----|------------|
| `IdentityBindingRef` | `khr-compat/identity-binding-{sessionId}` |
| `AccessLeaseRef` | `khr-compat/access-lease-{sessionId}` |
| `ShellEntitlementRef` | `khr-compat/shell-entitlement-{sessionId}` |
| `PrivilegeLeaseRef` | `khr-compat/privilege-lease-observe` (stub) |
| `PolicyBindingRef` | `khr-compat/policy-binding-session-{sessionId}` |

All refs carry `compatibilityOnly: true` in KHR-AQ stubs.

---

## TP invariants (KHR-AQ)

| ID | Invariant |
|----|-----------|
| IAL-01 | All lease/binding exports are **read-only** |
| IAL-02 | **No enforcement** — no policy apply, no auth middleware |
| IAL-03 | **No revoke** / **no disconnect** automation |
| IAL-04 | **No privilege escalation** — `PrivilegeLease` is observe-only stub |
| IAL-05 | **Compatibility mapping** allowed (`poolId`, legacy session ids) |
| IAL-06 | **No CRD/runtime** changes in KHR-AQ |
| IAL-07 | **No production enable** claims |

---

## Consumers

| Repo | Artifact |
|------|----------|
| **rdp-GW** | `accessLeaseRef`, `policyBindingRefs`, `privilegeLeaseRefs` on access graph export |
| **karl-directoryservice** | `build_identity_binding_stub`, `build_shell_entitlement_stub` |
| **KARL-APP** | `AccessLeaseDescriptor`, `ShellEntitlementDescriptor`, `PolicyBindingDescriptor` |
| **Karl-Dashboard** | `accessLeaseSummary`, `identityBindingSummary`, `policyBindingSummary` |
| **Karl-Hyperdensity** | This contract |

---

## Related

- `ACCESS_GRAPH_CONTRACT.md`
- `SHELLLEASE_GATEWAYROUTE_CONTRACT.md`
