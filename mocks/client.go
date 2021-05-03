package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/hirano00o/edy/client"
)

type MockDynamoDBAPI struct {
	mock.Mock
	DescribeTableAPIClient
	QueryAPIClient
	ScanAPIClient
}

func (_m *MockDynamoDBAPI) CreateInstance() client.DynamoDB {
	ret := _m.Called()
	return ret.Get(0).(client.DynamoDB)
}

func (_m *DescribeTableAPIClient) CreateInstance() client.DynamoDB {
	ret := _m.Called()
	return ret.Get(0).(client.DynamoDB)
}

func (_m *QueryAPIClient) CreateInstance() client.DynamoDB {
	ret := _m.Called()
	return ret.Get(0).(client.DynamoDB)
}

func (_m *ScanAPIClient) CreateInstance() client.DynamoDB {
	ret := _m.Called()
	return ret.Get(0).(client.DynamoDB)
}
