package gidari

import (
	"context"
	"fmt"
	"runtime"

	"github.com/alpstable/gidari/proto"
)

type responseWorkerJob struct {
	testString string
	done       chan bool
}

type webWorkerJob struct {
	testString string
}

// Iterator holds the request state of a gidari configuration and can be used to iterate over the request results.
type Iterator struct {
	cfg *Config
	err error

	// Current is a byte slice of the most recent request pushed onto the iterator by the "Next" method.
	Current *proto.IteratorResult

	// currentChan is a channel that holds the end-result of the response worker. The iterator's "Next" method is
	// used to push data from the currentChan onto the Current field.
	//
	// The size of the currentChan is partially non-deterministic. That is, the buffer size should be equal to the
	// number of results in an HTTP JSON response. However, the number of results is not known until the response
	// is received and decoded by the response worker. Therefore, this channel must remain unbuffered and the
	// closure of this channel is left to the response worker.
	currentChan chan *proto.IteratorResult

	// reponseWorkerJobChan is responsible for decoding an HTTP response body into a slice of
	// IteratorResults which will be pushed to the currentChan. This channel is buffered to be equal to
	// the number of responses we expect to receive, which should be equal to the number of requests made.
	responseWorkerJobChan chan responseWorkerJob

	// webWorkerJobChan is responsible for making HTTP requests and pushing the response body onto the
	// responseWorkerJobChan. This channel is buffered to be equal to the number of requests made.
	webWorkerJobChan chan webWorkerJob
}

// NewIterator returns an Iterator object for the given configuration.
func NewIterator(ctx context.Context, cfg *Config) (*Iterator, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	iter := &Iterator{
		cfg: cfg,
	}

	return iter, nil
}

// Close will close the iterator and release any resources.
func (iter *Iterator) Close(ctx context.Context) {
	close(iter.currentChan)
	close(iter.responseWorkerJobChan)
	close(iter.webWorkerJobChan)
}

// Err returns the last error encountered by the iterator.
func (iter *Iterator) Err() error {
	return iter.err
}

func (iter *Iterator) startResponseWorker(_ context.Context) {
	for job := range iter.responseWorkerJobChan {
		iter.currentChan <- &proto.IteratorResult{
			Data: []byte(job.testString),
		}
		//job.done <- true
	}

	// When the response worker has closed, we need to close the current channel.
	close(iter.currentChan)

	//for job := range cfg.jobs {
	//	if job == nil {
	//		cfg.done <- false

	//		continue
	//	}

	//	url := job.req.URL.String()

	//	results, err := proto.DecodeIteratorResults(url, job.b)
	//	if err != nil {
	//		panic(err)
	//	}

	//	// Make the current channel buffer size equal to the number of results.
	//	//iter.currentChan = make(chan *proto.IteratorResult, len(results))

	//	for _, result := range results {
	//		// If the current value on the iterator is nil, then we need to initialize it. Otherwise,
	//		// we push the data onto the currentChan.
	//		if iter.Current == nil {
	//			fmt.Println("data")
	//			iter.Current = result
	//		} else {
	//			fmt.Println("pre channel")
	//			iter.currentChan <- result
	//			fmt.Println("post channel")
	//		}
	//	}

	//	cfg.done <- true
	//}
}

// startWebWorker will start a worker that will make a request to the given URL and push the response onto the
// the repository channel. If the request fails, the error will be pushed onto the error channel which will
// propagate to the iterator's Next method.
func (iter *Iterator) startWebWorker(_ context.Context) {
	for job := range iter.webWorkerJobChan {
		iter.responseWorkerJobChan <- responseWorkerJob{
			testString: job.testString,
			//done:       make(chan bool),
		}
		//<-job.done
	}

	fmt.Println("closing response worker")

	// When the web worker has been closed, we need to close the response channel.
	close(iter.responseWorkerJobChan)

	//defer close(errChan)

	//for job := range jobs {
	//	fmt.Printf("fetch config: %+v\n", job.fetchConfig)
	//	rsp, err := web.Fetch(ctx, job.fetchConfig)
	//	if err != nil {
	//		errChan <- fmt.Errorf("error fetching data: %v", err)

	//		return
	//	}

	//	bytes, err := io.ReadAll(rsp.Body)
	//	if err != nil {
	//		errChan <- fmt.Errorf("error reading body: %v", err)

	//		return
	//	}

	//	if !json.Valid(bytes) {
	//		if job.flattenedRequest.clobColumn == "" {
	//			job.repoJobs <- nil

	//			continue
	//		}

	//		data := make(map[string]string)
	//		data[job.flattenedRequest.clobColumn] = string(bytes)

	//		bytes, err = json.Marshal(data)
	//		if err != nil {
	//			errChan <- fmt.Errorf("error marshaling data: %v", err)

	//			return
	//		}
	//	}

	//	job.repoJobs <- &repoJob{
	//		b:     bytes,
	//		req:   *rsp.Request,
	//		table: job.table,
	//	}

	//	errChan <- nil
	//}
}

