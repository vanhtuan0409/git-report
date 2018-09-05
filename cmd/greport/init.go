package main

import (
	"fmt"
	"path/filepath"

	gitreport "github.com/vanhtuan0409/git-report"
	cli "gopkg.in/urfave/cli.v1"
)

func initConfig(c *cli.Context) error {
	configPath := gitreport.GetDefaultConfigPath()
	_, err := gitreport.CreateDefaultConfig(configPath)
	if err != nil {
		return err
	}
	fmt.Printf("Config files created at: %s\n", filepath.Dir(configPath))
	return nil
}
