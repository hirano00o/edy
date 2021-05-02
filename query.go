package edy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/hirano00o/edy/client"
	"github.com/hirano00o/edy/model"
)

func query(ctx context.Context, tableName string, partitionValue string) ([]map[string]interface{}, error) {
	table, err := describeTable(ctx, tableName)
	if err != nil {
		return nil, err
	}
	cli := ctx.Value(newClientKey).(client.DynamoDB)

	var p expression.ValueBuilder
	switch table.PartitionKeyType {
	case model.N:
		pInt, err := strconv.Atoi(partitionValue)
		if err != nil {
			return nil, err
		}
		p = expression.Value(pInt)
	case model.S:
		p = expression.Value(partitionValue)
	case model.B:
		p = expression.Value(bytes.NewBufferString(partitionValue).Bytes())
	default:
		p = expression.Value(partitionValue)
	}

	// PartitionKey condition
	condition := expression.KeyEqual(expression.Key(table.PartitionKeyName), p)

	builder := expression.NewBuilder().WithKeyCondition(condition)
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	resMap := make([]map[string]interface{}, 0, table.ItemCount)
	paginator := dynamodb.NewQueryPaginator(cli, input)
	for paginator.HasMorePages() {
		res, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		v := make([]map[string]interface{}, 0, 25)
		err = attributevalue.UnmarshalListOfMaps(res.Items, &v)
		if err != nil {
			return nil, err
		}
		resMap = append(resMap, v...)
	}

	return resMap, nil
}

func (i *Instance) Query(
	ctx context.Context,
	w io.Writer,
	tableName string,
	partitionValue,
	sortValue string,
	filterCondition string,
) error {
	cli := i.NewClient.CreateInstance()
	ctx = context.WithValue(ctx, newClientKey, cli)

	res, err := query(ctx, tableName, partitionValue)
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
