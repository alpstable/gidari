[![docs](https://img.shields.io/static/v1?label=coinbase&message=reference&color=blue)](https://pkg.go.dev/github.com/alpine-hodler/web@v2.0.1-alpha/polygon)

# Polygon API Wrapper

- [Creating a Client](#creating-a-client)
- [Development](#development)
  - [Testing](#testing)

This package wraps the references defined by the [Polygon API](https://polygon.io/docs/stocks/getting-started), and can be installed using

```
go get github.com/alpine-hodler/web
```

## Creating a Client

The `polygon.Client` type is a wrapper for the Go [net/http](https://pkg.go.dev/net/http) standard library package.  An [`http.RoundTripper`](https://pkg.go.dev/net/http#RoundTripper) is required to authenticate for Polygon requests.  Currently the only method supported for Polygon authentication is via OAuth 2.0 bearer token.  See documentation examples for more information on best practices.

### Rate Limits

Per the [Polygon Knowledge Base](https://polygon.io/knowledge-base/article/what-is-the-request-limit-for-polygon-restful-apis):

> Our free tier subscriptions come with a limit of 5 API requests per minute. Paying customers are allowed unlimited API requests. While we do not have a concrete rate limit, we do monitor usage to ensure that no single user is affecting others’ quality of service. We recommend staying under 100 requests per second to avoid any throttling issues.

The HTTP package uses `"golang.org/x/time/rate"` to ensure that rate limits are honored.  This data is [lazy loaded](https://en.wikipedia.org/wiki/Lazy_loading) to cut down on memory consumption.

## Development

Notes on developing in this package.

### Testing

You will need to create an account with [Polygon](https://polygon.io/pricing).  Then populate the following data in `pkg/polygon/.test.env`:
```.env
POLYGON_API_KEY=
```

Note that `pkg/coinbase/.test.env` is an ignored file and should not be commitable to the repository.
