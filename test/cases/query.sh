#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb query --table-name User --key-condition-expression ID=:id \
#   --expression-attribute-values "{\":id\":{\"N\":\"1\"}}" --endpoint-url http://localhost:8000
CMD="edy q -t User -p 1 --local 8000"

. "${SCRIPT_ROOT_DIR}"/helper.sh

run_such_query_helper
