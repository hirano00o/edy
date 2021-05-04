package model

type Table struct {
	Arn          string                  `json:"tableArn"`
	Name         string                  `json:"tableName"`
	PartitionKey *Key                    `json:"partitionKey"`
	SortKey      *Key                    `json:"sortKey,omitempty"`
	GSI          []*GlobalSecondaryIndex `json:"gsi,omitempty"`
	ItemCount    int64                   `json:"itemCount"`
}

type Key struct {
	Name    string        `json:"name"`
	Type    AttributeType `json:"-"`
	TypeStr string        `json:"type"`
}

type GlobalSecondaryIndex struct {
	Name         string `json:"indexName"`
	PartitionKey *Key   `json:"partitionKey"`
	SortKey      *Key   `json:"sortKey,omitempty"`
}
