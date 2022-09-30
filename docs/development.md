# Development

- [Dependencies](#dependencies)
- [Build](#build)
- [Integration Testing](#integration-testing)
  - [Network Updates](#network-updates)
  - [Credential Setup](#credential-setup)
  - [Running Integration Tests](#running-integration-tests)

## Dependencies
To develop locally you need to install the following dependencies:

1. Docker: https://docs.docker.com/get-docker/
2. Go: https://go.dev/doc/install
3. protobuf: http://google.github.io/proto-lens/installing-protoc.html
4. protoc-gen-go: https://developers.google.com/protocol-buffers/docs/gotutorial#compiling-your-protocol-buffers
5. godotenv (test only): https://github.com/joho/godotenv#installation
6. golangci-lin (test only): https://golangci-lint.run/usage/install/#local-installation
7. gofumt (test only): https://github.com/mvdan/gofumpt

## Build

To build run the default make:

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

### Credential Setup

The integration tests uses the free [Coinbse Pro API](https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getaccounts-1). To use `make tests` you must configure an environment file at `/etc/alpine-hodler/auth.env` with the test keys. It should look like this:

```.env
CBP_PASSPHRASE=<YOUR_SANDBOX_PASSPHRASE>
CBP_KEY=<YOUR_SANDBOX_KEY>
CBP_SECRET=<YOUR_SANDBOX_SECRET>
```

To create test keys, follow the guide [here](https://help.coinbase.com/en/pro/other-topics/api/how-do-i-create-an-api-key-for-coinbase-pro) or reach out to a repository contributor. Create READ ONLY test keys for [Coinbase Pro Sandbox](https://public.sandbox.pro.coinbase.com/). DO NOT USE LIVE COINBASE PRO CREDENTIALS.

### Running Integration Tests

To test locally first build the containers for integration tests using `make containers`. Then run `make tests`.
