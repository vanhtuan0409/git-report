package gitreport

import (
	"fmt"
	"strings"
)

type IReportGenerator interface {
	GenerateFromResults([]*Result) string
	GenerateFromCommits(*Result) string
}

type reportGenerator struct{}

func NewReportGenerator() IReportGenerator {
	return &reportGenerator{}
}

func (r *reportGenerator) GenerateFromResults(results []*Result) string {
	sb := new(strings.Builder)
	for _, result := range results {
		fmt.Fprint(sb, r.GenerateFromCommits(result))
	}
	return sb.String()
}

func (r *reportGenerator) GenerateFromCommits(result *Result) string {
	sb := new(strings.Builder)
	fmt.Fprintf(sb, "Repository: %s\n", result.Repo)
	for _, commit := range result.Commits {
		fmt.Fprintf(sb, "    - %s: %s\n", "time", commit.Message())
	}
	return sb.String()
}
