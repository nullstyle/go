# go-electron

gopm bundles npm packages for gopherjs TODO

## Installation

```shell
$ go install github.com/nullstyle/go/tools/gopm

# to ensure prerequisites are installed
$ gopm check
```

## Prerequisites

This tool is a thin wrapper around `browserify`.  It is you're responsibility to install it (although a pull request that adds automatic installation would be great).  To help you out, `gopm check` can be used to see if your local shell has all of the requirements in pace.

## Usage

```
gopm bundles npm packages for gopherjs

Usage:
  gopm [command]

Available Commands:
  build       Build creates the gopm bundle for the app at PATH
  check       check the local execution environment

Flags:
      --config string   config file (default is $HOME/.gopm.yaml)

Use "gopm [command] --help" for more information about a command.
```
