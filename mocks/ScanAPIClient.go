// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	mock "github.com/stretchr/testify/mock"
)

// ScanAPIClient is an autogenerated mock type for the ScanAPIClient type
type ScanAPIClient struct {
	mock.Mock
}

// Scan provides a mock function with given fields: _a0, _a1, _a2
func (_m *ScanAPIClient) Scan(_a0 context.Context, _a1 *dynamodb.ScanInput, _a2 ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *dynamodb.ScanOutput
	if rf, ok := ret.Get(0).(func(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) *dynamodb.ScanOutput); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.ScanOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
