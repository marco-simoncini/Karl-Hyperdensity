# Bundle Append Write-Back

Replay CLI now supports deterministic append write-back:

- `-append-bundle`
- `-append-bundle-in <bundle.json>`
- `-bundle-out <output.json>`

Behavior:

- validates existing chain before append
- rejects broken chain
- rejects duplicate runs
- appends run with linked `previousRunHash`
- writes updated bundle index deterministically with fixed evaluation time
