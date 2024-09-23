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
  --config                 Specify the configuration file. Detected .virtual-node-env.json automatically if not specified.

ENVIRONMENT VARIABLES:
  NODE_MIRROR              The mirror of the nodejs download, defaults to: https://nodejs.org/dist/
                           Chinese users defaults to: https://registry.npmmirror.com/-/binary/node/
  NODE_ENV_DIR             The directory where the nodejs is stored, defaults to: $HOME/.virtual-node-env
  DEBUG                    Print debug information when set DEBUG=1

Configuration:
	The configuration file is a JSON file that contains the node version.
	By default, if there is no configuration in the current directory, it will automatically search for the configuration file upwards.

SOURCE CODE:
  https://github.com/axetroy/virtual-node-env`)
}

type Flag struct {
	Help    bool
	Version bool
	Config  *string
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
			case strings.HasPrefix(arg, "--config"):
				eqIndex := strings.Index(arg, "=")

				if eqIndex != -1 {
					val := arg[eqIndex+1:]

					f.Config = &val
				} else {
					if length+1 >= len(args) {
						panic("missing value for --config flag")
					}

					val := args[length+1] // take value from next arg
					f.Config = &val
					length++
				}
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

	// Detect the configuration file if not specified
	if f.Config == nil {
		cwd, err := os.Getwd()

		if err != nil {
			return nil, errors.WithStack(err)
		}

		configFilePath := util.LoopUpFile(cwd, ".virtual-node-env.json")

		if configFilePath != nil {
			util.Debug("Use configuration file: %s\n", *configFilePath)
		}

		f.Config = configFilePath
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

		var configPath *string

		if f.Config == nil {
			configPathInCwd := filepath.Join(cwd, ".virtual-node-env.json")

			if util.IsFileExist(configPathInCwd) {
				configPath = &configPathInCwd
			}
		} else {
			configPath = f.Config
		}

		// If the configuration file is found, then use the configuration file to run the command
		if configPath != nil {
			util.Debug("Use configuration file: %s\n", *configPath)

			return cli.RunWithConfig(*configPath, f.Cmd)
		}

		// If the package.json file is found, then use the node version in the package.json to run the command
		if util.IsFileExist(filepath.Join(cwd, "package.json")) {
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

		// Loop up the configuration file in the parent directory
		parentConfig := util.LoopUpFile(cwd, ".virtual-node-env.json")

		if parentConfig != nil {
			util.Debug("Use configuration file: %s\n", *parentConfig)
			return cli.RunWithConfig(*parentConfig, f.Cmd)
		} else {
			return errors.WithStack(errors.New("missing configuration file"))
		}
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
