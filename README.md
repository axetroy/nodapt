[![Build Status](https://github.com/axetroy/virtual-node-env/workflows/ci/badge.svg)](https://github.com/axetroy/virtual-node-env/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/virtual-node-env)](https://goreportcard.com/report/github.com/axetroy/virtual-node-env)
![Latest Version](https://img.shields.io/github/v/release/axetroy/virtual-node-env.svg)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/virtual-node-env.svg)

## virtual-node-env

A tool similar to virtualenv, used to set a specific node version for a specified project, which is used to meet the needs of different node versions of different projects.

### Features

- [x] Cross-platform supports. including windows
- [x] Automatically download node

### Usage

```bash
$ virtual-node-env --node=16.20.0 node -v
v16.20.0

$ virtual-node-env --node=16.20.0 npm -v
8.19.4

$ virtual-node-env --node=18.20.0 node -v
v18.20.0

$ virtual-node-env --node=18.20.0 npm -v
10.5.0
```

or put it into `package.json`

```json
 "scripts": {
    "build": "virtual-node-env --node=16.20.0 yarn build"
  },
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
$ virtual-node-env --clean
# then remove the binary file
```

### License

The [Anti-996 License](LICENSE)
