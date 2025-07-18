package run

import (
	"context"
	"fmt"
	"os"

	risoros "github.com/risor-io/risor/os"
	"github.com/urfave/cli/v3"

	"github.com/foohq/ren"
	"github.com/foohq/ren/filesystems/local"
	"github.com/foohq/ren/internal/ren/actions"
	"github.com/foohq/ren/modules"
	renos "github.com/foohq/ren/os"
)

const (
	FlagEnv = "env"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "Run Risor script from a package",
		ArgsUsage: "<pkg> [[arg] ...]",
		Flags:     []cli.Flag{
			// TODO: support env variables
		},
		Action:       action,
		OnUsageError: actions.UsageError,
	}
}

func action(ctx context.Context, c *cli.Command) error {
	localFS, err := local.NewFS()
	if err != nil {
		err := fmt.Errorf("cannot initialize local filesystem: %w", err)
		_, _ = fmt.Fprintln(os.Stderr, err)
		return err
	}

	filesystems := map[string]risoros.FS{
		"file": localFS,
	}
	return runAction(filesystems)(ctx, c)
}

func runAction(filesystems map[string]risoros.FS) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		if c.Args().Len() != 1 {
			err := fmt.Errorf("command expects the following arguments: %s", c.ArgsUsage)
			_, _ = fmt.Fprintln(os.Stderr, err)
			return err
		}

		pkg := c.Args().First()
		args := c.Args().Tail()

		ros := renos.New(
			renos.WithStdin(os.Stdin),
			renos.WithStdout(os.Stdout),
			renos.WithArgs(args),
			renos.WithFilesystems(filesystems),
		)

		err := ren.RunFile(
			ctx,
			pkg,
			ros,
			ren.WithGlobals(modules.Globals()),
		)
		if err != nil {
			err := fmt.Errorf("run error: %w", err)
			_, _ = fmt.Fprintln(os.Stderr, err)
			return err
		}

		return nil
	}
}
