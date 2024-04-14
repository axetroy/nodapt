[![Build Status](https://github.com/axetroy/virtual-node-env/workflows/ci/badge.svg)](https://github.com/axetroy/virtual-node-env/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/virtual-node-env)](https://goreportcard.com/report/github.com/axetroy/virtual-node-env)
![Latest Version](https://img.shields.io/github/v/release/axetroy/virtual-node-env.svg)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/virtual-node-env.svg)

## virtual-node-env

A tool similar to virtualenv, used to set a specific node version for a specified project, which is used to meet the needs of different node versions of different projects.

It is similar to nvm, but not a replacement for nvm. nvm manages the global node version for you, and can even switch node versions temporarily.

virtual-node-env can also switch node versions, but not globally, but following the project.

> [!NOTE]
>
> In actual project development, you may encounter situations where several
> projects depend on different node versions. For example 12.x.y / 16.x.y / 20.x.y.
>
> They are not fully compatible, so you need node version management tools, such as nvm.
> And it requires manual version switching, not automatic.
> In my project, there are some automatic CIs, so the node version needs to
> automatically follow the project.

### Features

- [x] Cross-platform supports. including windows
- [x] Automatically download node
- [x] Support `use` command to switch node version temporarily

### Usage

```bash
# Start a new shell with the specified node version
$ virtual-node-env use 16.20.0
$ node -v
16.20.0

# Specify the node version and run the specified command
$ virtual-node-env use 18.20.0 node -v
v18.20.0
```

or put it into `package.json`

```bash
npm install @axetroy/virtual-node-env --save-exact -D
```

```json
 "scripts": {
    "install-deps": "virtual-node-env use 16.20.0 npm install",
    "build": "virtual-node-env use 16.20.0 npm build"
  },
```

more usage, please run `virtual-node-env --help`

```
$ virtual-node-env --help
virtual-node-env - Virtual node environment, similar to nvm

USAGE:
virtual-node-env [OPTIONS] use <VERSION> [COMMAND]
virtual-node-env [OPTIONS] clean
virtual-node-env [OPTIONS] ls|list

COMMANDS:
  use <VERSION> [COMMAND]  Use the specified version of node to run the command
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
  DEBUG                    Print debug information

SOURCE CODE:
  https://github.com/axetroy/virtual-node-env

```

### Install

1. [Cask](https://github.com/cask-pkg/cask.rs) (Mac/Linux/Windows)

```bash
$ cask install github.com/axetroy/virtual-node-env
$ virtual-node-env --help
```

2. Install via npm

```sh
$ npm install @axetroy/virtual-node-env -g
$ virtual-node-env --help
```

### Uninstall

```bash
$ virtual-node-env clean
# then remove the binary file
```

### License

The [Anti-996 License](LICENSE)
