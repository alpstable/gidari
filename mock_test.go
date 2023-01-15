package gidari

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/alpstable/gidari/proto"
	"golang.org/x/time/rate"
)

type mockConfigOptions struct {
	stgCount    int
	reqCount    int
	rateLimiter *rate.Limiter
}

func newMockConfig(opts mockConfigOptions) *Config {
	cfg := &Config{
		RateLimiter: opts.rateLimiter,
	}

	// Default rate limit the requests to 100 per second.
	if opts.rateLimiter == nil {
		cfg.RateLimiter = rate.NewLimiter(rate.Limit(1*time.Second), 100)
	}

	// Create mock storages.
	cfg.Storage = make([]*Storage, opts.stgCount)
	for i := 0; i < opts.stgCount; i++ {
		cfg.Storage[i] = &Storage{
			Storage:  newMockStorage(),
			Database: fmt.Sprintf("database%d", i),
		}
	}

	// Create mock requests.
	cfg.Requests = make([]*Request, opts.reqCount)
	cfg.Client = newMockHTTPClient()
	for i := 0; i < opts.reqCount; i++ {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://example%d", i), nil)
		cfg.Requests[i] = &Request{
			Request: req,
			Table:   fmt.Sprintf("table%d", i),
		}

		// Add response to mock HTTP client.
		cfg.Client.(*mockHTTPClient).responses[req] = &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"id": 1}`))),
			Request:    req,
		}
	}

	return cfg
}

type mockHTTPClient struct {
	mutex     sync.Mutex
	responses map[*http.Request]*http.Response
}

func newMockHTTPClient() *mockHTTPClient {
	return &mockHTTPClient{
		responses: make(map[*http.Request]*http.Response),
	}
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

func newMockStorage() *mockStorage {
	return &mockStorage{}
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
