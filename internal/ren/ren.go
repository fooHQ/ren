package ren

import (
	"github.com/urfave/cli/v3"

	"github.com/foohq/ren"
	"github.com/foohq/ren/internal/ren/actions"
	"github.com/foohq/ren/internal/ren/commands/build"
	"github.com/foohq/ren/internal/ren/commands/run"
)

func New() *cli.Command {
	return &cli.Command{
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
}
