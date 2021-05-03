package edy

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

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
