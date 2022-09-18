# Development

- [Dependencies](#dependencies)
- [Build](#build)
- [Testing](#testing)

## Dependencies
To develop locally you need to install the following dependencies:

1. Docker: https://docs.docker.com/get-docker/
2. Go: https://go.dev/doc/install

## Build

To build run the default make:

```
make
```

## Testing

To test locally first build the containers for integration tests:

```
make containers
```

Then run the tests:

```
make test
```

### CI/CD

This repository uses [CircleCI](https://circleci.com/docs/executor-intro#docker) for it's CI/CD. To test the containerized integration test locally run `make containers` and then `make ctest-local`. Note that `make ctest` is not indended for local testing.
