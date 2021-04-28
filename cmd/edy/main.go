package main

import (
	"log"
	"os"

	"github.com/hirano00o/edy/meta"

	"github.com/urfave/cli/v2"
)

func main() {
	if err := run(os.Args); err != nil {
		log.SetFlags(0)
		log.Fatalln(err)
	}

}

func run(args []string) error {
	baseOptions := []cli.Flag{
		&cli.StringFlag{
			Name:    "table-name",
			Usage:   "DynamoDB table name",
			Aliases: []string{"t"},
		},
		&cli.StringFlag{
			Name:    "region",
			Usage:   "Region",
			Aliases: []string{"r"},
		},
	}
	scanQueryOptions := []cli.Flag{
		&cli.StringFlag{
			Name:    "filter",
			Usage:   "Filter by specified condition",
			Aliases: []string{"f"},
		},
	}
	app := &cli.App{
		Name:    meta.CliName,
		Version: meta.Version,
		Usage:   "Easy to use DynamoDB CLI",
		Flags:   append(baseOptions, scanQueryOptions...),
		Commands: []*cli.Command{
			{
				Name:    "scan",
				Usage:   "Scan specified table",
				Aliases: []string{"s"},
				Flags:   baseOptions,
				Action:  scanCmd(),
			},
			{
				Name:    "query",
				Usage:   "Query specified table",
				Aliases: []string{"q"},
				Flags:   baseOptions,
				Action:  queryCmd(),
			},
		},
	}
	return app.Run(args)
}

func queryCmd() cli.ActionFunc {
	return func(context *cli.Context) error {
		return nil
	}
}

func scanCmd() cli.ActionFunc {
	return func(context *cli.Context) error {
		return nil
	}
}
