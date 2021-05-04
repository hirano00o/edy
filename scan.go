package edy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/hirano00o/edy/client"
)

func (i *Instance) Scan(
	ctx context.Context,
	w io.Writer,
	tableName string,
	filterCondition string,
) error {
	cli := i.NewClient.CreateInstance()
	ctx = context.WithValue(ctx, newClientKey, cli)

	res, err := scan(ctx, tableName, filterCondition)
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

func scan(ctx context.Context, tableName string, filterCondition string) ([]map[string]interface{}, error) {
	table, err := describeTable(ctx, tableName)
	if err != nil {
		return nil, err
	}
	cli := ctx.Value(newClientKey).(client.DynamoDB)

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Filter condition
	if len(filterCondition) != 0 {
		c, err := analyseFilterCondition(filterCondition)
		if err != nil {
			return nil, err
		}
		builder := expression.NewBuilder().WithCondition(*c)
		expr, err := builder.Build()
		if err != nil {
			return nil, err
		}
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
		input.FilterExpression = expr.Condition()
	}

	resMap := make([]map[string]interface{}, 0, table.ItemCount)
	paginator := dynamodb.NewScanPaginator(cli, input)
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
