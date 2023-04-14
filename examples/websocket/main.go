// Copyright 2023 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alpstable/gidari"
	"golang.org/x/net/websocket"
	"google.golang.org/protobuf/types/known/structpb"
)

type OSListWriter struct{}

func (lw *OSListWriter) Write(ctx context.Context, list *structpb.ListValue) error {
	for _, v := range list.GetValues() {
		m, ok := v.GetKind().(*structpb.Value_StructValue)
		if !ok {
			continue
		}

		fmt.Println(m.StructValue.AsMap())
	}

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	svc, err := gidari.NewService(ctx)
	if err != nil {
		panic(err.Error())
	}

	defer svc.Socket.Close()

	// Create a connection to the Coinbase web socket.
	url := "wss://ws-feed.exchange.coinbase.com"
	conn, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		panic(err.Error())
	}

	// Create a socket with the connection that writes with the
	// OSListWriter.
	socket := gidari.NewSocket(conn,
		gidari.WithSocketWriters(&OSListWriter{}))

	// Add the socket to the service.
	svc.Socket.Connections(socket)

	// Define the subscription message
	msg := []byte(`{
    "type": "subscribe",
    "product_ids": [
        "BTC-USD",
    ],
    "channels": [
    	"ticker"
    ]
}`)

	// Subscribe to the Coinbase web socket.
	if _, err := conn.Write(msg); err != nil {
		panic(err.Error())
	}

	// Start storing the data.
	if err := svc.Socket.Store(ctx); err != nil {
		// If the error is not a deadline exceeded error, then panic.
		if !errors.Is(err, context.DeadlineExceeded) {
			panic(err.Error())
		}
	}

	msg = []byte(`{"type":"unsubscribe","channels":["heartbeat"]}`)

	// Unsubscribe from the Coinbase web socket.
	if _, err := conn.Write(msg); err != nil {
		panic(err.Error())
	}
}
