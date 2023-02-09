package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/alpstable/gidari"
	"google.golang.org/protobuf/types/known/structpb"
)

type ExampleWriter struct {
	lists []*structpb.ListValue
}

func (w *ExampleWriter) Write(ctx context.Context, list *structpb.ListValue) error {
	w.lists = append(w.lists, list)

	return nil
}

func main() {
	ctx := context.Background()

	const api = "https://anapioficeandfire.com/api"

	// First we create a service that can be used to make bulk HTTP
	// requests to the API.
	svc, err := gidari.NewService(ctx)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	// Create some requests and add them to the service.
	bookReq, _ := http.NewRequest(http.MethodGet, api+"/books", nil)
	charReq, _ := http.NewRequest(http.MethodGet, api+"/characters", nil)
	housReq, _ := http.NewRequest(http.MethodGet, api+"/houses", nil)

	w := &ExampleWriter{}

	svc.HTTP.
		Requests(&gidari.HTTPRequest{Request: bookReq, Writer: w}).
		Requests(&gidari.HTTPRequest{Request: charReq, Writer: w}).
		Requests(&gidari.HTTPRequest{Request: housReq, Writer: w})
	// Add a rate limiter to the service, 5 requests per second. This can
	// help avoid "429" errors.
	//svc.HTTP.RateLimiter(rate.NewLimiter(rate.Every(1*time.Second), 5))

	for _, list := range w.lists {
		fmt.Println(list.Values)
	}
}
