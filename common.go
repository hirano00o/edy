package edy

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

	"github.com/hirano00o/edy/client"
	"github.com/hirano00o/edy/model"
)

type clientKey string

const newClientKey clientKey = "client"

type Edy interface {
	Scan(
		ctx context.Context,
		w io.Writer,
		tableName string,
		filterCondition string,
	) error
	Query(
		ctx context.Context,
		w io.Writer,
		tableName string,
		partitionValue,
		sortCondition string,
		filterCondition string,
	) error
	DescribeTable(ctx context.Context, w io.Writer, tableName string) error
}

type Instance struct {
	client.NewClient
}

func NewEdyClient(c client.NewClient) Edy {
	return &Instance{
		c,
	}
}

type conditionalOperation int

const (
	notOperator conditionalOperation = iota + 1
	key
	operator
	value
	join
	logicalOperator
)

func analyseFilterCondition(
	condition string,
) (*expression.ConditionBuilder, error) {
	var c expression.ConditionBuilder
	var err error
	s := strings.Split(condition, " ")
	var op model.ComparisonOperator
	var conditionKey string
	var conditionKeyType model.AttributeType
	var conditionValue []string
	notCondition := false
	lOp := model.LogicalOperator(0)
	nextState := notOperator

	for i := 0; i < len(s); i++ {
		switch nextState {
		case notOperator:
			if strings.Compare(s[i], "not") == 0 {
				notCondition = true
			} else {
				i--
			}
			nextState = key
		case key:
			keyT := strings.Split(s[i], ",")
			if len(keyT) != 2 {
				return nil, fmt.Errorf("invalid condition, no key type specified: %s", s[i])
			}
			conditionKey = keyT[0]
			conditionKeyType = model.AttributeTypeStr(keyT[1]).Name()
			nextState = operator
		case operator:
			op, err = model.ConvertToComparisonOperator(s[i])
			if err != nil {
				return nil, err
			}
			if op == model.EXISTS {
				nextState = join
			} else {
				nextState = value
			}
		case value:
			switch {
			case op == model.ComparisonOperator(0):
				return nil, fmt.Errorf("unknown condition error: %s", condition)
			case op == model.EXISTS:
				i--
				nextState = join
			case op == model.BETWEEN:
				if len(conditionValue) < 2 {
					conditionValue = append(conditionValue, s[i])
				}
				if len(conditionValue) == 2 {
					nextState = join
				}
			case op == model.IN:
				_, err = model.ConvertToLogicalOperator(s[i])
				if err == nil {
					i--
					nextState = join
				} else {
					conditionValue = append(conditionValue, s[i])
				}
			default:
				conditionValue = append(conditionValue, s[i])
				nextState = join
			}
		case join:
			i--
			_c, err := makeExpression(op, conditionKeyType, conditionValue, conditionKey, notCondition)
			if err != nil {
				return nil, err
			}
			switch lOp {
			case model.AND:
				c = c.And(*_c)
			case model.OR:
				c = c.Or(*_c)
			default:
				c = *_c
			}
			op = model.ComparisonOperator(0)
			conditionKey = ""
			conditionKeyType = nil
			conditionValue = nil
			notCondition = false
			nextState = logicalOperator
		case logicalOperator:
			lOp, err = model.ConvertToLogicalOperator(s[i])
			if err != nil {
				return nil, fmt.Errorf("invalid condition, unknown logical operator: %s", s[i])
			}
			nextState = notOperator
		default:
			return nil, fmt.Errorf("unknown condition error: %s", condition)
		}
	}
	if nextState != join && op != model.IN {
		return nil, fmt.Errorf("invalid condition: %s", condition)
	}
	_c, err := makeExpression(op, conditionKeyType, conditionValue, conditionKey, notCondition)
	if err != nil {
		return nil, err
	}
	switch lOp {
	case model.AND:
		c = c.And(*_c)
	case model.OR:
		c = c.Or(*_c)
	default:
		c = *_c
	}

	return &c, nil
}

func makeExpression(
	op model.ComparisonOperator,
	conditionKeyType model.AttributeType,
	conditionValue []string,
	conditionKey string,
	notCondition bool,
) (*expression.ConditionBuilder, error) {
	var c expression.ConditionBuilder
	switch op {
	case model.EQ:
		v, err := conditionKeyType.Value(conditionValue[0])
		if err != nil {
			return nil, fmt.Errorf("invalid condition, %v", err)
		}
		c = expression.Equal(expression.Name(conditionKey), expression.Value(v))
	case model.LE:
		v, err := conditionKeyType.Value(conditionValue[0])
		if err != nil {
			return nil, fmt.Errorf("invalid condition, %v", err)
		}
		c = expression.LessThanEqual(expression.Name(conditionKey), expression.Value(v))
	case model.LT:
		v, err := conditionKeyType.Value(conditionValue[0])
		if err != nil {
			return nil, fmt.Errorf("invalid condition, %v", err)
		}
		c = expression.LessThan(expression.Name(conditionKey), expression.Value(v))
	case model.GE:
		v, err := conditionKeyType.Value(conditionValue[0])
		if err != nil {
			return nil, fmt.Errorf("invalid condition, %v", err)
		}
		c = expression.GreaterThanEqual(expression.Name(conditionKey), expression.Value(v))
	case model.GT:
		v, err := conditionKeyType.Value(conditionValue[0])
		if err != nil {
			return nil, fmt.Errorf("invalid condition, %v", err)
		}
		c = expression.GreaterThan(expression.Name(conditionKey), expression.Value(v))
	case model.BeginsWith:
		v, err := conditionKeyType.Value(conditionValue[0])
		if err != nil {
			return nil, fmt.Errorf("invalid condition, %v", err)
		}
		c = expression.BeginsWith(expression.Name(conditionKey), v.(string))
	case model.BETWEEN:
		v1, err := conditionKeyType.Value(conditionValue[0])
		if err != nil {
			return nil, fmt.Errorf("invalid condition, %v", err)
		}
		v2, err := conditionKeyType.Value(conditionValue[1])
		if err != nil {
			return nil, fmt.Errorf("invalid condition, %v", err)
		}
		c = expression.Between(expression.Name(conditionKey), expression.Value(v1), expression.Value(v2))
	case model.CONTAINS:
		v, err := conditionKeyType.Value(conditionValue[0])
		if err != nil {
			return nil, fmt.Errorf("invalid condition, %v", err)
		}
		c = expression.Contains(expression.Name(conditionKey), v.(string))
	case model.IN:
		v := make([]expression.OperandBuilder, len(conditionValue))
		for i := range conditionValue {
			_v, err := conditionKeyType.Value(conditionValue[i])
			if err != nil {
				return nil, fmt.Errorf("invalid condition, %v", err)
			}
			v[i] = expression.Value(_v)
		}
		c = expression.In(expression.Name(conditionKey), v[0], v...)
	case model.EXISTS:
		if notCondition {
			c = expression.AttributeNotExists(expression.Name(conditionKey))
		} else {
			c = expression.AttributeExists(expression.Name(conditionKey))
		}
	}
	if notCondition && op != model.EXISTS {
		c = c.Not()
	}

	return &c, nil
}
