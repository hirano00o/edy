package model

import (
	"bytes"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AttributeType interface {
	String() string
	Value(s string) (interface{}, error)
	ConvertValueMember(s string) types.AttributeValue
}

type S struct{}
type N struct{}
type B struct{}
type SS struct{}
type NS struct{}

func (S) String() string {
	return "S"
}

func (S) Value(s string) (interface{}, error) {
	return s, nil
}

func (S) ConvertValueMember(s string) types.AttributeValue {
	return &types.AttributeValueMemberS{
		Value: s,
	}
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

func (N) ConvertValueMember(s string) types.AttributeValue {
	return &types.AttributeValueMemberN{
		Value: s,
	}
}

func (B) String() string {
	return "B"
}

func (B) Value(s string) (interface{}, error) {
	return bytes.NewBufferString(s).Bytes(), nil
}

func (B) ConvertValueMember(s string) types.AttributeValue {
	return &types.AttributeValueMemberB{
		Value: bytes.NewBufferString(s).Bytes(),
	}
}

func (SS) String() string {
	return "SS"
}

func (SS) Value(s string) (interface{}, error) {
	return s, nil
}

func (SS) ConvertValueMember(s string) types.AttributeValue {
	return &types.AttributeValueMemberSS{
		Value: []string{s},
	}
}

func (NS) String() string {
	return "NS"
}

func (NS) Value(s string) (interface{}, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (NS) ConvertValueMember(s string) types.AttributeValue {
	return &types.AttributeValueMemberNS{
		Value: []string{s},
	}
}

var (
	s  S  = struct{}{}
	n  N  = struct{}{}
	b  B  = struct{}{}
	ss SS = struct{}{}
	ns NS = struct{}{}
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
	case "SS":
		return ss
	case "NS":
		return ns
	default:
		return s
	}
}
