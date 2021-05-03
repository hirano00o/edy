package edy

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/hirano00o/edy/mocks"
)

func TestInstance_DescribeTable(t *testing.T) {
	type args struct {
		ctx       context.Context
		tableName string
	}
	tests := []struct {
		name    string
		args    args
		mocking func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI
		wantW   string
		wantErr bool
	}{
		{
			name: "Describe TEST table",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()
				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(describeTableOutputFixture(t, false), nil)
				return m
			},
			wantW: "Table Arn:\tTEST_ARN\nTable Name:\tTEST\nPartition Key:\tTEST_PARTITION_ATTRIBUTE(S)\n" +
				"Sort Key:\tTEST_SORT_ATTRIBUTE(S)\nItemCount:\t1",
		},
		{
			name: "Describe TEST table with GSI",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()
				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(describeTableOutputFixture(t, true), nil)
				return m
			},
			wantW: "Table Arn:\tTEST_ARN\nTable Name:\tTEST\n" +
				"Partition Key:\tTEST_PARTITION_ATTRIBUTE(S)\nSort Key:\tTEST_SORT_ATTRIBUTE(S)\n" +
				"GSI:\n\tIndex:\tTEST_GSI\n" +
				"\tPartition Key:\tTEST_ATTRIBUTE_1(S)\n\tSort Key:\tTEST_ATTRIBUTE_2(N)\n" +
				"ItemCount:\t1",
		},
		{
			name: "DescribeTable error",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()
				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(nil, fmt.Errorf("DescribeTable error"))
				return m
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mocking(t, tt.args.ctx)
			i := &Instance{
				NewClient: mock,
			}
			w := &bytes.Buffer{}
			err := i.DescribeTable(tt.args.ctx, w, tt.args.tableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("DescribeTable() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func describeTableOutputFixture(t *testing.T, gsi bool) *dynamodb.DescribeTableOutput {
	t.Helper()

	output := &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			TableArn:  aws.String("TEST_ARN"),
			TableId:   aws.String("TEST_ID"),
			TableName: aws.String("TEST"),
			ItemCount: 1,
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("TEST_PARTITION_ATTRIBUTE"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("TEST_SORT_ATTRIBUTE"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("TEST_ATTRIBUTE_1"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("TEST_ATTRIBUTE_2"),
					AttributeType: types.ScalarAttributeTypeN,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("TEST_PARTITION_ATTRIBUTE"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("TEST_SORT_ATTRIBUTE"),
					KeyType:       types.KeyTypeRange,
				},
			},
		},
	}

	if gsi {
		output.Table.GlobalSecondaryIndexes = append(
			output.Table.GlobalSecondaryIndexes,
			[]types.GlobalSecondaryIndexDescription{
				{
					IndexName: aws.String("TEST_GSI"),
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("TEST_ATTRIBUTE_1"),
							KeyType:       types.KeyTypeHash,
						},
						{
							AttributeName: aws.String("TEST_ATTRIBUTE_2"),
							KeyType:       types.KeyTypeRange,
						},
					},
				},
			}...,
		)
	}

	return output
}
