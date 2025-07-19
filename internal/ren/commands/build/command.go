package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/foohq/ren/internal/ren/actions"
	"github.com/foohq/ren/packager"
)

const (
	FlagOutput = "output"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:      "build",
		Usage:     "Package Risor scripts",
		ArgsUsage: "<dir>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    FlagOutput,
				Usage:   "set output file",
				Aliases: []string{"o"},
			},
		},
		Action:       action,
		OnUsageError: actions.UsageError,
	}
}

func action(ctx context.Context, c *cli.Command) error {
	return buildAction()(ctx, c)
}

func buildAction() cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		if c.Args().Len() != 1 {
			err := fmt.Errorf("command expects the following arguments: %s", c.ArgsUsage)
			_, _ = fmt.Fprintln(os.Stderr, err)
			return err
		}

		srcDir := c.Args().First()
		outputName := c.String(FlagOutput)
		if outputName == "" {
			abs, err := filepath.Abs(srcDir)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				return err
			}
			outputName = packager.NewFilename(filepath.Base(abs))
		}

		err := packager.Build(srcDir, outputName)
		if err != nil {
			err := fmt.Errorf("build error: %w", err)
			_, _ = fmt.Fprintln(os.Stderr, err)
			return err
		}

		return nil
	}
}
