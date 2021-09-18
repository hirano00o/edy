package edy

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func analyseProjection(projection string) *expression.ProjectionBuilder {
	p := regexp.MustCompile(`,[\s]*|\s+`).Split(strings.TrimSpace(projection), -1)
	var pj expression.ProjectionBuilder
	for i := range p {
		pj = expression.AddNames(pj, expression.Name(p[i]))
	}
	return &pj
}
