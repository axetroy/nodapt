[English](README.md) | 中文简体

[![Build Status](https://github.com/axetroy/virtual-node-env/workflows/ci/badge.svg)](https://github.com/axetroy/virtual-node-env/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/virtual-node-env)](https://goreportcard.com/report/github.com/axetroy/virtual-node-env)
![Latest Version](https://img.shields.io/github/v/release/axetroy/virtual-node-env.svg)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/virtual-node-env.svg)

### 背景

在开发 NodeJS 项目时，我们经常需要切换 NodeJS 版本，例如，项目 A 需要`16.x.y`，而项目 B 使用`20.x.y`。

但是全局的版本管理工具类似 nvm 并不能满足，它有一下几个问题：

1. nvm 并不是跨平台的，windows 上使用起来并不是很方便。
2. nvm 需要提前安装好指定的版本才能切换，这对于 CI/CD 环境并不是很友好。
3. 在 Monorepo 中可能会存在 package A 需要`16.x.y`，而 package B 需要`20.x.y`，这时候 nvm 并不能很好的解决这个问题。

所以我开发了这个工具，用于解决这个问题。

它会根据 `packages.json` 中的 NodeJS 版本约束，使用合适的 NodeJS 版本运行命令。

### 用法

```bash
# 自动选择 NodeJS 版本运行命令
$ virtual-node-env node -v

# 指定 NodeJS 版本并运行指定命令
$ virtual-node-env use ^18 node -v
```

### 集成到你的 NodeJS 项目中

1. 在 `package.json` 中添加 NodeJS 版本约束。

```diff
+  "engines": {
+    "node": "^20.x.x"
+  },
  "scripts": {
    "dev": "vite dev"
  }
```

2. 使用 `virtual-node-env` 命令运行脚本。

```diff
- yarn dev
+ virtual-node-env yarn dev
```

运行 `--help` 查看更多选项。

```
$ virtual-node-env --help
virtual-node-env - A virtual node environment for node.js, node version manager for projects.

USAGE:
  virtual-node-env [OPTIONS] <ARGS...>
  virtual-node-env [OPTIONS] run <ARGS...>
  virtual-node-env [OPTIONS] use <CONSTRAINT> <ARGS...>
  virtual-node-env [OPTIONS] rm <CONSTRAINT> <ARGS...>
  virtual-node-env [OPTIONS] clean
  virtual-node-env [OPTIONS] ls
  virtual-node-env [OPTIONS] ls-remote

COMMANDS:
  run <ARGS...>               Automatically select node version to run commands
  use <CONSTRAINT> <ARGS...>  Use the specified version of node to run the command
  rm|remove <CONSTRAINT>      Remove the specified version of node that installed by virtual-node-env
  clean                       Remove all the node version that installed by virtual-node-env
  ls|list                     List all the installed node version
  ls-remote|list-remote       List all the available node version
  <ARGS...>                   Alias for 'run <ARGS...>' but shorter

OPTIONS:
  --help|-h                   Print help information
  --version|-v                Print version information

ENVIRONMENT VARIABLES:
  NODE_MIRROR                 The mirror of the nodejs download, defaults to: https://nodejs.org/dist/
                              Chinese users defaults to: https://registry.npmmirror.com/-/binary/node/
  NODE_ENV_DIR                The directory where the nodejs is stored, defaults to: $HOME/.virtual-node-env
  DEBUG                       Print debug information when set DEBUG=1

EXAMPLES:
  virtual-node-env node -v
  virtual-node-env run node -v
  virtual-node-env use v14.17.0 node -v

SOURCE CODE:
  https://github.com/axetroy/virtual-node-env
```

### 安装

1. 通过 [Cask](https://github.com/cask-pkg/cask.rs) 安装 (Mac/Linux/Windows)

```bash
$ cask install github.com/axetroy/virtual-node-env
$ virtual-node-env --help
```

2. 通过 npm 安装

```sh
$ npm install @axetroy/virtual-node-env -g
$ virtual-node-env --help
# 或者使用别名 vnode
$ vnode --help
```

### 卸载

```bash
$ virtual-node-env clean
# 然后移除可执行文件或者通过包管理器卸载
```

### NodeJS 版本选择算法

本节解释运行 `virtual-node-env` 时发生的情况以及它如何选择节点版本。

1. 检查 `package.json` 是否存在。
2. 如果 `package.json` 存在:
   1. 如果指定了 `engines.node` 字段，则使用指示的版本。
      1. 如果当前安装的版本与指定的版本匹配，则直接运行命令。
      2. 否则，从远程列表中选择匹配的最新版本，然后安装并运行命令。
   2. 否则, 直接运行命令。
3. 否则, 直接运行命令。

### 类似项目

[https://github.com/jdx/mise](https://github.com/jdx/mise)

[https://github.com/gvcgo/version-manager](https://github.com/gvcgo/version-manager)

[https://github.com/version-fox/vfox](https://github.com/version-fox/vfox)

### 开源许可

The [Anti-996 License](LICENSE)
