# Contributing to Gidari

- [Dependencies](#dependencies)
- [Build](#build)
- [Integration Testing](#integration-testing)
  - [Network Updates](#network-updates)
  - [Running Integration Tests](#running-integration-tests)
- [Socials](#socials)

Thank you for your intest in contributing to Gidari! Please make sure to fork this repository before working through issues.

## Bug Fixes and New Features

See the [Gidari MVP](https://github.com/orgs/alpine-hodler/projects/3) project list for open issues, please only focus on issues in the "Scheduled" column. Issues labeled with "good first issue" are excellent starting points for new engineers. If you have completed an issue:

1. Fork this repository
2. Create a pull request pointing to "main"
3. Add a reviewer

All pull requests are subject to the GitHub workflow CI defined in the Actions section of the repository.

## Dependencies

To develop locally you need to install the following dependencies:

1. Docker: https://docs.docker.com/get-docker/
2. Go: https://go.dev/doc/install
3. protobuf: http://google.github.io/proto-lens/installing-protoc.html
4. protoc-gen-go: https://developers.google.com/protocol-buffers/docs/gotutorial#compiling-your-protocol-buffers
5. golangci-lint (test only): https://golangci-lint.run/usage/install/#local-installation
6. gofumt (test only): https://github.com/mvdan/gofumpt

## Build

To build a binary (this is not a required step):

```
make
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
127.0.0.1 postgres1
```

### Running Integration Tests

To test locally first build the containers for integration tests using `make containers`. Then run `make tests`.
