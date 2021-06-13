# edy

[![Release](https://github.com/hirano00o/edy/actions/workflows/release.yml/badge.svg)](https://github.com/hirano00o/edy/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/hirano00o/edy/branch/master/graph/badge.svg?token=VYSBN9XRON)](https://codecov.io/gh/hirano00o/edy)
[![Go Report Card](https://goreportcard.com/badge/github.com/hirano00o/edy)](https://goreportcard.com/report/github.com/hirano00o/edy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

edy is a command line interface designed to make DynamoDB easy to use.
When you query or scan on DynamoDB with AWS CLI, you have to write a lot of keys and values and options.
If you run it many times, it's very hard. Also, the results are deeply nested and difficult to read.
We are developing `edy` to make the results easier to handle and in order to reduce writing.
Currently, `scan`, `query` (and `describe-table`), `put`, `delete` are available. Options support filter and projection, GSI.
Other commands and options are under development.

# Installation
## Download Binaries

#### macOS

```shell
$ curl -O -L https://github.com/hirano00o/edy/releases/latest/download/edy_darwin_amd64.tar.gz
$ tar zxvf edy_darwin_amd64.tar.gz
$ mv edy /usr/local/bin/
$ chmod +x /usr/local/bin/edy
$ edy --help
```

#### Linux x86-64

```shell
$ curl -O -L https://github.com/hirano00o/edy/releases/latest/download/edy_linux_amd64.tar.gz
$ tar zxvf edy_linux_amd64.tar.gz
$ sudo mv edy /usr/local/bin/
$ chmod +x /usr/local/bin/edy
$ edy --help
```

## go get

```shell
$ go get github.com/hirano00o/edy/cmd/edy
```

# Usage
## Prerequisites

You need to set up AWS credentials in your environment. For example, performs `aws configure` or export environment variables such as `AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`.
Please see [Configuration and credential file settings - AWS Command Line Interface](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) for details.

## Overview

Currently, available commands are `describe`, `scan`, `query`, `put`, `delete`.

### describe

The `describe` command behaves similarly to `aws dynamodb describe-table`.

```console
$ edy describe --table-name User  # Shortened version: edy desc -t User
{
  "tableArn": "arn:aws:dynamodb:ddblocal:000000000000:table/User",
  "tableName": "User",
  "partitionKey": {
    "name": "ID",
    "type": "N"
  },
  "sortKey": {
    "name": "Name",
    "type": "S"
  },
  "gsi": [
    {
      "indexName": "EmailGSI",
      "partitionKey": {
        "name": "Email",
        "type": "S"
      }
    }
  ],
  "itemCount": 7
}
```

### scan

The `scan` command behaves similarly to `aws dynamodb scan`.  You can filter the results by using the `filter` option. Please see `edy s -h` for details.

```console
$ edy scan --table-name User --filter "not Birthplace,S exists and Age,N > 25" # Shortened version: edy s -t User -f "not Birthplace,S exists and Age,N > 25"
[
  {
    "Age": 26,
    "Email": "eve@example.com",
    "ID": 7,
    "Name": "Eve"
  }
]
```

### query

The `query` command behaves similarly to `aws dynamodb query`.

```console
$ edy query --table-name User --partition 1 # Shortened version: edy q -t User -p 1
[
  {
    "Age": 20,
    "Birthplace": "Arkansas",
    "Email": "alice@example.com",
    "ID": 1,
    "Name": "Alice"
  }
]
```

It can also specify the sort key and filter condition and so on.

```console
$ edy q -h
NAME:
   edy query - Query table

USAGE:
   edy query [command options] [arguments...]

OPTIONS:
   --table-name value, -t value    DynamoDB table name
   --region value, -r value        AWS region
   --profile value                 AWS profile name
   --local value                   Port number or full URL if you connect such as dynamodb-local and LocalStack.
                                   ex. --local 8000
   --partition value, -p value     The value of partition key
   --sort value, -s value          The value and condition of sort key.
                                   ex1. --sort "> 20"
                                   ex2. --sort "between 20 25"
                                   Available operator is =,<=,<,>=,>,between,begins_with
   --index value, --idx value      Global secondary index name
   --filter value, -f value        The condition if you use filter.
                                   ex. --filter "Age,N >= 20 and Email,S in alice@example.com bob@example.com or not Birthplace,S exists"
                                   Available operator is =,<=,<,>=,>,between,begins_with,exists,in,contains
   --projection value, --pj value  Identifies and retrieve the attributes that you want.
   --help, -h                      show help (default: false)
```

### put

The `put` command behaves similarly to `aws dynamodb put-item` or `aws dynamodb batch-write-item` (only PutRequest).
It creates a record by passing json to the `--item(-i)` option or using `--input-file` from file.  
Supported type is `S`, `N`, `SS`, `NS`, `M`, `L`, `BOOL`, `NULL`.

```console
# When put 1 item.
$ edy put --table-name User --item '{"ID":3, "Name":"Alice", "Interest":{"SNS":["Twitter","Facebook"]}}'
# Shortened version: edy p -i '{"ID":3, "Name":"Alice", "Interest":{"SNS":["Twitter","Facebook"]}}'
{
  "unprocessed": []
}

# When put some items, specify the items in an array.
$ edy put --table-name User --item '[{"ID":3, "Name":"Alice", "Interest":{"SNS":["Twitter","Facebook"]}}, {"ID":4, "Name":"Bob", "Interest":{"SNS":["Facebook"]}}]'
# Shortened version: edy p -i '{"ID":3, "Name":"Alice", "Interest":{"SNS":["Twitter","Facebook"]}}'
{
  "unprocessed": []
}
```

### delete

The `delete` command behaves similarly to `aws dynamodb batch-write-item` (only DeleteRequest).

```console
$ edy delete --table-name User --partition 1 --sort "Alice" # Shortened version: edy del -t User -p 1 -s Alice
{
  "unprocessed": []
}
```

If you want to delete multiple matching records at once, you should use `--input-file` option. The input file format is json.
You specify `partition` key or `partition` and `sort` keys.
Example is as follows.

```json
[
  {
    "partition": 1,
    "sort": "Alice"
  },
  {
    "partition": 2
  }
]
```

When you saved the above to delete.json, execute as follows.

```console
$ edy delete --table-name User --input-file delete.json # Shortened version: edy del -t User --input-file delete.json
{
  "unprocessed": []
}
```

## If use DynamoDB Local or LocalStack

You can connect to the local application such as DynamoDB Local and LocalStack by using `--local` option.
This option can specify the port number or full url, such as `--local 8000` or `--local http://localhost:8000`.

# Future works

* Implement other commands such as create and put and batch.
* Modify `scan` and` query` not to call `describe` internally every time.
