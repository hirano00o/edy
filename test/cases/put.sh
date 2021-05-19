#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb put-item --table-name User \
#   --item "{\"ID\":{\"N\":\"11\"}, \"Name\":{\"S\":\"Isaac\"}, \"Email\":{\"S\":\"isaac@example.com\"}, \"Age\":{\"N\":\"30\"}, \"Birthday\":{\"M\":{\"Year\":{\"N\":\"1991\"}, \"Month\":{\"N\":\"9\"}, \"Day\":{\"N\":\"30\"}}}}" \
#   --endpoint-url http://localhost:8000
CMD="edy p -t User -i '{\"ID\":{\"N\":\"11\"}, \"Name\":{\"S\":\"Isaac\"}, \"Email\":{\"S\":\"isaac@example.com\"}, \"Age\":{\"N\":\"30\"}, \"Birthday\":{\"M\":{\"Year\":{\"N\":\"1991\"}, \"Month\":{\"N\":\"9\"}, \"Day\":{\"N\":\"30\"}}}}' --local 8080"

. "${SCRIPT_ROOT_DIR}"/helper.sh

run_such_put_helper
