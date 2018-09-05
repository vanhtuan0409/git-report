package main

import (
	"fmt"

	gitreport "github.com/vanhtuan0409/git-report"
	cli "gopkg.in/urfave/cli.v1"
)

func showConfig(c *cli.Context) error {
	configPath := gitreport.GetDefaultConfigPath()
	config, err := gitreport.ReadConfigFromFile(configPath)
	if err != nil {
		return err
	}

	fmt.Println(config.ToString())
	return nil
}
