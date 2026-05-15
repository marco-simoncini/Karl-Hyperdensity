# Risk Register

| Risk ID | Risk | Description | Impact | Likelihood | Mitigation |
|---|---|---|---|---|---|
| R1 | artifact bloat | Windows branch includes large artifact payload (499 artifact files in Hyperdensity diff) | High | High | selective import only; keep proofs as references; no bulk artifact merge |
| R2 | branch drift | Windows branches diverge from mainline (notably Dashboard stale) | High | High | integrate from fresh mainline branch with small PRs |
| R3 | product claim drift | technical preview language may drift into GA/production claims | High | Medium | enforce claim boundary doc + review checklist gates |
| R4 | duplicated contracts | parallel Windows docs may duplicate/contradict Linux contract semantics | Medium | Medium | contract alignment matrix + naming normalization before merge |
| R5 | stale dashboard UI | importing stale Dashboard branch can regress safety UX and baseline | High | High | never source UI from `The-Father-Windows`; build UI later from `The-Father` |
| R6 | secret hygiene | historical artifacts may accidentally include sensitive local material | Critical | Low-Medium | strict staging filters; no secret/local files; pre-commit scan |
| R7 | registry/delivery | delivery channels may fail or drift from reviewed artifact set | Medium | Medium | defer delivery work; gate via dedicated release milestone |
| R8 | runtime safety | accidental enablement of mutation paths/autonomous apply | Critical | Medium | keep executor planning-only; explicit blockers and kill switch enforcement |
| R9 | Windows proof overclaim | benchmark/proof data can be overstated as production guarantee | High | Medium | preserve evidence wording: preview-only, environment-scoped, non-GA |
