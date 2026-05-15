# cli_summary

CLI implemented: `cmd/karl-fluid-governance`.

Modes:

- fixture mode: `-fixture`
- bundle mode: `-admission` + `-bundle`
- optional: `-policy`, `-requested-action`, `-evaluation-time`

Safety:

- local file replay only
- no cluster interaction
- no QMP interaction
- no runtime mutation
