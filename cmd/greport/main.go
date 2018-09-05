package main

import (
	"fmt"
	"os"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "Git Report"
	app.Usage = "Collect git commit messages and organize by days to create a daily report"
	app.Version = "v0.1.0"
	app.Commands = []cli.Command{
		{
			Name:   "generate",
			Usage:  "Generate daily report from commit messages",
			Action: generateReport,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "from, f",
					Value: "",
					Usage: "Get commits starting from date. Date format: YYYY-MM-DD",
				},
				cli.StringFlag{
					Name:  "to, t",
					Value: "",
					Usage: "Get commits ending to date. Date format: YYYY-MM-DD",
				},
			},
		},
		{
			Name:   "init",
			Usage:  "Initialize config file",
			Action: initConfig,
		},
		{
			Name:   "config",
			Usage:  "Show config file",
			Action: showConfig,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Encountered error: %s\n", err.Error())
		os.Exit(1)
	}
}
