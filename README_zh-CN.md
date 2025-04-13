[English](README.md) | 中文简体

[![Build Status](https://github.com/axetroy/nodapt/workflows/ci/badge.svg)](https://github.com/axetroy/nodapt/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/nodapt)](https://goreportcard.com/report/github.com/axetroy/nodapt)
![Latest Version](https://img.shields.io/github/v/release/axetroy/nodapt.svg)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/nodapt.svg)

### 介绍

Nodapt (/noʊˈdæpt/) 是一个适配多个 Node.js 版本的命令行工具。它会根据 `package.json` 中的 Node.js 版本约束，自动选择并使用合适的 Node.js 版本运行命令。

### 背景

在开发 Node.js 项目时，我们经常需要切换 Node.js 版本。例如，项目 A 需要 `16.x.y`，而项目 B 使用 `20.x.y`。

然而，传统的全局版本管理工具（如 nvm）并不能很好地满足这些需求，主要存在以下问题：

1. **跨平台支持不足**：nvm 在 Windows 上使用不够方便。
2. **版本预安装要求**：nvm 需要提前安装指定版本，CI/CD 环境中使用不够友好。
3. **Monorepo 支持不足**：在 Monorepo 中，不同子项目可能需要不同的 Node.js 版本，nvm 无法很好地解决这一问题。

为了解决这些问题，开发了 Nodapt。它能够根据 `package.json` 中的 Node.js 版本约束，自动选择并安装合适的版本运行命令。

### 特性

- [x] 跨平台支持（Mac/Linux/Windows）
- [x] 自动选择并安装 Node.js 版本运行命令
- [x] 支持指定 Node.js 版本运行命令
- [x] 支持 `package.json` 中的 Node.js 版本约束
- [x] 支持 Monorepo 项目
- [x] 支持 CI/CD 环境
- [x] 兼容其他 Node.js 版本管理工具（如 nvm、n、fnm 等）
- [x] 支持 `nodapt use <version>` 命令开启新的 shell 会话

### 用法

```bash
# 自动选择 Node.js 版本运行命令
$ nodapt node -v

# 指定 Node.js 版本并运行命令
$ nodapt use ^18 node -v

# 指定版本范围并开启新的 shell 会话
$ nodapt use 20
```

### 集成到你的 Node.js 项目中

1. 在 `package.json` 中添加 Node.js 版本约束：

```diff
+  "engines": {
+    "node": "^20.x.x"
+  },
  "scripts": {
    "dev": "vite dev"
  }
```

2. 使用 `nodapt` 命令运行脚本：

```diff
- yarn dev
+ nodapt yarn dev
```

运行 `nodapt --help` 查看更多选项。

### 安装

#### 通过 [Cask](https://github.com/cask-pkg/cask.rs) 安装（Mac/Linux/Windows）

```bash
$ cask install github.com/axetroy/nodapt
$ nodapt --help
```

#### 通过 npm 安装

```bash
$ npm install @axetroy/nodapt -g
$ nodapt --help
```

### 卸载

```bash
$ nodapt clean
# 然后移除可执行文件，或者通过包管理器卸载
```

### Node.js 版本选择算法

本节解释运行 `nodapt` 时的行为以及它如何选择 Node.js 版本：

1. 检查当前目录下是否存在 `package.json` 文件。
2. 如果存在：
   1. 检查 `engines.node` 字段是否指定了版本约束：
      - 如果当前安装的版本符合约束，则直接使用。
      - 如果不符合，从远程列表中选择匹配的最新版本，安装后运行命令。
   2. 如果未指定 `engines.node`，直接运行命令。
3. 如果 `package.json` 不存在，直接运行命令。

### 类似项目

- [https://github.com/jdx/mise](https://github.com/jdx/mise)
- [https://github.com/gvcgo/version-manager](https://github.com/gvcgo/version-manager)
- [https://github.com/version-fox/vfox](https://github.com/version-fox/vfox)

### 开源许可

本项目采用 [Anti-996 License](LICENSE) 开源许可。
