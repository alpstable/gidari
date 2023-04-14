// Copyright 2023 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0

package gidari

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"testing"
	"time"
)

var defaultTestTimeout time.Duration = 5 * time.Second

func assertSocketWrites(t *testing.T, writers []ListWriter, data [][]byte) {
	t.Helper()

	writes := 0

	for _, writer := range writers {
		// Assert that the writer is a mockListWriter.
		mockWriter, ok := writer.(*mockListWriter)
		if !ok {
			t.Fatalf("writer is not a mockListWriter")
		}

		// Count the number of writes.
		writes += mockWriter.count

		// Assert that the data is correct.
		for idx, got := range mockWriter.data {
			want := data[idx]

			var wantJSON, gotJSON interface{}
			if err := json.Unmarshal(want, &wantJSON); err != nil {
				t.Fatalf("failed to unmarshal want: %v", err)
			}

			if err := json.Unmarshal(got, &gotJSON); err != nil {
				t.Fatalf("failed to unmarshal got: %v", err)
			}

			if !reflect.DeepEqual(gotJSON, wantJSON) {
				t.Errorf("writer[%d].data = %s; want %s", idx, got, want)
			}
		}
	}

	// Assert that the number of writes is correct.
	if writes != len(data)*len(writers) {
		t.Errorf("writes = %d; want %d", writes, len(data)*len(writers))
	}
}

// newTestSocket will create a new socket with a mock connection, given writers
// and data.
func newTestSocket(t *testing.T, conn io.ReadWriter, writers []ListWriter, data [][]byte) *Socket {
	t.Helper()

	if conn == nil {
		conn = &mockConn{readData: data}
	}

	return NewSocket(conn, WithSocketWriters(writers...))
}

func newTestSocketB(b *testing.B, conn io.ReadWriter, writers []ListWriter, data [][]byte) *Socket {
	b.Helper()

	if conn == nil {
		conn = &mockConn{readData: data}
	}

	return NewSocket(conn, WithSocketWriters(writers...))
}

