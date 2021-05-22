package edy

import (
	"context"
	"io"

	"github.com/hirano00o/edy/client"
)

type clientKey string

const newClientKey clientKey = "client"

type Edy interface {
	Scan(ctx context.Context, w io.Writer, tableName, filterCondition, projection string) error
	Query(
		ctx context.Context,
		w io.Writer,
		tableName,
		partitionValue,
		sortCondition,
		filterCondition,
		index,
		projection string,
	) error
	DescribeTable(ctx context.Context, w io.Writer, tableName string) error
	Put(ctx context.Context, w io.Writer, tableName, item string) error
	Delete(ctx context.Context, w io.Writer, tableName, partitionValue, sortValue string) error
}

type Instance struct {
	client.NewClient
}
