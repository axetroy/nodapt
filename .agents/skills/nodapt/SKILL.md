---
name: nodapt
description: Use this skill to execute commands when a user requests to run nodejs/npm commands, and you can specify the desired nodejs version.
---

# Nodapt

nodapt is a virtual Node.js environment designed to provide a simple and reliable version management tool for Node.js projects. It can automatically select and install the appropriate version of Node.js based on the version constraints in `package.json` and run commands.

If your `package.json` file contains Node.js version constraints, or if the user explicitly requests to run Node.js with a specific version or manage installed versions, then you should use this skill to run the command.

If the user requests to use or remove a specific Node.js version but does not provide a valid version format, ask the user to clarify the desired version before executing the command.

If the user does not have nodapt installed, you can suggest them to install it using the instructions in the [Installation](#installation) section below.

## When to Use This Skill

Use this skill when:

1. Running any Node.js or npm command in a project
2. The project's `package.json` specifies a Node.js version in `engines.node`
3. The user explicitly requests a specific Node.js version
4. Managing (installing, removing, listing) Node.js versions
5. The user mentions "nodapt" or asks about Node.js version management

## Basic Usage Patterns

### 1. Automatic Version Selection

Run commands with the Node.js version specified in `package.json`:

```bash
# Run node command with auto-detected version
nodapt node -v

# Run npm command
nodapt npm install

# Run any node-based tool
nodapt npx create-react-app my-app
```

### 2. Specify Node.js Version

Use a specific Node.js version for a command:

```bash
# Use exact version
nodapt use v14.17.0 node app.js

# Use major version (latest 22.x)
nodapt use 22 node script.js

# Use semantic versioning
nodapt use ^16.14.0 npm test
```

### 3. Version Management Commands

List installed versions:

```bash
nodapt ls
# Example output:
# v14.17.0
# v16.20.2
# v18.18.0
```

List available remote versions:

```bash
nodapt ls-remote
```

Remove a specific version:

```bash
nodapt rm v14.17.0
```

Clean all installed versions:

```bash
nodapt clean
```

## Common Scenarios

### Scenario 1: Running a Project with Version Constraints

When a user asks to run a project, first check if `package.json` exists and has Node.js version constraints:

```bash
# Example package.json
{
  "engines": {
    "node": ">=16.0.0 <19.0.0"
  }
}

# Run the project
nodapt npm start
```

### Scenario 2: Testing Across Node.js Versions

When a user needs to test code with different Node.js versions:

```bash
# Test with Node.js 14
nodapt use 14 npm test

# Test with Node.js 16
nodapt use 16 npm test

# Test with Node.js 18
nodapt use 18 npm test
```

### Scenario 3: Global Tool Installation

Install tools without affecting the system Node.js:

```bash
# Install a global tool
nodapt npm install -g yarn

# Use the installed tool
nodapt yarn build
```

## Error Handling

### Version Not Found

If the specified version is not available:

```bash
$ nodapt use v14.17.0 node -v
Error: Version v14.17.0 not found

# Check available versions
nodapt ls-remote | grep 14.17
```

### Invalid Version Format

If the user provides an invalid version format, ask for clarification:

```bash
# User: "use version 14"
# Assistant: "Would you like to use Node.js 14.x (latest minor version) or specifically v14.0.0?"

# User: "remove the old version"
# Assistant: "Please specify the exact version to remove, e.g., 'nodapt rm v14.17.0'"
```

### No Compatible Version Found

When auto-detection fails:

```bash
$ nodapt node -v
Error: No compatible Node.js version found for engines.node: ">=99.0.0"

# Suggested action: Ask user to modify package.json or specify a version manually
```

## Installation Instructions

### Method 1: Using Cask (Recommended)

```bash
$ cask install github.com/axetroy/nodapt
```

### Method 2: Using npm

```bash
$ npm install -g nodapt
```

### Method 3: Manual Installation

Check the [source repository](https://github.com/axetroy/nodapt) for binary releases.

### Verify Installation

```bash
nodapt --version
```

## Configuration

### Environment Variables

Set custom Node.js mirror (useful in China):

```bash
export NODE_MIRROR="https://registry.npmmirror.com/-/binary/node/"
```

Set custom installation directory:

```bash
export NODE_ENV_DIR="/custom/path/.nodapt"
```

Enable debug output:

```bash
export DEBUG=1
nodapt node -v
```

## Best Practices

1. **Always check for package.json**: Before running commands with nodapt, check if the project has a `package.json` with Node.js version constraints.

2. **Use version ranges wisely**: When specifying versions, use semantic versioning ranges (e.g., `^16.14.0`, `>=14.0.0 <17.0.0`).

3. **Prefer automatic selection**: Let nodapt auto-detect the version from `package.json` when possible.

4. **Clean up old versions**: Periodically run `nodapt clean` to remove unused Node.js versions and free disk space.

5. **Specify exact versions in CI/CD**: In automated environments, always specify exact versions to ensure reproducibility.

## Troubleshooting

### Issue: Command not found

```bash
# Ensure nodapt is installed
which nodapt

# If not found, reinstall
npm install -g nodapt
```

### Issue: Permission denied

```bash
# On Unix systems, check permissions
ls -la $(which nodapt)

# Or install without sudo using npm
npm config set prefix ~/.npm-global
export PATH=~/.npm-global/bin:$PATH
npm install -g nodapt
```

### Issue: Slow downloads

```bash
# Use Chinese mirror for faster downloads
export NODE_MIRROR="https://registry.npmmirror.com/-/binary/node/"
nodapt node -v
```

## Examples in Conversation

**User:** "Run my Node.js app"

```bash
# Assistant should check for package.json first
nodapt node app.js
```

**User:** "Test with Node.js 18"

```bash
nodapt use 18 npm test
```

**User:** "Install dependencies"

```bash
nodapt npm install
```

**User:** "Which Node.js versions do I have installed?"

```bash
nodapt ls
```

**User:** "Remove Node.js 14"

```bash
nodapt rm v14.17.0
```

## Help Command

For complete help information:

```bash
nodapt --help
```

This displays all available commands, options, and examples.

## Source Code

Find the source code and report issues at:
https://github.com/axetroy/nodapt
