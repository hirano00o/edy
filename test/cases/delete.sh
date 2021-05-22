#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb delete-item --table-name User \
#   --key "{\"ID\":{\"N\":\"2\"}, \"Name\":{\"S\":\"Bob\"}}" \
#   --endpoint-url http://localhost:8000
CMD="edy del -t User -p 99 -s DELETE_USER --local 8000"
EXPECTED_DELETE_ITEM_COUNT=1

. "${SCRIPT_ROOT_DIR}"/helper.sh

TABLE_NAME="User"

run_such_delete_helper
