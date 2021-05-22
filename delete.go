package edy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/hirano00o/edy/client"
)

func deleteItems(ctx context.Context, tableName, partitionValue, sortValue string) (map[string]interface{}, error) {
	if len(partitionValue) == 0 {
		return nil, fmt.Errorf("required partition value")
	}
	table, err := describeTable(ctx, tableName)
	if err != nil {
		return nil, err
	}
	cli := ctx.Value(newClientKey).(client.DynamoDB)

	partitionKeyName, partitionKeyType := table.PartitionKey.Name, table.PartitionKey.Type

	m := make(map[string]types.AttributeValue)

	// PartitionKey condition
	m[partitionKeyName] = partitionKeyType.ConvertValueMember(partitionValue)

	// SortKey condition
	if len(sortValue) != 0 && table.SortKey != nil {
		sortKeyName, sortKeyType := table.SortKey.Name, table.SortKey.Type
		m[sortKeyName] = sortKeyType.ConvertValueMember(sortValue)
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       m,
	}

	_, err = cli.DeleteItem(ctx, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"unprocessed": []string{}}, nil
}

func (i *Instance) Delete(ctx context.Context, w io.Writer, tableName, partitionValue, sortValue string) error {
	cli := i.NewClient.CreateInstance()
	ctx = context.WithValue(ctx, newClientKey, cli)

	res, err := deleteItems(ctx, tableName, partitionValue, sortValue)
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
