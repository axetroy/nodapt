[![Build Status](https://github.com/axetroy/virtual-node-env/workflows/ci/badge.svg)](https://github.com/axetroy/virtual-node-env/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/virtual-node-env)](https://goreportcard.com/report/github.com/axetroy/virtual-node-env)
![Latest Version](https://img.shields.io/github/v/release/axetroy/virtual-node-env.svg)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/virtual-node-env.svg)

### Intention

This project is used to switch node versions, but it does not want global node version management tools such as [nvm](https://github.com/nvm-sh/nvm).

It is not a global switch, but follows the project. For example, project A requires `16.x.y`, while project B uses `20.x.y`.

This tool allows you not to pay attention to which version of node you should use in the project.

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

2. Run the script with `virtual-node-env` command.

```diff
- yarn dev
+ virtual-node-env yarn dev
```

### Install

1. Install via [Cask](https://github.com/cask-pkg/cask.rs) (Mac/Linux/Windows)

```bash
$ cask install github.com/axetroy/virtual-node-env
$ virtual-node-env --help
```

2. Install via npm

```sh
$ npm install @axetroy/virtual-node-env -g
$ virtual-node-env --help
# or use the alias
$ vnode --help
```

### Uninstall

```bash
$ virtual-node-env clean
# then remove the binary file or uninstall vir package manager
```

### NodeJS version selection algorithm

This section explains what happens when you run `virtual-node-env` and how it selects the node version.

1. Check for the presence of `package.json`.
2. If `package.json` exists:
   1. If the `engines.node` field is specified, use the indicated version.
   2. Otherwise, run the command directly.
3. Otherwise, run the command directly.

### Similar Projects

[https://github.com/jdx/mise](https://github.com/jdx/mise)

[https://github.com/gvcgo/version-manager](https://github.com/gvcgo/version-manager)

[https://github.com/version-fox/vfox](https://github.com/version-fox/vfox)

### License

The [Anti-996 License](LICENSE)
