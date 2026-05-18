---
name: nodapt
description: Use this skill to execute commands when a user requests to run nodejs/npm commands, and you can specify the desired nodejs version.
---

# Nodapt

nodapt is a virtual Node.js environment designed to provide a simple and reliable version management tool for Node.js projects. It can automatically select and install the appropriate version of Node.js based on the version constraints in `package.json` and run commands.

If your `package.json` file contains Node.js version constraints, then you should use this skill to run the command.

## Example

Typically, you would add `nodapt` before the command you want to run, for example...

```bash
# Use nodapt to run the `node -v` command, automatically selecting the appropriate nodejs version to run the command
nodapt node -v
# Use nodapt to run the `node -v` command, specifying the use of v14.17.0 version of nodejs to run the command
nodapt use v14.17.0 node -v
# Use nodapt to run the `node -v` command, specifying the use of the latest version of Node.js 22.x
nodapt use 22 node -v
```

### nodapt help information

```
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
