package model

import (
	"bytes"
	"fmt"
	"strconv"
)

type AttributeType interface {
	String() string
	Value(s string) (interface{}, error)
}

type S struct{}
type N struct{}
type B struct{}

func (S) String() string {
	return "S"
}

func (S) Value(s string) (interface{}, error) {
	return s, nil
}

func (N) String() string {
	return "N"
}

func (N) Value(s string) (interface{}, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (B) String() string {
	return "B"
}

func (B) Value(s string) (interface{}, error) {
	return bytes.NewBufferString(s).Bytes(), nil
}

var (
	s S = struct{}{}
	n N = struct{}{}
	b B = struct{}{}
)

type AttributeTypeStr string

func (a AttributeTypeStr) Name() AttributeType {
	switch a {
	case "S":
		return s
	case "N":
		return n
	case "B":
		return b
	default:
		return s
	}
}

type ComparisonOperator int

const (
	EQ ComparisonOperator = iota + 1
	NE
	LE
	LT
	GE
	GT
	NOT_NULL
	NULL
	CONTAINS
	NOT_CONTAINS
	BEGINS_WITH
	IN
	BETWEEN
)

var mapOptionComparisonOperator map[string]ComparisonOperator = map[string]ComparisonOperator{
	"=":            EQ,
	"!=":           NE,
	"<=":           LE,
	"<":            LT,
	">=":           GE,
	">":            GT,
	"not_null":     NOT_NULL,
	"null":         NULL,
	"contains":     CONTAINS,
	"not_contains": NOT_CONTAINS,
	"begins_with":  BEGINS_WITH,
	"in":           IN,
	"between":      BETWEEN,
}

func ConvertToComparisonOperator(op string) (ComparisonOperator, error) {
	if v, ok := mapOptionComparisonOperator[op]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("invalid comparison operator: %s", op)
}
