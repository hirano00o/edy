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
	"github.com/hirano00o/edy/model"
)

func analyseSortCondition(
	sortCondition,
	sortKey string,
	sortKeyType model.AttributeType,
) (*expression.KeyConditionBuilder, error) {
	s := strings.Split(sortCondition, " ")
	if len(s) > 3 {
		return nil, fmt.Errorf("invalid condition, specified condition is a lot: %s", sortCondition)
	} else if len(s) < 2 {
		return nil, fmt.Errorf("invalid condition, specified condition is insufficient: %s", sortCondition)
	}

	op, err := model.ConvertToComparisonOperator(s[0])
	if err != nil {
		return nil, err
	}
	if op == model.BETWEEN {
		if len(s) != 3 {
			return nil, fmt.Errorf("invalid condition, specified condition is insufficient: %s", sortCondition)
		}
	}

	var c expression.KeyConditionBuilder
	v1, err := sortKeyType.Value(s[1])
	if err != nil {
		return nil, fmt.Errorf("invalid condition, %v", err)
	}

	switch op {
	case model.EQ:
		c = expression.KeyEqual(expression.Key(sortKey), expression.Value(v1))
	case model.LE:
		c = expression.KeyLessThanEqual(expression.Key(sortKey), expression.Value(v1))
	case model.LT:
		c = expression.KeyLessThan(expression.Key(sortKey), expression.Value(v1))
	case model.GE:
		c = expression.KeyGreaterThanEqual(expression.Key(sortKey), expression.Value(v1))
	case model.GT:
		c = expression.KeyGreaterThan(expression.Key(sortKey), expression.Value(v1))
	case model.BeginsWith:
		c = expression.KeyBeginsWith(expression.Key(sortKey), v1.(string))
	case model.BETWEEN:
		v2, err := sortKeyType.Value(s[2])
		if err != nil {
			return nil, fmt.Errorf("invalid condition, %v", err)
		}
		c = expression.KeyBetween(expression.Key(sortKey), expression.Value(v1), expression.Value(v2))
	default:
		return nil, fmt.Errorf("invalid condition, specified comparison operator can not use: %s", sortCondition)
	}

	return &c, nil
}

func query(
	ctx context.Context,
	tableName,
	partitionValue,
	sortCondition,
	filterCondition,
	index string,
) ([]map[string]interface{}, error) {
	table, err := describeTable(ctx, tableName)
	if err != nil {
		return nil, err
	}
	cli := ctx.Value(newClientKey).(client.DynamoDB)

	partitionKeyName, partitionKeyType := table.PartitionKey.Name, table.PartitionKey.Type
	var sortKeyName string
	var sortKeyType model.AttributeType
	if table.SortKey != nil {
		sortKeyName, sortKeyType = table.SortKey.Name, table.SortKey.Type
	}

	// Index
	if len(index) != 0 {
		indexIsGSI := false
		for i := range table.GSI {
			if table.GSI[i].Name == index {
				partitionKeyName, partitionKeyType = table.GSI[i].PartitionKey.Name, table.GSI[i].PartitionKey.Type
				if table.GSI[i].SortKey != nil {
					sortKeyName, sortKeyType = table.GSI[i].SortKey.Name, table.GSI[i].SortKey.Type
				}
				indexIsGSI = true
				break
			}
		}
		if !indexIsGSI {
			return nil, fmt.Errorf("there is no index: %s", index)
		}
	}
	v, err := partitionKeyType.Value(partitionValue)
	if err != nil {
		return nil, err
	}

	// PartitionKey condition
	condition := expression.KeyEqual(expression.Key(partitionKeyName), expression.Value(v))
	// SortKey condition
	if len(sortCondition) != 0 {
		c, err := analyseSortCondition(sortCondition, sortKeyName, sortKeyType)
		if err != nil {
			return nil, err
		}
		condition = condition.And(*c)
	}
	builder := expression.NewBuilder().WithKeyCondition(condition)
	// Filter condition
	if len(filterCondition) != 0 {
		c, err := analyseFilterCondition(filterCondition)
		if err != nil {
			return nil, err
		}
		builder = builder.WithCondition(*c)
	}

	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Condition(),
	}
	if len(index) != 0 {
		input.IndexName = aws.String(index)
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
	tableName,
	partitionValue,
	sortCondition,
	filterCondition,
	index string,
) error {
	cli := i.NewClient.CreateInstance()
	ctx = context.WithValue(ctx, newClientKey, cli)

	res, err := query(ctx, tableName, partitionValue, sortCondition, filterCondition, index)
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
