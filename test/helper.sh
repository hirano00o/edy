#!/bin/bash

# Declare SCRIPT_ROOT_DIR, CMD, TEST_NAME in advance.

run_such_query_helper() {
  CASE_DIR=${SCRIPT_ROOT_DIR}/cases
  EXPECTED_FILE=${SCRIPT_ROOT_DIR}/cases/expected/${TEST_NAME}.json
  if [ ! -e "${EXPECTED_FILE}" ];
  then
    EXPECTED_FILE=${SCRIPT_ROOT_DIR}/cases/expected/${TEST_NAME}.csv
  fi
  if ! eval "${SCRIPT_ROOT_DIR}/${CMD}" > "${CASE_DIR}/actual/${TEST_NAME}";
  then
    printf "\033[31m%s\033[m:\t%s\n" "=== FAILED" "${TEST_NAME} failed to execute ${CMD}"
    exit 1
  fi
  if ! diff -u "${CASE_DIR}/actual/${TEST_NAME}" "${EXPECTED_FILE}" > tmp.diff;
  then
    printf "\033[31m%s\033[m:\t%s\n" "=== FAILED" "${TEST_NAME}"
    sed -e "s@${CASE_DIR}/actual/${TEST_NAME}@actual@" -e "s@${EXPECTED_FILE}@expected@" tmp.diff
    rm tmp.diff
    exit 1
  fi
  rm tmp.diff

  printf "\033[32m%s\033[m:\t%s\n" "--- PASSED" "${TEST_NAME}"
}

run_such_put_helper() {
  CASE_DIR=${SCRIPT_ROOT_DIR}/cases
  EXPECTED_FILE=${SCRIPT_ROOT_DIR}/cases/expected/${TEST_NAME}.json
  if [ ! -e "${EXPECTED_FILE}" ];
  then
    EXPECTED_FILE=${SCRIPT_ROOT_DIR}/cases/expected/${TEST_NAME}.csv
  fi
  if ! eval "${SCRIPT_ROOT_DIR}/${CMD}" > "${CASE_DIR}/actual/${TEST_NAME}";
  then
    printf "\033[31m%s\033[m:\t%s\n" "=== FAILED" "${TEST_NAME} failed to execute ${CMD}"
    exit 1
  fi
  if ! printf "{\n  \"unprocessed\": []\n}\n" | diff -u "${CASE_DIR}/actual/${TEST_NAME}" - > tmp.diff;
  then
    printf "\033[31m%s\033[m:\t%s\n" "=== FAILED" "${TEST_NAME}"
    sed -e "s@${CASE_DIR}/actual/${TEST_NAME}@actual@" -e "s@+++ -@+++ expected@" tmp.diff
    rm tmp.diff
    exit 1
  fi

  "${SCRIPT_ROOT_DIR}/edy" s -t User -f "${FILTER_CONDITION}" --local 8000 > "${CASE_DIR}/actual/${TEST_NAME}";

  if ! diff -u "${CASE_DIR}/actual/${TEST_NAME}" "${EXPECTED_FILE}" > tmp.diff;
  then
    printf "\033[31m%s\033[m:\t%s\n" "=== FAILED" "${TEST_NAME}"
    sed -e "s@${CASE_DIR}/actual/${TEST_NAME}@actual@" -e "s@${EXPECTED_FILE}@expected@" tmp.diff
    rm tmp.diff
    exit 1
  fi
  rm tmp.diff

  printf "\033[32m%s\033[m:\t%s\n" "--- PASSED" "${TEST_NAME}"
}

run_such_delete_helper() {
  CASE_DIR=${SCRIPT_ROOT_DIR}/cases

  BEFORE_ITEM_COUNT=$(aws dynamodb scan --table-name "${TABLE_NAME}" \
    --endpoint-url http://localhost:8000 | jq ".Items | length")

  if ! eval "${SCRIPT_ROOT_DIR}/${CMD}" > "${CASE_DIR}/actual/${TEST_NAME}";
  then
    printf "\033[31m%s\033[m:\t%s\n" "=== FAILED" "${TEST_NAME} failed to execute ${CMD}"
    exit 1
  fi
  if ! printf "{\n  \"unprocessed\": []\n}\n" | diff -u "${CASE_DIR}/actual/${TEST_NAME}" - > tmp.diff;
  then
    printf "\033[31m%s\033[m:\t%s\n" "=== FAILED" "${TEST_NAME}"
    sed -e "s@${CASE_DIR}/actual/${TEST_NAME}@actual@" -e "s@+++ -@+++ expected@" tmp.diff
    rm tmp.diff
    exit 1
  fi

  AFTER_ITEM_COUNT=$(aws dynamodb scan --table-name "${TABLE_NAME}" \
    --endpoint-url http://localhost:8000 | jq ".Items | length")

  ACTUAL_COUNT=$(expr "${BEFORE_ITEM_COUNT}" - "${AFTER_ITEM_COUNT}")
  if [ "${ACTUAL_COUNT}" != "${EXPECTED_DELETE_ITEM_COUNT}" ]; then
    printf "\033[31m%s\033[m:\t%s\n" "=== FAILED" "${TEST_NAME}"
    printf "got %s, want %s" "${EXPECTED_DELETE_ITEM_COUNT}" "${ACTUAL_COUNT}"
    exit 1
  fi

  printf "\033[32m%s\033[m:\t%s\n" "--- PASSED" "${TEST_NAME}"
}
