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
	"net/url"
	"strconv"
)

type Statistics struct {
	TotalBandwidthUsed                     int64
	TotalRequestsServed                    int64
	CacheHitRate                           float32
	BandwidthUsedChart                     map[string]float32
	BandwidthCachedChart                   map[string]float32
	CacheHitRateChart                      map[string]float32
	RequestsServedChart                    map[string]float32
	PullRequestsPulledChart                map[string]float32
	OriginShieldBandwidthUsedChart         map[string]float32
	OriginShieldInternalBandwidthUsedChart map[string]float32
	UserBalanceHistoryChart                map[string]float32
	GeoTrafficDistribution                 map[string]float32
	Error3xxChart                          map[string]float32
	Error4xxChart                          map[string]float32
	Error5xxChart                          map[string]float32
}

func (c *Client) GetStatistics(dateFrom BunnyTime, dateTo BunnyTime, zoneID int64, serverZoneID int64, loadErrors bool, hourly bool) (*Statistics, error) {

	v := url.Values{}

	if !dateFrom.IsZero() {
		v.Set("dateFrom", dateFrom.Format("2006-01-02T15:04:05"))
	}

	if !dateTo.IsZero() {
		v.Set("dateTo", dateTo.Format("2006-01-02T15:04:05"))
	}

	if zoneID > 0 {
		v.Set("pullZone", strconv.FormatInt(zoneID, 10))
	}

	if serverZoneID > 0 {
		v.Set("serverZoneId", strconv.FormatInt(serverZoneID, 10))
	}

	if loadErrors {
		v.Set("loadErrors", "true")
	}

	if hourly {
		v.Set("hourly", "true")
	}

	var stats Statistics
	return &stats, c.doRequest("GET", "/statistics", v.Encode(), nil, &stats)
}
