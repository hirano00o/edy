package edy

import (
	"context"
	"io"

	"github.com/hirano00o/edy/client"
)

type clientKey string

const newClientKey clientKey = "client"

type Edy interface {
	Query(
		ctx context.Context,
		w io.Writer,
		tableName string,
		partitionValue,
		sortValue interface{},
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
