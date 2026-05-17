# shellcheck shell=bash
# KHR-BU: shared validation mode helpers (source only).
khr_validation_root() {
  cd "$(dirname "${BASH_SOURCE[1]:-${BASH_SOURCE[0]}}")/.." && pwd
}

khr_live_validate_enabled() {
  [[ "${KHR_LIVE_VALIDATE:-}" == "1" ]]
}

khr_offline_validate_mode() {
  ! khr_live_validate_enabled
}
