# Actuator Dry-Run

Evidence:

- `raw_logs_sanitized/actuator_dry_run_output.json`

Dry-run result:

- mode: `dry-run`
- allowed: `true`
- policyDecision: `dry-run-accepted`
- target file validated: `.../cpu.max`
- before cpu.max: `600000 100000`
- requested cpu.max (request): `600000 100000`
- blockers: none

Dry-run accepted with identity/path/allowlist checks.
