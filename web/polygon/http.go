package polygon

import "github.com/alpine-hodler/web/internal"

// * This is a generated file, do not edit

// Get aggregate bars for a stock over a given date range in custom time window sizes.
//
// source: https://polygon.io/docs/stocks/get_v2_aggs_ticker__stocksticker__range__multiplier___timespan___from___to
func (c *Client) AggregateBar(stocksTicker string, multiplier string, timespan Timespan, from string, to string) (m *Bar, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(AggregateBarPath),
		internal.HTTPWithParams(map[string]string{
			"stocks_ticker": stocksTicker,
			"multiplier":    multiplier,
			"timespan":      timespan.String(),
			"from":          from,
			"to":            to}),
		internal.HTTPWithRatelimiter(getRateLimiter(AggregateBarRatelimiter, 99)),
		internal.HTTPWithRequest(req))
}

// Upcoming gets market holidays and their open/close times.
//
// source: https://polygon.io/docs/crypto/get_v1_marketstatus_upcoming
func (c *Client) Upcoming() (m []*Upcoming, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(UpcomingPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(UpcomingRatelimiter, 99)),
		internal.HTTPWithRequest(req))
}
