package main

import (
	"os"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "Git Report"
	app.Usage = "Collect git commit messages and organize by days to create a daily report"
	app.Version = "v0.1.0"
	app.Action = generateReport

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
