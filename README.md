# Diver

This is a `cli` tool to interact with the APIs of the Docker Enterprise Edition products enabling an end user to manage:

- Docker Universal Control Plane
- Docker Trusted Registry
- Docker Store

## Downloading

Pre-built `diver` binaries are available on the releases page [https://github.com/joeabbey/diver/releases](https://github.com/joeabbey/diver/releases)

## Building

**QUICK USAGE** From your local machine `go get github.com/joeabbey/diver` however you may need to retrieve any dependancies or tooling. 

**Alternatively (or for dev work)**
Clone the repository and build with `make build`, the `make docker` will a build a local `scratch `container that only has the binary.

Alternatively you can manually compile `diver` through the use of `go build`.

## Documentation

All documentation is being migrated from the repository readme to the `/docs` folder.


## Debugging Issues

When errors are reported turn up the `--logLevel` to 5, which enables debugging output.
