# go-code

go-code is a tool to make the development experience of writing go code with visutal studio code more pleasant.  For now it helpers me debug tests using the visual studio code debugger

## Installation

```shell
$ go install github.com/nullstyle/go/tools/go-code

# to ensure prerequisites are installed
$ go-code check
```

## Prerequisites

This tool is a thin wrapper around other tools, namely `dlv`.  It is you're responsibility to install them (although a pull request that adds automatic installation would be great).  To help you out, `go-code check` can be used to see if your local shell has all of the requirements in pace.

## Usage

```
go-code test ./somepkg
```

