#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb batch-write-item \
#   --request-items file://${SCRIPT_ROOT_DIR}/cases/input/put_items_from_file.json \
#   --endpoint-url http://localhost:8000
CMD="edy p -t User -I ${SCRIPT_ROOT_DIR}/cases/input/put_items_from_file.json --local 8000"

. "${SCRIPT_ROOT_DIR}"/helper.sh

FILTER_CONDITION="ID,N = 15 or ID,N = 16"

run_such_put_helper
