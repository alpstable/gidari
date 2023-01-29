package gidari

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/alpstable/gidari/proto"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

// HTTPRequest represents a request to be made by the service to the client.
// This object wraps the "net/http" package request object.
type HTTPRequest struct {
	*http.Request

	// Table is an optional field nd the table name to be used for the
	// storage of data from this request. The default value for the table
	// will be the endpoint of the request URL.
	Table string
}

// HTTPResponse represents a response from the client to to the storage
// service.
type HTTPResponse struct{}

// Client is an interface that wraps the "Do" method of the "net/http" package's
// "Client" type.
type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPService struct {
	client Client
	svc    *Service

	// Iterator is a service that provides the functionality to asynchronously
	// iterate over a set of requests, handling them with a custom handler.
	// Each response in the request is achieved by calling the Iterator's
	// "Next" method, returning the "http.Response" object defined by the
	// "net/http" package.
	Iterator *HTTPIteratorService

	rlimiter *rate.Limiter
	requests []*HTTPRequest
}

func NewHTTPService(svc *Service) *HTTPService {
	httpSvc := &HTTPService{svc: svc}
	httpSvc.Iterator = NewHTTPIteratorService(httpSvc)
	httpSvc.client = http.DefaultClient

	return httpSvc
}

// RateLimiter sets the optional rate limiter for the service. A rate limiter
// will limit the request to a set of bursts per period, avoiding 429 errors.
func (svc *HTTPService) RateLimiter(rate *rate.Limiter) *HTTPService {
	svc.rlimiter = rate

	return svc
}

// Client sets the optional client to be used by the service. If no client is
// set, the default "http.DefaultClient" defined by the "net/http" package
// will be used.
func (svc *HTTPService) Client(client Client) *HTTPService {
	svc.client = client

	return svc
}

// Requests sets the option requests to be made by the service to the client.
// If no client has been set for the service, the default "http.DefaultClient"
// defined by the "net/http" package will be used.
func (svc *HTTPService) Requests(reqs ...*HTTPRequest) *HTTPService {
	svc.requests = append(svc.requests, reqs...)

	return svc
}

/*
svc, _ := NewService(context.Background, WithStorage("sqlite3", "test.db"))
if _, err := httpSvc := svc.HTTP.Client(x).Requests(x).Do(); err != nil {
	return err
}
*/

// Do will execute the requests set for the service. If no requests have been
// set, the service will do nothing and return nil.
func (svc *HTTPService) Do(ctx context.Context) (*HTTPResponse, error) {
	reqCount := len(svc.requests)

	// If there are no requests, do nothing.
	if reqCount == 0 {
		return nil, nil
	}

	// Create a channel to send requests to the worker.
	upsertWorkerJobs := make(chan upsertWorkerJob, reqCount)

	// done is a channel that will be closed when the worker is done.
	done := make(chan struct{}, reqCount)

	// Create repositories
	repositories, err := svc.svc.repositories(ctx)
	if err != nil {
		return nil, err
	}

	// Start the upsert worker.
	for i := 1; i <= runtime.NumCPU(); i++ {
		go startUpsertWorker(ctx, upsertWorkerConfig{
			jobs:         upsertWorkerJobs,
			repositories: repositories,
			done:         done,
		})
	}

	for svc.Iterator.Next(ctx) {
		rsp := svc.Iterator.Current.Response

		// If there is no response, then return an error.
		if rsp == nil {
			return nil, fmt.Errorf("no response")
		}

		// Read the response body of the request.
		body, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		// Get the best fit type for decoding the response body. If the
		// best fit is "Unknown", then return an error.
		bestFit := proto.BestFitDecodeType(rsp.Header.Get("Accept"))
		if bestFit == proto.DecodeTypeUnknown {
			return nil, fmt.Errorf("unknown decode type for %q", rsp.Request.URL.String())
		}

		upsertWorkerJobs <- upsertWorkerJob{
			table:    svc.Iterator.Current.Table,
			data:     body,
			dataType: bestFit,
		}
	}

	if err := svc.Iterator.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over requests: %w", err)
	}

	for w := 1; w <= reqCount; w++ {
		<-done
	}

	return nil, nil
}

// Current is a struct that represents the most recent response by calling the
// "Next" method on the HTTPIteratorService.
type Current struct {
	Response *http.Response // HTTP response from the request.
	Table    string         // Name of the table for storage.
}

/*
svc, _ := NewService(context.Background)
for svc.HTTP.Iterator.Next() {
    // Do something with iter.Current
}

// Check for errors
if err := iter.Err(); err != nil {
    // Handle error
}
*/

type HTTPIteratorService struct {
	svc *HTTPService

	Current *Current

	currentChan chan *Current

	// errCh is a channel that holds any errors encountered by the iterator.
	// It can be propagated to the user by the "Err" method.
	errCh chan error

	// err is the error that was encountered by the iterator and is
	// propagated to the user with the "Err" method.
	err error
}

func NewHTTPIteratorService(svc *HTTPService) *HTTPIteratorService {
	iter := &HTTPIteratorService{
		svc:   svc,
		errCh: make(chan error, 1),
	}

	return iter
}

// Err returns any error encountered by the iterator.
func (iter *HTTPIteratorService) Err() error {
	return iter.err
}

