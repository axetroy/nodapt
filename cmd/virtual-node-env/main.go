package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	VirtualNodeEnvironment "github.com/axetroy/virtual_node_env"
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
		pwd, err := os.Getwd()

		if err != nil {
			return nil, errors.WithStack(err)
		}

		defaultConfigFile := filepath.Join(pwd, ".virtual-node-env.json")

		if _, err := os.Stat(defaultConfigFile); err != nil {
			if !os.IsNotExist(err) {
				return nil, errors.WithStack(err)
			}
		}

		f.Config = &defaultConfigFile
	}

	return &f, nil
}

func run() error {
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
			return VirtualNodeEnvironment.Use(nodeVersion)
		} else {
			return VirtualNodeEnvironment.Run(&VirtualNodeEnvironment.Options{
				Version: nodeVersion,
				Cmd:     commands,
			})
		}

	case "ls", "list":
		return VirtualNodeEnvironment.List()
	case "ls-remote", "list-remote":
		return VirtualNodeEnvironment.ListRemote()
	case "clean":
		return VirtualNodeEnvironment.Clean()
	default:
		if f.Config != nil {
			config, err := VirtualNodeEnvironment.LoadConfig(*f.Config)
			if err != nil {
				return errors.WithStack(err)
			}

			if config.Node == "" {
				return errors.New("missing node field in the configuration file")
			}

			commands := f.Cmd

			// run command
			if len(commands) > 0 {
				return VirtualNodeEnvironment.Run(&VirtualNodeEnvironment.Options{
					Version: config.Node,
					Cmd:     commands,
				})
			}
		}
		return errors.Errorf("unknown command: %s", cmd)
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
