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

EVENT_DIR="${SPEC_PROJ_DIR}/run0001"
OF_DIR="${EVENT_DIR}/OUTPUT_FILES"
DATA_DIR="${EVENT_DIR}/DATA"

# nparts = number of partitions
NPROC=$APP_SPEC_EXEC_NPROC
# input_directory = directory containing mesh files mesh_file,nodes_coords_file,..
MESH_DIR="${SPEC_PROJ_DIR}/MESH-default"
# output_directory = directory for output files proc***_Databases
DB_MPI_DIR="${OF_DIR}/DATABASES_MPI"

####################################################################################

echo "Byebye"

exit 0

cd "$EVENT_DIR"

echo
echo "  decomposing mesh..."
echo


# checks exit code
if ! /app/bin/xdecompose_mesh "$NPROC" "$MESH_DIR" "$DB_MPI_DIR"
then
  echo
  echo "There was an error running xdecompose_mesh"
  exit 1
else
  echo
  echo "xdecompose_mesh completed successfully";
fi

##*****************************************##

echo
cp "$DATA_DIR/Par_file" "$DATA_DIR/decomp_mesh.Par_file"
echo "Done..., results saved in results in directory: $DB_MPI_DIR"

date
