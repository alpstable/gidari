package twitter

import "path"

// * This is a generated file, do not edit

type rawPath uint8

const (
	_ rawPath = iota
	AllTweetsPath
	BookmarksPath
	ComplianceJobsPath
	CreateBookmarkPath
	DeleteBookmarkPath
	MePath
	TweetsPath
)

// TODO
func getAllTweetsPath(params map[string]string) string {
	return path.Join("/2", "tweets", "search", "all")
}

// Bookmarks allows you to get information about a authenticated user’s 800 most recent bookmarked Tweets. This request
// requires OAuth 2.0 Authorization Code with PKCE.
func getBookmarksPath(params map[string]string) string {
	return path.Join("/2", "users", params["user_id"], "bookmarks")
}

// ComplianceJobs will return a list of recent compliance jobs.
func getComplianceJobsPath(params map[string]string) string {
	return path.Join("/2", "compliance", "jobs")
}

// CreateBookmarks causes the user ID of an authenticated user identified in the path parameter to Bookmark the target
// Tweet provided in the request body. This request requires OAuth 2.0 Authorization Code with PKCE.
func getCreateBookmarkPath(params map[string]string) string {
	return path.Join("/2", "users", params["user_id"], "bookmarks")
}

// DeleteBookmarks are a core feature of the Twitter app that allows you to “save” Tweets and easily access them later.
// With these endpoints, you can retrieve, create, delete or build solutions to manage your Bookmarks via the API. This
// request requires OAuth 2.0 Authorization Code with PKCE
func getDeleteBookmarkPath(params map[string]string) string {
	return path.Join("/2", "users", params["user_id"], "bookmarks", params["tweet_id"])
}

// TODO
func getMePath(params map[string]string) string {
	return path.Join("/2", "me")
}

// TODO
func getTweetsPath(params map[string]string) string {
	return path.Join("/2", "tweets")
}

// Get takes an rawPath const and rawPath arguments to parse the URL rawPath path.
func (p rawPath) Path(params map[string]string) string {
	return map[rawPath]func(map[string]string) string{
		AllTweetsPath:      getAllTweetsPath,
		BookmarksPath:      getBookmarksPath,
		ComplianceJobsPath: getComplianceJobsPath,
		CreateBookmarkPath: getCreateBookmarkPath,
		DeleteBookmarkPath: getDeleteBookmarkPath,
		MePath:             getMePath,
		TweetsPath:         getTweetsPath,
	}[p](params)
}

func (p rawPath) Scope() string {
	return map[rawPath]string{
		BookmarksPath: "tweet.read users.read bookmark.read",
	}[p]
}
