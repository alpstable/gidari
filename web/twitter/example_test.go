package twitter_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"fmt"

	"github.com/alpine-hodler/web/pkg/transport"
	"github.com/alpine-hodler/web/pkg/twitter"
	"github.com/alpine-hodler/web/tools"
	"github.com/joho/godotenv"
)

const refreshTokenFilename = "refresh_token.json"
const refreshTokenURI = "https://api.twitter.com/2/oauth2/token"

type refreshToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

func TestExamples(t *testing.T) {
	return // ! Twitter is not test-worthy functionality, yet.

	defer tools.Quiet()()

	godotenv.Load(".simple-test.env")
	newrt, err := fetchUserContextAccessToken()
	if err != nil {
		t.Fatalf("Error fetching user context access token: %v", err)
	}
	os.Setenv("TWITTER_OAUTH2_USER_CONTEXT", newrt.AccessToken)

	t.Run("NewClient_basic", func(t *testing.T) { ExampleNewClient_basic() })
	t.Run("NewClient_oauth1", func(t *testing.T) { ExampleNewClient_oauth1() })
	t.Run("NewClient_oauth2", func(t *testing.T) { ExampleNewClient_oauth2() })
	t.Run("CreateBookmark", func(t *testing.T) { ExampleClient_CreateBookmark() })
	t.Run("Bookmarks", func(t *testing.T) { ExampleClient_Bookmarks() })
	t.Run("DeleteBookmark", func(t *testing.T) { ExampleClient_DeleteBookmark() })

}

func ExampleNewClient_basic() {
	// Read credentials from environment variables.
	email := os.Getenv("TWITTER_ENTERPRISE_EMAIL")
	password := os.Getenv("TWITTER_ENTERPRISE_PASSWORD")
	url := os.Getenv("TWITTER_URL")

	// Initialize an Basic client transport
	basic := transport.NewBasic().SetEmail(email).SetPassword(password).SetURL(url)

	// Initialize client using the Auth2 client transport.
	client, _ := twitter.NewClient(context.TODO(), basic)
	fmt.Printf("Twitter client: %+v\n", client)
}

func ExampleNewClient_oauth1() {
	// Read credentials from environment variables.
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")
	url := os.Getenv("TWITTER_URL")

	// Initialize an Auth1 client transport.
	oauth1 := transport.NewAuth1().
		SetAccessToken(accessToken).
		SetAccessTokenSecret(accessSecret).
		SetConsumerKey(consumerKey).
		SetConsumerSecret(consumerSecret).
		SetURL(url)

	// Initialize client using the Auth1 client transport.
	client, err := twitter.NewClient(context.TODO(), oauth1)
	if err != nil {
		log.Fatalf("Error creating new client: %v", err)
	}

	// Fetch some tweets to test the connection.
	tweet, err := client.Tweets(new(twitter.TweetsOptions).SetIds([]string{"1261326399320715264"}))
	if err != nil {
		log.Fatalf("Error fetching MongoDB tweet: %v", err)
	}

	fmt.Printf("A tweet about MongoDB: %+v\n", tweet.Data[0])
}

func ExampleNewClient_oauth2() {
	// Read credentials from environment variables.
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
	url := os.Getenv("TWITTER_URL")

	// Initialize an Auth2 client transport
	oauth2 := transport.NewAuth2().SetBearer(bearerToken).SetURL(url)

	// Initialize client using the Auth2 client transport.
	client, err := twitter.NewClient(context.TODO(), oauth2)
	if err != nil {
		log.Fatalf("Error creating new client: %v", err)
	}

	// Fetch some tweets to test the connection.
	tweet, err := client.Tweets(new(twitter.TweetsOptions).SetIds([]string{"1261326399320715264"}))
	if err != nil {
		log.Fatalf("Error fetching MongoDB tweet: %v", err)
	}

	fmt.Printf("A tweet about MongoDB: %+v\n", tweet.Data[0])
}

