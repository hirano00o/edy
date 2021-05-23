package edy

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/hirano00o/edy/model"
)

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
			if (conditionKeyType.String() == new(model.SS).String() ||
				conditionKeyType.String() == new(model.NS).String()) &&
				!(op == model.IN || op == model.EQ || op == model.EXISTS) {
				return nil, fmt.Errorf("%s operand can not use type %s", op.String(), conditionKeyType.String())
			}
			if op == model.EXISTS {
				nextState = join
			} else {
				nextState = value
			}
		case value:
			switch {
			case op == model.BETWEEN:
				if len(conditionValue) < 2 {
					conditionValue = append(conditionValue, s[i])
				}
				if len(conditionValue) == 2 {
					nextState = join
				}
			case op == model.IN || (op == model.EQ &&
				(conditionKeyType.String() == new(model.SS).String() ||
					conditionKeyType.String() == new(model.NS).String())):
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
	if nextState != join && op != model.IN &&
		!(op == model.EQ &&
			(conditionKeyType.String() == new(model.SS).String() ||
				conditionKeyType.String() == new(model.NS).String())) {
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

func makeExpressionValue(
	op model.ComparisonOperator,
	conditionKeyType model.AttributeType,
	conditionValue []string,
) (interface{}, error) {
	switch conditionKeyType.String() {
	case new(model.SS).String():
		return []expression.OperandBuilder{
			expression.Value(
				&types.AttributeValueMemberSS{Value: conditionValue},
			)}, nil
	case new(model.NS).String():
		return []expression.OperandBuilder{
			expression.Value(
				&types.AttributeValueMemberNS{Value: conditionValue},
			)}, nil
	default:
		v := make([]expression.OperandBuilder, len(conditionValue))
		for i := range conditionValue {
			cv, err := conditionKeyType.Value(conditionValue[i])
			if err != nil {
				return nil, fmt.Errorf("invalid condition, cannot convert key type: %v", err)
			}
			if op == model.CONTAINS || op == model.BeginsWith {
				return cv, nil
			}
			v[i] = expression.Value(cv)
		}
		return v, nil
	}
}

func makeExpression(
	op model.ComparisonOperator,
	conditionKeyType model.AttributeType,
	conditionValue []string,
	conditionKey string,
	notCondition bool,
) (*expression.ConditionBuilder, error) {
	var c expression.ConditionBuilder
	v, err := makeExpressionValue(op, conditionKeyType, conditionValue)
	if err != nil {
		return nil, err
	}

	switch op {
	case model.EQ:
		c = expression.Equal(expression.Name(conditionKey), v.([]expression.OperandBuilder)[0])
	case model.NE:
		c = expression.NotEqual(expression.Name(conditionKey), v.([]expression.OperandBuilder)[0])
	case model.LE:
		c = expression.LessThanEqual(expression.Name(conditionKey), v.([]expression.OperandBuilder)[0])
	case model.LT:
		c = expression.LessThan(expression.Name(conditionKey), v.([]expression.OperandBuilder)[0])
	case model.GE:
		c = expression.GreaterThanEqual(expression.Name(conditionKey), v.([]expression.OperandBuilder)[0])
	case model.GT:
		c = expression.GreaterThan(expression.Name(conditionKey), v.([]expression.OperandBuilder)[0])
	case model.BeginsWith:
		c = expression.BeginsWith(expression.Name(conditionKey), v.(string))
	case model.BETWEEN:
		c = expression.Between(
			expression.Name(conditionKey),
			v.([]expression.OperandBuilder)[0],
			v.([]expression.OperandBuilder)[1],
		)
	case model.CONTAINS:
		c = expression.Contains(expression.Name(conditionKey), v.(string))
	case model.IN:
		vv := v.([]expression.OperandBuilder)
		if len(vv) == 1 {
			c = expression.In(expression.Name(conditionKey), vv[0])
		} else {
			c = expression.In(expression.Name(conditionKey), vv[0], vv...)
		}
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
