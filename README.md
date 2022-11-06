# Gidari

[![PkgGoDev](https://img.shields.io/badge/go.dev-docs-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/alpstable/gidari)
![Build Status](https://github.com/alpstable/gidari/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpstable/gidari)](https://goreportcard.com/report/github.com/alpstable/gidari)
[![Discord](https://img.shields.io/discord/987810353767403550)](https://discord.gg/3jGYQz74s7)

<p align="center"><img src="https://raw.githubusercontent.com/alpstable/gidari/main/etc/assets/gidari-gopher.png" width="300"></p>

Gidari is a "web-to-storage" tool for querying web APIs and persisting the resulting data onto local storage.

##

* [Installation](#installation)
* [Configurations](#configurations)
* [SQL](#sql)
* [No SQL](#nosql)
* [Contributing](#contributing)
* [Releases](#releases)
* [Resources](#resources)

## Installation

```sh
go get github.com/alpstable/gidari@latest
```

For information on using the CLI, see [here](https://github.com/alpstable/gidari-cli).

## SQL

Supported SQL options 

- [Postgres](https://github.com/alpstable/gpostgres) (WIP)

## NoSQL

Supported NoSQL options 

- [CSV](https://github.com/alpstable/gcsv) (WIP)
- [MongoDB](https://github.com/alpstable/gmongo)

## Contributing

Follow [this guide](docs/CONTRIBUTING.md) for information on contributing.

## Releases

See [here](docs/release_process.md) for the release process.

## Resources

- Public REST APIs from [Postman Documenter](https://documenter.getpostman.com/view/8854915/Szf7znEe)
- Go Gopher artwork by [Victoria Trum](https://www.fiverr.com/victoria_trum?source=order_page_user_message_link)
- The original Go gopher was designed by the awesome [Renee French](http://reneefrench.blogspot.com/)
