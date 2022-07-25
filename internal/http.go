package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// Encoders are optional data than can be encoded into a request's URL or Body.
type Encoder interface {
	EncodeBody() (io.Reader, error)
	EncodeQuery(*http.Request)
}

type Client interface {
	Self() *http.Client

	// RateLimiter will try to prevent 429 errors. Consistenly hitting 429 might result in an API ban.
	RateLimiter() *rate.Limiter
}

type httpFetchConfigs struct {
	clbk        func(uint8, http.Response) error
	client      http.Client
	req         *http.Request
	ratelimiter *rate.Limiter
	encoder     Encoder
	endpoint    endpoint
	params      map[string]string
	path        uint8
}

type endpoint interface {
	Path(map[string]string) string

	// Scope is the OAuth 2.0 scope required to make requests to this endpoint.
	Scope() string
}

// stringer is any type that can textualize it's value as a string.
type stringer interface {
	String() string
}

// HTTPWithCallback will set a callback function that takes a uint8 argument and returns an error.
func HTTPWithCallback(clbk func(uint8, http.Response) error) func(*httpFetchConfigs) {
	return func(c *httpFetchConfigs) {
		c.clbk = clbk
	}
}

// HTTPWithClient will set the http client on the fetch configurations.
func HTTPWithClient(client http.Client) func(*httpFetchConfigs) {
	return func(c *httpFetchConfigs) {
		c.client = client
	}
}

// HTTPWithEncoder will set the encoder for the http request body and query params.
func HTTPWithEncoder(enc Encoder) func(*httpFetchConfigs) {
	return func(c *httpFetchConfigs) {
		c.encoder = enc
	}
}

// HTTPWithEndpoint will set the endpoint data for the request on the configurations.
func HTTPWithEndpoint(ep endpoint) func(*httpFetchConfigs) {
	return func(c *httpFetchConfigs) {
		c.endpoint = ep
	}
}

// HTTPWithParams will set the query parameters for the request on the configurations.
func HTTPWithParams(params map[string]string) func(*httpFetchConfigs) {
	return func(c *httpFetchConfigs) {
		c.params = params
	}
}

// HTTPWithRatelimiter will set the rate limiting mechanism on the fetch configurations.
func HTTPWithRatelimiter(rl *rate.Limiter) func(*httpFetchConfigs) {
	return func(c *httpFetchConfigs) {
		c.ratelimiter = rl
	}
}

// HTTPWithRequest will set the http request on the fetch configurations.
func HTTPWithRequest(req *http.Request) func(*httpFetchConfigs) {
	return func(c *httpFetchConfigs) {
		c.req = req
	}
}

func httpQueryEncode(req *http.Request, key, val string) {
	q := req.URL.Query()
	q.Add(key, val)
	req.URL.RawQuery = q.Encode()
}

// HTTPQueryEncodeBool will convert a bool pointer into a string and then add it to the query parameters of an HTTP
// request's URL.
func HTTPQueryEncodeBool(req *http.Request, key string, val *bool) {
	if val != nil {
		httpQueryEncode(req, key, strconv.FormatBool(*val))
	}
}

func intPtrString(i *int) string {
	if i != nil {
		return strconv.Itoa(*i)
	}
	return ""
}

// HTTPQueryEncodeInt will convert an Int pointer into a string and then add it to the query parameters of an HTTP
// request's URL.
func HTTPQueryEncodeInt(req *http.Request, key string, val *int) {
	if val != nil {
		httpQueryEncode(req, key, intPtrString(val))
	}
}

func int32PtrString(val *int32) string {
	if val != nil {
		return strconv.Itoa(int(*val))
	}
	return ""
}

func uint8PtrString(val *uint8) string {
	if val != nil {
		return strconv.Itoa(int(uint8(*val)))
	}
	return ""
}

// HTTPQueryEncodeInt32 will convert an Int32 pointer into a string and then add it to the query parameters of an HTTP
// request's URL.
func HTTPQueryEncodeInt32(req *http.Request, key string, val *int32) {
	if val != nil {
		httpQueryEncode(req, key, int32PtrString(val))
	}
}

