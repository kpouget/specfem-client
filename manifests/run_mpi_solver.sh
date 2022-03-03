#! /bin/bash

set -o errexit
set -o pipefail
set -o nounset
set -x

if [[ -z "${SPEC_PROJ_DIR:-}" ]]; then
    echo "FATAL: SPEC_PROJ_DIR is not defined, cannot run :/"
    exit 1
fi

####################################################################################

EVENT_ID=run0001

EVENT_DIR="${SPEC_PROJ_DIR}/${EVENT_ID}"
OF_DIR="${EVENT_DIR}/OUTPUT_FILES"
DATA_DIR="${EVENT_DIR}/DATA"
SYN_DIR="${EVENT_DIR}/SYN"

####################################################################################

if [[ "${PMIX_RANK}" == 0 ]]; then
    date
    echo "Running solver for ${EVENT_DIR} using ${SPEC_NPROC} processors..."
    set -x
fi

####################################################################################

# set NPROC in Par_file

PAR_FILE="${DATA_DIR}/Par_file"

sed "s/^NPROC.*/NPROC = ${SPEC_NPROC} /g" -i "${PAR_FILE}"

grep 'NPROC' "${PAR_FILE}" | head -1;

####################################################################################

cd "$EVENT_DIR"

./utils/change_simulation_type.pl -F # missing file

# checks exit code
if ! ./bin/xspecfem3D; then
  echo
  echo "There was an error running xspecfem3D"
  exit 1
else
  if [[ "${PMIX_RANK}" != 0 ]]; then
      exit 0
  fi
  echo
  echo "xspecfem3D completed successfully"
fi

#
# /!\ only RANK 0 runs here
#

echo "move data in OUTPUT_FILES to SYN"

mv "${OF_DIR}"/*.semd "${SYN_DIR}"

##*****************************************##

# Save Par_file for record keeping

ls "${PAR_FILE}"

cp "${PAR_FILE}" "${DATA_DIR}/fwd_solver.Par_file"

echo "Done, see results in directory: OUTPUT_FILES/"

date
