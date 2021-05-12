#!/bin/bash

SCRIPT_DIR=$(dirname $0)

function setup() {
  docker run --rm --name integration_test -d -p 8000:8000 amazon/dynamodb-local
  sleep 3

  aws dynamodb create-table \
    --table-name User \
    --attribute-definitions AttributeName=ID,AttributeType=N AttributeName=Name,AttributeType=S AttributeName=Email,AttributeType=S \
    --key-schema AttributeName=ID,KeyType=HASH AttributeName=Name,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --global-secondary-indexes \
      "[
        {
          \"IndexName\": \"EmailIndex\",
          \"KeySchema\": [{\"AttributeName\": \"Email\",\"KeyType\": \"HASH\"}],
          \"Projection\": {\"ProjectionType\":\"ALL\"},
          \"ProvisionedThroughput\": {\"ReadCapacityUnits\": 5, \"WriteCapacityUnits\": 5}
        }
      ]" \
      --endpoint-url http://localhost:8000 >/dev/null
  # shellcheck disable=SC2181
  if [ $? -ne 0 ]; then exit 1; fi

  aws dynamodb batch-write-item --request-items file://"${SCRIPT_DIR}"/test_data.json --endpoint http://localhost:8000
  # shellcheck disable=SC2181
  if [ $? -ne 0 ]; then exit 1; fi

  go build -o "${SCRIPT_DIR/}"edy "${SCRIPT_DIR}"/../cmd/edy/main.go
}

function tearDown() {
  rm -f "${SCRIPT_DIR}"/edy
  docker stop integration_test
}

trap tearDown 0 1 2 3 15
setup

for f in ${SCRIPT_DIR}/cases/*
do
  sh "${f}" "${SCRIPT_DIR}"
done
