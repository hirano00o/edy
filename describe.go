package edy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/hirano00o/edy/client"
	"github.com/hirano00o/edy/model"
)

func describeTable(ctx context.Context, tableName string) (*model.Table, error) {
	cli := ctx.Value(newClientKey).(client.DynamoDB)
	res, err := cli.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}
	t := model.Table{
		Arn:       aws.ToString(res.Table.TableArn),
		Name:      aws.ToString(res.Table.TableName),
		ItemCount: res.Table.ItemCount,
	}
	attr := make(map[string]model.AttributeType)
	for _, a := range res.Table.AttributeDefinitions {
		attr[aws.ToString(a.AttributeName)] = model.AttributeTypeStr(a.AttributeType).Name()
	}
	for _, k := range res.Table.KeySchema {
		name := aws.ToString(k.AttributeName)
		switch k.KeyType {
		case "HASH":
			t.PartitionKey = new(model.Key)
			t.PartitionKey.Name = name
			t.PartitionKey.Type = attr[t.PartitionKey.Name]
			t.PartitionKey.TypeStr = attr[t.PartitionKey.Name].String()
		case "RANGE":
			t.SortKey = new(model.Key)
			t.SortKey.Name = name
			t.SortKey.Type = attr[t.SortKey.Name]
			t.SortKey.TypeStr = attr[t.SortKey.Name].String()
		}
	}
	t.GSI = make([]*model.GlobalSecondaryIndex, len(res.Table.GlobalSecondaryIndexes))
	for i, g := range res.Table.GlobalSecondaryIndexes {
		t.GSI[i] = new(model.GlobalSecondaryIndex)
		t.GSI[i].Name = aws.ToString(g.IndexName)
		for j := range g.KeySchema {
			name := aws.ToString(g.KeySchema[j].AttributeName)
			switch g.KeySchema[j].KeyType {
			case "HASH":
				t.GSI[i].PartitionKey = new(model.Key)
				t.GSI[i].PartitionKey.Name = name
				t.GSI[i].PartitionKey.Type = attr[t.GSI[i].PartitionKey.Name]
				t.GSI[i].PartitionKey.TypeStr = attr[t.GSI[i].PartitionKey.Name].String()
			case "RANGE":
				t.GSI[i].SortKey = new(model.Key)
				t.GSI[i].SortKey.Name = name
				t.GSI[i].SortKey.Type = attr[t.GSI[i].SortKey.Name]
				t.GSI[i].SortKey.TypeStr = attr[t.GSI[i].SortKey.Name].String()
			}
		}
	}

	return &t, nil
}

func (i *Instance) DescribeTable(ctx context.Context, w io.Writer, tableName string) error {
	cli := i.NewClient.CreateInstance()
	ctx = context.WithValue(ctx, newClientKey, cli)

	t, err := describeTable(ctx, tableName)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(t, "", strings.Repeat(" ", 2))
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", string(b))

	return nil
}
