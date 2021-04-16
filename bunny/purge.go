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

package bunny

import (
	"fmt"
	"net/url"
)

func (c *Client) PurgePullZoneCache(zoneID int64) error {
	return c.doRequest("POST", fmt.Sprintf("/pullzone/%v/purgeCache", zoneID), "", nil, nil)
}

func (c *Client) PurgeURL(purgeURL string, headerName string, headerValue string) error {
	// why is this a GET, bunny?

	v := url.Values{}
	v.Set("url", purgeURL)

	if headerName != "" && headerValue != "" {
		v.Set("headerName", headerName)
		v.Set("headerValue", headerValue)
	}

	return c.doRequest("GET", "/purge", v.Encode(), nil, nil)
}
