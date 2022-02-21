#! /bin/bash

set -o errexit
set -o pipefail
set -o nounset

if [[ -z "${SPEC_PROJ_DIR:-}" ]]; then
    echo "FATAL: SPEC_PROJ_DIR is not defined, cannot run :/"
    exit 1
fi

if [[ "${PMIX_RANK}" == 0 ]]; then
    date
    set -x
fi

####################################################################################

EVENT_DIR="${SPEC_PROJ_DIR}/run0001"
OF_DIR="${EVENT_DIR}/OUTPUT_FILES"
DATA_DIR="${EVENT_DIR}/DATA"

####################################################################################

echo "Byebye"

exit 0

cd "$SPEC_PROJ_DIR"

./bin/xgenerate_databases

if [[ "${PMIX_RANK}" == 0 ]]; then
    cp "$DATA_DIR/Par_file" "$DATA_DIR/gen_db.Par_file"
    echo "Done..., results saved in directory: $DB_MPI_DIR"

    date
fi