// HTTPQueryEncodeString will convert an String pointer into a string and then add it to the query parameters of an HTTP
// request's URL.
func HTTPQueryEncodeString(req *http.Request, key string, val *string) {
	if val != nil {
		httpQueryEncode(req, key, *val)
	}
}

// HTTPQueryEncodeStringer will convert a stringer type into a string and then add it to the query parameters of
// an HTTP request's URL.
func HTTPQueryEncodeStringer(req *http.Request, key string, val stringer) {
	if val != nil {
		if str := val.String(); len(str) > 0 {
			httpQueryEncode(req, key, str)
		}
	}
}

// HTTPQueryEncodeStringSlice will convert a slice of strings into a string and then add it to the query parameters of
// an HTTP request's URL.
func HTTPQueryEncodeStrings(req *http.Request, key string, val []string) {
	if val != nil {
		slice := []string{}
		for _, v := range val {
			slice = append(slice, v)
		}
		httpQueryEncode(req, key, strings.Join(slice, ", "))
	}
}

// HTTPQueryEncodeTime will convert a time.Time type into a string and then add it to the query parameters of an HTTP
// request's URL.
func HTTPQueryEncodeTime(req *http.Request, key string, val *time.Time) {
	if val != nil {
		httpQueryEncode(req, key, val.String())
	}
}

// HTTPQueryEncodeUint8 will convert auint type into a string and then add it to the query parameters of an HTTP
// request's URL.
func HTTPQueryEncodeUint8(req *http.Request, key string, val *uint8) {
	if val != nil {
		httpQueryEncode(req, key, uint8PtrString(val))
	}
}

// parseErrorMessage takes a response and a status and builder an error message to send to the server.
func parseErrorMessage(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Errorf("Status Code %v (%v): %v", resp.StatusCode, resp.Status, string(body))
}

// validateResponse is a switch condition that parses an error response
func validateResponse(res *http.Response) (err error) {
	if res == nil {
		err = fmt.Errorf("no response, check request and env file")
	} else {
		switch res.StatusCode {
		case
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusInternalServerError,
			http.StatusNotFound,
			http.StatusTooManyRequests,
			http.StatusForbidden:
			err = parseErrorMessage(res)
		}
	}
	return
}

// HTTPFetch will make an HTTP request given a http.Client and a partially formatted http.Request, it will then try to
// edecode the model.
func HTTPFetch(model interface{}, opts ...func(*httpFetchConfigs)) error {
	configs := new(httpFetchConfigs)
	for _, opt := range opts {
		opt(configs)
	}
	configs.req.URL.Path = configs.endpoint.Path(configs.params)
	if configs.encoder != nil {
		configs.encoder.EncodeQuery(configs.req)
	}
	if configs.req.Body != nil {
		configs.req.Header.Set("content-type", "application/json")
	}

	ctx := context.Background()
	err := configs.ratelimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return fmt.Errorf("error waiting on rate limiter: %v", err)
	}

	resp, err := configs.client.Do(configs.req)
	if err != nil {
		return fmt.Errorf("error making request %+v: %v", configs.req, err)
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return err
	}
	// bodyBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)
	if model != nil {
		if err := json.NewDecoder(resp.Body).Decode(&model); err != nil {
			fmt.Println(resp.Status, resp.StatusCode)
			return fmt.Errorf("error decoding response body: %v", err)
		}
	}
	return nil
}

// HTTPNewRequest will return a new request.  If the options are set, this function will encode a body if possible.
func HTTPNewRequest(method, url string, encoder Encoder) (*http.Request, error) {
	if encoder == nil {
		return http.NewRequest(method, url, nil)
	}
	buf, err := encoder.EncodeBody()
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, url, buf)
}

// HTTPBodyFragment will add a new fragment to the HTTP request body.
func HTTPBodyFragment(body map[string]interface{}, key string, val interface{}) {
	if val != nil {
		body[key] = val
	}
}
