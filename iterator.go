package gidari

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"runtime"

	"github.com/alpstable/gidari/internal/web"
	"github.com/alpstable/gidari/internal/web/auth"
	"github.com/alpstable/gidari/proto"
)

type responseWorkerJob struct {
	rsp       []byte
	uri       *url.URL
	tableName string
}

type responseJobFn func(context.Context, *responseWorkerJob) (*proto.IteratorResult, error)

type webWorkerJob struct {
	req *flattenedRequest
}

// Iterator holds the request state of a gidari configuration and can be used to iterate over the request results.
type Iterator struct {
	// Current is a byte slice of the most recent request pushed onto the iterator by the "Next" method.
	Current *proto.IteratorResult

	// cfg is the configuration for the iterator.
	cfg *Config

	// currentChan is a channel that holds the end-result of the response worker. The iterator's "Next" method is
	// used to push data from the currentChan onto the Current field.
	//
	// The size of the currentChan is partially non-deterministic. That is, the buffer size should be equal to the
	// number of results in an HTTP JSON response. However, the number of results is not known until the response
	// is received and decoded by the response worker. Therefore, this channel must remain unbuffered and the
	// closure of this channel is left to the response worker.
	currentChan chan *proto.IteratorResult

	// errCh is a channel that holds any errors encountered by the iterator. It can be propagated to the user
	// by the "Err" method.
	errCh chan error

	err error

	// responseJobFn is an optional function that can be used to run custom logic on the response worker job. By
	// default this function is nil and the response worker will attempt to decode the response body into an
	// IteratorResult.
	responseJobFn responseJobFn
}

// NewIterator returns an Iterator object for the given configuration.
func NewIterator(ctx context.Context, cfg *Config) (*Iterator, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	iter := &Iterator{
		cfg:   cfg,
		errCh: make(chan error, 1),
	}

	return iter, nil
}

// Close will close the iterator and release any resources.
func (iter *Iterator) Close(ctx context.Context) {
	//close(iter.currentChan)
}

// Err returns any error encountered by the iterator.
func (iter *Iterator) Err() error {
	return iter.err
}

type responseWorkerConfig struct {
	coreNum                      int
	jobs                         <-chan *responseWorkerJob
	currentIteratorResultJobChan chan<- *proto.IteratorResult
	done                         chan<- bool
	errCh                        chan<- error
	fn                           responseJobFn
}

func startResponseWorker(_ context.Context, cfg responseWorkerConfig) {
	defer func() {
		if cfg.coreNum == 1 {
			close(cfg.currentIteratorResultJobChan)
			close(cfg.done)
		}
	}()

	for job := range cfg.jobs {
		if job == nil {
			cfg.done <- true

			continue
		}

		// TODO: move this block to its own function.
		if cfg.fn != nil {
			result, err := cfg.fn(context.Background(), job)
			if err != nil {
				cfg.errCh <- fmt.Errorf("failed to run response job function: %w", err)
			} else {
				cfg.currentIteratorResultJobChan <- result
			}

			cfg.done <- true

			continue
		}

		results, err := proto.DecodeIteratorResults(job.uri.String(), job.rsp)
		if err != nil {
			cfg.errCh <- fmt.Errorf("error decoding iterator results: %w", err)

			continue
		}

		for _, result := range results {
			cfg.currentIteratorResultJobChan <- result
		}

		cfg.done <- true
	}

}

type webWorkerConfig struct {
	coreNum               int
	jobs                  <-chan webWorkerJob
	responseWorkerJobChan chan<- *responseWorkerJob
	done                  chan<- bool
	errCh                 chan<- error
}

// fetch is a helper function that fetches the given request and returns the response body.
func fetch(ctx context.Context, cfg *web.FetchConfig) ([]byte, error) {
	rsp, err := web.Fetch(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}

	bytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	defer rsp.Body.Close()

	return bytes, nil
}

// sanitizeInvalidJSON will sanitize invalid JSON by convering the source bytes into a string and using the clob column
// on the request to create a valid JSON object. This function will return "false" if there is an error creating the
// new JSON object or if the source is invalid JSON and no clob column is defined for the request.
func sanitizeInvalidJSON(src []byte, cfg webWorkerJob) ([]byte, bool, error) {
	// If the source is valid JSON, then we can return the source bytes.
	if json.Valid(src) {
		return src, true, nil
	}

	// If the source is invalid JSON and there is no clob column defined, then we cannot sanitize the JSON.
	if cfg.req.clobColumn == "" {
		return nil, false, nil
	}

	// If the source is invalid JSON and there is a clob column defined, then we can create a new JSON object
	// using the clob column.
	obj := map[string]string{
		cfg.req.clobColumn: string(src),
	}

	newSrc, err := json.Marshal(obj)
	if err != nil {
		return nil, false, fmt.Errorf("error marshaling new JSON object: %w", err)
	}

	return newSrc, true, nil
}

