# risk_register

- risk: stale Dashboard lane accidentally proposed as merge source
  - mitigation: explicit `dashboardMergeAllowed=false` and forbidden scope listing

- risk: runtime readiness misinterpretation from replay/boundary artifacts
  - mitigation: strict claim and safety attestation fields set to false

- risk: direct merge pressure from `The-Father-Windows`
  - mitigation: explicit `directMergeAllowed=false`, selective PR sequence only

- risk: OS-ISO packaging pressure before runtime MVP gate
  - mitigation: explicit `osIsoMergeAllowed=false` and deferred decision artifact

- risk: operational drift between Hyperdensity and Inventory integrations
  - mitigation: lock merge to verified branch heads and operational outcomes in this gate
