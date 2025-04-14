package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/axetroy/nodapt/internal/command"
	"github.com/axetroy/nodapt/internal/util"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func defaultBehaviorHandler(helpFlag, versionFlag bool, commandName string) {
	if helpFlag {
		printHelp()
		return
	}

	if versionFlag {
		fmt.Printf("%s %s %s\n", version, commit, date)
		return
	}

	fmt.Printf("Unknown command: %s\n", commandName)
	printHelp()
}

func printHelp() {
	fmt.Println(`nodapt - A virtual node environment for node.js, node version manager for projects.

USAGE:
  nodapt [OPTIONS] <ARGS...>
  nodapt [OPTIONS] run <ARGS...>
  nodapt [OPTIONS] use <CONSTRAINT> [ARGS...>
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
	// Define global flags
	helpLongFlag := flag.Bool("help", false, "Print help information")
	helpShortFlag := flag.Bool("h", false, "Print help information")
	versionLongFlag := flag.Bool("version", false, "Print version information")
	versionShortFlag := flag.Bool("v", false, "Print version information")

	flag.Parse()

	showHelp := *helpLongFlag || *helpShortFlag
	showVersion := *versionLongFlag || *versionShortFlag

	util.Debug("args %v\n", os.Args)

	args := flag.Args()

	if len(args) == 0 {
		defaultBehaviorHandler(showHelp, showVersion, "")
		return
	}

	commandName := args[0]
	switch commandName {
	case "use":
		if len(args) < 2 {
			fmt.Println("Error: 'use' command requires a version constraint and optional commands.")
			return
		}
		constraint := args[1]
		commands := args[2:]
		if len(commands) == 0 {
			if err := command.Use(&constraint); err != nil {
				handleError(err)
			}
		} else {
			if err := command.RunWithConstraint(constraint, commands); err != nil {
				handleError(err)
			}
		}
	case "remove", "rm":
		if len(args) < 2 {
			fmt.Println("Error: 'remove' command requires at least one version constraint.")
			return
		}
		for _, constraint := range args[1:] {
			if err := command.Remove(constraint); err != nil {
				handleError(err)
			}
		}
	case "clean":
		if err := command.Clean(); err != nil {
			handleError(err)
		}
	case "list", "ls":
		if err := command.List(); err != nil {
			handleError(err)
		}
	case "list-remote", "ls-remote":
		if err := command.ListRemote(); err != nil {
			handleError(err)
		}
	case "run":
		if err := command.Run(args[1:]); err != nil {
			handleError(err)
		}
	default:
		if err := command.Run(args[0:]); err != nil {
			handleError(err)
		}
	}
}

func handleError(err error) {
	if os.Getenv("DEBUG") == "1" {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		fmt.Fprintln(os.Stderr, "Print debug information when set DEBUG=1")
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		os.Exit(exitErr.ExitCode())
	} else if unwrappedErr := errors.Unwrap(err); unwrappedErr != nil {
		if exitErr, ok := unwrappedErr.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
	}

	os.Exit(1)
}
