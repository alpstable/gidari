package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/alpstable/gidari"
	"github.com/alpstable/gmongo"
	"golang.org/x/time/rate"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.TODO()

	// Create a MongoDB client using the official MongoDB Go Driver.
	clientOptions := options.Client().ApplyURI("mongodb://mongo1:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	// Plug the client into a Gidari MongoDB Storage adapater.
	mdbStorage, err := gmongo.New(ctx, client)
	if err != nil {
		log.Fatalf("failed to create a MongoDB storage adapter: %v", err)
	}

	// Transport the "earthquake" data from "https://earthquake.usgs.gov" into the MongoDB database.
	if err := gidari.Transport(ctx, earthquakeConfig(mdbStorage)); err != nil {
		log.Fatalf("failed to transport earthquake data: %v", err)
	}

	// Transport the "cryptonator" data from "https://api.cryptonator.com" into the MongoDB database.
	if err := gidari.Transport(ctx, invalidJSONResponseConfig(mdbStorage)); err != nil {
		log.Fatalf("failed to transport cryptonator data: %v", err)
	}

	// Transport the "zippopotam" data from "http://api.zippopotam.us" into the MongoDB database.
	if err := gidari.Transport(ctx, multipleRequestConfig(mdbStorage)); err != nil {
		log.Fatalf("failed to transport zippopotam data: %v", err)
	}
}

// earthquakeConfig is a sample configuration to fetch data from "https://earthquake.usgs.gov".
//
// This example illustrates basic usage of Gidari.
func earthquakeConfig(storage *gmongo.Mongo) *gidari.Config {
	return &gidari.Config{
		URL: func() *url.URL {
			url, _ := url.Parse("https://earthquake.usgs.gov")

			return url
		}(),
		Requests: []*gidari.Request{
			{
				Endpoint:    "/fdsnws/event/1/query",
				Method:      http.MethodGet,
				Table:       "earthquake_events",
				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 5),
				Query: map[string]string{
					"format":    "geojson",
					"starttime": "2020-01-01T00:00:00Z",
					"endtime":   "2020-01-04T19:50:02Z",
				},
			},
		},
		StorageOptions: []gidari.StorageOptions{
			{
				Storage:  storage,
				Database: "test",
			},
		},
	}
}

// invalidJSONResponseConfig is a sample configuration to fetch data from "https://api.cryptonator.com" which does not
// respond with valid JSON.
//
// This example illustrates how to handle invalid JSON responses using a "ClobColumn".
func invalidJSONResponseConfig(storage *gmongo.Mongo) *gidari.Config {
	return &gidari.Config{
		URL: func() *url.URL {
			url, _ := url.Parse("https://api.cryptonator.com")

			return url
		}(),
		Requests: []*gidari.Request{
			{
				Endpoint:    "/api/ticker/btc-usd",
				Table:       "cryptonator_btc_ticker",
				Method:      http.MethodGet,
				ClobColumn:  "data",
				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 1),
			},
		},
		StorageOptions: []gidari.StorageOptions{
			{
				Storage:  storage,
				Database: "test",
			},
		},
	}
}

// multipleRequestConfig is a sample configuration to fetch data from "http://api.zippopotam.us" which has multiple
// requests.
//
// This example illustrates how to handle multiple requests.
func multipleRequestConfig(storage *gmongo.Mongo) *gidari.Config {
	return &gidari.Config{
		URL: func() *url.URL {
			url, _ := url.Parse("http://api.zippopotam.us")

			return url
		}(),
		Requests: []*gidari.Request{
			{
				Endpoint:    "/us/98121",
				Method:      http.MethodGet,
				Table:       "seatle",
				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 1),
			},
			{
				Endpoint:    "/us/90210",
				Method:      http.MethodGet,
				Table:       "beverly_hills",
				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 1),
			},
		},
		StorageOptions: []gidari.StorageOptions{
			{
				Storage:  storage,
				Database: "test",
			},
		},
	}
}
