package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/hirano00o/edy"
	"github.com/hirano00o/edy/client"
	"github.com/hirano00o/edy/meta"
)

var baseOptions = []cli.Flag{
	&cli.StringFlag{
		Name:     "table-name",
		Usage:    "DynamoDB table name.",
		Aliases:  []string{"t"},
		Required: true,
	},
	&cli.StringFlag{
		Name:    "region",
		Usage:   "AWS region.",
		Aliases: []string{"r"},
	},
	&cli.StringFlag{
		Name:  "profile",
		Usage: "AWS profile name.",
	},
	&cli.StringFlag{
		Name: "local",
		Usage: "Port number or full URL if you connect such as dynamodb-local and LocalStack.\n" +
			"\tex. --local 8000",
	},
}

var queryOptions = []cli.Flag{
	&cli.StringFlag{
		Name:     "partition",
		Usage:    "The value of partition key.",
		Aliases:  []string{"p"},
		Required: true,
	},
	&cli.StringFlag{
		Name: "sort",
		Usage: "The value and condition of sort key.\n" +
			"\tex1. --sort \"> 20\"\n" +
			"\tex2. --sort \"between 20 25\"\n" +
			"\tAvailable operator is =,<=,<,>=,>,between,begins_with",
		Aliases: []string{"s"},
	},
	&cli.StringFlag{
		Name:    "index",
		Usage:   "Global secondary index name.",
		Aliases: []string{"idx"},
	},
}

var scanQueryOptions = []cli.Flag{
	&cli.StringFlag{
		Name: "filter",
		Usage: "The condition if you use filter.\n" +
			"\tex. --filter \"Age,N >= 20 and Email,S in alice@example.com bob@example.com or not Birthplace,S exists\"\n" +
			"\tAvailable operator is =,<=,<,>=,>,between,begins_with,exists,in,contains",
		Aliases: []string{"f"},
	},
	&cli.StringFlag{
		Name: "projection",
		Usage: "Identifies and retrieve the attributes that you want.\n" +
			"\tex. --projection \"Age, Email, Birthplace\"",
		Aliases: []string{"pj"},
	},
	&cli.StringFlag{
		Name: "output",
		Usage: "Output format to show the result.\n" +
			"\tAvailable format is JSON, csv. Default is JSON",
		Aliases: []string{"o"},
	},
}

var putOptions = []cli.Flag{
	&cli.StringFlag{
		Name: "item",
		Usage: "Specify the item you want to create.\n" +
			"\tex. --item '{\"ID\":3,\"Name\":\"Alice\",\"Interest\":{\"SNS\":[\"Twitter\",\"Facebook\"]}}'",
		Aliases: []string{"i"},
	},
	&cli.StringFlag{
		Name:    "input-file",
		Usage:   "Read item to put from json file. Use either the --item option or this option.",
		Aliases: []string{"I"},
	},
}

var deleteOptions = []cli.Flag{
	&cli.StringFlag{
		Name:    "partition",
		Usage:   "The value of partition key.",
		Aliases: []string{"p"},
	},
	&cli.StringFlag{
		Name:    "sort",
		Usage:   "The value and sort key.",
		Aliases: []string{"s"},
	},
	&cli.StringFlag{
		Name:    "input-file",
		Usage:   "Read item to delete from json file. Use either the --partition (and --sort) option or this option.",
		Aliases: []string{"I"},
	},
}

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		log.SetFlags(0)
		log.Fatalln(err)
	}

}

func run(w io.Writer, args []string) error {
	app := &cli.App{
		Name:    meta.CliName,
		Version: meta.Version,
		Usage:   "Easy to use DynamoDB CLI",
		Commands: []*cli.Command{
			{
				Name:    "describe",
				Usage:   "Describe table",
				Aliases: []string{"desc"},
				Flags:   baseOptions,
				Action:  cmd(w),
			},
			{
				Name:    "scan",
				Usage:   "Scan table",
				Aliases: []string{"s"},
				Flags:   append(baseOptions, scanQueryOptions...),
				Action:  cmd(w),
			},
			{
				Name:    "query",
				Usage:   "Query table",
				Aliases: []string{"q"},
				Flags:   append(append(baseOptions, queryOptions...), scanQueryOptions...),
				Action:  cmd(w),
			},
			{
				Name:    "put",
				Usage:   "Put item",
				Aliases: []string{"p"},
				Flags:   append(baseOptions, putOptions...),
				Action:  cmd(w),
			},
			{
				Name:    "delete",
				Usage:   "Delete item",
				Aliases: []string{"del"},
				Flags:   append(baseOptions, deleteOptions...),
				Action:  cmd(w),
			},
		},
	}
	return app.Run(args)
}

func cmd(w io.Writer) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		c, err := client.New(ctx.Context, getOptions(ctx))
		if err != nil {
			return err
		}
		f := func(fileName string) (string, error) {
			b, err := ioutil.ReadFile(fileName)
			if err != nil {
				return "", nil
			}
			return string(b), nil
		}
		switch ctx.Command.Name {
		case "describe":
			return newEdyClient(c).DescribeTable(ctx.Context, w, ctx.String("table-name"))
		case "scan":
			return newEdyClient(c).Scan(
				ctx.Context,
				w,
				ctx.String("table-name"),
				ctx.String("filter"),
				ctx.String("projection"),
				ctx.String("output"),
			)
		case "query":
			return newEdyClient(c).Query(
				ctx.Context,
				w,
				ctx.String("table-name"),
				ctx.String("partition"),
				ctx.String("sort"),
				ctx.String("filter"),
				ctx.String("index"),
				ctx.String("projection"),
				ctx.String("output"),
			)
		case "put":
			return newEdyClient(c).Put(
				ctx.Context,
				w,
				ctx.String("table-name"),
				ctx.String("item"),
				ctx.String("input-file"),
				f,
			)
		case "delete":
			return newEdyClient(c).Delete(
				ctx.Context,
				w,
				ctx.String("table-name"),
				ctx.String("partition"),
				ctx.String("sort"),
				ctx.String("input-file"),
				f,
			)
		default:
			return nil
		}
	}
}

func getOptions(ctx *cli.Context) map[string]string {
	o := make(map[string]string)

	// Get endpoint url.
	if p := ctx.String("local"); len(p) != 0 {
		o["local"] = p
	}

	// Get region.
	if r := ctx.String("region"); len(r) != 0 {
		o["region"] = r
	}

	// Get profile.
	if p := ctx.String("profile"); len(p) != 0 {
		o["profile"] = p
	}

	return o
}

func newEdyClient(c client.NewClient) edy.Edy {
	return &edy.Instance{
		NewClient: c,
	}
}
