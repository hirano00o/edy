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

func TestInstance_Delete(t *testing.T) {
	type args struct {
		ctx            context.Context
		tableName      string
		partitionValue string
		sortValue      string
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
				input := &dynamodb.DeleteItemInput{
					TableName: aws.String("TEST"),
					Key: map[string]types.AttributeValue{
						"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{
							Value: "TEST_VALUE1",
						},
					},
				}
				m.DeleteItemClient.On("DeleteItem", ctx, input).Return(&dynamodb.DeleteItemOutput{}, nil)

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
				input := &dynamodb.DeleteItemInput{
					TableName: aws.String("TEST"),
					Key: map[string]types.AttributeValue{
						"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{
							Value: "TEST_VALUE1",
						},
						"TEST_SORT_ATTRIBUTE": &types.AttributeValueMemberS{
							Value: "TEST_VALUE2",
						},
					},
				}
				m.DeleteItemClient.On("DeleteItem", ctx, input).Return(&dynamodb.DeleteItemOutput{}, nil)

				return m
			},
			wantW: "{\n  \"unprocessed\": []\n}\n",
		},
		{
			name: "Error DeleteItem",
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
				input := &dynamodb.DeleteItemInput{
					TableName: aws.String("TEST"),
					Key: map[string]types.AttributeValue{
						"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{
							Value: "TEST_VALUE1",
						},
						"TEST_SORT_ATTRIBUTE": &types.AttributeValueMemberS{
							Value: "TEST_VALUE2",
						},
					},
				}
				m.DeleteItemClient.On("DeleteItem", ctx, input).Return(nil, fmt.Errorf("cannot delete item"))

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
				ctx = context.WithValue(ctx, newClientKey, m)
				m.On("CreateInstance").Return(m)
				table := describeTableOutputFixture(t, false)
				m.DescribeTableAPIClient.On("DescribeTable", ctx, &dynamodb.DescribeTableInput{
					TableName: aws.String("TEST"),
				}).Return(table, nil)

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
			err := i.Delete(tt.args.ctx, w, tt.args.tableName, tt.args.partitionValue, tt.args.sortValue)
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
