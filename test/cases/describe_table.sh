#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb describe-table --table-name User --endpoint-url http://localhost:8000
CMD="edy d -t User --local 8000"

. "${SCRIPT_ROOT_DIR}"/helper.sh

run_such_query_helper