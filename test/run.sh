#!/bin/bash

SCRIPT_DIR=${1:-$(dirname "$0")}

initialise() {
  mkdir -p "${SCRIPT_DIR}"/cases/actual
  if ! aws dynamodb create-table \
    --region ap-northeast-1 \
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

  if ! aws dynamodb batch-write-item --region ap-northeast-1 --request-items file://"${SCRIPT_DIR}"/test_data.json --endpoint http://localhost:8000 >/dev/null;
  then
    exit 1
  fi

  go build -o "${SCRIPT_DIR}"/edy "${SCRIPT_DIR}"/../cmd/edy/main.go
}

initialise

for f in "${SCRIPT_DIR}"/cases/*.sh
do
  if ! sh "${f}" "${SCRIPT_DIR}";
  then
    exit 1
  fi
done
