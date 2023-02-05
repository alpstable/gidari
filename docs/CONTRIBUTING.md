# Contributing to Gidari

- [Dependencies](#dependencies)
- [Build](#build)
- [Testing](#testing)
- [Testing with the CLI](#testing-with-the-cli)

Thank you for your interest in contributing to Gidari! Please make sure to fork this repository before working through issues.

## Bug Fixes and New Features

See the [Gidari MVP](https://github.com/orgs/alpstable/projects/3) project list for open issues, please only focus on issues in the "Scheduled" column. Issues labeled with "good first issue" are excellent starting points for new engineers. If you have completed an issue:

1. Fork this repository
2. Create a pull request pointing to "main"
3. Add a reviewer

All pull requests are subject to the GitHub workflow CI defined in the Actions section of the repository.

## Dependencies

To develop locally you will need to install the following dependencies:

1. Go: https://go.dev/doc/install
2. Google protobuf compiler (protoc):

> ### Mac OS and Linux
>
> - http://google.github.io/proto-lens/installing-protoc.html

> ### Windows
>
> - Download the latest release (e.g., "protoc-21.8-win64.zip") under "Assets" https://github.com/protocolbuffers/protobuf/releases
>
> - Add to PATH by extracting to "C:\protoc-XX.X-winXX" (Be sure to replace 'X' with your appropriate release and system type)

3. protoc-gen-go: https://developers.google.com/protocol-buffers/docs/gotutorial#compiling-your-protocol-buffers
4. `gofumpt`: https://github.com/mvdan/gofumpt
5. `golangcli-lint`: https://github.com/golangci/golangci-lint#install-golangci-lint

## Testing

To test run

```
make tests
```

## Testing with the CLI

You may want to test the Gidari CLI with changes you make in Gidari. To do this we will use the go.mod [repalce directive](https://go.dev/ref/mod#go-mod-file-replace). You will need to fork the [github.com/alpstable/gidari-cli](https://github.com/alpstable/gidari-cli) repository and add the following to the go.mod file:

```go.mod
replace github.com/alpstable/gidari => your/local/gidari/fork
```

Then run `make` in the gidari-cli fork, this will create a binary that uses your local gidari fork.
