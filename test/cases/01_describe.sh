#!/bin/bash

SCRIPT_ROOT_DIR=$1
TEST_NAME=${0:%.*}

"${SCRIPT_ROOT_DIR}"/edy d -t User > 01_actual.json
# shellcheck disable=SC2181
if [ $? -ne 0 ]; then echo "ERROR: ${TEST_NAME}"; exit 1; fi
