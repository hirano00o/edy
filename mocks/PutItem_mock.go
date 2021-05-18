package mocks

import (
	"context"
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

type PutItemClient struct {
	mock.Mock
}

func (_m *PutItemClient) PutItem(
	_a0 context.Context,
	_a1 *dynamodb.PutItemInput,
	_a2 ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *dynamodb.PutItemOutput
	if rf, ok := ret.Get(0).(func(
		context.Context,
		*dynamodb.PutItemInput,
		...func(*dynamodb.Options,
		)) *dynamodb.PutItemOutput); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else if ret.Get(0) != nil {
		log.Println(reflect.TypeOf(ret.Get(0)))
		r0 = ret.Get(0).(*dynamodb.PutItemOutput)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
