# Gidari

[![PkgGoDev](https://img.shields.io/badge/go.dev-docs-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/alpstable/gidari)
![Build Status](https://github.com/alpstable/gidari/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpstable/gidari)](https://goreportcard.com/report/github.com/alpstable/gidari)
[![Discord](https://img.shields.io/discord/987810353767403550)](https://discord.gg/3jGYQz74s7)

<p align="center"><img src="https://raw.githubusercontent.com/alpstable/gidari/main/etc/assets/gidari-gopher.png" width="300"></p>

Gidari is a "web-to-storage" tool for batch querying web APIs and persisting the resulting data onto local storage.

## Installation

```sh
go get github.com/alpstable/gidari@latest
```

For information on using the CLI, see [here](https://github.com/alpstable/gidari-cli).

## Usage

At the moment, Gidari only supports an HTTP service. There are two ways to use the HTTP service:

1. Iterate over [`http.Response`](https://pkg.go.dev/net/http#Response) data, for pre-defined [`http.Requests`](https://pkg.go.dev/net/http#Request).
2. Use any number of "proto.UpsertWriter" to concurrently "write" response data for pre-defined `http.Requests`.

See the Go Docs for more information on these use-cases and examples of how to apply them.

### List Writers

A "list writer" is a library for writing lists of data (such an HTTP Response) to storage. Here is a list of examples for different "list writers"

| Library                                     | Example                             | Example Description                              |
|---------------------------------------------|-------------------------------------|--------------------------------------------------|
| [csvpb](https://github.com/alpstable/csvpb) | [examples/csvpb](examples/csvpb/main.go) | Use the HTTPService to write CSV data to stdout  |

## Contributing

Follow [this guide](docs/CONTRIBUTING.md) for information on contributing.

## Resources

- Public REST APIs from [Postman Documenter](https://documenter.getpostman.com/view/8854915/Szf7znEe)
- Go Gopher artwork by [Victoria Trum](https://www.fiverr.com/victoria_trum?source=order_page_user_message_link)
- The original Go gopher was designed by the awesome [Renee French](http://reneefrench.blogspot.com/)
