package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"github.com/kanjih/ratchet/handler"
)

func main() {
	app := cli.NewApp()
	app.Usage = "Spanner migration tool"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:  "init",
			Usage: "create 'Migrations' table (which manages migrations) if not exists",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     handler.FlagNameProjectId,
					Usage:    "gcp project id",
					Required: true,
				},
				&cli.StringFlag{
					Name:     handler.FlagNameInstanceName,
					Usage:    "spanner instance name",
					Required: true,
				},
				&cli.StringFlag{
					Name:     handler.FlagNameDatabaseName,
					Usage:    "spanner database name",
					Required: true,
				},
			},
			Action: handler.Init,
		},
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "make new migration",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  handler.FlagNameDml,
					Usage: "specify to make dml migration",
				},
				&cli.BoolFlag{
					Name:  handler.FlagNamePartitionedDml,
					Usage: "specify to make partitioned dml migration",
				},
			},
			Action: handler.New,
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run migrations",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     handler.FlagNameProjectId,
					Usage:    "gcp project id",
					Required: true,
				},
				&cli.StringFlag{
					Name:     handler.FlagNameInstanceName,
					Usage:    "spanner instance name",
					Required: true,
				},
				&cli.StringFlag{
					Name:     handler.FlagNameDatabaseName,
					Usage:    "spanner database name",
					Required: true,
				},
			},
			Action: handler.Run,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("[ERROR] " + err.Error())
		os.Exit(1)
	}
}
