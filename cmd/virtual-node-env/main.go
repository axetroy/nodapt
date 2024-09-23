package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/axetroy/virtual_node_env/internal/cli"
	"github.com/axetroy/virtual_node_env/internal/node"
	"github.com/axetroy/virtual_node_env/internal/util"
	"github.com/pkg/errors"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printHelp() {
	println(`virtual-node-env - Virtual node environment, similar to nvm

USAGE:
virtual-node-env [OPTIONS] [COMMAND]
virtual-node-env [OPTIONS] use <VERSION> [COMMAND]
virtual-node-env [OPTIONS] clean
virtual-node-env [OPTIONS] ls|list
virtual-node-env [OPTIONS] ls-remote|list-remote

COMMANDS:
  use <VERSION> [COMMAND]  Use the specified version of node to run the command
  rm|remove <VERSION>      Remove the specified version of node that installed by virtual-node-env
  clean                    Clean the virtual node environment
  ls|list                  List all the installed node version
  ls-remote|list-remote    List all the available node version

OPTIONS:
  --help                   Print help information
  --version                Print version information

ENVIRONMENT VARIABLES:
  NODE_MIRROR              The mirror of the nodejs download, defaults to: https://nodejs.org/dist/
                           Chinese users defaults to: https://registry.npmmirror.com/-/binary/node/
  NODE_ENV_DIR             The directory where the nodejs is stored, defaults to: $HOME/.virtual-node-env
  DEBUG                    Print debug information when set DEBUG=1

SOURCE CODE:
  https://github.com/axetroy/virtual-node-env`)
}

type Flag struct {
	Help    bool
	Version bool
	Cmd     []string
}

func parse() (*Flag, error) {
	args := os.Args[1:]

	f := Flag{}

	length := 0
	commandIndex := -1

	for length < len(args) {
		arg := args[length]

		if commandIndex >= 0 && length > commandIndex {
			f.Cmd = append(f.Cmd, arg)
			length++
			continue
		}

		if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
			switch true {
			case arg == "--help", arg == "-h":
				f.Help = true
			case arg == "--version", arg == "-v":
				f.Version = true
			case arg == "--":
				if commandIndex == -1 {
					commandIndex = length
				} else {
					f.Cmd = append(f.Cmd, arg)
				}
			}
		} else {
			commandIndex = length
			f.Cmd = append(f.Cmd, arg)
		}

		length++
	}

	return &f, nil
}

func run() error {
	cwd, err := os.Getwd()

	if err != nil {
		return errors.WithStack(err)
	}

	f, err := parse()
	if err != nil {
		return errors.WithStack(err)
	}

	if f.Help {
		printHelp()
		os.Exit(0)
	}

	if f.Version {
		println(fmt.Sprintf("%s %s %s", version, commit, date))
		os.Exit(0)
	}

	if len(f.Cmd) == 0 {
		return errors.New("missing command")
	}

	cmd := f.Cmd[0]

	switch cmd {
	case "use":
		if len(f.Cmd) == 1 {
			return fmt.Errorf("missing node version")
		}

		nodeVersion := strings.TrimPrefix(f.Cmd[1], "v")

		commands := f.Cmd[2:]

		if len(commands) == 0 {
			return cli.Use(nodeVersion)
		} else {
			return cli.Run(&cli.RunOptions{
				Version: nodeVersion,
				Cmd:     commands,
			})
		}
	case "rm", "remove":
		if len(f.Cmd) == 1 {
			return fmt.Errorf("missing node version")
		}

		nodeVersion := strings.TrimPrefix(f.Cmd[1], "v")

		return cli.Remove(nodeVersion)

	case "ls", "list":
		return cli.List()
	case "ls-remote", "list-remote":
		return cli.ListRemote()
	case "clean":
		return cli.Clean()
	default:
		if len(f.Cmd) == 0 {
			return errors.New("missing command")
		}

		packageJSONPath := util.LoopUpFile(cwd, "package.json")

		// If the package.json file is found, then use the node version in the package.json to run the command
		if packageJSONPath != nil {
			util.Debug("Use node version from %s\n", filepath.Join(cwd, "package.json"))

			semverVersionConstraint, err := node.GetVersionFromPackageJSON(filepath.Join(cwd, "package.json"))

			if err != nil {
				return errors.WithMessage(err, "failed to get node version from package.json")
			}

			matchVersion, err := node.GetMatchVersion(*semverVersionConstraint)

			if err != nil {
				return errors.WithMessage(err, "failed to get match version")
			}

			if matchVersion == nil {
				return errors.New("no match version found")
			}

			return cli.Run(&cli.RunOptions{
				Version: *matchVersion,
				Cmd:     f.Cmd,
			})
		}

		return errors.New("can not found package.json")
	}
}

func main() {
	if err := run(); err != nil {
		if os.Getenv("DEBUG") == "1" {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			fmt.Fprintln(os.Stderr, "Print debug information when set DEBUG=1")
		}
		os.Exit(1)
	}
}
