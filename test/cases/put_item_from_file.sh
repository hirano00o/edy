#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb put-item --table-name User \
#   --item file://${SCRIPT_ROOT_DIR}/cases/input/put_item_from_file.json \
#   --endpoint-url http://localhost:8000
CMD="edy p -t User -I ${SCRIPT_ROOT_DIR}/cases/input/put_item_from_file.json --local 8000"

. "${SCRIPT_ROOT_DIR}"/helper.sh

FILTER_CONDITION="ID,N = 14"

run_such_put_helper