type responseWorkerConfig struct {
	coreNum int

	// jobs are a channel of HTTP Response objects sent by the web worker.
	jobs <-chan *http.Response

	// responseChan
	responseChan chan<- *http.Response

	done  chan<- bool
	errCh chan<- error
}

type webWorkerJob struct {
	req      *HTTPRequest
	client   Client
	rlimiter *rate.Limiter
}

type webWorkerConfig struct {
	// id is a unique identifier for the worker. This value MUST be set in
	// order to start a web worker. One and only one web worker
	// configuration MUST have an ID of 1 in order to close the response
	// channel.
	id int

	jobs      chan webWorkerJob
	currentCh chan *Current
	done      chan bool
	errCh     chan error

	hasClosed  bool
	jobCounter int
}

// fetch will make an HTTP request using the underlying client and endpoint.
func fetch(ctx context.Context, job *webWorkerJob) (*http.Response, error) {
	// If the rate limiter is not set, set it with defaults.
	if rlimiter := job.rlimiter; rlimiter != nil {
		if err := job.rlimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter error: %w", err)
		}
	}

	rsp, err := job.client.Do(job.req.Request)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return rsp, nil
}

// startWebWorker will start a worker upto the given specifications of the
// configuration. The worker will listen for jobs defined by the confirugation,
// asynchronous make web requests, and then propagate them onto the response
// channel.
//
// This function should be the only function that sends to the response channel
// (i.e. "rspCh"). Because this function is meant to be used as a worker pool,
// it is important that the response channel is not closed until all workers
// have finished. Therefore, this function will close the response channel ONLY
// when the worker with ID 1 has finished. This works because all workers will
// be blocked from the "defer" method until the "jobs" channel is closed.
//
// If an error is encountered, the worker will push the error onto the error
// channel and then exit. Note that only the  most recent error will be
// propagated to the "errCh" channel, per the rules of "errgroup.Group". Also,
// regardless of errors encountered, the worker will always continue to process
// jobs until the jobs channel is closed.
func startWebWorker(ctx context.Context, cfg *webWorkerConfig) {
	errs, ctx := errgroup.WithContext(ctx)

	for job := range cfg.jobs {
		job := job

		errs.Go(func() error {
			select {
			case <-ctx.Done():
				cfg.jobCounter++

				return wrapErrDeadlineExceeded(ctx.Err())
			default:

				var httpRsp *http.Response
				defer func() {
					// Push the response onto the response
					// channel, if it exists.
					cfg.currentCh <- &Current{
						Response: httpRsp,
						Table:    job.req.Table,
					}

					// Alert the done channel that the job
					// is complete, if it exists.
					cfg.done <- true
				}()

				rsp, err := fetch(ctx, &job)
				if err != nil {
					return err
				}

				httpRsp = rsp

				return nil
			}
		})
	}

	if err := errs.Wait(); err != nil {
		cfg.errCh <- err
	}

	if cfg.id == 1 {
		close(cfg.currentCh)
		close(cfg.done)
		close(cfg.errCh)
	}
}

// startWorkers will start the iterator's web workers and response workers. This
// method can be used to lazy load the underlying buffered channels.
func (iter *HTTPIteratorService) startWorkers(ctx context.Context) {
	reqCount := len(iter.svc.requests)
	iter.currentChan = make(chan *Current, reqCount)

	// webWorkerJobChan is responsible for making HTTP requests and pushing
	// the response body onto the responseWorkerJobChan. This channel is
	// buffered to be equal to the number of requests made.
	webWorkerJobChan := make(chan webWorkerJob, reqCount)
	done := make(chan bool, reqCount)

	// Start the web workers.
	for i := 0; i < runtime.NumCPU(); i++ {
		go startWebWorker(ctx, &webWorkerConfig{
			id:        i + 1,
			jobs:      webWorkerJobChan,
			currentCh: iter.currentChan,
			done:      done,
			errCh:     iter.errCh,
		})
	}

	go func() {
		// Send the flattened requests to the web workers for processing.
		for _, req := range iter.svc.requests {
			webWorkerJobChan <- webWorkerJob{
				req:      req,
				client:   iter.svc.client,
				rlimiter: iter.svc.rlimiter,
			}
		}
	}()

	go func() {
		// Wait for all the web workers to finish.
		for i := 0; i < reqCount; i++ {
			<-done
		}

		close(webWorkerJobChan)
	}()
}

// Next will push the next response as a byte slice onto the Iterator. If there
// are no more responses, the returned boolean will be false. The user is
// responsible for decoding the response.
//
// The HTTP requests used to define the configuration will be fetched
// concurrently once the "Next" method is called for the first time.
func (iter *HTTPIteratorService) Next(ctx context.Context) bool {
	// If the current channel is nil, then we need to start the workers.
	// This will lazy load the web workers and the response workers, each
	// buffered by the number of flattened requests.
	if iter.currentChan == nil {
		iter.startWorkers(ctx)
	}

	for {
		select {
		// If the context has timed out or been canceled, then we return
		// false.
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

			// If the result is "nil," then something is wrong with
			// the response. We will skip this response and continue
			// to encounter the underlying error.
			if result == nil {
				continue
			}

			iter.Current = result

			return true
		}
	}
}
