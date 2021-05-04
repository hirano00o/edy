package model

import "fmt"

type ComparisonOperator int
type LogicalOperator int

const (
	EQ ComparisonOperator = iota + 1
	NE
	LE
	LT
	GE
	GT
	EXISTS
	CONTAINS
	BeginsWith
	IN
	BETWEEN
)

const (
	AND LogicalOperator = iota + 1
	OR
)

var mapComparisonOperator = map[string]ComparisonOperator{
	"=":           EQ,
	"!=":          NE,
	"<=":          LE,
	"<":           LT,
	">=":          GE,
	">":           GT,
	"exists":      EXISTS,
	"contains":    CONTAINS,
	"begins_with": BeginsWith,
	"in":          IN,
	"between":     BETWEEN,
}

var mapLogicalOperator = map[string]LogicalOperator{
	"and": AND,
	"or":  OR,
}

func ConvertToComparisonOperator(op string) (ComparisonOperator, error) {
	if v, ok := mapComparisonOperator[op]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("invalid comparison operator: %s", op)
}

func ConvertToLogicalOperator(op string) (LogicalOperator, error) {
	if v, ok := mapLogicalOperator[op]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("invalid logical operator: %s", op)
}
