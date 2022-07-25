[![docs](https://img.shields.io/static/v1?label=coinbase&message=reference&color=blue)](https://pkg.go.dev/github.com/alpine-hodler/web@v0.1.0-alpha/pkg/twitter)

# Twitter API Wrapper

- [Creating a Client](#creating-a-client)
- [Development](#development)
  - [Testing](#testing)
  - [OAuth 2.0 User Context](#oauth-20-user-context)

This package wraps the references defined by the [Twitter API](https://developer.twitter.com/en/docs/api-reference-index), and can be installed using

```
go get github.com/alpine-hodler/web/pkg/twitter
```

## Creating a Client

The `twitter.Client` type is a wrapper for the Go [net/http](https://pkg.go.dev/net/http) standard library package.  An [`http.RoundTripper`](https://pkg.go.dev/net/http#RoundTripper) is required to authenticate for Twitter requests.  All Twitter [authentication methods](https://developer.twitter.com/en/docs/authentication/overview) are supported.  See the documentation for examples.

Note that [basic authentication](https://developer.twitter.com/en/docs/authentication/basic-auth) is restricted to enterprise accounts

> The email and password combination are the same ones that you will use to access the enterprise API console

## Development

Notes on developing in this package.

### Testing

You will need to [Sign up for access to the Twitter API](https://developer.twitter.com/en/docs/api-reference-index) and generate the APP keys.  Then populate the following data in `pkg/twitter/.simple-test.env`:
```.env
TWITTER_URL=https://api.twitter.com
TWITTER_CLIENT_ID=
TWITTER_CLIENT_SECRET=

# Basic
TWITTER_ENTERPRISE_EMAIL=
TWITTER_ENTERPRISE_PASSWORD=

# OAuth2
TWITTER_BEARER_TOKEN=

# OAuth 1
TWITTER_ACCESS_TOKEN=
TWITTER_ACCESS_SECRET=
TWITTER_CONSUMER_KEY=
TWITTER_CONSUMER_SECRET=
```

Note that `pkg/twitter/.simple-test.env` is an ignored file and should not be commitable to the repository.

#### Priming a Refresh Token

To run the integration tests, you will need to prime a refresh token to make HTTP requests that require OAuth 2.0 User Context.  Just follow the guide [here](#oauth-20-user-context) to get a refresh token.  If this is your first time, you'll need to run the below curl command to generate the `refresh_token.json` template used for testing.  Then copy and paste the output of that file into `refresh_token.json`.

```sh
curl -u $CLIENT_ID:$CLIENT_SECRET \
-X POST 'https://api.twitter.com/2/oauth2/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode "refresh_token=$REFRESH_TOKEN" \
--data-urlencode 'grant_type=refresh_token' \
--data-urlencode "client_id=$CLIENT_ID" \
| jq
```

### OAuth 2.0 User Context

Some of the requests requires OAuth 2.0 User Context, which requires a refresh token to use programatically. To get an access token through postman, follow the guide [here](developer.twitter.com/en/docs/tutorials/postman-getting-started).  In general, the following should suffice:

- Grant Type: Authorization Code (With PKCE)
- Callback URL: `https://oauth.pstmn.io/v1/callback`
- Auth URL: `https://twitter.com/i/oauth2/authorize`
- Access Token URL: `https://api.twitter.com/2/oauth2/token`
- Client ID: _see developer portal_
- Client Secret: _see developer portal_
- Code Challenge Method: `SHA-256`
- State: `state`
- Client Authentication: Send as Basic Auth header

To get a refresh token add `offline.access` to the scope.
