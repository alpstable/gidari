# Gidari

[![PkgGoDev](https://img.shields.io/badge/go.dev-docs-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/alpstable/gidari)
![Build Status](https://github.com/alpstable/gidari/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpstable/gidari)](https://goreportcard.com/report/github.com/alpstable/gidari)
[![Discord](https://img.shields.io/discord/987810353767403550)](https://discord.gg/3jGYQz74s7)

<p align="center"><img src="https://raw.githubusercontent.com/alpstable/gidari/main/etc/assets/gidari-gopher.png" width="300"></p>

Gidari is a library for batch querying data and persisting the results local storage.

## Installation

```sh
go get github.com/alpstable/gidari@latest
```

For information on using the CLI, see [here](https://github.com/alpstable/gidari-cli).

## Usage

At the moment, Gidari only supports HTTP services. There are two ways to use an HTTP service:

1. Iterate over [`http.Response`](https://pkg.go.dev/net/http#Response) data, for pre-defined [`http.Requests`](https://pkg.go.dev/net/http#Request).
2. Define a writer to concurrently "write" response data for pre-defined `http.Requests`.

See the Go Docs or [Web-to-Storage Examples](#web-to-storage-examples) for more examples.

### Web-to-Storage Examples

| Data Type | Writer                                      | Example                                 | Example Description                              |
|-----------|---------------------------------------------|-----------------------------------------|--------------------------------------------------|
| csv       | [csvpb](https://github.com/alpstable/csvpb) | [examples/csvp](examples/csvpb/main.go) | Use the HTTPService to write CSV data to stdout  |

## Contributing

Follow [this guide](docs/CONTRIBUTING.md) for information on contributing.

## Resources

- Public REST APIs from [Postman Documenter](https://documenter.getpostman.com/view/8854915/Szf7znEe)
- Go Gopher artwork by [Victoria Trum](https://www.fiverr.com/victoria_trum?source=order_page_user_message_link)
- The original Go gopher was designed by the awesome [Renee French](http://reneefrench.blogspot.com/)
