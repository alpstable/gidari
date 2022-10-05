// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package transport

import (
	"net/url"
	"time"
)

// timeseries is a struct that contains the information needed to query a web API for timeseries data.
type timeseries struct {
	StartName string `yaml:"startName"`
	EndName   string `yaml:"endName"`

	// Period is the size of each chunk in seconds for which we can query the API. Some API will not allow us to
	// query all data within the start and end range.
	Period int32 `yaml:"period"`

	// Layout is the time layout for parsing the "Start" and "End" values into "time.Time". The default is assumed
	// to be RFC3339.
	Layout *string `yaml:"layout"`

	// chunks are the time ranges for which we can query the API. These are broken up into pieces for API requests
	// that only return a limited number of results.
	chunks [][2]time.Time
}

// chunk will attempt to use the query string of a URL to partition the timeseries into "chunks" of time for queying
// a web API.
func (ts *timeseries) chunk(rurl url.URL) error {
	// If layout is not set, then default it to be RFC3339
	if ts.Layout == nil {
		str := time.RFC3339
		ts.Layout = &str
	}

	query := rurl.Query()

	startSlice := query[ts.StartName]
	if len(startSlice) != 1 {
		return MissingTimeseriesFieldError("startName")
	}

	start, err := time.Parse(*ts.Layout, startSlice[0])
	if err != nil {
		return UnableToParseError("startTime")
	}

	endSlice := query[ts.EndName]
	if len(endSlice) != 1 {
		return MissingTimeseriesFieldError("endName")
	}

	end, err := time.Parse(*ts.Layout, endSlice[0])
	if err != nil {
		return UnableToParseError("endTime")
	}

	for start.Before(end) {
		next := start.Add(time.Second * time.Duration(ts.Period))
		if next.Before(end) {
			ts.chunks = append(ts.chunks, [2]time.Time{start, next})
		} else {
			ts.chunks = append(ts.chunks, [2]time.Time{start, end})
		}

		start = next
	}

	return nil
}
