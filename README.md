[English](README.md) | 中文简体

[![Build Status](https://github.com/axetroy/nodapt/workflows/ci/badge.svg)](https://github.com/axetroy/nodapt/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/nodapt)](https://goreportcard.com/report/github.com/axetroy/nodapt)
![Latest Version](https://img.shields.io/github/v/release/axetroy/nodapt.svg)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/nodapt.svg)

### Introduction

Nodapt (/noʊˈdæpt/) is a command-line tool designed to work with multiple Node.js versions. It automatically selects and uses the appropriate Node.js version to run commands based on the version constraints specified in the `package.json` file.

### Background

When developing Node.js projects, it is common to switch between different Node.js versions. For example, Project A might require `16.x.y`, while Project B uses `20.x.y`.

However, traditional global version management tools (e.g., nvm) often fall short in meeting these needs due to the following issues:

1. **Limited cross-platform support**: nvm is not very convenient to use on Windows.
2. **Pre-installation requirements**: nvm requires pre-installing specific versions, which is not ideal for CI/CD environments.
3. **Lack of Monorepo support**: In Monorepo setups, different subprojects may require different Node.js versions, which nvm cannot handle effectively.

To address these challenges, Nodapt was developed. It automatically selects and installs the appropriate Node.js version to run commands based on the version constraints in `package.json`.

### Features

- [x] Cross-platform support (Mac/Linux/Windows)
- [x] Automatically select and install the appropriate Node.js version to run commands
- [x] Support for running commands with a specified Node.js version
- [x] Support for Node.js version constraints in `package.json`
- [x] Monorepo project support
- [x] CI/CD environment support
- [x] Compatibility with other Node.js version managers (e.g., nvm, n, fnm)
- [x] Support for opening a new shell session with the `nodapt use <version>` command

### Usage

```bash
# Automatically select the appropriate Node.js version to run a command
$ nodapt node -v

# Run a command with a specified Node.js version
$ nodapt use ^18 node -v

# Specify a version range and open a new shell session
$ nodapt use 20
```

### Integrating with Your Node.js Project

1. Add Node.js version constraints to your `package.json` file:

```diff
+  "engines": {
+    "node": "^20.x.x"
+  },
  "scripts": {
    "dev": "vite dev"
  }
```

2. Use the `nodapt` command to run scripts:

```diff
- yarn dev
+ nodapt yarn dev
```

Run `nodapt --help` to see more options.

### Installation

#### Install via [Cask](https://github.com/cask-pkg/cask.rs) (Mac/Linux/Windows)

```bash
$ cask install github.com/axetroy/nodapt
$ nodapt --help
```

#### Install via npm

```bash
$ npm install @axetroy/nodapt -g
$ nodapt --help
```

### Uninstallation

```bash
$ nodapt clean
# Then remove the executable file or uninstall it via your package manager
```

### Node.js Version Selection Algorithm

This section explains how `nodapt` behaves and selects the appropriate Node.js version when executed:

1. Check if a `package.json` file exists in the current directory.
2. If it exists:
   1. Check if the `engines.node` field specifies a version constraint:
      - If the currently installed version satisfies the constraint, use it directly.
      - If not, select the latest matching version from the remote list, install it, and then run the command.
   2. If `engines.node` is not specified, run the command directly.
3. If `package.json` does not exist, run the command directly.

### Similar Projects

- [https://github.com/jdx/mise](https://github.com/jdx/mise)
- [https://github.com/gvcgo/version-manager](https://github.com/gvcgo/version-manager)
- [https://github.com/version-fox/vfox](https://github.com/version-fox/vfox)

### License

This project is licensed under the [Anti-996 License](LICENSE).
