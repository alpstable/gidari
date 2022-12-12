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
package decode

import "testing"

func TestParseAcceptHeader(t *testing.T) {
	type parseTest struct {
		input  string
		output AcceptSlice
	}

	parseTests := []parseTest{
		{ // 0
			// Empty/not sent header signals that everything is accepted.
			input: "",
			output: AcceptSlice{
				{ // 0
					typ:           "*",
					subtype:       "*",
					qualityFactor: 1,
					extensions:    map[string]string{},
				},
			},
		},
		{ // 1
			// Chrome is currently sending this.
			input: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			output: AcceptSlice{
				{ // 0
					typ:           "text",
					subtype:       "html",
					qualityFactor: 1,
					extensions:    map[string]string{},
				},
				{ // 1
					typ:           "application",
					subtype:       "xhtml+xml",
					qualityFactor: 1,
					extensions:    map[string]string{},
				},
				{ // 2
					typ:           "application",
					subtype:       "xml",
					qualityFactor: 0.9,
					extensions:    map[string]string{},
				},
				{ // 3
					typ:           "*",
					subtype:       "*",
					qualityFactor: 0.8,
					extensions:    map[string]string{},
				},
			},
		},
		{ // 2
			// Same as 1, except with crazy whitespacing.
			input: `text  /  html  ,	application	/	xhtml+xml	,
					application
					/
					xml
					;
					q
					=
					0.9
					,  *  /  *  ;  q  =  0.8`,
			output: AcceptSlice{
				{ // 0
					typ:           "text",
					subtype:       "html",
					qualityFactor: 1,
					extensions:    map[string]string{},
				},
				{ // 1
					typ:           "application",
					subtype:       "xhtml+xml",
					qualityFactor: 1,
					extensions:    map[string]string{},
				},
				{ // 2
					typ:           "application",
					subtype:       "xml",
					qualityFactor: 0.9,
					extensions:    map[string]string{},
				},
				{ // 3
					typ:           "*",
					subtype:       "*",
					qualityFactor: 0.8,
					extensions:    map[string]string{},
				},
			},
		},
		{ // 3
			// Same as 1, except with modified/invalid qvals.
			input: "text/html;q=1.05,application/xhtml+xml;q=-1.05,application/xml;q=1.0=0.5,*/*;q=INVALID",
			output: AcceptSlice{
				{ // 0
					typ:           "text",
					subtype:       "html",
					qualityFactor: 1,
					extensions:    map[string]string{},
				},
			},
		},
		{ // 4
			// Complex ordering of preference.
			input: "*/*,*/*;a=1,*/*;a=1;b=1,text/*,text/*;a=1,text/*;a=1;b=1,*/plain,*/plain;a=1,*/plain;a=1;b=1,text/plain,text/plain;a=1,text/plain;a=1;b=1",
			output: AcceptSlice{
				{ // 0
					typ:           "text",
					subtype:       "plain",
					qualityFactor: 1,
					extensions: map[string]string{
						"a": "1",
						"b": "1",
					},
				},
				{ // 1
					typ:           "text",
					subtype:       "plain",
					qualityFactor: 1,
					extensions: map[string]string{
						"a": "1",
					},
				},
				{ // 2
					typ:           "text",
					subtype:       "plain",
					qualityFactor: 1,
					extensions:    map[string]string{},
				},
				{ // 3
					typ:           "text",
					subtype:       "*",
					qualityFactor: 1,
					extensions: map[string]string{
						"a": "1",
						"b": "1",
					},
				},
				{ // 4
					typ:           "text",
					subtype:       "*",
					qualityFactor: 1,
					extensions: map[string]string{
						"a": "1",
					},
				},
				{ // 5
					typ:           "text",
					subtype:       "*",
					qualityFactor: 1,
					extensions:    map[string]string{},
				},
				{ // 6
					typ:           "*",
					subtype:       "plain",
					qualityFactor: 1,
					extensions: map[string]string{
						"a": "1",
						"b": "1",
					},
				},
				{ // 7
					typ:           "*",
					subtype:       "plain",
					qualityFactor: 1,
					extensions: map[string]string{
						"a": "1",
					},
				},
				{ // 8
					typ:           "*",
					subtype:       "plain",
					qualityFactor: 1,
					extensions:    map[string]string{},
				},
				{ // 9
					typ:           "*",
					subtype:       "*",
					qualityFactor: 1,
					extensions: map[string]string{
						"a": "1",
						"b": "1",
					},
				},
				{ // 10
					typ:           "*",
					subtype:       "*",
					qualityFactor: 1,
					extensions: map[string]string{
						"a": "1",
					},
				},
				{ // 11
					typ:           "*",
					subtype:       "*",
					qualityFactor: 1,
					extensions:    map[string]string{},
				},
			},
		},
	}

	var accepted AcceptSlice
	for testPos, test := range parseTests {
		accepted = parseAcceptHeader(test.input)
		if len(accepted) != len(test.output) {
			t.Errorf("Parse (%d): expected %d elements, received %d.", testPos, len(test.output), len(accepted))
			continue
		}
		for i, a := range accepted {
			if a.typ != test.output[i].typ {
				t.Errorf("Parse (%d.%d): expected type '%v', received '%v'.", testPos, i, test.output[i].typ, a.typ)
			}
			if a.subtype != test.output[i].subtype {
				t.Errorf("Parse (%d.%d): expected subtype '%v', received '%v'.", testPos, i, test.output[i].subtype, a.subtype)
			}
			if a.qualityFactor != test.output[i].qualityFactor {
				t.Errorf("Parse (%d.%d): expected qval '%v', received '%v'.", testPos, i, test.output[i].qualityFactor, a.qualityFactor)
			}
			if !mapsAreSimilar(t, a.extensions, test.output[i].extensions) {
				t.Errorf("Parse (%d.%d): expected extensions '%v', received '%v'.", testPos, i, test.output[i].extensions, a.extensions)
			}
		}
	}
}

func mapsAreSimilar(t *testing.T, a, b map[string]string) bool {
	t.Helper()

	if len(a) != len(b) {
		return false
	}

	for aKey, aVal := range a {
		if bVal, exists := b[aKey]; !exists || aVal != bVal {
			return false
		}
	}

	return true
}

func BenchmarkParseAcceptHeader(b *testing.B) {
	b.ReportAllocs()

	for _, tcase := range []struct {
		name   string
		header string
	}{
		{
			name:   "csv weighted",
			header: "*/*,*/*;a=1,*/*;a=1;b=1,text/*,text/*;a=1,text/*;a=1;b=1,*/plain,*/plain;a=1,*/plain;a=1;b=1,text/plain,text/plain;a=1,text/plain;a=1;b=1",
		},
	} {
		tcase := tcase

		b.Run(tcase.name, func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = parseAcceptHeader(tcase.header)
			}
		})
	}
}
