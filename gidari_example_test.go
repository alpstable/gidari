package gidari_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alpstable/gidari"
	"golang.org/x/time/rate"
)

func ExampleHTTPService_Iterator() {
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

	svc.HTTP.
		Requests(&gidari.HTTPRequest{Request: bookReq}).
		Requests(&gidari.HTTPRequest{Request: charReq}).
		Requests(&gidari.HTTPRequest{Request: housReq})

	// Add a rate limiter to the service, 5 requests per second. This can
	// help avoid "429" errors.
	svc.HTTP.RateLimiter(rate.NewLimiter(rate.Every(1*time.Second), 5))

	// byteSize will keep track of the sum of bytes for each HTTP Response's
	// body.
	var byteSize int
	for svc.HTTP.Iterator.Next(ctx) {
		current := svc.HTTP.Iterator.Current

		rsp := current.Response
		if rsp == nil {
			break
		}

		// Get the byte slice from the response body.
		body, err := io.ReadAll(current.Response.Body)
		if err != nil {
			log.Fatalf("failed to read response body: %v", err)
		}

		// Add the number of bytes to the sum.
		byteSize += len(body)
	}

	// Check to see if an error occurred.
	if err := svc.HTTP.Iterator.Err(); err != nil {
		log.Fatalf("failed to iterate over HTTP responses: %v", err)
	}

	fmt.Println("Total number of bytes:", byteSize)
	// Output:
	// Total number of bytes: 256179
}
