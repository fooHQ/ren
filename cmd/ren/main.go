package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/urfave/cli/v3"

	"github.com/foohq/ren"
	"github.com/foohq/ren/cmd/ren/actions"
	"github.com/foohq/ren/cmd/ren/commands/build"
	"github.com/foohq/ren/cmd/ren/commands/run"
)

var app = &cli.Command{
	Name:    "ren",
	Usage:   "Build, test, run Risor scripts",
	Version: ren.Version(),
	Flags:   []cli.Flag{
		// TODO
	},
	Commands: []*cli.Command{
		build.NewCommand(),
		run.NewCommand(),
	},
	CommandNotFound: actions.CommandNotFound,
	OnUsageError:    actions.UsageError,
	HideHelpCommand: true,
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Run(ctx, os.Args)
	if err != nil {
		os.Exit(1)
	}
}
