package main

import (
	"fmt"
	"os"
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

COMMANDS:
  use <VERSION> [COMMAND]  Use the specified version of node to run the command
  clean                    Clean the virtual node environment
  lslist                   List all the installed node version

OPTIONS:
  --help                   Print help information
  --version                Print version information

ENVIRONMENT VARIABLES:
	NODE_MIRROR              The mirror address of the nodejs download, default is https://nodejs.org/dist/
                           Chinese users use this mirror by default: https://registry.npmmirror.com/-/binary/node/

SOURCE CODE:
  https://github.com/axetroy/virtual-node-env`)
}

type Flag struct {
	Help    bool
	Version bool
	Cmd     []string
}

func parse() Flag {
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
			// case strings.HasPrefix(arg, "--node"):
			// 	eqIndex := strings.Index(arg, "=")

			// 	if eqIndex != -1 {
			// 		f.Node = arg[eqIndex+1:]
			// 	} else {
			// 		if length+1 >= len(args) {
			// 			panic("missing node version")
			// 		}

			// 		f.Node = args[length+1] // take value from next arg
			// 		length++
			// 	}
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

	return f
}

func run() error {
	f := parse()

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

		cmds := f.Cmd[2:]

		if len(cmds) == 0 {
			return VirtualNodeEnvironment.Use(nodeVersion)
		} else {
			return VirtualNodeEnvironment.Run(&VirtualNodeEnvironment.Options{
				Version: nodeVersion,
				Cmd:     cmds,
			})
		}

	case "ls", "list":
		return VirtualNodeEnvironment.List()
	case "clean":
		return VirtualNodeEnvironment.Clean()
	default:
		return errors.Errorf("unknown command: %s", cmd)
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(255)
	}
}
