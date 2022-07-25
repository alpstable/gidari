package coinbasepro

import (
	"github.com/alpine-hodler/web/pkg/websocket"

	"golang.org/x/sync/errgroup"
)

type TickerChannel chan Ticker

// AsyncTicker is an object that helps maintain state when trying to stream product ticker data.  It starts an
// underlying worker and queues versions of itself to stream from the websocket connection.
type AsyncTicker struct {
	Errors *errgroup.Group

	channel         *TickerChannel
	closed, closing chan struct{}
	conn            websocket.Connector

	// message is  thecoinbase pro websocket subscription that can be used to get a feed of real-time market data.
	message *websocket.Message

	// jobs are a channel of AsyncTicker objects that the user can enqueue by calling Open.  The only possible AT object
	// that can be enqueued is the value that the jobs channel is contained in.  We keep track of which state we're in by
	// using `AsyncTicker.iteration`
	jobs chan *AsyncTicker

	// iteration represents which job we're currently processing.  For instance, if a user called Open 100 times, the
	// first iteration would be 1 and the final iteration would be 100.
	iteration int
}

func (ticker *AsyncTicker) enqueue() { go func() { ticker.jobs <- ticker }() }

func newAsyncTicker(conn websocket.Connector, products ...string) *AsyncTicker {
	ticker := new(AsyncTicker)
	ticker.Errors = new(errgroup.Group)
	ticker.conn = conn

	channels := []websocket.Channel{{Name: "ticker"}}
	msg, _ := websocket.NewProductsMessage(products, channels)
	msg.Subscribe(conn)
	ticker.message = msg

	ticker.closing = make(chan struct{})

	// initialize the jobs channel and start the worker that listens for the jobs
	ticker.jobs = make(chan *AsyncTicker)
	go ticker.worker()

	return ticker
}

// streamAsyncTickerData enables streams of messages to be consumed by the underlying AsyncTicker.channel by using the
// open websocket connection to read rows as json and send them to AsyncTicker.channel.  Once the Close method is
// called, then the ticker channel will close and the row-loop terminated.
func streamAsyncTickerData(ticker *AsyncTicker) error {
	defer func(ticker *AsyncTicker) {
		close(*ticker.channel)
	}(ticker)
	for {
		var row Ticker
		if err := ticker.conn.ReadJSON(&row); err != nil {
			return err
		}
		select {
		case <-ticker.closing:
			return nil
		case *ticker.channel <- row:
		}
	}
}

// worker ranges over the AsyncTicker.jobs channel to asyncronously stream ticker data for each job the queue.
func (ticker *AsyncTicker) worker() error {
	for ticker := range ticker.jobs {
		// Only one ticker can be processed at a time, so we can add a job id to it here and can be assured that there will
		// not be any race conditions.
		ticker.iteration++

		// Make channels for the process
		tmp := make(TickerChannel)
		ticker.channel = &tmp
		ticker.closed = make(chan struct{})

		// Start the stream for async ticker data
		if err := streamAsyncTickerData(ticker); err != nil {
			return err
		}

		// We need to wait until the ticker has been closed, which can only happen when the ticker.closing channel is
		// channeled data.
		<-ticker.closed
	}
	return nil
}

// Open starts the websocket stream, streaming it into the AsyncTicker.channel.  This method is idempotent, so if you
// call it multiple times successively without closing the websocket it will close the ws for you in each successive run
// and re-make the channels to stream over.f the calls need to return without starting the go-routing.
func (ticker *AsyncTicker) Open() *AsyncTicker {
	ticker.enqueue()
	return ticker
}

// Channel returns the ticker channel for streaming
func (ticker *AsyncTicker) Channel() TickerChannel {
	if ch := ticker.channel; ch != nil {
		return *ticker.channel
	}
	return ticker.Channel()
}

// Close unsubscribes the message from the websocket and closes the channel. The Close routine can be called multiple
// times safely.  This method is run asyncronously, relying on the worker to not enqueue jobs where the previous
// iteration has not been closed.  The asynchronous nature of the function allows the user to "Close" the job without
// waiting on the closing channel to be resolved.
func (ticker *AsyncTicker) Close() error {
	if err := ticker.message.Unsubscribe(ticker.conn); err != nil {
		return err
	}
	go func() {
		select {
		case ticker.closing <- struct{}{}:
			close(ticker.closed)
		case <-ticker.closed:
			// In case the close function is called twice for two iterations before the uderlying go routine has resolved.
		}
	}()
	return nil
}
