package twitter

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/alpine-hodler/web/internal"
)

// * This is a generated file, do not edit

// CreateBookmarkOptions are options for API requests.
type CreateBookmarkOptions struct {
	// TweetID is the ID of the Tweet that you would like the user id to Bookmark.
	TweetID *string `bson:"tweet_id" json:"tweet_id" sql:"tweet_id"`
}

// BookmarksOptions are options for API requests.
type BookmarksOptions struct {
	// Expansions enable you to request additional data objects that relate to the originally returned Tweets. Submit a list
	// of desired expansions in a comma-separated list without spaces. The ID that represents the expanded data object will
	// be included directly in the Tweet data object, but the expanded object metadata will be returned within the includes
	// response object, and will also include the ID so that you can match this data object to the original Tweet object.
	Expansions *Expansions `bson:"expansions" json:"expansions" sql:"expansions"`

	// MaxResults are the maximum number of results to be returned per page. This can be a number between 1 and 100. By
	// default, each page will return 100 results.
	MaxResults *uint8 `bson:"max_results" json:"max_results" sql:"max_results"`

	// MediaField enables you to select which specific media fields will deliver in each returned Tweet. Specify the desired
	// fields in a comma-separated list without spaces between commas and fields. The Tweet will only return media fields if
	// the Tweet contains media and if you've also included the expansions=attachments.media_keys query parameter in your
	// request. While the media ID will be located in the Tweet object, you will find this ID and all additional media
	// fields in the includes data object.
	MediaFields *MediaFields `bson:"media.fields" json:"media.fields" sql:"media.fields"`

	// Pagination token is used to request the next page of results if all results weren't returned with the latest request,
	// or to go back to the previous page of results. To return the next page, pass the next_token returned in your previous
	// response. To go back one page, pass the previous_token returned in your previous response.
	PaginationToken *string `bson:"pagination_token" json:"pagination_token" sql:"pagination_token"`

	// PlaceFields enables you to select which specific place fields will deliver in each returned Tweet. Specify the
	// desired fields in a comma-separated list without spaces between commas and fields. The Tweet will only return place
	// fields if the Tweet contains a place and if you've also included the expansions=geo.place_id query parameter in your
	// request. While the place ID will be located in the Tweet object, you will find this ID and all additional place
	// fields in the includes data object.
	PlaceFields *PlaceFields `bson:"place.fields" json:"place.fields" sql:"place.fields"`

	// PollFields enables you to select which specific poll fields will deliver in each returned Tweet. Specify the desired
	// fields in a comma-separated list without spaces between commas and fields. The Tweet will only return poll fields if
	// the Tweet contains a poll and if you've also included the expansions=attachments.poll_ids query parameter in your
	// request. While the poll ID will be located in the Tweet object, you will find this ID and all additional poll fields
	// in the includes data object.
	PollFields *PollFields `bson:"poll.fields" json:"poll.fields" sql:"poll.fields"`

	// TweetFields enables you to select which specific Tweet fields will deliver in each returned Tweet object. Specify the
	// desired fields in a comma-separated list without spaces between commas and fields. You can also pass the
	// expansions=referenced_tweets.id expansion to return the specified fields for both the original Tweet and any included
	// referenced Tweets. The requested Tweet fields will display in both the original Tweet data object, as well as in the
	// referenced Tweet expanded data object that will be located in the includes data object.
	TweetFields *TweetFields `bson:"tweet.fields" json:"tweet.fields" sql:"tweet.fields"`

	// TweetFields enables you to select which specific user fields will deliver in each returned Tweet. Specify the desired
	// fields in a comma-separated list without spaces between commas and fields. While the user ID will be located in the
	// original Tweet object, you will find this ID and all additional user fields in the includes data object.
	UserFields *UserFields `bson:"user.fields" json:"user.fields" sql:"user.fields"`
}

// ComplianceJobsOptions are options for API requests.
type ComplianceJobsOptions struct {
	// Status allows to filter by job status. Only one filter can be specified per request. Default: `all`
	Status *Status `bson:"status" json:"status" sql:"status"`

	// Type allows to filter by job type - either by tweets or user ID. Only one filter (tweets or users) can be specified
	// per request.
	Type ComplianceJob `bson:"type" json:"type" sql:"type"`
}

