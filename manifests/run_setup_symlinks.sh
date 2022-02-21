#! /bin/bash

set -o errexit
set -o pipefail
set -o nounset
set -x

if [[ -z "${SPEC_PROJ_DIR:-}" ]]; then
    echo "FATAL: SPEC_PROJ_DIR is not defined, cannot run :/"
    exit 1
fi

date

####################################################################################

BASE_EVENT_DIR="${SPEC_PROJ_DIR}/run0001"
BASE_OF_DIR="${BASE_EVENT_DIR}/OUTPUT_FILES"
BASE_DB_MPI_DIR="${BASE_OF_DIR}/DATABASES_MPI"

####################################################################################

echo
echo "Creating Database links"
echo

cat <<EOF
####################################################
#
# First: make db and mesh links for each event directory
#
####################################################
EOF

for dir in "${SPEC_PROJ_DIR}"/run[0-9]*/; do
  RDIR=$(echo "$dir" | rev | cut -d '/' -f2 | rev)
  if [[ "$RDIR" == "run0001" ]]; then
    continue
  fi

  LOC_OF_DIR="$dir/OUTPUT_FILES";
  LOC_DB_MPI_DIR="$LOC_OF_DIR/DATABASES_MPI";

  echo
  echo "## remove old mesh and db links #FIXME: shouldn't cleanup here"

  rm "$LOC_DB_MPI_DIR"/proc*external*.bin "$LOC_DB_MPI_DIR"/proc*Database;

  echo
  echo "## create new mesh and db links"

  EM_CMD="ls $BASE_DB_MPI_DIR/proc*external*.bin"
  DB_CMD="ls $BASE_DB_MPI_DIR/proc*Database"
  XARG_CMD="xargs -t -P0 -I {} ln -srf $BASE_DB_MPI_DIR/{} $LOC_DB_MPI_DIR/{}"

  ${EM_CMD} | rev | cut -d '/' -f1 | rev | ${XARG_CMD}
  ${DB_CMD} | rev | cut -d '/' -f1 | rev | ${XARG_CMD}

  echo
  echo "## make mesh header links"

  VM_CMD="ls $BASE_OF_DIR/values_from_mesher.h"
  SM_CMD="ls $BASE_OF_DIR/surface_from_mesher.h"
  XARG_CMD="xargs -t -P0 -I {} ln -srf $BASE_OF_DIR/{} $LOC_OF_DIR/{}"
  ${VM_CMD} | rev | cut -d '/' -f1 | rev | ${XARG_CMD}
  ${SM_CMD} | rev | cut -d '/' -f1 | rev | ${XARG_CMD}
  echo
done

cat <<EOF
##########################################################
#
# Second: make mesh links for topo dir and model updates
#
##########################################################
EOF

EVENT_DIR="${SPEC_PROJ_DIR}/run0001"
OF_DIR="${EVENT_DIR}/OUTPUT_FILES"
DB_MPI_DIR="${OF_DIR}/DATABASES_MPI"
TOPO_DIR=${EVENT_DIR}/topo

MESH_CMD="ls $DB_MPI_DIR/proc*external_mesh.bin"
XARG_CMD="xargs -t -P0 -I {} ln -srf $DB_MPI_DIR/{} $TOPO_DIR/{}"

${MESH_CMD} | rev | cut -d '/' -f1 | rev | ${XARG_CMD}

cat <<EOF
##########################################################
#
# Fourth: create and set the current_iteration.sh
#
##########################################################
EOF


echo "Done..., see results in directory: $EVENT_DIR/"
date

echo
echo "Finished Creating Database links";
