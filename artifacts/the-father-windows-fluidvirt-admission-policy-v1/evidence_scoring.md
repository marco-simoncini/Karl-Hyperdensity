# evidence_scoring

Scoring function: `EvaluateAdmissionEvidenceScore`.

Output:

- numeric score (0-100)
- evidence level (`insufficient`, `partial`, `dryrun-ready`, `future-apply-admissible`)
- missing evidence
- hard blockers
- soft unknowns

Signals include annotations, shell validity, identity completeness, QMP readiness/read-only, guest ACK, reboot/identity proofs, continuity proofs, migration/recreate absence, rollback/return readiness, memory driver safety, critical events absence, and evidence freshness.

Hard blockers (P0/P1) override numeric score and prevent future-apply admission.
