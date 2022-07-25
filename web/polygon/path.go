package polygon

import "path"

// * This is a generated file, do not edit

type rawPath uint8

const (
	_ rawPath = iota
	AggregateBarPath
	UpcomingPath
)

// Get aggregate bars for a stock over a given date range in custom time window sizes.
func getAggregateBarPath(params map[string]string) string {
	return path.Join("/v2", "aggs", "ticker", params["stocks_ticker"], "range", params["multiplier"], params["timespan"], params["from"], params["to"])
}

// Upcoming gets market holidays and their open/close times.
func getUpcomingPath(params map[string]string) string {
	return path.Join("/v1", "marketstatus", "upcoming")
}

// Get takes an rawPath const and rawPath arguments to parse the URL rawPath path.
func (p rawPath) Path(params map[string]string) string {
	return map[rawPath]func(map[string]string) string{
		AggregateBarPath: getAggregateBarPath,
		UpcomingPath:     getUpcomingPath,
	}[p](params)
}

func (p rawPath) Scope() string {
	return map[rawPath]string{}[p]
}
