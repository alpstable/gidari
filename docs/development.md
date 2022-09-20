# Development

- [Dependencies](#dependencies)
- [Build](#build)
- [Testing](#testing)

## Dependencies
To develop locally you need to install the following dependencies:

1. Docker: https://docs.docker.com/get-docker/
2. Go: https://go.dev/doc/install
3. godotenv: https://github.com/joho/godotenv#installation
4. golangci-lin: https://golangci-lint.run/usage/install/#local-installation

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

You will also need to sync your /etc/hosts file with the docker containers, you only need to do this once:

```
make hosts
```

To use `make tests` you willl need to ndd an environment configuration at `/etc/alpine-hodler/auth.env` with the test keys. It should look like this:

```.env
CBP_PASSPHRASE=
CBP_KEY=
CBP_SECRET=
POL_BEARER_TOKEN=
```

