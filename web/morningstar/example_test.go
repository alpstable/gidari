package morningstar_test

import (
	"log"
	"os"
	"testing"

	"github.com/alpine-hodler/web/pkg/morningstar"
	"github.com/alpine-hodler/web/tools"
	"github.com/joho/godotenv"
)

func TestExamples(t *testing.T) {
	defer tools.Quiet()()

	godotenv.Load(".test.env")

	// Run Examples
	t.Run("NewBearerToken", func(t *testing.T) { ExampleNewBearerToken() })
}

func ExampleNewBearerToken() {
	url := "https://www.us-api.morningstar.com/token/oauth"
	username := os.Getenv("MORNINGSTAR_USERNAME")
	password := os.Getenv("MORNINGSTAR_PASSWORD")

	token, err := morningstar.NewBearerToken(url, username, password)
	if err != nil {
		log.Fatalf("error posting client credentials: %v", err)
	}

	// Do something with our token
	os.Setenv("MORNINGSTAR_BEARER_TOKEN", token.AccessToken)
}
