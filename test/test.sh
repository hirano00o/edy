#!/bin/bash

SCRIPT_DIR=$(dirname "$0")

setup () {
  echo "Setup integration test."
  echo ""
  docker run --rm --name integration_test -d -p 8000:8000 amazon/dynamodb-local >/dev/null
  sleep 1
}

tearDown () {
  rm -rf "${SCRIPT_DIR}"/edy "${SCRIPT_DIR}"/cases/actual
  docker stop integration_test >/dev/null
  echo ""
  echo "Finished integration test."
}

UNAME="$(uname)"
if [ "${UNAME}" == "Darwin" -o "${UNAME}" == "Linux" ]; then
  if [ "$(which jq)" == "" ]; then
    echo "Please install jq command"
    exit 1
  fi
else
  if [ "$(where jq)" == "" ]; then
    echo "Please install jq command"
    exit 1
  fi
fi

trap tearDown 0 1 2 3 15
setup

sh "${SCRIPT_DIR}/run.sh" "${SCRIPT_DIR}"
