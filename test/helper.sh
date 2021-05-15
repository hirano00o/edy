#!/bin/bash

run_such_query_helper() {
  CASE_DIR=${SCRIPT_ROOT_DIR}/cases
  EXPECTED_FILE=${SCRIPT_ROOT_DIR}/cases/expected/${TEST_NAME}.json
  if ! eval "${SCRIPT_ROOT_DIR}/${CMD}" > "${CASE_DIR}/${TEST_NAME}"_actual.json;
  then
    printf "\033[31m%s\033[m:\t%s\n" "=== FAILED" "${TEST_NAME} failed to execute ${CMD}"
    exit 1
  fi
  if ! diff -u "${CASE_DIR}/${TEST_NAME}"_actual.json "${EXPECTED_FILE}" > tmp.diff;
  then
    printf "\033[31m%s\033[m:\t%s\n" "=== FAILED" "${TEST_NAME}"
    sed -e "s@${CASE_DIR}/${TEST_NAME}_actual.json@actual@" -e "s@${EXPECTED_FILE}@expected@" tmp.diff
    rm tmp.diff
    exit 1
  fi
  rm tmp.diff

  printf "\033[32m%s\033[m:\t%s\n" "--- PASSED" "${TEST_NAME}"
}
