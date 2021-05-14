#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb scan --table-name User \
#   --projection-expression "Email,Address.City"
#   --filter-expression "contains(Interest.SNS,:sns) and attribute_exists(Interest.Video)" \
#   --expression-attribute-values "{\":sns\":{\"S\":\"Twitter\"}}" \
#   --endpoint-url http://localhost:8000
CMD="edy s -t User -f \"Interest.SNS,S contains Twitter and Interest.Video,SS exists\" \
  --pj \"Email Address.City\" --local 8000"

. "${SCRIPT_ROOT_DIR}"/helper.sh

run_such_query_helper