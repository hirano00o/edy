package edy

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/hirano00o/edy/mocks"
)

func TestInstance_Delete(t *testing.T) {
	type args struct {
		ctx            context.Context
		tableName      string
		partitionValue string
		sortValue      string
		fileName       string
		f              func(string) (string, error)
	}
	tests := []struct {
		name    string
		args    args
		mocking func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI
		wantW   string
		wantErr bool
	}{
		{
			name: "Delete with partition value",
			args: args{
				ctx:            context.Background(),
				tableName:      "TEST",
				partitionValue: "TEST_VALUE1",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				table := describeTableOutputFixture(t, false)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(table, nil)
				input := &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						"TEST": {
							{
								DeleteRequest: &types.DeleteRequest{
									Key: map[string]types.AttributeValue{
										"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{
											Value: "TEST_VALUE1",
										},
									},
								},
							},
						},
					},
				}
				m.BatchWriteItemClient.On("BatchWriteItem", ctx, input).Return(&dynamodb.BatchWriteItemOutput{
					UnprocessedItems: map[string][]types.WriteRequest{},
				}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Delete with partition and sort value",
			args: args{
				ctx:            context.Background(),
				tableName:      "TEST",
				partitionValue: "TEST_VALUE1",
				sortValue:      "TEST_VALUE2",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				table := describeTableOutputFixture(t, false)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(table, nil)
				input := &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						"TEST": {
							{
								DeleteRequest: &types.DeleteRequest{
									Key: map[string]types.AttributeValue{
										"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{
											Value: "TEST_VALUE1",
										},
										"TEST_SORT_ATTRIBUTE": &types.AttributeValueMemberS{
											Value: "TEST_VALUE2",
										},
									},
								},
							},
						},
					},
				}
				m.BatchWriteItemClient.On("BatchWriteItem", ctx, input).Return(&dynamodb.BatchWriteItemOutput{
					UnprocessedItems: map[string][]types.WriteRequest{},
				}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Delete from file",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				fileName:  "TEST.json",
				f: func(s string) (string, error) {
					return "{\"partition\":\"TEST_VALUE1\",\"sort\":\"TEST_VALUE2\"}", nil
				},
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				table := describeTableOutputFixture(t, false)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(table, nil)
				input := &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						"TEST": {
							{
								DeleteRequest: &types.DeleteRequest{
									Key: map[string]types.AttributeValue{
										"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{
											Value: "TEST_VALUE1",
										},
										"TEST_SORT_ATTRIBUTE": &types.AttributeValueMemberS{
											Value: "TEST_VALUE2",
										},
									},
								},
							},
						},
					},
				}
				m.BatchWriteItemClient.On("BatchWriteItem", ctx, input).Return(&dynamodb.BatchWriteItemOutput{
					UnprocessedItems: map[string][]types.WriteRequest{},
				}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Multiple delete from file",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				fileName:  "TEST.json",
				f: func(s string) (string, error) {
					return "[{\"partition\":\"TEST_VALUE11\",\"sort\":\"TEST_VALUE12\"}," +
						"{\"partition\":\"TEST_VALUE21\",\"sort\":\"TEST_VALUE22\"}]", nil
				},
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				table := describeTableOutputFixture(t, false)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(table, nil)
				input := &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						"TEST": {
							{
								DeleteRequest: &types.DeleteRequest{
									Key: map[string]types.AttributeValue{
										"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{
											Value: "TEST_VALUE11",
										},
										"TEST_SORT_ATTRIBUTE": &types.AttributeValueMemberS{
											Value: "TEST_VALUE12",
										},
									},
								},
							},
							{
								DeleteRequest: &types.DeleteRequest{
									Key: map[string]types.AttributeValue{
										"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{
											Value: "TEST_VALUE21",
										},
										"TEST_SORT_ATTRIBUTE": &types.AttributeValueMemberS{
											Value: "TEST_VALUE22",
										},
									},
								},
							},
						},
					},
				}
				m.BatchWriteItemClient.On("BatchWriteItem", ctx, input).Return(&dynamodb.BatchWriteItemOutput{
					UnprocessedItems: map[string][]types.WriteRequest{},
				}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Delete with retry once",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				fileName:  "TEST.json",
				f: func(s string) (string, error) {
					return "[{\"partition\":\"TEST_VALUE11\",\"sort\":\"TEST_VALUE12\"}]", nil
				},
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				table := describeTableOutputFixture(t, false)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(table, nil)
				reqItems := map[string][]types.WriteRequest{
					"TEST": {
						{
							DeleteRequest: &types.DeleteRequest{
								Key: map[string]types.AttributeValue{
									"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{
										Value: "TEST_VALUE11",
									},
									"TEST_SORT_ATTRIBUTE": &types.AttributeValueMemberS{
										Value: "TEST_VALUE12",
									},
								},
							},
						},
					},
				}
				input := &dynamodb.BatchWriteItemInput{
					RequestItems: reqItems,
				}
				m.BatchWriteItemClient.On("BatchWriteItem", ctx, input).Return(&dynamodb.BatchWriteItemOutput{
					UnprocessedItems: reqItems,
				}, nil).Once()
				m.BatchWriteItemClient.On("BatchWriteItem", ctx, input).Return(&dynamodb.BatchWriteItemOutput{
					UnprocessedItems: map[string][]types.WriteRequest{},
				}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Delete failed, unprocessed items leaves",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				fileName:  "TEST.json",
				f: func(s string) (string, error) {
					return "[{\"partition\":\"TEST_VALUE11\",\"sort\":\"TEST_VALUE12\"}]", nil
				},
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				table := describeTableOutputFixture(t, false)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(table, nil)
				reqItems := map[string][]types.WriteRequest{
					"TEST": {
						{
							DeleteRequest: &types.DeleteRequest{
								Key: map[string]types.AttributeValue{
									"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{
										Value: "TEST_VALUE11",
									},
									"TEST_SORT_ATTRIBUTE": &types.AttributeValueMemberS{
										Value: "TEST_VALUE12",
									},
								},
							},
						},
					},
				}
				input := &dynamodb.BatchWriteItemInput{
					RequestItems: reqItems,
				}
				m.BatchWriteItemClient.On("BatchWriteItem", ctx, input).Return(&dynamodb.BatchWriteItemOutput{
					UnprocessedItems: reqItems,
				}, nil)

				return m
			},
			wantW: "{\n" +
				strings.Repeat(" ", 2) + "\"unprocessed\": [\n" +
				strings.Repeat(" ", 4) + "{\n" +
				strings.Repeat(" ", 6) + "\"DeleteRequest\": {\n" +
				strings.Repeat(" ", 8) + "\"Key\": {\n" +
				strings.Repeat(" ", 10) + "\"TEST_PARTITION_ATTRIBUTE\": {\n" +
				strings.Repeat(" ", 12) + "\"Value\": \"TEST_VALUE11\"\n" +
				strings.Repeat(" ", 10) + "},\n" +
				strings.Repeat(" ", 10) + "\"TEST_SORT_ATTRIBUTE\": {\n" +
				strings.Repeat(" ", 12) + "\"Value\": \"TEST_VALUE12\"\n" +
				strings.Repeat(" ", 10) + "}\n" +
				strings.Repeat(" ", 8) + "}\n" +
				strings.Repeat(" ", 6) + "},\n" +
				strings.Repeat(" ", 6) + "\"PutRequest\": null\n" +
				strings.Repeat(" ", 4) + "}\n" +
				strings.Repeat(" ", 2) + "]\n" +
				"}\n",
		},
		{
			name: "Error Delete",
			args: args{
				ctx:            context.Background(),
				tableName:      "TEST",
				partitionValue: "TEST_VALUE1",
				sortValue:      "TEST_VALUE2",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				table := describeTableOutputFixture(t, false)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(table, nil)
				input := &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						"TEST": {
							{
								DeleteRequest: &types.DeleteRequest{
									Key: map[string]types.AttributeValue{
										"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{
											Value: "TEST_VALUE1",
										},
										"TEST_SORT_ATTRIBUTE": &types.AttributeValueMemberS{
											Value: "TEST_VALUE2",
										},
									},
								},
							},
						},
					},
				}
				m.BatchWriteItemClient.On("BatchWriteItem", ctx, input).Return(nil, fmt.Errorf("cannot delete"))

				return m
			},
			wantErr: true,
		},
		{
			name: "Error DescribeTable",
			args: args{
				ctx:            context.Background(),
				tableName:      "TEST",
				partitionValue: "TEST_VALUE1",
				sortValue:      "TEST_VALUE2",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(nil, fmt.Errorf("cannot describe table"))

				return m
			},
			wantErr: true,
		},
		{
			name: "Partition value is empty",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				return m
			},
			wantErr: true,
		},
		{
			name: "Both partition and file argument is not empty",
			args: args{
				ctx:            context.Background(),
				tableName:      "TEST",
				partitionValue: "TEST_VALUE1",
				fileName:       "TEST.json",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				return m
			},
			wantErr: true,
		},
		{
			name: "Error read file",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				fileName:  "TEST.json",
				f: func(string) (string, error) {
					return "", fmt.Errorf("cannot read file")
				},
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				return m
			},
			wantErr: true,
		},
		{
			name: "Invalid json in file",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				fileName:  "TEST.json",
				f: func(string) (string, error) {
					return "[{\"partition\"\"TEST_VALUE1\"}]", nil
				},
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				return m
			},
			wantErr: true,
		},
		{
			name: "Invalid key in json with list",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				fileName:  "TEST.json",
				f: func(string) (string, error) {
					return "[{\"invalid\":\"invalid\"}]", nil
				},
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				return m
			},
			wantErr: true,
		},
		{
			name: "Invalid key in json",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				fileName:  "TEST.json",
				f: func(string) (string, error) {
					return "{\"invalid\":\"invalid\"}", nil
				},
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
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
			err := i.Delete(
				tt.args.ctx,
				w,
				tt.args.tableName,
				tt.args.partitionValue,
				tt.args.sortValue,
				tt.args.fileName,
				tt.args.f,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Delete() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
