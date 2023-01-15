package gidari

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/alpstable/gidari/internal/web"
	"golang.org/x/time/rate"
)

type flattenedRequestConfig struct {
	table      string
	clobColumn string
}

// flattenedRequest wraps HTTP requests with additional information for storage.
// The number of flattened requests may not be equal to the number of requests
// on the configuration. For example, a large timeseries request will be
// flattened into multiple requests.
type flattenedRequest struct {
	// fetchConfig is the configuration for the HTTP request. Each request
	// gets it's own connection to ensure that the web worker can process
	// concurrently without locking. Despite this, however, all of the
	// requests should share a common rate limiter to prevent overloading
	// the web API and gettig a 429 response.
	fetchConfig *web.FetchConfig

	client web.Client
	cfg    *flattenedRequestConfig
}

// newFlattenedRequest will compress the request information into a "web.FetchConfig"
// request and a "table" name for storage interaction.
func newFlattenedRequest(req *Request, client web.Client, rateLimiter *rate.Limiter) *flattenedRequest {
	fetchConfig := newFetchConfig(req.Request, client, rateLimiter)

	// If the table is not set on the request, then set it using the last path segment of the endpoint.
	//if req.Table == "" {
	//	req.Table = path.Base(req.Endpoint)
	//}

	return &flattenedRequest{
		fetchConfig: fetchConfig,
		//rspHandler:  req.HttpResponseHandler,
		//reqHandler:  req.HttpRequestHandler,
		//table:       req.Table,
		//clobColumn:  req.ClobColumn,
	}
}

// chunkTimeseries will attempt to use the query string of a URL to partition the timeseries into "Chunks" of time for
// queying a web API.
func chunkTimeseries(timeseries *Timeseries, rurl url.URL) error {
	// If layout is not set, then default it to be RFC3339
	if timeseries.Layout == "" {
		timeseries.Layout = time.RFC3339
	}

	query := rurl.Query()

	startSlice := query[timeseries.StartName]
	if len(startSlice) != 1 {
		return ErrInvalidStartTimeSize
	}

	start, err := time.Parse(timeseries.Layout, startSlice[0])
	if err != nil {
		return fmt.Errorf("failed to parse start time: %w", err)
	}

	endSlice := query[timeseries.EndName]
	if len(endSlice) != 1 {
		return ErrInvalidEndTimeSize
	}

	end, err := time.Parse(timeseries.Layout, endSlice[0])
	if err != nil {
		return fmt.Errorf("unable to parse end time: %w", err)
	}

	for start.Before(end) {
		next := start.Add(time.Second * time.Duration(timeseries.Period))
		if next.Before(end) {
			timeseries.Chunks = append(timeseries.Chunks, [2]time.Time{start, next})
		} else {
			timeseries.Chunks = append(timeseries.Chunks, [2]time.Time{start, end})
		}

		start = next
	}

	return nil
}

// flattenRequestTimeseries will compress the request information into a "web.FetchConfig" request and a "table" name
// for storage interaction. This function will create a flattened request for each time series in the request. If no
// timeseries are defined, this function will return a single flattened request.
//func flattenRequestTimeseries(req *Request, rurl url.URL, client *web.Client) ([]*flattenedRequest, error) {
//	timeseries := req.Timeseries
//	if timeseries == nil {
//		flatReq := flattenRequest(req, rurl, client)
//
//		return []*flattenedRequest{flatReq}, nil
//	}
//
//	requests := make([]*flattenedRequest, 0, len(timeseries.Chunks))
//
//	// Add the query params to the URL.
//	if req.Query != nil {
//		query := rurl.Query()
//		for key, value := range req.Query {
//			query.Set(key, value)
//		}
//
//		rurl.RawQuery = query.Encode()
//	}
//
//	if err := chunkTimeseries(timeseries, rurl); err != nil {
//		return nil, fmt.Errorf("failed to set time series chunks: %w", err)
//	}
//
//	for _, chunk := range timeseries.Chunks {
//		// copy the request and update it to reflect the partitioned timeseries
//		chunkReq := req
//		chunkReq.Query[timeseries.StartName] = chunk[0].Format(timeseries.Layout)
//		chunkReq.Query[timeseries.EndName] = chunk[1].Format(timeseries.Layout)
//
//		fetchConfig := newFetchConfig(chunkReq, rurl, client)
//
//		requests = append(requests, &flattenedRequest{
//			fetchConfig: fetchConfig,
//			table:       req.Table,
//			clobColumn:  req.ClobColumn,
//		})
//	}
//
//	return requests, nil
//}

// flattenConfigRequests will flatten the requests into a single slice for HTTP requests.
func flattenConfigRequests(ctx context.Context, cfg *Config) ([]*flattenedRequest, error) {
	// If the client is not set on the configuration, then create a new one.
	client := cfg.Client
	if cfg.Client == nil {
		client = http.DefaultClient
	}

	var flattenedRequests []*flattenedRequest

	for _, req := range cfg.Requests {
		//flatReqs, err := flattenRequestTimeseries(req, *cfg.URL, client)
		//if err != nil {
		//	return nil, err
		//}

		//flattenedRequests = append(flattenedRequests, flatReqs...)

		flatReq := newFlattenedRequest(req, client, cfg.RateLimiter)
		flattenedRequests = append(flattenedRequests, flatReq)
	}

	if len(flattenedRequests) == 0 {
		return nil, ErrNoRequests
	}

	return flattenedRequests, nil
}
