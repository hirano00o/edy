package edy

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/hirano00o/edy/client"
	"github.com/hirano00o/edy/model"
	"io"
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
	attr := make(map[string]string)
	for _, a := range res.Table.AttributeDefinitions {
		attr[aws.ToString(a.AttributeName)] = string(a.AttributeType)
	}
	for _, k := range res.Table.KeySchema {
		name := aws.ToString(k.AttributeName)
		switch k.KeyType {
		case "HASH":
			t.PartitionKeyName = name
			t.PartitionKeyType = attr[t.PartitionKeyName]
		case "RANGE":
			t.SortKeyName = name
			t.SortKeyType = attr[t.SortKeyName]
		}
	}
	t.GSI = make([]*model.GlobalSecondaryIndex, len(res.Table.GlobalSecondaryIndexes))
	for i, g := range res.Table.GlobalSecondaryIndexes {
		t.GSI[i].Name = aws.ToString(g.IndexName)
		for j := range g.KeySchema {
			name := aws.ToString(g.KeySchema[j].AttributeName)
			switch g.KeySchema[j].KeyType {
			case "HASH":
				t.GSI[i].PartitionKeyName = name
				t.GSI[i].PartitionKeyType = attr[t.GSI[i].PartitionKeyName]
			case "RANGE":
				t.GSI[i].SortKeyName = name
				t.GSI[i].SortKeyType = attr[t.GSI[i].SortKeyName]
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
	fmt.Fprintf(w, "Table Arn:\t%s\n", t.Arn)
	fmt.Fprintf(w, "Table Name:\t%s\n", t.Name)
	fmt.Fprintf(w, "Partition Key:\t%s(%s)\n", t.PartitionKeyName, t.PartitionKeyType)
	if len(t.SortKeyName) != 0 {
		fmt.Fprintf(w, "Sort Key:\t%s(%s)\n", t.SortKeyName, t.SortKeyType)
	}
	if len(t.GSI) != 0 {
		fmt.Fprintf(w, "GSI:\n")
	}
	for i := range t.GSI {
		fmt.Fprintf(w, "\tIndex:\t%s\n", t.GSI[i].Name)
		fmt.Fprintf(w, "\tPartition Key:\t%s(%s)\n", t.GSI[i].PartitionKeyName, t.GSI[i].PartitionKeyType)
		if len(t.SortKeyName) != 0 {
			fmt.Fprintf(w, "\tSort Key:\t%s(%s)\n", t.GSI[i].SortKeyName, t.GSI[i].SortKeyType)
		}
	}
	fmt.Fprintf(w, "ItemCount:\t%d", t.ItemCount)

	return nil
}
