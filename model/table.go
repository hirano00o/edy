package model

type Table struct {
	Arn              string
	Name             string
	PartitionKeyName string
	PartitionKeyType string
	SortKeyName      string
	SortKeyType      string
	GSI              []*GlobalSecondaryIndex
	ItemCount        int64
}

type GlobalSecondaryIndex struct {
	Name             string
	PartitionKeyName string
	PartitionKeyType string
	SortKeyName      string
	SortKeyType      string
}
