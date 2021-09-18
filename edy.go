package edy

import (
	"context"
	"io"

	"github.com/hirano00o/edy/client"
)

type clientKey string

const newClientKey clientKey = "client"

type Edy interface {
	Scan(ctx context.Context, w io.Writer, tableName, filterCondition, projection, output string) error
	Query(
		ctx context.Context,
		w io.Writer,
		tableName,
		partitionValue,
		sortCondition,
		filterCondition,
		index,
		projection string,
		output string,
	) error
	DescribeTable(ctx context.Context, w io.Writer, tableName string) error
	Put(ctx context.Context, w io.Writer, tableName, item, fileName string, f func(string) (string, error)) error
	Delete(
		ctx context.Context,
		w io.Writer,
		tableName,
		partitionValue,
		sortValue,
		fileName string,
		f func(string) (string, error),
	) error
}

type Instance struct {
	client.NewClient
}