// AllTweetsOptions are options for API requests.
type AllTweetsOptions struct {
	// MaxResults are the maximum number of search results to be returned by a request. A number between 10 and the system
	// limit (currently 500). By default, a request response will return 10 results.
	MaxResults *int `bson:"max_results" json:"max_results" sql:"max_results"`
}

// TweetsOptions are options for API requests.
type TweetsOptions struct {
	// TODO
	Ids []string `bson:"ids" json:"ids" sql:"ids"`
}

func (opts *AllTweetsOptions) EncodeBody() (buf io.Reader, err error)      { return }
func (opts *BookmarksOptions) EncodeBody() (buf io.Reader, err error)      { return }
func (opts *ComplianceJobsOptions) EncodeBody() (buf io.Reader, err error) { return }
func (opts *TweetsOptions) EncodeBody() (buf io.Reader, err error)         { return }
func (opts *CreateBookmarkOptions) EncodeQuery(req *http.Request)          { return }

// SetMaxResults sets the MaxResults field on AllTweetsOptions. MaxResults are the maximum number of search results to
// be returned by a request. A number between 10 and the system limit (currently 500). By default, a request response
// will return 10 results.
func (opts *AllTweetsOptions) SetMaxResults(MaxResults int) *AllTweetsOptions {
	opts.MaxResults = &MaxResults
	return opts
}

// SetExpansions sets the Expansions field on BookmarksOptions. Expansions enable you to request additional data objects
// that relate to the originally returned Tweets. Submit a list of desired expansions in a comma-separated list without
// spaces. The ID that represents the expanded data object will be included directly in the Tweet data object, but the
// expanded object metadata will be returned within the includes response object, and will also include the ID so that
// you can match this data object to the original Tweet object.
func (opts *BookmarksOptions) SetExpansions(Expansions Expansions) *BookmarksOptions {
	opts.Expansions = &Expansions
	return opts
}

// SetMaxResults sets the MaxResults field on BookmarksOptions. MaxResults are the maximum number of results to be
// returned per page. This can be a number between 1 and 100. By default, each page will return 100 results.
func (opts *BookmarksOptions) SetMaxResults(MaxResults uint8) *BookmarksOptions {
	opts.MaxResults = &MaxResults
	return opts
}

// SetMediaFields sets the MediaFields field on BookmarksOptions. MediaField enables you to select which specific media
// fields will deliver in each returned Tweet. Specify the desired fields in a comma-separated list without spaces
// between commas and fields. The Tweet will only return media fields if the Tweet contains media and if you've also
// included the expansions=attachments.media_keys query parameter in your request. While the media ID will be located in
// the Tweet object, you will find this ID and all additional media fields in the includes data object.
func (opts *BookmarksOptions) SetMediaFields(MediaFields MediaFields) *BookmarksOptions {
	opts.MediaFields = &MediaFields
	return opts
}

// SetPaginationToken sets the PaginationToken field on BookmarksOptions. Pagination token is used to request the next
// page of results if all results weren't returned with the latest request, or to go back to the previous page of
// results. To return the next page, pass the next_token returned in your previous response. To go back one page, pass
// the previous_token returned in your previous response.
func (opts *BookmarksOptions) SetPaginationToken(PaginationToken string) *BookmarksOptions {
	opts.PaginationToken = &PaginationToken
	return opts
}

// SetPlaceFields sets the PlaceFields field on BookmarksOptions. PlaceFields enables you to select which specific place
// fields will deliver in each returned Tweet. Specify the desired fields in a comma-separated list without spaces
// between commas and fields. The Tweet will only return place fields if the Tweet contains a place and if you've also
// included the expansions=geo.place_id query parameter in your request. While the place ID will be located in the Tweet
// object, you will find this ID and all additional place fields in the includes data object.
func (opts *BookmarksOptions) SetPlaceFields(PlaceFields PlaceFields) *BookmarksOptions {
	opts.PlaceFields = &PlaceFields
	return opts
}

