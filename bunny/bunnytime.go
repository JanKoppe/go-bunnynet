// Copyright (c) 2021 Jan Koppe
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"fmt"
	"strings"
	"time"
)

// timestamps returned by the bunny api are not RFC3339/ISO8601,
// even though the API documentation suggests so. The actual responses
// are missing the local offset part. this special type handles
// (un-)marshalling into the time.Time type.

// this code was copied from https://stackoverflow.com/a/25088283

type BunnyTime struct {
	time.Time
}

// note the missing local offset at the end ("Z" or "(+|-)00:00")
const btLayout = "2006-01-02T15:04:05"

func (bt *BunnyTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		bt.Time = time.Time{}
		return
	}
	// for some reason, sometimes responses are RFC3399, sometimes not.
	// I don't know, who needs standards anyways, amirite?
	if strings.HasSuffix(s, "Z") {
		s = s[:len(s)-1]
	}
	bt.Time, err = time.Parse(btLayout, s)
	return
}

func (bt *BunnyTime) MarshalJSON() ([]byte, error) {
	if bt.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", bt.Time.Format(btLayout))), nil
}

var nilTime = (time.Time{}).UnixNano()

func (bt *BunnyTime) IsSet() bool {
	return bt.UnixNano() != nilTime
}
