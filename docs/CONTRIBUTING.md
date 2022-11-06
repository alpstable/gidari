# Contributing to Gidari

- [Dependencies](#dependencies)
- [Build](#build)
- [Testing](#testing)
- [Integration Testing](#integration-testing)
  - [Network Updates](#network-updates)
  - [Running Integration Tests](#running-integration-tests)
- [Testing with the CLI](#testing-with-the-cli)

Thank you for your interest in contributing to Gidari! Please make sure to fork this repository before working through issues.

## Bug Fixes and New Features

See the [Gidari MVP](https://github.com/orgs/alpstable/projects/3) project list for open issues, please only focus on issues in the "Scheduled" column. Issues labeled with "good first issue" are excellent starting points for new engineers. If you have completed an issue:

1. Fork this repository
2. Create a pull request pointing to "main"
3. Add a reviewer

All pull requests are subject to the GitHub workflow CI defined in the Actions section of the repository.

## Dependencies

To develop locally you need to install the following dependencies:

1. Docker: https://docs.docker.com/get-docker/
2. Go: https://go.dev/doc/install
3. Google protobuf compiler (protoc):

> ### Mac OS and Linux
>
> - http://google.github.io/proto-lens/installing-protoc.html

> ### Windows
>
> - Download the latest release (e.g., "protoc-21.8-win64.zip") under "Assets" https://github.com/protocolbuffers/protobuf/releases
>
> - Add to PATH by extracting to "C:\protoc-XX.X-winXX" (Be sure to replace 'X' with your appropriate release and system type)

4. protoc-gen-go: https://developers.google.com/protocol-buffers/docs/gotutorial#compiling-your-protocol-buffers

## Testing

To test run 

```
make tests
```

## Integration Testing

Gidari is a web-to-storage data transport, which means that integration tests are inevitable. This is an imperfect practice and any constructive feedback on improving the workflow is much appreciated.

### Network Updates

You will also need to sync your /etc/hosts file with the docker containers:

```
# Alpine Hodler Containers
127.0.0.1 mongo1
127.0.0.1 mongo2
127.0.0.1 mongo3
```

### Running Integration Tests

To test locally first build the containers for integration tests using `make containers`. Then run `make e2e`.

## Testing with the CLI

You may want to test the Gidari CLI with changes you make in Gidari. To do this we will use the go.mod [repalce directive](https://go.dev/ref/mod#go-mod-file-replace). You will need to fork the [github.com/alpstable/gidari-cli](https://github.com/alpstable/gidari-cli) repository and add the following to the go.mod file:

```go.mod
replace github.com/alpstable/gidari => your/local/gidari/fork
```

Then run `make` in the gidari-cli fork, this will create a binary that uses your local gidari fork.
