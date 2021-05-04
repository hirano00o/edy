package model

import (
	"bytes"
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