func TestSocketStore(t *testing.T) {
	t.Parallel()

	type socketCase struct {
		conn    io.ReadWriter
		writers []ListWriter
		data    [][]byte
		want    [][]byte
	}

	tests := []struct {
		name  string
		cases []socketCase
	}{
		{
			name: "one socket and one writer",
			cases: []socketCase{
				{
					writers: []ListWriter{&mockListWriter{}},
					data: [][]byte{
						[]byte(`[{"x":1}]`),
						[]byte(`[{"x":2}]`),
						[]byte(`[{"x":3}]`),
					},
					want: [][]byte{
						[]byte(`[{"x":1}]`),
						[]byte(`[{"x":2}]`),
						[]byte(`[{"x":3}]`),
					},
				},
			},
		},
		{
			name: "one socket and two writers",
			cases: []socketCase{
				{
					writers: []ListWriter{
						&mockListWriter{},
						&mockListWriter{},
					},
					data: [][]byte{
						[]byte(`[{"x":1}]`),
						[]byte(`[{"x":2}]`),
						[]byte(`[{"x":3}]`),
					},
					want: [][]byte{
						[]byte(`[{"x":1}]`),
						[]byte(`[{"x":2}]`),
						[]byte(`[{"x":3}]`),
					},
				},
			},
		},
		{
			name: "two sockets and one writer",
			cases: []socketCase{
				{
					writers: []ListWriter{&mockListWriter{}},
					data: [][]byte{
						[]byte(`[{"x":1}]`),
						[]byte(`[{"x":2}]`),
						[]byte(`[{"x":3}]`),
					},
					want: [][]byte{
						[]byte(`[{"x":1}]`),
						[]byte(`[{"x":2}]`),
						[]byte(`[{"x":3}]`),
					},
				},
				{
					writers: []ListWriter{&mockListWriter{}},
					data: [][]byte{
						[]byte(`[{"x":4}]`),
						[]byte(`[{"x":5}]`),
						[]byte(`[{"x":6}]`),
					},
					want: [][]byte{
						[]byte(`[{"x":4}]`),
						[]byte(`[{"x":5}]`),
						[]byte(`[{"x":6}]`),
					},
				},
			},
		},
		{
			name: "two sockets and two writers",
			cases: []socketCase{
				{
					writers: []ListWriter{
						&mockListWriter{},
						&mockListWriter{},
					},
					data: [][]byte{
						[]byte(`[{"x":1}]`),
						[]byte(`[{"x":2}]`),
						[]byte(`[{"x":3}]`),
					},
					want: [][]byte{
						[]byte(`[{"x":1}]`),
						[]byte(`[{"x":2}]`),
						[]byte(`[{"x":3}]`),
					},
				},
				{
					writers: []ListWriter{
						&mockListWriter{},
						&mockListWriter{},
					},
					data: [][]byte{
						[]byte(`[{"x":4}]`),
						[]byte(`[{"x":5}]`),
						[]byte(`[{"x":6}]`),
					},
					want: [][]byte{
						[]byte(`[{"x":4}]`),
						[]byte(`[{"x":5}]`),
						[]byte(`[{"x":6}]`),
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			sockets := make([]*Socket, 0, len(test.cases))
			wantMsgs := make(map[*Socket][][]byte, len(test.cases))
			for _, c := range test.cases {
				socket := newTestSocket(t, c.conn, c.writers, c.data)
				defer socket.close()

				sockets = append(sockets, socket)
				wantMsgs[socket] = c.want
			}

			// Construct the service.
			svc := NewSocketService(nil).Connections(sockets...)
			ctx := context.Background()

			// Start the service.
			svcErr := make(chan error, 1)
			go func() {
				defer close(svcErr)

				if err := svc.Store(ctx); err != nil {
					svcErr <- err
				}
			}()

			// Wait for the service to finish.
			select {
			case err := <-svcErr:
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			case <-time.After(defaultTestTimeout):
				t.Fatal("timed out waiting for response")
			}

			for _, socket := range sockets {
				want := wantMsgs[socket]
				assertSocketWrites(t, socket.writers, want)
			}
		})
	}
}

func TestSocketStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		readData [][]byte
		writers  []ListWriter
		want     [][]byte
		ctx      func() (context.Context, context.CancelFunc)
		conn     io.ReadWriter
		err      error
	}{
		{
			name:     "empty",
			readData: [][]byte{},
			writers:  []ListWriter{},
			want:     [][]byte{},
		},
		{
			name: "single job",
			readData: [][]byte{
				[]byte(`[{"x":1}]`),
			},
			writers: []ListWriter{&mockListWriter{}},
			want:    [][]byte{[]byte(`[{"x":1}]`)},
		},
		{
			name: "multiple jobs",
			readData: [][]byte{
				[]byte(`[{"x":1}]`),
				[]byte(`[{"x":2}]`),
				[]byte(`[{"x":3}]`),
			},
			writers: []ListWriter{&mockListWriter{}},
			want: [][]byte{
				[]byte(`[{"x":1}]`),
				[]byte(`[{"x":2}]`),
				[]byte(`[{"x":3}]`),
			},
		},
		{
			name: "multiple jobs with multiple writers",
			readData: [][]byte{
				[]byte(`[{"x":1}]`),
				[]byte(`[{"x":2}]`),
				[]byte(`[{"x":3}]`),
			},
			writers: []ListWriter{
				&mockListWriter{},
				&mockListWriter{},
			},
			want: [][]byte{
				[]byte(`[{"x":1}]`),
				[]byte(`[{"x":2}]`),
				[]byte(`[{"x":3}]`),
			},
		},
		{
			name: "partial messages",
			readData: [][]byte{
				[]byte(`[{"x":1}`),
				[]byte(`,{"x":2}`),
				[]byte(`,{"x":3}]`),
			},
			writers: []ListWriter{&mockListWriter{}},
			want:    [][]byte{[]byte(`[{"x":1},{"x":2},{"x":3}]`)},
		},
		{
			name:     "connection that blocks on read",
			readData: [][]byte{},
			conn:     &mockConnBlocker{},
			ctx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)

				return ctx, cancel
			},
			err: context.DeadlineExceeded,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			socket := newTestSocket(t, test.conn, test.writers, test.readData)
			defer socket.close()

			ctx := context.Background()
			if test.ctx != nil {
				var cancel context.CancelFunc
				ctx, cancel = test.ctx()

				defer cancel()
			}

			errCh := socket.start(ctx)

			select {
			case err, ok := <-errCh:
				if ok {
					if errors.Is(err, test.err) {
						return
					}

					t.Errorf("unexpected error: %v", err)
				}
			case <-time.After(defaultTestTimeout):
				t.Error("timeout waiting for response")
			}

			assertSocketWrites(t, test.writers, test.want)
		})
	}
}

func BenchmarkSocketStart(b *testing.B) {
	readData := [][]byte{}

	// Add 10_000 json messages.
	for i := 0; i < 10_000; i++ {
		readData = append(readData, []byte(`[{"x":1},{"x":2},{"x":3}]`))
	}

	writers := []ListWriter{&mockListWriter{}}

	socket := newTestSocketB(b, nil, writers, readData)
	defer socket.close()

	// Run the benchmark.
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			errCh := socket.start(ctx)

			select {
			case err, ok := <-errCh:
				if ok {
					b.Fatalf("unexpected error: %v", err)
				}
			case <-time.After(defaultTestTimeout):
				b.Fatal("timeout waiting for response")
			}
		}
	})
}
