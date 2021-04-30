package main

import (
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/hirano00o/edy"
	"github.com/hirano00o/edy/client"
	"github.com/hirano00o/edy/meta"
)

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		log.SetFlags(0)
		log.Fatalln(err)
	}

}

func run(w io.Writer, args []string) error {
	baseOptions := []cli.Flag{
		&cli.StringFlag{
			Name:     "table-name",
			Usage:    "DynamoDB table name",
			Aliases:  []string{"t"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "region",
			Usage:   "AWS region",
			Aliases: []string{"r"},
		},
		&cli.StringFlag{
			Name:    "profile",
			Usage:   "AWS profile name",
			Aliases: []string{"p"},
		},
		&cli.StringFlag{
			Name:  "local",
			Usage: "Connect to localhost of specified port number. --local 8000",
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
		Commands: []*cli.Command{
			{
				Name:    "describe",
				Usage:   "Describe specified table",
				Aliases: []string{"d"},
				Flags:   baseOptions,
				Action:  describeCmd(w),
			},
			{
				Name:    "scan",
				Usage:   "Scan specified table",
				Aliases: []string{"s"},
				Flags:   append(baseOptions, scanQueryOptions...),
				Action:  scanCmd(w),
			},
			{
				Name:    "query",
				Usage:   "Query specified table",
				Aliases: []string{"q"},
				Flags:   append(baseOptions, scanQueryOptions...),
				Action:  queryCmd(w),
			},
		},
	}
	return app.Run(args)
}

func queryCmd(w io.Writer) cli.ActionFunc {
	return func(context *cli.Context) error {
		return nil
	}
}

func scanCmd(w io.Writer) cli.ActionFunc {
	return func(context *cli.Context) error {
		return nil
	}
}

func describeCmd(w io.Writer) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		c, err := client.New(ctx.Context, getOptions(ctx))
		if err != nil {
			return err
		}
		return edy.NewEdyClient(c).DescribeTable(ctx.Context, w, ctx.String("table-name"))
	}
}

func getOptions(context *cli.Context) map[string]string {
	o := make(map[string]string)

	// Get endpoint url.
	if p := context.String("local"); len(p) != 0 {
		o["local"] = p
	}

	// Get region.
	if r := context.String("region"); len(r) != 0 {
		o["region"] = r
	}

	// Get profile.
	if p := context.String("profile"); len(p) != 0 {
		o["profile"] = p
	}

	return o
}
