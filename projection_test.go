package edy

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func Test_analyseProjection(t *testing.T) {
	type args struct {
		projection string
	}
	tests := []struct {
		name string
		args args
		want expression.ProjectionBuilder
	}{
		{
			name: "Space split case",
			args: args{
				projection: "PJ1 PJ2 PJ3",
			},
			want: expression.ProjectionBuilder{}.AddNames(
				expression.Name("PJ1"),
				expression.Name("PJ2"),
				expression.Name("PJ3"),
			),
		},
		{
			name: "Comma split case",
			args: args{
				projection: "PJ1,PJ2,PJ3",
			},
			want: expression.ProjectionBuilder{}.AddNames(
				expression.Name("PJ1"),
				expression.Name("PJ2"),
				expression.Name("PJ3"),
			),
		},
		{
			name: "Space and comma split case",
			args: args{
				projection: "PJ1,   PJ2, PJ3",
			},
			want: expression.ProjectionBuilder{}.AddNames(
				expression.Name("PJ1"),
				expression.Name("PJ2"),
				expression.Name("PJ3"),
			),
		},
		{
			name: "Last unnecessary space case",
			args: args{
				projection: "PJ1, PJ2, PJ3 ",
			},
			want: expression.ProjectionBuilder{}.AddNames(
				expression.Name("PJ1"),
				expression.Name("PJ2"),
				expression.Name("PJ3"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := analyseProjection(tt.args.projection); !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("analyseProjection() = %v, want %v", got, tt.want)
			}
		})
	}
}
