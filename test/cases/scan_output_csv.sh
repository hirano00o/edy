#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb scan --table-name User \
#   --filter-expression "ID = :id1 or ID = :id2" \
#   --expression-attribute-values "{\":id1\":{\"N\":\"12\"}, \":id2\":{\"N\":\"13\"}}" \
#   --endpoint-url http://localhost:8000
CMD="edy s -t User -f \"ID,N = 12 or ID,N = 13\" -o csv --local 8000"

. "${SCRIPT_ROOT_DIR}"/helper.sh

run_such_query_helper
