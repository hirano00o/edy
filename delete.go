package edy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/hirano00o/edy/client"
	"github.com/hirano00o/edy/model"
)

type dynamoDBValue struct {
	partitionValue string
	sortValue      string
}

func deleteItems(ctx context.Context, tableName string, items []*dynamoDBValue) (map[string]interface{}, error) {
	table, err := describeTable(ctx, tableName)
	if err != nil {
		return nil, err
	}
	cli := ctx.Value(newClientKey).(client.DynamoDB)

	partitionKeyName, partitionKeyType := table.PartitionKey.Name, table.PartitionKey.Type

	deleteRequest := make([]types.WriteRequest, len(items))
	for i := range items {
		m := make(map[string]types.AttributeValue)

		// PartitionKey condition
		m[partitionKeyName] = partitionKeyType.ConvertValueMember(items[i].partitionValue)

		// SortKey condition
		if len(items[i].sortValue) != 0 && table.SortKey != nil {
			sortKeyName, sortKeyType := table.SortKey.Name, table.SortKey.Type
			m[sortKeyName] = sortKeyType.ConvertValueMember(items[i].sortValue)
		}

		deleteRequest[i].DeleteRequest = &types.DeleteRequest{
			Key: m,
		}
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: deleteRequest,
		},
	}
	var res *dynamodb.BatchWriteItemOutput
	var unprocessedCount int
	for i := 0; i < 1+model.RetryMax; i++ {
		res, err = cli.BatchWriteItem(ctx, input)
		if err != nil {
			return nil, err
		}
		unprocessedCount = len(res.UnprocessedItems[tableName])
		if unprocessedCount == 0 {
			break
		}
		input.RequestItems = res.UnprocessedItems
	}
	if unprocessedCount > 0 {
		return map[string]interface{}{
			"unprocessed": res.UnprocessedItems[tableName],
		}, nil
	}

	return map[string]interface{}{"unprocessed": []string{}}, nil
}

func analyseDeleteRequestItem(requestJSONStr string) ([]*dynamoDBValue, error) {
	jsonItem, err := parseJSON(requestJSONStr)
	if err != nil {
		return nil, err
	}

	switch j := jsonItem.(type) {
	case []map[string]interface{}:
		deleteItems := make([]*dynamoDBValue, len(j))
		for i := range j {
			deleteItems[i], err = getValueFromRequestItems(j[i])
			if err != nil {
				return nil, err
			}
			if len(deleteItems[i].partitionValue) == 0 {
				return nil, fmt.Errorf("required partition value: %v", jsonItem)
			}
		}
		return deleteItems, nil
	case map[string]interface{}:
		deleteItems := make([]*dynamoDBValue, 1)
		deleteItems[0], err = getValueFromRequestItems(j)
		if err != nil {
			return nil, err
		}
		if len(deleteItems[0].partitionValue) == 0 {
			return nil, fmt.Errorf("required partition value: %v", jsonItem)
		}

		return deleteItems, nil
	default:
		return nil, fmt.Errorf("unknown error: %v", jsonItem)
	}
}

func getValueFromRequestItems(jsonItem map[string]interface{}) (*dynamoDBValue, error) {
	var partitionValue, sortValue string
	for k := range jsonItem {
		var v string
		// Primary key is allowed string, number, byte.
		if reflect.TypeOf(jsonItem[k]).Kind() == reflect.Float64 {
			v = strconv.FormatFloat(jsonItem[k].(float64), 'f', -1, 64)
		} else {
			v = jsonItem[k].(string)
		}
		switch {
		case k == "partition":
			partitionValue = v
		case k == "sort":
			sortValue = v
		default:
			return nil, fmt.Errorf("unknown key specified: %v", k)
		}
	}
	return &dynamoDBValue{
		partitionValue: partitionValue,
		sortValue:      sortValue,
	}, nil
}

func (i *Instance) Delete(
	ctx context.Context,
	w io.Writer,
	tableName,
	partitionValue,
	sortValue,
	fileName string,
	f func(string) (string, error),
) error {
	var items []*dynamoDBValue
	switch {
	case len(partitionValue) == 0 && len(fileName) == 0:
		return fmt.Errorf("required either --partition or --input-file option")
	case len(partitionValue) != 0 && len(fileName) != 0:
		return fmt.Errorf("use either --partition or --input-file option")
	case len(fileName) != 0:
		var err error
		strJSONItems, err := f(fileName)
		if err != nil {
			return err
		}
		items, err = analyseDeleteRequestItem(strJSONItems)
		if err != nil {
			return err
		}
	case len(partitionValue) != 0:
		items = append(items, &dynamoDBValue{
			partitionValue: partitionValue,
			sortValue:      sortValue,
		})
	}

	cli := i.NewClient.CreateInstance()
	ctx = context.WithValue(ctx, newClientKey, cli)

	res, err := deleteItems(ctx, tableName, items)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(res, "", strings.Repeat(" ", 2))
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", string(b))

	return nil
}
