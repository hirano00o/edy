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
			wantW: "{\n  \"succeeded\": 1\n}\n",
		},
		{
			name: "Put 2 items",
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
			wantW: "{\n  \"succeeded\": 1\n}\n",
		},
		{
			name: "Put items with map",
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
			wantW: "{\n  \"succeeded\": 1\n}\n",
		},
		{
			name: "Put items with list",
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
			wantW: "{\n  \"succeeded\": 1\n}\n",
		},
		{
			name: "Put items with null",
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
			wantW: "{\n  \"succeeded\": 1\n}\n",
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
			name: "Error put items",
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
