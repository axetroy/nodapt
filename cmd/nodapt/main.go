package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/axetroy/nodapt/internal/command"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printHelp() {
	fmt.Println(`nodapt - A virtual node environment for node.js, node version manager for projects.

USAGE:
  nodapt [OPTIONS] <ARGS...>
  nodapt [OPTIONS] run <ARGS...>
  nodapt [OPTIONS] use <CONSTRAINT> [ARGS...]
  nodapt [OPTIONS] rm <CONSTRAINT>
  nodapt [OPTIONS] clean
  nodapt [OPTIONS] ls
  nodapt [OPTIONS] ls-remote

COMMANDS:
  <ARGS...>                   Alias for 'run <ARGS...>' but shorter
  run <ARGS...>               Automatically select node version to run commands
  use <CONSTRAINT> <ARGS...>  Use the specified version of node to run the command
  rm|remove <CONSTRAINT>      Remove the specified version of node that installed by nodapt
  clean                       Remove all the node version that installed by nodapt
  ls|list                     List all the installed node version
  ls-remote|list-remote       List all the available node version
  help <COMMAND>              Print help information for the specified command

GLOBAL OPTIONS:
  --help|-h                   Print help information
  --version|-v                Print version information

GLOBAL ENVIRONMENT VARIABLES:
  NODE_MIRROR                 The mirror of the nodejs download, defaults to: https://nodejs.org/dist/
                              Chinese users defaults to: https://registry.npmmirror.com/-/binary/node/
  NODE_ENV_DIR                The directory where the nodejs is stored, defaults to: $HOME/.nodapt
  DEBUG                       Print debug information when set DEBUG=1

EXAMPLES:
  nodapt node -v
  nodapt run node -v
  nodapt use v14.17.0 node -v

SOURCE CODE:
  https://github.com/axetroy/nodapt`)
}

func main() {
	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Print help information",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print version information",
	}

	app := &cli.Command{
		Name:    "nodapt",
		Usage:   "A virtual node environment for node.js, node version manager for projects.",
		Version: version,
		Suggest: true,
		Commands: []*cli.Command{
			{
				Name:        "run",
				Usage:       `Automatically select node version to run commands`,
				Description: `Automatically select node version to run commands`,
				Arguments:   cli.AnyArguments,
				ArgsUsage:   `<COMMANDS...>`,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return command.Run(cmd.Args().Slice())
				},
			},
			{
				Name:        "use",
				Usage:       "Use the specified version of node to run the command",
				Description: "Use the specified version of node to run the command",
				Arguments:   cli.AnyArguments,
				ArgsUsage:   `<CONSTRAINT> [COMMANDS...]`,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					args := cmd.Args().Slice()

					length := len(args)

					switch length {
					case 0:
						return command.Use(nil)
					case 1:
						constraint := args[0]

						return command.Use(&constraint)
					default:
						constraint := args[0]
						commands := args[1:]

						return command.RunWithConstraint(constraint, commands)
					}
				},
			},
			{
				Name:        "remove",
				Usage:       "Remove the specified version of node that installed by nodapt",
				Description: "Remove the specified version of node that installed by nodapt",
				Aliases:     []string{"rm"},
				Arguments:   cli.AnyArguments,
				ArgsUsage:   `<CONSTRAINT...>`,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					for _, constraint := range cmd.Args().Slice() {
						if err := command.Remove(constraint); err != nil {
							return errors.WithStack(err)
						}
					}

					return nil
				},
			},
			{
				Name:        "clean",
				Usage:       "Remove all the node version that installed by nodapt",
				Description: "Remove all the node version that installed by nodapt",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return command.Clean()
				},
			},
			{
				Name:        "list",
				Usage:       "List all the installed node version",
				Description: "List all the installed node version",
				Aliases:     []string{"ls"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return command.List()
				},
			},
			{
				Name:        "list-remote",
				Usage:       "List all the available node version",
				Description: "List all the available node version",
				Aliases:     []string{"ls-remote"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return command.ListRemote()
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return command.Run(cmd.Args().Slice())
		},
	}

	cli.VersionPrinter = func(c *cli.Command) {
		fmt.Printf("%s %s %s\n", version, commit, date)
	}

	cli.HelpPrinter = func(w io.Writer, templ string, data any) {
		printHelp()
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		if os.Getenv("DEBUG") == "1" {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			fmt.Fprintf(os.Stderr, "current commit hash %s\n", commit)
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			fmt.Fprintln(os.Stderr, "Print debug information when set DEBUG=1")
		}

		unwrapError := errors.Unwrap(err)

		if err, ok := err.(*exec.ExitError); ok {
			os.Exit(err.ExitCode())
		} else if err, ok := unwrapError.(*exec.ExitError); ok {
			os.Exit(err.ExitCode())
		} else {
			os.Exit(1)
		}
	}
}
