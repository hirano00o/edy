package model

type Table struct {
	Arn              string
	Name             string
	PartitionKeyName string
	PartitionKeyType AttributeType
	SortKeyName      string
	SortKeyType      AttributeType
	GSI              []*GlobalSecondaryIndex
	ItemCount        int64
}

type GlobalSecondaryIndex struct {
	Name             string
	PartitionKeyName string
	PartitionKeyType AttributeType
	SortKeyName      string
	SortKeyType      AttributeType
}
