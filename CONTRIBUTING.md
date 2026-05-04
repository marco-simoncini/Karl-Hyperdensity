# Contributing

Before opening a PR, run:

- `./scripts/validate.sh`
- `go test ./...`

Contribution guardrails:

- Do not move runtime implementation code from `Karl-Dashboard` into this repository before planned extraction phases.
- Keep schemas and examples aligned with documented contracts.
- Keep policy posture unchanged (`recommendation_only`, `operator_controlled`, autonomous mode off).
