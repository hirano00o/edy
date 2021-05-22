#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb put-item --table-name User \
#   --item "{\"ID\":{\"N\":\"11\"}, \"Name\":{\"S\":\"Isaac\"}, \"Email\":{\"S\":\"isaac@example.com\"}, \"Age\":{\"N\":\"30\"}, \"Birthday\":{\"M\":{\"Year\":{\"N\":\"1991\"}, \"Month\":{\"N\":\"9\"}, \"Day\":{\"N\":\"30\"}}}}" \
#   --endpoint-url http://localhost:8000
CMD="edy p -t User -i '{\"ID\":11, \"Name\":\"Isaac\", \"Email\":\"isaac@example.com\", \"Age\":30, \"Birthday\":{\"Year\":1991, \"Month\":9, \"Day\":30}}' --local 8000"

. "${SCRIPT_ROOT_DIR}"/helper.sh

FILTER_CONDITION="ID,N = 11"

run_such_put_helper
