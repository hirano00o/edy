#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=$(basename "$0" | sed "s/\..*//")

# aws dynamodb batch-write-item \
#   --request-items "{\"User\":[{\"PutRequest\":{\"Item\":{\"ID\":{\"N\":\"12\"}, \"Name\":{\"S\":\"Ivan\"}, \"Email\":{\"S\":\"isaac@example.com\"}, \"Age\":{\"N\":\"30\"}, \"Birthday\":{\"M\":{\"Year\":{\"N\":\"1991\"}, \"Month\":{\"N\":\"9\"}, \"Day\":{\"N\":\"30\"}}}}}},
#   "{\"PutRequest\":{\"Item\":{\"ID\":{\"N\":\"12\"}, \"Name\":{\"S\":\"Ivan\"}, \"Email\":{\"S\":\"isaac@example.com\"}, \"Age\":{\"N\":\"30\"}, \"Birthday\":{\"M\":{\"Year\":{\"N\":\"1991\"}, \"Month\":{\"N\":\"9\"}, \"Day\":{\"N\":\"30\"}}}}}}]}" \
#   --endpoint-url http://localhost:8000
CMD="edy p -t User -i '[{\"ID\":12, \"Name\":\"Ivan\", \"Email\":\"ivan@example.com\", \"Age\":32, \"Birthday\":{\"Year\":1989, \"Month\":3, \"Day\":30}},{\"ID\":13, \"Name\":\"Justin\", \"Email\":\"justin@example.com\", \"Age\":32, \"Birthday\":{\"Year\":1989, \"Month\":2, \"Day\":28}}]' --local 8000"

. "${SCRIPT_ROOT_DIR}"/helper.sh

FILTER_CONDITION="ID,N = 12 or ID,N = 13"

run_such_put_helper
