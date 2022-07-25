package polygon_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alpine-hodler/web/pkg/polygon"
	"github.com/alpine-hodler/web/pkg/transport"
	"github.com/alpine-hodler/web/tools"
	"github.com/joho/godotenv"
)

func TestExamples(t *testing.T) {
	defer tools.Quiet()()

	godotenv.Load(".test.env")

	type testCase struct {
		name string
		fn   func(t *testing.T)
	}

	testCases := []testCase{
		{"Client.AggregateBar", func(t *testing.T) { ExampleClient_AggregateBar() }},
		{"Client.Upcoming", func(t *testing.T) { ExampleClient_Upcoming() }},
	}

	for idx, testCase := range testCases {
		t.Run(testCase.name, testCase.fn)

		// Polygon FREE is rate-limited to 5 requests per minute.
		if idx != 0 && idx%5 == 0 {
			time.Sleep(1 * time.Minute)
		}
	}
}

func ExampleClient_AggregateBar() {
	url := "https://api.polygon.io"
	apikey := os.Getenv("POLYGON_API_KEY")

	client, err := polygon.NewClient(context.TODO(), transport.NewAuth2().SetBearer(apikey).SetURL(url))
	if err != nil {
		log.Fatalf("error fetching client: %v", err)
	}

	bar, err := client.AggregateBar("AAPL", "1", polygon.TimespanDay, "2021-07-22", "2021-07-22")
	if err != nil {
		log.Fatalf("error fetching aggregate bar: %v", err)
	}

	fmt.Printf("an aggregate bar: %+v\n", bar.Results[0])
}

func ExampleClient_Upcoming() {
	url := "https://api.polygon.io"
	apikey := os.Getenv("POLYGON_API_KEY")

	client, err := polygon.NewClient(context.TODO(), transport.NewAuth2().SetBearer(apikey).SetURL(url))
	if err != nil {
		log.Fatalf("error fetching client: %v", err)
	}

	upcoming, err := client.Upcoming()
	if err != nil {
		log.Fatalf("error fetching upcoming holidays: %v", err)
	}

	fmt.Printf("the next holiday is %q\n", upcoming[0].Name)
}
