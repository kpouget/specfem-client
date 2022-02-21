#! /bin/bash

set -o errexit
set -o pipefail
set -o nounset
set -x

if [[ -z "${SPEC_PROJ_DIR:-}" ]]; then
    echo "FATAL: SPEC_PROJ_DIR is not defined, cannot run :/"
    exit 1
fi

touch "${DATA_DIR}/SOLVER_${PMIX_RANK}"
