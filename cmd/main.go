package main

import (
	"fmt"
	"log"
	"os"

	"github.com/marafu/nova8-scripts/cmd/checkmarx"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "checkmarx-cli",
		Usage: "Make an integration in Checkmarx API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `config.yml`",
			},
			&cli.StringFlag{
				Name:    "project-name",
				Aliases: []string{"pn"},
				Usage:   "Checkmarx project name",
			},
		},
		Action: func(ctx *cli.Context) error {

			if ctx.NArg() > 0 {
				log.Fatalln("usage -h")
			}

			config := ctx.String("config")
			projectName := ctx.String("project-name")

			if config == "" && projectName == "" {
				return fmt.Errorf("usage -h")
			}

			checkmarx.App(config, projectName)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
