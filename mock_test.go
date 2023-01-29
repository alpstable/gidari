package gidari

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/alpstable/gidari/proto"
	"golang.org/x/time/rate"
)

type mockServiceOptions struct {
	stgCount    int
	reqCount    int
	rateLimiter *rate.Limiter
}

func newMockService(opts mockServiceOptions) *Service {
	svc, err := NewService(context.Background(),
		WithStorage(newMockStorage(opts.stgCount)...))
	if err != nil {
		panic(err)
	}

	reqs := newHTTPRequests(opts.reqCount)

	svc.HTTP.
		Client(newMockHTTPClient(reqs)).
		RateLimiter(opts.rateLimiter).
		Requests(reqs...)

	return svc
}

func newHTTPRequests(volume int) []*HTTPRequest {
	requests := make([]*HTTPRequest, volume)
	for i := 0; i < volume; i++ {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://example%d", i), nil)
		requests[i] = &HTTPRequest{
			Request: req,
			Table:   fmt.Sprintf("table%d", i),
		}
	}

	return requests
}

type mockHTTPClient struct {
	mutex     sync.Mutex
	responses map[*http.Request]*http.Response
}

func newMockHTTPClient(reqs []*HTTPRequest) *mockHTTPClient {
	m := &mockHTTPClient{
		responses: make(map[*http.Request]*http.Response, len(reqs)),
	}

	for i, req := range reqs {
		m.responses[req.Request] = &http.Response{
			StatusCode: http.StatusOK,

			// Make the body JSON {x:1}
			Body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{"x":%d}`, i))),
		}
	}

	return m
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.responses) == 0 {
		return nil, nil
	}

	rsp := m.responses[req]
	delete(m.responses, req)

	return rsp, nil
}

type mockStorage struct {
	closeCount int
	closeMutex sync.Mutex

	listPrimaryKeysCount int
	listPrimaryKeysMutex sync.Mutex

	listTablesCount int
	listTablesMutex sync.Mutex

	isNoSQLCount int
	isNoSQLMutex sync.Mutex

	startTxCount int
	startTxMutex sync.Mutex

	truncateCount int
	truncateMutex sync.Mutex

	typeCount int
	typeMutex sync.Mutex

	upsertCount int
	upsertMutex sync.Mutex

	upsertBinaryCount int
	upsertBinaryMutex sync.Mutex

	pingCount int
	pingMutex sync.Mutex
}

func newMockStorage(volume int) []*Storage {
	stgs := make([]*Storage, volume)
	for i := 0; i < volume; i++ {
		stgs[i] = &Storage{
			Storage: &mockStorage{},
		}
	}

	return stgs
}

func (m *mockStorage) Close() {
	m.closeMutex.Lock()
	defer m.closeMutex.Unlock()

	m.closeCount++
}

func (m *mockStorage) ListPrimaryKeys(ctx context.Context) (*proto.ListPrimaryKeysResponse, error) {
	m.listPrimaryKeysMutex.Lock()
	defer m.listPrimaryKeysMutex.Unlock()

	m.listPrimaryKeysCount++

	return nil, nil
}

func (m *mockStorage) ListTables(ctx context.Context) (*proto.ListTablesResponse, error) {
	m.listTablesMutex.Lock()
	defer m.listTablesMutex.Unlock()

	m.listTablesCount++

	return nil, nil
}

func (m *mockStorage) IsNoSQL() bool {
	m.isNoSQLMutex.Lock()
	defer m.isNoSQLMutex.Unlock()

	m.isNoSQLCount++

	return true
}

func (m *mockStorage) StartTx(context.Context) (*proto.Txn, error) {
	m.startTxMutex.Lock()
	defer m.startTxMutex.Unlock()

	m.startTxCount++

	return nil, nil
}

func (m *mockStorage) Truncate(context.Context, *proto.TruncateRequest) (*proto.TruncateResponse, error) {
	m.truncateMutex.Lock()
	defer m.truncateMutex.Unlock()

	m.truncateCount++

	return nil, nil
}

func (m *mockStorage) Type() uint8 {
	m.typeMutex.Lock()
	defer m.typeMutex.Unlock()

	m.typeCount++

	return 0
}

func (m *mockStorage) Upsert(context.Context, *proto.UpsertRequest) (*proto.UpsertResponse, error) {
	m.upsertMutex.Lock()
	defer m.upsertMutex.Unlock()

	m.upsertCount++

	return nil, nil
}

func (m *mockStorage) UpsertBinary(context.Context, *proto.UpsertBinaryRequest) (*proto.UpsertBinaryResponse, error) {
	m.upsertBinaryMutex.Lock()
	defer m.upsertBinaryMutex.Unlock()

	m.upsertBinaryCount++

	return nil, nil
}

func (m *mockStorage) Ping() error {
	m.pingMutex.Lock()
	defer m.pingMutex.Unlock()

	m.pingCount++

	return nil
}
