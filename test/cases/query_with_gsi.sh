#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb query --table-name User --index EmailIndex --key-condition-expression Email=:email \
#   --expression-attribute-values "{\":email\":{\"S\":\"charlie@example.com\"}}" --endpoint-url http://localhost:8000
CMD="edy q -t User -p charlie@example.com --idx EmailIndex --local 8000"

. "${SCRIPT_ROOT_DIR}"/helper.sh

run_such_query_helper