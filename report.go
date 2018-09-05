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
	groups := groupByDay(result.Commits)
	for _, g := range groups {
		fmt.Fprintf(sb, "  + %s\n", g.name)
		for _, commit := range g.commits {
			fmt.Fprintf(sb, "    - %s\n", commit.Message())
		}
	}

	return sb.String()
}

type group struct {
	name    string
	commits []*GitCommit
}

func groupByDay(commits []*GitCommit) []*group {
	groups := []*group{}
	for _, commit := range commits {
		dayStr := commit.Author.Date.Format("2006-01-02")
		if len(groups) == 0 || groups[len(groups)-1].name != dayStr {
			groups = append(groups, &group{
				name:    dayStr,
				commits: []*GitCommit{commit},
			})
		} else {
			lastGroup := groups[len(groups)-1]
			lastGroup.commits = append(lastGroup.commits, commit)
		}
	}
	return groups
}
