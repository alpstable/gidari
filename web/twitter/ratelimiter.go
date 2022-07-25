package twitter

import (
	"time"

	"golang.org/x/time/rate"
)

// * This is a generated file, do not edit

type ratelimiter uint8

const (
	_ ratelimiter = iota
	AllTweetsRatelimiter
	BookmarksRatelimiter
	ComplianceJobsRatelimiter
	CreateBookmarkRatelimiter
	DeleteBookmarkRatelimiter
	MeRatelimiter
	TweetsRatelimiter
)

var ratelimiters = [uint8(8)]*rate.Limiter{}

func init() {
	ratelimiters[AllTweetsRatelimiter] = nil
	ratelimiters[BookmarksRatelimiter] = nil
	ratelimiters[ComplianceJobsRatelimiter] = nil
	ratelimiters[CreateBookmarkRatelimiter] = nil
	ratelimiters[DeleteBookmarkRatelimiter] = nil
	ratelimiters[MeRatelimiter] = nil
	ratelimiters[TweetsRatelimiter] = nil
}

// getRateLimiter will load the rate limiter for a specific request, lazy loaded.
func getRateLimiter(rl ratelimiter, b int) *rate.Limiter {
	if ratelimiters[rl] == nil {
		ratelimiters[rl] = rate.NewLimiter(rate.Every(1*time.Second), b)
	}
	return ratelimiters[rl]
}
