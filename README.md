# Gidari

[![PkgGoDev](https://img.shields.io/badge/go.dev-docs-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/alpstable/gidari)
![Build Status](https://github.com/alpstable/gidari/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpstable/gidari)](https://goreportcard.com/report/github.com/alpstable/gidari)
[![Discord](https://img.shields.io/discord/987810353767403550)](https://discord.gg/3jGYQz74s7)

<p align="center"><img src="https://raw.githubusercontent.com/alpstable/gidari/main/etc/assets/gidari-gopher.png" width="300"></p>

Gidari is a "web-to-storage" tool for querying web APIs and persisting the resulting data onto local storage. A configuration file is used to define how this querying and storing should occur. Once you have a configuration file, you can intiate this transport using the command `gidari --config <configuration.yml>`. See [here](https://youtu.be/NgeOJ50IWhY) for a quick demonstration.

##

* [Installation](#installation)  
* [Configurations](#configurations)  
* [SQL](#sql)  
* [No SQL](#nosql)  
* [Contributing](#contributing)  
* [Releases](#releases)  
* [Resources](#resources)  

## Installation

### Command Mode

```sh
go install github.com/alpstable/gidari/cmd/gidari@latest
```

### Library

```sh
go get github.com/alpstable/gidari/pkg@latest
```

Gidari is not currently available as a stable API. Library support is a [WIP](https://github.com/alpstable/gidari/milestone/5).

## Usage

Using Gidari in command mode is a two step process:

1. Create a configuraiton file to instruct the binary on how to make the RESful HTTP requests and where to store the data
2. Run `gidari --config your_configuration.yml --verbose`

The `configuration.yml` file is used to define a set of rules for making RESTful HTTP requests and where to store the data. See [here](https://github.com/alpstable/gidari/tree/main/e2e/testdata/upsert) for example configurations.

### Configurations

| Key                              | Required | Type   | Description                                                                                                      |
|----------------------------------|----------|--------|------------------------------------------------------------------------------------------------------------------|
| url                              | T        | string | The API base URL                                                                                                 |
| authentication                   | F        | map    | Data required for authenticating the web API HTTP Requests                                                       |
| authentication.apiKey.passphrase | T        | string |                                                                                                                  |
| authentication.apiKey.Key        | T        | string |                                                                                                                  |
| authentication.apiKey.Secret     | T        | string |                                                                                                                  |
| authentication.auth2.Bearer      | T        | string |                                                                                                                  |
| connectionString                 | T        | List   | List of connection strings for communication with storage                                                        |
| rateLimit                        | T        | map    | Data required for limiting the number of requests per second, avoiding 429 errors                                |
| rateLimit.burst                  | T        | uint   | Number of requests that can be made per second                                                                   |
| rateLimit.period                 | T        | uint   | Period for the rateLimit.burst                                                                                   |
| truncate                         | F        | bool   | Truncate all tables in the databse before performing upserts                                                     |
| requests                         | F        | list   | List of requests to receive data from the web API for upserting into storage                                     |
| request.endpoint                 | T        | string | Endpoint for making the RESTful API request                                                                      |
| request.table                    | F        | string | Name of the table in the storage for upserting data. This field defaults to the last string in the endpoint path |
| request.clobColumn               | F        | string | Name of the column where data will be stored if the response data is not a valid JSON. Logs a warning if invalid JSON is received and this field is not set, and no data is saved. | 
| request.timseries                | F        | map    | Data required for upserting timeseries data, which are batched and can be resource intensive                     |
| request.timeseries.startName     | T        | string | "Name of the query/path parameter for the "start" datetime of the timeseries"                                  |
| request.timeseries.endName       | T        | string | "Name of the query/path parameter for the "end" datetime of the timeseries"                                    |
| request.timeseries.period        | T        | uint   | How often (in seconds) to build a new datetime range to batch.                                                   |
| request.timeseries.layout        | T        | string | The layout for how to build a datetime to query over (e.g. RFC3339 would be "2006-01-02T15:04:05Z07:00")     |
| request.query                    | N        | map    | A hash of data that holds the query parameters for a request                                                     |

### SQL

TODO

### NoSQL

The NoSQL use case should require no overhead from the user. Just include the connection string in the `connectionString` list of the configuration file. Currently this project only supports [MongoDB](https://www.mongodb.com/docs/drivers/go/current/).

## Contributing

Follow [this guide](docs/CONTRIBUTING.md) for information on contributing.

## Releases

See [here](docs/release_process.md) for the release process.

## Resources

- Public REST APIs from [Postman Documenter](https://documenter.getpostman.com/view/8854915/Szf7znEe)
- Go Gopher artwork by [Victoria Trum](https://www.fiverr.com/victoria_trum?source=order_page_user_message_link)
- The original Go gopher was designed by the awesome [Renee French](http://reneefrench.blogspot.com/)
