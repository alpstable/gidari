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
	"errors"
	"fmt"
	"io"
)

// Socket is a wrapper around a connection that will read from the connection
// and send the data to the list writers. The socket will not close the
// underlying connection. That is the responsibility of the caller.
type Socket struct {
	conn    io.ReadWriter
	done    chan struct{}
	writers []ListWriter
}

// SocketOption is a function that will configure the socket.
type SocketOption func(*Socket)

// NewSocket will create a new socket with the given connection and options.
func NewSocket(conn io.ReadWriter, opts ...SocketOption) *Socket {
	sockets := &Socket{conn: conn}

	for _, opt := range opts {
		opt(sockets)
	}

	return sockets
}

// WithSocketWriters will set the list writers that the socket will write to.
func WithSocketWriters(writers ...ListWriter) SocketOption {
	return func(sockets *Socket) {
		sockets.writers = writers
	}
}

// close will signal the socket to stop reading from the connection. This will
// not close the underlying connection. That is the responsibility of the
// caller.
func (soc *Socket) close() {
	soc.done <- struct{}{}
}

func readWithContext(ctx context.Context, rw io.Reader) ([]byte, error) {
	const maxMessageSize = 1024

	msg := make(chan []byte)
	errs := make(chan error, 1)

	go func() {
		buf := make([]byte, maxMessageSize)

		n, err := rw.Read(buf)
		if err != nil {
			errs <- err

			return
		}

		msg <- buf[:n]
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context error: %w", ctx.Err())
	case err := <-errs:
		return nil, fmt.Errorf("unable to read message: %w", err)
	case buf := <-msg:
		return buf, nil
	}
}

func (soc *Socket) start(ctx context.Context) <-chan error {
	errs := make(chan error, 1)
	soc.done = make(chan struct{}, 1)

	go func() {
		defer close(errs)

		var buffer []byte

		for {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()

				return
			case <-soc.done:
				return
			default:
			}

			msg, err := readWithContext(ctx, soc.conn)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				errs <- err
			}

			// Append incoming message to buffer
			buffer = append(buffer, msg...)

			// Process incoming messages
			for {
				// Look for a complete message in the buffer
				if isPartialJSON(buffer) {
					// Incomplete message, wait for more
					// data
					break
				}

				job := &listWriterJob{writers: soc.writers}
				job.decFunc = decodeFuncJSONFromBytes(buffer)

				if err := <-writeList(ctx, job); err != nil {
					errs <- err

					return
				}

				// Empty buffer
				buffer = nil

				break
			}
		}
	}()

	return errs
}

// SocketService is a service that will listen to messages from sockets and
// send the data to their respective list writers for socket-to-storage
// operations.
type SocketService struct {
	svc *Service

	done    chan struct{}
	sockets []*Socket
}

// NewSocketService will create a new socket service that can listen to
// messages from multiple sockets and send the data to the list writers.
func NewSocketService(svc *Service) *SocketService {
	ws := &SocketService{svc: svc}

	return ws
}

// Close will close all sockets.
func (svc *SocketService) Close() {
	for _, socket := range svc.sockets {
		socket.close()
	}
}

// Connections will set the sockets that the service will listen to.
func (svc *SocketService) Connections(socs ...*Socket) *SocketService {
	svc.sockets = socs

	return svc
}

// Store will start listening to messages from the sockets and send the data
// to their respective list writers for socket-to-storage operations. This
// method will block until all sockets are closed, an error occurs, the context
// is canceled, or the service is closed.
func (svc *SocketService) Store(ctx context.Context) error {
	svc.done = make(chan struct{}, 1)
	socketErrors := make(chan error, len(svc.sockets))

	// Start the socket workers.
	for _, socket := range svc.sockets {
		go func(socket *Socket) {
			if err := <-socket.start(ctx); err != nil {
				socketErrors <- err

				return
			}

			socketErrors <- nil
		}(socket)
	}

	for i := 0; i < len(svc.sockets); i++ {
		err := <-socketErrors
		if err != nil {
			return err
		}
	}

	return nil
}
