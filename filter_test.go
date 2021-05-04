package edy

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func Test_analyseFilterCondition(t *testing.T) {
	type args struct {
		condition string
	}
	tests := []struct {
		name    string
		args    args
		want    expression.ConditionBuilder
		wantErr bool
	}{
		{
			name: "EQ case",
			args: args{
				condition: "ID,S = 1234",
			},
			want: expression.Equal(expression.Name("ID"), expression.Value("1234")),
		},
		{
			name: "LE case",
			args: args{
				condition: "ID,S <= 1234",
			},
			want: expression.LessThanEqual(expression.Name("ID"), expression.Value("1234")),
		},
		{
			name: "LT case",
			args: args{
				condition: "ID,S < 1234",
			},
			want: expression.LessThan(expression.Name("ID"), expression.Value("1234")),
		},
		{
			name: "GE case",
			args: args{
				condition: "ID,S >= 1234",
			},
			want: expression.GreaterThanEqual(expression.Name("ID"), expression.Value("1234")),
		},
		{
			name: "GT case",
			args: args{
				condition: "ID,S > 1234",
			},
			want: expression.GreaterThan(expression.Name("ID"), expression.Value("1234")),
		},
		{
			name: "BeginsWith case",
			args: args{
				condition: "ID,S begins_with 1234",
			},
			want: expression.BeginsWith(expression.Name("ID"), "1234"),
		},
		{
			name: "Between case",
			args: args{
				condition: "ID,S between 12 34",
			},
			want: expression.Between(expression.Name("ID"), expression.Value("12"), expression.Value("34")),
		},
		{
			name: "Contains case",
			args: args{
				condition: "ID,S contains 1234",
			},
			want: expression.Contains(expression.Name("ID"), "1234"),
		},
		{
			name: "IN case",
			args: args{
				condition: "ID,S in 1 2 3 4",
			},
			want: expression.In(
				expression.Name("ID"),
				expression.Value("1"),
				[]expression.OperandBuilder{
					expression.Value("1"),
					expression.Value("2"),
					expression.Value("3"),
					expression.Value("4"),
				}...,
			),
		},
		{
			name: "Exists case",
			args: args{
				condition: "ID,S exists",
			},
			want: expression.AttributeExists(expression.Name("ID")),
		},
		{
			name: "not EQ case",
			args: args{
				condition: "not ID,S = 1234",
			},
			want: expression.Equal(expression.Name("ID"), expression.Value("1234")).Not(),
		},
		{
			name: "not Between case",
			args: args{
				condition: "not ID,S between 12 34",
			},
			want: expression.Between(expression.Name("ID"), expression.Value("12"), expression.Value("34")).Not(),
		},
		{
			name: "not In case",
			args: args{
				condition: "not ID,S in 1 2 3 4",
			},
			want: expression.In(
				expression.Name("ID"),
				expression.Value("1"),
				[]expression.OperandBuilder{
					expression.Value("1"),
					expression.Value("2"),
					expression.Value("3"),
					expression.Value("4"),
				}...,
			).Not(),
		},
		{
			name: "not Exists case",
			args: args{
				condition: "not ID,S exists",
			},
			want: expression.AttributeNotExists(expression.Name("ID")),
		},
		{
			name: "not EQ and EQ case",
			args: args{
				condition: "not ID,S = 1234 and Name,S = user1",
			},
			want: expression.Equal(expression.Name("ID"), expression.Value("1234")).Not().And(
				expression.Equal(expression.Name("Name"), expression.Value("user1"))),
		},
		{
			name: "not EQ and not In case",
			args: args{
				condition: "not ID,S = 1234 and not Name,S in user1 user2 user3 user4",
			},
			want: expression.Equal(expression.Name("ID"), expression.Value("1234")).Not().And(
				expression.In(
					expression.Name("Name"),
					expression.Value("user1"),
					[]expression.OperandBuilder{
						expression.Value("user1"),
						expression.Value("user2"),
						expression.Value("user3"),
						expression.Value("user4"),
					}...,
				).Not(),
			),
		},
		{
			name: "not EQ and Exists case",
			args: args{
				condition: "not ID,S = 1234 and Name,S exists",
			},
			want: expression.Equal(expression.Name("ID"), expression.Value("1234")).Not().And(
				expression.AttributeExists(expression.Name("Name"))),
		},
		{
			name: "not In or not EQ case",
			args: args{
				condition: "not Name,S in user1 user2 user3 user4 or not ID,S = 1234",
			},
			want: expression.In(
				expression.Name("Name"),
				expression.Value("user1"),
				[]expression.OperandBuilder{
					expression.Value("user1"),
					expression.Value("user2"),
					expression.Value("user3"),
					expression.Value("user4"),
				}...,
			).Not().Or(
				expression.Equal(expression.Name("ID"), expression.Value("1234")).Not()),
		},
		{
			name: "Exists or not EQ case",
			args: args{
				condition: "Name,S exists or not ID,S = 1234",
			},
			want: expression.AttributeExists(expression.Name("Name")).Or(
				expression.Equal(expression.Name("ID"), expression.Value("1234")).Not()),
		},
		{
			name: "Exists or not EQ case and not In case",
			args: args{
				condition: "Age,N exists or not ID,S = 1234 and not Name,S in user1 user2 user3 user4",
			},
			want: expression.AttributeExists(expression.Name("Age")).Or(
				expression.Equal(expression.Name("ID"), expression.Value("1234")).Not()).And(
				expression.In(
					expression.Name("Name"),
					expression.Value("user1"),
					[]expression.OperandBuilder{
						expression.Value("user1"),
						expression.Value("user2"),
						expression.Value("user3"),
						expression.Value("user4"),
					}...,
				).Not(),
			),
		},
		{
			name: "Missing key type",
			args: args{
				condition: "ID = 1234",
			},
			wantErr: true,
		},
		{
			name: "Invalid comparison operator",
			args: args{
				condition: "ID,S => 1234",
			},
			wantErr: true,
		},
		{
			name: "Invalid logical operator",
			args: args{
				condition: "ID,S = 1234 && Name,S = user1",
			},
			wantErr: true,
		},
		{
			name: "Invalid number value",
			args: args{
				condition: "ID,N = a1234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := analyseFilterCondition(tt.args.condition)
			if (err != nil) != tt.wantErr {
				t.Errorf("analyseFilterCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got == nil {
				return
			}
			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("analyseFilterCondition() got = %v, want %v", got, tt.want)
			}
		})
	}
}