// startWebWorker will start a worker that will make a request to the given URL and push the response onto the
// the repository channel. If the request fails, the error will be pushed onto the error channel which will
// propagate to the iterator's Next method.
func startWebWorker(ctx context.Context, cfg webWorkerConfig) {
	defer func() {
		if cfg.coreNum == 1 {
			close(cfg.responseWorkerJobChan)
		}
	}()

	for job := range cfg.jobs {
		fetchConfig := job.req.fetchConfig

		rspBytes, err := fetch(ctx, fetchConfig)
		if err != nil {
			cfg.errCh <- fmt.Errorf("error fetching: %w", err)
		}

		// Sanitize the response bytes. If the response bytes are invalid JSON then return nil on the
		// responseWorkerJobChan.
		sanitizedRspBytes, ok, err := sanitizeInvalidJSON(rspBytes, job)
		if err != nil {
			cfg.errCh <- fmt.Errorf("error sanitizing invalid JSON: %w", err)

			continue
		}

		if !ok {
			cfg.responseWorkerJobChan <- nil

			continue
		}

		cfg.responseWorkerJobChan <- &responseWorkerJob{
			rsp:       sanitizedRspBytes,
			uri:       fetchConfig.URL,
			tableName: job.req.table,
		}

		cfg.done <- true
	}
}

// connect will attempt to connect to the web API client. Since there are multiple ways to build a transport given the
// authentication data, this method will exhaust every transport option in the "Authentication" struct.
//
// If the a client is defined on the configuration, then connect will return the user-defined client instead of
// instantiating one gracefully.
func connect(ctx context.Context, cfg *Config) (*web.Client, error) {
	if cfg.Client != nil {
		return &web.Client{
			Client: *cfg.Client,
		}, nil
	}

	if apiKey := cfg.Authentication.APIKey; apiKey != nil {
		client, err := web.NewClient(ctx, auth.NewAPIKey().
			SetURL(cfg.RawURL).
			SetKey(apiKey.Key).
			SetPassphrase(apiKey.Passphrase).
			SetSecret(apiKey.Secret))
		if err != nil {
			return nil, fmt.Errorf("failed to create API key client: %w", err)
		}

		return client, nil
	}

	if apiKey := cfg.Authentication.Auth2; apiKey != nil {
		client, err := web.NewClient(ctx, auth.NewAuth2().SetBearer(apiKey.Bearer).SetURL(cfg.RawURL))
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		}

		return client, nil
	}

	// In the case of no authentication, create a client without an auth transport.
	client, err := web.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

// startWorkers will start the iterator's web workers and response workers. This method can be used to lazy load the
// underlying buffered channels.
func (iter *Iterator) startWorkers(ctx context.Context) {
	iter.currentChan = make(chan *proto.IteratorResult)

	flattenedRequests, err := flattenConfigRequests(ctx, iter.cfg)
	if err != nil {
		panic(err)
	}

	reqCount := len(flattenedRequests)
	jobCount := reqCount * 2

	// reponseWorkerJobChan is responsible for decoding an HTTP response body into a slice of
	// IteratorResults which will be pushed to the currentChan. This channel is buffered to be equal to
	// the number of responses we expect to receive, which should be equal to the number of requests made.
	responseWorkerJobChan := make(chan *responseWorkerJob, reqCount)

	// webWorkerJobChan is responsible for making HTTP requests and pushing the response body onto the
	// responseWorkerJobChan. This channel is buffered to be equal to the number of requests made.
	webWorkerJobChan := make(chan webWorkerJob, reqCount)
	done := make(chan bool, jobCount)

	// Start the response workers.
	for i := 0; i < runtime.NumCPU(); i++ {
		go startResponseWorker(ctx, responseWorkerConfig{
			coreNum:                      i + 1,
			jobs:                         responseWorkerJobChan,
			done:                         done,
			currentIteratorResultJobChan: iter.currentChan,
			errCh:                        iter.errCh,
			fn:                           iter.responseJobFn,
		})
	}

	// Start the web workers.
	for i := 0; i < runtime.NumCPU(); i++ {
		go startWebWorker(ctx, webWorkerConfig{
			coreNum:               i + 1,
			jobs:                  webWorkerJobChan,
			responseWorkerJobChan: responseWorkerJobChan,
			done:                  done,
			errCh:                 iter.errCh,
		})
	}

	go func() {
		// Send the flattened requests to the web workers for processing.
		for _, req := range flattenedRequests {
			webWorkerJobChan <- webWorkerJob{
				req: req,
			}
		}
	}()

	go func() {

		// Wait for all the web workers to finish.
		for i := 0; i < jobCount; i++ {
			<-done
		}

		close(webWorkerJobChan)
	}()
}

// Next will push the next response as a byte slice onto the Iterator. If there are no more responses, the
// returned boolean will be false. The user is responsible for decoding the response.
func (iter *Iterator) Next(ctx context.Context) bool {
	// If the current channel is nil, then we need to start the workers. This will lazy load the web workers and
	// the response workers, each buffered by the number of flattened requests.
	if iter.currentChan == nil {
		iter.startWorkers(ctx)
	}

	for {
		select {
		// If the context has timed out or been canceled, then we return false.
		case <-ctx.Done():
			return false
		case err := <-iter.errCh:
			if err != nil {
				iter.err = err

				return false
			}
		case result, ok := <-iter.currentChan:
			if !ok {
				return false
			}

			iter.Current = result

			return true
		}
	}
}
