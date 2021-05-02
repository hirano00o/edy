package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Client struct {
	config   aws.Config
	endpoint string
}

type DynamoDB interface {
	dynamodb.DescribeTableAPIClient
	dynamodb.ScanAPIClient
	dynamodb.QueryAPIClient
}

type NewClient interface {
	CreateInstance() DynamoDB
}

func New(context context.Context, options map[string]string) (*Client, error) {
	var optFns []func(*config.LoadOptions) error
	cli := new(Client)
	for k := range options {
		switch k {
		case "local":
			url := fmt.Sprintf("http://localhost:%s", options[k])
			cli.endpoint = url
			// Issue: It does not work endpoint setting.
			optFns = append(optFns,
				config.WithEndpointResolver(
					aws.EndpointResolverFunc(
						func(service, region string) (aws.Endpoint, error) {
							return aws.Endpoint{URL: url}, nil
						},
					),
				),
			)
		case "region":
			optFns = append(optFns, config.WithRegion(options[k]))
		case "profile":
			optFns = append(optFns, config.WithSharedConfigProfile(options[k]))
		}
	}
	c, err := config.LoadDefaultConfig(context, optFns...)
	if err != nil {
		return nil, err
	}
	cli.config = c
	return cli, nil
}

func (c *Client) CreateInstance() DynamoDB {
	if len(c.endpoint) != 0 {
		return dynamodb.NewFromConfig(c.config,
			dynamodb.WithEndpointResolver(dynamodb.EndpointResolverFromURL(c.endpoint)),
		)
	}
	return dynamodb.NewFromConfig(c.config)
}
