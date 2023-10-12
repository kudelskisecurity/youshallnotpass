package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kudelskisecurity/youshallnotpass/pkg/cmd/validatetokencmd"
	"github.com/urfave/cli/v2"
)

var (
	Version   = "development"
	BuildTime = "unknown"
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v", "V"},
		Usage:   "Print the version",
	}

	app := &cli.App{}
	app.EnableBashCompletion = true

	app.Name = "youshallnotpass"
	app.Usage = "Secure Authenticated Pipelines"
	app.UsageText = "youshallnotpass [command] [arguments...]"
	app.Copyright = fmt.Sprintf(`(c) %d Kudelski Security.`, time.Now().Year())
	app.Version = fmt.Sprintf("%s (built %s)", Version, BuildTime)
	app.Description = `youshallnotpass allows for secure authenticated pipelines`

	app.Commands = commands()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func commands() []*cli.Command {
	cmds := []*cli.Command{
		{
			Name: "version",
			Action: func(c *cli.Context) (err error) {
				cli.VersionPrinter(c)
				return nil
			},
			Usage:       "Print the version",
			Description: "Print the version",
		},
	}

	cmds = append(cmds, validatetokencmd.Commands()...)

	return cmds
}
