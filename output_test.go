package edy

import "testing"

func Test_adjustSpecifiedFormat(t *testing.T) {
	type args struct {
		outputFormat string
		data         []map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Output JSON",
			args: args{
				outputFormat: "JSON",
				data: []map[string]interface{}{
					{
						"TEST_ATTRIBUTE_1": "TEST_ATTRIBUTE_1_VALUE_1",
						"TEST_ATTRIBUTE_2": 21,
						"TEST_ATTRIBUTE_3": []string{"VALUE_11", "VALUE_12", "VALUE_13"},
					},
					{
						"TEST_ATTRIBUTE_1": "TEST_ATTRIBUTE_1_VALUE_2",
						"TEST_ATTRIBUTE_2": 22,
						"TEST_ATTRIBUTE_3": []string{"VALUE_21", "VALUE_22", "VALUE_23"},
					},
					{
						"TEST_ATTRIBUTE_1": "TEST_ATTRIBUTE_1_VALUE_3",
						"TEST_ATTRIBUTE_2": 23,
						"TEST_ATTRIBUTE_3": []string{"VALUE_31", "VALUE_32", "VALUE_33"},
					},
				},
			},
			want: jsonFixture(t,
				[]map[string]interface{}{
					{
						"TEST_ATTRIBUTE_1": "TEST_ATTRIBUTE_1_VALUE_1",
						"TEST_ATTRIBUTE_2": 21,
						"TEST_ATTRIBUTE_3": []string{"VALUE_11", "VALUE_12", "VALUE_13"},
					},
					{
						"TEST_ATTRIBUTE_1": "TEST_ATTRIBUTE_1_VALUE_2",
						"TEST_ATTRIBUTE_2": 22,
						"TEST_ATTRIBUTE_3": []string{"VALUE_21", "VALUE_22", "VALUE_23"},
					},
					{
						"TEST_ATTRIBUTE_1": "TEST_ATTRIBUTE_1_VALUE_3",
						"TEST_ATTRIBUTE_2": 23,
						"TEST_ATTRIBUTE_3": []string{"VALUE_31", "VALUE_32", "VALUE_33"},
					},
				}),
		},
		{
			name: "Output csv",
			args: args{
				outputFormat: "CSV",
				data: []map[string]interface{}{
					{
						"TEST_ATTRIBUTE_1": "TEST_ATTRIBUTE_1_VALUE_1",
						"TEST_ATTRIBUTE_2": 21,
						"TEST_ATTRIBUTE_3": []string{"VALUE_11", "VALUE_12", "VALUE_13"},
					},
					{
						"TEST_ATTRIBUTE_1": "TEST_ATTRIBUTE_1_VALUE_2",
						"TEST_ATTRIBUTE_2": 22,
						"TEST_ATTRIBUTE_3": []string{"VALUE_21", "VALUE_22", "VALUE_23"},
					},
					{
						"TEST_ATTRIBUTE_1": "TEST_ATTRIBUTE_1_VALUE_3",
						"TEST_ATTRIBUTE_2": 23,
						"TEST_ATTRIBUTE_3": []string{"VALUE_31", "VALUE_32", "VALUE_33"},
					},
				},
			},
			want: "TEST_ATTRIBUTE_1,TEST_ATTRIBUTE_2,TEST_ATTRIBUTE_3\n" +
				"TEST_ATTRIBUTE_1_VALUE_1,21,[VALUE_11 VALUE_12 VALUE_13]\n" +
				"TEST_ATTRIBUTE_1_VALUE_2,22,[VALUE_21 VALUE_22 VALUE_23]\n" +
				"TEST_ATTRIBUTE_1_VALUE_3,23,[VALUE_31 VALUE_32 VALUE_33]\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := adjustSpecifiedFormat(tt.args.outputFormat, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("adjustSpecifiedFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("adjustSpecifiedFormat() got = %v, want %v", got, tt.want)
			}
		})
	}
}
