# Integration test

## How to run in your local machine

```shell
sh test.sh
```

## How to add the test

If you add commands or options, you should to add the integration test case.
Place the `TEST_CASE_NAME.sh` under the `cases` directory as a name that represents the test case.
You input the command to the `CMD` variable, and leave the `aws` command that has the same meaning as that command in the comments.
Then, call `run_such_query_helper` if the command receives the result such as json.
If such as put command, use `run_such_put_helper`.
Also, if you do not have enough test data, please add it at any time.

Place the expected result in the `expected` directory as `TEST_CASE_NAME.json`.

For example, if you test the query command, file places is as follows.

```shell
.
├── README.md
├── cases
│   ├── expected
│   │   └── query.json
│   └── query.sh
├── helper.sh
├── run.sh
├── test.sh
└── test_data.json
```
