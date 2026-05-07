# Safety Scan

Safety scan executed on modified files with patterns:

- `Bearer `
- `kubeconfig`
- `client-secret`
- `password`
- `BEGIN PRIVATE KEY`
- `kubectl apply`
- `kubectl patch`
- `helm upgrade`
- `frontend-next`
- `React`
- `tsx`
- `css`
- `:443`
- `port: 443`
- `:8888`
- `port: 8888`
- `device_add`
- `qom-set`
- `system_powerdown`
- `shutdown`
- `reboot`
- `cpu-add`
- `object-add`
- `object-del`
- `migrate`

Result:

- no secret leakage found
- no frontend files changed
- no Dashboard files changed
- no deploy command execution introduced
- forbidden QMP command strings appear only in denylist/docs/tests as rejected commands
