# go-electron

go-electron is a tool to write electron application in go.

## Installation

```shell
$ go install github.com/nullstyle/go/tools/go-electron

# to ensure prerequisites are installed
$ go-electron check
```

## Prerequisites

This tool is a thin wrapper around several other tools, such as `electron` and `electron-packager`.  It is you're responsibility to install them (although a pull request that adds automatic installation would be great).  To help you out, `go-electron check` can be used to see if your local shell has all of the requirements in pace.

## Usage

```
go-electron makes it easy to develop, test, and package electron applications built from go code and gopherjs.

Usage:
  go-electron [command]

Available Commands:
  build       Build the go-electron app at PATH
  check       check the local execution environment
  run         run the go-electron app at PATH

Flags:
      --config string   config file (default is $HOME/.go-electron.yaml)
  -h, --help            help for go-electron

Use "go-electron [command] --help" for more information about a command.
```

## Writing electron applications in go

go-electron aims to provide a minimal, sensible layer overtop the terrible, normal web development experience of nowadays.  You decide whether I've succeeded or not.  It's certainly not there yet.

Building an application with this tool requires that you write two go packages, one for electron main process (called the root package) and one for the browser process.  The build process assumes that the browser package is nested underneath the "main" package.  See [https://github.com/nullstyle/go/tree/master/apps/electron-example] for an example.

### The root package

The root package of an application app serves two purposes: it contains the code that executes within the main process of the electron app, but it also provides metadata to the build process.  Strangely enough, it gets run both using the native go runtime and using gopherjs.

The metadata needs of the package should be handled by using the [electron](https://github.com/nullstyle/go/tree/master/electron) package, instantiating an `App` struct and calling `Start()` on the instance.