// startWorkers will start the iterator's web workers and response workers. This method can be used to lazy load the
// underlying buffered channels.
func (iter *Iterator) startWorkers(ctx context.Context) {
	// TODO: actually flatten the requests here instead of using a const.
	frlen := len(iter.cfg.Requests)

	iter.currentChan = make(chan *proto.IteratorResult)
	iter.responseWorkerJobChan = make(chan responseWorkerJob, frlen)
	iter.webWorkerJobChan = make(chan webWorkerJob, frlen)

	// Start the response workers.
	for i := 0; i < runtime.NumCPU(); i++ {
		go iter.startResponseWorker(ctx)
	}

	// Start the web workers.
	for i := 0; i < runtime.NumCPU(); i++ {
		go iter.startWebWorker(ctx)
	}

	// Send the flattened requests to the web workers for processing.
	for _, req := range iter.cfg.Requests {
		iter.webWorkerJobChan <- webWorkerJob{
			testString: req.Endpoint,
		}
	}
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
		// If the current channel has a value, then we will attempt to set the current value on the iterator
		// and return true.
		case result, ok := <-iter.currentChan:
			if !ok {
				return false
			}

			iter.Current = result

			return true
		}
	}

	// Check if the current channel is not nil. If it isn't, then we need to drain it by writing to the Current
	// bytes slice.
	//if iter.currentChan != nil {
	//	select {
	//	case iteratorResult := <-iter.currentChan:
	//		iter.Current = iteratorResult

	//		return true
	//	default:
	//		iter.currentChan = nil
	//	}
	//}

	//iter.currentChan = make(chan *proto.IteratorResult)

	//threads := runtime.NumCPU()
	//cfg := iter.cfg

	//// If the done channel is nil we need to create it.
	//if iter.done == nil {
	//	iter.done = make(chan bool, 1)
	//}

	//// If the current channel is nil and the done channel has a value, then we need to return false.
	//select {
	//case <-iter.done:
	//	return false
	//default:
	//}

	//// Create the web worker channels.
	//webWorkerErrChan := make(chan error, 1)
	//webWorkerJobsChan := make(chan *webJob, len(iter.cfg.Requests))

	//// Start the web worker.
	//for webWorkerID := 1; webWorkerID <= threads; webWorkerID++ {
	//	go iter.startWebWorker(ctx, webWorkerJobsChan, webWorkerErrChan)
	//}

	//flattenedRequests, err := flattenConfigRequests(ctx, cfg)
	//if err != nil {
	//	iter.err = err

	//	return false
	//}

	//repoConfig, err := newRepoConfig(ctx, cfg, len(flattenedRequests))
	//if err != nil {
	//	iter.err = err

	//	return false
	//}

	//defer repoConfig.closeRepos()

	//// Start the repository workers.
	//for responseWorkerID := 1; responseWorkerID <= threads; responseWorkerID++ {
	//	go iter.startResponseWorker(ctx, repoConfig)
	//}

	//// Enqueue the worker jobs
	//for _, req := range flattenedRequests {
	//	webWorkerJobsChan <- newWebJob(cfg, req, repoConfig.jobs)
	//}

	//// Wait for all of the data to flush.
	//for a := 1; a <= len(flattenedRequests); a++ {
	//	<-repoConfig.done
	//}

	//fmt.Println("done")

	//// If its the first request, then set the current value directly.
	////iter.Current = []byte("hello")

	////iter.currentChan <- []byte("hello2")
	//iter.done <- true
}
