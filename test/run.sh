#!/bin/bash

SCRIPT_DIR=$(dirname "$0")

setup () {
  echo "Setup integration test."
  echo ""
  docker run --rm --name integration_test -d -p 8000:8000 amazon/dynamodb-local >/dev/null
  sleep 1

  if ! aws dynamodb create-table \
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
      --endpoint-url http://localhost:8000 >/dev/null;
  then
    exit 1
  fi

  if ! aws dynamodb batch-write-item --request-items file://"${SCRIPT_DIR}"/test_data.json --endpoint http://localhost:8000 >/dev/null;
  then
    exit 1
  fi

  go build -o "${SCRIPT_DIR}"/edy "${SCRIPT_DIR}"/../cmd/edy/main.go
}

tearDown () {
  rm -f "${SCRIPT_DIR}"/edy "${SCRIPT_DIR}"/cases/*.json
  docker stop integration_test >/dev/null
  echo ""
  echo "Finished integration test."
}

trap tearDown 0 1 2 3 15
setup

for f in "${SCRIPT_DIR}"/cases/*.sh
do
  if ! sh "${f}" "${SCRIPT_DIR}";
  then
    exit 1
  fi
done
