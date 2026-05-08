# Executor Input Boundary

Input boundary is planning-only and forbids raw runtime materials.

Forbidden input examples:
- raw_cgroup_path
- raw_qmp_command
- raw_qga_command
- raw_shell_command
- raw_secret
- kubeconfig/token
- unapproved or unaudited candidates
