package edy

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/hirano00o/edy/mocks"
)

func TestInstance_Scan(t *testing.T) {
	type args struct {
		ctx             context.Context
		tableName       string
		filterCondition string
		projection      string
		output          string
	}
	tests := []struct {
		name    string
		args    args
		mocking func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI
		wantW   string
		wantErr bool
	}{
		{
			name: "Scan",
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
				input := &dynamodb.ScanInput{
					TableName: aws.String("TEST"),
				}
				m.ScanAPIClient.On("Scan", ctx, input).Return(scanOutputFixture(t, ""), nil)
				return m
			},
			wantW: jsonFixture(t, []map[string]interface{}{
				{
					"TEST_PARTITION_ATTRIBUTE": "TEST_PARTITION_VALUE_1",
					"TEST_SORT_ATTRIBUTE":      "TEST_SORT_VALUE_1",
					"TEST_ATTRIBUTE_1":         "TEST_ATTRIBUTE_1_VALUE_1",
					"TEST_ATTRIBUTE_2":         1,
				},
			}),
		},
		{
			name: "Scan with filter",
			args: args{
				ctx:             context.Background(),
				tableName:       "TEST",
				filterCondition: "TEST_ATTRIBUTE_2,N = 1",
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
				filterCondition := expression.Equal(
					expression.Name("TEST_ATTRIBUTE_2"),
					expression.Value(1),
				)
				expr, err := expression.NewBuilder().WithCondition(filterCondition).Build()
				if err != nil {
					t.Fatalf("expression build error: %v", err)
				}
				input := &dynamodb.ScanInput{
					TableName:                 aws.String("TEST"),
					ExpressionAttributeNames:  expr.Names(),
					ExpressionAttributeValues: expr.Values(),
					FilterExpression:          expr.Condition(),
				}
				m.ScanAPIClient.On("Scan", ctx, input).Return(scanOutputFixture(t, ""), nil)
				return m
			},
			wantW: jsonFixture(t, []map[string]interface{}{
				{
					"TEST_PARTITION_ATTRIBUTE": "TEST_PARTITION_VALUE_1",
					"TEST_SORT_ATTRIBUTE":      "TEST_SORT_VALUE_1",
					"TEST_ATTRIBUTE_1":         "TEST_ATTRIBUTE_1_VALUE_1",
					"TEST_ATTRIBUTE_2":         1,
				},
			}),
		},
		{
			name: "Scan with projection",
			args: args{
				ctx:        context.Background(),
				tableName:  "TEST",
				projection: "TEST_ATTRIBUTE_1 TEST_ATTRIBUTE_2",
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
				pj := expression.NamesList(
					expression.Name("TEST_ATTRIBUTE_1"),
					expression.Name("TEST_ATTRIBUTE_2"),
				)
				expr, err := expression.NewBuilder().WithProjection(pj).Build()
				if err != nil {
					t.Fatalf("expression build error: %v", err)
				}
				input := &dynamodb.ScanInput{
					TableName:                 aws.String("TEST"),
					ExpressionAttributeNames:  expr.Names(),
					ExpressionAttributeValues: expr.Values(),
					ProjectionExpression:      expr.Projection(),
				}
				m.ScanAPIClient.On("Scan", ctx, input).
					Return(scanOutputFixture(t, "TEST_ATTRIBUTE_1 TEST_ATTRIBUTE_2"), nil)
				return m
			},
			wantW: jsonFixture(t, []map[string]interface{}{
				{
					"TEST_ATTRIBUTE_1": "TEST_ATTRIBUTE_1_VALUE_1",
					"TEST_ATTRIBUTE_2": 1,
				},
			}),
		},
		{
			name: "Scan error",
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
				input := &dynamodb.ScanInput{
					TableName: aws.String("TEST"),
				}
				m.ScanAPIClient.On("Scan", ctx, input).Return(nil, fmt.Errorf("scan error"))
				return m
			},
			wantErr: true,
		},
		{
			name: "Invalid filter condition",
			args: args{
				ctx:             context.Background(),
				tableName:       "TEST",
				filterCondition: "ERROR = ERROR",
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
		{
			name: "Error DescribeTable",
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
				}).Return(nil, fmt.Errorf("cannot describe table"))

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
			err := i.Scan(
				tt.args.ctx,
				w,
				tt.args.tableName,
				tt.args.filterCondition,
				tt.args.projection,
				tt.args.output,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Scan() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func scanOutputFixture(t *testing.T, pj string) *dynamodb.ScanOutput {
	t.Helper()
	if len(pj) != 0 {
		m := map[string]types.AttributeValue{
			"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{Value: "TEST_PARTITION_VALUE_1"},
			"TEST_SORT_ATTRIBUTE":      &types.AttributeValueMemberS{Value: "TEST_SORT_VALUE_1"},
			"TEST_ATTRIBUTE_1":         &types.AttributeValueMemberS{Value: "TEST_ATTRIBUTE_1_VALUE_1"},
			"TEST_ATTRIBUTE_2":         &types.AttributeValueMemberN{Value: "1"},
		}
		var rMap = map[string]types.AttributeValue{}
		for _, s := range strings.Split(pj, " ") {
			rMap[s] = m[s]
		}
		return &dynamodb.ScanOutput{
			Items: []map[string]types.AttributeValue{
				rMap,
			},
			Count: 1,
		}
	}
	return &dynamodb.ScanOutput{
		Items: []map[string]types.AttributeValue{
			{
				"TEST_PARTITION_ATTRIBUTE": &types.AttributeValueMemberS{Value: "TEST_PARTITION_VALUE_1"},
				"TEST_SORT_ATTRIBUTE":      &types.AttributeValueMemberS{Value: "TEST_SORT_VALUE_1"},
				"TEST_ATTRIBUTE_1":         &types.AttributeValueMemberS{Value: "TEST_ATTRIBUTE_1_VALUE_1"},
				"TEST_ATTRIBUTE_2":         &types.AttributeValueMemberN{Value: "1"},
			},
		},
		Count: 1,
	}
}
