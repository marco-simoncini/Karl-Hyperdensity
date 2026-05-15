# KHR cgroup path policy (Sprint 7)

This document defines how `khr-linux-agent` evaluates cgroup paths during **read-only discovery**.

## Scanned root

All candidates must resolve to a real directory whose **fully evaluated path** remains inside `scannedRoot` (after `filepath.EvalSymlinks`). If symlink resolution escapes the scan root, the candidate is **blocked** with a reason — the agent does not follow unchecked absolute links outside the jail.

## Optional allow prefix

When `-allow-path-prefix` is set:

- The prefix must parse to an absolute path.
- Any **selected** cgroup directory must also lie under that prefix (after symlink resolution).
- This is an operator safety belt for hosts with many cgroup subtrees unrelated to Karl workloads.

If the prefix is not under `scannedRoot`, discovery emits a **warning**; candidates may still be evaluated, but policy interactions can reject everything.

## Provider handle (`Cell.spec.providerHandle.cgroupPath`)

When present:

1. The path is parsed from JSON (invalid JSON produces a **warning**, not a crash).
2. The raw path and a **rebased** variant (host `/sys/fs/cgroup/...` → under `-cgroup-root`) are probed early in the candidate list.
3. No writes occur; missing directories are acceptable outcomes.

## Heuristic naming

Heuristic candidates intentionally prefer **specific** slice directories before the generic `karl.slice` root so that discovery selects a leaf scope when both exist.

## Relationship to apply

Path policy code (`pkg/khr/cgroup/path_policy.go`) is shared philosophy with future apply gates, but **Sprint 7 does not invoke apply**. `mutationsForbidden` remains `true` in JSON output.
