# ResourcePort CR preview and sandbox apply (KHR-K)

Controlled path from **JSON-only observation** to **local CR preview** and optional **sandbox kubectl apply**. Production remains JSON-only by default.

## Defaults

| Flag / setting | Default |
|----------------|---------|
| `--emit-cr` | `false` |
| `--apply-cr` | `false` |
| `--i-understand-this-is-sandbox` | `false` |
| `resourcePortLoopEnabled` | `false` in shipped config |
| Emission mode | `observed-json` |

## Emission modes

| Mode | Meaning |
|------|---------|
| `observed-json` | Loop JSON only; no CR files, no cluster apply |
| `cr-preview` | Stable local `ResourcePort` CR JSON under `--loop-output-dir` |
| `cr-applied-sandbox` | Preview + `kubectl apply` in sandbox (all gates pass) |

## Apply gates (all required)

1. `resourcePortLoopEnabled: true`
2. `sandboxMode` + `linuxOnly`
3. Namespace in `allowedNamespaces` (e.g. `khr-runtime-sandbox`)
4. Label allowlist match (`khr.karl.io/sandbox=true`)
5. Cluster context `karl-metal-01@ovh` when discovery is used
6. `--emit-cr=true`
7. `--apply-cr=true`
8. `--i-understand-this-is-sandbox`

ResourcePort CRs are **cluster-scoped**. Names are prefixed `khr-` and labeled with `karl.io/sandbox-namespace`.

## Owner metadata (annotations)

| Key | Example |
|-----|---------|
| `karl.io/source` | `karl-host-runtime` |
| `karl.io/runtime-version` | from `host.RuntimeVersion` |
| `karl.io/observed-at` | RFC3339 |
| `karl.io/safety-mode` | `sandbox` |
| `karl.io/emission-mode` | `cr-preview` / `cr-applied-sandbox` |

## CLI examples

```bash
# JSON-only (default)
karl-host-runtime -mode=resourceport-loop -config=examples/khr/runtime-sandbox/karl-host-runtime-config-loop.yaml \
  -namespace=khr-runtime-sandbox -cluster-context=karl-metal-01@ovh

# Local CR preview
karl-host-runtime -mode=resourceport-loop ... -emit-cr=true -loop-output-dir=/tmp/khr-cr-preview

# Sandbox apply (opt-in)
karl-host-runtime -mode=resourceport-loop ... -emit-cr=true -apply-cr=true \
  -i-understand-this-is-sandbox -loop-output-dir=/tmp/khr-cr-preview

# Cleanup applied sandbox CRs
karl-host-runtime -mode=resourceport-cleanup -config=... -namespace=khr-runtime-sandbox \
  -cluster-context=karl-metal-01@ovh
```

## Evidence

`./scripts/khr_resourceport_cr_preview_evidence.sh` on `karl-metal-01@ovh`.

## Safety

- No ResourceLease apply from this path.
- No production namespace mutation.
- ISO ships preview disabled; see Karl-OS-ISO `KHR_HOST_RUNTIME_PREVIEW.md`.
