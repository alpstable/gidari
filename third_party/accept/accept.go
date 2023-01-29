// Copyright (c) 2013, Ryan Rogers
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package accept

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var errInvalidTypeSubtype = "accept: Invalid type '%s'."

// Accept represents a parsed Accept(-Charset|-Encoding|-Language) header.
type Accept struct {
	Typ, Subtype  string
	QualityFactor float64
	Extensions    map[string]string
}

// AcceptSlice is a slice of Accept.
type AcceptSlice []Accept

// Len implements the Len() method of the Sort interface.
func (a AcceptSlice) Len() int {
	return len(a)
}

// Less implements the Less() method of the Sort interface.  Elements are
// sorted in order of decreasing preference.
func (a AcceptSlice) Less(i, j int) bool {
	// Higher qvalues come first.
	if a[i].QualityFactor > a[j].QualityFactor {
		return true
	} else if a[i].QualityFactor < a[j].QualityFactor {
		return false
	}

	// Specific types come before wildcard types.
	if a[i].Typ != "*" && a[j].Typ == "*" {
		return true
	} else if a[i].Typ == "*" && a[j].Typ != "*" {
		return false
	}

	// Specific subtypes come before wildcard subtypes.
	if a[i].Subtype != "*" && a[j].Subtype == "*" {
		return true
	} else if a[i].Subtype == "*" && a[j].Subtype != "*" {
		return false
	}

	// A lot of extensions comes before not a lot of extensions.
	if len(a[i].Extensions) > len(a[j].Extensions) {
		return true
	}

	return false
}

// Swap implements the Swap() method of the Sort interface.
func (a AcceptSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// parseMediaRange parses the provided media range, and on success returns the
// parsed range params and type/subtype pair.
func parseMediaRange(mediaRange string) (rangeParams, typeSubtype []string, err error) {
	rangeParams = strings.Split(mediaRange, ";")
	typeSubtype = strings.Split(rangeParams[0], "/")

	// typeSubtype should have a length of exactly two.
	if len(typeSubtype) > 2 {
		err = fmt.Errorf(errInvalidTypeSubtype, rangeParams[0])
		return
	} else {
		typeSubtype = append(typeSubtype, "*")
	}

	// Sanitize typeSubtype.
	typeSubtype[0] = strings.TrimSpace(typeSubtype[0])
	typeSubtype[1] = strings.TrimSpace(typeSubtype[1])
	if typeSubtype[0] == "" {
		typeSubtype[0] = "*"
	}
	if typeSubtype[1] == "" {
		typeSubtype[1] = "*"
	}

	return
}

// ParseAcceptHeader parses a HTTP Accept(-Charset|-Encoding|-Language) header and returns
// AcceptSlice, sorted in decreasing order of preference.  If the header lists
// multiple types that have the same level of preference (same specificity of
// type and subtype, same qvalue, and same number of extensions), the type
// that was listed in the header first comes first in the returned value.
//
// See http://www.w3.org/Protocols/rfc2616/rfc2616-sec14 for more information.
func ParseAcceptHeader(header string) AcceptSlice {
	mediaRanges := strings.Split(header, ",")
	accepted := make(AcceptSlice, 0, len(mediaRanges))

	for _, mediaRange := range mediaRanges {
		rangeParams, typeSubtype, err := parseMediaRange(mediaRange)
		if err != nil {
			continue
		}

		accept := Accept{
			Typ:           typeSubtype[0],
			Subtype:       typeSubtype[1],
			QualityFactor: 1.0,
			Extensions:    make(map[string]string),
		}

		// If there is only one rangeParams, we can stop here.
		if len(rangeParams) == 1 {
			accepted = append(accepted, accept)
			continue
		}

		// Validate the rangeParams.
		validParams := true
		for _, v := range rangeParams[1:] {
			nameVal := strings.SplitN(v, "=", 2)
			if len(nameVal) != 2 {
				validParams = false
				break
			}
			nameVal[1] = strings.TrimSpace(nameVal[1])
			if name := strings.TrimSpace(nameVal[0]); name == "q" {
				qval, err := strconv.ParseFloat(nameVal[1], 64)
				if err != nil || qval < 0 {
					validParams = false
					break
				}
				if qval > 1.0 {
					qval = 1.0
				}
				accept.QualityFactor = qval
			} else {
				accept.Extensions[name] = nameVal[1]
			}
		}

		if validParams {
			accepted = append(accepted, accept)
		}
	}

	sort.Sort(accepted)
	return accepted
}
