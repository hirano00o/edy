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

func TestInstance_Put(t *testing.T) {
	type args struct {
		ctx       context.Context
		tableName string
		item      string
	}
	tests := []struct {
		name    string
		args    args
		mocking func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI
		wantW   string
		wantErr bool
	}{
		{
			name: "Put key value",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item:      "{\"TEST_KEY1\":\"TEST_VALUE1\"}",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				input := &dynamodb.PutItemInput{
					TableName: aws.String("TEST"),
					Item: map[string]types.AttributeValue{
						"TEST_KEY1": &types.AttributeValueMemberS{
							Value: "TEST_VALUE1",
						},
					},
				}
				m.PutItemClient.On("PutItem", ctx, input).Return(&dynamodb.PutItemOutput{}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Put 2 data",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item:      "{\"TEST_KEY1\":\"TEST_VALUE1\",\"TEST_KEY2\":[1,2,3]}",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				input := &dynamodb.PutItemInput{
					TableName: aws.String("TEST"),
					Item: map[string]types.AttributeValue{
						"TEST_KEY1": &types.AttributeValueMemberS{
							Value: "TEST_VALUE1",
						},
						"TEST_KEY2": &types.AttributeValueMemberNS{
							Value: []string{"1", "2", "3"},
						},
					},
				}
				m.PutItemClient.On("PutItem", ctx, input).Return(&dynamodb.PutItemOutput{}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Put data with map",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item: "{\"TEST_KEY1\":\"T1\",\"TEST_KEY2\":{" +
					"\"TEST_KEY3\":{\"TEST_KEY4\":[1,2],\"TEST_KEY5\":[\"T5\"]},\"TEST_KEY6\":{\"TEST_KEY7\":1}" +
					"},\"TEST_KEY8\":false}",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				input := &dynamodb.PutItemInput{
					TableName: aws.String("TEST"),
					Item: map[string]types.AttributeValue{
						"TEST_KEY1": &types.AttributeValueMemberS{
							Value: "T1",
						},
						"TEST_KEY2": &types.AttributeValueMemberM{
							Value: map[string]types.AttributeValue{
								"TEST_KEY3": &types.AttributeValueMemberM{
									Value: map[string]types.AttributeValue{
										"TEST_KEY4": &types.AttributeValueMemberNS{
											Value: []string{"1", "2"},
										},
										"TEST_KEY5": &types.AttributeValueMemberSS{
											Value: []string{"T5"},
										},
									},
								},
								"TEST_KEY6": &types.AttributeValueMemberM{
									Value: map[string]types.AttributeValue{
										"TEST_KEY7": &types.AttributeValueMemberN{
											Value: "1",
										},
									},
								},
							},
						},
						"TEST_KEY8": &types.AttributeValueMemberBOOL{
							Value: false,
						},
					},
				}
				m.PutItemClient.On("PutItem", ctx, input).Return(&dynamodb.PutItemOutput{}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Put data with list",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item:      "{\"TEST_KEY1\":[\"T1\", true, 1]}",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				input := &dynamodb.PutItemInput{
					TableName: aws.String("TEST"),
					Item: map[string]types.AttributeValue{
						"TEST_KEY1": &types.AttributeValueMemberL{
							Value: []types.AttributeValue{
								&types.AttributeValueMemberS{
									Value: "T1",
								},
								&types.AttributeValueMemberBOOL{
									Value: true,
								},
								&types.AttributeValueMemberN{
									Value: "1",
								},
							},
						},
					},
				}
				m.PutItemClient.On("PutItem", ctx, input).Return(&dynamodb.PutItemOutput{}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Put data with list include null",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item:      "{\"TEST_KEY1\":[null, \"T1\"]}",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				input := &dynamodb.PutItemInput{
					TableName: aws.String("TEST"),
					Item: map[string]types.AttributeValue{
						"TEST_KEY1": &types.AttributeValueMemberL{
							Value: []types.AttributeValue{
								&types.AttributeValueMemberNULL{
									Value: true,
								},
								&types.AttributeValueMemberS{
									Value: "T1",
								},
							},
						},
					},
				}
				m.PutItemClient.On("PutItem", ctx, input).Return(&dynamodb.PutItemOutput{}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Put data with null",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item:      "{\"TEST_KEY1\":null}",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				input := &dynamodb.PutItemInput{
					TableName: aws.String("TEST"),
					Item: map[string]types.AttributeValue{
						"TEST_KEY1": &types.AttributeValueMemberNULL{
							Value: true,
						},
					},
				}
				m.PutItemClient.On("PutItem", ctx, input).Return(&dynamodb.PutItemOutput{}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Put 1 item with array",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item:      "[{\"TEST_KEY1\":\"TEST_VALUE1\"}]",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				input := &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						"TEST": {
							{
								PutRequest: &types.PutRequest{
									Item: map[string]types.AttributeValue{
										"TEST_KEY1": &types.AttributeValueMemberS{
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
			name: "Put 2 items",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item: "[{\"TEST_KEY1\":\"TEST_VALUE1\"}," +
					"{\"TEST_KEY2\":\"TEST_VALUE2\", \"TEST_KEY3\":[\"TEST_VALUE31\",32,true]}]",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				input := &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						"TEST": {
							{
								PutRequest: &types.PutRequest{
									Item: map[string]types.AttributeValue{
										"TEST_KEY1": &types.AttributeValueMemberS{
											Value: "TEST_VALUE1",
										},
									},
								},
							},
							{
								PutRequest: &types.PutRequest{
									Item: map[string]types.AttributeValue{
										"TEST_KEY2": &types.AttributeValueMemberS{
											Value: "TEST_VALUE2",
										},
										"TEST_KEY3": &types.AttributeValueMemberL{
											Value: []types.AttributeValue{
												&types.AttributeValueMemberS{
													Value: "TEST_VALUE31",
												},
												&types.AttributeValueMemberN{
													Value: "32",
												},
												&types.AttributeValueMemberBOOL{
													Value: true,
												},
											},
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
			name: "Put item with 1 retry",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item:      "[{\"TEST_KEY1\":\"TEST_VALUE1\"}]",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				reqItems := map[string][]types.WriteRequest{
					"TEST": {
						{
							PutRequest: &types.PutRequest{
								Item: map[string]types.AttributeValue{
									"TEST_KEY1": &types.AttributeValueMemberS{
										Value: "TEST_VALUE1",
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
			name: "Put item, failed retry",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item:      "[{\"TEST_KEY1\":\"TEST_VALUE1\"}]",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				reqItems := map[string][]types.WriteRequest{
					"TEST": {
						{
							PutRequest: &types.PutRequest{
								Item: map[string]types.AttributeValue{
									"TEST_KEY1": &types.AttributeValueMemberS{
										Value: "TEST_VALUE1",
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
				strings.Repeat(" ", 6) + "\"DeleteRequest\": null,\n" +
				strings.Repeat(" ", 6) + "\"PutRequest\": {\n" +
				strings.Repeat(" ", 8) + "\"Item\": {\n" +
				strings.Repeat(" ", 10) + "\"TEST_KEY1\": {\n" +
				strings.Repeat(" ", 12) + "\"Value\": \"TEST_VALUE1\"\n" +
				strings.Repeat(" ", 10) + "}\n" +
				strings.Repeat(" ", 8) + "}\n" +
				strings.Repeat(" ", 6) + "}\n" +
				strings.Repeat(" ", 4) + "}\n" +
				strings.Repeat(" ", 2) + "]\n" +
				"}\n",
		},
		{
			name: "Error unmarshal JSON",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item:      "ERROR",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				m.On("CreateInstance").Return(m)

				return m
			},
			wantErr: true,
		},
		{
			name: "Error put item",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item:      "{\"ERROR\":\"ERROR\"}",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				input := &dynamodb.PutItemInput{
					TableName: aws.String("TEST"),
					Item: map[string]types.AttributeValue{
						"ERROR": &types.AttributeValueMemberS{
							Value: "ERROR",
						},
					},
				}
				m.PutItemClient.On("PutItem", ctx, input).Return(nil, fmt.Errorf("cannot put items"))

				return m
			},
			wantErr: true,
		},
		{
			name: "Error batch write item",
			args: args{
				ctx:       context.Background(),
				tableName: "TEST",
				item:      "[{\"ERROR\":\"ERROR\"}]",
			},
			mocking: func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI {
				t.Helper()

				m := new(mocks.MockDynamoDBAPI)
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				input := &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						"TEST": {
							{
								PutRequest: &types.PutRequest{
									Item: map[string]types.AttributeValue{
										"ERROR": &types.AttributeValueMemberS{
											Value: "ERROR",
										},
									},
								},
							},
						},
					},
				}
				m.BatchWriteItemClient.On("BatchWriteItem", ctx, input).Return(nil, fmt.Errorf("cannot batch write items"))

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
			err := i.Put(tt.args.ctx, w, tt.args.tableName, tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Put() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
