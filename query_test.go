package edy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/hirano00o/edy/mocks"
	"github.com/hirano00o/edy/model"
)

func Test_analyseSortCondition(t *testing.T) {
	type args struct {
		sortCondition string
		sortKey       string
		sortKeyType   model.AttributeType
	}
	tests := []struct {
		name    string
		args    args
		want    expression.KeyConditionBuilder
		wantErr bool
	}{
		{
			name: "EQ case",
			args: args{
				sortCondition: "= 1234",
				sortKey:       "ID",
				sortKeyType:   model.S{},
			},
			want: expression.KeyEqual(expression.Key("ID"), expression.Value("1234")),
		},
		{
			name: "LE case",
			args: args{
				sortCondition: "<= 1234",
				sortKey:       "ID",
				sortKeyType:   model.S{},
			},
			want: expression.KeyLessThanEqual(expression.Key("ID"), expression.Value("1234")),
		},
		{
			name: "LT case",
			args: args{
				sortCondition: "< 1234",
				sortKey:       "ID",
				sortKeyType:   model.S{},
			},
			want: expression.KeyLessThan(expression.Key("ID"), expression.Value("1234")),
		},
		{
			name: "GE case",
			args: args{
				sortCondition: ">= 1234",
				sortKey:       "ID",
				sortKeyType:   model.N{},
			},
			want: expression.KeyGreaterThanEqual(expression.Key("ID"), expression.Value(1234)),
		},
		{
			name: "GT case",
			args: args{
				sortCondition: "> 1234",
				sortKey:       "ID",
				sortKeyType:   model.N{},
			},
			want: expression.KeyGreaterThan(expression.Key("ID"), expression.Value(1234)),
		},
		{
			name: "BeginsWith case",
			args: args{
				sortCondition: "begins_with 1234",
				sortKey:       "ID",
				sortKeyType:   model.S{},
			},
			want: expression.KeyBeginsWith(expression.Key("ID"), "1234"),
		},
		{
			name: "Between case",
			args: args{
				sortCondition: "between 12 34",
				sortKey:       "ID",
				sortKeyType:   model.N{},
			},
			want: expression.KeyBetween(expression.Key("ID"), expression.Value(12), expression.Value(34)),
		},
		{
			name: "Missing condition value",
			args: args{
				sortCondition: "EQ",
				sortKey:       "ID",
				sortKeyType:   model.N{},
			},
			wantErr: true,
		},
		{
			name: "Lot condition value",
			args: args{
				sortCondition: "EQ 12 34",
				sortKey:       "ID",
				sortKeyType:   model.N{},
			},
			wantErr: true,
		},
		{
			name: "Between lot condition value",
			args: args{
				sortCondition: "between 12 34 56",
				sortKey:       "ID",
				sortKeyType:   model.N{},
			},
			wantErr: true,
		},
		{
			name: "Missing between value",
			args: args{
				sortCondition: "between 12",
				sortKey:       "ID",
				sortKeyType:   model.N{},
			},
			wantErr: true,
		},
		{
			name: "Invalid number value",
			args: args{
				sortCondition: "= a12",
				sortKey:       "ID",
				sortKeyType:   model.N{},
			},
			wantErr: true,
		},
		{
			name: "Invalid number value when between",
			args: args{
				sortCondition: "between 12 a34",
				sortKey:       "ID",
				sortKeyType:   model.N{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := analyseSortCondition(tt.args.sortCondition, tt.args.sortKey, tt.args.sortKeyType)
			if (err != nil) != tt.wantErr {
				t.Errorf("analyseSortCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got == nil {
				return
			}
			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("analyseSortCondition() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstance_Query(t *testing.T) {
	type args struct {
		ctx             context.Context
		tableName       string
		partitionValue  string
		sortCondition   string
		filterCondition string
	}
	tests := []struct {
		name    string
		args    args
		mocking func(t *testing.T, ctx context.Context) *mocks.MockDynamoDBAPI
		wantW   string
		wantErr bool
	}{
		{
			name: "Query with partition key",
			args: args{
				ctx:            context.Background(),
				tableName:      "TEST",
				partitionValue: "TEST_PARTITION_VALUE_1",
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
				condition := expression.KeyEqual(
					expression.Key("TEST_PARTITION_ATTRIBUTE"),
					expression.Value("TEST_PARTITION_VALUE_1"),
				)
				builder := expression.NewBuilder().WithKeyCondition(condition)
				expr, err := builder.Build()
				if err != nil {
					t.Fatalf("expression build error: %v", err)
				}
				input := &dynamodb.QueryInput{
					TableName:                 aws.String("TEST"),
					ExpressionAttributeNames:  expr.Names(),
					ExpressionAttributeValues: expr.Values(),
					KeyConditionExpression:    expr.KeyCondition(),
				}
				m.QueryAPIClient.On("Query", ctx, input).Return(queryOutputFixture(t), nil)
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
			name: "Query with partition and sort key",
			args: args{
				ctx:            context.Background(),
				tableName:      "TEST",
				partitionValue: "TEST_PARTITION_VALUE_1",
				sortCondition:  "= TEST_SORT_VALUE_1",
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
				condition := expression.KeyEqual(
					expression.Key("TEST_PARTITION_ATTRIBUTE"),
					expression.Value("TEST_PARTITION_VALUE_1"),
				).And(expression.KeyEqual(
					expression.Key("TEST_SORT_ATTRIBUTE"),
					expression.Value("TEST_SORT_VALUE_1"),
				))
				builder := expression.NewBuilder().WithKeyCondition(condition)
				expr, err := builder.Build()
				if err != nil {
					t.Fatalf("expression build error: %v", err)
				}
				input := &dynamodb.QueryInput{
					TableName:                 aws.String("TEST"),
					ExpressionAttributeNames:  expr.Names(),
					ExpressionAttributeValues: expr.Values(),
					KeyConditionExpression:    expr.KeyCondition(),
				}
				m.QueryAPIClient.On("Query", ctx, input).Return(queryOutputFixture(t), nil)
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
			name: "Query with partition key and filter",
			args: args{
				ctx:             context.Background(),
				tableName:       "TEST",
				partitionValue:  "TEST_PARTITION_VALUE_1",
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
				keyCondition := expression.KeyEqual(
					expression.Key("TEST_PARTITION_ATTRIBUTE"),
					expression.Value("TEST_PARTITION_VALUE_1"),
				)
				filterCondition := expression.Equal(
					expression.Name("TEST_ATTRIBUTE_2"),
					expression.Value(1),
				)
				builder := expression.NewBuilder().WithKeyCondition(keyCondition).WithCondition(filterCondition)
				expr, err := builder.Build()
				if err != nil {
					t.Fatalf("expression build error: %v", err)
				}
				input := &dynamodb.QueryInput{
					TableName:                 aws.String("TEST"),
					ExpressionAttributeNames:  expr.Names(),
					ExpressionAttributeValues: expr.Values(),
					KeyConditionExpression:    expr.KeyCondition(),
					FilterExpression:          expr.Condition(),
				}
				m.QueryAPIClient.On("Query", ctx, input).Return(queryOutputFixture(t), nil)
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
			name: "Query error",
			args: args{
				ctx:            context.Background(),
				tableName:      "TEST",
				partitionValue: "TEST_PARTITION_VALUE_1",
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
				condition := expression.KeyEqual(
					expression.Key("TEST_PARTITION_ATTRIBUTE"),
					expression.Value("TEST_PARTITION_VALUE_1"),
				)
				builder := expression.NewBuilder().WithKeyCondition(condition)
				expr, err := builder.Build()
				if err != nil {
					t.Fatalf("expression build error: %v", err)
				}
				input := &dynamodb.QueryInput{
					TableName:                 aws.String("TEST"),
					ExpressionAttributeNames:  expr.Names(),
					ExpressionAttributeValues: expr.Values(),
					KeyConditionExpression:    expr.KeyCondition(),
				}
				m.QueryAPIClient.On("Query", ctx, input).Return(nil, fmt.Errorf("query error"))
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
			err := i.Query(
				tt.args.ctx,
				w,
				tt.args.tableName,
				tt.args.partitionValue,
				tt.args.sortCondition,
				tt.args.filterCondition,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Query() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func queryOutputFixture(t *testing.T) *dynamodb.QueryOutput {
	t.Helper()
	return &dynamodb.QueryOutput{
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

func jsonFixture(t *testing.T, m []map[string]interface{}) string {
	t.Helper()
	b, err := json.MarshalIndent(m, "", strings.Repeat(" ", 2))
	if err != nil {
		t.Fatalf("json marshal error: %v", err)
	}
	return string(b) + "\n"
}
