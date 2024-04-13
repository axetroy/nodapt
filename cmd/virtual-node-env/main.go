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
virtual-node-env [OPTIONS] <COMMAND>
virtual-node-env use <VERSION>

OPTIONS:
  --help                 Print help information
  --version              Print version information
  --clean                Clean the virtual node environment
  --node <version>       Specify the version of node to use
  --node <version>       Specify the version of node to use

ENVIRONMENT VARIABLES:
	NODE_MIRROR            The mirror address of the nodejs download, default is https://nodejs.org/dist/
                         Chinese users use this mirror by default: https://registry.npmmirror.com/-/binary/node/

SOURCE CODE:
  https://github.com/axetroy/virtual-node-env`)
}

type Flag struct {
	// TODO: support parse from meta data
	Help    bool     `json:"help" long:"help" short:"h"`
	Version bool     `json:"version" long:"version" short:"v"`
	Clean   bool     `json:"clean" long:"clean"`
	Node    string   `json:"node" long:"node"`
	Cmd     []string `json:"cmd"`
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
			case arg == "--clean":
				f.Clean = true
			case strings.HasPrefix(arg, "--node"):
				eqIndex := strings.Index(arg, "=")

				if eqIndex != -1 {
					f.Node = arg[eqIndex+1:]
				} else {
					if length+1 >= len(args) {
						panic("missing node version")
					}

					f.Node = args[length+1] // take value from next arg
					length++
				}
			case arg == "--":
				commandIndex = length
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

	if f.Clean {

		homeDir, err := os.UserHomeDir()

		if err != nil {
			return errors.WithStack(err)
		}

		dir := filepath.Join(homeDir, ".virtual-node-env")
		if err := os.RemoveAll(dir); err != nil {
			return errors.WithStack(err)
		}
		os.Exit(0)
	}

	if len(f.Cmd) == 0 {
		panic("missing command")
	}

	cmd := f.Cmd[0]

	switch cmd {
	case "use":
		if len(f.Cmd) == 1 {
			return fmt.Errorf("missing node version")
		}

		nodeVersion := strings.TrimPrefix(f.Cmd[1], "v")

		return VirtualNodeEnvironment.Use(nodeVersion)
	default:
		nodeVersion := f.Node

		if strings.TrimSpace(nodeVersion) == "" {
			return fmt.Errorf("node version is required")
		}

		nodeVersion = strings.TrimPrefix(nodeVersion, "v")

		return VirtualNodeEnvironment.Setup(&VirtualNodeEnvironment.Options{
			Version: nodeVersion,
			Cmd:     f.Cmd,
		})
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Println("error:")
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(255)
	}
}