// SetPollFields sets the PollFields field on BookmarksOptions. PollFields enables you to select which specific poll
// fields will deliver in each returned Tweet. Specify the desired fields in a comma-separated list without spaces
// between commas and fields. The Tweet will only return poll fields if the Tweet contains a poll and if you've also
// included the expansions=attachments.poll_ids query parameter in your request. While the poll ID will be located in
// the Tweet object, you will find this ID and all additional poll fields in the includes data object.
func (opts *BookmarksOptions) SetPollFields(PollFields PollFields) *BookmarksOptions {
	opts.PollFields = &PollFields
	return opts
}

// SetTweetFields sets the TweetFields field on BookmarksOptions. TweetFields enables you to select which specific Tweet
// fields will deliver in each returned Tweet object. Specify the desired fields in a comma-separated list without
// spaces between commas and fields. You can also pass the expansions=referenced_tweets.id expansion to return the
// specified fields for both the original Tweet and any included referenced Tweets. The requested Tweet fields will
// display in both the original Tweet data object, as well as in the referenced Tweet expanded data object that will be
// located in the includes data object.
func (opts *BookmarksOptions) SetTweetFields(TweetFields TweetFields) *BookmarksOptions {
	opts.TweetFields = &TweetFields
	return opts
}

// SetUserFields sets the UserFields field on BookmarksOptions. TweetFields enables you to select which specific user
// fields will deliver in each returned Tweet. Specify the desired fields in a comma-separated list without spaces
// between commas and fields. While the user ID will be located in the original Tweet object, you will find this ID and
// all additional user fields in the includes data object.
func (opts *BookmarksOptions) SetUserFields(UserFields UserFields) *BookmarksOptions {
	opts.UserFields = &UserFields
	return opts
}

// SetType sets the Type field on ComplianceJobsOptions. Type allows to filter by job type - either by tweets or user
// ID. Only one filter (tweets or users) can be specified per request.
func (opts *ComplianceJobsOptions) SetType(Type ComplianceJob) *ComplianceJobsOptions {
	opts.Type = Type
	return opts
}

// SetStatus sets the Status field on ComplianceJobsOptions. Status allows to filter by job status. Only one filter can
// be specified per request. Default: `all`
func (opts *ComplianceJobsOptions) SetStatus(Status Status) *ComplianceJobsOptions {
	opts.Status = &Status
	return opts
}

// SetTweetID sets the TweetID field on CreateBookmarkOptions. TweetID is the ID of the Tweet that you would like the
// user id to Bookmark.
func (opts *CreateBookmarkOptions) SetTweetID(TweetID string) *CreateBookmarkOptions {
	opts.TweetID = &TweetID
	return opts
}

// SetIds sets the Ids field on TweetsOptions. TODO
func (opts *TweetsOptions) SetIds(Ids []string) *TweetsOptions {
	opts.Ids = Ids
	return opts
}

func (opts *CreateBookmarkOptions) EncodeBody() (buf io.Reader, err error) {
	if opts != nil {
		body := make(map[string]interface{})
		internal.HTTPBodyFragment(body, "tweet_id", opts.TweetID)
		raw, err := json.Marshal(body)
		if err == nil {
			buf = bytes.NewBuffer(raw)
		}
	}
	return
}

func (opts *AllTweetsOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeInt(req, "max_results", opts.MaxResults)
	}
	return
}

func (opts *BookmarksOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeStringer(req, "expansions", opts.Expansions)
		internal.HTTPQueryEncodeUint8(req, "max_results", opts.MaxResults)
		internal.HTTPQueryEncodeStringer(req, "media.fields", opts.MediaFields)
		internal.HTTPQueryEncodeString(req, "pagination_token", opts.PaginationToken)
		internal.HTTPQueryEncodeStringer(req, "place.fields", opts.PlaceFields)
		internal.HTTPQueryEncodeStringer(req, "poll.fields", opts.PollFields)
		internal.HTTPQueryEncodeStringer(req, "tweet.fields", opts.TweetFields)
		internal.HTTPQueryEncodeStringer(req, "user.fields", opts.UserFields)
	}
	return
}

func (opts *ComplianceJobsOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeStringer(req, "type", &opts.Type)
		internal.HTTPQueryEncodeStringer(req, "status", opts.Status)
	}
	return
}

func (opts *TweetsOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeStrings(req, "ids", opts.Ids)
	}
	return
}
