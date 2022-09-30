# Gidari

[![PkgGoDev](https://img.shields.io/badge/go.dev-docs-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/alpine-hodler/gidari)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpine-hodler/gidari)](https://goreportcard.com/report/github.com/alpine-hodler/gidari)

<p align="center"><img src="https://raw.githubusercontent.com/alpine-hodler/gidari/main/etc/assets/gidari-gopher.png" width="300"></p>

Gidari is a "web-to-storage" tool for querying web APIs and persisting the resulting data onto local storage. A configuraiton file is used to define how this querying and storing should occur. Once you have a configuration file, you can intiate this transport using the command `gidari --config <configuration.yml>`.

## Installation

TODO

## Usage

TODO

### Configuration

| Key                 | Required | Type    | Description                                                                                                                   |
|---------------------|----------|---------|-------------------------------------------------------------------------------------------------------------------------------|
| `url`               | Y        | string  | The URL for the RESTful API for making requests                                                                               |
| `authentication`    | Y        | map     | Data required for authenticating the web API requests                                                                         |
| `connectionStrings` | Y        | list    | List of connection strings for communicating with local/remote storage                                                        |
| `rateLimit`         | Y        | map     | Data required for limiting the number of requests per second, avoiding 429 errors                                             |
| `rateLimit.burst`   | Y        | int     | Number of requests that can be made per second                                                                                |
| `rateLimit.period`  | Y        | int     | Period for the `rateLimit.burst`                                                                                              |
| `truncate`          | N        | boolean | Truncate all tables in the database before performing request upserts                                                         |
| `requests`          | N        | list    | List of requests to receive data from the web API for upserting into local/remote storage                                     |
| `request.endpoint`  | Y        | string  | Endpoint for making the RESTful API request                                                                                   |
| `table`             | N        | string  | Name of the table in the remote/local storage for upserting data. This field defaults to the last string in the endpoint path |

### SQL

TODO

### NoSQL

TODO

## Repository

The `repository` and `proto` packages are the only packages within the application that are public-facing stable API with the purpose of communicating CRUD requests to the storage devices used in the web-to-storage transfers.

## Contributing

Follow [this guide](docs/development.md) to configure a development environment. See the [Gidari MVP](https://github.com/orgs/alpine-hodler/projects/3) project list for open issues, please only focus on issues in the "Scheduled" column. Issues labeled with "good first issue" are excellent starting points for new engineers. If you have completed an issue:

1. Create a pull request pointing to "main"
2. Add a reviewer
3. Make sure the CI passes

