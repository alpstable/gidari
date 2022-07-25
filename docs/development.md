# Development

- [Dependencies](#dependencies)
	- [Bazel](#bazel)
	- [Docker](#docker)
	- [Go](#go)
- [Build](#build)
- [Packages](#packages)
- [Adding an API](#adding-an-api)

## Dependencies
To develop locally you need to install three dependencies:

### Bazel

Bazel is an open-source build and test tool similar to Make, Maven, and Gradle. It uses a human-readable, high-level build language. Bazel supports projects in multiple languages and builds outputs for multiple platforms. Bazel supports large codebases across multiple repositories, and large numbers of users.

To install, following the intructions [here](https://docs.bazel.build/versions/4.2.2/bazel-overview.html#how-do-i-use-bazel)

If you're on macOS, [you can install Bazel via Homebrew](https://docs.bazel.build/versions/4.2.2/install-os-x.html#step-2-install-bazel-via-homebrew):

```
brew install bazel
```

### Docker

https://docs.docker.com/get-docker/

### Go

https://go.dev/doc/install


## Build

To build simply run

```
make
```

## Packages

- [Coinbase Cloud API (Coinbase Pro)](https://github.com/alpine-hodler/driver/blob/main/web/coinbasepro#development/README.md)
- [Polygon API](https://github.com/alpine-hodler/driver/blob/main/web/polygon/README.md#development)
- [Twitter API](https://github.com/alpine-hodler/driver/blob/main/web/twitter#development/README.md)

## Adding a Web API

- Create a new file in `scripts/meta/schema` named `your_example`.
- Create the go package under `pkg` using [Go best practices](https://go.dev/blog/package-names#package-names).  In our example it should probably be `yourexample`.
- Update `meta.docker-compose.yaml` to include `- ./web/coinbasepro:/usr/src/yourexample` under `generate.volumes` and `test-generate.volumes`.
- Run `make build-meta`
- Then start adding your schemes to `scripts/meta/schema/your_example`.  To build the schemas run `make generate`.
