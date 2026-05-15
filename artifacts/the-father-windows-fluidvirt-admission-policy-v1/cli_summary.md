# cli_summary

CLI implemented: `cmd/karl-fluid-admission`.

Supported modes:

- `-fixture <path>` admission replay fixture input;
- `-bundle <path>` evidence bundle input;
- `-policy <path>` optional policy pack;
- `-requested-action <action>` optional action override;
- `-evaluation-time <RFC3339>` deterministic output.

Safety properties:

- local file input only;
- no cluster calls;
- no QMP calls;
- no mutations.
