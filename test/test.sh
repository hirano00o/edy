#!/bin/bash

SCRIPT_DIR=$(dirname "$0")

setup () {
  echo "Setup integration test."
  echo ""
  docker run --rm --name integration_test -d -p 8000:8000 amazon/dynamodb-local >/dev/null
  sleep 1
}

tearDown () {
  rm -f "${SCRIPT_DIR}"/edy "${SCRIPT_DIR}"/cases/*.json
  docker stop integration_test >/dev/null
  echo ""
  echo "Finished integration test."
}

trap tearDown 0 1 2 3 15
setup

sh "${SCRIPT_DIR}/run.sh" "${SCRIPT_DIR}"
