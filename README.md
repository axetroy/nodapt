English | [中文简体](README_zh-CN.md)

[![Build Status](https://github.com/axetroy/virtual-node-env/workflows/ci/badge.svg)](https://github.com/axetroy/virtual-node-env/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/virtual-node-env)](https://goreportcard.com/report/github.com/axetroy/virtual-node-env)
![Latest Version](https://img.shields.io/github/v/release/axetroy/virtual-node-env.svg)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/virtual-node-env.svg)

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
$ virtual-node-env node -v

# Specify the NodeJS version and run the specified command
$ virtual-node-env use 18.20.0 node -v
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

### Installation

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
# then remove the binary file or uninstall via package manager
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
