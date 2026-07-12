package run

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/foohq/ren"
	"github.com/foohq/ren/builtins"
	"github.com/foohq/ren/cmd/ren/actions"
	"github.com/foohq/ren/modules"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:         "run",
		Usage:        "Run Risor script from a package",
		ArgsUsage:    "<pkg> [[arg] ...]",
		Flags:        []cli.Flag{},
		Action:       action,
		OnUsageError: actions.UsageError,
	}
}

func action(ctx context.Context, c *cli.Command) error {
	return runAction()(ctx, c)
}

func runAction() cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		if c.Args().Len() == 0 {
			err := fmt.Errorf("command expects the following arguments: %s", c.ArgsUsage)
			_, _ = fmt.Fprintln(os.Stderr, err)
			return err
		}

		pkg := c.Args().First()
		args := c.Args().Tail()

		opts := []ren.Option{
			ren.WithArgs(args),
		}
		for _, builtin := range builtins.Builtins() {
			opts = append(opts, ren.WithBuiltin(builtin))
		}

		for _, module := range modules.Modules() {
			opts = append(opts, ren.WithModule(module))
		}

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		opts = append(opts, ren.WithExitHandler(func(c int) {
			cancel()
		}))

		err := ren.RunFile(
			ctx,
			pkg,
			opts...,
		)
		if err != nil {
			err := fmt.Errorf("run error: %w", err)
			_, _ = fmt.Fprintln(os.Stderr, err)
			return err
		}

		return nil
	}
}
