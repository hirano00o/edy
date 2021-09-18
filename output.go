package edy

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type formatType int

const (
	jsonType formatType = iota
	csvType
)

var formatTypeMap = map[string]formatType{
	"json": jsonType,
	"csv":  csvType,
}

func getKeyOrder(data []map[string]interface{}) []string {
	var order []string
	orderMap := make(map[string]struct{})
	for i := range data {
		for k := range data[i] {
			if _, ok := orderMap[k]; !ok {
				orderMap[k] = struct{}{}
				order = append(order, k)
			}
		}
	}
	sort.Strings(order)
	return order
}

func adjustSpecifiedFormat(outputFormat string, data []map[string]interface{}) (string, error) {
	switch formatTypeMap[strings.ToLower(outputFormat)] {
	case csvType:
		var b bytes.Buffer
		writer := csv.NewWriter(&b)
		keys := getKeyOrder(data)
		err := writer.Write(keys)
		if err != nil {
			return "", err
		}
		for i := range data {
			var record []string
			for k := range keys {
				record = append(record, fmt.Sprint(data[i][keys[k]]))
			}
			err := writer.Write(record)
			if err != nil {
				return "", err
			}
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			return "", err
		}
		return b.String(), nil
	default:
		b, err := json.MarshalIndent(data, "", strings.Repeat(" ", 2))
		if err != nil {
			return "", err
		}
		return string(b) + "\n", nil
	}
}
