package edy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/hirano00o/edy/client"
)

const (
	NONE = iota
	STRING
	INT
	FLOAT64
	BOOL
	NULL
	LIST
)

func listType(items []interface{}) int {
	typ := NONE
	for i := range items {
		switch items[i].(type) {
		case string:
			if typ == NONE {
				typ = STRING
			} else if typ != STRING {
				return LIST
			}
		case int:
			if typ == NONE {
				typ = INT
			} else if typ != INT {
				return LIST
			}
		case float64:
			if typ == NONE {
				typ = FLOAT64
			} else if typ != FLOAT64 {
				return LIST
			}
		case bool:
			if typ == NONE {
				typ = BOOL
			} else if typ != BOOL {
				return LIST
			}
		case nil:
			if typ == NONE {
				typ = NULL
			} else if typ != NULL {
				return LIST
			}
		}
	}
	return typ
}

func setAttrEachType(item interface{}) (types.AttributeValue, error) {
	switch t := item.(type) {
	case map[string]interface{}:
		mm, err := recursiveAnalyseJSON(t)
		if err != nil {
			return nil, err
		}
		return &types.AttributeValueMemberM{
			Value: mm,
		}, nil
	case string:
		return &types.AttributeValueMemberS{
			Value: t,
		}, nil
	case int:
		return &types.AttributeValueMemberN{
			Value: strconv.Itoa(t),
		}, nil
	case float64:
		return &types.AttributeValueMemberN{
			Value: strconv.FormatFloat(t, 'f', -1, 64),
		}, nil
	case bool:
		return &types.AttributeValueMemberBOOL{
			Value: t,
		}, nil
	case nil:
		return &types.AttributeValueMemberNULL{
			Value: true,
		}, nil
	case []interface{}:
		typ := listType(t)
		switch typ {
		case LIST, BOOL, NULL:
			m := make([]types.AttributeValue, len(t))
			for i := range t {
				v, err := setAttrEachType(t[i])
				if err != nil {
					return nil, err
				}
				m[i] = v
			}
			return &types.AttributeValueMemberL{
				Value: m,
			}, nil
		case STRING:
			ss := make([]string, len(t))
			for i := range t {
				ss[i] = t[i].(string)
			}
			return &types.AttributeValueMemberSS{
				Value: ss,
			}, nil
		case INT:
			ss := make([]string, len(t))
			for i := range t {
				ss[i] = strconv.Itoa(t[i].(int))
			}
			return &types.AttributeValueMemberNS{
				Value: ss,
			}, nil
		case FLOAT64:
			ss := make([]string, len(t))
			for i := range t {
				ss[i] = strconv.FormatFloat(t[i].(float64), 'f', -1, 64)
			}
			return &types.AttributeValueMemberNS{
				Value: ss,
			}, nil
		default:
			return nil, fmt.Errorf("unsupported type or invalid type: %v", t)
		}
	default:
		return nil, fmt.Errorf("unsupported type or invalid type: %v", t)
	}
}

func recursiveAnalyseJSON(items map[string]interface{}) (map[string]types.AttributeValue, error) {
	m := make(map[string]types.AttributeValue)
	for k := range items {
		v, err := setAttrEachType(items[k])
		if err != nil {
			return nil, err
		}
		m[k] = v
	}
	return m, nil
}

func analyseItem(item string) (map[string]types.AttributeValue, error) {
	var jsonItem map[string]interface{}
	err := json.Unmarshal(bytes.NewBufferString(item).Bytes(), &jsonItem)
	if err != nil {
		return nil, fmt.Errorf("invalid json format: %v", err)
	}

	return recursiveAnalyseJSON(jsonItem)
}

func put(ctx context.Context, tableName, item string) (map[string]int, error) {
	cli := ctx.Value(newClientKey).(client.DynamoDB)

	i, err := analyseItem(item)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      i,
	}
	res, err := cli.PutItem(ctx, input)
	if err != nil {
		return nil, err
	}
	log.Println(res.ResultMetadata)
	return map[string]int{"unprocessed": 0}, nil
}

func (i *Instance) Put(ctx context.Context, w io.Writer, tableName, item string) error {
	cli := i.NewClient.CreateInstance()
	ctx = context.WithValue(ctx, newClientKey, cli)

	res, err := put(ctx, tableName, item)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(res, "", strings.Repeat(" ", 2))
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", string(b))

	return nil
}
