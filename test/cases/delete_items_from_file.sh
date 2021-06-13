#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb batch-write-item \
#   --request-items file://${SCRIPT_ROOT_DIR}/cases/input/delete_items_from_file.json \
#   --endpoint-url http://localhost:8000
CMD="edy del -t User -I ${SCRIPT_ROOT_DIR}/cases/input/delete_items_from_file.json --local 8000"
EXPECTED_DELETE_ITEM_COUNT=2

. "${SCRIPT_ROOT_DIR}"/helper.sh

TABLE_NAME="User"

run_such_delete_helper
