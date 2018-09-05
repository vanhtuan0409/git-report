package main

import (
	"fmt"
	"time"

	gitreport "github.com/vanhtuan0409/git-report"
	cli "gopkg.in/urfave/cli.v1"
)

func generateReport(c *cli.Context) error {
	configPath := gitreport.GetDefaultConfigPath()
	config, err := gitreport.ReadConfigFromFile(configPath)
	if err != nil {
		return err
	}

	fromOption := c.String("from")
	fromValue, err := time.Parse("2006-01-02", fromOption)
	if err != nil {
		year, month, day := time.Now().AddDate(0, 0, -config.DefaultTimeRange).Date()
		fromValue = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	}

	toString := c.String("to")
	toValue, err := time.Parse("2006-01-02", toString)
	if err != nil {
		year, month, day := time.Now().Date()
		toValue = time.Date(year, month, day, 23, 59, 59, 0, time.Local)
	}

	resultChan := make(chan string)
	errChan := make(chan error)
	for _, repoPath := range config.Repos {
		go func(path string) {
			gitClient := gitreport.NewGitClient(path)
			result, err := gitClient.Log(&gitreport.LogOption{
				Authors:           config.FilterEmail,
				FetchAllBranch:    true,
				FilterMergeCommit: true,
				Since:             &fromValue,
				Until:             &toValue,
			})
			if err != nil {
				errChan <- fmt.Errorf("Cannot fetch git commits from url: %s. Original Error:\n%s", path, err.Error())
				return
			}

			generator := gitreport.NewReportGenerator()
			report := generator.GenerateFromCommits(result)
			resultChan <- report
		}(repoPath)
	}

	for i := 0; i < len(config.Repos); i++ {
		select {
		case result := <-resultChan:
			fmt.Println(result)
		case err := <-errChan:
			fmt.Println(err.Error())
		}
	}

	return nil
}
