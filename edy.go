package edy

import (
	"context"
	"io"

	"github.com/hirano00o/edy/client"
)

type clientKey string

const newClientKey clientKey = "client"

type Edy interface {
	Scan(
		ctx context.Context,
		w io.Writer,
		tableName,
		filterCondition string,
	) error
	Query(
		ctx context.Context,
		w io.Writer,
		tableName,
		partitionValue,
		sortCondition,
		filterCondition,
		index string,
	) error
	DescribeTable(ctx context.Context, w io.Writer, tableName string) error
}

type Instance struct {
	client.NewClient
}
