English | [中文简体](README_zh-CN.md)

[![Build Status](https://github.com/axetroy/nodapt/workflows/ci/badge.svg)](https://github.com/axetroy/nodapt/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/nodapt)](https://goreportcard.com/report/github.com/axetroy/nodapt)
![Latest Version](https://img.shields.io/github/v/release/axetroy/nodapt.svg)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/nodapt.svg)

### Introduction

Nodapt (/noʊˈdæpt/) is a command-line tool that adapts to multiple NodeJS versions. It will run commands with the appropriate NodeJS version based on the NodeJS version constraints in `packages.json`.

### Background

When developing NodeJS projects, we often need to switch NodeJS versions. For example, project A requires `16.x.y`, while project B uses `20.x.y`.

However, global version management tools like nvm cannot meet the requirements. It has the following problems:

1. nvm is not cross-platform, and it is not very convenient to use on Windows.
2. nvm needs to install the specified version in advance to switch, which is not very friendly to the CI/CD environment.
3. In a Monorepo, there may be a situation where package A requires `16.x.y`, while package B requires `20.x.y`. In this case, nvm cannot solve this problem well.

So I developed this tool to solve this problem.

It will run the command with the appropriate NodeJS version according to the NodeJS version constraint in `packages.json`.

### Usage

```bash
# Automatically select the NodeJS version to run the command
$ nodapt node -v

# Specify the NodeJS version and run the specified command
$ nodapt use ^18 node -v
```

### Integrate into your NodeJS project

1. Add NodeJS version constraint in `package.json`.

```diff
+  "engines": {
+    "node": "^20.x.x"
+  },
  "scripts": {
    "dev": "vite dev"
  }
```

2. Run the script with `nodapt` command.

```diff
- yarn dev
+ nodapt yarn dev
```

Run with `--help` to see more options.

```
$ nodapt --help
nodapt - A virtual node environment for node.js, node version manager for projects.

USAGE:
  nodapt [OPTIONS] <ARGS...>
  nodapt [OPTIONS] run <ARGS...>
  nodapt [OPTIONS] use <CONSTRAINT> <ARGS...>
  nodapt [OPTIONS] rm <CONSTRAINT> <ARGS...>
  nodapt [OPTIONS] clean
  nodapt [OPTIONS] ls
  nodapt [OPTIONS] ls-remote

COMMANDS:
  run <ARGS...>               Automatically select node version to run commands
  use <CONSTRAINT> <ARGS...>  Use the specified version of node to run the command
  rm|remove <CONSTRAINT>      Remove the specified version of node that installed by nodapt
  clean                       Remove all the node version that installed by nodapt
  ls|list                     List all the installed node version
  ls-remote|list-remote       List all the available node version
  <ARGS...>                   Alias for 'run <ARGS...>' but shorter

OPTIONS:
  --help|-h                   Print help information
  --version|-v                Print version information

ENVIRONMENT VARIABLES:
  NODE_MIRROR                 The mirror of the nodejs download, defaults to: https://nodejs.org/dist/
                              Chinese users defaults to: https://registry.npmmirror.com/-/binary/node/
  NODE_ENV_DIR                The directory where the nodejs is stored, defaults to: $HOME/.nodapt
  DEBUG                       Print debug information when set DEBUG=1

EXAMPLES:
  nodapt node -v
  nodapt run node -v
  nodapt use v14.17.0 node -v

SOURCE CODE:
  https://github.com/axetroy/nodapt
```

### Installation

1. Install via [Cask](https://github.com/cask-pkg/cask.rs) (Mac/Linux/Windows)

```bash
$ cask install github.com/axetroy/nodapt
$ nodapt --help
```

2. Install via npm

```sh
$ npm install @axetroy/nodapt -g
$ nodapt --help
```

### Uninstall

```bash
$ nodapt clean
# then remove the binary file or uninstall via package manager
```

### NodeJS version selection algorithm

This section explains what happens when you run `nodapt` and how it selects the node version.

1. Check for the presence of `package.json`.
2. If `package.json` exists:
   1. If the `engines.node` field is specified, use the indicated version.
      1. If the currently installed version matches `engines.node`, the command is run using the currently installed version.
      2. Otherwise, select the latest matching version from the remote list, install it, and run the command.
   2. Otherwise, run the command directly.
3. Otherwise, run the command directly.

### Similar Projects

[https://github.com/jdx/mise](https://github.com/jdx/mise)

[https://github.com/gvcgo/version-manager](https://github.com/gvcgo/version-manager)

[https://github.com/version-fox/vfox](https://github.com/version-fox/vfox)

### License

The [Anti-996 License](LICENSE)
