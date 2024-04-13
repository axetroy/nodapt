package main

import (
	"fmt"
	"os"
	"strings"

	VirtualNodeEnvironment "github.com/axetroy/virtual_node_env"
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

OPTIONS:
  --help                 Print help information
  --version              Print version information
  --node <version>       Specify the version of node to use

SOURCE CODE:
  https://github.com/axetroy/virtual-node-env`)
}

type Flag struct {
	// TODO: support parse from meta data
	Help    bool     `json:"help" long:"help" short:"h"`
	Version bool     `json:"version" long:"version" short:"v"`
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
			case arg == "--version", arg == "--version":
				f.Version = true
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

	nodeVersion := f.Node

	if strings.TrimSpace(nodeVersion) == "" {
		return fmt.Errorf("node version is required")
	}

	nodeVersion = strings.TrimPrefix(nodeVersion, "v")

	if len(f.Cmd) == 0 {
		panic("missing command")
	}

	return VirtualNodeEnvironment.Setup(&VirtualNodeEnvironment.Options{
		Version: nodeVersion,
		Cmd:     f.Cmd,
	})
}

func main() {
	if err := run(); err != nil {
		fmt.Println("error:")
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(255)
	}
}