func ExampleClient_CreateBookmark() {
	// Read credentials from environment variables.
	bearerToken := os.Getenv("TWITTER_OAUTH2_USER_CONTEXT")
	url := os.Getenv("TWITTER_URL")
	userID := os.Getenv("TWITER_USER_ID")

	// Initialize an Auth1 client transport.
	oauth2 := transport.NewAuth2().SetBearer(bearerToken).SetURL(url)

	// Initialize client using the Auth1 client transport.
	client, err := twitter.NewClient(context.TODO(), oauth2)
	if err != nil {
		log.Fatalf("Error creating new client: %v", err)
	}

	options := new(twitter.CreateBookmarkOptions).SetTweetID("1529619581140099072")

	// Fetch some tweets to test the connection.
	bookmarkWrite, err := client.CreateBookmark(userID, options)
	if err != nil {
		log.Fatalf("Error fetching MongoDB tweet: %v", err)
	}
	fmt.Printf("Bookmark created: %+v\n", bookmarkWrite)
}

func ExampleClient_Bookmarks() {
	// Read credentials from environment variables.
	bearerToken := os.Getenv("TWITTER_OAUTH2_USER_CONTEXT")
	url := os.Getenv("TWITTER_URL")
	userID := os.Getenv("TWITER_USER_ID")

	// Initialize an Auth1 client transport.
	oauth2 := transport.NewAuth2().SetBearer(bearerToken).SetURL(url)

	// Initialize client using the Auth1 client transport.
	client, err := twitter.NewClient(context.TODO(), oauth2)
	if err != nil {
		log.Fatalf("Error creating new client: %v", err)
	}

	options := new(twitter.BookmarksOptions).
		SetMaxResults(uint8(5)).
		SetTweetFields(twitter.TweetFields{
			twitter.TweetFieldCreatedAt,
			twitter.TweetFieldInReplyToUserID,
			twitter.TweetFieldReferencedTweets,
			twitter.TweetFieldAttachments,
			twitter.TweetFieldAuthorID,
			twitter.TweetFieldEntities}).
		SetExpansions(twitter.Expansions{
			twitter.ExpansionAttachmentsMediaKeys})

	// Fetch some tweets to test the connection.
	bookmarks, err := client.Bookmarks(userID, options)
	if err != nil {
		log.Fatalf("Error fetching MongoDB tweet: %v", err)
	}
	fmt.Printf("Bookmarks: %+v\n", bookmarks.Data[0])
}

func ExampleClient_DeleteBookmark() {
	// Read credentials from environment variables.
	bearerToken := os.Getenv("TWITTER_OAUTH2_USER_CONTEXT")
	url := os.Getenv("TWITTER_URL")
	userID := os.Getenv("TWITER_USER_ID")

	// Initialize an Auth1 client transport.
	oauth2 := transport.NewAuth2().SetBearer(bearerToken).SetURL(url)

	// Initialize client using the Auth1 client transport.
	client, err := twitter.NewClient(context.TODO(), oauth2)
	if err != nil {
		log.Fatalf("Error creating new client: %v", err)
	}

	// Fetch some tweets to test the connection.
	bookmarkWrite, err := client.DeleteBookmark(userID, "1529619581140099072")
	if err != nil {
		log.Fatalf("Error fetching MongoDB tweet: %v", err)
	}
	fmt.Printf("Bookmark deleted: %+v\n", bookmarkWrite)
}

// fetchUserContextAccessToken returns an OAuth 2.0 User Context bearer token to test with.
func fetchUserContextAccessToken() (*refreshToken, error) {
	file, err := os.Open(refreshTokenFilename)
	// file, err := os.OpenFile(refreshTokenFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	var old refreshToken
	if err := json.Unmarshal([]byte(bytes), &old); err != nil {
		return nil, err
	}

	clientID := os.Getenv("TWITTER_CLIENT_ID")
	clientSecret := os.Getenv("TWITTER_CLIENT_SECRET")

	client := &http.Client{}
	urlEncodings := fmt.Sprintf(`refresh_token=%s&grant_type=refresh_token&client_id=%s`, old.RefreshToken, clientID)
	var data = strings.NewReader(urlEncodings)
	req, err := http.NewRequest("POST", refreshTokenURI, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var newrt refreshToken
	if err := json.Unmarshal(body, &newrt); err != nil {
		return nil, err
	}
	bytes, err = json.Marshal(newrt)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(refreshTokenFilename, bytes, 0644); err != nil {
		return nil, err
	}
	return &newrt, err
}
