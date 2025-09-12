# Guayavita (gvc)

Guayavita is a new programming language. This repository contains the Guayavita Compiler and Tooling (CLI), also referred to as gvc.

Status: early development. Interfaces may change.

## Features
- Single binary CLI built with Go
- Version/build info embedded at build time (Makefile provided)

## Getting started

### Prerequisites
- Go 1.24+ installed

### Install
Option 1: go install

```
go install jmpeax.com/guayavita/gvc@latest
```

Option 2: build locally

```
make build
./bin/guayavita -- version
```

You can also run without building a binary:

```
make run
```

### Usage

```
# Show version/build metadata
guayavita -- version
```

## Development

- Print build variables that will be injected by the Makefile:

```
make print-vars
```

- Build the project:

```
make build
```

- Install the CLI to your GOPATH/bin:

```
make install
```

## Project layout
- main.go: program entrypoint
- cmd/: Cobra command definitions
- internal/commons: build/version variables injected via ldflags
- internal/term: terminal output helpers

## License
This project is licensed under the BSD 3-Clause License. See the LICENSE file for details.
